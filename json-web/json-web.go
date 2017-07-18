package jsonweb

import (
	"encoding/json"
	"net/http"
)

// JSON is a simple wrapper for json that is then returned to the user
type JSON struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

func (j JSON) Write(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(j)
}
