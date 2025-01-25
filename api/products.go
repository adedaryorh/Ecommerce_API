package api

func (p Product) router(server *Server) {
	p.server = server

	// Protect these routes with Role-Based Middleware (admin only)
	serverGroup := server.router.Group("/products", u.server.AuthenticatedMiddleware(), RoleBasedMiddleware("admin"))
	serverGroup.POST("", p.createProduct) // Admin only
	serverGroup.GET("/:id", p.getProduct)
	serverGroup.PUT("/:id", p.updateProduct)    // Admin only
	serverGroup.DELETE("/:id", p.deleteProduct) // Admin only
}
