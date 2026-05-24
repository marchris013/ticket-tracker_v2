# Ticket Tracker API v2
A simple REST API for tracking tickets, built using Go and PostgreSQL.

## How to Run
Make sure you have Go and PostgreSQL installed.

Create a database named `ticket_tracker` in PostgreSQL, then create the tickets table:

```sql
CREATE TABLE tickets (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    status TEXT NOT NULL DEFAULT 'open',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

Update the connection string in `main.go` with your PostgreSQL password, then:
```
go run main.go
```
The server starts on port 8080.

## Endpoints
**Create a ticket**
```
POST /tickets
```
Body:
```json
{
  "title": "Fix login bug",
  "description": "Users cant log in on mobile"
}
```

**Get all tickets**
```
GET /tickets
```

**Get one ticket**
```
GET /tickets/1
```

## Design Decisions
- PostgreSQL replaces in-memory storage — data persists across server restarts
- ID, status, and created_at are set by the database automatically, not by the server
- Database connection is verified with Ping on startup — server won't start if database is unreachable
- Concurrency is handled by PostgreSQL instead of a mutex
