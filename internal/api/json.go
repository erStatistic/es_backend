package rumiapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		fmt.Println(err)
	}

	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}

	respondWithJson(w, code, msg, nil)
}

func respondWithJson(w http.ResponseWriter, code int, msg string, payload any) {
	w.Header().Set("Content-Type", "application/json")
	Response := APIResponse{
		Code: code,
		Msg:  msg,
		Data: payload,
	}
	data, err := json.Marshal(Response)
	if err != nil {
		log.Printf("Error marshalling JSON: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	data = append(data, byte('\n'))
	_, _ = w.Write(data)
}
