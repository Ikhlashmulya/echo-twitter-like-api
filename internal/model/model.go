package model

type WebResponse[T any] struct {
	Data   T       `json:"data"`
	Paging *Paging `json:"paging,omitempty"`
}

type Paging struct {
	Page  int   `json:"page,omitempty"`
	Size  int   `json:"size,omitempty"`
	Total int64 `json:"total,omitempty"`
}
