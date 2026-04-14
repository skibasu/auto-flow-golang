package repository

import (
	"reflect"
	"testing"

	"github.com/skibasu/auto-flow-api/internal/dto"
)

func TestGetUserFilters(t *testing.T) {
	tests := []struct {
		name               string
		filter             dto.UsersFilterRequest
		expectedConditions []string
		expectedArgs       []any
	}{
		{
			name:               "no filters",
			filter:             dto.UsersFilterRequest{},
			expectedConditions: nil,
			expectedArgs:       nil,
		},
		{
			name: "email filter only",
			filter: dto.UsersFilterRequest{
				Email: "john@test.pl",
			},
			expectedConditions: []string{"u.email ILIKE $1"},
			expectedArgs:       []any{"%john@test.pl%"},
		},
		{
			name: "roles filter only",
			filter: dto.UsersFilterRequest{
				Roles: []string{"ADMIN", "MANAGER"},
			},
			expectedConditions: []string{"ur.role = ANY($1)"},
			expectedArgs:       []any{[]string{"ADMIN", "MANAGER"}},
		},
		{
			name: "multiple filters keep sequential placeholders",
			filter: dto.UsersFilterRequest{
				Email:       "test@test.pl",
				FirstName:   "Jan",
				PhoneNumber: "123",
				LastName:    "Kow",
				Roles:       []string{"CLIENT"},
			},
			expectedConditions: []string{
				"u.email ILIKE $1",
				"u.first_name ILIKE $2",
				"u.phone_number ILIKE $3",
				"u.last_name ILIKE $4",
				"ur.role = ANY($5)",
			},
			expectedArgs: []any{"%test@test.pl%", "%Jan%", "%123%", "%Kow%", []string{"CLIENT"}},
		},
		{
			name: "sparse fields still use compact numbering",
			filter: dto.UsersFilterRequest{
				FirstName: "Ada",
				Roles:     []string{"ADMIN"},
			},
			expectedConditions: []string{
				"u.first_name ILIKE $1",
				"ur.role = ANY($2)",
			},
			expectedArgs: []any{"%Ada%", []string{"ADMIN"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conditions, args := getUserFilters(tt.filter)

			if !reflect.DeepEqual(conditions, tt.expectedConditions) {
				t.Fatalf("conditions mismatch\nexpected: %#v\nactual:   %#v", tt.expectedConditions, conditions)
			}

			if !reflect.DeepEqual(args, tt.expectedArgs) {
				t.Fatalf("args mismatch\nexpected: %#v\nactual:   %#v", tt.expectedArgs, args)
			}
		})
	}
}
