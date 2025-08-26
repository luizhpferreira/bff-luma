package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bff-luma/internal/config"
	"bff-luma/internal/database"
	"bff-luma/internal/handlers"
	authmiddleware "bff-luma/internal/middleware"
	"bff-luma/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	// Carrega configurações
	cfg := config.LoadConfig()

	// Inicializa banco de dados
	db, err := database.NewDatabase(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Erro ao inicializar banco de dados: %v", err)
	}
	defer db.Close()

	// Inicializa serviços
	lnbitsService := services.NewLNBitsService(cfg.LNBitsBaseURL, cfg.LNBitsAPIToken, cfg.LNBitsWebhookSecret)
	jwtService := services.NewJWTService(cfg.JWTSecret)
	emailService := services.NewEmailService(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUsername, cfg.SMTPPassword, cfg.SMTPFromEmail, cfg.SMTPFromName, cfg.SMTPUseTLS, cfg.AppDomain, cfg.AppProtocol)
	passwordService := services.NewPasswordService()
	cleanupService := services.NewCleanupService(db)
	rateLimiter := services.NewRateLimiter()
	walletService := services.NewWalletService(db, lnbitsService, jwtService, emailService, passwordService)

	// Inicializa handlers
	walletHandler := handlers.NewWalletHandler(walletService, cleanupService, rateLimiter)

	// Inicia serviço de limpeza automática
	cleanupService.Start()

	// Configura roteador
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Rotas
	r.Get("/health", walletHandler.HealthCheck)
	r.Get("/confirm-email", walletHandler.ConfirmEmailPage)
	r.Get("/reset-password", walletHandler.ResetPasswordPage)

	// API v1
	r.Route("/api/v1", func(r chi.Router) {
		// Aplica rate limiting global
		r.Use(authmiddleware.RateLimitMiddleware(rateLimiter))
		
		// Rotas públicas (não precisam de autenticação)
		r.Post("/wallets", walletHandler.CreateWallet)
		r.Post("/login", walletHandler.Login)
		r.Post("/refresh", walletHandler.RefreshToken)
		r.Post("/forgot-password", walletHandler.ForgotPassword)
		r.Post("/reset-password", walletHandler.ResetPassword)
		r.Post("/validate-reset-token", walletHandler.ValidateResetToken)
		r.Post("/confirm-email", walletHandler.ConfirmEmail)
		
		// Rotas de administração (limpeza e rate limit stats)
		r.Post("/admin/cleanup", walletHandler.CleanupTokens)
		r.Get("/admin/cleanup/stats", walletHandler.GetCleanupStats)
		r.Get("/admin/rate-limit/stats", walletHandler.GetRateLimitStats)
		
		// Rotas protegidas (precisam de autenticação)
		r.Group(func(r chi.Router) {
			r.Use(authmiddleware.AuthMiddleware(jwtService))
			
			r.Get("/wallets", walletHandler.GetWalletInfo)
			r.Post("/invoices", walletHandler.CreateInvoice)
			r.Get("/payments/status", walletHandler.CheckPaymentStatus)

		})
	})

	// Configura servidor
	srv := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Canal para receber sinais de interrupção
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Inicia servidor em goroutine
	go func() {
		log.Printf("🚀 Servidor BFF Luma iniciado na porta %s", cfg.AppPort)
		log.Printf("📊 Health check: http://localhost:%s/health", cfg.AppPort)
		log.Printf("💳 API v1: http://localhost:%s/api/v1", cfg.AppPort)
		
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro ao iniciar servidor: %v", err)
		}
	}()

	// Aguarda sinal de interrupção
	<-stop
	log.Println("🛑 Recebido sinal de interrupção, encerrando servidor...")

	// Para o serviço de limpeza
	cleanupService.Stop()

	// Shutdown graceful
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Erro durante shutdown: %v", err)
	}

	log.Println("✅ Servidor encerrado com sucesso")
}
