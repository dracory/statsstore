package visitoractivity

import (
	"strings"

	"github.com/dracory/statsstore"
	"github.com/gouniverse/hb"
)

// deviceIcon returns an icon representing the device type
func deviceIcon(visitor statsstore.VisitorInterface) hb.TagInterface {
	deviceType := strings.ToLower(visitor.UserDeviceType())

	iconClass := "bi bi-question-circle"
	color := "text-secondary"

	switch {
	case strings.Contains(deviceType, "desktop"):
		iconClass = "bi bi-display"
		color = "text-primary"
	case strings.Contains(deviceType, "mobile"):
		iconClass = "bi bi-phone"
		color = "text-success"
	case strings.Contains(deviceType, "tablet"):
		iconClass = "bi bi-tablet"
		color = "text-info"
	case strings.Contains(deviceType, "bot"):
		iconClass = "bi bi-robot"
		color = "text-warning"
	}

	return hb.I().Class(iconClass+" "+color).Attr("title", visitor.UserDevice())
}

// osIcon returns an icon representing the operating system
func osIcon(visitor statsstore.VisitorInterface) hb.TagInterface {
	os := strings.ToLower(visitor.UserOs())

	iconClass := "bi bi-circle"
	color := "text-secondary"

	switch {
	case strings.Contains(os, "windows"):
		iconClass = "bi bi-windows"
		color = "text-primary"
	case strings.Contains(os, "mac"), strings.Contains(os, "ios"):
		iconClass = "bi bi-apple"
		color = "text-dark"
	case strings.Contains(os, "android"):
		iconClass = "bi bi-android2"
		color = "text-success"
	case strings.Contains(os, "linux"):
		iconClass = "bi bi-ubuntu"
		color = "text-warning"
	}

	return hb.I().
		Class(iconClass + " " + color).
		Title(visitor.UserOs() + " " + visitor.UserOsVersion())
}
