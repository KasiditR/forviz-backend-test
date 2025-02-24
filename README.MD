# Forviz Book RESTful API

```bash
docker-compose up -d
go run ./cmd/main.go
```

### <span style="color:yellow">POST</span> /api/v1/auth/register

Request

```json
{
    "username": "kasidit",
    "password": "123456"
}
```

Response

```json
{
    "message": "User registered successfully"
}
```

### <span style="color:yellow">POST</span> /api/v1/auth/login

Request

```json
{
    "username": "kasidit",
    "password": "123456"
}
```

Response

```json
{
    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6IjY3YmM3MDYzODJhMzc5NDk3YjAyYjEwNyIsIlVzZXJOYW1lIjoiSm9obiIsIkltYWdlIjoiIiwiZXhwIjoxNzQwNDg5NTcxfQ.xY_XOIgtNEdGxjwE2XBG6DuaUcEL7XngHekaTB8zUbY",
    "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6IiIsIlVzZXJOYW1lIjoiIiwiSW1hZ2UiOiIiLCJleHAiOjE3NDEwMDc5NzF9.1AcCHFWAlPqhWO9yWeZUoFeaqhyhccP2Ga21cG1TRrE",
    "user_id": "67bc706382a379497b02b107"
}
```

### <span style="color:yellow">POST</span> /api/v1/auth/refresh-token

Request

```json
{
    "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6IiIsIlVzZXJOYW1lIjoiIiwiSW1hZ2UiOiIiLCJleHAiOjE3NDEwMDc5NzF9.1AcCHFWAlPqhWO9yWeZUoFeaqhyhccP2Ga21cG1TRrE"
}
```

Response

```json
{
    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6IiIsIlVzZXJOYW1lIjoiIiwiZXhwIjoxNzQwNDk1NDU1fQ.kSIStQJn4Ye4-NFvW0qOF5Z7Xt_0sMphISxXG-BCsxU"
}
```

### <span style="color:yellow">POST</span> /api/v1/book/add

Request

```json
{
    "book_name":"Game of Thrones",
    "author":"George R. R. Martin",
    "category":"Novel"
}
```

Response

```json
{
    "message": "Add book successfully"
}
```

### <span style="color:red">DELETE</span> /api/v1/book/delete/:id

Authorization : Bearer JWT

Response

```json
{
    "message": "Delete book successfully"
}
```

### <span style="color:rgb(71, 133, 247)">PUT</span> /api/v1/book/edit

Authorization : Bearer JWT

Request

```json
{
    "book_id":"67bc8101112ff4d3c99effbc",
    "book_name":"johnn",
    "category":"cartoon"
}
```

Response

```json
{
    "message": "Edit book successfully"
}
```

### <span style="color:rgb(63, 235, 48)">GET</span> /api/v1/book/detail/:id

Authorization : Bearer JWT

Response

```json
{
    "book_id": "67bc73ec9c3de977b373af92",
    "book_name": "block",
    "author": "boss",
    "category": "cartoon"
}
```

### <span style="color:rgb(63, 235, 48)">GET</span> /api/v1/book/search?bookName=&authorName=&categoryName=

Authorization : Bearer JWT

Response

```json
{
    "books": [
        {
            "book_id": "67bc8101112ff4d3c99effbc",
            "book_name": "johnn",
            "author": "boss",
            "category": "cartoon"
        }
    ]
}
```

### <span style="color:rgb(63, 235, 48)">GET</span> /api/v1/book/most-borrow

Authorization : Bearer JWT

Response

```json
{
    "books": [
        {
            "book_id": "67bc8101112ff4d3c99effbc",
            "book_name": "johnn",
            "author": "boss",
            "category": "cartoon",
            "borrowed_count": 1
        }
    ]
}
```
### <span style="color:yellow">POST</span> /api/v1/book/borrow

Authorization : Bearer JWT

Request

```json
{
    "book_id":"67bc8101112ff4d3c99effbc"
}
```

Response

```json
{
    "message": "Borrow book successfully"
}
```

### <span style="color:yellow">POST</span> /api/v1/book/return

Authorization : Bearer JWT

Request

```json
{
    "book_id":"67bc8101112ff4d3c99effbc"
}
```

Response

```json
{
    "message": "Returned book successfully"
}
```