# ESB OpenAPI 通用代理接口说明

ESB 当前对外提供一组通用代理接口，不再为每个第三方接口单独写路由，例如不再需要配置 `SubmitLeavePath`、`LeaveTypesPath` 这类固定地址。

调用方把第三方完整 URL 放在 `/api/esb/v1/` 后面，ESB 会根据 URL 的域名匹配 `conf/app.ini` 中的 `[openapi.*]` 配置，然后完成认证、转发和响应包装。

## 通用接口

```http
GET    /api/esb/v1/{第三方完整URL}
POST   /api/esb/v1/{第三方完整URL}
PUT    /api/esb/v1/{第三方完整URL}
DELETE /api/esb/v1/{第三方完整URL}
```

示例：

```http
GET /api/esb/v1/https://gaiaopenapi-s.copm.com.cn/atd-webapi/api/v1/workflow/apply/list
X-API-Token: <APP_TOKEN>
```

`/api/esb/v1/` 后面必须是完整第三方 URL，包含 `https://`、域名和路径。

项目也提供 `/swagger/index.html` 页面查看接口文档。由于 Swagger 对“路径参数里再放完整 URL”这种场景支持有限，页面里会显示成 `/api/esb/v1/{url}`。实际联调时建议优先用 Apifox、Postman 或 curl 按完整 URL 调用。

## 认证方式

每次调用 ESB 代理接口都需要传 ESB 自己的 APP TOKEN。

推荐写法：

```http
X-API-Token: <APP_TOKEN>
```

APP TOKEN 是普通应用令牌，不是 JWT，不要加 `Bearer`。

兼容写法：

```http
Authorization: <APP_TOKEN>
```

不推荐在 APP TOKEN 前加 `Bearer`，避免和用户登录 token、第三方 Gaia token 混淆。

## agent 插件配置

如果在 agent 或插件管理页面配置本服务：

插件 URL 填 ESB 代理地址，例如：

```text
http://localhost:8000/api/esb/v1/https://gaiaopenapi-s.copm.com.cn/atd-webapi/api/v1/workflow/apply/list
```

授权方式选择：

```text
API Key
```

添加位置选择：

```text
Header
```

Key 填：

```text
X-API-Token
```

Value 填：

```text
<分配给这个agent或应用的APP_TOKEN>
```

如果接口是 POST、PUT，并且请求体是 JSON，再在 Header 列表里加：

```text
Content-Type: application/json
```

## OpenAPI 应用配置

每个第三方 OpenAPI 应用使用一个独立配置段：

```ini
[openapi.gaiastandard]
BaseUrl = https://gaiaopenapi-s.copm.com.cn
AuthURL = https://gaiaopenapi-s.copm.com.cn/identity/api/v1/oauth
GrantType = client_credentials
ClientSecret = <client_secret>
CorpID = zhwytest01
TokenTTLSeconds = 7200
TokenPrefix = Bearer
FixedHeaders = tenant={CorpID}
UserAgent = Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36
```

字段说明：

| 字段 | 说明 |
| --- | --- |
| `BaseUrl` | 第三方 OpenAPI 根地址。ESB 根据请求 URL 的协议和域名匹配这个配置。 |
| `AuthURL` | 第三方 token 接口。不需要第三方认证的应用可以不配。 |
| `GrantType` | 第三方鉴权参数，默认 `client_credentials`。 |
| `ClientSecret` | 第三方分配的密钥。 |
| `CorpID` | 第三方租户或企业 ID。Gaia 当前是 `zhwytest01`。 |
| `TokenTTLSeconds` | 第三方 token 在 Redis 中缓存的时间。Gaia 当前配置 7200 秒，也就是 2 小时。 |
| `TokenPrefix` | ESB 调第三方接口时使用的 Authorization 前缀。Gaia 要求 `Bearer`。 |
| `FixedHeaders` | ESB 转发到第三方时固定追加的请求头，多个用英文逗号分隔。 |
| `UserAgent` | ESB 调第三方时使用的 User-Agent。当前使用标准浏览器 UA，避免被第三方安全策略误拦截。 |

## Gaia 的 tenant 请求头

Gaia 的部分接口，例如 `/apply-leave-list/leaveapplysubmit`，要求请求头带：

```http
tenant: zhwytest01
```

所以 Gaia 配置里有：

```ini
FixedHeaders = tenant={CorpID}
```

ESB 转发到 Gaia 时会自动把 `{CorpID}` 替换成配置中的：

```ini
CorpID = zhwytest01
```

