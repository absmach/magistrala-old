// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/absmach/magistrala"
	"github.com/absmach/magistrala/internal/postgres"
	"github.com/absmach/magistrala/pkg/clients"
	"github.com/absmach/magistrala/pkg/errors"
	svcerr "github.com/absmach/magistrala/pkg/errors/service"
)

const recoveryDuration = 5 * time.Minute

var (
	errRollbackPolicy     = errors.New("failed to rollback policy")
	errRemoveLocalPolicy  = errors.New("failed to remove from local policy copy")
	errRemovePolicyEngine = errors.New("failed to remove from policy engine")
)

var (
	// ErrFailedToRetrieveMembers failed to retrieve group members.
	ErrFailedToRetrieveMembers = errors.New("failed to retrieve group members")

	// ErrFailedToRetrieveMembership failed to retrieve memberships.
	ErrFailedToRetrieveMembership = errors.New("failed to retrieve memberships")

	// ErrFailedToRetrieveAll failed to retrieve groups.
	ErrFailedToRetrieveAll = errors.New("failed to retrieve all groups")

	// ErrFailedToRetrieveParents failed to retrieve groups.
	ErrFailedToRetrieveParents = errors.New("failed to retrieve all groups")

	// ErrFailedToRetrieveChildren failed to retrieve groups.
	ErrFailedToRetrieveChildren = errors.New("failed to retrieve all groups")

	errIssueUser          = errors.New("failed to issue new login key")
	errIssueTmp           = errors.New("failed to issue new temporary key")
	errRevoke             = errors.New("failed to remove key")
	errRetrieve           = errors.New("failed to retrieve key data")
	errIdentify           = errors.New("failed to validate token")
	errPlatform           = errors.New("invalid platform id")
	errCreateDomainPolicy = errors.New("failed to create domain policy")
	errAddPolicies        = errors.New("failed to add policies")
	errRemovePolicies     = errors.New("failed to remove the policies")
	errInvalidPolicy      = errors.New("failed to validate policy")
)

// Authn specifies an API that must be fullfiled by the domain service
// implementation, and all of its decorators (e.g. logging & metrics).
// Token is a string value of the actual Key and is used to authenticate
// an Auth service request.
type Authn interface {
	// Issue issues a new Key, returning its token value alongside.
	Issue(ctx context.Context, token string, key Key) (Token, error)

	// Revoke removes the Key with the provided id that is
	// issued by the user identified by the provided key.
	Revoke(ctx context.Context, token, id string) error

	// RetrieveKey retrieves data for the Key identified by the provided
	// ID, that is issued by the user identified by the provided key.
	RetrieveKey(ctx context.Context, token, id string) (Key, error)

	// Identify validates token token. If token is valid, content
	// is returned. If token is invalid, or invocation failed for some
	// other reason, non-nil error value is returned in response.
	Identify(ctx context.Context, token string) (Key, error)
}

// Service specifies an API that must be fulfilled by the domain service
// implementation, and all of its decorators (e.g. logging & metrics).
// Token is a string value of the actual Key and is used to authenticate
// an Auth service request.
type Service interface {
	Authn
	Authz
	Domains
}

var _ Service = (*service)(nil)

type service struct {
	keys               KeyRepository
	domains            DomainsRepository
	idProvider         magistrala.IDProvider
	agent              PolicyAgent
	tokenizer          Tokenizer
	loginDuration      time.Duration
	refreshDuration    time.Duration
	invitationDuration time.Duration
}

// New instantiates the auth service implementation.
func New(keys KeyRepository, domains DomainsRepository, idp magistrala.IDProvider, tokenizer Tokenizer, policyAgent PolicyAgent, loginDuration, refreshDuration, invitationDuration time.Duration) Service {
	return &service{
		tokenizer:          tokenizer,
		domains:            domains,
		keys:               keys,
		idProvider:         idp,
		agent:              policyAgent,
		loginDuration:      loginDuration,
		refreshDuration:    refreshDuration,
		invitationDuration: invitationDuration,
	}
}

