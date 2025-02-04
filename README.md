# Anne Hub

Routing and processing server for anne wear and anne companion.

## Setup the PostgreSQL DB and Migrations

Make sure you have PostgreSQL 14 installed and running.

### Install PostgreSQL

```sh
brew install postgres@14
```

### Start PostgreSQL Service

```sh
brew services start postgresql@14
```

### Install Migration Tool

Get the CLI tool for migrations (with brew on Linux or macOS, otherwise you can get it from the releases of the [official repository](https://github.com/golang-migrate/migrate)):

```sh
brew install golang-migrate
```

### Create the Database

Replace `<username>` with your PostgreSQL username.

```sh
psql -U <username> -tc "SELECT 1 FROM pg_database WHERE datname = 'anne_hub';" | grep -q 1 || psql -U <username> -c "CREATE DATABASE anne_hub;"
```

### Create a Migration

```sh
migrate create -ext sql -dir db/migrations -seq migration_name
```

### Apply Migrations

Set the `ANNE_HUB_DB` environment variable with your database connection string and run:

```sh
migrate -database $ANNE_HUB_DB -path db/migrations up
```

## Environment Variables

Create a `.env` file in the root directory with the following content:

```env
GROQ_API_KEY=your_groq_api_key
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=your_db_username
DB_NAME=anne_hub
DB_PASSWORD=your_db_password
DB_SSLMODE=disable
```

Ensure you replace the placeholders with your actual configuration.

## Quickstart with Docker

For building:

```sh
sudo docker build -t anne-hub .
```

To run:

```sh
sudo docker run --env-file .env -p 1323:1323 anne-hub
```

## Quickstart with Go CLI

```sh
go mod tidy
```

```sh
go build -o bin/anne-hub
```

```sh
./bin/anne-hub
```

## API Documentation

### General Routes

#### GET `/ok`

- **Description**: Health check endpoint.
- **Response**:
  - Status: `200 OK`
  - Body:

    ```json
    {
      "message": "OK"
    }
    ```

#### GET `/gh-actions-test`

- **Description**: Endpoint for testing GitHub Actions.
- **Response**:
  - Status: `200 OK`
  - Body:

    ```json
    {
      "message": "GH Actions Test"
    }
    ```

#### GET `/uuid`

- **Description**: Generates a new UUID.
- **Response**:
  - Status: `200 OK`
  - Body:

    ```json
    {
      "uuid": "generated-uuid-string"
    }
    ```

### User Routes

#### GET 

users



- **Description**: Retrieve all users along with their interests and tasks.
- **Response**:
  - Status: `200 OK`
  - Body: Array of user data with interests and tasks.

#### GET `/users/:id`

- **Description**: Retrieve a user by ID.
- **Parameters**:
  - 

id

 (path): UUID of the user.
- **Response**:
  - Status: `200 OK`
  - Body: User details.

#### POST 

users



- **Description**: Create a new user.
- **Request Body**:

  ```json
  {
    "username": "string",
    "first_name": "string",
    "last_name": "string",
    "email": "string",
    "password": "string",
    "age": integer,
    "country": "string",
    "city": "string"
  }
  ```

- **Response**:
  - Status: `201 Created`
  - Body: Created user details.

#### PUT `/users/:id`

- **Description**: Update an existing user.
- **Parameters**:
  - 

id

 (path): UUID of the user.
- **Request Body**:

  ```json
  {
    "username": "string",
    "first_name": "string",
    "last_name": "string",
    "email": "string",
    "password": "string",
    "age": integer,
    "country": "string",
    "city": "string"
  }
  ```

- **Response**:
  - Status: `200 OK`
  - Body: Updated user details.

#### DELETE `/users/:id`

- **Description**: Delete a user by ID.
- **Parameters**:
  - 

id

 (path): UUID of the user.
- **Response**:
  - Status: `200 OK`
  - Body:

    ```json
    {
      "message": "User deleted successfully."
    }
    ```

### Task Routes

#### GET `/tasks`

- **Description**: Retrieve all tasks.
- **Response**:
  - Status: `200 OK`
  - Body: Array of tasks.

#### GET `/tasks/:id`

- **Description**: Retrieve all tasks for a specific user.
- **Parameters**:
  - 

id

 (path): UUID of the user.
- **Response**:
  - Status: `200 OK`
  - Body: Array of tasks for the user.

#### POST `/tasks`

- **Description**: Create a new task.
- **Request Body**:

  ```json
  {
    "user_id": "uuid",
    "title": "string",
    "description": "string",
    "due_date": "YYYY-MM-DDTHH:MM:SSZ",  // ISO8601 format
    "completed": boolean
  }
  ```

- **Response**:
  - Status: `201 Created`
  - Body: Created task details.

#### PUT `/tasks/:id`

- **Description**: Update an existing task.
- **Parameters**:
  - 

id

 (path): ID of the task.
- **Request Body**:

  ```json
  {
    "title": "string",
    "description": "string",
    "due_date": "YYYY-MM-DDTHH:MM:SSZ",  // ISO8601 format
    "completed": boolean
  }
  ```

- **Response**:
  - Status: `200 OK`
  - Body: Updated task details.

#### DELETE `/tasks/:id`

- **Description**: Delete a task by ID.
- **Parameters**:
  - 

id

 (path): ID of the task.
- **Response**:
  - Status: `200 OK`
  - Body:

    ```json
    {
      "message": "Task deleted successfully."
    }
    ```

### Interest Routes

#### GET `/interests`

- **Description**: Retrieve all interests.
- **Response**:
  - Status: `200 OK`
  - Body: Array of interests.

#### GET `/interests/:id`

- **Description**: Retrieve an interest by ID.
- **Parameters**:
  - 

id

 (path): ID of the interest.
- **Response**:
  - Status: `200 OK`
  - Body: Interest details.

#### POST `/interests`

- **Description**: Create a new interest.
- **Request Body**:

  ```json
  {
    "user_id": "uuid",
    "name": "string",
    "description": "string",
    "level": integer,
    "level_accuracy": "string"
  }
  ```

- **Response**:
  - Status: `201 Created`
  - Body: Created interest details.

#### PUT `/interests/:id`

- **Description**: Update an existing interest.
- **Parameters**:
  - 

id

 (path): ID of the interest.
- **Request Body**:

  ```json
  {
    "name": "string",
    "description": "string",
    "level": integer,
    "level_accuracy": "string"
  }
  ```

- **Response**:
  - Status: `200 OK`
  - Body: Updated interest details.

#### DELETE `/interests/:id`

- **Description**: Delete an interest by ID.
- **Parameters**:
  - 

id

 (path): ID of the interest.
- **Response**:
  - Status: `200 OK`
  - Body:

    ```json
    {
      "message": "Interest deleted successfully."
    }
    ```

### WebSocket Routes

#### WebSocket `/ws/conversation`

- **Description**: WebSocket endpoint for handling real-time conversations.

- **Headers**:
  - `X-User-ID`: UUID of the user.
  - `X-Device-ID`: Device ID.
  - `X-Language`: Language code (e.g., `en`, `de`).

- **Protocol**:
  1. **Connection**: Establish a WebSocket connection to `/ws/conversation`.
  2. **Send Headers**: Clients must first send a JSON message containing the required headers.

     ```json
     {
       "X-User-ID": "uuid",
       "X-Device-ID": "device_id",
       "X-Language": "language_code"
     }
     ```

  3. **Send Audio Data**: After headers are accepted, clients can send binary messages with audio data.
  4. **End of Stream**: To indicate the end of the audio stream, send the text message `"EOS"`.

- **Response**:
  - The server will process the audio data and send back responses as text messages, including any assistant responses and actions.

## Additional Notes

- Current setup doesn't contain any authentication nor any middleware, due to its `development` status.

## License

This project is licensed under the MIT License. See the 

LICENSE

 file for details.
