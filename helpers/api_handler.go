package helpers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// APIRenderList will add the Result-Count header to the response before returning the JSON
func APIRenderList(w http.ResponseWriter, req *http.Request, resultCount int, result interface{}) error {
	w.Header().Add("Result-Count", strconv.Itoa(resultCount))
	err := json.NewEncoder(w).Encode(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "error rendering JSON", http.StatusInternalServerError)
		return err
	}
	return nil
}
