package pagination

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	defaultPage  = 1
	defaultLimit = 20
	maxLimit     = 100
)

type Params struct {
	Page  int
	Limit int
}

// Offset вычисляет смещение для SQL LIMIT/OFFSET.
func (p Params) Offset() int {
	return (p.Page - 1) * p.Limit
}

// Response — конверт для любого списка с метаданными пагинации.
type Response struct {
	Data       any `json:"data"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalPages int   `json:"total_pages"`
}

func NewResponse(data any, total int64, p Params) Response {
	totalPages := int(total) / p.Limit
	if int(total)%p.Limit != 0 {
		totalPages++
	}
	return Response{
		Data:       data,
		Total:      total,
		Page:       p.Page,
		Limit:      p.Limit,
		TotalPages: totalPages,
	}
}

// FromContext читает page и limit из gin.Context query-параметров.
func FromContext(c *gin.Context) Params {
	page := parsePositiveInt(c.Query("page"), defaultPage)
	limit := parsePositiveInt(c.Query("limit"), defaultLimit)

	if limit > maxLimit {
		limit = maxLimit
	}

	return Params{Page: page, Limit: limit}
}

func parsePositiveInt(s string, fallback int) int {
	if s == "" {
		return fallback
	}
	v, err := strconv.Atoi(s)
	if err != nil || v < 1 {
		return fallback
	}
	return v
}
