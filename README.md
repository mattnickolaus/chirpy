# Chirpy

A backend HTTP server in Go built without use of a framework built as a part of [Boot.dev's](www.boot.dev) backend course.

## Why You Should Care

This project, developed as part of a comprehensive backend course, demonstrates proficiency in several key areas of modern web development:

*   **Go Web Servers:** Deep understanding of how to build production-style HTTP servers in Go from scratch, without relying on heavy frameworks.

*   **RESTful APIs:** Implemented robust RESTful APIs using JSON, proper HTTP headers, and status codes for effective client-server communication.

*   **Database Management:** Utilized type-safe SQL with PostgreSQL to efficiently store and retrieve data, showcasing practical database integration skills.

*   **Secure Authentication & Authorization:** Developed a secure system using well-tested cryptography libraries for user authentication and authorization, including token-based mechanisms.

*   **Webhooks & API Keys:** Gained hands-on experience with webhooks for inter-service communication and implemented secure API key management.

*   **API Documentation:** Maintained clear and concise API documentation, a crucial skill for collaborative development and usability.


This project showcases the ability to architect, build, and secure a complete backend application, highlighting a strong foundation in modern backend engineering principles.

## Prerequisites

Before you begin, ensure you have the following installed:

*   [Go](https://golang.org/doc/install)
*   [PostgreSQL](https://www.postgresql.org/download/)

## Installation

To install simply clone this repo: 

```bash
git clone github.com/mattnickolaus/chirpy
```

## Configuration

Before running chirpy, you need to setup a Postgres database and create a `.env` file.

### Postgres Schema

Enter the `psql` shell:
- Mac: `psql postgres`
- Linux: `sudo -u postgres psql`

Create a new database called gator: 

``` sql
CREATE DATABASE gator;
```

Connect to the database:

``` sql
\c gator
```

Create a new user and grant it privileges to the gator database. You can do this with the following commands in `psql`:

```sql
CREATE USER gator_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE gator TO gator_user;
```

This username and password will then be used in the PostgreSQL connection string.

```
"postgres://username:password@localhost:5432/gator"
```

### `.env` File

Next you need to create a `.env` file at the root of the project directory. 

``` bash
touch path_to_project/.env 
```
This file should contain the following: 

```
DB_URL="postgres://username:password@localhost:5432/chirpy"
PLATFORM="dev" # change to current environment
SECRET="generated-random-string"
POLKA_KEY="f271c81ff7084ee5b99a5091b42d486e" # webhook api key
```

Replace the `DB_URL` with your PostgreSQL connection string. For the `SECRET` variable you can generate a long random string with the below command and replace the current string contents.

```
openssl rand -base64 64
```

### Set Up Migration

Once you have defined your environment variables, you can run the following command from the root of the project directory to migrate to the proper Postgres schema:

```bash
cd sql/schema && goose postgres DB_URL up && cd ../..
```
Replace the DB_URL again with your PostgreSQL connection string.

**From there you are ready to run the Chirpy server!!**

## API Documentation

### `GET /api/healthz`

A health check endpoint.

**Responses:**

- `200 OK`: with a `text/plain` body `OK`.

---

### `POST /api/users`

Creates a new user.

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Responses:**

- `201 Created`: with the created user object (without password).
  ```json
  {
    "id": "...",
    "email": "user@example.com",
    "is_chirpy_red": false
  }
  ```
- `400 Bad Request`: on malformed JSON or validation error.
- `500 Internal Server Error`: on other errors.

---

### `PUT /api/users`

Updates an existing user's email or password. Requires authentication.

**Headers:**

- `Authorization: Bearer <token>`

**Request Body:**

```json
{
  "email": "new.email@example.com",
  "password": "newpassword123"
}
```

**Responses:**

- `200 OK`: with the updated user object.
  ```json
  {
    "id": "...",
    "email": "new.email@example.com"
  }
  ```
- `401 Unauthorized`: if the token is invalid or not provided.
- `400 Bad Request`: on malformed JSON or validation error.
- `500 Internal Server Error`: on other errors.

---

### `POST /api/login`

Logs in a user.

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "password123",
}
```

**Responses:**

- `200 OK`: with user details and tokens.
  ```json
  {
    "id": "...",
    "email": "user@example.com",
    "is_chirpy_red": false,
    "token": "...",
    "refresh_token": "..."
  }
  ```
- `401 Unauthorized`: if credentials are incorrect.
- `400 Bad Request`: on malformed JSON.
- `500 Internal Server Error`: on other errors.

---

### `POST /api/polka/webhooks`

A webhook endpoint for Polka to upgrade a user to Chirpy Red.

**Headers:**

- `Authorization: ApiKey <polka_api_key>`

**Request Body:**

```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": 123
  }
}
```

**Responses:**

- `200 OK`
- `401 Unauthorized`: if the API key is invalid or not for the correct event type.
- `400 Bad Request`: on malformed JSON.
- `500 Internal Server Error`: on other errors.

---

### `POST /api/chirps`

Creates a new chirp. Requires authentication.

**Headers:**

- `Authorization: Bearer <token>`

**Request Body:**

```json
{
  "body": "This is a new chirp!"
}
```

**Responses:**

- `201 Created`: with the created chirp object.
  ```json
  {
    "id": "...",
    "body": "This is a new chirp!",
    "user_id": "..."
  }
  ```
- `400 Bad Request`: if the chirp is too long or empty.
- `401 Unauthorized`: if the token is invalid or not provided.
- `500 Internal Server Error`: on other errors.

---

### `GET /api/chirps`

Gets all chirps.

**Query Parameters:**

- `author_id`: (optional) filter chirps by author ID.
- `sort`: (optional) `asc` or `desc`. Defaults to `asc`.

**Responses:**

- `200 OK`: with an array of chirp objects.
  ```json
  [
    {
      "id": "...",
      "body": "...",
      "user_id": "..."
    }
  ]
  ```
- `500 Internal Server Error`: on other errors.

---

### `GET /api/chirps/{chirpID}`

Gets a single chirp by its ID.

**Path Parameters:**

- `chirpID`: The ID of the chirp to retrieve.

**Responses:**

- `200 OK`: with the chirp object.
- `404 Not Found`: if the chirp doesn't exist.
- `500 Internal Server Error`: on other errors.

---

### `DELETE /api/chirps/{chirpID}`

Deletes a chirp. Requires authentication and the authenticated user must be the author of the chirp.

**Headers:**

- `Authorization: Bearer <token>`

**Path Parameters:**

- `chirpID`: The ID of the chirp to delete.

**Responses:**

- `200 OK`
- `403 Forbidden`: if the user is not the author of the chirp.
- `401 Unauthorized`: if the token is invalid or not provided.
- `404 Not Found`: if the chirp doesn't exist.
- `500 Internal Server Error`: on other errors.

---

### `POST /api/refresh`

Refreshes an access token using a refresh token.

**Headers:**

- `Authorization: Bearer <refresh_token>`

**Responses:**

- `200 OK`: with a new access token.
  ```json
  {
    "token": "..."
  }
  ```
- `401 Unauthorized`: if the refresh token is invalid or expired.
- `500 Internal Server Error`: on other errors.

---

### `POST /api/revoke`

Revokes a refresh token.

**Headers:**

- `Authorization: Bearer <refresh_token>`

**Responses:**

- `200 OK`
- `401 Unauthorized`: if the refresh token is invalid.
- `500 Internal Server Error`: on other errors.

---

## Admin Endpoints

These are internal admin endpoints.

- `POST /admin/reset`: Resets the API hit counter.
- `GET /admin/metrics`: Returns the number of API hits in an HTML format.
