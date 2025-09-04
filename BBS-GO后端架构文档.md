# BBS-GO 后端架构文档

## 1. 项目概述

BBS-GO 是一个基于 Go 语言开发的现代化论坛系统，采用前后端分离架构，后端使用 Iris 框架提供 REST API 服务。

## 2. 整体架构

### 2.1 架构模式
- **分层架构**：采用经典的MVC分层架构
- **前后端分离**：后端提供REST API，前端独立部署
- **缓存策略**：内存缓存 + 数据库持久化

### 2.2 技术栈
- **Web框架**：Iris v12
- **ORM**：GORM
- **数据库**：支持MySQL/PostgreSQL等
- **缓存**：内存缓存
- **搜索引擎**：MeiliSearch（替换了ES）
- **邮件服务**：SMTP
- **文件上传**：本地/云存储

## 3. 目录结构

```
server/
├── cmd/                    # 命令行工具
│   └── test/              # 测试工具
├── internal/              # 内部包，不对外暴露
│   ├── cache/            # 缓存层
│   ├── controllers/      # 控制器层
│   │   ├── admin/       # 管理后台API
│   │   ├── api/         # 前端API
│   │   └── render/      # 数据渲染层
│   ├── install/         # 安装配置
│   ├── middleware/      # 中间件
│   ├── models/          # 数据模型
│   │   ├── constants/   # 常量定义
│   │   └── dto/        # 数据传输对象
│   ├── pkg/            # 内部工具包
│   ├── repositories/   # 数据访问层
│   ├── scheduler/      # 定时任务
│   ├── server/        # 服务器配置
│   ├── services/      # 业务逻辑层
│   │   └── eventhandler/ # 事件处理
│   └── spam/          # 反垃圾
├── locales/           # 国际化文件
├── logs/              # 日志文件
└── migrations/        # 数据库迁移
```

## 4. 核心组件架构

### 4.1 分层架构图

```mermaid
graph TB
    A[HTTP请求] --> B[路由层Router]
    B --> C[中间件层Middleware]
    C --> D[控制器层Controller]
    D --> E[业务服务层Service]
    E --> F[数据访问层Repository]
    F --> G[数据库Database]
    
    D --> H[渲染层Render]
    E --> I[缓存层Cache]
    E --> J[外部服务External]
    
    subgraph "中间件"
    C1[CORS中间件]
    C2[认证中间件]
    C3[日志中间件]
    C4[错误恢复中间件]
    end
    
    subgraph "外部服务"
    J1[邮件服务]
    J2[搜索引擎]
    J3[文件存储]
    end
```

### 4.2 核心模型关系图

```mermaid
erDiagram
    User ||--o{ Topic : "发布"
    User ||--o{ Comment : "评论"
    User ||--o{ Article : "创作"
    User ||--o{ UserToken : "认证"
    User ||--o{ Message : "接收"
    User ||--o{ Favorite : "收藏"
    User ||--o{ UserLike : "点赞"
    
    Topic ||--o{ Comment : "包含"
    Topic ||--o{ TopicTag : "标记"
    Topic }o--|| TopicNode : "归属"
    
    Article ||--o{ Comment : "包含"
    Article ||--o{ ArticleTag : "标记"
    
    Tag ||--o{ TopicTag : "关联"
    Tag ||--o{ ArticleTag : "关联"
    
    User {
        int64 id
        string username
        string email
        string nickname
        string avatar
        string password
        int status
        int64 createTime
    }
    
    Topic {
        int64 id
        int64 userId
        int64 nodeId
        string title
        string content
        int status
        int64 viewCount
        int64 commentCount
        int64 createTime
    }
    
    Comment {
        int64 id
        int64 userId
        string entityType
        int64 entityId
        string content
        int64 quoteId
        int status
        int64 createTime
    }
```

## 5. 功能模块详细架构

### 5.1 用户认证系统

#### 5.1.1 认证流程时序图

