package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ListHandler struct{}

func NewListHandler() *ListHandler {
	return &ListHandler{}
}

func (h *ListHandler) HandleGetListById(w http.ResponseWriter, r *http.Request) {
	paramsListID := chi.URLParam(r, "id")
	if paramsListID == "" {
		http.NotFound(w, r)
		return
	}

	listID, err := strconv.ParseInt(paramsListID, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "this is the list id %d\n", listID)
}

func (h *ListHandler) HandleCreateListById(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "created a list\n")
}
