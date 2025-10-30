package statsstore

import (
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
)

func (store *Store) visitorQuery(options VisitorQueryOptions) *goqu.SelectDataset {
	q := goqu.Dialect(store.dbDriverName).From(store.visitorTableName)

	if options.ID != "" {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID))
	}

	if len(options.IDIn) > 0 {
		q = q.Where(goqu.C(COLUMN_ID).In(options.IDIn))
	}

	// if options.Status != "" {
	// 	q = q.Where(goqu.C(COLUMN_STATUS).Eq(options.Status))
	// }

	// if len(options.StatusIn) > 0 {
	// 	q = q.Where(goqu.C(COLUMN_STATUS).In(options.StatusIn))
	// }

	// if options.Email != "" {
	// 	q = q.Where(goqu.C(COLUMN_EMAIL).Eq(options.Email))
	// }

	if options.Country == "empty" {
		q = q.Where(goqu.C(COLUMN_COUNTRY).Eq(""))
	} else if options.Country != "" {
		q = q.Where(goqu.C(COLUMN_COUNTRY).Eq(options.Country))
	}

	if options.PathExact != "" {
		q = q.Where(goqu.C(COLUMN_PATH).Eq(options.PathExact))
	} else if options.PathContains != "" {
		pattern := "%" + options.PathContains + "%"
		q = q.Where(goqu.C(COLUMN_PATH).ILike(pattern))
	}

	if options.DeviceType != "" {
		if strings.EqualFold(options.DeviceType, "empty") {
			q = q.Where(goqu.C(COLUMN_USER_DEVICE_TYPE).Eq(""))
		} else {
			q = q.Where(goqu.C(COLUMN_USER_DEVICE_TYPE).Eq(options.DeviceType))
		}
	}

	if options.CreatedAtGte != "" && options.CreatedAtLte != "" {
		q = q.Where(
			goqu.C(COLUMN_CREATED_AT).Gte(options.CreatedAtGte),
			goqu.C(COLUMN_CREATED_AT).Lte(options.CreatedAtLte),
		)
	} else if options.CreatedAtGte != "" {
		q = q.Where(goqu.C(COLUMN_CREATED_AT).Gte(options.CreatedAtGte))
	} else if options.CreatedAtLte != "" {
		q = q.Where(goqu.C(COLUMN_CREATED_AT).Lte(options.CreatedAtLte))
	}

	if !options.CountOnly {
		if options.Limit > 0 {
			q = q.Limit(uint(options.Limit))
		}

		if options.Offset > 0 {
			q = q.Offset(uint(options.Offset))
		}
	}

	sortOrder := sb.DESC
	if options.SortOrder != "" {
		sortOrder = options.SortOrder
	}

	if options.OrderBy != "" {
		if strings.EqualFold(sortOrder, sb.ASC) {
			q = q.Order(goqu.I(options.OrderBy).Asc())
		} else {
			q = q.Order(goqu.I(options.OrderBy).Desc())
		}
	}

	if options.WithDeleted {
		return q
	}

	softDeleted := goqu.C(COLUMN_DELETED_AT).
		Gt(carbon.Now(carbon.UTC).ToDateTimeString())

	return q.Where(softDeleted)
}
