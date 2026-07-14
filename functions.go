package statsstore

import (
	"fmt"
	"strings"

	"github.com/LumenResearch/uasurfer"
)

// == USER AGENT PARSING =======================================================

// userAgentInfo holds parsed user-agent data.
type userAgentInfo struct {
	Browser        string
	BrowserVersion string
	Os             string
	OsVersion      string
	Device         string
	DeviceType     string
}

// parseUserAgent extracts browser, OS, device, and device-type from a
// user-agent string using the uasurfer library. Unknown values are left
// as empty strings.
func parseUserAgent(ua string) userAgentInfo {
	info := userAgentInfo{}
	if ua == "" {
		return info
	}

	parsed := uasurfer.Parse(ua)

	// Browser — uasurfer lacks an Edge constant, so detect it from the
	// Edg/ token that Edge (Chromium) includes in its UA string.
	if strings.Contains(ua, "Edg/") {
		info.Browser = "Edge"
		if parsed.Browser.Version.Major > 0 {
			info.BrowserVersion = versionToString(parsed.Browser.Version)
		}
	} else if parsed.Browser.Name != uasurfer.BrowserUnknown {
		info.Browser = parsed.Browser.Name.StringTrimPrefix()
		info.BrowserVersion = versionToString(parsed.Browser.Version)
	}

	// OS
	if parsed.OS.Name != uasurfer.OSUnknown {
		info.Os = osNameToString(parsed.OS.Name)
	}

	// OS version
	info.OsVersion = versionToString(parsed.OS.Version)

	// Device name (from platform)
	switch parsed.OS.Platform {
	case uasurfer.PlatformiPhone:
		info.Device = "iPhone"
	case uasurfer.PlatformiPad:
		info.Device = "iPad"
	case uasurfer.PlatformiPod:
		info.Device = "iPod"
	}

	// Device type
	info.DeviceType = deviceTypeToString(parsed.DeviceType)

	return info
}

// osNameToString maps uasurfer OSName constants to human-readable strings.
func osNameToString(name uasurfer.OSName) string {
	switch name {
	case uasurfer.OSMacOSX:
		return "macOS"
	case uasurfer.OSiPadOS:
		return "iOS"
	default:
		return name.StringTrimPrefix()
	}
}

// versionToString formats a uasurfer.Version as "major.minor" or
// "major.minor.patch" (when patch > 0). Returns empty string if major is 0.
func versionToString(v uasurfer.Version) string {
	if v.Major == 0 {
		return ""
	}
	if v.Patch > 0 {
		return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	}
	return fmt.Sprintf("%d.%d", v.Major, v.Minor)
}

// deviceTypeToString maps uasurfer DeviceType constants to the lowercase
// values expected by the admin UI (desktop, mobile, tablet).
func deviceTypeToString(dt uasurfer.DeviceType) string {
	switch dt {
	case uasurfer.DeviceComputer:
		return "desktop"
	case uasurfer.DevicePhone:
		return "mobile"
	case uasurfer.DeviceTablet:
		return "tablet"
	case uasurfer.DeviceUnknown:
		return ""
	default:
		return strings.ToLower(dt.StringTrimPrefix())
	}
}