func (svc service) Issue(ctx context.Context, token string, key Key) (Token, error) {
	key.IssuedAt = time.Now().UTC()
	switch key.Type {
	case APIKey:
		return svc.userKey(ctx, token, key)
	case RefreshKey:
		return svc.refreshKey(ctx, token, key)
	case RecoveryKey:
		return svc.tmpKey(recoveryDuration, key)
	case InvitationKey:
		return svc.invitationKey(ctx, key)
	default:
		return svc.accessKey(ctx, key)
	}
}

func (svc service) Revoke(ctx context.Context, token, id string) error {
	issuerID, _, err := svc.authenticate(token)
	if err != nil {
		return errors.Wrap(errRevoke, err)
	}
	if err := svc.keys.Remove(ctx, issuerID, id); err != nil {
		return errors.Wrap(errRevoke, err)
	}
	return nil
}

func (svc service) RetrieveKey(ctx context.Context, token, id string) (Key, error) {
	issuerID, _, err := svc.authenticate(token)
	if err != nil {
		return Key{}, errors.Wrap(errRetrieve, err)
	}

	return svc.keys.Retrieve(ctx, issuerID, id)
}

func (svc service) Identify(ctx context.Context, token string) (Key, error) {
	key, err := svc.tokenizer.Parse(token)
	if err == ErrAPIKeyExpired {
		err = svc.keys.Remove(ctx, key.Issuer, key.ID)
		return Key{}, errors.Wrap(ErrAPIKeyExpired, err)
	}
	if err != nil {
		return Key{}, errors.Wrap(errIdentify, err)
	}

	switch key.Type {
	case RecoveryKey, AccessKey, InvitationKey:
		return key, nil
	case APIKey:
		_, err := svc.keys.Retrieve(ctx, key.Issuer, key.ID)
		if err != nil {
			return Key{}, svcerr.ErrAuthentication
		}
		return key, nil
	default:
		return Key{}, svcerr.ErrAuthentication
	}
}

func (svc service) Authorize(ctx context.Context, pr PolicyReq) error {
	if err := svc.PolicyValidation(pr); err != nil {
		return errors.Wrap(errInvalidPolicy, err)
	}
	if pr.SubjectKind == TokenKind {
		key, err := svc.Identify(ctx, pr.Subject)
		if err != nil {
			return errors.Wrap(svcerr.ErrAuthentication, err)
		}
		if key.Subject == "" {
			if pr.ObjectType == GroupType || pr.ObjectType == ThingType || pr.ObjectType == DomainType {
				return errors.ErrDomainAuthorization
			}
			return svcerr.ErrAuthentication
		}
		pr.Subject = key.Subject
	}
	if err := svc.agent.CheckPolicy(ctx, pr); err != nil {
		return errors.Wrap(svcerr.ErrAuthorization, err)
	}
	return nil
}

func (svc service) AddPolicy(ctx context.Context, pr PolicyReq) error {
	if err := svc.PolicyValidation(pr); err != nil {
		return errors.Wrap(errInvalidPolicy, err)
	}
	return svc.agent.AddPolicy(ctx, pr)
}

func (svc service) PolicyValidation(pr PolicyReq) error {
	if pr.ObjectType == PlatformType && pr.Object != MagistralaObject {
		return errPlatform
	}
	return nil
}

func (svc service) AddPolicies(ctx context.Context, prs []PolicyReq) error {
	for _, pr := range prs {
		if err := svc.PolicyValidation(pr); err != nil {
			return errors.Wrap(errInvalidPolicy, err)
		}
	}
	return svc.agent.AddPolicies(ctx, prs)
}

func (svc service) DeletePolicy(ctx context.Context, pr PolicyReq) error {
	return svc.agent.DeletePolicy(ctx, pr)
}

func (svc service) DeletePolicies(ctx context.Context, prs []PolicyReq) error {
	for _, pr := range prs {
		if err := svc.PolicyValidation(pr); err != nil {
			return errors.Wrap(errInvalidPolicy, err)
		}
	}
	return svc.agent.DeletePolicies(ctx, prs)
}

