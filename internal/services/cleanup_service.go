package services

import (
	"fmt"
	"log"
	"time"

	"bff-luma/internal/database"
)

// CleanupService gerencia a limpeza automática de dados expirados
type CleanupService struct {
	db           *database.Database
	interval     time.Duration
	stopChan     chan bool
	isRunning    bool
}

// NewCleanupService cria um novo serviço de limpeza
func NewCleanupService(db *database.Database) *CleanupService {
	return &CleanupService{
		db:       db,
		interval: 1 * time.Hour, // Limpa a cada hora
		stopChan: make(chan bool),
	}
}

// Start inicia o serviço de limpeza automática
func (s *CleanupService) Start() {
	if s.isRunning {
		log.Println("⚠️ Serviço de limpeza já está rodando")
		return
	}

	s.isRunning = true
	log.Println("🧹 Iniciando serviço de limpeza automática...")

	go s.run()
}

// Stop para o serviço de limpeza
func (s *CleanupService) Stop() {
	if !s.isRunning {
		return
	}

	log.Println("🛑 Parando serviço de limpeza automática...")
	s.stopChan <- true
	s.isRunning = false
}

// run executa o loop principal de limpeza
func (s *CleanupService) run() {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	// Executa limpeza imediatamente na primeira vez
	s.cleanup()

	for {
		select {
		case <-ticker.C:
			s.cleanup()
		case <-s.stopChan:
			log.Println("✅ Serviço de limpeza parado")
			return
		}
	}
}

// cleanup executa a limpeza de tokens expirados
func (s *CleanupService) cleanup() {
	// Conta tokens expirados antes da limpeza
	expiredCount, err := s.db.GetExpiredTokensCount()
	if err != nil {
		log.Printf("❌ Erro ao contar tokens expirados: %v", err)
		return
	}

	if expiredCount == 0 {
		log.Println("✨ Nenhum token expirado encontrado")
		return
	}

	log.Printf("🧹 Encontrados %d tokens expirados, iniciando limpeza...", expiredCount)

	// Executa a limpeza
	cleanedCount, err := s.db.CleanExpiredTokens()
	if err != nil {
		log.Printf("❌ Erro ao limpar tokens expirados: %v", err)
		return
	}

	log.Printf("✅ Limpeza concluída: %d tokens removidos", cleanedCount)
}

// CleanupNow executa uma limpeza imediata (para testes ou limpeza manual)
func (s *CleanupService) CleanupNow() error {
	log.Println("🧹 Executando limpeza manual...")
	
	expiredCount, err := s.db.GetExpiredTokensCount()
	if err != nil {
		return fmt.Errorf("erro ao contar tokens expirados: %w", err)
	}

	if expiredCount == 0 {
		log.Println("✨ Nenhum token expirado encontrado")
		return nil
	}

	cleanedCount, err := s.db.CleanExpiredTokens()
	if err != nil {
		return fmt.Errorf("erro ao limpar tokens expirados: %w", err)
	}

	log.Printf("✅ Limpeza manual concluída: %d tokens removidos", cleanedCount)
	return nil
}

// GetStats retorna estatísticas de limpeza
func (s *CleanupService) GetStats() (map[string]interface{}, error) {
	expiredCount, err := s.db.GetExpiredTokensCount()
	if err != nil {
		return nil, fmt.Errorf("erro ao obter estatísticas: %w", err)
	}

	return map[string]interface{}{
		"expired_tokens": expiredCount,
		"is_running":     s.isRunning,
		"interval":       s.interval.String(),
	}, nil
}
