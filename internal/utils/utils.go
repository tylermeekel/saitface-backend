package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Data any `json:"data"`
}

func SendJSON(w http.ResponseWriter, data any) {
	var resp Response

	resp.Data = data

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
