package utils

type QueryParams struct {
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	Sort   string `json:"sort"`
	Filter string `json:"filter"`
}
