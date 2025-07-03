package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mikemcavoydev/list-api/internal/store"
)

type ListHandler struct {
	listStore store.ListStore
}

func NewListHandler(listStore store.ListStore) *ListHandler {
	return &ListHandler{
		listStore: listStore,
	}
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
	var list store.List
	err := json.NewDecoder(r.Body).Decode(&list)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "failed to create list", http.StatusInternalServerError)
		return
	}

	createdList, err := h.listStore.CreateList(&list)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "failed to create list", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdList)
}
