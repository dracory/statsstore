package statsstore

import "github.com/golang-module/carbon/v2"

type VisitorInterface interface {
	// From dataobject

	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()
	Path() string
	SetPath(path string) VisitorInterface
	Country() string
	SetCountry(country string) VisitorInterface
	CreatedAt() string
	CreatedAtCarbon() carbon.Carbon
	SetCreatedAt(createdAt string) VisitorInterface
	DeletedAt() string
	SetDeletedAt(deletedAt string) VisitorInterface
	ID() string
	SetID(id string) VisitorInterface
	IpAddress() string
	SetIpAddress(ip string) VisitorInterface
	UserAcceptLanguage() string
	SetUserAcceptLanguage(userAcceptLanguage string) VisitorInterface
	UserAcceptEncoding() string
	SetUserAcceptEncoding(userAcceptEncoding string) VisitorInterface
	UserAgent() string
	SetUserAgent(userAgent string) VisitorInterface
	UserBrowser() string
	SetUserBrowser(userBrowser string) VisitorInterface
	UserBrowserVersion() string
	SetUserBrowserVersion(userBrowserVersion string) VisitorInterface
	UserDevice() string
	SetUserDevice(userDevice string) VisitorInterface
	UserDeviceType() string
	SetUserDeviceType(userDeviceType string) VisitorInterface
	UserOs() string
	SetUserOs(userOs string) VisitorInterface
	UserOsVersion() string
	SetUserOsVersion(userOsVersion string) VisitorInterface
	UserReferrer() string
	SetUserReferrer(userReferrer string) VisitorInterface
	// Memo() string
	// SetMemo(memo string) UserInterface
	// Meta(name string) string
	// SetMeta(name string, value string) error
	// Metas() (map[string]string, error)
	// SetMetas(metas map[string]string) error
	// Status() string
	// SetStatus(status string) UserInterface
	UpdatedAt() string
	UpdatedAtCarbon() carbon.Carbon
	SetUpdatedAt(updatedAt string) VisitorInterface
}