func (svc service) ListObjects(ctx context.Context, pr PolicyReq, nextPageToken string, limit int32) (PolicyPage, error) {
	if limit <= 0 {
		limit = 100
	}
	res, npt, err := svc.agent.RetrieveObjects(ctx, pr, nextPageToken, limit)
	if err != nil {
		return PolicyPage{}, errors.Wrap(svcerr.ErrNotFound, err)
	}
	var page PolicyPage
	for _, tuple := range res {
		page.Policies = append(page.Policies, tuple.Object)
	}
	page.NextPageToken = npt
	return page, nil
}

func (svc service) ListAllObjects(ctx context.Context, pr PolicyReq) (PolicyPage, error) {
	res, err := svc.agent.RetrieveAllObjects(ctx, pr)
	if err != nil {
		return PolicyPage{}, errors.Wrap(svcerr.ErrNotFound, err)
	}
	var page PolicyPage
	for _, tuple := range res {
		page.Policies = append(page.Policies, tuple.Object)
	}
	return page, nil
}

func (svc service) CountObjects(ctx context.Context, pr PolicyReq) (int, error) {
	return svc.agent.RetrieveAllObjectsCount(ctx, pr)
}

func (svc service) ListSubjects(ctx context.Context, pr PolicyReq, nextPageToken string, limit int32) (PolicyPage, error) {
	if limit <= 0 {
		limit = 100
	}
	res, npt, err := svc.agent.RetrieveSubjects(ctx, pr, nextPageToken, limit)
	if err != nil {
		return PolicyPage{}, errors.Wrap(svcerr.ErrNotFound, err)
	}
	var page PolicyPage
	for _, tuple := range res {
		page.Policies = append(page.Policies, tuple.Subject)
	}
	page.NextPageToken = npt
	return page, nil
}

func (svc service) ListAllSubjects(ctx context.Context, pr PolicyReq) (PolicyPage, error) {
	res, err := svc.agent.RetrieveAllSubjects(ctx, pr)
	if err != nil {
		return PolicyPage{}, errors.Wrap(svcerr.ErrNotFound, err)
	}
	var page PolicyPage
	for _, tuple := range res {
		page.Policies = append(page.Policies, tuple.Subject)
	}
	return page, nil
}

func (svc service) CountSubjects(ctx context.Context, pr PolicyReq) (int, error) {
	return svc.agent.RetrieveAllSubjectsCount(ctx, pr)
}

func (svc service) ListPermissions(ctx context.Context, pr PolicyReq, filterPermisions []string) (Permissions, error) {
	return svc.agent.RetrievePermissions(ctx, pr, filterPermisions)
}

func (svc service) tmpKey(duration time.Duration, key Key) (Token, error) {
	key.ExpiresAt = time.Now().Add(duration)
	value, err := svc.tokenizer.Issue(key)
	if err != nil {
		return Token{}, errors.Wrap(errIssueTmp, err)
	}

	return Token{AccessToken: value}, nil
}

func (svc service) accessKey(ctx context.Context, key Key) (Token, error) {
	var err error
	key.Type = AccessKey
	key.ExpiresAt = time.Now().Add(svc.loginDuration)

	key.Subject, err = svc.checkUserDomain(ctx, key)
	if err != nil {
		return Token{}, errors.Wrap(svcerr.ErrAuthorization, err)
	}

	access, err := svc.tokenizer.Issue(key)
	if err != nil {
		return Token{}, errors.Wrap(errIssueTmp, err)
	}
	key.ExpiresAt = time.Now().Add(svc.refreshDuration)
	key.Type = RefreshKey
	refresh, err := svc.tokenizer.Issue(key)
	if err != nil {
		return Token{}, errors.Wrap(errIssueTmp, err)
	}

	return Token{AccessToken: access, RefreshToken: refresh}, nil
}

