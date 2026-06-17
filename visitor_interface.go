package statsstore

import "github.com/dromara/carbon/v2"

// VisitorInterface defines the interface for a visitor record.
type VisitorInterface interface {
	// Methods
	FingerprintCalculate() string
	IsSoftDeleted() bool

	// Setters and Getters

	GetID() string
	SetID(id string) VisitorInterface

	GetPath() string
	SetPath(path string) VisitorInterface

	GetCountry() string
	SetCountry(country string) VisitorInterface

	GetCreatedAt() string
	GetCreatedAtCarbon() *carbon.Carbon
	SetCreatedAt(createdAt string) VisitorInterface

	GetSoftDeletedAt() string
	GetSoftDeletedAtCarbon() *carbon.Carbon
	SetSoftDeletedAt(deletedAt string) VisitorInterface

	GetFingerprint() string
	SetFingerprint(fingerprint string) VisitorInterface

	GetIpAddress() string
	SetIpAddress(ipAddress string) VisitorInterface

	GetUpdatedAt() string
	GetUpdatedAtCarbon() *carbon.Carbon
	SetUpdatedAt(updatedAt string) VisitorInterface

	GetUserAcceptLanguage() string
	SetUserAcceptLanguage(userAcceptLanguage string) VisitorInterface

	GetUserAcceptEncoding() string
	SetUserAcceptEncoding(userAcceptEncoding string) VisitorInterface

	GetUserAgent() string
	SetUserAgent(userAgent string) VisitorInterface

	GetUserBrowser() string
	SetUserBrowser(userBrowser string) VisitorInterface

	GetUserBrowserVersion() string
	SetUserBrowserVersion(userBrowserVersion string) VisitorInterface

	GetUserDevice() string
	SetUserDevice(userDevice string) VisitorInterface

	GetUserDeviceType() string
	SetUserDeviceType(userDeviceType string) VisitorInterface

	GetUserOs() string
	SetUserOs(userOs string) VisitorInterface

	GetUserOsVersion() string
	SetUserOsVersion(userOsVersion string) VisitorInterface

	GetUserReferrer() string
	SetUserReferrer(userReferrer string) VisitorInterface
}
