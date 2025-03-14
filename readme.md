
# REST API for Social Media App

- The user can create a post
- The Post can contains comments from users
- User authentication handled by JWT

## Tech Stack
- Golang
- Gin Framework
- PostgreSQL


## Installation

Clone this repo

```bash
  git clone [this repo]
  cd gotweet
```

Run the backend
```bash
  go run cmd/main.go
```

## API Endpoint Request Example

### **Signup User**

Request Body
```json
{
  "username": "Mike",
  "email": "mike@email.com",
  "password": "mike123"
}
```

Response Body
```json
{
  "message": "user successfully created"
}
```