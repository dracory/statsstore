package shared

import (
	"net/http"

	"github.com/gouniverse/hb"
	"github.com/gouniverse/statsstore"
)

// UIContext represents the interface for UI context needed by controllers
type UIContext interface {
	GetResponse() http.ResponseWriter
	GetRequest() *http.Request
	GetStore() statsstore.StoreInterface
	GetLayout() LayoutInterface
	GetHomeURL() string
	GetPathHome() string
	GetPathVisitorActivity() string
	GetPathVisitorPaths() string
	URL(path string, params map[string]string) string
	Breadcrumbs(items []Breadcrumb) hb.TagInterface
	AdminHeader() hb.TagInterface
}

// LayoutInterface defines the layout methods needed by controllers
type LayoutInterface interface {
	SetTitle(title string)
	SetScriptURLs(scripts []string)
	SetScripts(scripts []string)
	SetStyleURLs(styles []string)
	SetStyles(styles []string)
	SetBody(string)
	Render(w http.ResponseWriter, r *http.Request) string
}

// Breadcrumb represents a navigation breadcrumb
type Breadcrumb struct {
	Name string
	URL  string
}

// // UIOptions contains the options for creating a new admin UI
// type UIOptions struct {
// 	ResponseWriter http.ResponseWriter
// 	Request        *http.Request
// 	Logger         *slog.Logger
// 	Store          statsstore.StoreInterface
// 	Layout         LayoutInterface
// 	HomeURL        string
// 	WebsiteUrl     string
// 	Endpoint       string
// }
