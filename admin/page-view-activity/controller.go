package pageviewactivity

import (
	"net/http"

	"github.com/dracory/hb"
	"github.com/dracory/statsstore/admin/shared"
)

// New creates a new page view activity controller.
func New(ui ControllerOptions) http.Handler {
	return &Controller{ui: ui}
}

// Controller handles rendering the page view activity screen.
type Controller struct {
	ui ControllerOptions
}

// ServeHTTP implements the http.Handler interface.
func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(c.Handler(w, r)))
}

// Handler prepares the layout and returns the rendered HTML.
func (c *Controller) Handler(w http.ResponseWriter, r *http.Request) string {
	data, errorMessage := buildControllerData(r, c.ui.Store)

	c.ui.Layout.SetTitle("Page View Activity | Visitor Analytics")

	if errorMessage != "" {
		c.ui.Layout.SetBody(hb.Div().
			Class("alert alert-danger").
			Text(errorMessage).ToHTML())

		return c.ui.Layout.Render(w, r)
	}

	c.ui.Layout.SetBody(c.page(data).ToHTML())

	return c.ui.Layout.Render(w, r)
}

// ToTag renders the controller to an HTML tag (useful for embedding).
func (c *Controller) ToTag(w http.ResponseWriter, r *http.Request) hb.TagInterface {
	return hb.Raw(c.Handler(w, r))
}

// page constructs the main page container.
func (c *Controller) page(data ControllerData) hb.TagInterface {
	breadcrumbs := shared.Breadcrumbs(data.Request, []shared.Breadcrumb{
		{
			Name: "Home",
			URL:  shared.UrlHome(data.Request),
		},
		{
			Name: "Visitor Analytics",
			URL:  shared.UrlHome(data.Request),
		},
		{
			Name: "Page View Activity",
			URL:  shared.UrlPageViewActivity(data.Request),
		},
	})

	title := hb.Heading1().
		Class("mt-3 mb-4 text-primary").
		HTML("Page View Activity")

	body := hb.Div().
		Class("alert alert-info").
		Text("Page View Activity UI implementation is in progress.")

	return hb.Div().
		Class("container").
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(shared.AdminHeaderUI(data.Request, c.ui.HomeURL)).
		Child(hb.HR()).
		Child(title).
		Child(body)
}
