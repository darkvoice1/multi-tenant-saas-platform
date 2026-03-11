# 权限矩阵与接口清单

## 角色定义
- `admin`：租户管理员，拥有全量权限。
- `manager`：项目/任务管理权限。
- `member`：任务执行与部分写权限。
- `guest`：只读。

## 权限矩阵（示例）
| 权限 | admin | manager | member | guest |
| --- | --- | --- | --- | --- |
| tenant:read | ✅ | ✅ | ❌ | ❌ |
| tenant:write | ✅ | ❌ | ❌ | ❌ |
| user:read | ✅ | ✅ | ❌ | ❌ |
| user:write | ✅ | ❌ | ❌ | ❌ |
| project:read | ✅ | ✅ | ✅ | ✅ |
| project:write | ✅ | ✅ | ❌ | ❌ |
| task:read | ✅ | ✅ | ✅ | ✅ |
| task:write | ✅ | ✅ | ✅ | ❌ |
| audit:read | ✅ | ❌ | ❌ | ❌ |
| audit:write | ✅ | ❌ | ❌ | ❌ |
| admin:ping | ✅ | ❌ | ❌ | ❌ |

## 接口权限清单（当前阶段）
- `POST /auth/login`：公开
- `POST /auth/refresh`：公开（仅刷新令牌）
- `POST /auth/logout`：公开（需要 refresh_token）
- `GET /auth/me`：登录可访问
- `GET /api/tenant/echo`：登录可访问（需租户上下文）
- `GET /api/admin/ping`：仅 admin

> 后续业务接口按权限矩阵继续扩展。
