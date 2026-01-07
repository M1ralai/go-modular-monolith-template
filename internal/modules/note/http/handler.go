package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/M1ralai/go-modular-monolith-template/internal/common/utils"
	"github.com/M1ralai/go-modular-monolith-template/internal/common/validation"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/note/dto"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/note/service"
	"github.com/gorilla/mux"
)

type Handler struct {
	service service.NoteService
}

func NewHandler(service service.NoteService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/notes", h.GetAll).Methods("GET")
	router.HandleFunc("/notes/favorites", h.GetFavorites).Methods("GET")
	router.HandleFunc("/notes/search", h.Search).Methods("GET")
	router.HandleFunc("/notes", h.Create).Methods("POST")
	router.HandleFunc("/notes/{id}", h.GetByID).Methods("GET")
	router.HandleFunc("/notes/{id}", h.Update).Methods("PUT", "PATCH")
	router.HandleFunc("/notes/{id}", h.Delete).Methods("DELETE")
	router.HandleFunc("/notes/{id}/backlinks", h.GetBacklinks).Methods("GET")
	router.HandleFunc("/notes/{id}/links", h.CreateLink).Methods("POST")
	router.HandleFunc("/notes/links/{linkId}", h.DeleteLink).Methods("DELETE")
}

func (h *Handler) getUserID(r *http.Request) int {
	return utils.GetUserIDFromContext(r.Context())
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz istek formatı", err.Error())
		return
	}

	if err := validation.Get().Struct(req); err != nil {
		utils.ReturnError(w, "VALIDATION_ERROR", "Doğrulama hatası", validation.FormatErr(err))
		return
	}

	userID := h.getUserID(r)
	note, err := h.service.Create(r.Context(), &req, userID)
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Not oluşturulamadı", err.Error())
		return
	}

	utils.WriteJson(w, note, http.StatusCreated, "Not oluşturuldu")
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz ID", err.Error())
		return
	}

	userID := h.getUserID(r)
	note, err := h.service.GetByID(r.Context(), id, userID)
	if err != nil {
		if err.Error() == "note not found" {
			utils.ReturnError(w, "NOT_FOUND", "Not bulunamadı", err.Error())
			return
		}
		if err.Error() == "unauthorized" {
			utils.ReturnError(w, "FORBIDDEN", "Bu işlem için yetkiniz yok", err.Error())
			return
		}
		utils.ReturnError(w, "INTERNAL_ERROR", "Not getirilemedi", err.Error())
		return
	}

	utils.WriteJson(w, note, http.StatusOK, "Not getirildi")
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID(r)
	notes, err := h.service.GetAll(r.Context(), userID)
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Notlar getirilemedi", err.Error())
		return
	}

	utils.WriteJson(w, notes, http.StatusOK, "Notlar getirildi")
}

func (h *Handler) GetFavorites(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID(r)
	notes, err := h.service.GetFavorites(r.Context(), userID)
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Favori notlar getirilemedi", err.Error())
		return
	}

	utils.WriteJson(w, notes, http.StatusOK, "Favori notlar getirildi")
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		utils.ReturnError(w, "BAD_REQUEST", "Arama terimi gerekli", "q parametresi boş")
		return
	}

	userID := h.getUserID(r)
	notes, err := h.service.Search(r.Context(), userID, query)
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Arama yapılamadı", err.Error())
		return
	}

	utils.WriteJson(w, notes, http.StatusOK, "Arama sonuçları")
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz ID", err.Error())
		return
	}

	var req dto.UpdateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz istek formatı", err.Error())
		return
	}

	if err := validation.Get().Struct(req); err != nil {
		utils.ReturnError(w, "VALIDATION_ERROR", "Doğrulama hatası", validation.FormatErr(err))
		return
	}

	userID := h.getUserID(r)
	note, err := h.service.Update(r.Context(), id, &req, userID)
	if err != nil {
		if err.Error() == "note not found" {
			utils.ReturnError(w, "NOT_FOUND", "Not bulunamadı", err.Error())
			return
		}
		if err.Error() == "unauthorized" {
			utils.ReturnError(w, "FORBIDDEN", "Bu işlem için yetkiniz yok", err.Error())
			return
		}
		utils.ReturnError(w, "INTERNAL_ERROR", "Not güncellenemedi", err.Error())
		return
	}

	utils.WriteJson(w, note, http.StatusOK, "Not güncellendi")
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz ID", err.Error())
		return
	}

	userID := h.getUserID(r)
	if err := h.service.Delete(r.Context(), id, userID); err != nil {
		if err.Error() == "note not found" {
			utils.ReturnError(w, "NOT_FOUND", "Not bulunamadı", err.Error())
			return
		}
		if err.Error() == "unauthorized" {
			utils.ReturnError(w, "FORBIDDEN", "Bu işlem için yetkiniz yok", err.Error())
			return
		}
		utils.ReturnError(w, "INTERNAL_ERROR", "Not silinemedi", err.Error())
		return
	}

	utils.WriteJson(w, nil, http.StatusOK, "Not silindi")
}

func (h *Handler) GetBacklinks(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz ID", err.Error())
		return
	}

	userID := h.getUserID(r)
	backlinks, err := h.service.GetBacklinks(r.Context(), id, userID)
	if err != nil {
		if err.Error() == "note not found" {
			utils.ReturnError(w, "NOT_FOUND", "Not bulunamadı", err.Error())
			return
		}
		if err.Error() == "unauthorized" {
			utils.ReturnError(w, "FORBIDDEN", "Bu işlem için yetkiniz yok", err.Error())
			return
		}
		utils.ReturnError(w, "INTERNAL_ERROR", "Backlinkler getirilemedi", err.Error())
		return
	}

	utils.WriteJson(w, backlinks, http.StatusOK, "Backlinkler getirildi")
}

func (h *Handler) CreateLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sourceID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz ID", err.Error())
		return
	}

	var req dto.CreateNoteLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz istek formatı", err.Error())
		return
	}

	if err := validation.Get().Struct(req); err != nil {
		utils.ReturnError(w, "VALIDATION_ERROR", "Doğrulama hatası", validation.FormatErr(err))
		return
	}

	userID := h.getUserID(r)
	link, err := h.service.CreateLink(r.Context(), sourceID, &req, userID)
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Link oluşturulamadı", err.Error())
		return
	}

	utils.WriteJson(w, link, http.StatusCreated, "Link oluşturuldu")
}

func (h *Handler) DeleteLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	linkID, err := strconv.Atoi(vars["linkId"])
	if err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz link ID", err.Error())
		return
	}

	userID := h.getUserID(r)
	if err := h.service.DeleteLink(r.Context(), linkID, userID); err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Link silinemedi", err.Error())
		return
	}

	utils.WriteJson(w, nil, http.StatusOK, "Link silindi")
}
