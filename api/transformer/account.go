package transformer

import "github.com/fantasticatif/health_monitor/data"

func AccountAPIResponse(acc data.Account) map[string]interface{} {
	return map[string]interface{}{
		"uuid":       acc.UUID,
		"created_at": dateToString(acc.CreatedAt),
		"name":       acc.Name,
		"updated_at": dateToString(acc.UpdatedAt),
	}
}
