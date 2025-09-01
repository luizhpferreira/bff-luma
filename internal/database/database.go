package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"bff-luma/internal/models"

	_ "github.com/lib/pq"
)

// Database representa a conexão com o banco de dados
type Database struct {
	db *sql.DB
}

// NewDatabase cria uma nova conexão com o banco de dados
func NewDatabase(dbURL string) (*Database, error) {
	db, err := sql.Open("postgres", dbURL)
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
		id SERIAL PRIMARY KEY,
		cpf TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		wallet_id TEXT NOT NULL UNIQUE,
		admin_key TEXT NOT NULL,
		invoice_key TEXT NOT NULL,
		email_confirmed BOOLEAN DEFAULT FALSE,
		email_confirmation_token TEXT,
		email_confirmation_expires_at TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS reset_tokens (
		id SERIAL PRIMARY KEY,
		email TEXT NOT NULL,
		token TEXT NOT NULL UNIQUE,
		expires_at TIMESTAMP NOT NULL,
		used BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS temp_passwords (
		id SERIAL PRIMARY KEY,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		expires_at TIMESTAMP NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_wallets_cpf ON wallets(cpf);
	CREATE INDEX IF NOT EXISTS idx_wallets_email ON wallets(email);
	CREATE INDEX IF NOT EXISTS idx_wallets_wallet_id ON wallets(wallet_id);
	CREATE INDEX IF NOT EXISTS idx_reset_tokens_token ON reset_tokens(token);
	CREATE INDEX IF NOT EXISTS idx_reset_tokens_email ON reset_tokens(email);
	CREATE INDEX IF NOT EXISTS idx_temp_passwords_email ON temp_passwords(email);
	`

	_, err := d.db.Exec(query)
	return err
}

// CreateWallet cria uma nova carteira no banco de dados
func (d *Database) CreateWallet(wallet *models.Wallet) error {
	query := `
	INSERT INTO wallets (cpf, email, password, wallet_id, admin_key, invoice_key, email_confirmed, email_confirmation_token, email_confirmation_expires_at, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	now := time.Now()
	_, err := d.db.Exec(query, 
		wallet.CPF, 
		wallet.Email, 
		wallet.Password, 
		wallet.WalletID, 
		wallet.AdminKey, 
		wallet.InvoiceKey, 
		wallet.EmailConfirmed,
		wallet.EmailConfirmationToken,
		wallet.EmailConfirmationExpiresAt,
		now, 
		now,
	)
	if err != nil {
		return fmt.Errorf("erro ao criar carteira: %w", err)
	}

	return nil
}

// CreateEmailConfirmationToken cria um token de confirmação de email
func (d *Database) CreateEmailConfirmationToken(email, token string, expiresAt time.Time) error {
	query := `
	UPDATE wallets 
	SET email_confirmation_token = $1, email_confirmation_expires_at = $2, updated_at = $3
	WHERE email = $4
	`

	_, err := d.db.Exec(query, token, expiresAt, time.Now(), email)
	if err != nil {
		return fmt.Errorf("erro ao criar token de confirmação de email: %w", err)
	}

	return nil
}

// ConfirmEmail confirma o email usando o token
func (d *Database) ConfirmEmail(token string) error {
	query := `
	UPDATE wallets 
	SET email_confirmed = TRUE, email_confirmation_token = NULL, email_confirmation_expires_at = NULL, updated_at = $1
	WHERE email_confirmation_token = $2 AND email_confirmation_expires_at > $3
	`

	result, err := d.db.Exec(query, time.Now(), token, time.Now())
	if err != nil {
		return fmt.Errorf("erro ao confirmar email: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("token inválido ou expirado")
	}

	return nil
}

// GetWalletByEmailConfirmationToken obtém uma carteira pelo token de confirmação
func (d *Database) GetWalletByEmailConfirmationToken(token string) (*models.Wallet, error) {
	query := `
	SELECT id, cpf, email, password, wallet_id, admin_key, invoice_key, email_confirmed, email_confirmation_token, email_confirmation_expires_at, created_at, updated_at
	FROM wallets WHERE email_confirmation_token = $1
	`

	wallet := &models.Wallet{}
	err := d.db.QueryRow(query, token).Scan(
		&wallet.ID,
		&wallet.CPF,
		&wallet.Email,
		&wallet.Password,
		&wallet.WalletID,
		&wallet.AdminKey,
		&wallet.InvoiceKey,
		&wallet.EmailConfirmed,
		&wallet.EmailConfirmationToken,
		&wallet.EmailConfirmationExpiresAt,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar carteira por token: %w", err)
	}

	return wallet, nil
}

// GetWalletByEmailAndTokenHistory verifica se uma carteira foi confirmada recentemente com um token específico
func (d *Database) GetWalletByEmailAndTokenHistory(email, token string) (*models.Wallet, error) {
	// Busca a carteira pelo email e verifica se foi confirmada recentemente
	// (dentro das últimas 24 horas) e se o token corresponde
	query := `
	SELECT id, cpf, email, password, wallet_id, admin_key, invoice_key, email_confirmed, email_confirmation_token, email_confirmation_expires_at, created_at, updated_at
	FROM wallets 
	WHERE email = $1 
	AND email_confirmed = true 
	AND updated_at > NOW() - INTERVAL '24 hours'
	`

	wallet := &models.Wallet{}
	err := d.db.QueryRow(query, email).Scan(
		&wallet.ID,
		&wallet.CPF,
		&wallet.Email,
		&wallet.Password,
		&wallet.WalletID,
		&wallet.AdminKey,
		&wallet.InvoiceKey,
		&wallet.EmailConfirmed,
		&wallet.EmailConfirmationToken,
		&wallet.EmailConfirmationExpiresAt,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar carteira por email e histórico: %w", err)
	}

	return wallet, nil
}

// GetWalletByEmail obtém uma carteira pelo email
func (d *Database) GetWalletByEmail(email string) (*models.Wallet, error) {
	query := `
	SELECT id, cpf, email, password, wallet_id, admin_key, invoice_key, email_confirmed, email_confirmation_token, email_confirmation_expires_at, created_at, updated_at
	FROM wallets WHERE email = $1
	`

	wallet := &models.Wallet{}
	err := d.db.QueryRow(query, email).Scan(
		&wallet.ID,
		&wallet.CPF,
		&wallet.Email,
		&wallet.Password,
		&wallet.WalletID,
		&wallet.AdminKey,
		&wallet.InvoiceKey,
		&wallet.EmailConfirmed,
		&wallet.EmailConfirmationToken,
		&wallet.EmailConfirmationExpiresAt,
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

// GetWalletByCPF obtém uma carteira pelo CPF
func (d *Database) GetWalletByCPF(cpf string) (*models.Wallet, error) {
	query := `
	SELECT id, cpf, email, password, wallet_id, admin_key, invoice_key, email_confirmed, email_confirmation_token, email_confirmation_expires_at, created_at, updated_at
	FROM wallets WHERE cpf = $1
	`

	wallet := &models.Wallet{}
	err := d.db.QueryRow(query, cpf).Scan(
		&wallet.ID,
		&wallet.CPF,
		&wallet.Email,
		&wallet.Password,
		&wallet.WalletID,
		&wallet.AdminKey,
		&wallet.InvoiceKey,
		&wallet.EmailConfirmed,
		&wallet.EmailConfirmationToken,
		&wallet.EmailConfirmationExpiresAt,
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
	SELECT id, cpf, email, password, wallet_id, admin_key, invoice_key, email_confirmed, email_confirmation_token, email_confirmation_expires_at, created_at, updated_at
	FROM wallets WHERE wallet_id = $1
	`

	wallet := &models.Wallet{}
	err := d.db.QueryRow(query, walletID).Scan(
		&wallet.ID,
		&wallet.CPF,
		&wallet.Email,
		&wallet.Password,
		&wallet.WalletID,
		&wallet.AdminKey,
		&wallet.InvoiceKey,
		&wallet.EmailConfirmed,
		&wallet.EmailConfirmationToken,
		&wallet.EmailConfirmationExpiresAt,
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

// WalletExists verifica se uma carteira existe pelo CPF
func (d *Database) WalletExists(cpf string) (bool, error) {
	query := `SELECT COUNT(*) FROM wallets WHERE cpf = $1`
	
	var count int
	err := d.db.QueryRow(query, cpf).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("erro ao verificar existência da carteira: %w", err)
	}

	return count > 0, nil
}

// WalletExistsByEmail verifica se uma carteira existe pelo email
func (d *Database) WalletExistsByEmail(email string) (bool, error) {
	query := `SELECT COUNT(*) FROM wallets WHERE email = $1`
	
	var count int
	err := d.db.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("erro ao verificar existência da carteira por email: %w", err)
	}

	return count > 0, nil
}

// CreateResetToken cria um token de reset de senha
func (d *Database) CreateResetToken(email, token string, expiresAt time.Time) error {
	query := `
	INSERT INTO reset_tokens (email, token, expires_at, created_at)
	VALUES ($1, $2, $3, $4)
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
	WHERE token = $1 AND used = FALSE AND expires_at > $2
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
	query := `UPDATE reset_tokens SET used = TRUE WHERE token = $1`

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

// UnconfirmEmail desconfirma o email de um usuário (usado em caso de erro no LNBits)
func (d *Database) UnconfirmEmail(email string) error {
	query := `
	UPDATE wallets 
	SET email_confirmed = FALSE, updated_at = $1
	WHERE email = $2
	`

	result, err := d.db.Exec(query, time.Now(), email)
	if err != nil {
		return fmt.Errorf("erro ao desconfirmar email: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("email não encontrado")
	}

	return nil
}

// UpdateWalletLNBitsData atualiza os dados do LNBits na carteira
func (d *Database) UpdateWalletLNBitsData(email, walletID, adminKey, invoiceKey string) error {
	query := `
	UPDATE wallets 
	SET wallet_id = $1, admin_key = $2, invoice_key = $3, updated_at = $4
	WHERE email = $5
	`

	result, err := d.db.Exec(query, walletID, adminKey, invoiceKey, time.Now(), email)
	if err != nil {
		return fmt.Errorf("erro ao atualizar dados do LNBits: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("carteira não encontrada para o email %s", email)
	}

	return nil
}

// SaveOriginalPassword salva a senha original temporariamente
func (d *Database) SaveOriginalPassword(email, password string) error {
	query := `
	INSERT INTO temp_passwords (email, password, expires_at)
	VALUES ($1, $2, $3)
	ON CONFLICT (email) DO UPDATE SET
		password = EXCLUDED.password,
		expires_at = EXCLUDED.expires_at,
		created_at = CURRENT_TIMESTAMP
	`

	expiresAt := time.Now().Add(1 * time.Hour)
	_, err := d.db.Exec(query, email, password, expiresAt)
	if err != nil {
		return fmt.Errorf("erro ao salvar senha original: %w", err)
	}

	return nil
}

// GetOriginalPassword obtém a senha original temporária
func (d *Database) GetOriginalPassword(email string) (string, error) {
	query := `
	SELECT password
	FROM temp_passwords 
	WHERE email = $1 AND expires_at > $2
	`

	var password string
	err := d.db.QueryRow(query, email, time.Now()).Scan(&password)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("senha original não encontrada ou expirada")
		}
		return "", fmt.Errorf("erro ao buscar senha original: %w", err)
	}

	return password, nil
}

// RemoveOriginalPassword remove a senha original temporária
func (d *Database) RemoveOriginalPassword(email string) error {
	query := `DELETE FROM temp_passwords WHERE email = $1`

	_, err := d.db.Exec(query, email)
	if err != nil {
		return fmt.Errorf("erro ao remover senha original: %w", err)
	}

	return nil
}

// CleanUnconfirmedAccounts remove contas não confirmadas após 24h
func (d *Database) CleanUnconfirmedAccounts() (int64, error) {
	query := `
	DELETE FROM wallets 
	WHERE email_confirmed = FALSE 
	AND created_at < NOW() - INTERVAL '24 hours'
	`

	result, err := d.db.Exec(query)
	if err != nil {
		return 0, fmt.Errorf("erro ao limpar contas não confirmadas: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("erro ao obter número de linhas afetadas: %w", err)
	}

	return rowsAffected, nil
}
