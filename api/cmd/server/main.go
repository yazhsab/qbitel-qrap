package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/quantun-opensource/qrap/api/internal/config"
	"github.com/quantun-opensource/qrap/api/internal/handler"
	"github.com/quantun-opensource/qrap/api/internal/repository"
	"github.com/quantun-opensource/qrap/api/internal/service"
	qdb "github.com/quantun-opensource/qrap/shared/go/database"
	qmw "github.com/quantun-opensource/qrap/shared/go/middleware"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	ctx := context.Background()
	poolCfg := qdb.DefaultPoolConfig(cfg.DatabaseURL)
	poolCfg.Logger = logger
	pool, err := qdb.NewPool(ctx, poolCfg)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer pool.Close()

	// Repositories
	orgRepo := repository.NewOrganizationRepository(pool)
	assessmentRepo := repository.NewAssessmentRepository(pool)
	findingRepo := repository.NewFindingRepository(pool)

	// Services
	orgSvc := service.NewOrganizationService(orgRepo, logger)
	assessmentSvc := service.NewAssessmentService(assessmentRepo, findingRepo, logger)
	findingSvc := service.NewFindingService(findingRepo, logger)

	// Handlers
	healthH := handler.NewHealthHandler()
	orgH := handler.NewOrganizationHandler(orgSvc, logger)
	assessmentH := handler.NewAssessmentHandler(assessmentSvc, logger)
	findingH := handler.NewFindingHandler(findingSvc, logger)

	// Router
	r := chi.NewRouter()

	// --- Infrastructure middleware ---
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.Timeout(30 * time.Second))

	// --- Security middleware ---
	r.Use(qmw.SecurityHeaders(qmw.DefaultSecurityHeadersConfig()))
	r.Use(qmw.MaxBodySize(cfg.MaxBodyBytes))

	// CORS (only if origins are configured)
	if len(cfg.CORSOrigins) > 0 {
		corsConfig := qmw.DefaultCORSConfig()
		corsConfig.AllowedOrigins = cfg.CORSOrigins
		r.Use(qmw.CORS(corsConfig))
	}

	// Rate limiting (100 req/min per IP)
	rateLimiter := qmw.NewRateLimiter(qmw.DefaultRateLimitConfig())
	defer rateLimiter.Stop()
	r.Use(rateLimiter.Middleware())

	// --- Auth middleware ---
	authConfig := qmw.AuthConfig{
		JWTSecret: cfg.JWTSecret,
		JWTIssuer: cfg.JWTIssuer,
		APIKeys:   qmw.ParseAPIKeyEntries(cfg.APIKeys),
		SkipPaths: []string{"/health"},
		Logger:    logger,
	}

	// Health endpoint (unauthenticated)
	r.Get("/health", healthH.Check)

	// API routes (authenticated)
	r.Route("/api/v1", func(r chi.Router) {
		// Apply auth only if JWT secret or API keys are configured
		if cfg.JWTSecret != "" || len(cfg.APIKeys) > 0 {
			r.Use(qmw.Auth(authConfig))
		}

		r.Mount("/organizations", orgH.Routes())
		r.Mount("/assessments", assessmentH.Routes())
		r.Mount("/findings", findingH.Routes())
	})

	addr := fmt.Sprintf(":%s", cfg.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("starting QRAP API server",
			zap.String("addr", addr),
			zap.Bool("auth_enabled", cfg.JWTSecret != "" || len(cfg.APIKeys) > 0),
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("server forced to shutdown", zap.Error(err))
	}
	logger.Info("server stopped")
}
