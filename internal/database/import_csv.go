package database

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/rchaser53/fx-data-analysis/internal/model"
)

// ImportTradesFromCSV reads a CSV file and inserts trades into the database
func (db *DB) ImportTradesFromCSV(csvPath string) error {
	file, err := os.Open(csvPath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	r := csv.NewReader(file)
	lineNum := 0
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read CSV: %w", err)
		}
		lineNum++
		if len(record) < 3 {
			return fmt.Errorf("invalid record at line %d: %v", lineNum, record)
		}

		// Parse trade_time
		tradeTime, err := time.Parse(time.RFC3339, record[0])
		if err != nil {
			return fmt.Errorf("invalid trade_time at line %d: %v", lineNum, err)
		}

		// Parse lot_size
		lotSize, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return fmt.Errorf("invalid lot_size at line %d: %v", lineNum, err)
		}

		// Parse purchase_rate
		purchaseRate, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return fmt.Errorf("invalid purchase_rate at line %d: %v", lineNum, err)
		}

		trade := &model.CreateTradeRequest{
			TradeTime:    tradeTime,
			LotSize:      lotSize,
			PurchaseRate: purchaseRate,
		}

		_, err = db.CreateTrade(trade)
		if err != nil {
			return fmt.Errorf("failed to insert trade at line %d: %v", lineNum, err)
		}
	}

	return nil
}
