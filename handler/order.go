package handler

import (
	"intern-project-v2/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type orderHandler struct {
	orderUsecase domain.OrderUsecase
}

func NewOrderHandler(orderUsecase domain.OrderUsecase) *orderHandler {
	return &orderHandler{
		orderUsecase: orderUsecase,
	}
}

// GetAll godoc
// @Summary Get all orders
// @Description Retrieve all orders
// @Tags Orders
// @Accept json
// @Produce json
// @Success 200 {array} domain.Order
// @Failure 500
// @Failure 400
// @Router /orders [get]
func (oh *orderHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()
	orders, err := oh.orderUsecase.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve orders"})
		return
	}
	if len(orders) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No orders found"})
		return
	}
	c.JSON(http.StatusOK,
		gin.H{
			"message": "Orders retrieved successfully",
			"orders":  orders})
}

// GetByID godoc
// @Summary Get order by ID
// @Description Retrieve an order by its ID
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} domain.Order
// @Failure 400
// @Failure 500
// @Router /orders/{id} [get]
func (oh *orderHandler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()
	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
		return
	}
	order, err := oh.orderUsecase.GetByID(ctx, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve order", "details": err.Error()})
		return
	}
	if order == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Order not found"})
		return
	}
	c.JSON(http.StatusOK,
		gin.H{
			"message": "Order retrieved successfully",
			"order":   order,
		})
}

// Create godoc
// @Summary Create a new order
// @Description Create a new order with the provided details
// @Tags Orders
// @Accept json
// @Produce json
// @Param order body domain.OrderRequest true "Order Request"
// @Success 201 {object} domain.Order
// @Failure 400
// @Failure 500
// @Router /orders [post]
func (oh *orderHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var orderReq domain.OrderRequest
	if err := c.ShouldBindJSON(&orderReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	order, err := oh.orderUsecase.Create(ctx, &orderReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order", "details": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, order)
}

// Update godoc
// @Summary Update an existing order
// @Description Update an order with the provided details
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param order body domain.OrderRequest true "Order Request"
// @Success 200 {object} domain.Order
// @Failure 400
// @Failure 500
// @Router /orders/{id} [put]
func (oh *orderHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()
	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
		return
	}
	var orderReq domain.OrderRequest
	if err := c.ShouldBindJSON(&orderReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	order, err := oh.orderUsecase.Update(ctx, orderID, &orderReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order", "details": err.Error()})
		return
	}
	if order == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Order not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order updated successfully", "order": order})
}

// Delete godoc
// @Summary Delete an order
// @Description Delete an order by its ID
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /orders/{id} [delete]
func (oh *orderHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
		return
	}
	order, err := oh.orderUsecase.Delete(ctx, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order", "details": err.Error()})
		return
	}
	if order == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Order not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully", "order": order})
}
