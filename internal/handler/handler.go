package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rchaser53/fx-data-analysis/internal/database"
	"github.com/rchaser53/fx-data-analysis/internal/model"
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
