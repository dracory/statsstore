package statsstore

import (
	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/dataobject"
	"github.com/gouniverse/sb"
	"github.com/gouniverse/uid"
	"github.com/gouniverse/utils"
)

type Visitor struct {
	dataobject.DataObject
}

var _ VisitorInterface = (*Visitor)(nil)

func NewVisitor() VisitorInterface {
	o := (&Visitor{}).
		SetID(uid.HumanUid()).
		SetPath("").
		SetCountry("").
		SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetIpAddress("").
		SetFingerprint("").
		SetUserAcceptEncoding("").
		SetUserAcceptLanguage("").
		SetUserAgent("").
		SetUserBrowser("").
		SetUserBrowserVersion("").
		SetUserDevice("").
		SetUserDeviceType("").
		SetUserOs("").
		SetUserOsVersion("").
		SetUserReferrer("").
		SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetDeletedAt(sb.MAX_DATETIME)
	return o
}

func NewVisitorFromExistingData(data map[string]string) VisitorInterface {
	o := &Visitor{}
	o.Hydrate(data)
	return o
}

func (visitor *Visitor) FingerprintCalculate() string {
	fingerprint := visitor.IpAddress() + visitor.UserAgent()
	hash := utils.StrToMD5Hash(fingerprint)
	return hash
}

func (visitor *Visitor) Country() string {
	return visitor.Get(COLUMN_COUNTRY)
}

func (visitor *Visitor) SetCountry(country string) VisitorInterface {
	visitor.Set(COLUMN_COUNTRY, country)
	return visitor
}

func (visitor *Visitor) CreatedAt() string {
	return visitor.Get(COLUMN_CREATED_AT)
}

func (visitor *Visitor) CreatedAtCarbon() carbon.Carbon {
	return carbon.Parse(visitor.CreatedAt(), carbon.UTC)
}

func (visitor *Visitor) SetCreatedAt(createdAt string) VisitorInterface {
	visitor.Set(COLUMN_CREATED_AT, createdAt)
	return visitor
}

func (visitor *Visitor) DeletedAt() string {
	return visitor.Get(COLUMN_DELETED_AT)
}

func (visitor *Visitor) DeletedAtCarbon() carbon.Carbon {
	return carbon.Parse(visitor.DeletedAt(), carbon.UTC)
}

func (visitor *Visitor) SetDeletedAt(deletedAt string) VisitorInterface {
	visitor.Set(COLUMN_DELETED_AT, deletedAt)
	return visitor
}

func (visitor *Visitor) Fingerprint() string {
	return visitor.Get(COLUMN_FINGERPRINT)
}

func (visitor *Visitor) SetFingerprint(fingerprint string) VisitorInterface {
	visitor.Set(COLUMN_FINGERPRINT, fingerprint)
	return visitor
}

func (visitor *Visitor) ID() string {
	return visitor.Get(COLUMN_ID)
}

func (visitor *Visitor) SetID(id string) VisitorInterface {
	visitor.Set(COLUMN_ID, id)
	return visitor
}

func (visitor *Visitor) IpAddress() string {
	return visitor.Get(COLUMN_IP_ADDRESS)
}

func (visitor *Visitor) SetIpAddress(ipAddress string) VisitorInterface {
	visitor.Set(COLUMN_IP_ADDRESS, ipAddress)
	return visitor
}

func (visitor *Visitor) Path() string {
	return visitor.Get(COLUMN_PATH)
}

func (visitor *Visitor) SetPath(path string) VisitorInterface {
	visitor.Set(COLUMN_PATH, path)
	return visitor
}

func (visitor *Visitor) UpdatedAt() string {
	return visitor.Get(COLUMN_UPDATED_AT)
}

func (visitor *Visitor) UpdatedAtCarbon() carbon.Carbon {
	return carbon.Parse(visitor.UpdatedAt(), carbon.UTC)
}

func (visitor *Visitor) SetUpdatedAt(updatedAt string) VisitorInterface {
	visitor.Set(COLUMN_UPDATED_AT, updatedAt)
	return visitor
}

func (visitor *Visitor) UserAcceptLanguage() string {
	return visitor.Get(COLUMN_USER_ACCEPT_LANGUAGE)
}

func (visitor *Visitor) SetUserAcceptLanguage(userAcceptLanguage string) VisitorInterface {
	visitor.Set(COLUMN_USER_ACCEPT_LANGUAGE, userAcceptLanguage)
	return visitor
}

func (visitor *Visitor) UserAcceptEncoding() string {
	return visitor.Get(COLUMN_USER_ACCEPT_ENCODING)
}

func (visitor *Visitor) SetUserAcceptEncoding(userAcceptEncoding string) VisitorInterface {
	visitor.Set(COLUMN_USER_ACCEPT_ENCODING, userAcceptEncoding)
	return visitor
}

