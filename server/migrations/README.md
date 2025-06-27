# 数据库迁移文件结构

本项目支持多语言数据库初始化，migration文件按以下结构组织：

```
migrations/
├── base/                           # 基础表结构迁移
│   ├── 001_init_schema.up.sql     # 创建所有表结构
│   ├── 001_init_schema.down.sql   # 删除表结构（回滚）
│   ├── 002_add_indexes.up.sql     # 添加索引
│   └── 002_add_indexes.down.sql   # 删除索引（回滚）
└── data/                          # 数据迁移（多语言）
    ├── zh-CN/                     # 中文数据
    │   ├── 001_init_roles.up.sql
    │   ├── 001_init_roles.down.sql
    │   ├── 002_init_nodes.up.sql
    │   └── 002_init_nodes.down.sql
    └── en-US/                     # 英文数据
        ├── 001_init_roles.up.sql
        ├── 001_init_roles.down.sql
        ├── 002_init_nodes.up.sql
        └── 002_init_nodes.down.sql
```

## 执行顺序

1. 首先执行 `base/` 目录下的结构迁移
2. 然后根据用户选择的语言执行对应的数据迁移

## 文件命名规范

- 文件名格式：`{序号}_{描述}.{up|down}.sql`
- 序号：3位数字，从001开始
- 描述：英文，用下划线分隔
- up：升级脚本
- down：回滚脚本

## 多语言支持

- 表结构在 `base/` 目录，语言无关
- 初始数据在 `data/{language}/` 目录，支持多语言
- 如果指定语言的迁移文件不存在，会回退到 `zh-CN`

## 示例

安装时指定语言参数：
```json
{
  "language": "en-US",
  "siteTitle": "My BBS",
  "siteDescription": "A community website",
  "dbConfig": {...},
  "username": "admin",
  "password": "123456"
}
``` 