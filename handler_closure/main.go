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

	handler := Handler{
		logger:  logger,
		service: service,
	}

	mux.HandleFunc("Get /person", HandlerHandlerClosure(&handler))

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
	}
}

// HandlerHandlerClosure - handler struct passed through closure
func HandlerHandlerClosure(handler *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve data from service.
		person, err := handler.service.GetPerson()
		if err != nil {
			handler.logger.Error("Error while retrieving data", "err", err)
			returnJSON(w, handler.logger, http.StatusInternalServerError, responseErr{
				Error: "Internal server error",
			})
			return
		}

		// 	Return response data.
		returnJSON(w, handler.logger, http.StatusInternalServerError, responsePerson{
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

type Handler struct {
	logger  *slog.Logger
	service Service
}

// HandlerClosureMethod - method on handler struct
func (h Handler) HandlerClosureMethod() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.logger.Info("Health check called")
	}
}
