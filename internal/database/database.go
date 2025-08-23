package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"bff-luma/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

// Database representa a conexão com o banco de dados
type Database struct {
	db *sql.DB
}

// NewDatabase cria uma nova conexão com o banco de dados
func NewDatabase(dbURL string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbURL)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir banco de dados: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao conectar com banco de dados: %w", err)
	}

	database := &Database{db: db}
	if err := database.createTables(); err != nil {
		return nil, fmt.Errorf("erro ao criar tabelas: %w", err)
	}

	return database, nil
}

// createTables cria as tabelas necessárias
func (d *Database) createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS wallets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		wallet_id TEXT NOT NULL UNIQUE,
		admin_key TEXT NOT NULL,
		invoice_key TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS reset_tokens (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL,
		token TEXT NOT NULL UNIQUE,
		expires_at DATETIME NOT NULL,
		used BOOLEAN DEFAULT FALSE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_wallets_email ON wallets(email);
	CREATE INDEX IF NOT EXISTS idx_wallets_wallet_id ON wallets(wallet_id);
	CREATE INDEX IF NOT EXISTS idx_reset_tokens_token ON reset_tokens(token);
	CREATE INDEX IF NOT EXISTS idx_reset_tokens_email ON reset_tokens(email);
	`

	_, err := d.db.Exec(query)
	return err
}

// CreateWallet cria uma nova carteira no banco de dados
func (d *Database) CreateWallet(wallet *models.Wallet) error {
	query := `
	INSERT INTO wallets (email, password, wallet_id, admin_key, invoice_key, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	_, err := d.db.Exec(query, wallet.Email, wallet.Password, wallet.WalletID, wallet.AdminKey, wallet.InvoiceKey, now, now)
	if err != nil {
		return fmt.Errorf("erro ao criar carteira: %w", err)
	}

	return nil
}

// GetWalletByEmail obtém uma carteira pelo email
func (d *Database) GetWalletByEmail(email string) (*models.Wallet, error) {
	query := `
	SELECT id, email, password, wallet_id, admin_key, invoice_key, created_at, updated_at
	FROM wallets WHERE email = ?
	`

	wallet := &models.Wallet{}
	err := d.db.QueryRow(query, email).Scan(
		&wallet.ID,
		&wallet.Email,
		&wallet.Password,
		&wallet.WalletID,
		&wallet.AdminKey,
		&wallet.InvoiceKey,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar carteira: %w", err)
	}

	return wallet, nil
}

// GetWalletByWalletID obtém uma carteira pelo wallet_id
func (d *Database) GetWalletByWalletID(walletID string) (*models.Wallet, error) {
	query := `
	SELECT id, email, password, wallet_id, admin_key, invoice_key, created_at, updated_at
	FROM wallets WHERE wallet_id = ?
	`

	wallet := &models.Wallet{}
	err := d.db.QueryRow(query, walletID).Scan(
		&wallet.ID,
		&wallet.Email,
		&wallet.Password,
		&wallet.WalletID,
		&wallet.AdminKey,
		&wallet.InvoiceKey,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar carteira: %w", err)
	}

	return wallet, nil
}

// WalletExists verifica se uma carteira existe pelo email
func (d *Database) WalletExists(email string) (bool, error) {
	query := `SELECT COUNT(*) FROM wallets WHERE email = ?`
	
	var count int
	err := d.db.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("erro ao verificar existência da carteira: %w", err)
	}

	return count > 0, nil
}

// CreateResetToken cria um token de reset de senha
func (d *Database) CreateResetToken(email, token string, expiresAt time.Time) error {
	query := `
	INSERT INTO reset_tokens (email, token, expires_at, created_at)
	VALUES (?, ?, ?, ?)
	`

	_, err := d.db.Exec(query, email, token, expiresAt, time.Now())
	if err != nil {
		return fmt.Errorf("erro ao criar token de reset: %w", err)
	}

	return nil
}

// GetResetToken obtém um token de reset válido
func (d *Database) GetResetToken(token string) (string, time.Time, bool, error) {
	query := `
	SELECT email, expires_at, used
	FROM reset_tokens 
	WHERE token = ? AND used = FALSE AND expires_at > ?
	`

	var email string
	var expiresAt time.Time
	var used bool

	err := d.db.QueryRow(query, token, time.Now()).Scan(&email, &expiresAt, &used)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", time.Time{}, false, nil
		}
		return "", time.Time{}, false, fmt.Errorf("erro ao buscar token de reset: %w", err)
	}

	return email, expiresAt, used, nil
}

// MarkResetTokenAsUsed marca um token como usado
func (d *Database) MarkResetTokenAsUsed(token string) error {
	query := `UPDATE reset_tokens SET used = TRUE WHERE token = ?`

	_, err := d.db.Exec(query, token)
	if err != nil {
		return fmt.Errorf("erro ao marcar token como usado: %w", err)
	}

	return nil
}

// UpdatePassword atualiza a senha de um usuário
func (d *Database) UpdatePassword(email, newPassword string) error {
	query := `UPDATE wallets SET password = $1, updated_at = $2 WHERE email = $3`

	_, err := d.db.Exec(query, newPassword, time.Now(), email)
	if err != nil {
		return fmt.Errorf("erro ao atualizar senha: %w", err)
	}

	return nil
}

// CleanExpiredTokens remove tokens de reset expirados
func (d *Database) CleanExpiredTokens() (int64, error) {
	query := `DELETE FROM reset_tokens WHERE expires_at < $1`

	result, err := d.db.Exec(query, time.Now())
	if err != nil {
		return 0, fmt.Errorf("erro ao limpar tokens expirados: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("erro ao obter número de linhas afetadas: %w", err)
	}

	return rowsAffected, nil
}

// GetExpiredTokensCount retorna o número de tokens expirados
func (d *Database) GetExpiredTokensCount() (int64, error) {
	query := `SELECT COUNT(*) FROM reset_tokens WHERE expires_at < $1`

	var count int64
	err := d.db.QueryRow(query, time.Now()).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("erro ao contar tokens expirados: %w", err)
	}

	return count, nil
}

// Close fecha a conexão com o banco de dados
func (d *Database) Close() error {
	log.Println("Fechando conexão com banco de dados...")
	return d.db.Close()
}
