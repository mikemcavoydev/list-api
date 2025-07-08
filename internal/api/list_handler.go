package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/mikemcavoydev/list-api/internal/middleware"
	"github.com/mikemcavoydev/list-api/internal/store"
	"github.com/mikemcavoydev/list-api/internal/utils"
)

type ListHandler struct {
	listStore store.ListStore
	logger    *log.Logger
}

func NewListHandler(listStore store.ListStore, logger *log.Logger) *ListHandler {
	return &ListHandler{
		listStore: listStore,
		logger:    logger,
	}
}

func (h *ListHandler) HandleGetListById(w http.ResponseWriter, r *http.Request) {
	listID, err := utils.ReadIDParam(r)
	if err != nil {
		h.logger.Printf("ERROR: readIDParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid list id"})
		return
	}

	list, err := h.listStore.GetListByID(listID)
	if err != nil {
		h.logger.Printf("ERROR: getListByID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if list == nil {
		h.logger.Printf("ERROR: getListByID: %v", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "list not found"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"list": list})
}

func (h *ListHandler) HandleCreateListById(w http.ResponseWriter, r *http.Request) {
	var list store.List
	err := json.NewDecoder(r.Body).Decode(&list)
	if err != nil {
		h.logger.Printf("ERROR: decodingCreateList: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request sent"})
		return
	}

	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser == store.AnonymousUser {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "you must be logged in"})
		return
	}

	list.UserID = currentUser.ID

	createdList, err := h.listStore.CreateList(&list)
	if err != nil {
		h.logger.Printf("ERROR: createList: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to create list"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"list": createdList})
}

func (h *ListHandler) HandleUpdateListById(w http.ResponseWriter, r *http.Request) {
	listID, err := utils.ReadIDParam(r)
	if err != nil {
		h.logger.Printf("ERROR: readIDParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid list id"})
		return
	}

	existingList, err := h.listStore.GetListByID(listID)
	if err != nil {
		h.logger.Printf("ERROR: getListById: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to fetch list"})
		return
	}

	if existingList == nil {
		h.logger.Printf("ERROR: getListByID: %v", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "list not found"})
		return
	}

	var updateListRequest struct {
		Title       *string           `json:"title"`
		Description *string           `json:"description"`
		Entries     []store.ListEntry `json:"entries"`
	}

	err = json.NewDecoder(r.Body).Decode(&updateListRequest)
	if err != nil {
		h.logger.Printf("ERROR: decodingUpdateRequest: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
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

	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser == store.AnonymousUser {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "you must be logged in"})
		return
	}

	listOwner, err := h.listStore.GetListOwner(listID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "list does not exist"})
			return
		}

		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if listOwner != currentUser.ID {
		utils.WriteJSON(w, http.StatusForbidden, utils.Envelope{"error": "you are not authorized to update this list"})
		return
	}

	err = h.listStore.UpdateList(existingList)
	if err != nil {
		h.logger.Printf("ERROR: updatingList: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to update list"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"list": existingList})
}

func (h *ListHandler) HandleDeleteList(w http.ResponseWriter, r *http.Request) {
	listID, err := utils.ReadIDParam(r)
	if err != nil {
		h.logger.Printf("ERROR: readIDParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid list id"})
		return
	}

	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser == store.AnonymousUser {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "you must be logged in"})
		return
	}

	listOwner, err := h.listStore.GetListOwner(listID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "list does not exist"})
			return
		}

		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if listOwner != currentUser.ID {
		utils.WriteJSON(w, http.StatusForbidden, utils.Envelope{"error": "you are not authorized to delete this list"})
		return
	}

	err = h.listStore.DeleteList(listID)
	if err == sql.ErrNoRows {
		h.logger.Printf("ERROR: deletingList: %v", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "list does not exist"})
		return
	}

	if err != nil {
		h.logger.Printf("ERROR: deletingList: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to delete list"})
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, utils.Envelope{"list": "deleted successfully"})
}
