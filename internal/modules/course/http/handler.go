package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/M1ralai/go-modular-monolith-template/internal/common/utils"
	"github.com/M1ralai/go-modular-monolith-template/internal/common/validation"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/course/dto"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/course/service"
	"github.com/gorilla/mux"
)

type Handler struct {
	service service.CourseService
}

func NewHandler(service service.CourseService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/courses", h.GetAll).Methods("GET")
	router.HandleFunc("/courses/active", h.GetActive).Methods("GET")
	router.HandleFunc("/courses", h.Create).Methods("POST")
	router.HandleFunc("/courses/components", h.CreateComponent).Methods("POST")
	router.HandleFunc("/courses/{id}/components", h.GetComponents).Methods("GET")
	router.HandleFunc("/courses/components/{id}", h.UpdateComponent).Methods("PUT", "PATCH")
	router.HandleFunc("/courses/components/{id}", h.DeleteComponent).Methods("DELETE")
	router.HandleFunc("/courses/schedules", h.CreateSchedule).Methods("POST")
	router.HandleFunc("/courses/{id}/schedules", h.GetSchedules).Methods("GET")
	router.HandleFunc("/courses/schedules/{id}", h.UpdateSchedule).Methods("PUT", "PATCH")
	router.HandleFunc("/courses/schedules/{id}", h.DeleteSchedule).Methods("DELETE")
	router.HandleFunc("/courses/{id}", h.GetByID).Methods("GET")
	router.HandleFunc("/courses/{id}", h.Update).Methods("PUT", "PATCH")
	router.HandleFunc("/courses/{id}", h.Delete).Methods("DELETE")
}

func (h *Handler) getUserID(r *http.Request) int {
	return utils.GetUserIDFromContext(r.Context())
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCourseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz istek formatı", err.Error())
		return
	}

	if err := validation.Get().Struct(req); err != nil {
		utils.ReturnError(w, "VALIDATION_ERROR", "Doğrulama hatası", validation.FormatErr(err))
		return
	}

	userID := h.getUserID(r)
	course, err := h.service.Create(r.Context(), &req, userID)
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Ders oluşturulamadı", err.Error())
		return
	}

	utils.WriteJson(w, course, http.StatusCreated, "Ders oluşturuldu")
}

func (h *Handler) CreateComponent(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateComponentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz istek formatı", err.Error())
		return
	}

	if err := validation.Get().Struct(req); err != nil {
		utils.ReturnError(w, "VALIDATION_ERROR", "Doğrulama hatası", validation.FormatErr(err))
		return
	}

	userID := h.getUserID(r)
	component, err := h.service.CreateComponent(r.Context(), &req, userID)
	if err != nil {
		if err.Error() == "course not found" {
			utils.ReturnError(w, "NOT_FOUND", "Ders bulunamadı", err.Error())
			return
		}
		if err.Error() == "unauthorized" {
			utils.ReturnError(w, "FORBIDDEN", "Bu işlem için yetkiniz yok", err.Error())
			return
		}
		utils.ReturnError(w, "INTERNAL_ERROR", "Component oluşturulamadı", err.Error())
		return
	}

	utils.WriteJson(w, component, http.StatusCreated, "Component oluşturuldu")
}

func (h *Handler) GetComponents(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz ID", err.Error())
		return
	}

	userID := h.getUserID(r)
	components, err := h.service.GetComponents(r.Context(), courseID, userID)
	if err != nil {
		if err.Error() == "course not found" {
			utils.ReturnError(w, "NOT_FOUND", "Ders bulunamadı", err.Error())
			return
		}
		if err.Error() == "unauthorized" {
			utils.ReturnError(w, "FORBIDDEN", "Bu işlem için yetkiniz yok", err.Error())
			return
		}
		utils.ReturnError(w, "INTERNAL_ERROR", "Components getirilemedi", err.Error())
		return
	}

	utils.WriteJson(w, components, http.StatusOK, "Components getirildi")
}

