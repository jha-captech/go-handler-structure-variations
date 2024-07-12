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

	mux.HandleFunc("Get /person", handler.HandlerMethodClosure())

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

type Handler struct {
	logger  *slog.Logger
	service Service
}

func (h *Handler) returnJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Error while marshaling data", "err", err, "data", data)
		http.Error(w, `{"Error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}
}

// HandlerMethodClosure - handler struct passed through closure
func (h *Handler) HandlerMethodClosure() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve data from service.
		person, err := h.service.GetPerson()
		if err != nil {
			h.logger.Error("Error while retrieving data", "err", err)
			h.returnJSON(w, http.StatusInternalServerError, responseErr{
				Error: "Internal server error",
			})
			return
		}

		// 	Return response data.
		h.returnJSON(w, http.StatusInternalServerError, responsePerson{
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
