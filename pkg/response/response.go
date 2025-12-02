package response

import (
	"encoding/json"
	"net/http"
)

func WriteResponse(w http.ResponseWriter, data any, StatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(StatusCode)
	_ = json.NewEncoder(w).Encode(data)
}
