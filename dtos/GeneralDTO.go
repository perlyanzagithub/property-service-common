package dtos

type RequestDTO struct {
	Size    int                    `json:"size"`
	Page    int                    `json:"page"`
	Filter  map[string]interface{} `json:"filter"`
	OrderBy string                 `json:"order_by"`
}
