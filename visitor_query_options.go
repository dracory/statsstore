package statsstore

type VisitorQueryOptions struct {
	ID   string
	IDIn []string
	// Status       string
	// StatusIn     []string
	Distinct     string // distinct select column
	Country      string
	PathContains string
	PathExact    string
	DeviceType   string
	CreatedAtGte string
	CreatedAtLte string
	Offset       int
	Limit        int
	SortOrder    string
	OrderBy      string
	CountOnly    bool
	WithDeleted  bool
}
