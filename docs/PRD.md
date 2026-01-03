# Product Requirements Document (PRD)
## Karima Store - Fashion E-commerce Backend

---

## 📋 Table of Contents

1. [Executive Summary](#1-executive-summary)
2. [Product Overview](#2-product-overview)
3. [Target Audience](#3-target-audience)
4. [Core Features](#4-core-features)
5. [Technical Architecture](#5-technical-architecture)
6. [User Stories](#6-user-stories)
7. [Functional Requirements](#7-functional-requirements)
8. [Non-Functional Requirements](#8-non-functional-requirements)
9. [Integration Requirements](#9-integration-requirements)
10. [Security Requirements](#10-security-requirements)
11. [Performance Requirements](#11-performance-requirements)
12. [Deployment Architecture](#12-deployment-architecture)
13. [Success Metrics](#13-success-metrics)
14. [Roadmap](#14-roadmap)
15. [Appendices](#15-appendices)

---

## 1. Executive Summary

### 1.1 Purpose
This Product Requirements Document (PRD) defines the requirements for **Karima Store**, a modern, scalable e-commerce backend API designed for fashion retail businesses. The system provides a comprehensive foundation for building online stores with advanced features including product management, order processing, payment integration, shipping, and customer engagement.

### 1.2 Vision
To provide a robust, high-performance e-commerce backend that enables fashion retailers to launch and scale their online business with minimal technical overhead while maintaining enterprise-grade security and reliability.

### 1.3 Key Objectives
- **Scalability**: Support high traffic volumes and concurrent users
- **Reliability**: Ensure 99.9% uptime with robust error handling
- **Security**: Implement enterprise-grade security standards
- **Flexibility**: Modular architecture for easy customization
- **Performance**: Sub-100ms response times for critical APIs
- **Developer Experience**: Clean API design with comprehensive documentation

### 1.4 Current Status
- **Status**: In Progress (70% Complete)
- **Last Updated**: January 2, 2026
- **Version**: 1.0.0
- **Framework**: Go 1.24.0 + Fiber v2 + PostgreSQL + Redis

---

## 2. Product Overview

### 2.1 Problem Statement
Fashion retailers face challenges when building e-commerce platforms:
- Complex inventory management with variants (size, color, etc.)
- Dynamic pricing strategies (retail, reseller, flash sales)
- Multiple payment gateway integrations
- Complex shipping calculations across regions
- Real-time inventory synchronization
- Customer engagement through notifications
- Security and compliance requirements

### 2.2 Solution
Karima Store provides a comprehensive backend API that addresses these challenges:
- **Product Management**: Full CRUD with variant support
- **Dynamic Pricing Engine**: Retail, reseller, and flash sale pricing
- **Payment Integration**: Midtrans gateway with webhook support
- **Shipping Integration**: Komerce API for real-time shipping costs
- **Inventory Management**: Real-time stock tracking and logging
- **Authentication**: Ory Kratos for secure identity management
- **Notifications**: WhatsApp integration via Fonnte
- **Caching**: Redis-based performance optimization

### 2.3 Value Proposition
- **Time to Market**: Launch e-commerce store in weeks, not months
- **Cost Effective**: Open-source with minimal infrastructure costs
- **Scalable**: Built to handle growth from startup to enterprise
- **Secure**: Enterprise-grade security out of the box
- **Flexible**: Modular architecture for custom integrations

---

## 3. Target Audience

### 3.1 Primary Users
1. **E-commerce Business Owners**
   - Fashion retailers launching online stores
   - Multi-brand fashion marketplaces
   - Dropshipping businesses

2. **Developers**
   - Backend API developers
   - Frontend developers building store interfaces
   - DevOps engineers managing deployments

3. **Administrators**
   - Store managers managing inventory
   - Customer support handling orders
   - Marketing teams managing promotions

### 3.2 Secondary Users
1. **Customers**
   - End-users browsing and purchasing products
   - Mobile app users
   - Web store visitors

2. **Partners**
   - Payment gateway providers (Midtrans)
   - Shipping providers (Komerce)
   - Notification providers (Fonnte)

---

## 4. Core Features

### 4.1 Product Management
- **Product CRUD**: Create, read, update, delete products
- **Variant Management**: Support for size, color, and custom variants
- **SKU Management**: Unique SKU tracking for each variant
- **Media Management**: Upload and manage product images/videos
- **Category Management**: Hierarchical product categorization
- **Stock Management**: Real-time inventory tracking
- **Stock Logs**: Complete audit trail of stock movements

### 4.2 Pricing Engine
- **Retail Pricing**: Standard customer pricing
- **Reseller Pricing**: Tiered pricing for resellers
- **Flash Sale Pricing**: Time-limited promotional pricing
- **Coupon System**: Discount code validation and application
- **Tax Calculation**: Automatic tax calculation (11% VAT)
- **Free Shipping Thresholds**: Configurable free shipping rules

### 4.3 Order Management
- **Shopping Cart**: Add/update/remove cart items
- **Checkout Process**: Complete order flow with validation
- **Order Status Tracking**: Pending, confirmed, shipped, delivered, cancelled, refunded
- **Payment Status**: Pending, paid, failed, refunded
- **Order History**: Customer order history
- **Stock Management**: Automatic stock adjustment on order events

### 4.4 Payment Integration
- **Midtrans Gateway**: Secure payment processing
- **Snap Token**: Seamless checkout experience
- **Webhook Handling**: Real-time payment notifications
- **Refund Processing**: Automated refund handling
- **Signature Verification**: Secure webhook validation

### 4.5 Shipping Integration
- **Komerce API**: Real-time shipping cost calculation
- **Destination Search**: Search shipping destinations
- **Multiple Couriers**: Support for multiple shipping providers
- **Shipping Zones**: Configurable shipping regions
- **Weight-based Pricing**: Calculate costs based on total weight

### 4.6 Authentication & Authorization
- **Ory Kratos**: Enterprise identity management
- **User Registration**: Self-service user registration
- **User Login**: Secure authentication flow
- **Session Management**: Secure session handling
- **RBAC**: Role-Based Access Control (Planned)
- **API Key Authentication**: For external integrations

### 4.7 Notifications
- **WhatsApp Notifications**: Order confirmations and updates
- **Fonnte Integration**: Reliable message delivery
- **Async Processing**: Non-blocking notification delivery
- **Customizable Templates**: Flexible message formatting

### 4.8 Media Management
- **File Upload**: Support for images and videos
- **Cloudflare R2**: Cloud storage integration
- **Local Storage**: Development-friendly local storage
- **Primary Image**: Designate main product image
- **Soft Delete**: Safe media deletion

### 4.9 Caching & Performance
- **Redis Caching**: Cache-Aside pattern implementation
- **Product Catalog Caching**: Optimized product listing
- **Rate Limiting**: API abuse protection
- **Connection Pooling**: Database connection optimization

### 4.10 Security Features
- **CORS Protection**: Cross-origin resource sharing control
- **CSRF Protection**: Cross-site request forgery prevention
- **Rate Limiting**: DDoS protection
- **Security Headers**: HTTP security headers
- **Request Validation**: Input sanitization and validation
- **API Key Authentication**: Secure external API access
- **Request Tracing**: Audit logging

---

## 5. Technical Architecture

### 5.1 System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Client Layer                            │
├─────────────────────────────────────────────────────────────────┤
│  Customer Storefront  │  Admin Panel  │  Mobile App  │  Partners│
│  (Next.js/React)      │  (Next.js)    │  (React Native)        │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      API Gateway Layer                          │
├─────────────────────────────────────────────────────────────────┤
│  HTTPS/TLS  │  Load Balancer  │  Cloudflare CDN                │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Application Layer (Go + Fiber)               │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐             │
│  │  Handlers   │  │  Services   │  │ Repository  │             │
│  │  (HTTP)     │  │  (Business  │  │  (Data)     │             │
│  │             │  │   Logic)    │  │             │             │
│  └─────────────┘  └─────────────┘  └─────────────┘             │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │            Middleware Layer                             │   │
│  │  Auth │ CORS │ CSRF │ Rate Limit │ Validation │ Logging│   │
│  └─────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        ▼                     ▼                     ▼
┌───────────────┐   ┌───────────────┐   ┌───────────────┐
│  PostgreSQL   │   │     Redis     │   │  Ory Kratos   │
│  (Primary DB) │   │   (Cache)     │   │  (Auth)       │
└───────────────┘   └───────────────┘   └───────────────┘
        │
        └─────────────────────────────────────────────┐
                                                      ▼
┌─────────────────────────────────────────────────────────────────┐
│                    External Services                            │
├─────────────────────────────────────────────────────────────────┤
│  Midtrans  │  Komerce  │  Fonnte  │  Cloudflare R2  │  SMTP   │
│  (Payment) │ (Shipping)│ (WhatsApp)│  (Storage)    │ (Email) │
└─────────────────────────────────────────────────────────────────┘
```

### 5.2 Technology Stack

#### Backend
- **Language**: Go 1.24.0
- **Web Framework**: Fiber v2 (High-performance HTTP framework)
- **ORM**: GORM (Go Object Relational Mapping)
- **Database**: PostgreSQL 12+
- **Cache**: Redis 6+
- **Container**: Docker/Podman

#### Authentication
- **Identity Management**: Ory Kratos v1.1.0
- **Session Management**: HTTP-only cookies
- **Token-based Auth**: JWT for API access

#### External Integrations
- **Payment**: Midtrans (Snap API)
- **Shipping**: Komerce API
- **Notifications**: Fonnte (WhatsApp)
- **Storage**: Cloudflare R2 / Local
- **Email**: SMTP

#### Development Tools
- **API Documentation**: Swagger/OpenAPI
- **Testing**: Go testing framework + Testify
- **Database Migrations**: golang-migrate
- **Build Tool**: Make

### 5.3 Architecture Patterns

#### Layered Architecture
1. **Handler Layer**: HTTP request/response handling
2. **Service Layer**: Business logic implementation
3. **Repository Layer**: Data access operations
4. **Model Layer**: Data structures and entities

#### Design Patterns
- **Repository Pattern**: Data access abstraction
- **Service Layer Pattern**: Business logic encapsulation
- **Middleware Pattern**: Cross-cutting concerns
- **Factory Pattern**: Storage provider selection
- **Strategy Pattern**: Pricing calculation strategies

#### Caching Strategy
- **Cache-Aside Pattern**: Application-managed caching
- **TTL-based Expiration**: Automatic cache invalidation
- **Write-through**: Cache updates on data changes

### 5.4 Database Schema

#### Core Tables
- **products**: Product information
- **product_skus**: Product variants
- **product_media**: Product images/videos
- **categories**: Product categories
- **users**: User accounts
- **carts**: Shopping cart items
- **orders**: Order information
- **order_items**: Order line items
- **coupons**: Discount codes
- **flash_sales**: Promotional campaigns
- **reviews**: Customer reviews
- **wishlists**: Saved products
- **stock_logs**: Inventory audit trail
- **shipping_zones**: Shipping regions
- **tax_rates**: Tax configuration

---

## 6. User Stories

### 6.1 Product Management

**As a store manager**, I want to create products with multiple variants so that customers can choose from different sizes and colors.

**As a store manager**, I want to upload product images so that customers can see what they're buying.

**As a store manager**, I want to manage inventory levels so that I don't oversell products.

**As a store manager**, I want to view stock movement logs so that I can track inventory changes.

### 6.2 Pricing & Discounts

**As a store manager**, I want to set different prices for retail and reseller customers so that I can offer wholesale pricing.

**As a marketing manager**, I want to create flash sales so that I can run time-limited promotions.

**As a marketing manager**, I want to create coupon codes so that customers can get discounts.

**As a customer**, I want to see the final price including taxes so that I know exactly how much I'll pay.

### 6.3 Shopping & Checkout

**As a customer**, I want to add products to my cart so that I can purchase multiple items at once.

**As a customer**, I want to calculate shipping costs before checkout so that I know the total cost.

**As a customer**, I want to pay securely using multiple payment methods so that I can choose my preferred option.

**As a customer**, I want to receive order confirmation via WhatsApp so that I know my order is confirmed.

### 6.4 Order Management

**As a customer**, I want to track my order status so that I know when my order will arrive.

**As a store manager**, I want to update order status so that customers know their order progress.

**As a store manager**, I want to process refunds so that I can handle returns and cancellations.

**As a customer support agent**, I want to view order history so that I can help customers with their orders.

### 6.5 Authentication

**As a customer**, I want to register an account so that I can save my information for future purchases.

**As a customer**, I want to login securely so that my account is protected.

**As a store manager**, I want to manage user roles so that I can control access to admin features.

**As a developer**, I want to use API keys so that I can integrate external services.

### 6.6 Notifications

**As a customer**, I want to receive WhatsApp notifications for order updates so that I stay informed.

**As a store manager**, I want to send promotional messages so that I can engage customers.

**As a customer**, I want to receive payment confirmation so that I know my payment was successful.

**As a customer**, I want to receive shipping notifications so that I know when my order is shipped.

---

## 7. Functional Requirements

### 7.1 Product Management (FR-001 to FR-020)

#### FR-001: Create Product
- **Description**: System must allow creation of new products
- **Acceptance Criteria**:
  - Product can be created with name, description, price, weight, and category
  - Product slug is auto-generated from name
  - Product is assigned a unique ID
  - Product is saved to database
- **Priority**: High

#### FR-002: Update Product
- **Description**: System must allow updating existing products
- **Acceptance Criteria**:
  - All product fields can be updated
  - Product slug is regenerated if name changes
  - Update timestamp is recorded
  - Changes are persisted to database
- **Priority**: High

#### FR-003: Delete Product
- **Description**: System must allow deletion of products
- **Acceptance Criteria**:
  - Product is soft-deleted (deleted_at timestamp set)
  - Product is no longer visible in listings
  - Historical data is preserved
- **Priority**: Medium

#### FR-004: List Products
- **Description**: System must allow listing products with filtering and pagination
- **Acceptance Criteria**:
  - Products can be filtered by category, price range, availability
  - Results are paginated (default: 20 per page)
  - Products are sorted by creation date (newest first)
  - Response includes total count
- **Priority**: High

#### FR-005: Get Product Details
- **Description**: System must allow retrieving detailed product information
- **Acceptance Criteria**:
  - Product includes all fields
  - Product includes all variants
  - Product includes all media
  - Product includes category information
- **Priority**: High

#### FR-006: Create Product Variant (SKU)
- **Description**: System must allow creation of product variants
- **Acceptance Criteria**:
  - SKU can be created with size, color, stock, and price override
  - SKU is assigned a unique ID
  - SKU is linked to parent product
  - SKU stock is tracked separately
- **Priority**: High

#### FR-007: Update SKU Stock
- **Description**: System must allow manual stock updates for SKUs
- **Acceptance Criteria**:
  - Stock can be increased or decreased
  - Stock change is logged in stock_logs table
  - Stock cannot go below zero
  - Stock update timestamp is recorded
- **Priority**: High

#### FR-008: Upload Product Media
- **Description**: System must allow uploading product images and videos
- **Acceptance Criteria**:
  - Images can be uploaded (JPG, PNG, WEBP)
  - Videos can be uploaded (MP4)
  - File size limit: 10MB
  - Media is stored in Cloudflare R2 or local storage
  - Primary image can be designated
- **Priority**: High

#### FR-009: Delete Product Media
- **Description**: System must allow deletion of product media
- **Acceptance Criteria**:
  - Media is soft-deleted
  - Media file is removed from storage
  - Product media list is updated
- **Priority**: Medium

#### FR-010: Create Category
- **Description**: System must allow creation of product categories
- **Acceptance Criteria**:
  - Category can be created with name and description
  - Category slug is auto-generated
  - Category can have parent category (optional)
- **Priority**: High

#### FR-011: List Categories
- **Description**: System must allow listing categories
- **Acceptance Criteria**:
  - Categories are returned in hierarchical structure
  - Category includes product count
  - Categories are sorted alphabetically
- **Priority**: High

#### FR-012: Stock Log Tracking
- **Description**: System must track all stock movements
- **Acceptance Criteria**:
  - Every stock change is logged
  - Log includes: SKU ID, quantity change, reason, timestamp
  - Logs are queryable by SKU and date range
- **Priority**: High

#### FR-013: Flash Sale Management
- **Description**: System must support flash sale campaigns
- **Acceptance Criteria**:
  - Flash sale can be created with start/end time and discount
  - Flash sale can be applied to specific products
  - Flash sale price is automatically calculated
  - Flash sale is only active during specified time
- **Priority**: Medium

#### FR-014: Product Search
- **Description**: System must allow searching products
- **Acceptance Criteria**:
  - Search by product name
  - Search by SKU
  - Search is case-insensitive
  - Results are ranked by relevance
- **Priority**: High

#### FR-015: Product Reviews
- **Description**: System must support product reviews
- **Acceptance Criteria**:
  - Customers can submit reviews (1-5 stars, comment)
  - Reviews are linked to orders
  - Reviews can be moderated
  - Average rating is calculated
- **Priority**: Medium

#### FR-016: Wishlist Management
- **Description**: System must allow customers to save products to wishlist
- **Acceptance Criteria**:
  - Customers can add products to wishlist
  - Customers can remove products from wishlist
  - Wishlist is per customer
  - Wishlist persists across sessions
- **Priority**: Low

#### FR-017: Product Visibility Control
- **Description**: System must allow controlling product visibility
- **Acceptance Criteria**:
  - Products can be marked as active/inactive
  - Inactive products are not shown in listings
  - Inactive products can still be accessed by direct link
- **Priority**: Medium

#### FR-018: Bulk Product Operations
- **Description**: System must support bulk product operations
- **Acceptance Criteria**:
  - Multiple products can be updated at once
  - Multiple products can be deleted at once
  - Operations are transactional (all or nothing)
- **Priority**: Low

#### FR-019: Product Import/Export
- **Description**: System must support product import/export
- **Acceptance Criteria**:
  - Products can be exported to CSV
  - Products can be imported from CSV
  - Import validates data before insertion
- **Priority**: Low

#### FR-020: Product Analytics
- **Description**: System must track product analytics
- **Acceptance Criteria**:
  - Track product views
  - Track add-to-cart events
  - Track conversion rates
  - Analytics are queryable
- **Priority**: Low

### 7.2 Pricing Engine (FR-021 to FR-030)

#### FR-021: Retail Price Calculation
- **Description**: System must calculate retail prices
- **Acceptance Criteria**:
  - Returns base product price
  - Applies SKU price override if set
  - Returns final price
- **Priority**: High

#### FR-022: Reseller Price Calculation
- **Description**: System must calculate reseller prices
- **Acceptance Criteria**:
  - Applies reseller discount tier
  - Supports multiple reseller tiers
  - Returns discounted price
- **Priority**: High

#### FR-023: Flash Sale Price Calculation
- **Description**: System must calculate flash sale prices
- **Acceptance Criteria**:
  - Checks if flash sale is active
  - Applies flash sale discount
  - Returns discounted price
- **Priority**: High

#### FR-024: Coupon Validation
- **Description**: System must validate coupon codes
- **Acceptance Criteria**:
  - Checks if coupon exists
  - Checks if coupon is active
  - Checks if coupon has usage limit
  - Checks if coupon is expired
  - Returns validation result
- **Priority**: High

#### FR-025: Coupon Application
- **Description**: System must apply coupon discounts
- **Acceptance Criteria**:
  - Supports percentage discounts
  - Supports fixed amount discounts
  - Applies to cart total or specific products
  - Calculates final discount amount
- **Priority**: High

#### FR-026: Tax Calculation
- **Description**: System must calculate taxes
- **Acceptance Criteria**:
  - Applies 11% VAT by default
  - Supports configurable tax rates
  - Calculates tax on subtotal
  - Returns tax amount
- **Priority**: High

#### FR-027: Order Summary Calculation
- **Description**: System must calculate complete order summary
- **Acceptance Criteria**:
  - Calculates subtotal (items × quantity)
  - Applies discounts (coupons, flash sales)
  - Calculates tax
  - Calculates shipping cost
  - Returns total amount
- **Priority**: High

#### FR-028: Free Shipping Calculation
- **Description**: System must check for free shipping eligibility
- **Acceptance Criteria**:
  - Checks if order meets minimum threshold
  - Checks if free shipping coupon is applied
  - Returns shipping cost (0 if eligible)
- **Priority**: Medium

#### FR-029: Price History Tracking
- **Description**: System must track price changes
- **Acceptance Criteria**:
  - Records price changes
  - Includes timestamp and reason
  - History is queryable
- **Priority**: Low

#### FR-030: Dynamic Pricing Rules
- **Description**: System must support dynamic pricing rules
- **Acceptance Criteria**:
  - Rules can be created based on conditions
  - Rules can be applied automatically
  - Rules can be enabled/disabled
- **Priority**: Low

### 7.3 Order Management (FR-031 to FR-045)

#### FR-031: Add to Cart
- **Description**: System must allow adding products to cart
- **Acceptance Criteria**:
  - Product can be added to cart
  - SKU can be specified
  - Quantity can be specified
  - Stock availability is checked
  - Cart is per customer
- **Priority**: High

#### FR-032: Update Cart Item
- **Description**: System must allow updating cart item quantities
- **Acceptance Criteria**:
  - Quantity can be increased
  - Quantity can be decreased
  - Stock availability is checked
  - Item is removed if quantity is 0
- **Priority**: High

#### FR-033: Remove from Cart
- **Description**: System must allow removing items from cart
- **Acceptance Criteria**:
  - Item is removed from cart
  - Cart is updated
  - Stock is not affected
- **Priority**: High

#### FR-034: View Cart
- **Description**: System must allow viewing cart contents
- **Acceptance Criteria**:
  - Returns all cart items
  - Includes product details
  - Includes SKU details
  - Calculates subtotal
- **Priority**: High

#### FR-035: Checkout Initiation
- **Description**: System must initiate checkout process
- **Acceptance Criteria**:
  - Validates cart contents
  - Validates stock availability
  - Calculates order summary
  - Creates pending order
  - Generates payment token
- **Priority**: High

#### FR-036: Payment Processing
- **Description**: System must process payments via Midtrans
- **Acceptance Criteria**:
  - Generates Snap token
  - Redirects to Midtrans payment page
  - Handles payment success/failure
  - Updates order status
- **Priority**: High

#### FR-037: Payment Webhook Handling
- **Description**: System must handle Midtrans webhooks
- **Acceptance Criteria**:
  - Validates webhook signature
  - Updates payment status
  - Updates order status
  - Adjusts stock on payment success
  - Restores stock on payment failure
- **Priority**: High

#### FR-038: Order Status Updates
- **Description**: System must allow updating order status
- **Acceptance Criteria**:
  - Status can be: pending, confirmed, shipped, delivered, cancelled, refunded
  - Status changes are logged
  - Notifications are sent on status changes
- **Priority**: High

#### FR-039: Order History
- **Description**: System must provide order history for customers
- **Acceptance Criteria**:
  - Returns all customer orders
  - Orders are sorted by date (newest first)
  - Includes order details and status
- **Priority**: High

#### FR-040: Order Details
- **Description**: System must provide detailed order information
- **Acceptance Criteria**:
  - Returns order information
  - Returns all order items
  - Returns payment status
  - Returns shipping information
- **Priority**: High

#### FR-041: Order Cancellation
- **Description**: System must allow order cancellation
- **Acceptance Criteria**:
  - Orders can be cancelled if not yet shipped
  - Stock is restored on cancellation
  - Refund is initiated if payment was made
  - Status is updated to cancelled
- **Priority**: High

#### FR-042: Order Refund
- **Description**: System must process refunds
- **Acceptance Criteria**:
  - Refund can be initiated for paid orders
  - Refund is processed via Midtrans
  - Stock is restored on refund
  - Status is updated to refunded
- **Priority**: High

#### FR-043: Order Search
- **Description**: System must allow searching orders
- **Acceptance Criteria**:
  - Search by order ID
  - Search by customer email
  - Search by date range
  - Search by status
- **Priority**: Medium

#### FR-044: Order Analytics
- **Description**: System must provide order analytics
- **Acceptance Criteria**:
  - Tracks total orders
  - Tracks total revenue
  - Tracks average order value
  - Tracks conversion rate
- **Priority**: Low

#### FR-045: Bulk Order Operations
- **Description**: System must support bulk order operations
- **Acceptance Criteria**:
  - Multiple orders can be updated at once
  - Multiple orders can be exported
  - Operations are transactional
- **Priority**: Low

### 7.4 Shipping Integration (FR-046 to FR-055)

#### FR-046: Destination Search
- **Description**: System must allow searching shipping destinations
- **Acceptance Criteria**:
  - Search by city name
  - Search by district name
  - Returns destination details
  - Results are paginated
- **Priority**: High

#### FR-047: Shipping Cost Calculation
- **Description**: System must calculate shipping costs
- **Acceptance Criteria**:
  - Calculates based on origin and destination
  - Calculates based on total weight
  - Returns multiple courier options
  - Returns estimated delivery time
- **Priority**: High

#### FR-048: Shipping Zone Management
- **Description**: System must support shipping zones
- **Acceptance Criteria**:
  - Zones can be created
  - Zones can have custom rates
  - Zones can have free shipping thresholds
- **Priority**: Medium

#### FR-049: Courier Selection
- **Description**: System must allow courier selection
- **Acceptance Criteria**:
  - Multiple couriers are available
  - Customer can select preferred courier
  - Selection is saved to order
- **Priority**: High

#### FR-050: Shipping Tracking
- **Description**: System must track shipping status
- **Acceptance Criteria**:
  - Tracking number is stored
  - Tracking URL is provided
  - Status is updated via webhook
- **Priority**: Medium

#### FR-051: Weight Calculation
- **Description**: System must calculate total order weight
- **Acceptance Criteria**:
  - Calculates based on SKU weight × quantity
  - Returns total weight in grams
- **Priority**: High

#### FR-052: Shipping Address Validation
- **Description**: System must validate shipping addresses
- **Acceptance Criteria**:
  - Validates required fields
  - Validates postal code
  - Validates phone number
- **Priority**: High

#### FR-053: Multiple Shipping Addresses
- **Description**: System must support multiple shipping addresses
- **Acceptance Criteria**:
  - Customers can save multiple addresses
  - Default address can be set
  - Address can be selected at checkout
- **Priority**: Medium

#### FR-054: Shipping Insurance
- **Description**: System must support shipping insurance
- **Acceptance Criteria**:
  - Insurance can be added to order
  - Insurance cost is calculated
  - Insurance is optional
- **Priority**: Low

#### FR-055: International Shipping
- **Description**: System must support international shipping (Future)
- **Acceptance Criteria**:
  - International destinations are supported
  - Customs documentation is generated
  - International rates are calculated
- **Priority**: Low

### 7.5 Authentication (FR-056 to FR-065)

#### FR-056: User Registration
- **Description**: System must allow user registration via Kratos
- **Acceptance Criteria**:
  - Users can register with email and password
  - Email verification is required
  - User profile is created
  - Session is established
- **Priority**: High

#### FR-057: User Login
- **Description**: System must allow user login via Kratos
- **Acceptance Criteria**:
  - Users can login with email and password
  - Session is established
  - Session cookie is set
  - Login is secure
- **Priority**: High

#### FR-058: User Logout
- **Description**: System must allow user logout
- **Acceptance Criteria**:
  - Session is terminated
  - Session cookie is cleared
  - User is redirected
- **Priority**: High

#### FR-059: Session Management
- **Description**: System must manage user sessions
- **Acceptance Criteria**:
  - Sessions are validated on each request
  - Sessions expire after inactivity
  - Sessions can be revoked
- **Priority**: High

#### FR-060: Password Reset
- **Description**: System must support password reset
- **Acceptance Criteria**:
  - Users can request password reset
  - Reset link is sent via email
  - Password can be updated via link
- **Priority**: High

#### FR-061: Profile Management
- **Description**: System must allow profile management
- **Acceptance Criteria**:
  - Users can update profile information
  - Users can change password
  - Changes are persisted
- **Priority**: Medium

#### FR-062: Role-Based Access Control
- **Description**: System must implement RBAC (Planned)
- **Acceptance Criteria**:
  - Roles can be assigned to users
  - Permissions can be assigned to roles
  - Access is checked per endpoint
  - Unauthorized access is denied
- **Priority**: High

#### FR-063: API Key Authentication
- **Description**: System must support API key authentication
- **Acceptance Criteria**:
  - API keys can be generated
  - API keys can be revoked
  - API keys are validated
  - API keys are scoped to permissions
- **Priority**: Medium

#### FR-064: Social Login
- **Description**: System must support social login (Future)
- **Acceptance Criteria**:
  - Google login is supported
  - Facebook login is supported
  - User profile is created/linked
- **Priority**: Low

#### FR-065: Two-Factor Authentication
- **Description**: System must support 2FA (Future)
- **Acceptance Criteria**:
  - 2FA can be enabled
  - OTP is sent via SMS/Email
  - OTP is verified on login
- **Priority**: Low

### 7.6 Notifications (FR-066 to FR-075)

#### FR-066: Order Confirmation Notification
- **Description**: System must send order confirmation via WhatsApp
- **Acceptance Criteria**:
  - Notification is sent after order creation
  - Includes order details
  - Includes payment information
  - Sent via Fonnte
- **Priority**: High

#### FR-067: Payment Success Notification
- **Description**: System must send payment success notification
- **Acceptance Criteria**:
  - Notification is sent after payment confirmation
  - Includes payment details
  - Sent via WhatsApp
- **Priority**: High

#### FR-068: Shipping Notification
- **Description**: System must send shipping notification
- **Acceptance Criteria**:
  - Notification is sent when order is shipped
  - Includes tracking number
  - Includes estimated delivery
  - Sent via WhatsApp
- **Priority**: High

#### FR-069: Delivery Notification
- **Description**: System must send delivery notification
- **Acceptance Criteria**:
  - Notification is sent when order is delivered
  - Includes delivery confirmation
  - Sent via WhatsApp
- **Priority**: Medium

#### FR-070: Cancellation Notification
- **Description**: System must send cancellation notification
- **Acceptance Criteria**:
  - Notification is sent when order is cancelled
  - Includes cancellation reason
  - Sent via WhatsApp
- **Priority**: Medium

#### FR-071: Refund Notification
- **Description**: System must send refund notification
- **Acceptance Criteria**:
  - Notification is sent when refund is processed
  - Includes refund amount
  - Sent via WhatsApp
- **Priority**: Medium

#### FR-072: Promotional Notifications
- **Description**: System must send promotional messages
- **Acceptance Criteria**:
  - Promotional messages can be sent
  - Messages can be targeted
  - Opt-out option is available
- **Priority**: Low

#### FR-073: Email Notifications
- **Description**: System must support email notifications
- **Acceptance Criteria**:
  - Emails are sent via SMTP
  - Email templates are customizable
  - Notifications are queued
- **Priority**: Medium

#### FR-074: Notification Preferences
- **Description**: System must allow notification preferences
- **Acceptance Criteria**:
  - Users can opt-in/opt-out
  - Preferences are saved
  - Preferences are respected
- **Priority**: Low

#### FR-075: Notification History
- **Description**: System must track notification history
- **Acceptance Criteria**:
  - All notifications are logged
  - Delivery status is tracked
  - History is queryable
- **Priority**: Low

### 7.7 Media Management (FR-076 to FR-080)

#### FR-076: File Upload
- **Description**: System must allow file uploads
- **Acceptance Criteria**:
  - Images can be uploaded (JPG, PNG, WEBP)
  - Videos can be uploaded (MP4)
  - File size limit: 10MB
  - Upload is validated
- **Priority**: High

#### FR-077: Cloud Storage Integration
- **Description**: System must integrate with Cloudflare R2
- **Acceptance Criteria**:
  - Files are stored in R2
  - Public URLs are generated
  - Files are accessible via CDN
- **Priority**: High

#### FR-078: Local Storage Support
- **Description**: System must support local storage for development
- **Acceptance Criteria**:
  - Files can be stored locally
  - Files are accessible via HTTP
  - Storage is configurable
- **Priority**: Medium

#### FR-079: Image Optimization
- **Description**: System must optimize images (Future)
- **Acceptance Criteria**:
  - Images are resized
  - Images are compressed
  - Multiple formats are generated
- **Priority**: Low

#### FR-080: Media Gallery
- **Description**: System must provide media gallery
- **Acceptance Criteria**:
  - All media can be viewed
  - Media can be filtered
  - Media can be deleted
- **Priority**: Low

---

## 8. Non-Functional Requirements

### 8.1 Performance Requirements (NFR-001 to NFR-010)

#### NFR-001: API Response Time
- **Description**: API endpoints must respond within specified time limits
- **Requirements**:
  - Simple GET requests: < 100ms (p95)
  - Complex queries: < 500ms (p95)
  - Write operations: < 200ms (p95)
  - Checkout process: < 1s (p95)
- **Measurement**: Response time monitoring with APM
- **Priority**: Critical

#### NFR-002: Throughput
- **Description**: System must handle specified request volume
- **Requirements**:
  - 1,000 requests/second sustained
  - 5,000 requests/second peak
  - 10,000 concurrent users
- **Measurement**: Load testing with k6/JMeter
- **Priority**: Critical

#### NFR-003: Database Performance
- **Description**: Database queries must be optimized
- **Requirements**:
  - Query execution time: < 50ms (p95)
  - Connection pool: 100 connections
  - Index optimization for all queries
- **Measurement**: Database query logging and analysis
- **Priority**: High

#### NFR-004: Cache Performance
- **Description**: Caching must improve response times
- **Requirements**:
  - Cache hit rate: > 80%
  - Cache response time: < 10ms
  - Cache TTL: Configurable per endpoint
- **Measurement**: Redis monitoring
- **Priority**: High

#### NFR-005: Image Loading
- **Description**: Product images must load quickly
- **Requirements**:
  - Image load time: < 2s
  - Image optimization: WebP format
  - CDN delivery: < 100ms
- **Measurement**: Page speed tools
- **Priority**: Medium

### 8.2 Availability Requirements (NFR-011 to NFR-015)

#### NFR-011: Uptime
- **Description**: System must maintain high availability
- **Requirements**:
  - 99.9% uptime (8.76 hours downtime/year)
  - Scheduled maintenance: < 4 hours/month
  - Graceful degradation on failures
- **Measurement**: Uptime monitoring (UptimeRobot, Pingdom)
- **Priority**: Critical

#### NFR-012: Disaster Recovery
- **Description**: System must recover from failures
- **Requirements**:
  - RTO (Recovery Time Objective): < 1 hour
  - RPO (Recovery Point Objective): < 15 minutes
  - Automated backups: Daily
  - Backup retention: 30 days
- **Measurement**: Disaster recovery testing
- **Priority**: High

#### NFR-013: Database Replication
- **Description**: Database must have replication for high availability
- **Requirements**:
  - Primary-replica setup
  - Automatic failover
  - Read replicas for scaling
- **Measurement**: Replication lag monitoring
- **Priority**: High

#### NFR-014: Load Balancing
- **Description**: System must distribute load across instances
- **Requirements**:
  - Horizontal scaling support
  - Session persistence (sticky sessions)
  - Health checks for instances
- **Measurement**: Load balancer metrics
- **Priority**: High

#### NFR-015: Graceful Shutdown
- **Description**: System must shutdown gracefully
- **Requirements**:
  - Complete in-flight requests
  - Close database connections
  - Flush cache
  - Maximum shutdown time: 30s
- **Measurement**: Shutdown testing
- **Priority**: Medium

### 8.3 Scalability Requirements (NFR-016 to NFR-020)

#### NFR-016: Horizontal Scaling
- **Description**: System must support horizontal scaling
- **Requirements**:
  - Stateless application design
  - Shared storage (R2, Redis)
  - Load balancer support
  - Auto-scaling capability
- **Measurement**: Scaling tests
- **Priority**: High

#### NFR-017: Database Scaling
- **Description**: Database must support scaling
- **Requirements**:
  - Connection pooling
  - Read replicas
  - Database sharding (future)
  - Query optimization
- **Measurement**: Database performance monitoring
- **Priority**: High

#### NFR-018: Storage Scaling
- **Description**: Storage must scale automatically
- **Requirements**:
  - Cloudflare R2 (unlimited storage)
  - CDN delivery
  - Automatic optimization
- **Measurement**: Storage usage monitoring
- **Priority**: Medium

#### NFR-019: Cache Scaling
- **Description**: Cache must support scaling
- **Requirements**:
  - Redis Cluster support
  - Automatic failover
  - Memory optimization
- **Measurement**: Cache performance monitoring
- **Priority**: Medium

#### NFR-020: Microservices Readiness
- **Description**: Architecture must support microservices (future)
- **Requirements**:
  - Service boundaries defined
  - API Gateway support
  - Service discovery
  - Inter-service communication
- **Measurement**: Architecture review
- **Priority**: Low

### 8.4 Security Requirements (NFR-021 to NFR-030)

#### NFR-021: Authentication Security
- **Description**: Authentication must be secure
- **Requirements**:
  - Ory Kratos integration
  - HTTP-only session cookies
  - Secure flag on cookies (HTTPS only)
  - SameSite attribute: Lax
- **Measurement**: Security audit
- **Priority**: Critical

#### NFR-022: Authorization Security
- **Description**: Authorization must be properly implemented
- **Requirements**:
  - RBAC implementation
  - Permission checks on all endpoints
  - Principle of least privilege
  - Regular permission audits
- **Measurement**: Access control testing
- **Priority**: Critical

#### NFR-023: Data Encryption
- **Description**: Sensitive data must be encrypted
- **Requirements**:
  - TLS 1.3 for data in transit
  - AES-256 for data at rest
  - Encrypted backups
  - Environment variable encryption
- **Measurement**: Encryption audit
- **Priority**: Critical

#### NFR-024: Input Validation
- **Description**: All inputs must be validated
- **Requirements**:
  - Server-side validation
  - SQL injection prevention
  - XSS prevention
  - CSRF protection
- **Measurement**: Security testing (OWASP ZAP)
- **Priority**: Critical

#### NFR-025: API Security
- **Description**: API endpoints must be secure
- **Requirements**:
  - Rate limiting: 100 req/min per IP
  - API key authentication for external access
  - Request size limits: 10MB
  - Request timeout: 30s
- **Measurement**: API security testing
- **Priority**: High

#### NFR-026: OWASP Compliance
- **Description**: System must comply with OWASP standards
- **Requirements**:
  - OWASP Top 10 mitigation
  - Security headers implementation
  - Regular security scans
  - Vulnerability management
- **Measurement**: OWASP testing
- **Priority**: High

#### NFR-027: Payment Security
- **Description**: Payment processing must be secure
- **Requirements**:
  - PCI DSS compliance
  - Midtrans integration (PCI compliant)
  - No card data storage
  - Secure webhook verification
- **Measurement**: PCI DSS audit
- **Priority**: Critical

#### NFR-028: Logging and Monitoring
- **Description**: Security events must be logged
- **Requirements**:
  - Audit logging for all actions
  - Failed login attempts
  - Suspicious activity detection
  - Log retention: 90 days
- **Measurement**: Log analysis
- **Priority**: High

#### NFR-029: Secrets Management
- **Description**: Secrets must be securely managed
- **Requirements**:
  - Environment variables for secrets
  - No secrets in code
  - Regular secret rotation
  - Secret vault (future: HashiCorp Vault)
- **Measurement**: Secrets audit
- **Priority**: High

#### NFR-030: Dependency Security
- **Description**: Dependencies must be secure
- **Requirements**:
  - Regular dependency updates
  - Vulnerability scanning
  - SCA (Software Composition Analysis)
  - SBOM (Software Bill of Materials)
- **Measurement**: Dependency scanning (Snyk, Dependabot)
- **Priority**: Medium

### 8.5 Maintainability Requirements (NFR-031 to NFR-035)

#### NFR-031: Code Quality
- **Description**: Code must follow best practices
- **Requirements**:
  - Go coding standards
  - Code review process
  - Linting (golangci-lint)
  - Code coverage: > 80%
- **Measurement**: Code quality tools
- **Priority**: High

#### NFR-032: Documentation
- **Description**: System must be well-documented
- **Requirements**:
  - API documentation (Swagger)
  - Code comments
  - Architecture documentation
  - Deployment documentation
- **Measurement**: Documentation review
- **Priority**: High

#### NFR-033: Testing
- **Description**: System must have comprehensive tests
- **Requirements**:
  - Unit tests for all services
  - Integration tests for critical flows
  - E2E tests for user journeys
  - Test coverage: > 80%
- **Measurement**: Test coverage reports
- **Priority**: High

#### NFR-034: CI/CD Pipeline
- **Description**: System must have automated CI/CD
- **Requirements**:
  - Automated testing on PR
  - Automated deployment on merge
  - Rollback capability
  - Deployment tracking
- **Measurement**: CI/CD metrics
- **Priority**: High

#### NFR-035: Monitoring and Alerting
- **Description**: System must have comprehensive monitoring
- **Requirements**:
  - Application performance monitoring
  - Error tracking
  - Log aggregation
  - Alerting on critical issues
- **Measurement**: Monitoring setup
- **Priority**: High

### 8.6 Usability Requirements (NFR-036 to NFR-040)

#### NFR-036: API Design
- **Description**: API must be intuitive and consistent
- **Requirements**:
  - RESTful design principles
  - Consistent naming conventions
  - Clear error messages
  - Versioning support
- **Measurement**: API review
- **Priority**: High

#### NFR-037: Error Handling
- **Description**: Errors must be handled gracefully
- **Requirements**:
  - Consistent error response format
  - Appropriate HTTP status codes
  - Detailed error messages (dev)
  - Generic error messages (prod)
- **Measurement**: Error handling testing
- **Priority**: High

#### NFR-038: Developer Experience
- **Description**: System must be developer-friendly
- **Requirements**:
  - Comprehensive API documentation
  - SDK examples
  - Postman collections
  - Quick start guide
- **Measurement**: Developer feedback
- **Priority**: Medium

#### NFR-039: Internationalization
- **Description**: System must support multiple languages (Future)
- **Requirements**:
  - i18n support
  - Multi-language content
  - Currency localization
  - Date/time localization
- **Measurement**: i18n testing
- **Priority**: Low

#### NFR-040: Accessibility
- **Description**: System must be accessible (Future)
- **Requirements**:
  - WCAG 2.1 compliance
  - Screen reader support
  - Keyboard navigation
  - Color contrast compliance
- **Measurement**: Accessibility testing
- **Priority**: Low

---

## 9. Integration Requirements

### 9.1 Payment Gateway Integration (IR-001 to IR-005)

#### IR-001: Midtrans Integration
- **Description**: Integrate Midtrans payment gateway
- **Requirements**:
  - Snap token generation
  - Payment status tracking
  - Webhook handling
  - Refund processing
  - Signature verification
- **API Version**: v1
- **Environment**: Sandbox → Production
- **Priority**: Critical

#### IR-002: Payment Methods
- **Description**: Support multiple payment methods
- **Requirements**:
  - Credit/Debit Cards
  - Bank Transfer (VA)
  - E-Wallets (GoPay, OVO, Dana)
  - QRIS
  - Retail outlets (Alfamart, Indomaret)
- **Priority**: High

#### IR-003: Payment Webhooks
- **Description**: Handle Midtrans webhooks
- **Requirements**:
  - Payment success webhook
  - Payment failure webhook
  - Payment pending webhook
  - Refund webhook
  - Signature verification
- **Priority**: Critical

#### IR-004: Payment Security
- **Description**: Ensure payment security
- **Requirements**:
  - PCI DSS compliance
  - No card data storage
  - Secure webhook endpoints
  - HTTPS only
- **Priority**: Critical

#### IR-005: Payment Analytics
- **Description**: Track payment analytics
- **Requirements**:
  - Payment success rate
  - Payment method distribution
  - Average payment time
  - Failed payment reasons
- **Priority**: Medium

### 9.2 Shipping Integration (IR-006 to IR-010)

#### IR-006: Komerce Integration
- **Description**: Integrate Komerce shipping API
- **Requirements**:
  - Destination search
  - Shipping cost calculation
  - Multiple couriers support
  - Tracking integration
- **API Version**: v1
- **Environment**: Sandbox → Production
- **Priority**: Critical

#### IR-007: Courier Support
- **Description**: Support multiple couriers
- **Requirements**:
  - JNE
  - J&T
  - SiCepat
  - AnterAja
  - Wahana
- **Priority**: High

#### IR-008: Shipping Zones
- **Description**: Configure shipping zones
- **Requirements**:
  - Zone-based pricing
  - Free shipping thresholds
  - Zone-specific couriers
  - Weight-based pricing
- **Priority**: Medium

#### IR-009: Tracking Integration
- **Description**: Integrate shipment tracking
- **Requirements**:
  - Tracking number storage
  - Tracking status updates
  - Tracking URL generation
  - Customer notifications
- **Priority**: Medium

#### IR-010: Shipping Analytics
- **Description**: Track shipping analytics
- **Requirements**:
  - Shipping cost distribution
  - Courier usage statistics
  - Delivery time analysis
  - Shipping error tracking
- **Priority**: Low

### 9.3 Notification Integration (IR-011 to IR-015)

#### IR-011: Fonnte Integration
- **Description**: Integrate Fonnte WhatsApp API
- **Requirements**:
  - Send text messages
  - Send media messages
  - Message delivery tracking
  - Queue management
- **API Version**: v1
- **Priority**: High

#### IR-012: Notification Templates
- **Description**: Create notification templates
- **Requirements**:
  - Order confirmation template
  - Payment success template
  - Shipping notification template
  - Delivery notification template
- **Priority**: High

#### IR-013: Email Integration
- **Description**: Integrate SMTP for email notifications
- **Requirements**:
  - SMTP configuration
  - Email templates
  - Email queue
  - Delivery tracking
- **Priority**: Medium

#### IR-014: Notification Preferences
- **Description**: Allow notification preferences
- **Requirements**:
  - Opt-in/opt-out
  - Channel selection (WhatsApp, Email)
  - Frequency control
  - Preference storage
- **Priority**: Low

#### IR-015: Notification Analytics
- **Description**: Track notification analytics
- **Requirements**:
  - Delivery rate
  - Open rate (email)
  - Click-through rate
  - Failed notifications
- **Priority**: Low

### 9.4 Storage Integration (IR-016 to IR-020)

#### IR-016: Cloudflare R2 Integration
- **Description**: Integrate Cloudflare R2 storage
- **Requirements**:
  - File upload
  - File retrieval
  - File deletion
  - Public URL generation
- **API Version**: S3-compatible
- **Priority**: High

#### IR-017: CDN Integration
- **Description**: Use Cloudflare CDN for delivery
- **Requirements**:
  - CDN configuration
  - Cache control headers
  - Image optimization
  - Global edge delivery
- **Priority**: High

#### IR-018: Local Storage Support
- **Description**: Support local storage for development
- **Requirements**:
  - File system storage
  - HTTP serving
  - Configuration switch
  - Development convenience
- **Priority**: Medium

#### IR-019: Image Optimization
- **Description**: Optimize images automatically
- **Requirements**:
  - Resize images
  - Compress images
  - Generate thumbnails
  - Convert formats (WebP)
- **Priority**: Low

#### IR-020: Storage Analytics
- **Description**: Track storage usage
- **Requirements**:
  - Storage usage metrics
  - Bandwidth usage
  - Popular files
  - Storage costs
- **Priority**: Low

### 9.5 Authentication Integration (IR-021 to IR-025)

#### IR-021: Ory Kratos Integration
- **Description**: Integrate Ory Kratos for authentication
- **Requirements**:
  - User registration
  - User login
  - Session management
  - Password reset
- **Version**: v1.1.0
- **Priority**: Critical

#### IR-022: Kratos Configuration
- **Description**: Configure Kratos settings
- **Requirements**:
  - Identity schema
  - Self-service flows
  - Session configuration
  - Security settings
- **Priority**: Critical

#### IR-023: Social Login
- **Description**: Integrate social login providers (Future)
- **Requirements**:
  - Google OAuth
  - Facebook OAuth
  - Apple Sign-in
  - Profile mapping
- **Priority**: Low

#### IR-024: 2FA Integration
- **Description**: Integrate two-factor authentication (Future)
- **Requirements**:
  - TOTP support
  - SMS verification
  - Email verification
  - Backup codes
- **Priority**: Low

#### IR-025: User Sync
- **Description**: Sync Kratos users to local database
- **Requirements**:
  - User profile sync
  - Role sync
  - Real-time updates
  - Conflict resolution
- **Priority**: High

---

## 10. Security Requirements

### 10.1 Authentication Security (SR-001 to SR-010)

#### SR-001: Password Policy
- **Description**: Enforce strong password policy
- **Requirements**:
  - Minimum 8 characters
  - At least 1 uppercase letter
  - At least 1 lowercase letter
  - At least 1 number
  - At least 1 special character
- **Priority**: Critical

#### SR-002: Password Hashing
- **Description**: Hash passwords securely
- **Requirements**:
  - Use Argon2 or bcrypt
  - Minimum cost factor: 10
  - Salt per password
  - No reversible encryption
- **Priority**: Critical

#### SR-003: Session Management
- **Description**: Manage sessions securely
- **Requirements**:
  - HTTP-only session cookies
  - Secure flag (HTTPS only)
  - SameSite attribute: Lax
  - Session expiration: 24 hours
- **Priority**: Critical

#### SR-004: Session Fixation Prevention
- **Description**: Prevent session fixation attacks
- **Requirements**:
  - Regenerate session ID on login
  - Invalidate old sessions
  - Bind session to IP (optional)
- **Priority**: High

#### SR-005: Multi-Factor Authentication
- **Description**: Support MFA (Future)
- **Requirements**:
  - TOTP-based MFA
  - SMS-based MFA
  - Email-based MFA
  - Backup codes
- **Priority**: Medium

#### SR-006: Login Rate Limiting
- **Description**: Limit login attempts
- **Requirements**:
  - 5 failed attempts: 5-minute lockout
  - 10 failed attempts: 30-minute lockout
  - IP-based tracking
  - Account-based tracking
- **Priority**: High

#### SR-007: Password Reset Security
- **Description**: Secure password reset process
- **Requirements**:
  - Time-limited reset links (1 hour)
  - Single-use reset links
  - Email notification on reset
  - Invalidate existing sessions
- **Priority**: High

#### SR-008: Account Lockout
- **Description**: Lock out suspicious accounts
- **Requirements**:
  - Automatic lockout after failed attempts
  - Manual unlock by admin
  - Email notification on lockout
  - Audit logging
- **Priority**: High

#### SR-009: Social Login Security
- **Description**: Secure social login (Future)
- **Requirements**:
  - OAuth 2.0 / OpenID Connect
  - State parameter validation
  - Token validation
  - Profile mapping
- **Priority**: Low

#### SR-010: Session Timeout
- **Description**: Implement session timeout
- **Requirements**:
  - Inactivity timeout: 30 minutes
  - Absolute timeout: 24 hours
  - Warning before timeout
  - Auto-logout on timeout
- **Priority**: Medium

### 10.2 Authorization Security (SR-011 to SR-020)

#### SR-011: Role-Based Access Control
- **Description**: Implement RBAC
- **Requirements**:
  - Define roles (admin, manager, customer)
  - Define permissions per role
  - Enforce permissions on all endpoints
  - Audit permission changes
- **Priority**: Critical

#### SR-012: Principle of Least Privilege
- **Description**: Apply least privilege principle
- **Requirements**:
  - Default deny access
  - Grant minimum necessary permissions
  - Regular permission audits
  - Revoke unused permissions
- **Priority**: High

#### SR-013: API Authorization
- **Description**: Authorize API requests
- **Requirements**:
  - Validate session on each request
  - Check permissions per endpoint
  - Log authorization failures
  - Return 403 for unauthorized access
- **Priority**: Critical

#### SR-014: Resource Ownership
- **Description**: Enforce resource ownership
- **Requirements**:
  - Customers can only access their data
  - Admins can access all data
  - Managers can access assigned data
  - Audit access attempts
- **Priority**: High

#### SR-015: Admin Access Control
- **Description**: Secure admin endpoints
- **Requirements**:
  - Admin role required
  - Additional authentication (optional)
  - Audit all admin actions
  - IP whitelisting (optional)
- **Priority**: High

#### SR-016: API Key Security
- **Description**: Secure API key usage
- **Requirements**:
  - Generate strong API keys
  - Scope API keys to permissions
  - Rotate API keys regularly
  - Revoke compromised keys
- **Priority**: High

#### SR-017: Permission Caching
- **Description**: Cache permissions efficiently
- **Requirements**:
  - Cache user permissions
  - Invalidate on role changes
  - Cache timeout: 5 minutes
  - Fallback to database
- **Priority**: Medium

#### SR-018: Audit Logging
- **Description**: Log all authorization events
- **Requirements**:
  - Log successful access
  - Log failed access attempts
  - Log permission changes
  - Retain logs: 90 days
- **Priority**: High

#### SR-019: Cross-Tenant Isolation
- **Description**: Isolate multi-tenant data (Future)
- **Requirements**:
  - Tenant-specific data
  - Tenant isolation at DB level
  - Tenant-specific permissions
  - Audit cross-tenant access
- **Priority**: Low

#### SR-020: Dynamic Permissions
- **Description**: Support dynamic permissions (Future)
- **Requirements**:
  - Attribute-based access control (ABAC)
  - Policy-based permissions
  - Real-time permission evaluation
  - Policy versioning
- **Priority**: Low

### 10.3 Data Security (SR-021 to SR-030)

#### SR-021: Data Encryption in Transit
- **Description**: Encrypt all data in transit
- **Requirements**:
  - TLS 1.3 for all connections
  - Strong cipher suites
  - HSTS enabled
  - Certificate rotation
- **Priority**: Critical

#### SR-022: Data Encryption at Rest
- **Description**: Encrypt sensitive data at rest
- **Requirements**:
  - AES-256 encryption
  - Database encryption
  - Backup encryption
  - Key management
- **Priority**: Critical

#### SR-023: PII Protection
- **Description**: Protect personally identifiable information
- **Requirements**:
  - Identify all PII
  - Encrypt PII fields
  - Limit PII access
  - PII retention policy
- **Priority**: Critical

#### SR-024: Payment Data Protection
- **Description**: Protect payment data
- **Requirements**:
  - PCI DSS compliance
  - No card data storage
  - Tokenization via Midtrans
  - Secure payment flow
- **Priority**: Critical

#### SR-025: Data Masking
- **Description**: Mask sensitive data in logs
- **Requirements**:
  - Mask email addresses
  - Mask phone numbers
  - Mask credit card numbers
  - No sensitive data in logs
- **Priority**: High

#### SR-026: Data Minimization
- **Description**: Collect only necessary data
- **Requirements**:
  - Define data collection requirements
  - Avoid collecting unnecessary data
  - Delete unused data
  - Regular data audits
- **Priority**: Medium

#### SR-027: Data Retention
- **Description**: Define data retention policies
- **Requirements**:
  - Order data: 7 years
  - Customer data: 5 years after last activity
  - Logs: 90 days
  - Backups: 30 days
- **Priority**: Medium

#### SR-028: Data Backup Security
- **Description**: Secure data backups
- **Requirements**:
  - Encrypted backups
  - Off-site backup storage
  - Regular backup testing
  - Backup access controls
- **Priority**: High

#### SR-029: Data Erasure
- **Description**: Support data erasure requests
- **Requirements**:
  - Right to be forgotten
  - Complete data deletion
  - Backup deletion
  - Confirmation of deletion
- **Priority**: Medium

#### SR-030: Data Integrity
- **Description**: Ensure data integrity
- **Requirements**:
  - Database constraints
  - Transaction integrity
  - Checksum verification
  - Audit trail
- **Priority**: High

### 10.4 Network Security (SR-031 to SR-040)

#### SR-031: Firewall Configuration
- **Description**: Configure firewall rules
- **Requirements**:
  - Allow only necessary ports
  - Block all inbound traffic by default
  - Rate limit connections
  - Geo-blocking (optional)
- **Priority**: High

#### SR-032: DDoS Protection
- **Description**: Protect against DDoS attacks
- **Requirements**:
  - Cloudflare DDoS protection
  - Rate limiting per IP
  - Challenge suspicious traffic
  - Anomaly detection
- **Priority**: High

#### SR-033: CORS Configuration
- **Description**: Configure CORS properly
- **Requirements**:
  - Whitelist allowed origins
  - Allow only necessary methods
  - Allow only necessary headers
  - Credentials: true (for cookies)
- **Priority**: High

#### SR-034: CSRF Protection
- **Description**: Implement CSRF protection
- **Requirements**:
  - CSRF tokens for state-changing requests
  - Token validation
  - SameSite cookies
  - Double-submit cookie pattern
- **Priority**: High

#### SR-035: Security Headers
- **Description**: Implement security headers
- **Requirements**:
  - X-Frame-Options: DENY
  - X-Content-Type-Options: nosniff
  - X-XSS-Protection: 1; mode=block
  - Strict-Transport-Security: max-age=31536000
  - Content-Security-Policy
- **Priority**: High

#### SR-036: API Rate Limiting
- **Description**: Implement API rate limiting
- **Requirements**:
  - 100 requests/minute per IP
  - 1000 requests/minute per user
  - Different limits per endpoint
  - Rate limit headers
- **Priority**: High

#### SR-037: Request Size Limits
- **Description**: Limit request sizes
- **Requirements**:
  - Max request body: 10MB
  - Max header size: 8KB
  - Max URL length: 2048 characters
  - Reject oversized requests
- **Priority**: Medium

#### SR-038: Request Timeout
- **Description**: Set request timeouts
- **Requirements**:
  - Read timeout: 30s
  - Write timeout: 30s
  - Idle timeout: 60s
  - Graceful timeout handling
- **Priority**: Medium

#### SR-039: IP Whitelisting
- **Description**: Support IP whitelisting (optional)
- **Requirements**:
  - Whitelist admin IPs
  - Whitelist service IPs
  - Whitelist partner IPs
  - Audit whitelist changes
- **Priority**: Low

#### SR-040: Network Segmentation
- **Description**: Segment network traffic
- **Requirements**:
  - Separate public and private networks
  - Database in private network
  - Redis in private network
  - VPC configuration
- **Priority**: Medium

### 10.5 Application Security (SR-041 to SR-050)

#### SR-041: Input Validation
- **Description**: Validate all inputs
- **Requirements**:
  - Server-side validation
  - Type checking
  - Length limits
  - Format validation
- **Priority**: Critical

#### SR-042: SQL Injection Prevention
- **Description**: Prevent SQL injection attacks
- **Requirements**:
  - Use parameterized queries
  - ORM usage (GORM)
  - Input sanitization
  - Query logging
- **Priority**: Critical

#### SR-043: XSS Prevention
- **Description**: Prevent XSS attacks
- **Requirements**:
  - Output encoding
  - Input sanitization
  - Content Security Policy
  - HTTP-only cookies
- **Priority**: Critical

#### SR-044: Command Injection Prevention
- **Description**: Prevent command injection
- **Requirements**:
  - Avoid shell commands
  - Use safe APIs
  - Input validation
  - Whitelist approach
- **Priority**: High

#### SR-045: Path Traversal Prevention
- **Description**: Prevent path traversal attacks
- **Requirements**:
  - Validate file paths
  - Use canonical paths
  - Restrict file access
  - Sandboxing
- **Priority**: High

#### SR-046: File Upload Security
- **Description**: Secure file uploads
- **Requirements**:
  - File type validation
  - File size limits
  - Virus scanning (optional)
  - Rename uploaded files
- **Priority**: High

#### SR-047: Dependency Security
- **Description**: Secure dependencies
- **Requirements**:
  - Regular updates
  - Vulnerability scanning
  - SCA tools
  - SBOM generation
- **Priority**: High

#### SR-048: Error Handling
- **Description**: Handle errors securely
- **Requirements**:
  - Generic error messages (production)
  - Detailed errors (development)
  - No stack traces in production
  - Log all errors
- **Priority**: High

#### SR-049: Logging Security
- **Description**: Secure logging practices
- **Requirements**:
  - No sensitive data in logs
  - Log rotation
  - Secure log storage
  - Log access controls
- **Priority**: Medium

#### SR-050: Security Testing
- **Description**: Regular security testing
- **Requirements**:
  - Penetration testing
  - Vulnerability scanning
  - Code review
  - Security audits
- **Priority**: High

---

## 11. Performance Requirements

### 11.1 Response Time Requirements (PR-001 to PR-010)

#### PR-001: API Response Time
- **Description**: API endpoints must respond quickly
- **Requirements**:
  - Simple GET: < 100ms (p95)
  - Complex queries: < 500ms (p95)
  - Write operations: < 200ms (p95)
  - Checkout: < 1s (p95)
- **Measurement**: APM monitoring
- **Priority**: Critical

#### PR-002: Database Query Time
- **Description**: Database queries must be fast
- **Requirements**:
  - Simple queries: < 10ms (p95)
  - Complex queries: < 50ms (p95)
  - Joins: < 100ms (p95)
  - Aggregations: < 200ms (p95)
- **Measurement**: Query logging
- **Priority**: Critical

#### PR-003: Cache Response Time
- **Description**: Cache must be fast
- **Requirements**:
  - Cache hit: < 10ms (p95)
  - Cache miss: < 50ms (p95)
  - Cache write: < 20ms (p95)
  - Cache hit rate: > 80%
- **Measurement**: Redis monitoring
- **Priority**: High

#### PR-004: External API Time
- **Description**: External API calls must be optimized
- **Requirements**:
  - Midtrans: < 2s (p95)
  - Komerce: < 3s (p95)
  - Fonnte: < 5s (p95)
  - Timeout: 10s
- **Measurement**: API monitoring
- **Priority**: High

#### PR-005: File Upload Time
- **Description**: File uploads must be efficient
- **Requirements**:
  - Small files (< 1MB): < 2s
  - Medium files (1-5MB): < 5s
  - Large files (5-10MB): < 10s
  - Progress indication
- **Measurement**: Upload monitoring
- **Priority**: Medium

#### PR-006: Image Loading Time
- **Description**: Images must load quickly
- **Requirements**:
  - Thumbnails: < 500ms
  - Medium images: < 1s
  - Full images: < 2s
  - CDN delivery
- **Measurement**: Page speed tools
- **Priority**: Medium

#### PR-007: Page Load Time
- **Description**: Pages must load quickly
- **Requirements**:
  - First Contentful Paint: < 1s
  - Time to Interactive: < 3s
  - Largest Contentful Paint: < 2.5s
  - Cumulative Layout Shift: < 0.1
- **Measurement**: Web Vitals
- **Priority**: High

#### PR-008: Search Response Time
- **Description**: Search must be fast
- **Requirements**:
  - Simple search: < 200ms (p95)
  - Advanced search: < 500ms (p95)
  - Faceted search: < 1s (p95)
  - Autocomplete: < 100ms (p95)
- **Measurement**: Search monitoring
- **Priority**: High

#### PR-009: Checkout Time
- **Description**: Checkout must be fast
- **Requirements**:
  - Cart loading: < 500ms
  - Shipping calculation: < 1s
  - Payment initiation: < 2s
  - Total checkout: < 5s
- **Measurement**: Checkout analytics
- **Priority**: Critical

#### PR-010: Dashboard Load Time
- **Description**: Admin dashboard must load quickly
- **Requirements**:
  - Dashboard overview: < 1s
  - Product list: < 500ms
  - Order list: < 500ms
  - Reports: < 2s
- **Measurement**: Dashboard monitoring
- **Priority**: Medium

### 11.2 Throughput Requirements (PR-011 to PR-020)

#### PR-011: API Throughput
- **Description**: API must handle high request volume
- **Requirements**:
  - Sustained: 1,000 req/s
  - Peak: 5,000 req/s
  - Burst: 10,000 req/s (30s)
  - Error rate: < 0.1%
- **Measurement**: Load testing
- **Priority**: Critical

#### PR-012: Database Throughput
- **Description**: Database must handle high query volume
- **Requirements**:
  - Reads: 5,000 queries/s
  - Writes: 1,000 queries/s
  - Mixed: 3,000 queries/s
  - Connection pool: 100
- **Measurement**: Database monitoring
- **Priority**: Critical

#### PR-013: Cache Throughput
- **Description**: Cache must handle high request volume
- **Requirements**:
  - Reads: 10,000 ops/s
  - Writes: 5,000 ops/s
  - Memory: 4GB minimum
  - Eviction policy: LRU
- **Measurement**: Redis monitoring
- **Priority**: High

#### PR-014: Concurrent Users
- **Description**: System must support many concurrent users
- **Requirements**:
  - Active users: 10,000
  - Concurrent sessions: 5,000
  - Concurrent checkouts: 100
  - Session timeout: 30 min
- **Measurement**: User analytics
- **Priority**: High

#### PR-015: Order Processing
- **Description**: System must process orders efficiently
- **Requirements**:
  - Orders per minute: 100
  - Orders per hour: 6,000
  - Orders per day: 100,000
  - Processing time: < 5s
- **Measurement**: Order analytics
- **Priority**: Critical

#### PR-016: Payment Processing
- **Description**: Payment processing must be efficient
- **Requirements**:
  - Payments per minute: 50
  - Payment success rate: > 95%
  - Payment retry: 3 attempts
  - Webhook processing: < 1s
- **Measurement**: Payment analytics
- **Priority**: Critical

#### PR-017: Notification Processing
- **Description**: Notifications must be processed efficiently
- **Requirements**:
  - Notifications per minute: 200
  - Queue depth: 1,000
  - Delivery rate: > 98%
  - Processing time: < 2s
- **Measurement**: Notification analytics
- **Priority**: High

#### PR-018: File Upload Throughput
- **Description**: File uploads must be efficient
- **Requirements**:
  - Uploads per minute: 50
  - Concurrent uploads: 10
  - Bandwidth: 1 Gbps
  - Storage: Unlimited (R2)
- **Measurement**: Upload monitoring
- **Priority**: Medium

#### PR-019: Search Throughput
- **Description**: Search must handle high query volume
- **Requirements**:
  - Searches per second: 100
  - Search latency: < 500ms (p95)
  - Index size: < 10GB
  - Update latency: < 1s
- **Measurement**: Search monitoring
- **Priority**: High

#### PR-020: Report Generation
- **Description**: Reports must be generated efficiently
- **Requirements**:
  - Daily reports: < 1 min
  - Weekly reports: < 5 min
  - Monthly reports: < 15 min
  - Export formats: CSV, PDF
- **Measurement**: Report analytics
- **Priority**: Medium

### 11.3 Scalability Requirements (PR-021 to PR-030)

#### PR-021: Horizontal Scaling
- **Description**: System must scale horizontally
- **Requirements**:
  - Stateless application
  - Load balancer support
  - Auto-scaling capability
  - Max instances: 20
- **Measurement**: Scaling tests
- **Priority**: Critical

#### PR-022: Database Scaling
- **Description**: Database must scale
- **Requirements**:
  - Read replicas: 3
  - Connection pooling
  - Query optimization
  - Index optimization
- **Measurement**: Database monitoring
- **Priority**: Critical

#### PR-023: Cache Scaling
- **Description**: Cache must scale
- **Requirements**:
  - Redis Cluster
  - Sharding support
  - Automatic failover
  - Memory: 16GB max
- **Measurement**: Cache monitoring
- **Priority**: High

#### PR-024: Storage Scaling
- **Description**: Storage must scale
- **Requirements**:
  - Cloudflare R2 (unlimited)
  - CDN delivery
  - Automatic optimization
  - Bandwidth: 10 TB/month
- **Measurement**: Storage monitoring
- **Priority**: Medium

#### PR-025: Queue Scaling
- **Description**: Message queues must scale
- **Requirements**:
  - Redis-based queues
  - Worker scaling
  - Queue depth: 10,000
  - Workers: 10 max
- **Measurement**: Queue monitoring
- **Priority**: Medium

#### PR-026: CDN Scaling
- **Description**: CDN must handle traffic
- **Requirements**:
  - Edge locations: Global
  - Bandwidth: 10 TB/month
  - Cache hit rate: > 90%
  - TTL: Configurable
- **Measurement**: CDN analytics
- **Priority**: Medium

#### PR-027: Load Balancing
- **Description**: Load must be distributed evenly
- **Requirements**:
  - Round-robin algorithm
  - Health checks
  - Session persistence
  - Sticky sessions
- **Measurement**: Load balancer metrics
- **Priority**: High

#### PR-028: Auto-scaling Rules
- **Description**: Define auto-scaling rules
- **Requirements**:
  - CPU > 70%: Scale up
  - CPU < 30%: Scale down
  - Memory > 80%: Scale up
  - Response time > 1s: Scale up
- **Measurement**: Auto-scaling logs
- **Priority**: High

#### PR-029: Graceful Scaling
- **Description**: Scaling must be graceful
- **Requirements**:
  - Zero-downtime deployment
  - Rolling updates
  - Connection draining
  - Health checks
- **Measurement**: Deployment monitoring
- **Priority**: High

#### PR-030: Capacity Planning
- **Description**: Plan for future growth
- **Requirements**:
  - 6-month forecast
  - 12-month forecast
  - Resource allocation
  - Budget planning
- **Measurement**: Capacity reports
- **Priority**: Medium

---

## 12. Deployment Architecture

### 12.1 Development Environment

#### Infrastructure
- **Container Engine**: Podman (rootless)
- **Orchestration**: podman-compose
- **OS**: Windows with WSL2 (Ubuntu/Debian)
- **Local Development**: localhost:8080

#### Services
```
┌─────────────────────────────────────────────────────────────┐
│                    Development Stack                         │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │  Backend    │  │  Kratos     │  │  Kratos UI  │        │
│  │  :8080      │  │  :4433/4434 │  │  :4455      │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
│  ┌─────────────┐  ┌─────────────┐                         │
│  │ PostgreSQL  │  │   Redis     │                         │
│  │  :5432      │  │  :6379      │                         │
│  └─────────────┘  └─────────────┘                         │
└─────────────────────────────────────────────────────────────┘
```

#### Configuration
- **Environment File**: `.env`
- **Database**: PostgreSQL (local)
- **Cache**: Redis (local)
- **Storage**: Local filesystem
- **Payment**: Midtrans Sandbox
- **Shipping**: Komerce Sandbox

### 12.2 Production Environment

#### Infrastructure
- **VPS Provider**: Cloud provider (AWS/GCP/DigitalOcean)
- **Container Engine**: Docker
- **Orchestration**: Docker Compose / Kubernetes (future)
- **Load Balancer**: Cloudflare / Nginx
- **CDN**: Cloudflare

#### Domains
```
┌─────────────────────────────────────────────────────────────┐
│                    Production Domains                        │
├─────────────────────────────────────────────────────────────┤
│  Customer Storefront: https://karima.com                    │
│  Admin Panel:         https://admin.ks-backend.cloud        │
│  API Backend:         https://api.ks-backend.cloud          │
│  Auth Service:        https://auth.ks-backend.cloud         │
└─────────────────────────────────────────────────────────────┘
```

#### Services
```
┌─────────────────────────────────────────────────────────────┐
│                    Production Stack                          │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────────────────────────────────────────────┐   │
│  │              Cloudflare CDN / Load Balancer          │   │
│  └─────────────────────────────────────────────────────┘   │
│                              │                               │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │  Backend    │  │  Kratos     │  │  Frontend   │        │
│  │  (Docker)   │  │  (Docker)   │  │  (Pages)    │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
│  ┌─────────────┐  ┌─────────────┐                         │
│  │ PostgreSQL  │  │   Redis     │                         │
│  │  (VPS)      │  │  (VPS)      │                         │
│  └─────────────┘  └─────────────┘                         │
└─────────────────────────────────────────────────────────────┘
```

#### Configuration
- **Environment File**: `.env.production`
- **Database**: PostgreSQL (VPS with SSL)
- **Cache**: Redis (VPS with password)
- **Storage**: Cloudflare R2
- **Payment**: Midtrans Production
- **Shipping**: Komerce Production

### 12.3 Deployment Process

#### CI/CD Pipeline
```
┌─────────────────────────────────────────────────────────────┐
│                    CI/CD Pipeline                           │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  1. Code Push → GitHub                                      │
│  2. Trigger → GitHub Actions                               │
│  3. Build → Docker Image                                    │
│  4. Test → Unit + Integration Tests                         │
│  5. Security Scan → Snyk / Dependabot                       │
│  6. Deploy → Staging Environment                            │
│  7. E2E Tests → Playwright / Cypress                       │
│  8. Manual Review → QA Team                                 │
│  9. Deploy → Production (Blue-Green)                        │
│ 10. Health Check → Automated Monitoring                    │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

#### Deployment Commands
```bash
# Development
make kratos-up          # Start all services
make kratos-down        # Stop all services
make logs               # View logs
make migrate            # Run migrations

# Production
docker build -t karima_store_backend:v1.0.0 .
docker-compose -f docker-compose.prod.yml up -d
docker-compose -f docker-compose.prod.yml ps
docker-compose -f docker-compose.prod.yml logs -f
```

### 12.4 Monitoring & Observability

#### Monitoring Stack
```
┌─────────────────────────────────────────────────────────────┐
│                    Monitoring Stack                          │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │ Prometheus  │  │ Grafana     │  │ Alertmanager│        │
│  │ (Metrics)   │  │ (Dashboards)│  │ (Alerts)    │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
│  ┌─────────────┐  ┌─────────────┐                         │
│  │ Loki        │  │ Jaeger      │                         │
│  │ (Logs)      │  │ (Tracing)   │                         │
│  └─────────────┘  └─────────────┘                         │
└─────────────────────────────────────────────────────────────┘
```

#### Metrics to Monitor
- **Application Metrics**:
  - Request rate
  - Response time (p50, p95, p99)
  - Error rate
  - Active connections

- **Database Metrics**:
  - Query performance
  - Connection pool usage
  - Replication lag
  - Disk usage

- **Cache Metrics**:
  - Hit rate
  - Memory usage
  - Eviction rate
  - Operations per second

- **Business Metrics**:
  - Orders per minute
  - Revenue per hour
  - Conversion rate
  - Customer satisfaction

#### Alerting Rules
- **Critical Alerts**:
  - Error rate > 1%
  - Response time > 1s (p95)
  - Database down
  - Redis down
  - Payment gateway down

- **Warning Alerts**:
  - Error rate > 0.1%
  - Response time > 500ms (p95)
  - CPU > 80%
  - Memory > 80%
  - Disk space < 20%

### 12.5 Backup & Recovery

#### Backup Strategy
```
┌─────────────────────────────────────────────────────────────┐
│                    Backup Strategy                           │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  Database Backups:                                           │
│  - Daily full backups: 2:00 AM UTC                          │
│  - Hourly incremental backups                               │
│  - Retention: 30 days                                        │
│  - Off-site storage: Cloudflare R2                          │
│                                                              │
│  Media Backups:                                              │
│  - Real-time sync to R2                                     │
│  - Versioning enabled                                       │
│  - Retention: 90 days                                       │
│                                                              │
│  Configuration Backups:                                     │
│  - Environment variables                                    │
│  - Docker configurations                                    │
│  - Kratos configurations                                   │
│  - Stored in Git (encrypted)                                │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

#### Recovery Procedures
- **RTO (Recovery Time Objective)**: < 1 hour
- **RPO (Recovery Point Objective)**: < 15 minutes
- **Recovery Steps**:
  1. Identify failure point
  2. Restore from latest backup
  3. Verify data integrity
  4. Restart services
  5. Run health checks
  6. Monitor for issues

---

## 13. Success Metrics

### 13.1 Business Metrics

#### Revenue Metrics
- **Total Revenue**: Target: $1M/year
- **Average Order Value**: Target: $50
- **Conversion Rate**: Target: 3%
- **Repeat Purchase Rate**: Target: 30%

#### Customer Metrics
- **Active Customers**: Target: 10,000
- **Customer Acquisition Cost**: Target: <$10
- **Customer Lifetime Value**: Target: $200
- **Customer Satisfaction**: Target: 4.5/5

#### Operational Metrics
- **Order Processing Time**: Target: < 5 minutes
- **Order Fulfillment Rate**: Target: 98%
- **Return Rate**: Target: < 5%
- **Support Ticket Volume**: Target: < 100/month

### 13.2 Technical Metrics

#### Performance Metrics
- **API Response Time**: Target: < 100ms (p95)
- **Page Load Time**: Target: < 2s
- **Uptime**: Target: 99.9%
- **Error Rate**: Target: < 0.1%

#### Scalability Metrics
- **Concurrent Users**: Target: 10,000
- **Requests per Second**: Target: 1,000
- **Database Query Time**: Target: < 50ms (p95)
- **Cache Hit Rate**: Target: > 80%

#### Security Metrics
- **Security Incidents**: Target: 0
- **Vulnerabilities**: Target: 0 critical, < 5 high
- **Failed Login Attempts**: Target: < 1%
- **Successful Penetration Tests**: Target: 100%

### 13.3 Development Metrics

#### Code Quality Metrics
- **Test Coverage**: Target: > 80%
- **Code Review Rate**: Target: 100%
- **Linting Errors**: Target: 0
- **Technical Debt**: Target: Low

#### Delivery Metrics
- **Deployment Frequency**: Target: Weekly
- **Lead Time**: Target: < 2 days
- **Mean Time to Recovery**: Target: < 1 hour
- **Change Failure Rate**: Target: < 5%

---

## 14. Roadmap

### 14.1 Phase 1: Foundation (Completed - Q4 2025)
**Status**: ✅ Completed

**Deliverables**:
- ✅ Project infrastructure setup
- ✅ Database schema design
- ✅ Core models and repositories
- ✅ Product management API
- ✅ Category management API
- ✅ Media management API
- ✅ Basic authentication setup

### 14.2 Phase 2: Core Features (Completed - Q4 2025)
**Status**: ✅ Completed

**Deliverables**:
- ✅ Pricing engine implementation
- ✅ Coupon system
- ✅ Shopping cart functionality
- ✅ Checkout process
- ✅ Midtrans payment integration
- ✅ Komerce shipping integration
- ✅ Order management
- ✅ Stock management
- ✅ Redis caching
- ✅ WhatsApp notifications

### 14.3 Phase 3: Authentication & Security (In Progress - Q1 2026)
**Status**: 🔄 In Progress (70%)

**Deliverables**:
- ✅ Ory Kratos integration
- ✅ Kratos middleware
- ✅ Identity schema
- ⏳ User registration flow
- ⏳ User login flow
- ⏳ Session management
- ⏳ RBAC implementation
- ⏳ API key authentication
- ⏳ Security headers
- ⏳ Rate limiting

**Remaining Tasks**:
- Complete user registration and login flows
- Implement RBAC (Role-Based Access Control)
- Sync user data from Kratos to local database
- Security audit and penetration testing

### 14.4 Phase 4: Advanced Features (Planned - Q1 2026)
**Status**: 📋 Planned

**Deliverables**:
- Order history and tracking
- Refund processing
- Product reviews and ratings
- Wishlist functionality
- Flash sale management
- Advanced reporting
- Admin dashboard
- Customer dashboard

### 14.5 Phase 5: Optimization & Scaling (Planned - Q2 2026)
**Status**: 📋 Planned

**Deliverables**:
- Performance optimization
- Database query optimization
- Cache optimization
- Load testing and tuning
- Auto-scaling configuration
- CDN optimization
- Image optimization
- Database read replicas

### 14.6 Phase 6: Integrations & Extensions (Planned - Q2 2026)
**Status**: 📋 Planned

**Deliverables**:
- Social login (Google, Facebook)
- Email notifications (SMTP)
- SMS notifications
- Push notifications
- Analytics integration (Google Analytics)
- CRM integration
- Accounting integration
- Multi-language support

### 14.7 Phase 7: Production Launch (Planned - Q3 2026)
**Status**: 📋 Planned

**Deliverables**:
- Production environment setup
- Domain configuration
- SSL certificates
- Load balancer setup
- Monitoring and alerting
- Backup and recovery
- Security hardening
- Performance testing
- User acceptance testing
- Launch preparation

### 14.8 Phase 8: Post-Launch Support (Planned - Q3 2026 onwards)
**Status**: 📋 Planned

**Deliverables**:
- Bug fixes and patches
- Feature enhancements
- Performance improvements
- Security updates
- Customer support
- Documentation updates
- Training and onboarding

---

## 15. Appendices

### 15.1 Glossary

| Term | Definition |
|------|------------|
| **API** | Application Programming Interface |
| **RBAC** | Role-Based Access Control |
| **SKU** | Stock Keeping Unit |
| **PCI DSS** | Payment Card Industry Data Security Standard |
| **TLS** | Transport Layer Security |
| **CDN** | Content Delivery Network |
| **CI/CD** | Continuous Integration/Continuous Deployment |
| **APM** | Application Performance Monitoring |
| **RTO** | Recovery Time Objective |
| **RPO** | Recovery Point Objective |
| **OWASP** | Open Web Application Security Project |
| **JWT** | JSON Web Token |
| **ORM** | Object-Relational Mapping |
| **TTL** | Time To Live |
| **LRU** | Least Recently Used |
| **PII** | Personally Identifiable Information |

### 15.2 API Endpoints Summary

#### Product Management
- `GET /api/v1/products` - List products
- `POST /api/v1/products` - Create product
- `GET /api/v1/products/:id` - Get product details
- `PUT /api/v1/products/:id` - Update product
- `DELETE /api/v1/products/:id` - Delete product
- `POST /api/v1/products/:id/stock` - Update stock

#### Category Management
- `GET /api/v1/categories` - List categories
- `POST /api/v1/categories` - Create category
- `GET /api/v1/categories/:id` - Get category details
- `PUT /api/v1/categories/:id` - Update category
- `DELETE /api/v1/categories/:id` - Delete category

#### Cart & Checkout
- `POST /api/v1/cart` - Add to cart
- `PUT /api/v1/cart/:id` - Update cart item
- `DELETE /api/v1/cart/:id` - Remove from cart
- `GET /api/v1/cart` - View cart
- `POST /api/v1/checkout` - Initiate checkout
- `POST /api/v1/webhooks/midtrans` - Payment webhook

#### Order Management
- `GET /api/v1/orders` - List orders
- `GET /api/v1/orders/:id` - Get order details
- `PUT /api/v1/orders/:id/status` - Update order status
- `POST /api/v1/orders/:id/cancel` - Cancel order
- `POST /api/v1/orders/:id/refund` - Process refund

#### Pricing
- `POST /api/v1/pricing/calculate` - Calculate price
- `POST /api/v1/pricing/order-summary` - Calculate order summary
- `POST /api/v1/coupons/validate` - Validate coupon

#### Shipping
- `GET /api/v1/shipping/destination/search` - Search destinations
- `GET /api/v1/shipping/calculate` - Calculate shipping cost

#### Media
- `POST /api/v1/media/upload` - Upload media
- `DELETE /api/v1/media/:id` - Delete media
- `PUT /api/v1/media/:id/primary` - Set primary image

#### Authentication
- `POST /api/v1/auth/register` - Register user
- `POST /api/v1/auth/login` - Login user
- `POST /api/v1/auth/logout` - Logout user
- `POST /api/v1/auth/refresh` - Refresh token

#### Notifications
- `POST /api/v1/notifications/send` - Send notification
- `GET /api/v1/notifications` - List notifications

#### Health & Metrics
- `GET /health` - Health check
- `GET /metrics` - Application metrics

### 15.3 Database Schema Overview

#### Core Tables
- `products` - Product information
- `product_skus` - Product variants
- `product_media` - Product images/videos
- `categories` - Product categories
- `users` - User accounts
- `carts` - Shopping cart items
- `orders` - Order information
- `order_items` - Order line items
- `coupons` - Discount codes
- `flash_sales` - Promotional campaigns
- `reviews` - Customer reviews
- `wishlists` - Saved products
- `stock_logs` - Inventory audit trail
- `shipping_zones` - Shipping regions
- `tax_rates` - Tax configuration

#### Relationships
- Products → Categories (Many-to-One)
- Products → Product SKUs (One-to-Many)
- Products → Product Media (One-to-Many)
- Users → Carts (One-to-Many)
- Users → Orders (One-to-Many)
- Orders → Order Items (One-to-Many)
- Orders → Shipping (One-to-One)
- Orders → Payments (One-to-Many)

### 15.4 Environment Variables Reference

#### Application Configuration
```env
APP_ENV=development|production
APP_PORT=8080
API_VERSION=v1
```

#### Database Configuration
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=karima_store
DB_PASSWORD=your_password
DB_NAME=karima_db
DB_SSL_MODE=disable|require
```

#### Redis Configuration
```env
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
```

#### Ory Kratos Configuration
```env
KRATOS_PUBLIC_URL=http://127.0.0.1:4433
KRATOS_ADMIN_URL=http://127.0.0.1:4434
KRATOS_UI_URL=http://127.0.0.1:4455
JWT_SECRET=your_jwt_secret
JWT_EXPIRATION=24h
```

#### Midtrans Configuration
```env
MIDTRANS_SERVER_KEY=your_server_key
MIDTRANS_CLIENT_KEY=your_client_key
MIDTRANS_IS_PRODUCTION=false|true
MIDTRANS_API_BASE_URL=https://app.sandbox.midtrans.com/snap/v1
```

#### Komerce Configuration
```env
RAJAONGKIR_API_KEY=your_api_key
RAJAONGKIR_BASE_URL=https://api-sandbox.collaborator.komerce.id/tariff/api/v1/
```

#### Cloudflare R2 Configuration
```env
FILE_STORAGE=local|r2
FILE_UPLOAD_MAX_SIZE=10MB
R2_ACCOUNT_ID=your_account_id
R2_ENDPOINT=https://your_account_id.r2.cloudflarestorage.com
R2_ACCESS_KEY_ID=your_access_key_id
R2_SECRET_ACCESS_KEY=your_secret_access_key
R2_BUCKET_NAME=karima-media
R2_PUBLIC_URL=https://your-custom-domain.com
```

#### Fonnte Configuration
```env
FONNTE_TOKEN=your_fonnte_token
FONNTE_URL=https://api.fonnte.com/send
```

#### CORS Configuration
```env
CORS_ORIGIN=http://localhost:3000,https://karima.com
```

### 15.5 References

#### Documentation
- [Fiber Documentation](https://docs.gofiber.io/)
- [GORM Documentation](https://gorm.io/docs/)
- [Ory Kratos Documentation](https://www.ory.sh/docs/kratos)
- [Midtrans Documentation](https://docs.midtrans.com/)
- [Komerce API Documentation](https://api.komerce.id/)
- [Fonnte Documentation](https://fonnte.com/api)
- [Cloudflare R2 Documentation](https://developers.cloudflare.com/r2/)
- [Go Documentation](https://golang.org/doc/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Redis Documentation](https://redis.io/documentation)

#### Standards
- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [PCI DSS](https://www.pcisecuritystandards.org/)
- [REST API Design Best Practices](https://restfulapi.net/)
- [JSON API Specification](https://jsonapi.org/)
- [OpenAPI Specification](https://swagger.io/specification/)

#### Tools
- [Swagger UI](https://swagger.io/tools/swagger-ui/)
- [Postman](https://www.postman.com/)
- [Docker](https://www.docker.com/)
- [Podman](https://podman.io/)
- [GitHub Actions](https://github.com/features/actions)
- [Prometheus](https://prometheus.io/)
- [Grafana](https://grafana.com/)
- [Jaeger](https://www.jaegertracing.io/)

---

## Document Information

| Field | Value |
|-------|-------|
| **Document Title** | Product Requirements Document (PRD) |
| **Project Name** | Karima Store - Fashion E-commerce Backend |
| **Version** | 1.0.0 |
| **Date Created** | January 2, 2026 |
| **Last Updated** | January 2, 2026 |
| **Author** | Development Team |
| **Status** | In Progress |
| **Review Status** | Pending Review |

---

## Change History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0.0 | 2026-01-02 | Development Team | Initial PRD creation |

---

**End of Document**
