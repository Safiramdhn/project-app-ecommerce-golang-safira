package helper

import (
	"encoding/json"
	"net/http"

	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/model"
)

type JSONResponse struct{}

// SendSuccess sends a successful JSON response
func (j *JSONResponse) SendSuccess(w http.ResponseWriter, data interface{}, message ...string) {
	response := model.StandardResponse{
		Status: model.StatusSuccess,
		Data:   data,
	}

	if len(message) > 0 {
		response.Message = message[0]
	}

	j.sendJSON(w, http.StatusOK, response)
}

// SendCreated sends a JSON response for resource creation
func (j *JSONResponse) SendCreated(w http.ResponseWriter, data interface{}, message ...string) {
	response := model.StandardResponse{
		Status: model.StatusSuccess,
		Data:   data,
	}

	defaultMessage := "Resource created successfully"
	if len(message) > 0 {
		defaultMessage = message[0]
	}
	response.Message = defaultMessage

	j.sendJSON(w, http.StatusCreated, response)
}

// SendError sends an error JSON response
func (j *JSONResponse) SendError(w http.ResponseWriter, statusCode int, message string, errors ...interface{}) {
	response := model.StandardResponse{
		Status:  model.StatusError,
		Message: message,
	}

	if len(errors) > 0 {
		response.Errors = errors[0]
	}

	j.sendJSON(w, statusCode, response)
}

// SendPaginatedResponse sends a paginated JSON response
func (j *JSONResponse) SendPaginatedResponse(
	w http.ResponseWriter,
	data interface{},
	page, limit, totalItems, totalPages int,
	message ...string,
) {
	response := model.PaginatedResponse{
		StandardResponse: model.StandardResponse{
			Status: model.StatusSuccess,
			Data:   data,
		},
		Page:       page,
		Limit:      limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}

	if len(message) > 0 {
		response.Message = message[0]
	}

	j.sendJSON(w, http.StatusOK, response)
}

// ValidationErrorResponse generates a structured validation error response
func (j *JSONResponse) ValidationErrorResponse(w http.ResponseWriter, validationErrors map[string]string) {
	j.SendError(
		w,
		http.StatusUnprocessableEntity,
		"Validation failed",
		validationErrors,
	)
}

// sendJSON is an internal method to send JSON response
func (j *JSONResponse) sendJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	// Set Content-Type header
	w.Header().Set("Content-Type", "application/json")

	// Set status code
	w.WriteHeader(statusCode)

	// Encode and send JSON
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}
