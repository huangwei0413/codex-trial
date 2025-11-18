# Student Service API

This document describes the RESTful endpoints exposed by the Student microservice. The API follows JSON over HTTP with standard status codes and a consistent response envelope.

- **Base URL**: `http://localhost:8080`
- **API version prefix**: `/api/v1`
- **Content-Type**: `application/json`

## Response Envelope

Successful responses wrap payloads inside `"data"`. Errors return an `"error"` string and, when relevant, `"message"` context.

```json
// Success
{
  "data": { ... }
}

// Error
{
  "error": "Student not found"
}
```

## Domain Model

| Field       | Type      | Notes                              |
|-------------|-----------|------------------------------------|
| `id`        | UUID      | Auto-generated primary key         |
| `name`      | string    | Required                           |
| `email`     | string    | Required, unique                   |
| `age`       | integer   | Optional                           |
| `major`     | string    | Optional                           |
| `gpa`       | float     | Optional                           |
| `created_at`| timestamp | Managed by the service             |
| `updated_at`| timestamp | Managed by the service             |

## Endpoints

### Health Check

- **GET** `/health`
- **200 OK**

```json
{
  "status": "healthy"
}
```

### List Students

- **GET** `/api/v1/students`
- **200 OK**

```json
{
  "data": [
    {
      "id": "22f26a59-8b4f-4dfe-b1be-489a78ed6d6f",
      "name": "Jane Doe",
      "email": "jane@example.com",
      "age": 21,
      "major": "Computer Science",
      "gpa": 3.8,
      "created_at": "2024-01-23T12:00:00Z",
      "updated_at": "2024-01-23T12:00:00Z"
    }
  ]
}
```

### Get Student By ID

- **GET** `/api/v1/students/{id}`
- **200 OK**

```json
{
  "data": {
    "id": "22f26a59-8b4f-4dfe-b1be-489a78ed6d6f",
    "name": "Jane Doe",
    "email": "jane@example.com",
    "age": 21,
    "major": "Computer Science",
    "gpa": 3.8,
    "created_at": "2024-01-23T12:00:00Z",
    "updated_at": "2024-01-23T12:00:00Z"
  }
}
```

- **404 Not Found** when the ID does not exist:

```json
{
  "error": "Student not found"
}
```

### Create Student

- **POST** `/api/v1/students`
- **Request body**

```json
{
  "name": "Jane Doe",
  "email": "jane@example.com",
  "age": 21,
  "major": "Computer Science",
  "gpa": 3.8
}
```

- **201 Created**

```json
{
  "data": {
    "id": "22f26a59-8b4f-4dfe-b1be-489a78ed6d6f",
    "name": "Jane Doe",
    "email": "jane@example.com",
    "age": 21,
    "major": "Computer Science",
    "gpa": 3.8,
    "created_at": "2024-01-23T12:00:00Z",
    "updated_at": "2024-01-23T12:00:00Z"
  }
}
```

- **400 Bad Request** when the payload is invalid:

```json
{
  "error": "Key: 'Student.Name' Error:Field validation for 'Name' failed on the 'required' tag"
}
```

- **409 Conflict** is not currently returned explicitly; duplicate emails will raise a `500` with `"Failed to create student"`.

### Update Student

- **PUT** `/api/v1/students/{id}`
- **Request body**: any subset of mutable fields (`name`, `email`, `age`, `major`, `gpa`)

```json
{
  "major": "Mathematics",
  "gpa": 3.9
}
```

- **200 OK** returns the updated record

```json
{
  "data": {
    "id": "22f26a59-8b4f-4dfe-b1be-489a78ed6d6f",
    "name": "Jane Doe",
    "email": "jane@example.com",
    "age": 21,
    "major": "Mathematics",
    "gpa": 3.9,
    "created_at": "2024-01-23T12:00:00Z",
    "updated_at": "2024-02-01T09:15:00Z"
  }
}
```

- **404 Not Found** if the student does not exist
- **400 Bad Request** for invalid JSON payloads

### Delete Student

- **DELETE** `/api/v1/students/{id}`
- **200 OK**

```json
{
  "message": "Student deleted successfully"
}
```

- **404 Not Found** if the student does not exist (handled implicitly when the ORM cannot match any rows)

## Error Handling

The service currently returns concise error strings. Common status codes include:

| Status | Scenario                          |
|--------|-----------------------------------|
| 200    | Successful read or delete         |
| 201    | Student created                   |
| 400    | Invalid JSON payload              |
| 404    | Student not found                 |
| 500    | Database or server failure        |

## Versioning

All resource endpoints live under `/api/v1/`. Breaking changes should be shipped behind a new version prefix (`/api/v2`). The `/health` endpoint is left at the root for platform probes.
