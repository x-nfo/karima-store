package swagger

import "github.com/swaggo/swag"

// APIInfo holds the API information
var APIInfo = &swag.Spec{
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
}

// SwaggerInfo holds the swagger info
var SwaggerInfo = &swag.Spec{
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