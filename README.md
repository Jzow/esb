# ESB 服务

这是一个基于 Gin 的 ESB OpenAPI 代理服务。当前对外只暴露通用代理入口，由调用方把第三方完整 URL 放到 `/api/esb/v1/` 后面，ESB 根据 URL 自动识别要转发到哪个 OpenAPI 应用。

## 快速入口

- 服务认证说明：[docs/auth.md](docs/auth.md)
- OpenAPI 通用代理说明：[docs/esb-openapi.md](docs/esb-openapi.md)

## 当前核心规则

1. 客户端调用 ESB 时使用 APP TOKEN，推荐放在请求头 `X-API-Token`。
2. APP TOKEN 是普通应用令牌，不是 JWT，不需要加 `Bearer`。
3. ESB 调用 Gaia 等第三方 OpenAPI 时，会自动获取第三方 token，并按第三方要求添加 `Authorization: Bearer <token>`。
4. 第三方 token 只缓存到 Redis，不使用进程内存缓存；Redis 过期后才重新请求第三方鉴权接口。
5. 第三方应用配置统一放在 `conf/app.ini` 的 `[openapi.{appName}]` 段中。
6. 如果某个第三方需要固定请求头，例如 Gaia 的 `tenant`，只在对应的 `[openapi.gaiastandard]` 里配置 `FixedHeaders`，其他应用不配置即可。

## 通用代理地址

```http
GET    /api/esb/v1/{第三方完整URL}
POST   /api/esb/v1/{第三方完整URL}
PUT    /api/esb/v1/{第三方完整URL}
DELETE /api/esb/v1/{第三方完整URL}
```

示例：

```http
GET /api/esb/v1/https://gaiaopenapi-s.copm.com.cn/atd-webapi/api/v1/workflow/apply/list
X-API-Token: <分配给调用方的APP_TOKEN>
```

更多配置、agent 填写方式、Gaia 示例和排错说明见 [docs/esb-openapi.md](docs/esb-openapi.md)。
