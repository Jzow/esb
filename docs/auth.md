# ESB API Authentication

The ESB service supports two authentication modes for its own public APIs.

1. User login token: call the login API with username and password. The generated token is stored in Redis and expires after 2 hours by default.
2. Application API token: configure static tokens in `conf/app.ini`. These tokens do not expire and are intended for service-to-service calls.

## Configuration

```ini
[service_auth]
Username = admin
Password = admin123
ApiTokens = esb-app-token-change-me
TokenTTLSeconds = 7200
```

`ApiTokens` supports multiple tokens separated by commas.

Gaia configuration must keep `BaseUrl` as the Gaia API root, not the OAuth endpoint:

```ini
[gaia-openapi]
BaseUrl = https://gaiaopenapi-s.copm.com.cn
AuthURL = https://gaiaopenapi-s.copm.com.cn/identity/api/v1/oauth
TokenPrefix = Bearer
```

If `BaseUrl` is set to the OAuth endpoint, business paths such as `LeaveTypesPath` will be appended after `/oauth` and Gaia will reject the request.

## Login

```http
POST /api/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin123"
}
```

Success response:

```json
{
  "code": 200,
  "msg": "ok",
  "data": {
    "token": "<user-token>",
    "token_type": "Bearer",
    "expires_in": 7200,
    "username": "admin"
  }
}
```

## Call Gaia Proxy APIs

Use a user login token:

```http
GET /api/gaia/v1/leave/types
Authorization: Bearer <user-token>
```

Use an application API token:

```http
GET /api/gaia/v1/leave/types
Authorization: Bearer esb-app-token-change-me
```

You can also pass the application token with:

```http
X-API-Token: esb-app-token-change-me
```

## Logout

```http
POST /api/auth/logout
Authorization: Bearer <user-token>
```

Logout removes user login tokens from Redis. Application API tokens are static config values and are not removed by logout.