如果以后接入新的第三方应用，比如 `robot`，它不需要 `tenant`，就不要在 `[openapi.robot]` 里配置 `FixedHeaders`：

```ini
[openapi.robot]
BaseUrl = https://robot.example.com
AuthURL =
TokenTTLSeconds = 3600
```

`FixedHeaders` 是按应用单独生效的，不会影响其他 OpenAPI 应用。

## APP TOKEN 和应用授权

APP TOKEN 在 `[service_auth]` 中配置：

```ini
[service_auth]
ApiTokens = gaia-token-001:gaiastandard,robot-token-001:robot,admin-token:gaiastandard|robot
```

含义：

| 写法 | 含义 |
| --- | --- |
| `gaia-token-001:gaiastandard` | 只能调用 Gaia。 |
| `robot-token-001:robot` | 只能调用 robot。 |
| `admin-token:gaiastandard|robot` | 可以调用 Gaia 和 robot。 |

ESB 会先根据完整 URL 找到应用名，再检查当前 APP TOKEN 是否有权限调用这个应用。

例如：

```http
GET /api/esb/v1/https://gaiaopenapi-s.copm.com.cn/atd-webapi/api/v1/workflow/apply/list
X-API-Token: gaia-token-001
```

这个 URL 的域名匹配 `[openapi.gaiastandard]`，所以 `gaia-token-001:gaiastandard` 可以调用。

## 请求转发流程

以这个请求为例：

```http
POST /api/esb/v1/https://gaiaopenapi-s.copm.com.cn/atd-appapi/api/v1/gaiastandard/getemployeeleaveremaindata/{tenant}
X-API-Token: gaia-token-001
Content-Type: application/json
```

ESB 会执行：

1. 读取 `/api/esb/v1/` 后面的完整 URL。
2. 通过 URL 域名匹配到 `[openapi.gaiastandard]`。
3. 校验 `X-API-Token` 是否允许访问 `gaiastandard`。
4. 从 Redis 读取 Gaia 第三方 token。
5. 如果 Redis 中没有 token，或者 token 已经过期，才调用 Gaia 的 `AuthURL` 获取新 token。
6. 新 token 会写入 Redis，TTL 使用 `[openapi.gaiastandard]` 的 `TokenTTLSeconds`，当前是 7200 秒。
7. 转发请求时自动添加 `Authorization: Bearer <Gaia token>`。
8. 根据 `FixedHeaders` 自动添加 `tenant: zhwytest01`。
9. 把 URL 中的 `{tenant}` 替换成 `CorpID`。
10. 把请求方法、query、body 转发给第三方。

Gaia token 只使用 Redis 缓存，不使用进程内存缓存。Redis key 格式为：

```text
esb:openapi:token:{appName}
```

当前 Gaia 的 key 是：

```text
esb:openapi:token:gaiastandard
```

最终转发到 Gaia 的地址会变成：

```text
https://gaiaopenapi-s.copm.com.cn/atd-appapi/api/v1/gaiastandard/getemployeeleaveremaindata/zhwytest01
```

## 响应格式

ESB 对外统一返回：

```json
{
  "code": 200,
  "msg": "ok",
  "data": {}
}
```

第三方响应映射规则：

| 第三方响应 | ESB 的 `data` |
| --- | --- |
| JSON 中有 `details` | 返回 `details` |
| JSON 中没有 `details` 但有 `data` | 返回 `data` |
| JSON 中没有 `details` 和 `data` | 返回整个 JSON 对象 |
| 第三方返回非 JSON | 返回原始字符串 |

例如 Gaia 文档里成功示例是：

```json
{
  "reason": "success",
  "code": 200,
  "details": [
    {
      "name": "年假"
    }
  ]
}
```

ESB 返回：

```json
{
  "code": 200,
  "msg": "ok",
  "data": [
    {
      "name": "年假"
    }
  ]
}
```

也就是说，你业务上真正需要的数据对应 Gaia 的 `details`，不是外层的 `code`。

## Gaia 调用示例

以下示例都需要带：

```http
X-API-Token: <分配给gaiastandard的APP_TOKEN>
```

### 获取请假类别

```http
GET /api/esb/v1/https://gaiaopenapi-s.copm.com.cn/atd-webapi/api/v1/attendance-type/leave/list?emploeeId=E00001&date=2026-06-02&batchApplyLeave=0
X-API-Token: <APP_TOKEN>
```

### 查询额度

