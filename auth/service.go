// Copyright (c) Magistrala
// SPDX-License-Identifier: Apache-2.0

package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/absmach/magistrala"
	"github.com/absmach/magistrala/internal/apiutil"
	"github.com/absmach/magistrala/pkg/clients"
	"github.com/absmach/magistrala/pkg/errors"
)

const (
	recoveryDuration   = 5 * time.Minute
	invitationDuration = 24 * time.Hour

	refreshToken = "refresh_token"
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

	errIssueUser = errors.New("failed to issue new login key")
	errIssueTmp  = errors.New("failed to issue new temporary key")
	errRevoke    = errors.New("failed to remove key")
	errRetrieve  = errors.New("failed to retrieve key data")
	errIdentify  = errors.New("failed to validate token")
	errPlatform  = errors.New("invalid platform id")
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
	keys            KeyRepository
	domains         DomainsRepository
	idProvider      magistrala.IDProvider
	agent           PolicyAgent
	tokenizer       Tokenizer
	loginDuration   time.Duration
	refreshDuration time.Duration
}

// New instantiates the auth service implementation.
func New(keys KeyRepository, domains DomainsRepository, idp magistrala.IDProvider, tokenizer Tokenizer, policyAgent PolicyAgent, loginDuration, refreshDuration time.Duration) Service {
	return &service{
		tokenizer:       tokenizer,
		domains:         domains,
		keys:            keys,
		idProvider:      idp,
		agent:           policyAgent,
		loginDuration:   loginDuration,
		refreshDuration: refreshDuration,
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
		return svc.tmpKey(invitationDuration, key)
	default:
		return svc.accessKey(key)
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
			return Key{}, errors.ErrAuthentication
		}
		return key, nil
	default:
		return Key{}, errors.ErrAuthentication
	}
}

func (svc service) Authorize(ctx context.Context, pr PolicyReq) error {
	if err := svc.PolicyValidation(pr); err != nil {
		return err
	}
	if pr.SubjectKind == TokenKind {
		key, err := svc.Identify(ctx, pr.Subject)
		if err != nil {
			return err
		}
		if key.Subject == "" {
			return errors.ErrAuthorization
		}
		pr.Subject = key.Subject
	}
	if err := svc.agent.CheckPolicy(ctx, pr); err != nil {
		return errors.Wrap(errors.ErrAuthorization, err)
	}
	return nil
}

func (svc service) AddPolicy(ctx context.Context, pr PolicyReq) error {
	if err := svc.PolicyValidation(pr); err != nil {
		return err
	}
	return svc.agent.AddPolicy(ctx, pr)
}

func (svc service) PolicyValidation(pr PolicyReq) error {
	if pr.ObjectType == PlatformType && pr.Object != MagistralaObject {
		return errPlatform
	}
	return nil
}

// Yet to do.
func (svc service) AddPolicies(ctx context.Context, token, object string, subjectIDs, relations []string) error {
	key, err := svc.Identify(ctx, token)
	if err != nil {
		return err
	}

	if err := svc.Authorize(ctx, PolicyReq{Object: MagistralaObject, Subject: key.Subject}); err != nil {
		return err
	}

	var errs error
	for _, subjectID := range subjectIDs {
		for _, relation := range relations {
			if err := svc.AddPolicy(ctx, PolicyReq{Object: object, Relation: relation, Subject: subjectID}); err != nil {
				errs = errors.Wrap(fmt.Errorf("cannot add '%s' policy on object '%s' for subject '%s': %w", relation, object, subjectID, err), errs)
			}
		}
	}
	return errs
}

func (svc service) DeletePolicy(ctx context.Context, pr PolicyReq) error {
	return svc.agent.DeletePolicy(ctx, pr)
}

// Yet to do.
func (svc service) DeletePolicies(ctx context.Context, token, object string, subjectIDs, relations []string) error {
	key, err := svc.Identify(ctx, token)
	if err != nil {
		return err
	}

	// Check if the user identified by token is the admin.
	if err := svc.Authorize(ctx, PolicyReq{Object: MagistralaObject, Subject: key.Subject}); err != nil {
		return err
	}

	var errs error
	for _, subjectID := range subjectIDs {
		for _, relation := range relations {
			if err := svc.DeletePolicy(ctx, PolicyReq{Object: object, Relation: relation, Subject: subjectID}); err != nil {
				errs = errors.Wrap(fmt.Errorf("cannot delete '%s' policy on object '%s' for subject '%s': %w", relation, object, subjectID, err), errs)
			}
		}
	}
	return errs
}

