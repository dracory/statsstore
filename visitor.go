package statsstore

import (
	"github.com/dracory/neat/database/orm"
	"github.com/dracory/neat/database/soft_delete"
	neatuid "github.com/dracory/neat/support/uid"
	"github.com/dracory/str"
	"github.com/dromara/carbon/v2"
)

// == TYPE =====================================================================

type visitorImplementation struct {
	orm.ShortID

	PathField               string `db:"path"`
	FingerprintField        string `db:"fingerprint"`
	IPAddressField          string `db:"ip_address"`
	CountryField            string `db:"country"`
	UserAcceptLanguageField string `db:"user_accept_language"`
	UserAcceptEncodingField string `db:"user_accept_encoding"`
	UserAgentField          string `db:"user_agent"`
	UserOsField             string `db:"user_os"`
	UserOsVersionField      string `db:"user_os_version"`
	UserDeviceField         string `db:"user_device"`
	UserDeviceTypeField     string `db:"user_device_type"`
	UserBrowserField        string `db:"user_browser"`
	UserBrowserVersionField string `db:"user_browser_version"`
	UserReferrerField       string `db:"user_referrer"`
	CreatedAtField          orm.CreatedAt
	UpdatedAtField          orm.UpdatedAt
	soft_delete.SoftDeletesMaxDate
}

var _ VisitorInterface = (*visitorImplementation)(nil)

// == CONSTRUCTORS =============================================================

// NewVisitor creates a new visitor.
func NewVisitor() VisitorInterface {
	o := &visitorImplementation{}
	o.SetID(neatuid.GenerateShortID())
	o.SetPath("")
	o.SetCountry("")
	o.SetIpAddress("")
	o.SetFingerprint("")
	o.SetUserAcceptEncoding("")
	o.SetUserAcceptLanguage("")
	o.SetUserAgent("")
	o.SetUserBrowser("")
	o.SetUserBrowserVersion("")
	o.SetUserDevice("")
	o.SetUserDeviceType("")
	o.SetUserOs("")
	o.SetUserOsVersion("")
	o.SetUserReferrer("")
	o.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetSoftDeletedAt(MAX_DATETIME)
	return o
}

// NewVisitorFromExistingData creates a new visitor from a raw column map.
func NewVisitorFromExistingData(data map[string]string) VisitorInterface {
	o := &visitorImplementation{}
	if v, ok := data[COLUMN_ID]; ok {
		o.SetID(v)
	}
	if v, ok := data[COLUMN_PATH]; ok {
		o.SetPath(v)
	}
	if v, ok := data[COLUMN_FINGERPRINT]; ok {
		o.SetFingerprint(v)
	}
	if v, ok := data[COLUMN_IP_ADDRESS]; ok {
		o.SetIpAddress(v)
	}
	if v, ok := data[COLUMN_COUNTRY]; ok {
		o.SetCountry(v)
	}
	if v, ok := data[COLUMN_USER_ACCEPT_LANGUAGE]; ok {
		o.SetUserAcceptLanguage(v)
	}
	if v, ok := data[COLUMN_USER_ACCEPT_ENCODING]; ok {
		o.SetUserAcceptEncoding(v)
	}
	if v, ok := data[COLUMN_USER_AGENT]; ok {
		o.SetUserAgent(v)
	}
	if v, ok := data[COLUMN_USER_OS]; ok {
		o.SetUserOs(v)
	}
	if v, ok := data[COLUMN_USER_OS_VERSION]; ok {
		o.SetUserOsVersion(v)
	}
	if v, ok := data[COLUMN_USER_DEVICE]; ok {
		o.SetUserDevice(v)
	}
	if v, ok := data[COLUMN_USER_DEVICE_TYPE]; ok {
		o.SetUserDeviceType(v)
	}
	if v, ok := data[COLUMN_USER_BROWSER]; ok {
		o.SetUserBrowser(v)
	}
	if v, ok := data[COLUMN_USER_BROWSER_VERSION]; ok {
		o.SetUserBrowserVersion(v)
	}
	if v, ok := data[COLUMN_USER_REFERRER]; ok {
		o.SetUserReferrer(v)
	}
	if v, ok := data[COLUMN_CREATED_AT]; ok {
		o.SetCreatedAt(v)
	}
	if v, ok := data[COLUMN_UPDATED_AT]; ok {
		o.SetUpdatedAt(v)
	}
	if v, ok := data[COLUMN_SOFT_DELETED_AT]; ok {
		o.SetSoftDeletedAt(v)
	}
	return o
}

