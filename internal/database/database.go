package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rchaser53/fx-data-analysis/internal/model"
)

type DB struct {
	*sql.DB
}

// NewDB initializes a new database connection
func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{db}, nil
}

// InitSchema creates the necessary tables
func (db *DB) InitSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS trades (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		trade_time DATETIME NOT NULL,
		lot_size REAL NOT NULL,
		purchase_rate REAL NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_trade_time ON trades(trade_time);
	`

	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}

	return nil
}

// CreateTrade inserts a new trade record
func (db *DB) CreateTrade(trade *model.CreateTradeRequest) (*model.Trade, error) {
	now := time.Now()
	query := `
		INSERT INTO trades (trade_time, lot_size, purchase_rate, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := db.Exec(query, trade.TradeTime, trade.LotSize, trade.PurchaseRate, now, now)
	if err != nil {
		return nil, fmt.Errorf("failed to create trade: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return &model.Trade{
		ID:           int(id),
		TradeTime:    trade.TradeTime,
		LotSize:      trade.LotSize,
		PurchaseRate: trade.PurchaseRate,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// GetTrade retrieves a trade by ID
func (db *DB) GetTrade(id int) (*model.Trade, error) {
	query := `
		SELECT id, trade_time, lot_size, purchase_rate, created_at, updated_at
		FROM trades
		WHERE id = ?
	`

	var trade model.Trade
	err := db.QueryRow(query, id).Scan(
		&trade.ID,
		&trade.TradeTime,
		&trade.LotSize,
		&trade.PurchaseRate,
		&trade.CreatedAt,
		&trade.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("trade not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get trade: %w", err)
	}

	return &trade, nil
}

// GetAllTrades retrieves all trades
func (db *DB) GetAllTrades() ([]model.Trade, error) {
	query := `
		SELECT id, trade_time, lot_size, purchase_rate, created_at, updated_at
		FROM trades
		ORDER BY trade_time DESC
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get trades: %w", err)
	}
	defer rows.Close()

	trades := make([]model.Trade, 0)
	for rows.Next() {
		var trade model.Trade
		err := rows.Scan(
			&trade.ID,
			&trade.TradeTime,
			&trade.LotSize,
			&trade.PurchaseRate,
			&trade.CreatedAt,
			&trade.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan trade: %w", err)
		}
		trades = append(trades, trade)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating trades: %w", err)
	}

	return trades, nil
}

// UpdateTrade updates an existing trade
func (db *DB) UpdateTrade(id int, req *model.UpdateTradeRequest) (*model.Trade, error) {
	// First, get the existing trade
	existingTrade, err := db.GetTrade(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.TradeTime != nil {
		existingTrade.TradeTime = *req.TradeTime
	}
	if req.LotSize != nil {
		existingTrade.LotSize = *req.LotSize
	}
	if req.PurchaseRate != nil {
		existingTrade.PurchaseRate = *req.PurchaseRate
	}
	existingTrade.UpdatedAt = time.Now()

	query := `
		UPDATE trades
		SET trade_time = ?, lot_size = ?, purchase_rate = ?, updated_at = ?
		WHERE id = ?
	`

	_, err = db.Exec(query,
		existingTrade.TradeTime,
		existingTrade.LotSize,
		existingTrade.PurchaseRate,
		existingTrade.UpdatedAt,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update trade: %w", err)
	}

	return existingTrade, nil
}

// DeleteTrade deletes a trade by ID
func (db *DB) DeleteTrade(id int) error {
	query := `DELETE FROM trades WHERE id = ?`

	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete trade: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("trade not found")
	}

	return nil
}