```mermaid
sequenceDiagram
    participant C as 客户端
    participant LC as LoginController
    participant US as UserService
    participant UTS as UserTokenService
    participant DB as 数据库
    participant Cache as 缓存

    Note over C,Cache: 用户注册流程
    C->>LC: POST /api/login/signup
    LC->>LC: 验证验证码
    LC->>US: SignUp(username, email, password)
    US->>US: 验证用户名、邮箱唯一性
    US->>US: 密码加密
    US->>DB: 创建用户记录
    DB-->>US: 返回用户信息
    US-->>LC: 返回新用户
    LC->>UTS: Generate(userId)
    UTS->>DB: 创建Token记录
    UTS-->>LC: 返回Token
    LC-->>C: 返回用户信息和Token

    Note over C,Cache: 用户登录流程
    C->>LC: POST /api/login/signin
    LC->>LC: 验证验证码
    LC->>US: SignIn(username, password)
    US->>DB: 根据用户名/邮箱查询用户
    US->>US: 验证密码
    US-->>LC: 返回用户信息
    LC->>UTS: Generate(userId)
    UTS->>DB: 创建Token记录
    UTS-->>LC: 返回Token
    LC-->>C: 返回用户信息和Token

    Note over C,Cache: Token验证流程
    C->>+UTS: 携带Token请求
    UTS->>Cache: 从缓存获取Token
    alt Token缓存命中
        Cache-->>UTS: 返回Token信息
    else Token缓存未命中
        UTS->>DB: 查询Token
        DB-->>UTS: 返回Token信息
        UTS->>Cache: 缓存Token信息
    end
    UTS->>UTS: 验证Token有效性
    UTS->>Cache: 获取用户信息
    UTS-->>-C: 返回当前用户
```

#### 5.1.2 Token管理机制

- **Token生成**：使用UUID生成32位随机字符串
- **Token存储**：数据库持久化 + 内存缓存
- **Token验证**：支持Cookie、Header（Authorization/X-User-Token）多种方式
- **Token过期**：可配置过期时间，默认30天
- **Token失效**：登出时标记为删除状态

### 5.2 消息通知系统

#### 5.2.1 消息通知架构

**重要发现：该系统采用数据库轮询方式实现消息通知，无WebSocket推送，无消息队列**

```mermaid
graph TD
    A[用户操作] --> B[业务服务]
    B --> C[MessageService.SendMsg]
    C --> D[数据库存储消息]
    C --> E[发送邮件通知]
    
    F[前端轮询] --> G[GET /api/user/messages]
    G --> H[MessageService.Find]
    H --> I[返回消息列表]
    I --> J[标记已读]
    
    K[未读消息数] --> L[GET /api/user/msg/recent]
    L --> M[返回最近3条未读消息]
```

#### 5.2.2 消息通知时序图

```mermaid
sequenceDiagram
    participant U1 as 用户A
    participant S as 业务服务
    participant MS as MessageService
    participant ES as EmailService
    participant DB as 数据库
    participant U2 as 用户B
    participant Frontend as 前端

    Note over U1,Frontend: 消息发送流程
    U1->>S: 执行操作(评论/点赞/收藏)
    S->>MS: SendMsg(from, to, msgType, content)
    MS->>DB: 创建消息记录
    MS->>ES: SendEmailNotice(message)
    ES->>ES: 发送邮件通知
    
    Note over U1,Frontend: 消息接收流程（前端轮询）
    Frontend->>Frontend: 定时轮询（如每30秒）
    Frontend->>MS: GET /api/user/msg/recent
    MS->>DB: 查询未读消息数量和最近3条
    MS-->>Frontend: 返回未读消息信息
    Frontend->>Frontend: 更新UI提示
    
    Note over U1,Frontend: 消息查看流程
    U2->>Frontend: 点击查看消息
    Frontend->>MS: GET /api/user/messages
    MS->>DB: 查询用户所有消息
    MS->>DB: 标记所有消息为已读
    MS-->>Frontend: 返回消息列表
    Frontend->>Frontend: 显示消息内容
```

