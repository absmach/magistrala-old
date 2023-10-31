package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/absmach/magistrala/auth"
	"github.com/absmach/magistrala/auth/postgres"
	"github.com/absmach/magistrala/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestAddPolicyCopy(t *testing.T) {
	repo := postgres.NewDomainRepository(database)
	cases := []struct {
		desc string
		pc   auth.PolicyCopy
		err  error
	}{
		{
			desc: "add a  policy copy",
			pc: auth.PolicyCopy{
				SubjectType: "unknown",
				SubjectID:   "unknown",
				Relation:    "unknown",
				ObjectType:  "unknown",
				ObjectID:    "unknown",
			},
			err: nil,
		},
		{
			desc: "add again same policy copy",
			pc: auth.PolicyCopy{
				SubjectType: "unknown",
				SubjectID:   "unknown",
				Relation:    "unknown",
				ObjectType:  "unknown",
				ObjectID:    "unknown",
			},
			err: errors.ErrConflict,
		},
	}

	for _, tc := range cases {
		err := repo.SavePolicyCopy(context.Background(), tc.pc)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %v got %v\n", tc.desc, tc.err, err))
	}
}

func TestDeletePolicyCopy(t *testing.T) {
	repo := postgres.NewDomainRepository(database)
	cases := []struct {
		desc string
		pc   auth.PolicyCopy
		err  error
	}{
		{
			desc: "delete a  policy copy",
			pc: auth.PolicyCopy{
				SubjectType: "unknown",
				SubjectID:   "unknown",
				Relation:    "unknown",
				ObjectType:  "unknown",
				ObjectID:    "unknown",
			},
			err: nil,
		},
	}

	for _, tc := range cases {
		err := repo.DeletePolicyCopy(context.Background(), tc.pc)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %v got %v\n", tc.desc, tc.err, err))
	}
}
