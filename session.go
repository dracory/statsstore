package statsstore

import (
	"time"

	"github.com/dromara/carbon/v2"
)

// SessionInterface defines a visitor session: a group of page views sharing a
// session id (same fingerprint, consecutive views within SESSION_GAP_MINUTES).
type SessionInterface interface {
	GetSessionID() string
	GetFingerprint() string
	GetIpAddress() string
	GetCountry() string
	GetUserDeviceType() string
	GetUserBrowser() string
	GetUserOs() string

	GetStartedAt() string
	GetStartedAtCarbon() *carbon.Carbon
	GetEndedAt() string
	GetEndedAtCarbon() *carbon.Carbon

	// GetDurationSeconds returns ended_at - started_at in seconds.
	GetDurationSeconds() int64

	GetPageViews() int64
	GetBotRows() int64

	GetEntryPath() string
	SetEntryPath(path string) SessionInterface
	GetExitPath() string
	SetExitPath(path string) SessionInterface

	// IsBounce reports whether the session consisted of a single page view.
	IsBounce() bool
}

var _ SessionInterface = (*sessionImplementation)(nil)

type sessionImplementation struct {
	sessionID      string
	fingerprint    string
	ipAddress      string
	country        string
	userDeviceType string
	userBrowser    string
	userOs         string
	startedAt      time.Time
	endedAt        time.Time
	pageViews      int64
	botRows        int64
	entryPath      string
	exitPath       string
}

func (o *sessionImplementation) GetSessionID() string      { return o.sessionID }
func (o *sessionImplementation) GetFingerprint() string    { return o.fingerprint }
func (o *sessionImplementation) GetIpAddress() string      { return o.ipAddress }
func (o *sessionImplementation) GetCountry() string        { return o.country }
func (o *sessionImplementation) GetUserDeviceType() string { return o.userDeviceType }
func (o *sessionImplementation) GetUserBrowser() string    { return o.userBrowser }
func (o *sessionImplementation) GetUserOs() string         { return o.userOs }

func (o *sessionImplementation) GetStartedAt() string {
	return carbon.CreateFromStdTime(o.startedAt, carbon.UTC).ToDateTimeString(carbon.UTC)
}

func (o *sessionImplementation) GetStartedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.startedAt, carbon.UTC)
}

func (o *sessionImplementation) GetEndedAt() string {
	return carbon.CreateFromStdTime(o.endedAt, carbon.UTC).ToDateTimeString(carbon.UTC)
}

func (o *sessionImplementation) GetEndedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.endedAt, carbon.UTC)
}

func (o *sessionImplementation) GetDurationSeconds() int64 {
	return int64(o.endedAt.Sub(o.startedAt).Seconds())
}

func (o *sessionImplementation) GetPageViews() int64 { return o.pageViews }
func (o *sessionImplementation) GetBotRows() int64   { return o.botRows }

func (o *sessionImplementation) GetEntryPath() string { return o.entryPath }
func (o *sessionImplementation) SetEntryPath(path string) SessionInterface {
	o.entryPath = path
	return o
}
func (o *sessionImplementation) GetExitPath() string { return o.exitPath }
func (o *sessionImplementation) SetExitPath(path string) SessionInterface {
	o.exitPath = path
	return o
}

func (o *sessionImplementation) IsBounce() bool {
	return o.pageViews <= 1
}
