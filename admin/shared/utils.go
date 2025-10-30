package shared

import (
	"net/http"
	urlpkg "net/url"

	"github.com/dracory/hb"
	"github.com/samber/lo"
)

func UrlHome(r *http.Request, params ...map[string]string) string {
	endpoint := lo.IfF(r.Context().Value(KeyEndpoint) != nil, func() string { return r.Context().Value(KeyEndpoint).(string) }).Else("/")

	p := lo.IfF(len(params) > 0, func() map[string]string { return params[0] }).Else(map[string]string{})

	p["path"] = PathHome

	return URL(r, endpoint, p)
}

func UrlVisitorActivity(r *http.Request, params ...map[string]string) string {
	endpoint := lo.IfF(r.Context().Value(KeyEndpoint) != nil, func() string { return r.Context().Value(KeyEndpoint).(string) }).Else("/")

	p := lo.IfF(len(params) > 0, func() map[string]string { return params[0] }).Else(map[string]string{})

	p["path"] = PathVisitorActivity

	return URL(r, endpoint, p)
}

func UrlVisitorPaths(r *http.Request, params ...map[string]string) string {
	endpoint := lo.IfF(r.Context().Value(KeyEndpoint) != nil, func() string { return r.Context().Value(KeyEndpoint).(string) }).Else("/")

	p := lo.IfF(len(params) > 0, func() map[string]string { return params[0] }).Else(map[string]string{})

	p["path"] = PathVisitorPaths

	return URL(r, endpoint, p)
}

func UrlPageViewActivity(r *http.Request, params ...map[string]string) string {
	endpoint := lo.IfF(r.Context().Value(KeyEndpoint) != nil, func() string { return r.Context().Value(KeyEndpoint).(string) }).Else("/")

	p := lo.IfF(len(params) > 0, func() map[string]string { return params[0] }).Else(map[string]string{})

	p["path"] = PathPageViewActivity

	return URL(r, endpoint, p)
}

// URL generates a URL with the given path and parameters
func URL(r *http.Request, path string, params map[string]string) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	baseURL := scheme + "://" + r.Host
	url := baseURL + path

	if params != nil {
		url += "?" + Query(params)
	}

	return url
}

// Query generates a query string from a map of parameters
func Query(queryData map[string]string) string {
	values := urlpkg.Values{}

	for key, value := range queryData {
		values.Add(key, value)
	}

	return HTTPBuildQuery(values)
}

// HTTPBuildQuery converts URL values to a query string
func HTTPBuildQuery(queryData urlpkg.Values) string {
	return queryData.Encode()
}

// Breadcrumbs creates a breadcrumb navigation component
func Breadcrumbs(r *http.Request, pageBreadcrumbs []Breadcrumb) hb.TagInterface {
	breadcrumbs := []Breadcrumb{}
	breadcrumbs = append(breadcrumbs, pageBreadcrumbs...)
	return BreadcrumbsUI(breadcrumbs)
}

// BreadcrumbsUI generates the UI for breadcrumbs
func BreadcrumbsUI(breadcrumbs []Breadcrumb) hb.TagInterface {
	if len(breadcrumbs) == 0 {
		return hb.Div()
	}

	nav := hb.Nav().
		Attr("aria-label", "breadcrumb").
		Child(hb.OL().
			Class("breadcrumb").
			Children(lo.Map(breadcrumbs, func(breadcrumb Breadcrumb, index int) hb.TagInterface {
				if index == len(breadcrumbs)-1 {
					return hb.LI().
						Class("breadcrumb-item active").
						Attr("aria-current", "page").
						HTML(breadcrumb.Name)
				}

				return hb.LI().
					Class("breadcrumb-item").
					Child(hb.A().
						Href(breadcrumb.URL).
						HTML(breadcrumb.Name))
			})))

	return nav
}

// StrTruncate truncates a string to the specified length and adds ellipsis if needed
func StrTruncate(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength] + "..."
}
