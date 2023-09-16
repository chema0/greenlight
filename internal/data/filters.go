package data

import (
	"strings"

	"github.com/chema0/greenlight/internal/validator"
)

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafeList []string
}

func (f Filters) sortColumn() string {
	for _, safeValue := range f.SortSafeList {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}

	panic("unsafe sort parameter: " + f.Sort)
}

func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}

	return "ASC"
}

func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(v.InRange(f.Page, 0, 10_000_000), "page", "must be greater than zero and less or equal than 10_000_000")
	v.Check(v.InRange(f.PageSize, 0, 100), "page", "must be greater than zero and less or equal than 100")

	v.Check(validator.PermittedValue(f.Sort, f.SortSafeList...), "sort", "invalid sort value")
}
