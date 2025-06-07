package data

import (
	"math"
	"strings"

	"github.com/VladimirArtyom/rest_eiga_api/internal/validator"
)

type Filters struct {
	Page              int
	PageSize          int
	Sort              string
	SupportedSortList []string
}

type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

func ValidateFilters(v *validator.Validator, f Filters) {

	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "a maximum page is 10 million")

	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "a maximum page is 100")

	v.Check(validator.In(f.Sort, f.SupportedSortList...), "sort", "invalid sort value")
}

func calculateMetadata(totalRecords int, page int, pageSize int) Metadata {

	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		FirstPage:    1,
		CurrentPage:  page,
		PageSize:     pageSize,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}

}

func (f Filters) sortColumn() string {
	for _, safeValue := range f.SupportedSortList {
		if f.Sort == safeValue {
			return strings.TrimPrefix(safeValue, "-")
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

func (f Filters) limit() int {
	return f.PageSize
}

func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}
