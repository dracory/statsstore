package shared

import (
	"log/slog"
	"net/http"

	"github.com/dracory/statsstore"
)

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

// ControllerOptions contains the options for creating a new admin controller
type ControllerOptions struct {
	Logger     *slog.Logger
	Store      statsstore.StoreInterface
	Layout     LayoutInterface
	HomeURL    string
	WebsiteUrl string
}
