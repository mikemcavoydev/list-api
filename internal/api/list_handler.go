package api

import (
	"database/sql"
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

	list, err := h.listStore.GetListByID(listID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "failed to fetch the list", http.StatusNotFound)
		return
	}

	if list == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(list)
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

func (h *ListHandler) HandleUpdateListById(w http.ResponseWriter, r *http.Request) {
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

	existingList, err := h.listStore.GetListByID(listID)
	if err != nil {
		http.Error(w, "failed to fetch list", http.StatusInternalServerError)
		return
	}

	if existingList == nil {
		http.NotFound(w, r)
		return
	}

	var updateListRequest struct {
		Title       *string           `json:"title"`
		Description *string           `json:"description"`
		Entries     []store.ListEntry `json:"entries"`
	}

	err = json.NewDecoder(r.Body).Decode(&updateListRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if updateListRequest.Title != nil {
		existingList.Title = *updateListRequest.Title
	}

	if updateListRequest.Description != nil {
		existingList.Description = *updateListRequest.Description
	}

	if updateListRequest.Entries != nil {
		existingList.Entries = updateListRequest.Entries
	}

	err = h.listStore.UpdateList(existingList)
	if err != nil {
		fmt.Println("update list error", err)
		http.Error(w, "failed to update the list", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existingList)
}

func (h *ListHandler) HandleDeleteList(w http.ResponseWriter, r *http.Request) {
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

	err = h.listStore.DeleteList(listID)
	if err == sql.ErrNoRows {
		http.Error(w, "list not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "failed to delete list", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
