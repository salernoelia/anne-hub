# Anne Hub

Routing and processing server for Anne wearable and companion app.

[Project Info](https://eliasalerno.ch/work/anne)

## Quick Setup

### Database
```sh
brew install postgres@14 golang-migrate
brew services start postgresql@14

# Create database
psql -U <username> -tc "SELECT 1 FROM pg_database WHERE datname = 'anne_hub';" | grep -q 1 || psql -U <username> -c "CREATE DATABASE anne_hub;"

# Run migrations
migrate -database $ANNE_HUB_DB -path db/migrations up
```

### Environment
Create `.env`:
```env
GROQ_API_KEY=your_groq_api_key
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=your_db_username
DB_NAME=anne_hub
DB_PASSWORD=your_db_password
DB_SSLMODE=disable
```

### Run
**Docker:**
```sh
docker build -t anne-hub .
docker run --env-file .env -p 1323:1323 anne-hub
```

**Go:**
```sh
go mod tidy
go build -o bin/anne-hub
./bin/anne-hub
```

## Key Endpoints

- `GET /ok` - Health check
- `GET /users` - All users with interests/tasks
- `POST /users` - Create user
- `GET /tasks/:user_id` - User's tasks
- `POST /tasks` - Create task
- `WebSocket /ws/conversation` - Real-time chat with audio

### WebSocket Usage
1. Connect to `/ws/conversation`
2. Send headers as JSON:
   ```json
   {
     "X-User-ID": "uuid",
     "X-Device-ID": "device_id", 
     "X-Language": "en"
   }
   ```
3. Send binary audio data
4. Send `"EOS"` to end stream

## License

MIT