func (svc service) invitationKey(ctx context.Context, key Key) (Token, error) {
	var err error
	key.Type = InvitationKey
	key.ExpiresAt = time.Now().Add(svc.invitationDuration)

	key.Subject, err = svc.checkUserDomain(ctx, key)
	if err != nil {
		return Token{}, err
	}

	access, err := svc.tokenizer.Issue(key)
	if err != nil {
		return Token{}, errors.Wrap(errIssueTmp, err)
	}

	return Token{AccessToken: access}, nil
}

func (svc service) refreshKey(ctx context.Context, token string, key Key) (Token, error) {
	k, err := svc.tokenizer.Parse(token)
	if err != nil {
		return Token{}, errors.Wrap(errRetrieve, err)
	}
	if k.Type != RefreshKey {
		return Token{}, errIssueUser
	}
	key.ID = k.ID
	if key.Domain == "" {
		key.Domain = k.Domain
	}
	key.User = k.User
	key.Type = AccessKey

	key.Subject, err = svc.checkUserDomain(ctx, key)
	if err != nil {
		return Token{}, errors.Wrap(svcerr.ErrAuthorization, err)
	}

	key.ExpiresAt = time.Now().Add(svc.loginDuration)
	access, err := svc.tokenizer.Issue(key)
	if err != nil {
		return Token{}, errors.Wrap(errIssueTmp, err)
	}
	key.ExpiresAt = time.Now().Add(svc.refreshDuration)
	key.Type = RefreshKey
	refresh, err := svc.tokenizer.Issue(key)
	if err != nil {
		return Token{}, errors.Wrap(errIssueTmp, err)
	}

	return Token{AccessToken: access, RefreshToken: refresh}, nil
}

func (svc service) checkUserDomain(ctx context.Context, key Key) (subject string, err error) {
	if key.Domain != "" {
		// Check user is platform admin.
		if err = svc.Authorize(ctx, PolicyReq{
			Subject:     key.User,
			SubjectType: UserType,
			Permission:  AdminPermission,
			Object:      MagistralaObject,
			ObjectType:  PlatformType,
		}); err == nil {
			return key.User, nil
		}
		// Check user is domain member.
		domainUserSubject := EncodeDomainUserID(key.Domain, key.User)
		if err = svc.Authorize(ctx, PolicyReq{
			Subject:     domainUserSubject,
			SubjectType: UserType,
			Permission:  MembershipPermission,
			Object:      key.Domain,
			ObjectType:  DomainType,
		}); err != nil {
			return "", err
		}
		return domainUserSubject, nil
	}
	return "", nil
}

func (svc service) userKey(ctx context.Context, token string, key Key) (Token, error) {
	id, sub, err := svc.authenticate(token)
	if err != nil {
		return Token{}, errors.Wrap(errIssueUser, err)
	}

	key.Issuer = id
	if key.Subject == "" {
		key.Subject = sub
	}

	keyID, err := svc.idProvider.ID()
	if err != nil {
		return Token{}, errors.Wrap(errIssueUser, err)
	}
	key.ID = keyID

	if _, err := svc.keys.Save(ctx, key); err != nil {
		return Token{}, errors.Wrap(errIssueUser, err)
	}

	tkn, err := svc.tokenizer.Issue(key)
	if err != nil {
		return Token{}, errors.Wrap(errIssueUser, err)
	}

	return Token{AccessToken: tkn}, nil
}

func (svc service) authenticate(token string) (string, string, error) {
	key, err := svc.tokenizer.Parse(token)
	if err != nil {
		return "", "", errors.Wrap(svcerr.ErrAuthentication, err)
	}
	// Only login key token is valid for login.
	if key.Type != AccessKey || key.Issuer == "" {
		return "", "", svcerr.ErrAuthentication
	}

	return key.Issuer, key.Subject, nil
}

// Switch the relative permission for the relation.
func SwitchToPermission(relation string) string {
	switch relation {
	case AdministratorRelation:
		return AdminPermission
	case EditorRelation:
		return EditPermission
	case ViewerRelation:
		return ViewPermission
	case MemberRelation:
		return MembershipPermission
	default:
		return relation
	}
}

