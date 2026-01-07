package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/common/utils"
	"github.com/M1ralai/go-modular-monolith-template/internal/common/validation"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/finance/dto"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/finance/service"
	"github.com/gorilla/mux"
)

type Handler struct{ service service.FinanceService }

func NewHandler(service service.FinanceService) *Handler { return &Handler{service: service} }

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/finance", h.GetAll).Methods("GET")
	router.HandleFunc("/finance/summary", h.GetSummary).Methods("GET")
	router.HandleFunc("/finance", h.Create).Methods("POST")
	router.HandleFunc("/finance/{id}", h.GetByID).Methods("GET")
	router.HandleFunc("/finance/{id}", h.Update).Methods("PUT", "PATCH")
	router.HandleFunc("/finance/{id}", h.Delete).Methods("DELETE")
}

func (h *Handler) getUserID(r *http.Request) int {
	return utils.GetUserIDFromContext(r.Context())
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz istek", err.Error())
		return
	}
	if err := validation.Get().Struct(req); err != nil {
		utils.ReturnError(w, "VALIDATION_ERROR", "Doğrulama hatası", validation.FormatErr(err))
		return
	}
	tx, err := h.service.Create(r.Context(), &req, h.getUserID(r))
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "İşlem oluşturulamadı", err.Error())
		return
	}
	utils.WriteJson(w, tx, http.StatusCreated, "İşlem oluşturuldu")
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	tx, err := h.service.GetByID(r.Context(), id, h.getUserID(r))
	if err != nil {
		utils.ReturnError(w, "NOT_FOUND", "İşlem bulunamadı", err.Error())
		return
	}
	utils.WriteJson(w, tx, http.StatusOK, "İşlem getirildi")
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	txs, err := h.service.GetAll(r.Context(), h.getUserID(r))
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "İşlemler getirilemedi", err.Error())
		return
	}
	utils.WriteJson(w, txs, http.StatusOK, "İşlemler getirildi")
}

func (h *Handler) GetSummary(w http.ResponseWriter, r *http.Request) {
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")
	start, _ := time.Parse("2006-01-02", startStr)
	end, _ := time.Parse("2006-01-02", endStr)
	if start.IsZero() {
		start = time.Now().AddDate(0, -1, 0)
	}
	if end.IsZero() {
		end = time.Now()
	}
	summary, err := h.service.GetSummary(r.Context(), h.getUserID(r), start, end)
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Özet getirilemedi", err.Error())
		return
	}
	utils.WriteJson(w, summary, http.StatusOK, "Finansal özet")
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var req dto.UpdateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz istek", err.Error())
		return
	}
	tx, err := h.service.Update(r.Context(), id, &req, h.getUserID(r))
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "İşlem güncellenemedi", err.Error())
		return
	}
	utils.WriteJson(w, tx, http.StatusOK, "İşlem güncellendi")
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if err := h.service.Delete(r.Context(), id, h.getUserID(r)); err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "İşlem silinemedi", err.Error())
		return
	}
	utils.WriteJson(w, nil, http.StatusOK, "İşlem silindi")
}