func (svc service) ListObjects(ctx context.Context, pr PolicyReq, nextPageToken string, limit int32) (PolicyPage, error) {
	if limit <= 0 {
		limit = 100
	}
	res, npt, err := svc.agent.RetrieveObjects(ctx, pr, nextPageToken, limit)
	if err != nil {
		return PolicyPage{}, err
	}
	var page PolicyPage
	for _, tuple := range res {
		page.Policies = append(page.Policies, tuple.Object)
	}
	page.NextPageToken = npt
	return page, err
}

func (svc service) ListAllObjects(ctx context.Context, pr PolicyReq) (PolicyPage, error) {
	res, err := svc.agent.RetrieveAllObjects(ctx, pr)
	if err != nil {
		return PolicyPage{}, err
	}
	var page PolicyPage
	for _, tuple := range res {
		page.Policies = append(page.Policies, tuple.Object)
	}
	return page, err
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
		return PolicyPage{}, err
	}
	var page PolicyPage
	for _, tuple := range res {
		page.Policies = append(page.Policies, tuple.Subject)
	}
	page.NextPageToken = npt
	return page, err
}

func (svc service) ListAllSubjects(ctx context.Context, pr PolicyReq) (PolicyPage, error) {
	res, err := svc.agent.RetrieveAllSubjects(ctx, pr)
	if err != nil {
		return PolicyPage{}, err
	}
	var page PolicyPage
	for _, tuple := range res {
		page.Policies = append(page.Policies, tuple.Subject)
	}
	return page, err
}

func (svc service) CountSubjects(ctx context.Context, pr PolicyReq) (int, error) {
	return svc.agent.RetrieveAllSubjectsCount(ctx, pr)
}

func (svc service) tmpKey(duration time.Duration, key Key) (Token, error) {
	key.ExpiresAt = time.Now().Add(duration)
	value, err := svc.tokenizer.Issue(key)
	if err != nil {
		return Token{}, errors.Wrap(errIssueTmp, err)
	}

	return Token{AccessToken: value}, nil
}

