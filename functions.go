package statsstore

import "github.com/dracory/botdetect"

// == USER AGENT PARSING =======================================================

// UserAgentInfo holds parsed user-agent data. Alias of
// botdetect.UserAgentInfo — all name fields are lowercase.
type UserAgentInfo = botdetect.UserAgentInfo

// ParseUserAgent extracts browser, OS, device, and device-type from a
// user-agent string. Delegates to botdetect.ParseUserAgent; unknown values
// are left as empty strings, and bots report DeviceType "bot".
func ParseUserAgent(ua string) UserAgentInfo {
	return botdetect.ParseUserAgent(ua)
}