#### 5.2.3 消息类型和触发机制

```go
// 消息类型定义
const (
    TypeTopicComment   Type = 0 // 收到话题评论
    TypeCommentReply   Type = 1 // 收到他人回复  
    TypeTopicLike      Type = 2 // 收到点赞
    TypeTopicFavorite  Type = 3 // 话题被收藏
    TypeTopicRecommend Type = 4 // 话题被设为推荐
    TypeTopicDelete    Type = 5 // 话题被删除
    TypeArticleComment Type = 6 // 收到文章评论
)
```

**通知机制特点：**
- ✅ **数据库存储**：所有消息持久化存储在数据库
- ✅ **邮件通知**：支持邮件推送（除话题删除外）
- ❌ **实时推送**：无WebSocket实时推送
- ❌ **消息队列**：无异步消息队列处理
- ✅ **前端轮询**：前端通过定时请求获取未读消息

### 5.3 内容管理系统

#### 5.3.1 话题发布流程

```mermaid
sequenceDiagram
    participant User as 用户
    participant TC as TopicController
    participant TS as TopicService
    participant US as UserService
    participant SS as SpamService
    participant DB as 数据库
    participant Cache as 缓存

    User->>TC: POST /api/topic/create
    TC->>US: CheckPostStatus(user)
    US-->>TC: 验证用户状态
    TC->>SS: CheckPost(content)
    SS-->>TC: 反垃圾检查
    TC->>TS: Create(topic)
    TS->>DB: 保存话题
    TS->>US: IncrTopicCount(userId)
    TS->>Cache: 清除相关缓存
    TS-->>TC: 返回话题信息
    TC-->>User: 返回成功响应
```

#### 5.3.2 评论系统流程

```mermaid
sequenceDiagram
    participant User as 用户
    participant CC as CommentController
    participant CS as CommentService
    participant MS as MessageService
    participant DB as 数据库

    User->>CC: POST /api/comment/create
    CC->>CS: Create(comment)
    CS->>DB: 保存评论
    CS->>CS: 更新实体评论数
    CS->>MS: 发送消息通知
    MS->>DB: 创建消息记录
    MS->>MS: 发送邮件通知
    CS-->>CC: 返回评论信息
    CC-->>User: 返回成功响应
```

### 5.4 搜索系统

```mermaid
graph TD
    A[搜索请求] --> B[SearchController]
    B --> C[SearchService]
    C --> D[MeiliSearch引擎]
    D --> E[返回搜索结果]
    
    F[内容更新] --> G[自动索引更新]
    G --> D
```

## 6. 数据库设计

### 6.1 核心表结构

#### 用户表 (t_user)
- id: 主键
- username: 用户名（唯一）
- email: 邮箱（唯一）
- password: 加密密码
- nickname: 昵称
- avatar: 头像
- status: 状态
- score: 积分
- create_time: 创建时间

#### 话题表 (t_topic)
- id: 主键
- user_id: 发布用户ID
- node_id: 节点ID
- title: 标题
- content: 内容
- content_type: 内容类型
- status: 状态
- view_count: 浏览数
- comment_count: 评论数
- create_time: 创建时间

#### 消息表 (t_message)
- id: 主键
- from_id: 发送者ID
- user_id: 接收者ID
- title: 标题
- content: 内容
- type: 消息类型
- status: 读取状态
- create_time: 创建时间

## 7. 缓存策略

### 7.1 缓存架构

```mermaid
graph TD
    A[应用层] --> B[缓存层]
    B --> C[用户缓存UserCache]
    B --> D[Token缓存UserTokenCache]
    B --> E[系统配置缓存SysConfigCache]
    B --> F[标签缓存TagCache]
    B --> G[话题缓存TopicCache]
    
    B --> H[数据库]
    
    I[缓存更新策略]
    I --> J[写入时失效]
    I --> K[定时刷新]
    I --> L[LRU淘汰]
```

