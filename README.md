# 🛒 小米商城 · Go 微服务电商系统

以小米商城为原型的 B2C 电商练手项目：**Gin 接入层 + go-micro 微服务 + Consul 注册发现 + MySQL + Redis**，覆盖「商品浏览 → 购物车 → 下单 → 支付 → 订单查询」完整闭环，后台有一套 RBAC 权限管理。

<div align="center">

![Go](https://img.shields.io/badge/Go-1.17+-00ADD8?logo=go)
![Gin](https://img.shields.io/badge/Gin-1.7.7-00ADD8)
![go-micro](https://img.shields.io/badge/go--micro-v4.9.0-326CE5)
![Consul](https://img.shields.io/badge/Consul-注册中心-F24C53)
![MySQL](https://img.shields.io/badge/MySQL-8.0+-4479A1)
![Redis](https://img.shields.io/badge/Redis-6.0+-DC382D)

</div>

---

## 📖 关于文档

仓库只提交 **README.md** 这一份文档。

开发过程中的另外几份（代码约定、修复记录、待优化清单、面试背诵稿）属于本地笔记，已加进 `.gitignore`，不随仓库上传。

---

## 🚀 快速开始

### 0. 环境要求

| 依赖 | 版本 | 说明 |
|---|---|---|
| Go | 1.17+ | |
| MySQL | 8.0+ | 需要建 **两个库** |
| Redis | 6.0+ | 缓存（配置里可关） |
| Consul | 任意 | **必须**，服务注册发现 |

### 1. 建库（两个库，别漏）

> ⚠️ 仓库里的 `.sql` 是 Navicat 导出的**纯表结构 + 数据**，**没有 `CREATE DATABASE`**，所以必须先建库再导入。

```bash
mysql -u root -p -e "CREATE DATABASE ginxiaomi DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
mysql -u root -p -e "CREATE DATABASE ginxiaomirbac DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 注意：PowerShell 不支持 < 重定向，请用 cmd，或改用下面的 source 写法
mysql -u root -p ginxiaomi     < ginxiaomi.sql
mysql -u root -p ginxiaomirbac < ginxiaomirbac.sql
```

PowerShell 写法：

```powershell
mysql -u root -p ginxiaomi     -e "source D:/workspace/xiaomishaop_go/xiaomishop-main/ginxiaomi.sql"
mysql -u root -p ginxiaomirbac -e "source D:/workspace/xiaomishaop_go/xiaomishop-main/ginxiaomirbac.sql"
```

| 库 | 表数 | 归属 |
|---|---|---|
| `ginxiaomi` | 15 张（商品 / 订单 / 用户 / 运营） | 主站 ginshop |
| `ginxiaomirbac` | 4 张（manager / role / access / role_access） | rbac 服务 |

> 两个库**没有外键关联**，跨库需求走 gRPC。`admin` 只存在于 `manager` 表，前台手机号账号只存在于 `user` 表，**互不通用**。

### 2. 启动 Consul（必须最先）

```bash
docker run -d --name consul -p 8500:8500 consul
# 或本地已安装 consul
consul agent -dev
```

验证：浏览器打开 http://localhost:8500/ui 能看到 Consul 界面。

### 3. 启动微服务

```bash
# 终端 A —— 验证码服务
cd server/captcha && go run main.go

# 终端 B —— 权限服务
cd server/rbac && go run main.go
```

看到 `Registry [consul] Registering node: captcha-xxxx` 才算注册成功。

### 4. 启动主站

```bash
# 终端 C —— 必须在 client/ginshop 目录下启动！
cd client/ginshop && go run main.go
```

> ⚠️ **一定要在 `client/ginshop` 目录执行**。`main.go` 里 `r.Static("/static", "./static")` 和 `LoadHTMLGlob("templates/...")` 都是相对路径，在仓库根目录启动会导致静态资源和模板全部 404，表现为样式错乱、按钮点不动。

### 5. 访问

| 入口 | 地址 |
|---|---|
| 前台商城 | http://localhost:8080 |
| 后台管理 | http://localhost:8080/admin/login |
| API 服务（独立进程） | http://localhost:3000 |
| Consul 控制台 | http://localhost:8500/ui |

---

## 🔑 测试账号

| 端 | 入口 | 账号 | 密码 |
|---|---|---|---|
| **后台管理** | `/admin/login` | `admin` | `123456` |
| **前台会员** | `/pass/login` | `11111111111`（11 个 1） | `123456` |

> 后台账号在 `ginxiaomirbac.manager` 表，前台账号在 `ginxiaomi.user` 表，两套体系不互通。
> 想注册新会员：`/pass/registerStep1` → 手机号 + 图形验证码 → **短信验证码在 ginshop 控制台里打印**（见下方「已知限制」）。

---

## 🏗 架构

```
                    浏览器
                      │
        ┌─────────────┴──────────────┐
        │                            │
   ginshop :8080                ginWebApi :3000
   （前台 + 后台，HTML 渲染）      （纯 JSON API）
        │                            │
        └──────────┬─────────────────┘
                   │  gRPC（按服务名调用，不写 IP）
        ┌──────────┴──────────┐
        │                     │
   rbac 服务              captcha 服务
  （权限/登录）            （图形验证码）
        │                     │
   ginxiaomirbac 库      内存 store（无 DB）
        │
        └──► Consul :8500  ◄── 注册与发现
        │
   ginxiaomi 库 + Redis
```

### 模块职责

| 模块 | 目录 | 端口 | 职责 | 数据 |
|---|---|---|---|---|
| 主站 | `client/ginshop` | 8080 | 前台商城 + 后台管理 | `ginxiaomi`、Redis |
| API | `client/ginWebApi` | 3000 | 纯 JSON API 示例 | 只调 captcha |
| 验证码客户端 | `client/captcha-client` | — | 独立调用示例 | — |
| 权限服务 | `server/rbac` | gRPC | 管理员登录、角色、权限 | `ginxiaomirbac` |
| 验证码服务 | `server/captcha` | gRPC | 图形验证码生成/校验 | 内存 store |

**每个模块有独立的 `go.mod`**（共 5 个），import 路径分别是 `ginshop/...`、`rbac/...`、`captcha/...`。

---

## 📁 目录结构

```
xiaomishop/
├── client/
│   ├── ginshop/                 # 主站（前台 + 后台）
│   │   ├── conf/app.ini         # MySQL / Redis / OSS 配置
│   │   ├── controllers/
│   │   │   ├── admin/           # 后台：登录、商品、分类、导航、角色、权限…
│   │   │   ├── itying/          # 前台：首页、商品、购物车、下单、支付、用户
│   │   │   └── api/             # 简单 JSON 接口
│   │   ├── middlewares/         # adminAuth（RBAC 鉴权）、userAuth（登录校验）
│   │   ├── models/              # 模型 + DB/Redis/Cookie/验证码封装
│   │   ├── proto/               # gRPC 生成代码
│   │   ├── routers/             # 路由注册
│   │   ├── static/              # CSS / JS / 图片
│   │   └── templates/           # HTML 模板
│   ├── ginWebApi/               # 独立 API 服务
│   └── captcha-client/          # 验证码调用示例
├── server/
│   ├── rbac/                    # 权限微服务（handler / models / proto）
│   └── captcha/                 # 验证码微服务
├── ginxiaomi.sql                # 主站库脚本（根目录）
├── ginxiaomirbac.sql            # 权限库脚本（根目录）
└── docker-compose.yml
```

---

## ✨ 功能清单

**前台**
- 首页轮播 / 分类导航 / 商品推荐
- 商品分类列表（分页）、商品详情（多规格、多图、关联商品/赠品/配件）
- 购物车（Cookie 存储，增删改查、勾选、数量增减、实时算价）
- 下单结算（收货地址、防重复提交、订单 + 明细落库、库存扣减）
- 支付（支付宝网页支付 / 微信 Native 扫码）
- 会员注册（图形验证码 + 短信验证码三步）、登录、订单列表与详情

**后台**
- RBAC：管理员 / 角色 / 权限（模块-菜单-操作三级节点，URL 级鉴权中间件）
- 商品管理（多图上传、规格属性、缩略图）、商品分类、商品类型
- 导航、轮播图、站点设置
- 图形验证码登录

---

## ⚙️ 配置说明

`client/ginshop/conf/app.ini`：

```ini
app_name   = 小米商城
admin_path = /admin
excludeAuthPath="/,/welcome,/loginOut"   # 免鉴权路径

[mysql]
ip = 127.0.0.1
database = ginxiaomi            # 主站库

[redis]
redisEnable = true              # 设为 false 则不使用缓存

[oss]
status = 0                      # 0=本地存储，1=阿里云 OSS
```

`server/rbac/conf/app.ini`：

```ini
[mysql]
database = ginxiaomirbac        # 权限库

[consul]
addr = 127.0.0.1:8500
```

---

## 🩹 已修复问题

项目完成后做了一轮代码审计，按「会不会造成资损 / 数据不可恢复」排序，修复了以下问题：

| # | 时间 | 问题 | 修复 |
|---|---|---|---|
| 1 | 2026-10-18 | 支付回调只验签、从不更新订单，付了钱系统不知道 | 补金额校验 + 幂等更新 + 返回 ACK，回调闭环 |
| 2 | 2026-10-18 | 下单金额取自 Cookie，可解密改价后下单 | 只信任商品 ID 与数量，价格一律服务端查库重算 |
| 3 | 2026-10-18 | 下单无事务、库存字段从不扣减，会超卖 | 事务包裹「建单 + 明细 + 扣库存」，扣减用乐观锁 |
| 4 | 2026-10-18 | 3 处 SQL 注入（表名/字段名直拼、username 拼接、LIKE 拼接） | 表名字段名走白名单，值全部参数化 |
| 5 | 2026-10-18 | 鉴权中间件只 `return` 不 `Abort`，未登录仍执行到业务 handler | 所有拦截分支统一 `Abort() + return` |
| 6 | 2026-10-18 | 订单搜索无结果时拼出 `id in ()`，整页 500 | `Pluck` 取订单 ID + `IN ?`，无结果直接返回空列表 |
| 7 | 2026-10-18 | rbac / captcha 不可用时 `res` 为 nil，空指针 panic | 加超时、判 err、判 nil，鉴权失败按 fail-close 拒绝 |
| 8 | 2026-10-18 | 后台登录验证码只有 2 位（参数顺序误解） | 改为 4 位，图片尺寸同步放大到 120×40 |
| 9 | 2026-10-18 | 注册页「立即注册」点击无反应，失败时无提示 | 加全局 ajax 错误提示 + 验证码未加载校验 |

> 详细记录（改前 / 改后 / 为什么 / 怎么验证）保留在本地开发笔记中，未纳入本仓库。

---

## ⚠️ 已知限制（务必先看）

1. **短信验证码是模拟的**
   注册第 2 步的验证码不会真的发短信，而是打印在 ginshop 的控制台：
   ```
   ----------自己集成发送短信的接口--------
   1234
   ```
   看到这行就直接用那 4 位数字。

2. **支付回调地址写死成内网 IP**（`118.123.14.36:8005`，代码里有 `TODO`）。
   真实联调必须换成公网可访问地址，否则支付宝/微信调不到你的服务。

3. **库存为 0 的商品无法下单**。
   下单流程已做库存扣减，`goods.goods_number = 0` 的商品会提示「库存不足」，演示前先把库存设成大于 0。

4. **密钥硬编码**：支付宝私钥、微信商户密钥、OSS AK/SK、DES 密钥、session 密钥都在源码里。
   生产必须改为环境变量 + 平台侧密钥轮换。

5. **`docker compose up` 起不来完整业务**：
   compose 里建库是 `xiaomishop`，应用连的是 `ginxiaomi`/`ginxiaomirbac`，且没有挂初始化 SQL，也没有 ginshop 的 Dockerfile。目前只建议用它起 consul。

6. **单实例设计**：session 用 cookie store、验证码用内存 store、图片传本地磁盘——多实例部署会失效，需要改成 Redis 存储 + 对象存储。

---

## 🔌 主要接口

### REST（ginshop :8080）

| 分类 | 接口 |
|---|---|
| 商品 | `GET /`、`GET /category:id`、`GET /detail?id=`、`GET /product/getImgList` |
| 购物车 | `GET /cart`、`/cart/addCart`、`/cart/incCart`、`/cart/decCart`、`/cart/delCart` |
| 下单 | `GET /buy/checkout`、`POST /buy/doCheckout`、`GET /buy/pay`、`GET /buy/orderPayStatus` |
| 支付 | `GET /alipay?id=`、`POST /alipayNotify`、`GET /wxpay?id=`、`POST /wxpay/notify` |
| 会员 | `GET /pass/login`、`/pass/registerStep1~3`、`POST /pass/doLogin`、`POST /pass/doRegister` |
| 后台 | `/admin/**`（全部经过 RBAC 鉴权中间件） |

### gRPC（服务间）

```protobuf
service RbacLogin   { rpc Login(LoginRequest) returns (LoginResponse); }
service RbacRole    { rpc RoleGet / RoleAdd / RoleAuth / MiddlewaresAuth ... }
service RbacManager { rpc ManagerGet / ManagerAdd ... }
service Captcha     { rpc MakeCaptcha / VerifyCaptcha }
```

---

## 🧪 构建与测试

```bash
# 编译（每个模块独立，根目录没有 go.mod）
cd client/ginshop        && go build ./...
cd client/ginWebApi      && go build ./...
cd client/captcha-client && go build ./...
cd server/rbac           && go build ./...
cd server/captcha        && go build ./...

# 测试（暂无自动化测试文件）
go test ./... -race
```

---

## ❓ 常见问题

| 现象 | 原因 | 处理 |
|---|---|---|
| 登录页验证码空白 / 500 | Consul 或 captcha 没启动 | 先起 Consul，再起 captcha；验证码框会显示「验证码服务未启动」 |
| 注册页「立即注册」点了没反应 | 静态资源 404，jQuery 没加载 | 必须在 `client/ginshop` 目录启动；F12 看 Console |
| 提示「权限服务不可用」 | rbac 没启动，或启动顺序反了没注册上 | 起 rbac；若在 Consul 之前启动过，重启一次 rbac |
| `Error 1049: Unknown database` | 库没建（SQL 文件不含 CREATE DATABASE） | 按第 1 步建库并导入 |
| 登录提示「用户名或密码错误」 | 两套账号体系搞混了 | 后台 `admin/123456`，前台 `11111111111/123456` |
| 页面样式全丢、图片 404 | 工作目录不对 | `cd client/ginshop` 后重启 |

---

## 📄 License

仅用于学习交流。