func (visitor *Visitor) UserAgent() string {
	return visitor.Get(COLUMN_USER_AGENT)
}

func (visitor *Visitor) SetUserAgent(userAgent string) VisitorInterface {
	visitor.Set(COLUMN_USER_AGENT, userAgent)
	return visitor
}

func (visitor *Visitor) UserBrowser() string {
	return visitor.Get(COLUMN_USER_BROWSER)
}

func (visitor *Visitor) SetUserBrowser(userBrowser string) VisitorInterface {
	visitor.Set(COLUMN_USER_BROWSER, userBrowser)
	return visitor
}

func (visitor *Visitor) UserBrowserVersion() string {
	return visitor.Get(COLUMN_USER_BROWSER_VERSION)
}

func (visitor *Visitor) SetUserBrowserVersion(userBrowserVersion string) VisitorInterface {
	visitor.Set(COLUMN_USER_BROWSER_VERSION, userBrowserVersion)
	return visitor
}

func (visitor *Visitor) UserDevice() string {
	return visitor.Get(COLUMN_USER_DEVICE)
}

func (visitor *Visitor) SetUserDevice(userDevice string) VisitorInterface {
	visitor.Set(COLUMN_USER_DEVICE, userDevice)
	return visitor
}

func (visitor *Visitor) UserDeviceType() string {
	return visitor.Get(COLUMN_USER_DEVICE_TYPE)
}

func (visitor *Visitor) SetUserDeviceType(userDeviceType string) VisitorInterface {
	visitor.Set(COLUMN_USER_DEVICE_TYPE, userDeviceType)
	return visitor
}

func (visitor *Visitor) UserOs() string {
	return visitor.Get(COLUMN_USER_OS)
}

func (visitor *Visitor) SetUserOs(userOs string) VisitorInterface {
	visitor.Set(COLUMN_USER_OS, userOs)
	return visitor
}

func (visitor *Visitor) UserOsVersion() string {
	return visitor.Get(COLUMN_USER_OS_VERSION)
}

func (visitor *Visitor) SetUserOsVersion(userOsVersion string) VisitorInterface {
	visitor.Set(COLUMN_USER_OS_VERSION, userOsVersion)
	return visitor
}

func (visitor *Visitor) UserReferrer() string {
	return visitor.Get(COLUMN_USER_REFERRER)
}

func (visitor *Visitor) SetUserReferrer(userReferrer string) VisitorInterface {
	visitor.Set(COLUMN_USER_REFERRER, userReferrer)
	return visitor
}

// type NewVistorParameters struct {
// 	Id                 string
// 	UrlID              string
// 	IPAddress          string
// 	Country            string
// 	UserAgent          string
// 	UserAcceptLanguage string
// 	UserBrowser        string
// 	UserBrowserVersion string
// 	UserDevice         string
// 	UserDeviceType     string
// 	UserOs             string
// 	UserOsVersion      string
// 	CreatedAt          time.Time
// 	UpdatedAt          time.Time
// 	DeletedAt          *time.Time
// }

// // NewPrestagedShortID create new from map
// func (st *Store) NewVisitor(params NewVistorParameters) Visitor {
// 	visitor := Visitor{}
// 	visitor.SetID(utils.ToString(params.Id))
// 	visitor.SetUrlID(params.UrlID)
// 	visitor.SetCountry(params.Country)
// 	visitor.SetIPAddress(params.IPAddress)
// 	visitor.SetUserAcceptLanguage(params.UserAcceptLanguage)
// 	visitor.SetUserAgent(params.UserAgent)
// 	visitor.SetUserBrowser(params.UserBrowser)
// 	visitor.SetUserBrowserVersion(params.UserBrowserVersion)
// 	visitor.SetUserDevice(params.UserDevice)
// 	visitor.SetUserDeviceType(params.UserDeviceType)
// 	visitor.SetUserOs(params.UserOs)
// 	visitor.SetUserOsVersion(params.UserOsVersion)
// 	visitor.SetCreatedAt(params.CreatedAt)
// 	visitor.SetUpdatedAt(params.UpdatedAt)
// 	if params.DeletedAt != nil {
// 		visitor.SetDeletedAt(*params.DeletedAt)
// 	}
// 	//url.SetCreatedAt(carbon.Parse(utils.ToString(m["created_at"]), carbon.UTC).ToStdTime())
// 	//url.SetUpdatedAt(carbon.Parse(utils.ToString(m["updated_at"]), carbon.UTC).ToStdTime())
// 	//url.SetDeletedAt(carbon.Parse(utils.ToString(m["deleted_at"]), carbon.UTC).ToStdTime())
// 	return visitor
// }

