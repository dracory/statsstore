package statsstore

import "testing"

func TestParseUserAgent(t *testing.T) {
	tests := []struct {
		name string
		ua   string
		want UserAgentInfo
	}{
		{
			name: "Firefox on Windows",
			ua:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:153.0) Gecko/20100101 Firefox/153.0",
			want: UserAgentInfo{
				Browser:        "Firefox",
				BrowserVersion: "153.0",
				Os:             "Windows",
				OsVersion:      "10.0",
				DeviceType:     "desktop",
			},
		},
		{
			name: "Chrome on Windows",
			ua:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			want: UserAgentInfo{
				Browser:        "Chrome",
				BrowserVersion: "120.0",
				Os:             "Windows",
				OsVersion:      "10.0",
				DeviceType:     "desktop",
			},
		},
		{
			name: "Edge on Windows",
			ua:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0",
			want: UserAgentInfo{
				Browser:        "Edge",
				BrowserVersion: "120.0",
				Os:             "Windows",
				OsVersion:      "10.0",
				DeviceType:     "desktop",
			},
		},
		{
			name: "Safari on macOS",
			ua:   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15",
			want: UserAgentInfo{
				Browser:        "Safari",
				BrowserVersion: "17.1",
				Os:             "macOS",
				OsVersion:      "10.15.7",
				DeviceType:     "desktop",
			},
		},
		{
			name: "Chrome on Android mobile",
			ua:   "Mozilla/5.0 (Linux; Android 13; Pixel 7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36",
			want: UserAgentInfo{
				Browser:        "Chrome",
				BrowserVersion: "120.0",
				Os:             "Android",
				OsVersion:      "13.0",
				DeviceType:     "mobile",
			},
		},
		{
			name: "Safari on iPhone",
			ua:   "Mozilla/5.0 (iPhone; CPU iPhone OS 17_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Mobile/15E148 Safari/604.1",
			want: UserAgentInfo{
				Browser:        "Safari",
				BrowserVersion: "17.2",
				Os:             "iOS",
				OsVersion:      "17.2",
				Device:         "iPhone",
				DeviceType:     "mobile",
			},
		},
		{
			name: "Safari on iPad",
			ua:   "Mozilla/5.0 (iPad; CPU OS 17_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Mobile/15E148 Safari/604.1",
			want: UserAgentInfo{
				Browser:        "Safari",
				BrowserVersion: "17.2",
				Os:             "iOS",
				OsVersion:      "17.2",
				Device:         "iPad",
				DeviceType:     "tablet",
			},
		},
		{
			name: "Empty UA",
			ua:   "",
			want: UserAgentInfo{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseUserAgent(tt.ua)
			if got != tt.want {
				t.Errorf("ParseUserAgent(%q)\n  got:  %+v\n  want: %+v", tt.ua, got, tt.want)
			}
		})
	}
}

func TestVisitorStringTrimming(t *testing.T) {
	v := NewVisitor()

	// Truncate ASCII and multibyte strings safely
	longStr := "a" + string([]rune{'😀', '😁', '😂'}) + "b"
	v.SetUserBrowser(longStr)
	if v.GetUserBrowser() != longStr {
		t.Fatalf("expected string within limit to remain untouched")
	}

	longBrowser := "SuperBrowserWithExceedinglyLongNameThatExceedsTheLimitByFar" + "123456789012345678901234567890123456789012345678901234567890"
	v.SetUserBrowser(longBrowser)
	if len(v.GetUserBrowser()) > MAX_LEN_USER_BROWSER {
		t.Fatalf("expected UserBrowser truncated to <= %d, got %d", MAX_LEN_USER_BROWSER, len(v.GetUserBrowser()))
	}

	longOS := "SuperOperatingSystemWithLongName" + "1234567890123456789012345678901234567890"
	v.SetUserOs(longOS)
	if len(v.GetUserOs()) > MAX_LEN_USER_OS {
		t.Fatalf("expected UserOs truncated to <= %d, got %d", MAX_LEN_USER_OS, len(v.GetUserOs()))
	}

	longDevice := "SuperDeviceWithLongName" + "12345678901234567890123456789012345678901234567890123456789012345678901234567890"
	v.SetUserDevice(longDevice)
	if len(v.GetUserDevice()) > MAX_LEN_USER_DEVICE {
		t.Fatalf("expected UserDevice truncated to <= %d, got %d", MAX_LEN_USER_DEVICE, len(v.GetUserDevice()))
	}

	// Multibyte string truncation check
	multibyteStr := ""
	for i := 0; i < 60; i++ {
		multibyteStr += "🔥"
	}
	v.SetUserBrowser(multibyteStr)
	if len(v.GetUserBrowser()) > MAX_LEN_USER_BROWSER {
		t.Fatalf("expected byte length of UserBrowser to be <= %d, got %d", MAX_LEN_USER_BROWSER, len(v.GetUserBrowser()))
	}
}