// == METHODS ==================================================================

// FingerprintCalculate calculates a fingerprint from IP and UserAgent.
func (o *visitorImplementation) FingerprintCalculate() string {
	fingerprint := o.IPAddressField + o.UserAgentField
	hash := str.MD5(fingerprint)
	return hash
}

// IsSoftDeleted returns true if the visitor is soft deleted.
func (o *visitorImplementation) IsSoftDeleted() bool {
	return o.SoftDeletesMaxDate.IsSoftDeleted()
}

// == SETTERS AND GETTERS ======================================================

// GetID returns the id of the visitor.
func (o *visitorImplementation) GetID() string {
	return o.ShortID.ID
}

// SetID sets the id of the visitor.
func (o *visitorImplementation) SetID(id string) VisitorInterface {
	o.ShortID.ID = id
	return o
}

// GetPath returns the path of the visitor.
func (o *visitorImplementation) GetPath() string {
	return o.PathField
}

// SetPath sets the path of the visitor.
func (o *visitorImplementation) SetPath(path string) VisitorInterface {
	o.PathField = path
	return o
}

// GetCountry returns the country of the visitor.
func (o *visitorImplementation) GetCountry() string {
	return o.CountryField
}

// SetCountry sets the country of the visitor.
func (o *visitorImplementation) SetCountry(country string) VisitorInterface {
	o.CountryField = country
	return o
}

// GetCreatedAt returns the created at time of the visitor.
func (o *visitorImplementation) GetCreatedAt() string {
	if o.CreatedAtField.CreatedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt).ToDateTimeString()
}

// GetCreatedAtCarbon returns the created at time of the visitor as a carbon object.
func (o *visitorImplementation) GetCreatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt)
}

// SetCreatedAt sets the created at time of the visitor.
func (o *visitorImplementation) SetCreatedAt(createdAt string) VisitorInterface {
	if createdAt == "" {
		return o
	}
	o.CreatedAtField.CreatedAt = carbon.Parse(createdAt, carbon.UTC).StdTime()
	return o
}

// GetSoftDeletedAt returns the soft deleted at time of the visitor.
func (o *visitorImplementation) GetSoftDeletedAt() string {
	if o.SoftDeletesMaxDate.SoftDeletedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.SoftDeletesMaxDate.SoftDeletedAt).ToDateTimeString()
}

// GetSoftDeletedAtCarbon returns the soft deleted at time of the visitor as a carbon object.
func (o *visitorImplementation) GetSoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.SoftDeletesMaxDate.SoftDeletedAt)
}

// SetSoftDeletedAt sets the soft deleted at time of the visitor.
func (o *visitorImplementation) SetSoftDeletedAt(deletedAt string) VisitorInterface {
	if deletedAt == "" {
		return o
	}
	o.SoftDeletesMaxDate.SoftDeletedAt = carbon.Parse(deletedAt, carbon.UTC).StdTime()
	return o
}

// GetFingerprint returns the fingerprint of the visitor.
func (o *visitorImplementation) GetFingerprint() string {
	return o.FingerprintField
}

// SetFingerprint sets the fingerprint of the visitor.
func (o *visitorImplementation) SetFingerprint(fingerprint string) VisitorInterface {
	o.FingerprintField = fingerprint
	return o
}

// GetIpAddress returns the IP address of the visitor.
func (o *visitorImplementation) GetIpAddress() string {
	return o.IPAddressField
}

// SetIpAddress sets the IP address of the visitor.
func (o *visitorImplementation) SetIpAddress(ipAddress string) VisitorInterface {
	o.IPAddressField = ipAddress
	return o
}

// GetUpdatedAt returns the updated at time of the visitor.
func (o *visitorImplementation) GetUpdatedAt() string {
	if o.UpdatedAtField.UpdatedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.UpdatedAtField.UpdatedAt).ToDateTimeString()
}

// GetUpdatedAtCarbon returns the updated at time of the visitor as a carbon object.
func (o *visitorImplementation) GetUpdatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.UpdatedAtField.UpdatedAt)
}

