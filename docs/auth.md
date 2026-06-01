# ESB 服务认证说明

ESB 自己的接口认证分为两类：

1. 用户登录 token：通过账号密码登录获取，token 存在 Redis，默认 2 小时有效。
2. 应用 APP TOKEN：在 `conf/app.ini` 中配置，适合 agent、插件、系统对系统调用，没有过期时间。

当前对接 agent 或第三方调用方时，建议使用第二种：应用 APP TOKEN。

## 认证配置

```ini
[service_auth]
Username = admin
Password = admin123
ApiTokens = nXGWDADLnVeVz3jBWRfHeMUh6CB4KADRTTd8Jb5LUuL:gaiastandard
TokenTTLSeconds = 7200
```

字段说明：

| 字段 | 说明 |
| --- | --- |
| `Username` | 用户登录接口的账号。 |
| `Password` | 用户登录接口的密码。 |
| `ApiTokens` | 分配给应用或 agent 的 APP TOKEN。支持多个 token，用英文逗号分隔。 |
| `TokenTTLSeconds` | 用户登录 token 的 Redis 有效期，默认 7200 秒。对 APP TOKEN 不生效。 |

## APP TOKEN 分配方式

`ApiTokens` 支持按 OpenAPI 应用授权。这里的应用名来自配置段 `[openapi.{appName}]`，例如 `[openapi.gaiastandard]` 的应用名就是 `gaiastandard`。

```ini
[service_auth]
ApiTokens = token-a:gaiastandard,token-b:robot|gaiastandard,gaiastandard=token-c,legacy-token
```

| 写法 | 含义 |
| --- | --- |
| `token-a:gaiastandard` | `token-a` 只能调用匹配到 `gaiastandard` 的第三方 URL。 |
| `token-b:robot|gaiastandard` | `token-b` 可以调用 `robot` 和 `gaiastandard` 两个应用。 |
| `gaiastandard=token-c` | `token-c` 只能调用 `gaiastandard`。 |
| `legacy-token` | 兼容旧写法，可以调用所有已配置的 OpenAPI 应用。 |

建议新接入的应用都使用带应用名的写法，例如：

```ini
ApiTokens = app-token-for-gaia:gaiastandard,app-token-for-robot:robot
```

这样 Gaia 的 token 不能调用 robot，robot 的 token 也不能调用 Gaia。

## agent 插件里怎么填写

在 agent 或插件配置页面中，授权方式选择 `API Key`。

添加位置选择：`Header`

Key 填：

```text
X-API-Token
```

Value 填：

```text
<分配给这个应用的APP_TOKEN>
```

注意：

- 这里填的是 ESB 的 APP TOKEN。
- 不需要加 `Bearer`。
- 不要把 Gaia 第三方返回的 token 填到这里。
- 如果请求体是 JSON，普通 Header 里再加一项 `Content-Type: application/json`。

## 用户登录 token

用户登录适合后台管理、临时调试等场景。

```http
POST /api/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin123"
}
```

成功返回：

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

使用用户登录 token 调用 ESB：

```http
GET /api/esb/v1/https://gaiaopenapi-s.copm.com.cn/atd-webapi/api/v1/attendance-type/leave/list
Authorization: Bearer <user-token>
```

退出登录：

```http
POST /api/auth/logout
Authorization: Bearer <user-token>
```

退出只会删除 Redis 里的用户登录 token，不会影响 `conf/app.ini` 中的 APP TOKEN。

## ESB token 和第三方 token 的区别

这两个 token 不要混在一起：

| token | 谁使用 | 放在哪里 | 是否 Bearer |
| --- | --- | --- | --- |
| ESB APP TOKEN | agent、插件、业务系统调用 ESB | `X-API-Token` | 不需要 |
| ESB 用户登录 token | 登录用户调用 ESB | `Authorization` | 需要 `Bearer` |
| Gaia 第三方 token | ESB 转发请求到 Gaia 时自动使用 | ESB 内部自动添加 `Authorization`，token 缓存在 Redis | Gaia 要求 `Bearer` |

所以配置 agent 时，只需要配置 `X-API-Token: <APP_TOKEN>`。

Gaia 第三方 token 的 Redis key 是：

```text
esb:openapi:token:gaiastandard
```

当前 `[openapi.gaiastandard]` 配置的 `TokenTTLSeconds = 7200`，也就是 2 小时。Redis key 过期后，ESB 才会重新调用 Gaia 鉴权接口获取新 token。
