package model

type Page[V any] struct {
	Page    int   `json:"page"`
	Count   int64 `json:"count"`
	Content []*V  `json:"content"`
}
