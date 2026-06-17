package statsstore

import "errors"

// VisitorQueryInterface defines the interface for visitor query operations.
type VisitorQueryInterface interface {
	Validate() error

	HasCountry() bool
	Country() string
	SetCountry(country string) VisitorQueryInterface

	HasCreatedAtGte() bool
	CreatedAtGte() string
	SetCreatedAtGte(createdAtGte string) VisitorQueryInterface

	HasCreatedAtLte() bool
	CreatedAtLte() string
	SetCreatedAtLte(createdAtLte string) VisitorQueryInterface

	HasDeviceType() bool
	DeviceType() string
	SetDeviceType(deviceType string) VisitorQueryInterface

	HasDistinct() bool
	Distinct() string
	SetDistinct(distinct string) VisitorQueryInterface

	HasID() bool
	ID() string
	SetID(id string) VisitorQueryInterface

	HasIDIn() bool
	IDIn() []string
	SetIDIn(idIn []string) VisitorQueryInterface

	HasLimit() bool
	Limit() int
	SetLimit(limit int) VisitorQueryInterface

	HasOffset() bool
	Offset() int
	SetOffset(offset int) VisitorQueryInterface

	HasOrderBy() bool
	OrderBy() string
	SetOrderBy(orderBy string) VisitorQueryInterface

	HasPathContains() bool
	PathContains() string
	SetPathContains(pathContains string) VisitorQueryInterface

	HasPathExact() bool
	PathExact() string
	SetPathExact(pathExact string) VisitorQueryInterface

	HasSortOrder() bool
	SortOrder() string
	SetSortOrder(sortOrder string) VisitorQueryInterface

	HasSoftDeletedIncluded() bool
	SoftDeletedIncluded() bool
	SetSoftDeletedIncluded(withSoftDeleted bool) VisitorQueryInterface
}

// VisitorQuery is a shortcut for NewVisitorQuery.
func VisitorQuery() VisitorQueryInterface {
	return NewVisitorQuery()
}

// NewVisitorQuery creates a new visitor query.
func NewVisitorQuery() VisitorQueryInterface {
	return &visitorQuery{
		properties: make(map[string]interface{}),
	}
}

var _ VisitorQueryInterface = (*visitorQuery)(nil)

type visitorQuery struct {
	properties map[string]interface{}
}

func (q *visitorQuery) Validate() error {
	if q.HasCreatedAtGte() && q.CreatedAtGte() == "" {
		return errors.New("visitor query: created_at_gte cannot be empty")
	}
	if q.HasCreatedAtLte() && q.CreatedAtLte() == "" {
		return errors.New("visitor query: created_at_lte cannot be empty")
	}
	if q.HasID() && q.ID() == "" {
		return errors.New("visitor query: id cannot be empty")
	}
	if q.HasIDIn() && len(q.IDIn()) < 1 {
		return errors.New("visitor query: id_in cannot be empty array")
	}
	if q.HasLimit() && q.Limit() < 0 {
		return errors.New("visitor query: limit cannot be negative")
	}
	if q.HasOffset() && q.Offset() < 0 {
		return errors.New("visitor query: offset cannot be negative")
	}
	return nil
}

func (q *visitorQuery) hasProperty(key string) bool {
	_, ok := q.properties[key]
	return ok
}

func (q *visitorQuery) HasCountry() bool { return q.hasProperty("country") }
func (q *visitorQuery) Country() string  { return q.properties["country"].(string) }
func (q *visitorQuery) SetCountry(v string) VisitorQueryInterface {
	q.properties["country"] = v
	return q
}

func (q *visitorQuery) HasCreatedAtGte() bool { return q.hasProperty("created_at_gte") }
func (q *visitorQuery) CreatedAtGte() string  { return q.properties["created_at_gte"].(string) }
func (q *visitorQuery) SetCreatedAtGte(v string) VisitorQueryInterface {
	q.properties["created_at_gte"] = v
	return q
}