### 7.2 缓存管理
- **用户缓存**：缓存用户基本信息，更新时失效
- **Token缓存**：缓存用户Token，提高认证性能
- **配置缓存**：缓存系统配置，减少数据库查询
- **内容缓存**：缓存热门内容，提高访问速度

## 8. 安全机制

### 8.1 认证授权
- **Token认证**：基于UUID的Token机制
- **权限控制**：基于角色的权限管理
- **管理员认证**：独立的管理员认证中间件

### 8.2 安全防护
- **验证码**：登录注册需要验证码
- **反垃圾**：多种反垃圾策略
- **输入验证**：严格的参数验证
- **SQL注入防护**：ORM框架防护

### 8.3 反垃圾系统

```mermaid
graph TD
    A[用户发布内容] --> B[SpamService检查]
    B --> C[频率限制策略]
    B --> D[验证码策略]  
    B --> E[邮箱验证策略]
    B --> F[违禁词过滤]
    
    C --> G{检查通过?}
    D --> G
    E --> G
    F --> G
    
    G -->|是| H[发布成功]
    G -->|否| I[拒绝发布]
```

## 9. 性能优化

### 9.1 数据库优化
- **索引优化**：合理设计数据库索引
- **查询优化**：使用条件构建器优化查询
- **分页查询**：游标分页提高性能

### 9.2 缓存优化
- **多级缓存**：内存缓存 + 数据库缓存
- **缓存预热**：系统启动时预热热点数据
- **缓存更新**：写入时失效策略

### 9.3 其他优化
- **静态资源**：支持静态资源压缩和缓存
- **数据库连接池**：GORM连接池管理
- **异步处理**：邮件发送等异步处理

## 10. 部署架构

### 10.1 推荐部署架构

```mermaid
graph TD
    A[负载均衡器] --> B[Web服务器1]
    A --> C[Web服务器2]
    A --> D[Web服务器N]
    
    B --> E[数据库主库]
    C --> E
    D --> E
    
    E --> F[数据库从库]
    
    G[Redis缓存] --> B
    G --> C
    G --> D
    
    H[MeiliSearch] --> B
    H --> C
    H --> D
    
    I[文件存储] --> B
    I --> C
    I --> D
```

### 10.2 环境配置
- **开发环境**：单机部署，SQLite数据库
- **测试环境**：独立数据库，完整功能测试
- **生产环境**：集群部署，读写分离，缓存集群

## 11. 监控和日志

### 11.1 日志系统
- **访问日志**：记录所有API请求
- **错误日志**：记录系统错误和异常
- **业务日志**：记录重要业务操作
- **操作日志**：管理员操作审计

### 11.2 监控指标
- **性能监控**：响应时间、吞吐量
- **错误监控**：错误率、异常统计
- **业务监控**：用户活跃度、内容统计
- **系统监控**：CPU、内存、磁盘使用率

## 12. 总结

BBS-GO后端采用现代化的Go技术栈，具有以下特点：

### 优点：
- **架构清晰**：分层架构，职责明确
- **性能优良**：Go语言高性能，合理的缓存策略
- **功能完整**：涵盖论坛所需的核心功能
- **安全可靠**：多重安全防护机制
- **易于维护**：代码结构清晰，注释完善

### 可优化点：
- **消息通知**：可考虑引入WebSocket实现实时推送
- **消息队列**：可引入消息队列处理异步任务
- **微服务化**：可考虑将大模块拆分为微服务
- **容器化**：可使用Docker容器化部署

### 技术特色：
1. **无实时推送**：消息通知采用前端轮询 + 邮件通知方式
2. **简化架构**：避免了复杂的消息队列，降低运维成本
3. **缓存优先**：大量使用内存缓存提高性能
4. **安全为先**：多层安全防护，防范常见安全威胁

该架构适合中小型论坛社区，在保证功能完整性的同时，保持了较低的技术复杂度和运维成本。