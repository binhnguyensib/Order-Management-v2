# RESTful API CRUD: Order Management System with Golang + Gin + MongoDB + Swagger

## 🌟 Overview

Build a simple RESTful API for managing a sales system with the following entities:

* **Customer**
* **Product**
* **Order**

Use the following stack:

* **Golang** with **Gin framework**
* **MongoDB** for database
* **Swagger UI** for API documentation

---

## 📊 Entities and CRUD Features

### 1. Customer

#### Endpoints

* `POST /customers`: Create a new customer
* `GET /customers`: Get list of customers
* `GET /customers/:id`: Get customer by ID
* `PUT /customers/:id`: Update customer
* `DELETE /customers/:id`: Delete customer

#### Model

```json
{
  "id": "string",
  "name": "string",
  "email": "string",
  "phone": "string"
}
```

---

### 2. Product

#### Endpoints

* `POST /products`: Add a new product
* `GET /products`: Get list of products
* `GET /products/:id`: Get product by ID
* `PUT /products/:id`: Update product
* `DELETE /products/:id`: Delete product

#### Model

```json
{
  "id": "string",
  "name": "string",
  "price": 100000,
  "stock": 20
}
```

---

### 3. Order

#### Endpoints

* `POST /orders`: Create a new order
* `GET /orders`: Get list of orders
* `GET /orders/:id`: Get order by ID
* `PUT /orders/:id`: Update order
* `DELETE /orders/:id`: Delete order

#### Model

```json
{
  "id": "string",
  "customerId": "string",
  "productIds": ["string"],
  "totalAmount": 200000,
  "createdAt": "datetime"
}
```

> ✨ Note: `totalAmount` can be auto-calculated from the sum of product prices.

---

## 📃 Technical Requirements

* Use **Gin** for HTTP routing and handling
* Use **MongoDB** to store all entities
* Use **Swagger** (`swaggo/swag`) for API documentation
* Validate inputs and handle errors gracefully
* MongoDB URI should be configurable using environment variables

---

## 🔧 Recommended Packages

| Purpose        | Package                              |
| -------------- | ------------------------------------ |
| Web framework  | github.com/gin-gonic/gin             |
| MongoDB driver | go.mongodb.org/mongo-driver/mongo    |
| Swagger docs   | github.com/swaggo/gin-swagger + swag |
| Env config     | github.com/joho/godotenv             |

---

## ✅ Expected Output

* Fully working REST API for managing `Customer`, `Product`, `Order`
* All data saved and queried from MongoDB
* Swagger UI available at: `GET /swagger/index.html`
* Project includes `README.md` for setup and running instructions

---

## 📂 Suggested Project Structure

```
.
├── main.go
├── controllers/
│   ├── customer.go
│   ├── product.go
│   └── order.go
├── models/
│   ├── customer.go
│   ├── product.go
│   └── order.go
├── routes/
│   └── routes.go
├── database/
│   └── mongodb.go
├── docs/                 # swagger generated
├── .env
└── go.mod
```

---
