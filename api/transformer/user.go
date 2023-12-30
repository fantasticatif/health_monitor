package transformer

import "github.com/fantasticatif/health_monitor/data"

func UserAPIResponse(user data.User) map[string]interface{} {
	return map[string]interface{}{
		"uuid":       user.UUID,
		"created_at": dateToString(user.CreatedAt),
		"updated_at": dateToString(user.UpdatedAt),
		"name":       user.Name,
		"email":      user.Email,
	}
}
