package model

type WebResponse[T any] struct {
	Data    T             `json:"data"`
	Paging  *PageMetadata `json:"paging,omitempty"`
	Token   string        `json:"token,omitempty"`
	Summary string        `json:"summary,omitempty"`
	Errors  string        `json:"errors,omitempty"`
}

type PageMetadata struct {
	Page      int   `json:"page"`
	PerPage   int   `json:"perPage"`
	TotalItem int64 `json:"totalItems"`
	TotalPage int64 `json:"totalPages"`
}