func (svc service) accessKey(key Key) (Token, error) {
	key.Type = AccessKey
	key.ExpiresAt = time.Now().Add(svc.loginDuration)
	key.Subject = EncodeDomainUserID(key.Domain, key.User)
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

func (svc service) refreshKey(ctx context.Context, token string, key Key) (Token, error) {
	k, err := svc.tokenizer.Parse(token)
	if err != nil {
		return Token{}, err
	}
	if k.Type != RefreshKey {
		return Token{}, errIssueUser
	}
	key.ID = k.ID
	if key.Domain == "" {
		key.Domain = k.Domain
	}
	key.User = k.User
	key.Subject = EncodeDomainUserID(key.Domain, key.User)
	key.Type = AccessKey
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
		return "", "", err
	}
	// Only login key token is valid for login.
	if key.Type != AccessKey || key.Issuer == "" {
		return "", "", errors.ErrAuthentication
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
		return Domain{}, err
	}
	d.CreatedBy = key.User

	domainID, err := svc.idProvider.ID()
	if err != nil {
		return Domain{}, err
	}
	d.ID = domainID

	if d.Status != clients.DisabledStatus && d.Status != clients.EnabledStatus {
		return Domain{}, apiutil.ErrInvalidStatus
	}

	d.CreatedAt = time.Now()

	if err := svc.addDomainPolicy(ctx, key.User, domainID, AdministratorRelation); err != nil {
		return Domain{}, err
	}
	defer func() {
		if err != nil {
			if errRollBack := svc.addDomainPolicyRollback(ctx, key.User, domainID, AdministratorRelation); errRollBack != nil {
				err = errors.Wrap(err, fmt.Errorf("failed to rollback policy %w", errRollBack))
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
		return Domain{}, err
	}

	return svc.domains.RetrieveByID(ctx, id)
}

func (svc service) UpdateDomain(ctx context.Context, token string, id string, d DomainReq) (Domain, error) {
	key, err := svc.Identify(ctx, token)
	if err != nil {
		return Domain{}, err
	}
	if err := svc.Authorize(ctx, PolicyReq{
		Subject:     key.Subject,
		SubjectType: UserType,
		SubjectKind: UsersKind,
		Object:      id,
		ObjectType:  DomainType,
		Permission:  EditPermission,
	}); err != nil {
		return Domain{}, err
	}
	return svc.domains.Update(ctx, id, key.User, d)
}

func (svc service) ChangeDomainStatus(ctx context.Context, token string, id string, d DomainReq) (Domain, error) {
	key, err := svc.Identify(ctx, token)
	if err != nil {
		return Domain{}, err
	}
	if err := svc.Authorize(ctx, PolicyReq{
		Subject:     key.Subject,
		SubjectType: UserType,
		SubjectKind: UsersKind,
		Object:      id,
		ObjectType:  DomainType,
		Permission:  AdminPermission,
	}); err != nil {
		return Domain{}, err
	}
	return svc.domains.Update(ctx, id, key.User, d)
}

func (svc service) ListDomains(ctx context.Context, token string, p Page) (DomainsPage, error) {
	key, err := svc.Identify(ctx, token)
	if err != nil {
		return DomainsPage{}, err
	}
	p.SubjectID = key.User
	return svc.domains.ListDomains(ctx, p)
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
		if err := svc.addDomainPolicy(ctx, userID, id, relation); err != nil {
			return err
		}
	}
	return nil
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

	for _, userID := range userIds {
		if err := svc.removeDomainPolicy(ctx, userID, id, relation); err != nil {
			return err
		}
	}
	return nil
}

func (svc service) ListUserDomains(ctx context.Context, token string, userID string, p Page) (DomainsPage, error) {
	return DomainsPage{}, nil
}

func (svc service) addDomainPolicy(ctx context.Context, userID, domainID, relation string) (err error) {
	pr := PolicyReq{
		Subject:     EncodeDomainUserID(domainID, userID),
		SubjectType: UserType,
		SubjectKind: UsersKind,
		Relation:    relation,
		Object:      domainID,
		ObjectType:  DomainType,
	}
	if err := svc.agent.AddPolicy(ctx, pr); err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if errDel := svc.agent.DeletePolicy(ctx, pr); errDel != nil {
				err = errors.Wrap(err, fmt.Errorf("failed to rollback policy %w", errDel))
			}
		}
	}()
	return svc.domains.SavePolicyCopy(ctx, PolicyCopy{
		SubjectType: UserType,
		SubjectID:   userID,
		Relation:    relation,
		ObjectType:  DomainType,
		ObjectID:    domainID,
	})
}

func (svc service) addDomainPolicyRollback(ctx context.Context, userID, domainID, relation string) error {
	var err error
	pr := PolicyReq{
		Subject:     EncodeDomainUserID(domainID, userID),
		SubjectType: UserType,
		SubjectKind: UsersKind,
		Relation:    relation,
		Object:      domainID,
		ObjectType:  DomainType,
	}
	if errPolicy := svc.agent.DeletePolicy(ctx, pr); errPolicy != nil {
		err = fmt.Errorf("failed to remove from policy engine %w", errPolicy)
	}
	errPolicyCopy := svc.domains.SavePolicyCopy(ctx, PolicyCopy{
		SubjectType: UserType,
		SubjectID:   userID,
		Relation:    relation,
		ObjectType:  DomainType,
		ObjectID:    domainID,
	})
	if errPolicyCopy != nil {
		err = errors.Wrap(err, fmt.Errorf("failed to remove from local policy copy %w", errPolicyCopy))
	}
	return err
}

func (svc service) removeDomainPolicy(ctx context.Context, userID, domainID, relation string) (err error) {
	pr := PolicyReq{
		Subject:     EncodeDomainUserID(domainID, userID),
		SubjectType: UserType,
		SubjectKind: UsersKind,
		Relation:    relation,
		Object:      domainID,
		ObjectType:  DomainType,
	}
	if err := svc.agent.DeletePolicy(ctx, pr); err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if errAdd := svc.agent.AddPolicy(ctx, pr); errAdd != nil {
				err = errors.Wrap(err, fmt.Errorf("failed to add back policy %w", errAdd))
			}
		}
	}()

	return svc.domains.DeletePolicyCopy(ctx, PolicyCopy{
		SubjectType: UserType,
		SubjectID:   userID,
		Relation:    relation,
		ObjectType:  DomainType,
		ObjectID:    domainID,
	})

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

	if len(duid) > 1 {
		return duid[0], duid[1]
	}
	return duid[0], ""
}
