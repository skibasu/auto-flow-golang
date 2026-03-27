package repository

import (
	"fmt"

	"github.com/skibasu/auto-flow-api/internal/dto"
)

// MOVE TO HELPERS
func getUserFilters(filter dto.UsersFilterRequest) (conditions []string, args []any) {

	argID := 1

	if filter.Email != "" {
		conditions = append(conditions, fmt.Sprintf("u.email ILIKE $%d", argID))
		args = append(args, "%"+filter.Email+"%")
		argID++
	}

	if filter.FirstName != "" {
		conditions = append(conditions, fmt.Sprintf("u.first_name ILIKE $%d", argID))
		args = append(args, "%"+filter.FirstName+"%")
		argID++
	}
	if filter.PhoneNumber != "" {
		conditions = append(conditions, fmt.Sprintf("u.phone_number ILIKE $%d", argID))
		args = append(args, "%"+filter.PhoneNumber+"%")
		argID++
	}
	if filter.LastName != "" {
		conditions = append(conditions, fmt.Sprintf("u.last_name ILIKE $%d", argID))
		args = append(args, "%"+filter.LastName+"%")
		argID++
	}

	// 🔥 roles filter
	if len(filter.Roles) > 0 {
		conditions = append(conditions, fmt.Sprintf("ur.role = ANY($%d)", argID))
		args = append(args, filter.Roles)
		argID++
	}
	return conditions, args
}
