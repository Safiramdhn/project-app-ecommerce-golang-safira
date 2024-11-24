# **E-Commerce API README**

## **Overview**

This is an E-Commerce API built using Go (Golang) that provides various functionalities for managing products, users, orders, and recommendations. The API is designed to facilitate the operations of an e-commerce platform.

## **API Endpoints**

### **User Management**

### **Register User**

- **Endpoint:** **`POST /api/register`**
- **Request Body:**

```
{
    "name": "John Doe",
    "emailOrPhoneNumber": "john.doe@example.com",
    "password": "password123"
}

```

- **Response:**

```
{
    "status": "success",
    "message": "User  created successfully",
    "data": "user_id"
}
```

### **User Login**

- **Endpoint:** **`GET /api/login`**
- **Request Body:**

```
{
    "emailOrPhoneNumber": "john.doe@example.com",
    "password": "password123"
}
```

- **Response:**

```
{
    "status": "success",
    "message": "Login successful",
    "data": {
        "token": "your_jwt_token"
    }
}

```

### **Product Management**

### **Get All Products**

- **Endpoint:** **`GET /api/products`**
- **Query Parameters:**
    - **`page`**: (optional, default: 1)
    - **`perPage`**: (optional, default: 5)
    - **`name`**: (optional, filter by product name)
    - **`categoryId`**: (optional, filter by category ID)
- **Response:**

```
{
    "status": "success",
    "message": "Products successfully retrieved",
   "data": [
        {
            "id": 1,
            "name": "Product 1",
            "description": "Description of Product 1",
            "price": 100,
            "discount": 10,
            "rating": 4.5,
            "photo_url": "http://example.com/product1.jpg"
        }
        // ... more products
    ],
    "pagination": {
        "page": 1,
        "limit": 5,
        "total_items": 100,
        "total_pages": 20
    }
}

```

### **Get Product By ID**

- **Endpoint:** **`GET /api/products/{id}`**
- **Response:**

```
{
    "status": "success",
    "message": "Product successfully retrieved",
    "data": {
        "id": 1,
        "name": "Product 1",
        "description": "Description of Product 1",
        "price": 100,
        "discount": 10,
        "rating": 4.5,
        "photo_url": "http://example.com/product1.jpg"
    }
}

```

### **Order Management**

### **Create Order**

- **Endpoint:** **`POST /api/orders`**
- **Request Body:**

```
{
    "cartID": 1,
	  "addressID": 1,
    "shippingType": "standard",
    "shippingCost": 5,
    "paymentMethod": "credit_card"
}

```

- **Response:**

```
{
    "status": "success",
    "message": "Order created successfully"
}

```

### **Get Order History**

- **Endpoint:** **`GET /api/orders`**
- **Response:**

```
{
    "status": "success",
    "message": "Order history successfully retrieved",
    "data": [
        {
            "id": 1,
            "total_amount": 150,
            "total_price": 155,
            "order_status": "completed"
        }
        // ... more orders
    ]
}

```

### **Wishlist Management**

### **Add Product to Wishlist**

- **Endpoint:** **`POST /api/wishlist/add`**
- **Request Body:**

```
{
    "productID": 1
}

```

- **Response:**

```
{
    "status": "success",
    "message": "Product added to wishlist successfully"
}

```

### **Get Wishlist**

- **Endpoint:** **`GET /api/wishlist`**
- **Response:**

```
{
    "status": "success",
    "message": "Wishlist successfully retrieved",
   "data": [
        {
            "id": 1,
           "productID": 1,
            "product": {
                "id": 1,
                "name": "Product 1",
                "photo_url": "http://example.com/product1.jpg"
            }
        }
        // ... more wishlist items
    ],
    "pagination": {
        "page": 1,
        "limit": 5,
        "total_items": 10,
        "total_pages": 2
    }
}

```

## **Running the API**

1. Clone the repository:
    
    ```
    git clone https://github .com/yourusername/ecommerce-api.git
    
    ```
    
2. Navigate to the project directory:
    
    ```
    cd ecommerce-api
    
    ```
    
3. Install the dependencies:
    
    ```
    go mod tidy
    
    ```
    
4. Run the application:
    
    ```
    go run main.go
    
    ```
    

## **Sample API Requests**

### **Register User**

```
curl -X POST http://localhost:8080/api/register \
-H "Content-Type: application/json" \
d '{
   "name": "John Doe",
    "emailOrPhoneNumber": "john.doe@example.com",
    "password": "password123"
}'

```

### **User Login**

```
curl -X GET http://localhost:8080/api/login \
-H "Content-Type: application/json" \
-d '{
    "emailOrPhoneNumber": "john.doe@example.com",
    "password": "password123"
}'

```

### **Get All Products**

```
url -X GET "http://localhost:8080/api/products?page=1&perPage=5"

```

### **Create Order**

```
curl -X POST http://localhost:8080/api/orders \
-H "Content-Type: application/json" \
-d '{
    "cartID": 1,
    "addressID": 1,
    "shippingType": "standard",
    "shippingCost": 5,
    "paymentMethod": "credit_card"
}'

```

### **Add Product to Wishlist**

```
curl -X POST http://localhost:8080/api/wishlist/add \
-H "Content-Type: application/json" \
-d '{
    "productID": 1
}'

```

## **Conclusion**

This README provides an overview of the E-Commerce API, including its endpoints, request and response formats, and sample API requests. For further details, please refer to the codebase or the documentation within the repository.