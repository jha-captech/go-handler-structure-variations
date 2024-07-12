package main

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	service := Service{}

	logger := slog.Default()

	mux.HandleFunc("Get /person", HandlerSimpleClosure(logger, &service))

	if err := http.ListenAndServe("localhost:8080", mux); err != nil {
		log.Fatal(err)
	}
}

// ── Handler ──────────────────────────────────────────────────────────────────────────────────────

type responseErr struct {
	Error string `json:"error"`
}

type responsePerson struct {
	Person Person `json:"person"`
}

func returnJSON(w http.ResponseWriter, logger *slog.Logger, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error("Error while marshaling data", "err", err, "data", data)
		http.Error(w, `{"Error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}
}

// HandlerSimpleClosure - all resources passed as parameters to handler.
func HandlerSimpleClosure(logger *slog.Logger, service *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve data from service.
		person, err := service.GetPerson()
		if err != nil {
			logger.Error("Error while retrieving data", "err", err)
			returnJSON(w, logger, http.StatusInternalServerError, responseErr{
				Error: "Internal server error",
			})
			return
		}

		// 	Return response data.
		returnJSON(w, logger, http.StatusInternalServerError, responsePerson{
			Person: person,
		})
	}
}

// ── Service ──────────────────────────────────────────────────────────────────────────────────────

type Service struct{}

type Person struct {
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Email  string `json:"email"`
	Active bool   `json:"active"`
}

func (s Service) GetPerson() (Person, error) {
	return Person{
		Name:   "John Doe",
		Age:    30,
		Email:  "john.doe@example.com",
		Active: true,
	}, nil
}
