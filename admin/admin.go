package admin

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/dracory/req"
	"github.com/dracory/statsstore"
	"github.com/dracory/statsstore/admin/home"
	pageviewactivity "github.com/dracory/statsstore/admin/page-view-activity"
	"github.com/dracory/statsstore/admin/shared"
	visitoractivity "github.com/dracory/statsstore/admin/visitor-activity"
	visitorpaths "github.com/dracory/statsstore/admin/visitor-paths"
)

type admin struct {
	store             statsstore.StoreInterface
	logger            *slog.Logger
	layout            shared.LayoutInterface
	homeURL           string
	websiteUrl        string
	endpoint          string
	countryNameByIso2 func(iso2Code string) (string, error)
}

var _ http.Handler = (*admin)(nil)

// ============================================================================
// == INTERFACE IMPLEMENTATION
// ============================================================================

// ServeHTTP implements the http.Handler interface
func (a *admin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := req.GetStringOr(r, "path", "home")

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
	options := shared.ControllerOptions{
		Store:             a.store,
		Logger:            a.logger,
		Layout:            a.layout,
		HomeURL:           a.homeURL,
		WebsiteUrl:        a.websiteUrl,
		CountryNameByIso2: a.countryNameByIso2,
	}

	routes := map[string]http.Handler{
		shared.PathHome:             home.New(options),
		shared.PathVisitorActivity:  visitoractivity.New(options),
		shared.PathVisitorPaths:     visitorpaths.New(options),
		shared.PathPageViewActivity: pageviewactivity.New(options),
	}

	if val, ok := routes[path]; ok {
		return val
	}

	return home.New(options)
}
