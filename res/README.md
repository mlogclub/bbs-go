## res 目录结构

```text
res/
├── images/
│   └── avatars/              # 内置默认头像资源
└── uploads/                  # 本地上传文件目录（运行时自动创建子目录）
    └── YYYY/MM/DD/<md5>.<ext>
```

### 说明

- 安装包中只包含 `res/images/` 和 `res/README.md`。
- `images/` 为随程序发布的静态内置资源目录。
- `uploads/` 通过服务端路由 `/res` 直接对外访问。
- 上传接口返回的本地文件 URL 统一为：`/res/uploads/...`。
- `uploads/` 为运行时数据目录，不参与构建打包，且已在 `server/.gitignore` 中忽略。
- 不要将用户上传文件、测试文件或本地生成的数据文件提交到仓库。
