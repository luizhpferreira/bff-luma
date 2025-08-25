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

// cleanup executa a limpeza de tokens expirados e contas não confirmadas
func (s *CleanupService) cleanup() {
	// 1. Limpeza de tokens expirados
	s.cleanupExpiredTokens()
	
	// 2. Limpeza de contas não confirmadas (a cada 6 horas)
	if time.Now().Hour()%6 == 0 {
		s.cleanupUnconfirmedAccounts()
	}
}

// cleanupExpiredTokens limpa tokens de reset expirados
func (s *CleanupService) cleanupExpiredTokens() {
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

	log.Printf("✅ Limpeza de tokens concluída: %d tokens removidos", cleanedCount)
}

// cleanupUnconfirmedAccounts limpa contas não confirmadas após 24h
func (s *CleanupService) cleanupUnconfirmedAccounts() {
	log.Println("🧹 Verificando contas não confirmadas...")
	
	// Executa a limpeza
	cleanedCount, err := s.db.CleanUnconfirmedAccounts()
	if err != nil {
		log.Printf("❌ Erro ao limpar contas não confirmadas: %v", err)
		return
	}

	if cleanedCount > 0 {
		log.Printf("✅ Limpeza de contas não confirmadas: %d contas removidas", cleanedCount)
	} else {
		log.Println("✨ Nenhuma conta não confirmada encontrada para remoção")
	}
}

// CleanupNow executa uma limpeza imediata (para testes ou limpeza manual)
func (s *CleanupService) CleanupNow() error {
	log.Println("🧹 Executando limpeza manual...")
	
	// 1. Limpeza de tokens expirados
	s.cleanupExpiredTokens()
	
	// 2. Limpeza de contas não confirmadas
	s.cleanupUnconfirmedAccounts()
	
	log.Println("✅ Limpeza manual concluída")
	return nil
}

// CleanUnconfirmedAccounts remove contas não confirmadas após 24h
func (s *CleanupService) CleanUnconfirmedAccounts() error {
	log.Println("🧹 Iniciando limpeza de contas não confirmadas...")
	
	// Remove contas não confirmadas criadas há mais de 24h
	deletedCount, err := s.db.CleanUnconfirmedAccounts()
	if err != nil {
		log.Printf("❌ Erro ao limpar contas não confirmadas: %v", err)
		return err
	}
	
	if deletedCount > 0 {
		log.Printf("✅ Removidas %d contas não confirmadas", deletedCount)
	} else {
		log.Println("✨ Nenhuma conta não confirmada encontrada para remoção")
	}
	
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
