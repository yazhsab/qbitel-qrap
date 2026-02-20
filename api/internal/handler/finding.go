package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/quantun-opensource/qrap/api/internal/model"
	"github.com/quantun-opensource/qrap/api/internal/service"
	qmw "github.com/quantun-opensource/qrap/shared/go/middleware"
)

type FindingHandler struct {
	svc    *service.FindingService
	logger *zap.Logger
}

func NewFindingHandler(svc *service.FindingService, logger *zap.Logger) *FindingHandler {
	return &FindingHandler{svc: svc, logger: logger}
}

func (h *FindingHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/{id}", h.Get)
	r.Get("/", h.List)
	return r
}

func (h *FindingHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid finding ID")
		return
	}

	finding, err := h.svc.Get(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "finding not found")
		return
	}
	writeJSON(w, http.StatusOK, finding.ToResponse())
}

func (h *FindingHandler) List(w http.ResponseWriter, r *http.Request) {
	assessmentIDStr := r.URL.Query().Get("assessment_id")
	if assessmentIDStr == "" {
		writeError(w, http.StatusBadRequest, "assessment_id is required")
		return
	}
	assessmentID, err := uuid.Parse(assessmentIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid assessment_id")
		return
	}

	pg := qmw.ParsePagination(r)
	riskLevel := r.URL.Query().Get("risk_level")
	category := r.URL.Query().Get("category")

	findings, total, err := h.svc.ListByAssessment(r.Context(), assessmentID, riskLevel, category, pg.Offset, pg.Limit)
	if err != nil {
		h.logger.Error("failed to list findings", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "failed to list findings")
		return
	}

	resp := model.FindingListResponse{
		TotalCount: total,
		Offset:     pg.Offset,
		Limit:      pg.Limit,
	}
	for _, f := range findings {
		resp.Findings = append(resp.Findings, f.ToResponse())
	}
	writeJSON(w, http.StatusOK, resp)
}