// func (v *Visitor) ToMap() map[string]any {
// 	entry := map[string]interface{}{}
// 	entry["id"] = v.ID()
// 	entry["url_id"] = v.UrlID()
// 	entry["country"] = v.Country()
// 	entry["ip_address"] = v.IPAddress()
// 	entry["user_accept_language"] = v.UserAcceptLanguage()
// 	entry["user_agent"] = v.UserAgent()
// 	entry["user_browser"] = v.UserBrowser()
// 	entry["user_browser_version"] = v.UserBrowserVersion()
// 	entry["user_device"] = v.UserDevice()
// 	entry["user_device_type"] = v.UserDeviceType()
// 	entry["user_os"] = v.UserOs()
// 	entry["user_os_version"] = v.UserOsVersion()
// 	entry["created_at"] = v.CreatedAt()
// 	entry["updated_at"] = v.UpdatedAt()
// 	entry["deleted_at"] = v.DeletedAt()
// 	return entry
// }

// func (v *Visitor) ID() string {
// 	return v.id
// }

// func (v *Visitor) Country() string {
// 	return v.country
// }

// func (v *Visitor) IPAddress() string {
// 	return v.ipAddress
// }

// func (v *Visitor) UrlID() string {
// 	return v.urlID
// }

// func (v *Visitor) UserAcceptLanguage() string {
// 	return v.userAcceptLanguage
// }

// func (v *Visitor) UserAgent() string {
// 	return v.userAgent
// }

// func (v *Visitor) UserBrowser() string {
// 	return v.userBrowser
// }

// func (v *Visitor) UserBrowserVersion() string {
// 	return v.userBrowserVersion
// }

// func (v *Visitor) UserDevice() string {
// 	return v.userDevice
// }

// func (v *Visitor) UserDeviceType() string {
// 	return v.userDeviceType
// }

// func (v *Visitor) UserOs() string {
// 	return v.userOs
// }

// func (v *Visitor) UserOsVersion() string {
// 	return v.userOsVersion
// }

// func (v *Visitor) CreatedAt() time.Time {
// 	return v.createdAt
// }

// func (v *Visitor) CreatedAtCarbon() carbon.Carbon {
// 	return carbon.CreateFromStdTime(v.createdAt, carbon.UTC)
// }

// func (v *Visitor) UpdatedAt() time.Time {
// 	return v.updatedAt
// }

// func (v *Visitor) DeletedAt() *time.Time {
// 	return v.deletedAt
// }

// func (v *Visitor) SetID(id string) *Visitor {
// 	v.id = id
// 	return v
// }

// func (v *Visitor) SetCountry(country string) *Visitor {
// 	v.country = country
// 	return v
// }

// func (v *Visitor) SetIPAddress(ipAddress string) *Visitor {
// 	v.ipAddress = ipAddress
// 	return v
// }

// func (v *Visitor) SetUrlID(urlID string) *Visitor {
// 	v.urlID = urlID
// 	return v
// }

// func (v *Visitor) SetUserAcceptLanguage(userAcceptLanguage string) *Visitor {
// 	v.userAcceptLanguage = userAcceptLanguage
// 	return v
// }

// func (v *Visitor) SetUserAgent(userAgent string) *Visitor {
// 	v.userAgent = userAgent
// 	return v
// }

// func (v *Visitor) SetUserBrowser(userBrowser string) *Visitor {
// 	v.userBrowser = userBrowser
// 	return v
// }

// func (v *Visitor) SetUserBrowserVersion(userBrowserVersion string) *Visitor {
// 	v.userBrowserVersion = userBrowserVersion
// 	return v
// }

// func (v *Visitor) SetUserDevice(userDevice string) *Visitor {
// 	v.userDevice = userDevice
// 	return v
// }

// func (v *Visitor) SetUserDeviceType(userDeviceType string) *Visitor {
// 	v.userDeviceType = userDeviceType
// 	return v
// }

// func (v *Visitor) SetUserOs(userOs string) *Visitor {
// 	v.userOs = userOs
// 	return v
// }

// func (v *Visitor) SetUserOsVersion(userOsVersion string) *Visitor {
// 	v.userOsVersion = userOsVersion
// 	return v
// }

// func (v *Visitor) SetCreatedAt(createdAt time.Time) *Visitor {
// 	v.createdAt = createdAt
// 	return v
// }

// func (v *Visitor) SetUpdatedAt(updatedAt time.Time) *Visitor {
// 	v.updatedAt = updatedAt
// 	return v
// }

// func (v *Visitor) SetDeletedAt(deletedAt time.Time) *Visitor {
// 	v.deletedAt = &deletedAt
// 	return v
// }