func (svc service) CreateDomain(ctx context.Context, token string, d Domain) (do Domain, err error) {
	key, err := svc.Identify(ctx, token)
	if err != nil {
		return Domain{}, errors.Wrap(svcerr.ErrAuthentication, err)
	}
	d.CreatedBy = key.User

	domainID, err := svc.idProvider.ID()
	if err != nil {
		return Domain{}, errors.Wrap(svcerr.ErrUniqueID, err)
	}
	d.ID = domainID

	if d.Status != clients.DisabledStatus && d.Status != clients.EnabledStatus {
		return Domain{}, svcerr.ErrInvalidStatus
	}

	d.CreatedAt = time.Now()

	if err := svc.createDomainPolicy(ctx, key.User, domainID, AdministratorRelation); err != nil {
		return Domain{}, errors.Wrap(errCreateDomainPolicy, err)
	}
	defer func() {
		if err != nil {
			if errRollBack := svc.createDomainPolicyRollback(ctx, key.User, domainID, AdministratorRelation); errRollBack != nil {
				err = errors.Wrap(err, errors.Wrap(errRollbackPolicy, errRollBack))
			}
		}
	}()

	return svc.domains.Save(ctx, d)
}

func (svc service) RetrieveDomain(ctx context.Context, token string, id string) (Domain, error) {
	if err := svc.Authorize(ctx, PolicyReq{
		Subject:     token,
		SubjectType: UserType,
		SubjectKind: TokenKind,
		Object:      id,
		ObjectType:  DomainType,
		Permission:  ViewPermission,
	}); err != nil {
		return Domain{}, errors.Wrap(svcerr.ErrAuthorization, err)
	}

	return svc.domains.RetrieveByID(ctx, id)
}

func (svc service) UpdateDomain(ctx context.Context, token string, id string, d DomainReq) (Domain, error) {
	key, err := svc.Identify(ctx, token)
	if err != nil {
		return Domain{}, errors.Wrap(svcerr.ErrAuthentication, err)
	}
	if err := svc.Authorize(ctx, PolicyReq{
		Subject:     key.Subject,
		SubjectType: UserType,
		SubjectKind: UsersKind,
		Object:      id,
		ObjectType:  DomainType,
		Permission:  EditPermission,
	}); err != nil {
		return Domain{}, errors.Wrap(svcerr.ErrAuthorization, err)
	}
	return svc.domains.Update(ctx, id, key.User, d)
}

func (svc service) ChangeDomainStatus(ctx context.Context, token string, id string, d DomainReq) (Domain, error) {
	key, err := svc.Identify(ctx, token)
	if err != nil {
		return Domain{}, errors.Wrap(svcerr.ErrAuthentication, err)
	}
	if err := svc.Authorize(ctx, PolicyReq{
		Subject:     key.Subject,
		SubjectType: UserType,
		SubjectKind: UsersKind,
		Object:      id,
		ObjectType:  DomainType,
		Permission:  AdminPermission,
	}); err != nil {
		return Domain{}, errors.Wrap(svcerr.ErrAuthorization, err)
	}
	return svc.domains.Update(ctx, id, key.User, d)
}

func (svc service) ListDomains(ctx context.Context, token string, p Page) (DomainsPage, error) {
	key, err := svc.Identify(ctx, token)
	if err != nil {
		return DomainsPage{}, errors.Wrap(svcerr.ErrAuthentication, err)
	}
	p.SubjectID = key.User
	if err := svc.Authorize(ctx, PolicyReq{
		Subject:     key.User,
		SubjectType: UserType,
		Permission:  AdminPermission,
		ObjectType:  PlatformType,
		Object:      MagistralaObject,
	}); err == nil {
		p.SubjectID = ""
	}
	dp, err := svc.domains.ListDomains(ctx, p)
	if err != nil {
		return DomainsPage{}, postgres.HandleError(svcerr.ErrViewEntity, err)
	}
	if p.SubjectID == "" {
		for i := range dp.Domains {
			dp.Domains[i].Permission = AdministratorRelation
		}
	}
	return dp, nil
}

