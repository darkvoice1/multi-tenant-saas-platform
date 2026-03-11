# SSO 接入说明（设计）

## 目标
- 支持企业 SSO（SAML 2.0 / OIDC）
- 统一通过后端完成断言校验与用户映射

## 方案摘要
- OIDC：优先推荐，复用 OAuth2/OIDC 流程。
- SAML：适用于部分传统企业 IdP。

## 对接点（本项目）
- 接入层：`/auth/sso/login`、`/auth/sso/callback`
- 用户映射：`email` 或 `employee_id` 映射到租户用户
- 多租户场景：`tenant_slug` 与 IdP 配置绑定

## 数据存储建议
- 新增 `sso_providers` 表：保存租户级 IdP 配置（issuer、sso_url、cert、client_id 等）
- 新增 `user_identities` 表：绑定本地用户与外部身份标识

## 安全要点
- 强制校验 `state` 和 `nonce`
- 强制校验回调域名白名单
- 记录登录审计日志
