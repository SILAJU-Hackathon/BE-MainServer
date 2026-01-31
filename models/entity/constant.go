package entity

const (
	// Report Status
	STATUS_PENDING          = "pending"
	STATUS_COMPLETED        = "complete"
	STATUS_ASSIGNED         = "assigned"
	STATUS_FINISH_BY_WORKER = "finish by worker"
	STATUS_VERIFIED         = "verified"

	// Roles
	ROLE_ADMIN  = "admin"
	ROLE_WORKER = "worker"
	ROLE_USER   = "user"

	// Destruct Class
	DESTRUCT_CLASS_GOOD = "good"

	// Notification Types
	NOTIF_TYPE_REPORT_STATUS   = "report_status"
	NOTIF_TYPE_ACHIEVEMENT     = "achievement"
	NOTIF_TYPE_LEVEL_UP        = "level_up"
	NOTIF_TYPE_REPORT_ASSIGNED = "report_assigned"
	NOTIF_TYPE_REPORT_VERIFIED = "report_verified"
	NOTIF_TYPE_SYSTEM          = "system"
)
