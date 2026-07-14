package statsstore

import (
	"regexp"
	"strings"
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

// Precompiled regexes for user-agent parsing.
var (
	reEdge    = regexp.MustCompile(`Edg(?:e|A|iOS)?/([\d.]+)`)
	reOpera   = regexp.MustCompile(`OPR/([\d.]+)`)
	reFirefox = regexp.MustCompile(`Firefox/([\d.]+)`)
	reChrome  = regexp.MustCompile(`Chrome/([\d.]+)`)
	reSafari  = regexp.MustCompile(`Version/([\d.]+).*Safari`)
	reAndroid = regexp.MustCompile(`Android ([\d.]+)`)
	reWinNT   = regexp.MustCompile(`Windows NT ([\d.]+)`)
	reIOS     = regexp.MustCompile(`OS ([\d_]+)`)
	reMacOS   = regexp.MustCompile(`Mac OS X ([\d_]+)`)
)

// parseUserAgent extracts browser, OS, device, and device-type from a
// user-agent string. Unknown values are left as empty strings.
func parseUserAgent(ua string) userAgentInfo {
	info := userAgentInfo{}
	if ua == "" {
		return info
	}

	uaLower := strings.ToLower(ua)

	// --- Browser detection (order matters: specific first) ---
	switch {
	case reEdge.MatchString(ua):
		if m := reEdge.FindStringSubmatch(ua); m != nil {
			info.Browser = "Edge"
			info.BrowserVersion = m[1]
		}
	case reOpera.MatchString(ua):
		if m := reOpera.FindStringSubmatch(ua); m != nil {
			info.Browser = "Opera"
			info.BrowserVersion = m[1]
		}
	case reFirefox.MatchString(ua):
		if m := reFirefox.FindStringSubmatch(ua); m != nil {
			info.Browser = "Firefox"
			info.BrowserVersion = m[1]
		}
	case reChrome.MatchString(ua):
		if m := reChrome.FindStringSubmatch(ua); m != nil {
			info.Browser = "Chrome"
			info.BrowserVersion = m[1]
		}
	case reSafari.MatchString(ua):
		if m := reSafari.FindStringSubmatch(ua); m != nil {
			info.Browser = "Safari"
			info.BrowserVersion = m[1]
		}
	case strings.Contains(uaLower, "safari"):
		info.Browser = "Safari"
	}

	// --- OS detection ---
	switch {
	case reAndroid.MatchString(ua):
		if m := reAndroid.FindStringSubmatch(ua); m != nil {
			info.Os = "Android"
			info.OsVersion = m[1]
		}
	case reWinNT.MatchString(ua):
		if m := reWinNT.FindStringSubmatch(ua); m != nil {
			info.Os = "Windows"
			switch m[1] {
			case "10.0":
				info.OsVersion = "10"
			case "6.3":
				info.OsVersion = "8.1"
			case "6.2":
				info.OsVersion = "8"
			case "6.1":
				info.OsVersion = "7"
			case "6.0":
				info.OsVersion = "Vista"
			case "5.1":
				info.OsVersion = "XP"
			default:
				info.OsVersion = m[1]
			}
		}
	case strings.Contains(uaLower, "iphone") || strings.Contains(uaLower, "ipad") || strings.Contains(uaLower, "ipod"):
		info.Os = "iOS"
		if m := reIOS.FindStringSubmatch(ua); m != nil {
			info.OsVersion = strings.ReplaceAll(m[1], "_", ".")
		}
	case reMacOS.MatchString(ua):
		if m := reMacOS.FindStringSubmatch(ua); m != nil {
			info.Os = "macOS"
			info.OsVersion = strings.ReplaceAll(m[1], "_", ".")
		}
	case strings.Contains(uaLower, "linux"):
		info.Os = "Linux"
	}

	// --- Device & device-type detection ---
	switch {
	case strings.Contains(uaLower, "ipad"):
		info.DeviceType = "tablet"
		info.Device = "iPad"
	case strings.Contains(uaLower, "iphone"):
		info.DeviceType = "mobile"
		info.Device = "iPhone"
	case strings.Contains(uaLower, "android"):
		if strings.Contains(uaLower, "tablet") || !strings.Contains(uaLower, "mobile") {
			info.DeviceType = "tablet"
		} else {
			info.DeviceType = "mobile"
		}
	case strings.Contains(uaLower, "mobile"):
		info.DeviceType = "mobile"
	case info.Os != "" && info.Os != "Android" && info.Os != "iOS":
		info.DeviceType = "desktop"
	}

	return info
}
