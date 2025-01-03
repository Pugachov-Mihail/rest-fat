package source

import "fucking-fat/internal/models"

func MapPermissions(r models.Role) string {
	switch {
	case r.Permission == 0:
		return models.Permissions[r.Permission]
	case r.Permission == 2:
		return models.Permissions[r.Permission]
	case r.Permission == 4:
		return models.Permissions[r.Permission]
	}
	return ""
}
