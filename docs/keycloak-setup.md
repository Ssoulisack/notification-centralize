# Keycloak Setup for Notification Center

Config values used: realm `notification-center`, client `ithq-notification-center`

---

## 1. Create a Realm

- [ ] Go to **Keycloak Admin Console**
- [ ] Click **Create Realm**
- [ ] Set **Realm name**: `notification-center`
- [ ] Click **Create**

---

## 2. Create a Client

- [ ] Go to **Clients** → **Create client**
- [ ] **Client type**: `OpenID Connect`
- [ ] **Client ID**: `ithq-notification-center`
- [ ] Click **Next**

### Capability Config
- [ ] **Client authentication**: `ON` (makes it a confidential client)
- [ ] **Authorization**: `OFF`
- [ ] **Authentication flow**: check `Standard flow` and `Direct access grants`
- [ ] Click **Next**

### Login Settings
- [ ] **Valid redirect URIs**: `http://localhost:8080/*`
- [ ] **Web origins**: `http://localhost:8080`
- [ ] Click **Save**

---

## 3. Copy the Client Secret

- [ ] Go to **Clients** → `ithq-notification-center` → **Credentials** tab
- [ ] Copy **Client secret**
- [ ] Paste it into `config.yaml`:

```yaml
keycloak:
  base_url: http://<your-keycloak-host>
  realm: notification-center
  client_id: ithq-notification-center
  client_secret: <paste-secret-here>
```

---

## 4. Configure Token Claims

The app needs `email`, `given_name`, `family_name`, `preferred_username` in the JWT.

- [ ] Go to **Clients** → `ithq-notification-center` → **Client scopes** tab
- [ ] Click `ithq-notification-center-dedicated`
- [ ] Click **Add mapper** → **By configuration**
- [ ] Verify these mappers exist (they are usually added by default via `profile` and `email` scopes):

| Claim | Mapper type | Token claim name |
|---|---|---|
| `email` | User Property | `email` |
| `preferred_username` | User Property | `preferred_username` |
| `given_name` | User Property | `given_name` |
| `family_name` | User Property | `family_name` |

- [ ] Go to **Clients** → `ithq-notification-center` → **Client scopes** tab (top)
- [ ] Confirm `profile` and `email` are in the **Assigned default client scopes** list

---

## 5. Create a User

- [ ] Go to **Users** → **Create new user**
- [ ] Fill in:
  - **Username**: e.g. `admin`
  - **Email**: your email
  - **First name** / **Last name**
  - **Email verified**: `ON`
- [ ] Click **Create**
- [ ] Go to **Credentials** tab → **Set password**
- [ ] Enter a password, set **Temporary**: `OFF`
- [ ] Click **Save**

---

## 6. Test — Get a Token

```bash
curl -s -X POST \
  http://<your-keycloak-host>/realms/notification-center/protocol/openid-connect/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "client_id=ithq-notification-center" \
  -d "client_secret=<your-client-secret>" \
  -d "username=<your-username>" \
  -d "password=<your-password>" \
  -d "grant_type=password" | jq '.access_token'
```

Copy the `access_token` and use it as `Bearer <token>` in the `Authorization` header.

---

## 7. Verify JWT Contains Required Claims

```bash
# Decode the token (replace TOKEN with your access_token)
echo "TOKEN" | cut -d. -f2 | base64 -d 2>/dev/null | jq '{sub, email, preferred_username, given_name, family_name}'
```

Expected output:
```json
{
  "sub": "some-uuid",
  "email": "you@example.com",
  "preferred_username": "admin",
  "given_name": "John",
  "family_name": "Doe"
}
```

The `sub` field becomes `keycloak_id` in the local users table.