func (svc service) AssignUsers(ctx context.Context, token string, id string, userIds []string, relation string) error {
	if err := svc.Authorize(ctx, PolicyReq{
		Subject:     token,
		SubjectType: UserType,
		SubjectKind: TokenKind,
		Object:      id,
		ObjectType:  DomainType,
		Permission:  SharePermission,
	}); err != nil {
		return err
	}

	if err := svc.Authorize(ctx, PolicyReq{
		Subject:     token,
		SubjectType: UserType,
		SubjectKind: TokenKind,
		Object:      id,
		ObjectType:  DomainType,
		Permission:  SwitchToPermission(relation),
	}); err != nil {
		return err
	}

	for _, userID := range userIds {
		if err := svc.Authorize(ctx, PolicyReq{
			Subject:     userID,
			SubjectType: UserType,
			Permission:  MembershipPermission,
			Object:      MagistralaObject,
			ObjectType:  PlatformType,
		}); err != nil {
			return errors.Wrap(svcerr.ErrMalformedEntity, fmt.Errorf("invalid user id : %s ", userID))
		}
	}

	return svc.addDomainPolicies(ctx, id, relation, userIds...)
}

func (svc service) UnassignUsers(ctx context.Context, token string, id string, userIds []string, relation string) error {
	if err := svc.Authorize(ctx, PolicyReq{
		Subject:     token,
		SubjectType: UserType,
		SubjectKind: TokenKind,
		Object:      id,
		ObjectType:  DomainType,
		Permission:  SharePermission,
	}); err != nil {
		return err
	}

	if err := svc.Authorize(ctx, PolicyReq{
		Subject:     token,
		SubjectType: UserType,
		SubjectKind: TokenKind,
		Object:      id,
		ObjectType:  DomainType,
		Permission:  SwitchToPermission(relation),
	}); err != nil {
		return err
	}

	return svc.removeDomainPolicies(ctx, id, relation, userIds...)
}

// IMPROVEMENT NOTE: Take decision: Only Patform admin or both Patform and domain admins can see others users domain.
func (svc service) ListUserDomains(ctx context.Context, token string, userID string, p Page) (DomainsPage, error) {
	res, err := svc.Identify(ctx, token)
	if err != nil {
		return DomainsPage{}, errors.Wrap(svcerr.ErrAuthentication, err)
	}
	if err := svc.Authorize(ctx, PolicyReq{
		Subject:     res.User,
		SubjectType: UserType,
		Permission:  AdminPermission,
		Object:      MagistralaObject,
		ObjectType:  PlatformType,
	}); err != nil {
		return DomainsPage{}, errors.Wrap(svcerr.ErrAuthorization, err)
	}
	if userID != "" && res.User != userID {
		p.SubjectID = userID
	} else {
		p.SubjectID = res.User
	}
	return svc.domains.ListDomains(ctx, p)
}

func (svc service) addDomainPolicies(ctx context.Context, domainID, relation string, userIDs ...string) (err error) {
	var prs []PolicyReq
	var pcs []Policy

	for _, userID := range userIDs {
		prs = append(prs, PolicyReq{
			Subject:     EncodeDomainUserID(domainID, userID),
			SubjectType: UserType,
			SubjectKind: UsersKind,
			Relation:    relation,
			Object:      domainID,
			ObjectType:  DomainType,
		})
		pcs = append(pcs, Policy{
			SubjectType: UserType,
			SubjectID:   userID,
			Relation:    relation,
			ObjectType:  DomainType,
			ObjectID:    domainID,
		})
	}
	if err := svc.agent.AddPolicies(ctx, prs); err != nil {
		return errors.Wrap(errAddPolicies, err)
	}
	defer func() {
		if err != nil {
			if errDel := svc.agent.DeletePolicies(ctx, prs); errDel != nil {
				err = errors.Wrap(err, errors.Wrap(errRollbackPolicy, errDel))
			}
		}
	}()
	return svc.domains.SavePolicies(ctx, pcs...)
}

