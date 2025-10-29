package visitoractivity

import (
	"net/http"

	"github.com/dracory/statsstore/admin/shared"
	"github.com/gouniverse/hb"
)

type TableView struct {
	ui shared.ControllerOptions
}

func NewTableView(ui shared.ControllerOptions) *TableView {
	return &TableView{ui: ui}
}

func (tv *TableView) Render(w http.ResponseWriter, r *http.Request) hb.TagInterface {
	data, errorMessage := tv.prepareData(r)
	if errorMessage != "" {
		return hb.Div().Class("alert alert-danger").Text(errorMessage)
	}

	tv.ui.Layout.SetTitle("Visitor Activity | Visitor Analytics")

	breadcrumbs := shared.Breadcrumbs(r, []shared.Breadcrumb{
		{Name: "Home", URL: tv.ui.HomeURL},
		{Name: "Visitor Analytics", URL: shared.UrlHome(r)},
		{Name: "Visitor Activity", URL: shared.UrlVisitorActivity(r)},
	})

	title := hb.Heading1().
		Class("mt-3 mb-4 text-primary").
		HTML("Visitor Activity")

	return hb.Div().
		Class("container").
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(shared.AdminHeaderUI(r, tv.ui.HomeURL)).
		Child(hb.HR()).
		Child(title).
		Child(CardVisitorActivity(*data))
}

func (tv *TableView) prepareData(r *http.Request) (*ControllerData, string) {
	data, errMsg := buildControllerData(r, tv.ui.Store)
	if errMsg != "" {
		return &data, errMsg
	}
	return &data, ""
}
