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

// maxNameLength is the maximum allowed length for name fields.
const maxNameLength = 255

type OrganizationHandler struct {
	svc    *service.OrganizationService
	logger *zap.Logger
}

func NewOrganizationHandler(svc *service.OrganizationService, logger *zap.Logger) *OrganizationHandler {
	return &OrganizationHandler{svc: svc, logger: logger}
}

func (h *OrganizationHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.Create)
	r.Get("/", h.List)
	r.Get("/{id}", h.Get)
	return r
}

func (h *OrganizationHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateOrganizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	if len(req.Name) > maxNameLength {
		writeError(w, http.StatusBadRequest, "name exceeds maximum length")
		return
	}
	if req.CreatedBy == "" {
		req.CreatedBy = actorFromRequest(r)
	}

	org, err := h.svc.Create(r.Context(), &req)
	if err != nil {
		h.logger.Error("failed to create organization", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "failed to create organization")
		return
	}
	writeJSON(w, http.StatusCreated, org.ToResponse())
}

func (h *OrganizationHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid organization ID")
		return
	}
	org, err := h.svc.Get(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "organization not found")
		return
	}
	writeJSON(w, http.StatusOK, org.ToResponse())
}

func (h *OrganizationHandler) List(w http.ResponseWriter, r *http.Request) {
	pg := qmw.ParsePagination(r)

	orgs, total, err := h.svc.List(r.Context(), pg.Offset, pg.Limit)
	if err != nil {
		h.logger.Error("failed to list organizations", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "failed to list organizations")
		return
	}

	resp := model.OrganizationListResponse{
		TotalCount: total,
		Offset:     pg.Offset,
		Limit:      pg.Limit,
	}
	for _, o := range orgs {
		resp.Organizations = append(resp.Organizations, o.ToResponse())
	}
	writeJSON(w, http.StatusOK, resp)
}