func (h *Handler) UpdateComponent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz ID", err.Error())
		return
	}

	var req dto.UpdateComponentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz istek formatı", err.Error())
		return
	}

	if err := validation.Get().Struct(req); err != nil {
		utils.ReturnError(w, "VALIDATION_ERROR", "Doğrulama hatası", validation.FormatErr(err))
		return
	}

	userID := h.getUserID(r)
	component, err := h.service.UpdateComponent(r.Context(), id, &req, userID)
	if err != nil {
		if err.Error() == "component not found" || err.Error() == "course not found" {
			utils.ReturnError(w, "NOT_FOUND", "Component bulunamadı", err.Error())
			return
		}
		if err.Error() == "unauthorized" {
			utils.ReturnError(w, "FORBIDDEN", "Bu işlem için yetkiniz yok", err.Error())
			return
		}
		utils.ReturnError(w, "INTERNAL_ERROR", "Component güncellenemedi", err.Error())
		return
	}

	utils.WriteJson(w, component, http.StatusOK, "Component güncellendi")
}

func (h *Handler) DeleteComponent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz ID", err.Error())
		return
	}

	userID := h.getUserID(r)
	if err := h.service.DeleteComponent(r.Context(), id, userID); err != nil {
		if err.Error() == "component not found" || err.Error() == "course not found" {
			utils.ReturnError(w, "NOT_FOUND", "Component bulunamadı", err.Error())
			return
		}
		if err.Error() == "unauthorized" {
			utils.ReturnError(w, "FORBIDDEN", "Bu işlem için yetkiniz yok", err.Error())
			return
		}
		utils.ReturnError(w, "INTERNAL_ERROR", "Component silinemedi", err.Error())
		return
	}

	utils.WriteJson(w, nil, http.StatusOK, "Component silindi")
}

func (h *Handler) CreateSchedule(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz istek formatı", err.Error())
		return
	}

	if err := validation.Get().Struct(req); err != nil {
		utils.ReturnError(w, "VALIDATION_ERROR", "Doğrulama hatası", validation.FormatErr(err))
		return
	}

	userID := h.getUserID(r)
	schedule, err := h.service.CreateSchedule(r.Context(), &req, userID)
	if err != nil {
		if err.Error() == "course not found" {
			utils.ReturnError(w, "NOT_FOUND", "Ders bulunamadı", err.Error())
			return
		}
		if err.Error() == "unauthorized" {
			utils.ReturnError(w, "FORBIDDEN", "Bu işlem için yetkiniz yok", err.Error())
			return
		}
		utils.ReturnError(w, "INTERNAL_ERROR", "Schedule oluşturulamadı", err.Error())
		return
	}

	utils.WriteJson(w, schedule, http.StatusCreated, "Schedule oluşturuldu")
}

func (h *Handler) GetSchedules(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz ID", err.Error())
		return
	}

	userID := h.getUserID(r)
	schedules, err := h.service.GetSchedules(r.Context(), courseID, userID)
	if err != nil {
		if err.Error() == "course not found" {
			utils.ReturnError(w, "NOT_FOUND", "Ders bulunamadı", err.Error())
			return
		}
		if err.Error() == "unauthorized" {
			utils.ReturnError(w, "FORBIDDEN", "Bu işlem için yetkiniz yok", err.Error())
			return
		}
		utils.ReturnError(w, "INTERNAL_ERROR", "Schedules getirilemedi", err.Error())
		return
	}

	utils.WriteJson(w, schedules, http.StatusOK, "Schedules getirildi")
}

