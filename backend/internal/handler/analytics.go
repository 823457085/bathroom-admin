package handler

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct {
	db *sql.DB
}

func NewAnalyticsHandler(db *sql.DB) *AnalyticsHandler {
	return &AnalyticsHandler{db: db}
}

type DashboardStats struct {
	TotalOrders    int     `json:"total_orders"`
	TodayOrders   int     `json:"today_orders"`
	TotalRevenue  float64 `json:"total_revenue"`
	TodayRevenue  float64 `json:"today_revenue"`
	TotalUsers    int     `json:"total_users"`
	TodayUsers    int     `json:"today_users"`
	TotalProducts int     `json:"total_products"`
}

func (h *AnalyticsHandler) Dashboard(c *gin.Context) {
	var stats DashboardStats

	h.db.QueryRow("SELECT COUNT(*) FROM orders").Scan(&stats.TotalOrders)
	h.db.QueryRow("SELECT COUNT(*) FROM orders WHERE DATE(created_at) = CURDATE()").Scan(&stats.TodayOrders)
	h.db.QueryRow("SELECT COALESCE(SUM(total_amount), 0) FROM orders").Scan(&stats.TotalRevenue)
	h.db.QueryRow("SELECT COALESCE(SUM(total_amount), 0) FROM orders WHERE DATE(created_at) = CURDATE()").Scan(&stats.TodayRevenue)
	h.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&stats.TotalUsers)
	h.db.QueryRow("SELECT COUNT(*) FROM users WHERE DATE(created_at) = CURDATE()").Scan(&stats.TodayUsers)
	h.db.QueryRow("SELECT COUNT(*) FROM products").Scan(&stats.TotalProducts)

	c.JSON(http.StatusOK, stats)
}

func (h *AnalyticsHandler) SalesTrend(c *gin.Context) {
	rows, err := h.db.Query(`
		SELECT DATE(created_at) as date, COUNT(*) as count, COALESCE(SUM(total_amount), 0) as amount
		FROM orders
		WHERE created_at >= DATE_SUB(CURDATE(), INTERVAL 30 DAY)
		GROUP BY DATE(created_at)
		ORDER BY date ASC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get trend"})
		return
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var date string
		var count int
		var amount float64
		rows.Scan(&date, &count, &amount)
		result = append(result, map[string]interface{}{"date": date, "count": count, "amount": amount})
	}
	c.JSON(http.StatusOK, result)
}

func (h *AnalyticsHandler) TopProducts(c *gin.Context) {
	rows, err := h.db.Query(`
		SELECT oi.product_id, oi.product_name, COUNT(*) as order_count
		FROM order_items oi
		JOIN orders o ON oi.order_id = o.id
		GROUP BY oi.product_id, oi.product_name
		ORDER BY order_count DESC
		LIMIT 10
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get top products"})
		return
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var pid int64
		var name string
		var cnt int
		rows.Scan(&pid, &name, &cnt)
		result = append(result, map[string]interface{}{"product_id": pid, "product_name": name, "order_count": cnt})
	}
	c.JSON(http.StatusOK, result)
}
