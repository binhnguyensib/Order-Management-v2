package handler

import (
	"context"
	"intern-project-v2/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderUsecase interface {
	GetAll(ctx context.Context) ([]*domain.Order, error)
	GetByID(ctx context.Context, id string) (*domain.Order, error)
	Create(ctx context.Context, order *domain.OrderRequest) (*domain.Order, error)
	Update(ctx context.Context, id string, orderReq *domain.OrderRequest) (*domain.Order, error)
	Delete(ctx context.Context, id string) (*domain.Order, error)
}

type OrderHandler struct {
	orderUsecase OrderUsecase
}

func NewOrderHandler(orderUsecase OrderUsecase) *OrderHandler {
	return &OrderHandler{
		orderUsecase: orderUsecase,
	}
}

func (oh *OrderHandler) GetAll(c *gin.Context) {
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

func (oh *OrderHandler) GetByID(c *gin.Context) {
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

func (oh *OrderHandler) Create(c *gin.Context) {
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

func (oh *OrderHandler) Update(c *gin.Context) {
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

func (oh *OrderHandler) Delete(c *gin.Context) {
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
