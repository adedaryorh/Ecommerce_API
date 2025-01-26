package api_errors

import (
	"context"
	db "github.com/adedaryorh/ecommerceapi/db/sqlc"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Create Order
// @Description Create a new order with order items
// @Tags Orders
// @Accept json
// @Produce json
// @Param order body object true "Order Creation Details"
// @Success 201 {object} db.Order "Successfully created order" // Corrected reference
// @Failure 400 {object} api_errors.ApiError "Bad Request"
// @Failure 500 {object} api_errors.ApiError "Internal Server Error"
// @Security BearerAuth
// @Router /orders [post]
func (s *Server) CreateOrder(c *gin.Context) {
	var orderParams struct {
		UserID      int64  `json:"user_id"`
		TotalAmount string `json:"total_amount"`
		Status      string `json:"status"`
		OrderItems  []struct {
			ProductID int64  `json:"product_id"`
			Quantity  int32  `json:"quantity"`
			Price     string `json:"price"`
		} `json:"order_items"`
	}

	if err := c.ShouldBindJSON(&orderParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := s.queries.CreateOrder(context.Background(), db.CreateOrderParams{
		UserID:      orderParams.UserID,
		TotalAmount: orderParams.TotalAmount,
		Status:      orderParams.Status,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, item := range orderParams.OrderItems {
		_, err := s.queries.AddOrderItem(context.Background(), db.AddOrderItemParams{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusCreated, order)
}

// @Summary List User Orders
// @Description Retrieve paginated orders for the authenticated user
// @Tags Orders
// @Produce json
// @Param limit query int false "Number of orders to retrieve" default(10)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {array} db.Order "List of orders" // Corrected reference for array of db.Order
// @Failure 401 {object} api_errors.ApiError "Unauthorized"
// @Failure 500 {object} api_errors.ApiError "Internal Server Error"
// @Security BearerAuth
// @Router /orders [get]
func (s *Server) ListUserOrders(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userIDInt64, ok := userID.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}
	limit := 10
	offset := 0
	if val, exists := c.GetQuery("limit"); exists {
		limit, _ = strconv.Atoi(val)
	}
	if val, exists := c.GetQuery("offset"); exists {
		offset, _ = strconv.Atoi(val)
	}

	orders, err := s.queries.ListUserOrders(context.Background(), db.ListUserOrdersParams{
		UserID: userIDInt64,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// @Summary Cancel Order
// @Description Cancel an order (admin only)
// @Tags Orders
// @Param id path string true "Order ID"
// @Success 200 {object} api_errors.ProductResponse "Success response" // No need to change this if ProductResponse is correctly defined
// @Failure 400 {object} api_errors.ApiError "Bad Request"
// @Failure 403 {object} api_errors.ApiError "Forbidden"
// @Failure 500 {object} api_errors.ApiError "Internal Server Error"
// @Security BearerAuth
// @Router /admin/orders/{id}/cancel [post]
func (s *Server) CancelOrder(c *gin.Context) {
	orderID := c.Param("id")

	// Verify if the user has admin privileges
	userID, _ := c.Get("user_id")
	user, err := s.queries.GetUserByID(context.Background(), userID.(int64))
	if err != nil || user.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can cancel orders"})
		return
	}

	// Convert orderID to int64 if required by your SQL method
	orderIDInt64, err := strconv.ParseInt(orderID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	// Cancel the order
	err = s.queries.CancelOrder(context.Background(), orderIDInt64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order cancelled successfully"})
}

// @Summary Update Order Status
// @Description Update the status of an order (admin only)
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param status body object true "New Order Status"
// @Success 200 {object} db.Order "Updated order status" // Corrected reference
// @Failure 400 {object} api_errors.ApiError "Bad Request"
// @Failure 403 {object} api_errors.ApiError "Forbidden"
// @Failure 500 {object} api_errors.ApiError "Internal Server Error"
// @Security BearerAuth
// @Router /admin/orders/{id}/status [patch]
func (s *Server) UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("id")

	var statusUpdate struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	user, err := s.queries.GetUserByID(context.Background(), userID.(int64))
	if err != nil || user.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can update order status"})
		return
	}

	orderIDInt64, err := strconv.ParseInt(orderID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	order, err := s.queries.UpdateOrderStatus(context.Background(), db.UpdateOrderStatusParams{
		Status: statusUpdate.Status,
		ID:     orderIDInt64,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}