func (svc service) createDomainPolicy(ctx context.Context, userID, domainID, relation string) (err error) {
	prs := []PolicyReq{
		{
			Subject:     EncodeDomainUserID(domainID, userID),
			SubjectType: UserType,
			SubjectKind: UsersKind,
			Relation:    relation,
			Object:      domainID,
			ObjectType:  DomainType,
		},
		{
			Subject:     MagistralaObject,
			SubjectType: PlatformType,
			Relation:    PlatformRelation,
			Object:      domainID,
			ObjectType:  DomainType,
		},
	}
	if err := svc.agent.AddPolicies(ctx, prs); err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if errDel := svc.agent.DeletePolicies(ctx, prs); errDel != nil {
				err = errors.Wrap(err, errors.Wrap(errRollbackPolicy, errDel))
			}
		}
	}()
	return svc.domains.SavePolicies(ctx, Policy{
		SubjectType: UserType,
		SubjectID:   userID,
		Relation:    relation,
		ObjectType:  DomainType,
		ObjectID:    domainID,
	})
}

func (svc service) createDomainPolicyRollback(ctx context.Context, userID, domainID, relation string) error {
	var err error
	prs := []PolicyReq{
		{
			Subject:     EncodeDomainUserID(domainID, userID),
			SubjectType: UserType,
			SubjectKind: UsersKind,
			Relation:    relation,
			Object:      domainID,
			ObjectType:  DomainType,
		},
		{
			Subject:     MagistralaObject,
			SubjectType: PlatformType,
			Relation:    PlatformRelation,
			Object:      domainID,
			ObjectType:  DomainType,
		},
	}
	if errPolicy := svc.agent.DeletePolicies(ctx, prs); errPolicy != nil {
		err = errors.Wrap(errRemovePolicyEngine, errPolicy)
	}
	errPolicyCopy := svc.domains.DeletePolicies(ctx, Policy{
		SubjectType: UserType,
		SubjectID:   userID,
		Relation:    relation,
		ObjectType:  DomainType,
		ObjectID:    domainID,
	})
	if errPolicyCopy != nil {
		err = errors.Wrap(err, errors.Wrap(errRemoveLocalPolicy, errPolicyCopy))
	}
	return err
}

func (svc service) removeDomainPolicies(ctx context.Context, domainID, relation string, userIDs ...string) (err error) {
	var prs []PolicyReq
	var pcs []Policy

	for _, userID := range userIDs {
		prs = append(prs, PolicyReq{
			Subject:     EncodeDomainUserID(domainID, userID),
			SubjectType: UserType,
			SubjectKind: UsersKind,
			Relation:    relation,
			Object:      domainID,
			ObjectType:  DomainType,
		})
		pcs = append(pcs, Policy{
			SubjectType: UserType,
			SubjectID:   userID,
			Relation:    relation,
			ObjectType:  DomainType,
			ObjectID:    domainID,
		})
	}
	if err := svc.agent.DeletePolicies(ctx, prs); err != nil {
		return errors.Wrap(errRemovePolicies, err)
	}

	return svc.domains.DeletePolicies(ctx, pcs...)
}

func EncodeDomainUserID(domainID string, userID string) string {
	if domainID == "" || userID == "" {
		return ""
	}
	return domainID + "_" + userID
}

func DecodeDomainUserID(domainUserID string) (string, string) {
	if domainUserID == "" {
		return domainUserID, domainUserID
	}
	duid := strings.Split(domainUserID, "_")

	switch {
	case len(duid) == 2:
		return duid[0], duid[1]
	case len(duid) == 1:
		return duid[0], ""
	case len(duid) <= 0 || len(duid) > 2:
		fallthrough
	default:
		return "", ""
	}
}
