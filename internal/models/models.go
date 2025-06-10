package models

type ModelResponse struct {
	ID      string `json:"id" validate:"required"`
	Created int64  `json:"created" validate:"required"`
	OwnedBy string `json:"owned_by" validate:"required"`
	Object  string `json:"object" validate:"required,oneof=model"`
}

type ModelListResponse struct {
	Type string          `json:"type" validate:"required,oneof=list"`
	Data []ModelResponse `json:"data" validate:"required,dive"`
}
