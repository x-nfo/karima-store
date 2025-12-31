package swagger

import "github.com/swaggo/swag"

// @title Karima Store API
// @version 1.0.0
// @description Karima Store - E-commerce API Documentation
// @termsOfService https://karima-store.com/terms
// @contact.name Karima Store Team
// @contact.email team@karima-store.com
// @contact.url https://karima-store.com
// @license.name MIT
// @license.url https://github.com/karima-store/karima-store/blob/main/LICENSE
// @host localhost:8080
// @BasePath /
// @schemes http https
// @securityDefinitions.BearerAuth BearerAuth
// @security.BearerAuth type:apiKey name:Authorization in:header description:JWT Bearer Token

// KarimaStoreAPI is the main struct for the API documentation
type KarimaStoreAPI struct{}

// APIInfo returns the API information
func (a *KarimaStoreAPI) APIInfo() *swag.Spec {
    return &swag.Spec{
        Version:          "1.0.0",
        Title:            "Karima Store API",
        Description:      "Karima Store - E-commerce API Documentation",
        TermsOfService:   "https://karima-store.com/terms",
        Contact: &swag.Contact{
            Name:      "Karima Store Team",
            Email:     "team@karima-store.com",
            URL:       "https://karima-store.com",
        },
        License: &swag.License{
            Name:      "MIT",
            URL:       "https://github.com/karima-store/karima-store/blob/main/LICENSE",
        },
        Host:             "localhost:8080",
        BasePath:         "/",
        Schemes:          []string{"http", "https"},
        Consumes:         []string{"application/json"},
        Produces:         []string{"application/json"},
        SecurityDefinitions: map[string]swag.SecurityScheme{
            "BearerAuth": {
                Type: "apiKey",
                Name: "Authorization",
                In:   "header",
                Description: "JWT Bearer Token",
            },
        },
    }
}

// GetSwagger returns the Swagger specification
func GetSwagger() *swag.Spec {
    return &swag.Spec{
        Info:             APIInfo,
        Paths:            map[string]swag.PathItem{},
        Definitions:      map[string]swag.Schema{},
        SecurityDefinitions: map[string]swag.SecurityScheme{
            "BearerAuth": {
                Type: "apiKey",
                Name: "Authorization",
                In:   "header",
                Description: "JWT Bearer Token",
            },
        },
    }
}