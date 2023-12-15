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

	json.NewEncoder(w).Encode(resp)
}
