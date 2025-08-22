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
func NewDatabase(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
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

	CREATE INDEX IF NOT EXISTS idx_wallets_email ON wallets(email);
	CREATE INDEX IF NOT EXISTS idx_wallets_wallet_id ON wallets(wallet_id);
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

// Close fecha a conexão com o banco de dados
func (d *Database) Close() error {
	log.Println("Fechando conexão com banco de dados...")
	return d.db.Close()
}
