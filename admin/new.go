package admin

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gouniverse/statsstore"
	"github.com/gouniverse/statsstore/admin/shared"
)

type Options struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Logger         *slog.Logger
	Store          statsstore.StoreInterface
	Layout         shared.LayoutInterface
	HomeURL        string
	WebsiteUrl     string
	Endpoint       string
}

func New(options Options) (http.Handler, error) {
	if options.ResponseWriter == nil {
		return nil, errors.New("response writer is required")
	}

	if options.Request == nil {
		return nil, errors.New("request is required")
	}

	if options.Store == nil {
		return nil, errors.New("store is required")
	}

	if options.Layout == nil {
		return nil, errors.New("layout is required")
	}

	if options.HomeURL == "" {
		return nil, errors.New("home URL is required")
	}

	logger := slog.Default()
	if options.Logger != nil {
		logger = options.Logger
	}

	adminInstance := &admin{
		response:   options.ResponseWriter,
		request:    options.Request,
		store:      options.Store,
		logger:     logger,
		layout:     options.Layout,
		homeURL:    options.HomeURL,
		websiteUrl: options.WebsiteUrl,
		endpoint:   options.Endpoint,
	}

	return adminInstance, nil
}
