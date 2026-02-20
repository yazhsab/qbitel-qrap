package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/quantun-opensource/qrap/api/internal/model"
	"github.com/quantun-opensource/qrap/api/internal/service"
	qmw "github.com/quantun-opensource/qrap/shared/go/middleware"
)

type AssessmentHandler struct {
	svc    *service.AssessmentService
	logger *zap.Logger
}

func NewAssessmentHandler(svc *service.AssessmentService, logger *zap.Logger) *AssessmentHandler {
	return &AssessmentHandler{svc: svc, logger: logger}
}

func (h *AssessmentHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.Create)
	r.Get("/", h.List)
	r.Get("/{id}", h.Get)
	r.Post("/{id}/run", h.Run)
	return r
}

func (h *AssessmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateAssessmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" || req.OrganizationID == "" {
		writeError(w, http.StatusBadRequest, "name and organization_id are required")
		return
	}
	if len(req.Name) > maxNameLength {
		writeError(w, http.StatusBadRequest, "name exceeds maximum length")
		return
	}
	if req.CreatedBy == "" {
		req.CreatedBy = actorFromRequest(r)
	}

	assessment, err := h.svc.Create(r.Context(), &req)
	if err != nil {
		h.logger.Error("failed to create assessment", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "failed to create assessment")
		return
	}
	writeJSON(w, http.StatusCreated, assessment.ToResponse())
}

func (h *AssessmentHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid assessment ID")
		return
	}

	resp, err := h.svc.GetWithSummary(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "assessment not found")
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *AssessmentHandler) List(w http.ResponseWriter, r *http.Request) {
	pg := qmw.ParsePagination(r)
	status := r.URL.Query().Get("status")

	var orgID *uuid.UUID
	if orgStr := r.URL.Query().Get("organization_id"); orgStr != "" {
		parsed, err := uuid.Parse(orgStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid organization_id")
			return
		}
		orgID = &parsed
	}

	assessments, total, err := h.svc.List(r.Context(), orgID, status, pg.Offset, pg.Limit)
	if err != nil {
		h.logger.Error("failed to list assessments", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "failed to list assessments")
		return
	}

	resp := model.AssessmentListResponse{
		TotalCount: total,
		Offset:     pg.Offset,
		Limit:      pg.Limit,
	}
	for _, a := range assessments {
		resp.Assessments = append(resp.Assessments, a.ToResponse())
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *AssessmentHandler) Run(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid assessment ID")
		return
	}

	assessment, err := h.svc.Run(r.Context(), id)
	if err != nil {
		h.logger.Error("failed to run assessment", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "failed to run assessment")
		return
	}
	writeJSON(w, http.StatusOK, assessment.ToResponse())
}
