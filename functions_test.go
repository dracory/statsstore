package statsstore

import "testing"

func TestParseUserAgent(t *testing.T) {
	tests := []struct {
		name string
		ua   string
		want userAgentInfo
	}{
		{
			name: "Firefox on Windows",
			ua:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:153.0) Gecko/20100101 Firefox/153.0",
			want: userAgentInfo{
				Browser:        "Firefox",
				BrowserVersion: "153.0",
				Os:             "Windows",
				OsVersion:      "10",
				DeviceType:     "desktop",
			},
		},
		{
			name: "Chrome on Windows",
			ua:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			want: userAgentInfo{
				Browser:        "Chrome",
				BrowserVersion: "120.0.0.0",
				Os:             "Windows",
				OsVersion:      "10",
				DeviceType:     "desktop",
			},
		},
		{
			name: "Edge on Windows",
			ua:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0",
			want: userAgentInfo{
				Browser:        "Edge",
				BrowserVersion: "120.0.0.0",
				Os:             "Windows",
				OsVersion:      "10",
				DeviceType:     "desktop",
			},
		},
		{
			name: "Safari on macOS",
			ua:   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15",
			want: userAgentInfo{
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
			want: userAgentInfo{
				Browser:        "Chrome",
				BrowserVersion: "120.0.0.0",
				Os:             "Android",
				OsVersion:      "13",
				DeviceType:     "mobile",
			},
		},
		{
			name: "Safari on iPhone",
			ua:   "Mozilla/5.0 (iPhone; CPU iPhone OS 17_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Mobile/15E148 Safari/604.1",
			want: userAgentInfo{
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
			want: userAgentInfo{
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
			want: userAgentInfo{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseUserAgent(tt.ua)
			if got != tt.want {
				t.Errorf("parseUserAgent(%q)\n  got:  %+v\n  want: %+v", tt.ua, got, tt.want)
			}
		})
	}
}
