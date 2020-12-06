package web

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteAsJSON(writer http.ResponseWriter, data interface{}) error {
	itemJSON, err := json.Marshal(data)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return err
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(itemJSON)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}
