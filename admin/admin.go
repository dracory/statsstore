package admin

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gouniverse/base/req"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/statsstore"
	"github.com/gouniverse/statsstore/admin/home"
	"github.com/gouniverse/statsstore/admin/shared"
	visitoractivity "github.com/gouniverse/statsstore/admin/visitor-activity"
	visitorpaths "github.com/gouniverse/statsstore/admin/visitor-paths"
)

type admin struct {
	response   http.ResponseWriter
	request    *http.Request
	store      statsstore.StoreInterface
	logger     *slog.Logger
	layout     shared.LayoutInterface
	homeURL    string
	websiteUrl string
	endpoint   string
}

var _ http.Handler = (*admin)(nil)

// ============================================================================
// == INTERFACE IMPLEMENTATION
// ============================================================================

// ServeHTTP implements the http.Handler interface
func (a *admin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := req.ValueOr(r, "path", "home")

	if path == "" {
		path = shared.PathHome
	}

	// Use the custom ContextKey type for context values
	ctx := context.WithValue(r.Context(), shared.KeyEndpoint, r.URL.Path)
	ctx = context.WithValue(ctx, shared.KeyAdminHomeURL, a.homeURL)

	a.findHandlerFromPath(path).ServeHTTP(w, r.WithContext(ctx))
}

// ============================================================================
// == PRIVATE METHODS
// ============================================================================

func (a *admin) findHandlerFromPath(path string) http.Handler {
	routes := map[string]http.Handler{
		shared.PathHome:            home.New(a),
		shared.PathVisitorActivity: visitoractivity.New(a),
		shared.PathVisitorPaths:    visitorpaths.New(a),
	}

	if val, ok := routes[path]; ok {
		return val
	}

	return home.New(a)
}

// Implement the shared.UIContext interface
func (a *admin) GetResponse() http.ResponseWriter {
	return a.response
}

func (a *admin) GetRequest() *http.Request {
	return a.request
}

func (a *admin) GetStore() statsstore.StoreInterface {
	return a.store
}

func (a *admin) GetLayout() shared.LayoutInterface {
	return a.layout
}

func (a *admin) GetHomeURL() string {
	return a.homeURL
}

func (a *admin) GetPathHome() string {
	return shared.ControllerHome
}

func (a *admin) GetPathVisitorActivity() string {
	return shared.ControllerVisitorActivity
}

func (a *admin) GetPathVisitorPaths() string {
	return shared.ControllerVisitorPaths
}

func (a *admin) URL(path string, params map[string]string) string {
	return shared.URL(a.request, path, params)
}

func (a *admin) Breadcrumbs(items []shared.Breadcrumb) hb.TagInterface {
	return shared.Breadcrumbs(a.request, items)
}

func (a *admin) AdminHeader() hb.TagInterface {
	return shared.AdminHeaderUI(a)
}
