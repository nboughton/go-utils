package jsonweb

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// JSON is a simple wrapper for json that is then returned to the user
type JSON struct {
	Status int         `json:"-"` // Do not encode status, this goes in the response header
	Data   interface{} `json:"data"`
}

// New creates a new JSON object which can variably accept data
func New(status int, data interface{}) *JSON {
	return &JSON{Status: status, Data: data}
}

// Write sends the JSON object back to the client via browser via the original ResponseWriter
func (j *JSON) Write(w http.ResponseWriter) error {
	w.Header().Set("Status", strconv.Itoa(j.Status))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(j)
}