func (q *visitorQuery) HasCreatedAtLte() bool { return q.hasProperty("created_at_lte") }
func (q *visitorQuery) CreatedAtLte() string  { return q.properties["created_at_lte"].(string) }
func (q *visitorQuery) SetCreatedAtLte(v string) VisitorQueryInterface {
	q.properties["created_at_lte"] = v
	return q
}

func (q *visitorQuery) HasDeviceType() bool { return q.hasProperty("device_type") }
func (q *visitorQuery) DeviceType() string  { return q.properties["device_type"].(string) }
func (q *visitorQuery) SetDeviceType(v string) VisitorQueryInterface {
	q.properties["device_type"] = v
	return q
}

func (q *visitorQuery) HasDistinct() bool { return q.hasProperty("distinct") }
func (q *visitorQuery) Distinct() string  { return q.properties["distinct"].(string) }
func (q *visitorQuery) SetDistinct(v string) VisitorQueryInterface {
	q.properties["distinct"] = v
	return q
}

func (q *visitorQuery) HasID() bool { return q.hasProperty("id") }
func (q *visitorQuery) ID() string  { return q.properties["id"].(string) }
func (q *visitorQuery) SetID(id string) VisitorQueryInterface {
	q.properties["id"] = id
	return q
}

func (q *visitorQuery) HasIDIn() bool  { return q.hasProperty("id_in") }
func (q *visitorQuery) IDIn() []string { return q.properties["id_in"].([]string) }
func (q *visitorQuery) SetIDIn(v []string) VisitorQueryInterface {
	q.properties["id_in"] = v
	return q
}

func (q *visitorQuery) HasLimit() bool { return q.hasProperty("limit") }
func (q *visitorQuery) Limit() int     { return q.properties["limit"].(int) }
func (q *visitorQuery) SetLimit(v int) VisitorQueryInterface {
	q.properties["limit"] = v
	return q
}

func (q *visitorQuery) HasOffset() bool { return q.hasProperty("offset") }
func (q *visitorQuery) Offset() int     { return q.properties["offset"].(int) }
func (q *visitorQuery) SetOffset(v int) VisitorQueryInterface {
	q.properties["offset"] = v
	return q
}

func (q *visitorQuery) HasOrderBy() bool { return q.hasProperty("order_by") }
func (q *visitorQuery) OrderBy() string  { return q.properties["order_by"].(string) }
func (q *visitorQuery) SetOrderBy(v string) VisitorQueryInterface {
	q.properties["order_by"] = v
	return q
}

func (q *visitorQuery) HasPathContains() bool { return q.hasProperty("path_contains") }
func (q *visitorQuery) PathContains() string  { return q.properties["path_contains"].(string) }
func (q *visitorQuery) SetPathContains(v string) VisitorQueryInterface {
	q.properties["path_contains"] = v
	return q
}

func (q *visitorQuery) HasPathExact() bool { return q.hasProperty("path_exact") }
func (q *visitorQuery) PathExact() string  { return q.properties["path_exact"].(string) }
func (q *visitorQuery) SetPathExact(v string) VisitorQueryInterface {
	q.properties["path_exact"] = v
	return q
}

func (q *visitorQuery) HasSortOrder() bool { return q.hasProperty("sort_order") }
func (q *visitorQuery) SortOrder() string  { return q.properties["sort_order"].(string) }
func (q *visitorQuery) SetSortOrder(v string) VisitorQueryInterface {
	q.properties["sort_order"] = v
	return q
}

func (q *visitorQuery) HasSoftDeletedIncluded() bool { return q.hasProperty("soft_deleted_included") }
func (q *visitorQuery) SoftDeletedIncluded() bool {
	if !q.HasSoftDeletedIncluded() {
		return false
	}
	return q.properties["soft_deleted_included"].(bool)
}
func (q *visitorQuery) SetSoftDeletedIncluded(v bool) VisitorQueryInterface {
	q.properties["soft_deleted_included"] = v
	return q
}
