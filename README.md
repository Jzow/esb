#ESB

## Service Authentication

See `docs/auth.md`.

## Gaia Proxy APIs

See `docs/gaia.md`.

## Gaia OpenAPI（流程）对接配置

在 `conf/app.ini` 中新增配置段：

```ini
[gaia-openapi]
BaseURL = https://gaiaopenapi-s.copm.com.cn
AuthURL = https://gaiaopenapi-s.copm.com.cn/identity/api/v1/oauth
GrantType = client_credentials
ClientSecret = <your client_secret>
CorpID = <your corp_id>
TokenTTLSeconds = 3600

SubmitLeavePath = <请假提交接口路径>
MyApplicationsPath = <我的申请接口路径>
LeaveQuotaPath = <查询额度接口路径>
LeaveTypesPath = <获取请假类别接口路径>
LeaveHoursPath = <获取请假时数接口路径>
ExceptionListPath = <查询异常列表接口路径>
```

支持通过不同的 `app-*.ini` 配置文件来区分环境（开发/测试/生产），启动时替换为对应配置文件即可。