```http
POST /api/esb/v1/https://gaiaopenapi-s.copm.com.cn/atd-appapi/api/v1/gaiastandard/getemployeeleaveremaindata/{tenant}
X-API-Token: <APP_TOKEN>
Content-Type: application/json

{
  "size": 10,
  "endDate": "2026-06-30",
  "unitCode": "U001",
  "employeeId": "E00001",
  "page": 1,
  "isIncludeSubUnit": false,
  "startDate": "2026-06-01"
}
```

### 获取请假时数

```http
POST /api/esb/v1/https://gaiaopenapi-s.copm.com.cn/atd-webapi/api/v1/workflow/leave/hours
X-API-Token: <APP_TOKEN>
Content-Type: application/json

{
  "packetId": "uuid",
  "endDate": "2026-06-05",
  "leaveModel": 1,
  "startTime": "08:00",
  "endTime": "18:00",
  "autoCalc": true,
  "personIds": ["E00001"],
  "personId": "E00001",
  "startDate": "2026-06-05"
}
```

### 请假提交

```http
POST /api/esb/v1/https://gaiaopenapi-s.copm.com.cn/atd-webapi/api/v1/gaiastandard/openapi/apply-leave-list/leaveapplysubmit
X-API-Token: <APP_TOKEN>
Content-Type: application/json

{
  "processInstanceId": "xxx",
  "submitCheck": "2",
  "multipleApplyMode": "2",
  "saveData": [
    {
      "classCode": "AnnualLeave",
      "hours": 8,
      "leaveModel": 0,
      "endDate": "2026-06-20",
      "customFields": {},
      "employeeId": "E00002",
      "times": 1,
      "leaveReason": "reasonId",
      "repeat": 1,
      "files": [],
      "startTime": "09:00",
      "endTime": "18:00",
      "reasonCode": "ANNUAL",
      "detail": "请假原因",
      "autoCalc": false,
      "startDate": "2026-06-20"
    }
  ],
  "applyEmployeeId": "E00001",
  "formNo": "202606200001"
}
```

这个接口需要 `tenant` 请求头，ESB 会根据 `[openapi.gaiastandard]` 的 `FixedHeaders = tenant={CorpID}` 自动添加，不需要调用方手动传。

### 我的申请

```http
GET /api/esb/v1/https://gaiaopenapi-s.copm.com.cn/atd-webapi/api/v1/workflow/apply/list?startDate=2026-06-01&endDate=2026-06-30&formType=leave&employeeId=E00001&status=APPROVING
X-API-Token: <APP_TOKEN>
```

## 常见问题

### 为什么返回 404

通常是路径没有命中 ESB 的通用路由。正确格式必须是：

```text
/api/esb/v1/https://第三方域名/第三方路径
```

不是：

```text
/api/esb/gaiastandard/https://第三方域名/第三方路径
```

也不是：

```text
/api/esb/v1/gaiastandard/https://第三方域名/第三方路径
```

当前路由里已经没有 `:openapi` 这一段，应用是通过完整 URL 的域名自动识别的。

### 为什么提示 unknown openapi url

说明完整 URL 的域名没有匹配到任何 `[openapi.*]` 的 `BaseUrl`。

例如请求：

```text
https://gaiaopenapi-s.copm.com.cn/xxx
```

必须有配置：

```ini
[openapi.gaiastandard]
BaseUrl = https://gaiaopenapi-s.copm.com.cn
```

### 为什么返回 401

常见原因：

- 没有传 `X-API-Token`。
- APP TOKEN 写错。
- APP TOKEN 没有当前 OpenAPI 应用权限。

例如 `robot-token-001:robot` 不能调用 Gaia URL。

### 为什么第三方提示 require Bearer token

这是第三方 Gaia 要求 ESB 转发时带 `Authorization: Bearer <Gaia token>`。

请确认 `[openapi.gaiastandard]` 中配置了：

```ini
AuthURL = https://gaiaopenapi-s.copm.com.cn/identity/api/v1/oauth
TokenPrefix = Bearer
```

这和客户端调用 ESB 的 APP TOKEN 不是一回事。客户端仍然使用：

```http
X-API-Token: <APP_TOKEN>
```

### 为什么 Gaia 安全策略拦截

如果第三方返回类似“访问请求可能对网站造成安全威胁，请求已被阻断”，可能是第三方安全策略识别了请求客户端。

当前配置已使用标准浏览器 User-Agent：

```ini
UserAgent = Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36
```

如果仍被拦截，需要第三方确认是否有 IP 白名单、WAF 策略、环境域名或鉴权参数限制。
