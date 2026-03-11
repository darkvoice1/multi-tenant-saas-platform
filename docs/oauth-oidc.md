# OAuth2 / OIDC 接入方案（设计说明）

## 目标
- 兼容第三方 IdP（如 Azure AD / Auth0 / 自建 Keycloak）
- 统一在后端完成授权码交换并签发本系统 JWT

## 推荐流程（Authorization Code + PKCE）
1. 前端跳转至 IdP 授权页（带 `client_id`、`redirect_uri`、`scope`、`state`、`code_challenge`）。
2. IdP 回调前端/后端的 `redirect_uri`，返回 `code`。
3. 后端使用 `code` + `code_verifier` 向 IdP 交换 `id_token`/`access_token`。
4. 后端校验 `id_token`（签名、aud、iss、exp）后，映射到本地用户。
5. 后端签发本系统的 JWT（access + refresh）。

## Mock 方案（已实现）
- `GET /auth/oidc/mock/authorize?tenant_id=<tenant_uuid>&email=<email>&state=demo`
- `POST /auth/oidc/callback` with `{ "code": "<code>", "state": "demo" }`

## 对接点（本项目）
- 新增配置：`OIDC_ISSUER`、`OIDC_CLIENT_ID`、`OIDC_CLIENT_SECRET`、`OIDC_REDIRECT_URL`。
- 新增接口：`POST /auth/oidc/callback` 接收 `code` 并完成 token exchange。
- 用户映射策略：优先使用 `email` 或 `sub` 作为唯一标识，映射到租户用户。

## Mock 说明
- Mock 用于演示授权码流程，不依赖外部 IdP。
- 生产环境替换为真实 IdP 接入即可。