func (h *Handler) UpdateSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz ID", err.Error())
		return
	}

	var req dto.UpdateScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz istek formatı", err.Error())
		return
	}

	if err := validation.Get().Struct(req); err != nil {
		utils.ReturnError(w, "VALIDATION_ERROR", "Doğrulama hatası", validation.FormatErr(err))
		return
	}

	userID := h.getUserID(r)
	schedule, err := h.service.UpdateSchedule(r.Context(), id, &req, userID)
	if err != nil {
		if err.Error() == "schedule not found" || err.Error() == "course not found" {
			utils.ReturnError(w, "NOT_FOUND", "Schedule bulunamadı", err.Error())
			return
		}
		if err.Error() == "unauthorized" {
			utils.ReturnError(w, "FORBIDDEN", "Bu işlem için yetkiniz yok", err.Error())
			return
		}
		utils.ReturnError(w, "INTERNAL_ERROR", "Schedule güncellenemedi", err.Error())
		return
	}

	utils.WriteJson(w, schedule, http.StatusOK, "Schedule güncellendi")
}

func (h *Handler) DeleteSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz ID", err.Error())
		return
	}

	userID := h.getUserID(r)
	if err := h.service.DeleteSchedule(r.Context(), id, userID); err != nil {
		if err.Error() == "schedule not found" || err.Error() == "course not found" {
			utils.ReturnError(w, "NOT_FOUND", "Schedule bulunamadı", err.Error())
			return
		}
		if err.Error() == "unauthorized" {
			utils.ReturnError(w, "FORBIDDEN", "Bu işlem için yetkiniz yok", err.Error())
			return
		}
		utils.ReturnError(w, "INTERNAL_ERROR", "Schedule silinemedi", err.Error())
		return
	}

	utils.WriteJson(w, nil, http.StatusOK, "Schedule silindi")
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz ID", err.Error())
		return
	}

	userID := h.getUserID(r)
	course, err := h.service.GetByID(r.Context(), id, userID)
	if err != nil {
		if err.Error() == "course not found" {
			utils.ReturnError(w, "NOT_FOUND", "Ders bulunamadı", err.Error())
			return
		}
		if err.Error() == "unauthorized" {
			utils.ReturnError(w, "FORBIDDEN", "Bu işlem için yetkiniz yok", err.Error())
			return
		}
		utils.ReturnError(w, "INTERNAL_ERROR", "Ders getirilemedi", err.Error())
		return
	}

	utils.WriteJson(w, course, http.StatusOK, "Ders getirildi")
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID(r)
	courses, err := h.service.GetAll(r.Context(), userID)
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Dersler getirilemedi", err.Error())
		return
	}

	utils.WriteJson(w, courses, http.StatusOK, "Dersler getirildi")
}

func (h *Handler) GetActive(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserID(r)
	courses, err := h.service.GetActive(r.Context(), userID)
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Aktif dersler getirilemedi", err.Error())
		return
	}

	utils.WriteJson(w, courses, http.StatusOK, "Aktif dersler getirildi")
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz ID", err.Error())
		return
	}

	var req dto.UpdateCourseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz istek formatı", err.Error())
		return
	}

	if err := validation.Get().Struct(req); err != nil {
		utils.ReturnError(w, "VALIDATION_ERROR", "Doğrulama hatası", validation.FormatErr(err))
		return
	}

	userID := h.getUserID(r)
	course, err := h.service.Update(r.Context(), id, &req, userID)
	if err != nil {
		if err.Error() == "course not found" {
			utils.ReturnError(w, "NOT_FOUND", "Ders bulunamadı", err.Error())
			return
		}
		if err.Error() == "unauthorized" {
			utils.ReturnError(w, "FORBIDDEN", "Bu işlem için yetkiniz yok", err.Error())
			return
		}
		utils.ReturnError(w, "INTERNAL_ERROR", "Ders güncellenemedi", err.Error())
		return
	}

	utils.WriteJson(w, course, http.StatusOK, "Ders güncellendi")
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
		if err.Error() == "course not found" {
			utils.ReturnError(w, "NOT_FOUND", "Ders bulunamadı", err.Error())
			return
		}
		if err.Error() == "unauthorized" {
			utils.ReturnError(w, "FORBIDDEN", "Bu işlem için yetkiniz yok", err.Error())
			return
		}
		utils.ReturnError(w, "INTERNAL_ERROR", "Ders silinemedi", err.Error())
		return
	}

	utils.WriteJson(w, nil, http.StatusOK, "Ders silindi")
}