// SetUpdatedAt sets the updated at time of the visitor.
func (o *visitorImplementation) SetUpdatedAt(updatedAt string) VisitorInterface {
	if updatedAt == "" {
		return o
	}
	o.UpdatedAtField.UpdatedAt = carbon.Parse(updatedAt, carbon.UTC).StdTime()
	return o
}

// GetUserAcceptLanguage returns the user accept language of the visitor.
func (o *visitorImplementation) GetUserAcceptLanguage() string {
	return o.UserAcceptLanguageField
}

// SetUserAcceptLanguage sets the user accept language of the visitor.
func (o *visitorImplementation) SetUserAcceptLanguage(userAcceptLanguage string) VisitorInterface {
	o.UserAcceptLanguageField = userAcceptLanguage
	return o
}

// GetUserAcceptEncoding returns the user accept encoding of the visitor.
func (o *visitorImplementation) GetUserAcceptEncoding() string {
	return o.UserAcceptEncodingField
}

// SetUserAcceptEncoding sets the user accept encoding of the visitor.
func (o *visitorImplementation) SetUserAcceptEncoding(userAcceptEncoding string) VisitorInterface {
	o.UserAcceptEncodingField = userAcceptEncoding
	return o
}

// GetUserAgent returns the user agent of the visitor.
func (o *visitorImplementation) GetUserAgent() string {
	return o.UserAgentField
}

// SetUserAgent sets the user agent of the visitor.
func (o *visitorImplementation) SetUserAgent(userAgent string) VisitorInterface {
	o.UserAgentField = userAgent
	return o
}

// GetUserBrowser returns the user browser of the visitor.
func (o *visitorImplementation) GetUserBrowser() string {
	return o.UserBrowserField
}

// SetUserBrowser sets the user browser of the visitor.
func (o *visitorImplementation) SetUserBrowser(userBrowser string) VisitorInterface {
	o.UserBrowserField = userBrowser
	return o
}

// GetUserBrowserVersion returns the user browser version of the visitor.
func (o *visitorImplementation) GetUserBrowserVersion() string {
	return o.UserBrowserVersionField
}

// SetUserBrowserVersion sets the user browser version of the visitor.
func (o *visitorImplementation) SetUserBrowserVersion(userBrowserVersion string) VisitorInterface {
	o.UserBrowserVersionField = userBrowserVersion
	return o
}

// GetUserDevice returns the user device of the visitor.
func (o *visitorImplementation) GetUserDevice() string {
	return o.UserDeviceField
}

// SetUserDevice sets the user device of the visitor.
func (o *visitorImplementation) SetUserDevice(userDevice string) VisitorInterface {
	o.UserDeviceField = userDevice
	return o
}

// GetUserDeviceType returns the user device type of the visitor.
func (o *visitorImplementation) GetUserDeviceType() string {
	return o.UserDeviceTypeField
}

// SetUserDeviceType sets the user device type of the visitor.
func (o *visitorImplementation) SetUserDeviceType(userDeviceType string) VisitorInterface {
	o.UserDeviceTypeField = userDeviceType
	return o
}

// GetUserOs returns the user OS of the visitor.
func (o *visitorImplementation) GetUserOs() string {
	return o.UserOsField
}

// SetUserOs sets the user OS of the visitor.
func (o *visitorImplementation) SetUserOs(userOs string) VisitorInterface {
	o.UserOsField = userOs
	return o
}

// GetUserOsVersion returns the user OS version of the visitor.
func (o *visitorImplementation) GetUserOsVersion() string {
	return o.UserOsVersionField
}

// SetUserOsVersion sets the user OS version of the visitor.
func (o *visitorImplementation) SetUserOsVersion(userOsVersion string) VisitorInterface {
	o.UserOsVersionField = userOsVersion
	return o
}

// GetUserReferrer returns the user referrer of the visitor.
func (o *visitorImplementation) GetUserReferrer() string {
	return o.UserReferrerField
}

// SetUserReferrer sets the user referrer of the visitor.
func (o *visitorImplementation) SetUserReferrer(userReferrer string) VisitorInterface {
	o.UserReferrerField = userReferrer
	return o
}
