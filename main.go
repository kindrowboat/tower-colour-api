package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type ColourMessage struct {
	Red     int    `json:"red"`
	Green   int    `   json:"green"`
	Blue    int    `json:"blue"`
	Message string `json:"message"`
}

func tomHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		// Just send out the JSON version of 'tom'
		colourMessage := ColourMessage{
			Red:     255,
			Green:   255,
			Blue:    255,
			Message: "a kind note",
		}
		j, _ := json.Marshal(colourMessage)
		w.Write(j)
	case "POST":
		// Decode the JSON in the body and overwrite 'tom' with it
		decoder := json.NewDecoder(r.Body)
		colourMessage := &ColourMessage{}
		err := decoder.Decode(colourMessage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		jsonMessage, _ := json.Marshal(colourMessage)
		w.Write(jsonMessage)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/", tomHandler)

	addr := ":3010"
	log.Printf("Serving on %s", addr)
	// err := http.ListenAndServe(fmt.Sprintf(":,", port), nil)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf(err.Error())
	}
}
