package api_errors

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	db "github.com/adedaryorh/ecommerceapi/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Product struct {
	server *Server
}

// Set up routes for product-related APIs.
func (p *Product) router(server *Server) {
	p.server = server

	serverGroup := server.router.Group("/products", server.AuthenticatedMiddleware(), RoleBasedMiddleware(server, "admin"))
	serverGroup.POST("/createProduct", p.createProduct)
	serverGroup.GET("/:id", p.getProduct)
	serverGroup.GET("", p.listProducts)
	serverGroup.PUT("/:id", p.updateProduct)
	serverGroup.DELETE("/:id", p.deleteProduct)
}

// ProductParams defines the expected input for product operations.
type ProductParams struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"` // Nullable description
	Price       string  `json:"price" binding:"required,gt=0"`
	Stock       int32   `json:"stock" binding:"required,gt=0"`
}

// ProductResponse defines the response structure for product data.
type ProductResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"` // Nullable description
	Price       float64   `json:"price"`       // Now a float64
	Stock       int32     `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Converts a db.Product to a ProductResponse.
// Converts a db.Product to a ProductResponse.
func (p ProductResponse) toProductResponse(product *db.Product) ProductResponse {
	var description *string
	if product.Description.Valid {
		description = &product.Description.String
	}
	// Convert the product price from string to float64
	price, err := strconv.ParseFloat(product.Price, 64)
	if err != nil {
		// Handle error gracefully if conversion fails
		fmt.Println("Error converting price:", err)
		price = 0.0 // Set default value in case of error
	}
	// Accessing the price as a string, not Float64
	return ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: description,
		Price:       price,
		Stock:       product.Stock,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}

// @Summary Create Product
// @Description Create a new product (admin only)
// @Tags Products
// @Accept json
// @Produce json
// @Param product body ProductParams true "Product Details"
// @Success 201 {object} api_errors.ProductResponse
// @Failure 400 {object} api_errors.ApiError
// @Failure 500 {object} api_errors.ApiError
// @Security BearerAuth
// @Router /products/createProduct [post]
func (p *Product) createProduct(c *gin.Context) {
	var params ProductParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert the string price to float64
	price, err := strconv.ParseFloat(params.Price, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price format"})
		return
	}

	arg := db.CreateProductParams{
		Name:        params.Name,
		Description: sql.NullString{String: *params.Description, Valid: params.Description != nil},
		Price:       price, // Pass the converted price
		Stock:       params.Stock,
	}

	product, err := p.server.queries.CreateProduct(context.Background(), arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ProductResponse{}.toProductResponse(&product))
}

// @Summary Get Product
// @Description Retrieve a specific product by ID
// @Tags Products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} api_errors.ProductResponse
// @Failure 400 {object} api_errors.ApiError
// @Failure 404 {object} api_errors.ApiError
// @Failure 500 {object} api_errors.ApiError
// @Security BearerAuth
// @Router /products/{id} [get]
func (p *Product) getProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	product, err := p.server.queries.GetProductByID(context.Background(), id)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, ProductResponse{}.toProductResponse(&product))
}

// @Summary List Products
// @Description Retrieve paginated list of products
// @Tags Products
// @Produce json
// @Param limit query int false "Number of products to retrieve" default(10)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} api_errors.ProductResponse
// @Failure 500 {object} api_errors.ApiError
// @Security BearerAuth
// @Router /products [get]
func (p *Product) listProducts(c *gin.Context) {
	limit := int32(10)
	offset := int32(0)

	if l, err := strconv.Atoi(c.Query("limit")); err == nil {
		limit = int32(l)
	}
	if o, err := strconv.Atoi(c.Query("offset")); err == nil {
		offset = int32(o)
	}

	arg := db.ListProductsParams{
		Limit:  limit,
		Offset: offset,
	}

	products, err := p.server.queries.ListProducts(context.Background(), arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list products: " + err.Error()})
		return
	}

	response := []ProductResponse{}
	for _, product := range products {
		response = append(response, ProductResponse{}.toProductResponse(&product))
	}

	c.JSON(http.StatusOK, gin.H{"products": response, "limit": limit, "offset": offset})
}

// @Summary Update Product
// @Description Update an existing product (admin only)
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body ProductParams true "Product Update Details"
// @Success 200 {object} api_errors.ProductResponse
// @Failure 400 {object} api_errors.ApiError
// @Failure 404 {object} api_errors.ApiError
// @Failure 500 {object} api_errors.ApiError
// @Security BearerAuth
// @Router /products/{id} [put]
func (p *Product) updateProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var params ProductParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert the string price to float64
	price, err := strconv.ParseFloat(params.Price, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price format"})
		return
	}

	arg := db.UpdateProductParams{
		Name:        params.Name,
		Description: sql.NullString{String: *params.Description, Valid: params.Description != nil},
		Price:       price, // Pass the converted price
		Stock:       params.Stock,
		UpdatedAt:   time.Now(),
		ID:          id,
	}

	product, err := p.server.queries.UpdateProduct(context.Background(), arg)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, ProductResponse{}.toProductResponse(&product))
}

// @Summary Delete Product
// @Description Delete a product by ID (admin only)
// @Tags Products
// @Param id path string true "Product ID"
// @Success 200 {object} api_errors.ApiError
// @Failure 400 {object} api_errors.ApiError
// @Failure 404 {object} api_errors.ApiError
// @Failure 500 {object} api_errors.ApiError
// @Security BearerAuth
// @Router /products/{id} [delete]
func (p *Product) deleteProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}
	err = p.server.queries.DeleteProduct(context.Background(), id)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
