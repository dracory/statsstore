// Package admin provides a self-hosted visitor analytics admin dashboard.
//
// It is a self-contained, reusable library package: consumers inject their own
// LayoutInterface, a statsstore.StoreInterface, and optional callbacks via
// admin.New(admin.Options{...}), and get the full dashboard for free.
//
// The package follows the folder-per-controller pattern: each controller lives
// in its own subfolder and renders via the injected LayoutInterface.
package admin

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/dracory/statsstore"
	"github.com/dracory/statsstore/admin/shared"
)

// Options contains all dependencies and configuration for the admin dashboard.
// It extends the original statsstore admin Options with optional callbacks for
// auth, flash messages, and base URL configuration so the package remains
// self-contained and reusable by any project.
type Options struct {
	// Store is the statsstore instance used for all visitor queries. Required.
	Store statsstore.StoreInterface

	// Layout renders the full HTML page. Required.
	Layout shared.LayoutInterface

	// HomeURL is the URL for the admin home page (e.g. "/admin").
	HomeURL string

	// WebsiteUrl is the public site URL.
	WebsiteUrl string

	// BaseURL is the base URL for the stats admin (e.g. "/admin/stats").
	// Defaults to "/admin/stats" when empty.
	BaseURL string

	// CountryNameByIso2 maps an ISO2 code to a human-readable country name.
	// Optional; when nil, raw ISO2 codes are displayed.
	CountryNameByIso2 func(iso2Code string) (string, error)

	// AuthUserID returns the authenticated user ID from the request, or ""
	// when unauthenticated. When nil, auth checks are skipped.
	AuthUserID func(r *http.Request) string

	// FlashError renders a flash error message and returns the redirect
	// response body. When nil, controllers fall back to http.Redirect.
	FlashError func(w http.ResponseWriter, r *http.Request, message string, redirectURL string, delaySeconds int) string

	// Logger is used for internal logging. Defaults to slog.Default().
	Logger *slog.Logger
}

// New creates a new admin instance. It validates required fields and returns
// an http.Handler that dispatches to the appropriate controller based on the
// "controller" query parameter.
func New(opts Options) (http.Handler, error) {
	if opts.Store == nil {
		return nil, errors.New("store is required")
	}
	if opts.Layout == nil {
		return nil, errors.New("layout is required")
	}
	if opts.HomeURL == "" {
		return nil, errors.New("home URL is required")
	}
	if opts.BaseURL == "" {
		opts.BaseURL = "/admin/stats"
	}

	logger := slog.Default()
	if opts.Logger != nil {
		logger = opts.Logger
	}
	_ = logger // reserved for future use

	return &adminHandler{opts: opts}, nil
}

// adminHandler implements http.Handler.
type adminHandler struct {
	opts Options
}

// ServeHTTP implements the http.Handler interface.
func (a *adminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	controllerOpts := shared.ControllerOptions{
		Store:             a.opts.Store,
		Layout:            a.opts.Layout,
		HomeURL:           a.opts.HomeURL,
		WebsiteUrl:        a.opts.WebsiteUrl,
		BaseURL:           a.opts.BaseURL,
		CountryNameByIso2: a.opts.CountryNameByIso2,
		AuthUserID:        a.opts.AuthUserID,
		FlashError:        a.opts.FlashError,
		Logger:            a.opts.Logger,
	}

	Handler(controllerOpts, w, r)
}
