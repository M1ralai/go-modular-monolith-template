package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/M1ralai/go-modular-monolith-template/internal/common/utils"
	"github.com/M1ralai/go-modular-monolith-template/internal/common/validation"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/user/dto"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/user/service"
	"github.com/gorilla/mux"
)

type Handler struct {
	service service.UserService
}

func NewHandler(service service.UserService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/users", h.GetAllUsers).Methods("GET")
	router.HandleFunc("/users", h.CreateUser).Methods("POST")
	router.HandleFunc("/users/{id}", h.GetUser).Methods("GET")
	router.HandleFunc("/users/{id}", h.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", h.DeleteUser).Methods("DELETE")
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz istek formatı", err.Error())
		return
	}

	if err := validation.Get().Struct(req); err != nil {
		utils.ReturnError(w, "VALIDATION_ERROR", "Doğrulama hatası", validation.FormatErr(err))
		return
	}

	user, err := h.service.CreateUser(r.Context(), &req)
	if err != nil {
		if err.Error() == "email already exists" {
			utils.ReturnError(w, "BAD_REQUEST", "Bu e-posta adresi zaten kullanımda", err.Error())
			return
		}
		utils.ReturnError(w, "INTERNAL_ERROR", "Kullanıcı oluşturulamadı", err.Error())
		return
	}

	utils.WriteJson(w, user, http.StatusCreated, "Kullanıcı başarıyla oluşturuldu")
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz kullanıcı ID", err.Error())
		return
	}

	user, err := h.service.GetUser(r.Context(), id)
	if err != nil {
		if err.Error() == "user not found" {
			utils.ReturnError(w, "NOT_FOUND", "Kullanıcı bulunamadı", err.Error())
			return
		}
		utils.ReturnError(w, "INTERNAL_ERROR", "Kullanıcı getirilemedi", err.Error())
		return
	}

	utils.WriteJson(w, user, http.StatusOK, "Kullanıcı getirildi")
}

func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAllUsers(r.Context())
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Kullanıcılar getirilemedi", err.Error())
		return
	}

	utils.WriteJson(w, users, http.StatusOK, "Kullanıcılar getirildi")
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz kullanıcı ID", err.Error())
		return
	}

	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz istek formatı", err.Error())
		return
	}

	if err := validation.Get().Struct(req); err != nil {
		utils.ReturnError(w, "VALIDATION_ERROR", "Doğrulama hatası", validation.FormatErr(err))
		return
	}

	user, err := h.service.UpdateUser(r.Context(), id, &req)
	if err != nil {
		if err.Error() == "user not found" {
			utils.ReturnError(w, "NOT_FOUND", "Kullanıcı bulunamadı", err.Error())
			return
		}
		if err.Error() == "email already exists" {
			utils.ReturnError(w, "BAD_REQUEST", "Bu e-posta adresi zaten kullanımda", err.Error())
			return
		}
		utils.ReturnError(w, "INTERNAL_ERROR", "Kullanıcı güncellenemedi", err.Error())
		return
	}

	utils.WriteJson(w, user, http.StatusOK, "Kullanıcı güncellendi")
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz kullanıcı ID", err.Error())
		return
	}

	if err := h.service.DeleteUser(r.Context(), id); err != nil {
		if err.Error() == "user not found" {
			utils.ReturnError(w, "NOT_FOUND", "Kullanıcı bulunamadı", err.Error())
			return
		}
		utils.ReturnError(w, "INTERNAL_ERROR", "Kullanıcı silinemedi", err.Error())
		return
	}

	utils.WriteJson(w, nil, http.StatusOK, "Kullanıcı silindi")
}
