package http

import (
	"net/http"
	"strconv"

	"github.com/M1ralai/go-modular-monolith-template/internal/common/utils"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/job/service"
	"github.com/gorilla/mux"
)

// Handler handles job-related HTTP requests
type Handler struct {
	service service.JobService
}

// NewHandler creates a new job handler
func NewHandler(service service.JobService) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers job routes
func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/jobs", h.ListJobs).Methods("GET")
	r.HandleFunc("/jobs/{job_name}/trigger", h.TriggerJob).Methods("POST")
	r.HandleFunc("/jobs/{job_name}/status", h.GetJobStatus).Methods("GET")
	r.HandleFunc("/jobs/{job_name}/history", h.GetJobHistory).Methods("GET")
}

// ListJobs lists all registered jobs
// GET /jobs
func (h *Handler) ListJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.service.ListJobs(r.Context())
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "İşler listelenemedi", err.Error())
		return
	}

	utils.WriteJson(w, jobs, http.StatusOK, "İşler listelendi")
}

// TriggerJob manually triggers a job
// POST /jobs/{job_name}/trigger
func (h *Handler) TriggerJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobName := vars["job_name"]

	execution, err := h.service.TriggerJob(r.Context(), jobName)
	if err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "İş tetiklenemedi", err.Error())
		return
	}

	utils.WriteJson(w, map[string]interface{}{
		"message":      "Job triggered",
		"job_name":     jobName,
		"execution_id": execution.ID,
	}, http.StatusAccepted, "İş tetiklendi")
}

// GetJobStatus returns the current status of a job
// GET /jobs/{job_name}/status
func (h *Handler) GetJobStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobName := vars["job_name"]

	execution, err := h.service.GetJobStatus(r.Context(), jobName)
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "İş durumu alınamadı", err.Error())
		return
	}

	if execution == nil {
		utils.WriteJson(w, map[string]interface{}{
			"job_name": jobName,
			"status":   "never_run",
		}, http.StatusOK, "İş hiç çalıştırılmamış")
		return
	}

	utils.WriteJson(w, map[string]interface{}{
		"job_name":     execution.JobName,
		"status":       execution.Status,
		"started_at":   execution.StartedAt,
		"completed_at": execution.CompletedAt,
		"duration_ms":  execution.DurationMs,
		"error":        execution.Error,
	}, http.StatusOK, "İş durumu")
}

// GetJobHistory returns the execution history for a job
// GET /jobs/{job_name}/history?limit=20
func (h *Handler) GetJobHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobName := vars["job_name"]

	limit := 20
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	executions, err := h.service.GetJobHistory(r.Context(), jobName, limit)
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "İş geçmişi alınamadı", err.Error())
		return
	}

	result := make([]map[string]interface{}, len(executions))
	for i, exec := range executions {
		result[i] = map[string]interface{}{
			"id":           exec.ID,
			"status":       exec.Status,
			"started_at":   exec.StartedAt,
			"completed_at": exec.CompletedAt,
			"duration_ms":  exec.DurationMs,
			"error":        exec.Error,
		}
	}

	utils.WriteJson(w, map[string]interface{}{
		"job_name":   jobName,
		"executions": result,
	}, http.StatusOK, "İş geçmişi")
}
