package store

import (
	"math"
	"strings"
)

type Filters struct {
	Front string `validate:"max=300"`
	Sort string `validate:"oneof=due -due last_review -last_review created_at -created_at id"`
	Page int64  `validate:"min=1,max=10000000"`
	PageSize int64  `validate:"min=1,max=100"`
}

func (f Filters) sortColumn() string {
	return strings.TrimPrefix(f.Sort, "-")
}

func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}

	return "ASC"
}

func (f Filters) offset() int64 {
	return (f.Page - 1) * f.PageSize
}

type Metadata struct {
	CurrentPage int64 `json:"current_page,omitempty"`
	PageSize int64 `json:"page_size,omitempty"`
	FirstPage int64 `json:"first_page,omitempty"` 
	LastPage int64  `json:"last_page,omitempty"`
	TotalRecords int64 `json:"total_records,omitempty"`
}

func calculateMetadata(totalRecords, page, pageSize int64) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage: page,
		PageSize: pageSize,
		FirstPage: 1,
		LastPage: int64(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}


