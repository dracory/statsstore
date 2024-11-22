package admin

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gouniverse/hb"
	"github.com/gouniverse/statsstore"
	"github.com/gouniverse/utils"
)

func UI(options UIOptions) (hb.TagInterface, error) {
	if options.ResponseWriter == nil {
		return nil, errors.New("options.ResponseWriter is required")
	}

	if options.Request == nil {
		return nil, errors.New("options.Request is required")
	}

	if options.Store == nil {
		return nil, errors.New("options.Store is required")
	}

	if options.Logger == nil {
		return nil, errors.New("options.Logger is required")
	}

	if options.Layout == nil {
		return nil, errors.New("options.Layout is required")
	}

	ui := &ui{
		response:   options.ResponseWriter,
		request:    options.Request,
		store:      options.Store,
		logger:     *options.Logger,
		layout:     options.Layout,
		homeURL:    options.HomeURL,
		websiteUrl: options.WebsiteUrl,
	}

	return ui.handler(), nil
}

type ui struct {
	response   http.ResponseWriter
	request    *http.Request
	store      statsstore.StoreInterface
	logger     slog.Logger
	layout     Layout
	homeURL    string
	websiteUrl string
}

func (ui *ui) handler() hb.TagInterface {
	controller := utils.Req(ui.request, "controller", "")

	if controller == "" {
		controller = pathHome
	}

	if controller == pathHome {
		return home(*ui)
	}

	if controller == pathVisitorActivity {
		return visitorActivity(*ui)
	}

	if controller == pathVisitorPaths {
		return visitorPaths(*ui)
	}

	ui.layout.SetBody(hb.H1().HTML(controller).ToHTML())
	return hb.Raw(ui.layout.Render(ui.response, ui.request))
	// redirect(a.response, a.request, url(a.request, pathQueueManager, map[string]string{}))
	// return nil
}

type Layout interface {
	SetTitle(title string)
	SetScriptURLs(scripts []string)
	SetScripts(scripts []string)
	SetStyleURLs(styles []string)
	SetStyles(styles []string)
	SetBody(string)
	Render(w http.ResponseWriter, r *http.Request) string
}

type UIOptions struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Logger         *slog.Logger
	Store          statsstore.StoreInterface
	Layout         Layout
	HomeURL        string
	WebsiteUrl     string
}

type PageInterface interface {
	hb.TagInterface
	ToTag(w http.ResponseWriter, r *http.Request) hb.TagInterface
}
