# Quickstart — Send Your First Notification

## Prerequisites

- Server is running: `make run`
- Keycloak is running and realm `notification-center` is configured
- PostgreSQL is reachable (check `config.yaml`)

---

## Step 1 — Get a Keycloak Access Token

```bash
curl -s -X POST \
  http://<keycloak-host>/realms/notification-center/protocol/openid-connect/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "client_id=ithq-notification-center" \
  -d "client_secret=<your-client-secret>" \
  -d "username=<your-username>" \
  -d "password=<your-password>" \
  -d "grant_type=password"
```

**Response:**
```json
{
  "access_token": "eyJhbGci...",
  "expires_in": 300
}
```

Save the `access_token` — you will use it as `Bearer <token>` in all JWT-protected requests.

---

## Step 2 — Sync Your User (first login)

This creates your user record in the local database from Keycloak claims.

```bash
curl -s http://localhost:8080/auth/me \
  -H "Authorization: Bearer <access_token>"
```

**Response:**
```json
{
  "id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
  "email": "you@example.com",
  "username": "admin"
}
```

> This must be called at least once before doing anything else. Your user record must exist in the DB.

---

## Step 3 — Create a Project

```bash
curl -s -X POST http://localhost:8080/projects \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My App",
    "description": "My first project"
  }'
```

**Response:**
```json
{
  "project": {
    "id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
    "name": "My App",
    "slug": "my-app"
  },
  "api_key": "nc_live_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
}
```

Save both:
- `project.id` → **PROJECT_ID**
- `api_key` → **API_KEY** (shown only once, save it now)

---

## Step 4 — Find the Recipient User ID

You need the database UUID of the user you want to notify.  
Your own user ID was returned in Step 2. You can also call:

```bash
curl -s http://localhost:8080/users/me \
  -H "Authorization: Bearer <access_token>"
```

Save the `id` field → **USER_ID**

---

## Step 5 — Send a Notification

### Option A — Using JWT (as the project owner)

```bash
curl -s -X POST http://localhost:8080/projects/<PROJECT_ID>/notifications \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Hello!",
    "body": "Your first notification from Notification Center.",
    "priority": "normal",
    "recipients": [
      {
        "user_id": "<USER_ID>",
        "channels": ["in_app"]
      }
    ]
  }'
```

### Option B — Using API Key (from your backend/service)

```bash
curl -s -X POST http://localhost:8080/api/v1/notifications \
  -H "X-API-Key: <API_KEY>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Hello!",
    "body": "Your first notification from Notification Center.",
    "priority": "normal",
    "recipients": [
      {
        "user_id": "<USER_ID>",
        "channels": ["in_app"]
      }
    ]
  }'
```

**Response:**
```json
{
  "id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
  "title": "Hello!",
  "body": "Your first notification from Notification Center.",
  "priority": "normal",
  "recipients": [
    {
      "id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
      "channel": "in_app",
      "status": "pending"
    }
  ]
}
```

---

## Step 6 — Check the Inbox

```bash
curl -s http://localhost:8080/inbox \
  -H "Authorization: Bearer <access_token>"
```

**Response:**
```json
{
  "data": [...],
  "total": 1,
  "limit": 20,
  "offset": 0
}
```

---

## Step 7 — Mark as Read

```bash
curl -s -X POST http://localhost:8080/inbox/<NOTIFICATION_ID>/read \
  -H "Authorization: Bearer <access_token>"
```

---

## Available Channels

| Channel | Value | Requires |
|---|---|---|
| In-app inbox | `in_app` | nothing extra |
| Email | `email` | user has email in DB + worker running |
| Push | `push` | device token registered + worker running |
| SMS | `sms` | worker running + SMS provider configured |

For `in_app` everything works without the worker. For other channels start the worker with `make run-worker`.

---

## Summary

```
1. Get token from Keycloak
2. GET  /auth/me                          ← sync user to DB (first time only)
3. POST /projects                         ← create project, get PROJECT_ID + API_KEY
4. GET  /users/me                         ← get your USER_ID
5. POST /api/v1/notifications             ← send notification (use API_KEY)
6. GET  /inbox                            ← read notifications (use JWT)
7. POST /inbox/:id/read                   ← mark as read
```
