package user

import (
	"encoding/json"
	"net/http"
	"errors"
	"strconv"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
		http.Error(
			w,
			"invalid user id",
			http.StatusBadRequest,
		)
		return
	}

	user, err := h.service.GetUser(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set(
		"Content-Type",
		"application/json",
	)
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) CreateUser(
	w http.ResponseWriter,
	r *http.Request,
) {
	var user User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	err = h.service.CreateUser(
		r.Context(),
		&user,
	)

	if err != nil {
		if errors.Is(err, ErrUserAlreadyExists) {
			http.Error(
				w,
				err.Error(),
				http.StatusConflict,
			)
			return
		}

		if errors.Is(err, ErrInvalidUser) {
			http.Error(
				w,
				err.Error(),
				http.StatusBadRequest,
			)
			return
		}

		http.Error(
			w,
			"internal server error",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(user)
}