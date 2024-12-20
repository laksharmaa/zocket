# Product Management System
A robust product management API built with Go, featuring image processing, caching, and message queuing capabilities.
🚀 Features
RESTful API endpoints for product management
Image processing and compression using Cloudinary
Caching with Redis
Message queuing with RabbitMQ
PostgreSQL database for data persistence
Graceful shutdown handling
Structured logging
📋 Prerequisites
Go 1.19 or higher
Docker and Docker Compose
PostgreSQL
Redis
RabbitMQ
Cloudinary account
🛠️ Installation
Clone the repository:
bash

```
git clone https://github.com/yourusername/product-management-system.git
cd product-management-system
```
## Create a .env file:
```
env


DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=lakshya0005
DB_NAME=product_management
REDIS_URL=localhost:6379
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
CLOUDINARY_NAME=your_cloudinary_name
CLOUDINARY_KEY=your_cloudinary_key
CLOUDINARY_SECRET=your_cloudinary_secret
```
## Start the required services using Docker Compose:

```
docker-compose up -d
Initialize the database:
sql


CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    product_description TEXT,
    product_price DECIMAL(10,2) NOT NULL,
    product_images TEXT[],
    compressed_product_images TEXT[],
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_products_user_id ON products(user_id);
```
## Run the application:


```
go run cmd/api/main.go
```
## 🔄 API Endpoints
### Root Endpoint
```

GET /
Response: Welcome message and API status
```
##Products
```


POST /api/v1/products
```
## Create a new product
```
GET /api/v1/products/{id}
Get a specific product by ID

GET /api/v1/products?user_id={id}&min_price={price}&max_price={price}&product_name={name}
```

## 📝 API Usage Examples
### Create Product
```


curl -X POST http://localhost:8080/api/v1/products \
-H "Content-Type: application/json" \
-d '{
    "user_id": 1,
    "product_name": "Gaming Laptop",
    "product_description": "High-performance gaming laptop",
    "product_price": 1999.99,
    "product_images": ["https://example.com/laptop1.jpg"]
}'
```
## Get Product

```

curl http://localhost:8080/api/v1/products/1
```
## Get Products with Filters
```


curl "http://localhost:8080/api/v1/products?user_id=1&min_price=1000&max_price=2000&product_name=Gaming"
```
## 🏗️ Project Structure
```


.
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   ├── middleware/
│   │   └── routes/
│   ├── config/
│   ├── models/
│   └── service/
├── pkg/
│   ├── cache/
│   ├── cloudinary/
│   ├── database/
│   └── queue/
├── docker-compose.yml
└── README.md
```
## ⚙️ Configuration
The application uses environment variables for configuration. See the .env file section above for required variables.
## 🔒 Security
Password authentication for PostgreSQL
AMQP authentication for RabbitMQ
API rate limiting
Secure image processing
🛟 Error Handling
The API returns appropriate HTTP status codes:
200: Success
201: Created
400: Bad Request
404: Not Found
500: Internal Server Error