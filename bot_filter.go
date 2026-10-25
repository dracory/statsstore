package statsstore

import "github.com/dracory/botdetect"

// Bot detection lives in github.com/dracory/botdetect — this file keeps
// statsstore's exported function names working by delegating to it, so
// existing consumers (statsadmin and downstream apps) don't break.

// IsBot checks whether a user-agent string belongs to a bot, crawler, or
// automated tool. An empty user-agent returns false.
func IsBot(userAgent string) bool {
	return botdetect.IsBot(userAgent)
}

// BotPathPatterns returns the list of path substrings for bot-only files
// (robots.txt, ads.txt, sitemap.xml, etc.).
func BotPathPatterns() []string {
	return botdetect.BotPathPatterns()
}

// MaliciousPathPatterns returns the list of path substrings for universally
// malicious endpoints (.env, .git/, .svn/, .htpasswd, shell.php).
func MaliciousPathPatterns() []string {
	return botdetect.MaliciousPathPatterns()
}

// IsBotPath reports whether a request path targets a bot-only file.
func IsBotPath(path string) bool {
	return botdetect.IsBotPath(path)
}

// IsMaliciousPath reports whether a request path targets a universally
// malicious endpoint.
func IsMaliciousPath(path string) bool {
	return botdetect.IsMaliciousPath(path)
}

// IsReferrerSpam checks whether a referrer URL host matches a known spam
// domain.
func IsReferrerSpam(referrer string) bool {
	return botdetect.IsReferrerSpam(referrer)
}

// IsDataCenterIP checks whether an IP address falls within known data
// center CIDR ranges.
func IsDataCenterIP(ip string) bool {
	return botdetect.IsDataCenterIP(ip)
}
