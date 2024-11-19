package model

type ResponseStatus string

const (
	StatusSuccess ResponseStatus = "success"
	StatusError   ResponseStatus = "error"
)

type StandardResponse struct {
	Status  ResponseStatus `json:"status"`
	Message string         `json:"message,omitempty"`
	Data    interface{}    `json:"data,omitempty"`
	Errors  interface{}    `json:"errors,omitempty"`
}
