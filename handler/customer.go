package handler

import (
	"context"
	"intern-project-v2/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CustomerUsecase interface {
	GetAll(ctx context.Context) ([]*domain.Customer, error)
	GetByID(ctx context.Context, id string) (*domain.Customer, error)
	Create(ctx context.Context, customer *domain.CustomerRequest) (*domain.Customer, error)
	Update(ctx context.Context, id string, customerReq *domain.CustomerRequest) (*domain.Customer, error)
	Delete(ctx context.Context, id string) (*domain.Customer, error)
}

type CustomerHandler struct {
	customerUsecase CustomerUsecase
}

func NewCustomerHandler(customerUsecase CustomerUsecase) *CustomerHandler {
	return &CustomerHandler{
		customerUsecase: customerUsecase,
	}

}

func (ch *CustomerHandler) GetAll(c *gin.Context) {
	customers, err := ch.customerUsecase.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve customers"})
		return
	}

	if len(customers) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No customers found"})
		return
	}

	c.JSON(http.StatusOK, customers)
}

func (ch *CustomerHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	customer, err := ch.customerUsecase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve customer",
			"details": err.Error()})
		return
	}

	if customer == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Customer not found"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

func (ch *CustomerHandler) Create(c *gin.Context) {
	var customerReq domain.CustomerRequest
	if err := c.ShouldBindJSON(&customerReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	customer, err := ch.customerUsecase.Create(c.Request.Context(), &customerReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create customer",
			"details": err.Error()})
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Customer created successfully",
		"customer": customer})
}

func (ch *CustomerHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	var customerReq domain.CustomerRequest
	if err := c.ShouldBindJSON(&customerReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if customerReq.Name == "" && customerReq.Email == "" && customerReq.Phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one field must be provided for update"})
		return
	}

	customer, err := ch.customerUsecase.Update(c.Request.Context(), id, &customerReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update customer",
			"details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Customer updated successfully",
		"customer": customer})
}

func (ch *CustomerHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	customer, err := ch.customerUsecase.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete customer",
			"details": err.Error()})
		return
	}

	if customer == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Customer not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Customer deleted successfully",
		"customer": customer})
}
