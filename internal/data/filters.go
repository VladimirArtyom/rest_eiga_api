package data

import "github.com/VladimirArtyom/rest_eiga_api/internal/validator"

type Filters struct {
	Page              int
	PageSize          int
	Sort              string
	SupportedSortList []string
}

func ValidateFilters(v *validator.Validator, f Filters) {

	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "a maximum page is 10 million")

	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "a maximum page is 100")

	v.Check(validator.In(f.Sort, f.SupportedSortList...), "sort", "invalid sort value")
}
