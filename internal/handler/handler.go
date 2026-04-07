package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rchaser53/fx-data-analysis/internal/database"
	"github.com/rchaser53/fx-data-analysis/internal/model"
)

const usdjpyDataDir = "./data/usdjpy"

const (
	usdJPYTimeframeDaily  = "daily"
	usdJPYTimeframeWeekly = "weekly"
)

type Handler struct {
	db *database.DB
}

// NewHandler creates a new handler
func NewHandler(db *database.DB) *Handler {
	return &Handler{db: db}
}

// CreateTrade creates a new trade
func (h *Handler) CreateTrade(c *gin.Context) {
	var req model.CreateTradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trade, err := h.db.CreateTrade(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, trade)
}

// GetTrade retrieves a trade by ID
func (h *Handler) GetTrade(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trade ID"})
		return
	}

	trade, err := h.db.GetTrade(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, trade)
}

// GetAllTrades retrieves all trades
func (h *Handler) GetAllTrades(c *gin.Context) {
	trades, err := h.db.GetAllTrades()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, trades)
}

// UpdateTrade updates an existing trade
func (h *Handler) UpdateTrade(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trade ID"})
		return
	}

	var req model.UpdateTradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trade, err := h.db.UpdateTrade(id, &req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, trade)
}

// DeleteTrade deletes a trade by ID
func (h *Handler) DeleteTrade(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trade ID"})
		return
	}

	if err := h.db.DeleteTrade(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "trade deleted successfully"})
}

// GetUSDJPYRates retrieves all USDJPY rate files and returns them sorted by date
func (h *Handler) GetUSDJPYRates(c *gin.Context) {
	timeframe := normalizeUSDJPYTimeframe(c.DefaultQuery("timeframe", usdJPYTimeframeDaily))
	if timeframe == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid timeframe"})
		return
	}

	filePaths, err := filepath.Glob(filepath.Join(usdjpyDataDir, "*.json"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list rate files"})
		return
	}

	rates := make([]model.USDJPYRate, 0, len(filePaths))
	for _, filePath := range filePaths {
		body, err := os.ReadFile(filePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read rate file: " + filepath.Base(filePath)})
			return
		}

		var rate model.USDJPYRate
		if err := json.Unmarshal(body, &rate); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse rate file: " + filepath.Base(filePath)})
			return
		}

		rates = append(rates, rate)
	}

	sort.Slice(rates, func(i, j int) bool {
		return rates[i].Date < rates[j].Date
	})

	rates, err = filterUSDJPYTradingDays(rates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to filter rate dates"})
		return
	}

	rates = withDailyLabels(rates)
	if timeframe == usdJPYTimeframeWeekly {
		rates, err = aggregateUSDJPYRatesByWeek(rates)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to aggregate weekly rates"})
			return
		}
	}

	c.JSON(http.StatusOK, model.USDJPYRatesResponse{
		Pair:      "USDJPY",
		Timeframe: timeframe,
		Rates:     rates,
	})
}

func normalizeUSDJPYTimeframe(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", usdJPYTimeframeDaily:
		return usdJPYTimeframeDaily
	case usdJPYTimeframeWeekly:
		return usdJPYTimeframeWeekly
	default:
		return ""
	}
}

func withDailyLabels(rates []model.USDJPYRate) []model.USDJPYRate {
	if len(rates) == 0 {
		return rates
	}

	result := make([]model.USDJPYRate, len(rates))
	for i, rate := range rates {
		rate.Label = rate.Date
		result[i] = rate
	}

	return result
}

func filterUSDJPYTradingDays(rates []model.USDJPYRate) ([]model.USDJPYRate, error) {
	if len(rates) == 0 {
		return rates, nil
	}

	filtered := make([]model.USDJPYRate, 0, len(rates))
	for _, rate := range rates {
		parsedDate, err := time.Parse("2006-01-02", rate.Date)
		if err != nil {
			return nil, err
		}

		switch parsedDate.Weekday() {
		case time.Saturday, time.Sunday:
			continue
		default:
			filtered = append(filtered, rate)
		}
	}

	return filtered, nil
}

func aggregateUSDJPYRatesByWeek(rates []model.USDJPYRate) ([]model.USDJPYRate, error) {
	if len(rates) == 0 {
		return rates, nil
	}

	type weekKey struct {
		year int
		week int
	}

	grouped := make([][]model.USDJPYRate, 0)
	currentKey := weekKey{}
	currentGroup := make([]model.USDJPYRate, 0)

	for _, rate := range rates {
		parsedDate, err := time.Parse("2006-01-02", rate.Date)
		if err != nil {
			return nil, err
		}

		year, week := parsedDate.ISOWeek()
		key := weekKey{year: year, week: week}
		if len(currentGroup) == 0 || key == currentKey {
			currentKey = key
			currentGroup = append(currentGroup, rate)
			continue
		}

		grouped = append(grouped, currentGroup)
		currentKey = key
		currentGroup = []model.USDJPYRate{rate}
	}

	if len(currentGroup) > 0 {
		grouped = append(grouped, currentGroup)
	}

	weeklyRates := make([]model.USDJPYRate, 0, len(grouped))
	for _, group := range grouped {
		first := group[0]
		last := group[len(group)-1]
		weeklyRate := model.USDJPYRate{
			Date:  first.Date,
			Label: first.Date + " - " + last.Date,
			Pair:  first.Pair,
			Bid:   last.Bid,
			Ask:   last.Ask,
			Open:  first.Open,
			High:  first.High,
			Low:   first.Low,
			Close: last.Close,
		}

		for _, rate := range group[1:] {
			if rate.High > weeklyRate.High {
				weeklyRate.High = rate.High
			}
			if rate.Low < weeklyRate.Low {
				weeklyRate.Low = rate.Low
			}
		}

		weeklyRate.Diff = weeklyRate.Close - weeklyRate.Open
		weeklyRates = append(weeklyRates, weeklyRate)
	}

	return weeklyRates, nil
}
