package shared

// ContextKey is a custom type for context keys to avoid collisions
type ContextKey string

// Context keys
const (
	KeyAdminHomeURL ContextKey = "admin_home_url"
	KeyEndpoint     ContextKey = "endpoint"
)

// Controller name constants
const (
	ControllerHome             = "home"
	ControllerVisitorActivity  = "visitor-activity"
	ControllerVisitorPaths     = "visitor-paths"
	ControllerPageViewActivity = "page-view-activity"
)

// Path constants for admin routes
const (
	PathHome             = "/admin/home"
	PathVisitorActivity  = "/admin/visitor-activity"
	PathVisitorPaths     = "/admin/visitor-paths"
	PathPageViewActivity = "/admin/page-view-activity"
)
