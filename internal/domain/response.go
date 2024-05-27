package domain

type Page struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

type SuccessResponse struct {
	Data any  `json:"data"`
	Meta Page `json:"meta"`
}
