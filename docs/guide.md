---
sidebar: auto
---

# 使用帮助

## 交流群
- QQ群号：653248175
- 扫码进群：

![bbs-go-qq.jpg](https://i.loli.net/2019/09/12/eiKSXycnDB7R6Gw.jpg)

## 简介
bbs-go 是一款基于Go语言开发的论坛系统，采用前后端分离技术，Go语言提供api进行数据支撑，用户界面使用Nuxt.js进行渲染，后台界面基于element-ui。主要功能如下：

- 用户中心（用户注册、登录、个人资料管理）
- 论坛功能
- 多人博客
- 站内消息
- 收藏功能
- 站内消息

## 演示
[https://mlog.club](https://mlog.club)

## 源码
- Github：[https://github.com/mlogclub/bbs-go](https://github.com/mlogclub/bbs-go)
- 码云：[https://gitee.com/mlog/bbs-go](https://gitee.com/mlog/bbs-go)


## 技术栈
- iris (https://github.com/kataras/iris) Go语言 mvc 框架
- gorm (http://gorm.io/) Go语言 orm 框架
- resty (https://github.com/go-resty/resty) Go语言好用的 http-client
- cron (https://github.com/robfig/cron) 定时任务
- goquery（https://github.com/PuerkitoBio/goquery）html dom 元素解析
- Nuxt.js (https://nuxtjs.org) 基于Vue的服务端渲染框架
- Element-UI (https://element.eleme.cn) 饿了么开源的基于 vue.js 的前端库
- Vditor (https://github.com/b3log/vditor) Markdown 编辑器

## Startup
### 安装依赖
mlog-club 的依赖是使用go mod来进行管理的，go mod使用帮助看这里：https://mlog.club/topic/9
```shell
# 第一步 clone 代码
git clone https://github.com/mlogclub/mlog.git

# 第二步 安装依赖
cd mlog
go mod tidy
```

### 配置

###配置简介
启动前需要先了解 mlog-club 的配置项，mlog-club 的示例配置文件为`bbs-go.example.yaml`，文件在server目录中，请详细看下该文件：
```yaml

Env: prod # 环境，线上环境：prod、测试环境：dev
BaseUrl: https://mlog.club # 网站域名
SiteTitle: M-LOG # 网站标题
Port: '8082' # 端口
ShowSql: false # 是否打印sql
ViewsPath: "./web/views" # views模版文件目录，可配置绝对路径
StaticPath: "./web/static" # 静态文件目录，可配置绝对路径

MySqlUrl: username:password@tcp(localhost:3306)/mlog_db?charset=utf8&parseTime=True&loc=Local  # 数据库链接
RedisAddr: 127.0.0.1:6379 # redis链接

# oauth服务端配置
OauthServer:
  AuthUrl: https://mlog.club/oauth/authorize
  TokenUrl: https://mlog.club/oauth/token

# oauth客户端配置
OauthClient:
  ClientId: xxx
  ClientSecret: xxx
  ClientRedirectUrl: https://mlog.club/oauth/client/callback
  ClientSuccessUrl: https://admin.mlog.club/mlog/login_success.html

# github登录配置
Github:
  ClientID:
  ClientSecret:

# 阿里云oss配置
AliyunOss:
  Host: oss-cn.aliyuncs.com
  Bucket: bucket-name
  Endpoint: xx
  AccessId: xx
  AccessSecret: xx

# 邮件服务器配置
Smtp:
  Addr: smtp.qq.com
  Port: '25'
  Username: xxx
  Password: xxx
```

#### 数据库配置
mlog-club使用的`gorm`打开了`AutoMigrate`功能系统会在启动的时候自动根据我们定义的实体类来初始化表结构，所以我们要做的就是正确创建和配置数据库，建表、建索引功能交个`gorm`即可。

建表后数据初始化：

```sql
-- 初始化用户（用户名：admin、密码：123456）
INSERT INTO `t_user`(`id`, `username`, `nickname`, `avatar`, `email`, `password`, `status`, `create_time`, `update_time`, `roles`, `type`, `description`) VALUES (1, 'admin', '管理员', '', '', '$2a$10$ofA39bAFMpYpIX/Xiz7jtOMH9JnPvYfPRlzHXqAtLPFpbE/cLdjmS', 0, 1555419028975, 1555419028975, '管理员', 0, '轻轻地我走了，正如我轻轻的来。');
```

#### Github 登录配置
首先前往 Github 新建一个`Oauth Application`，填写`Application Name`、`Homepage URL`和`Authorization callback URL`；

`Authorization callback URL`为：https://yourhost.com//user/github/callback， 例如 https://mlog.club 的配置 callback url 为：https://mlog.club/user/github/callback

然后复制Oauth Application的 ClientID 和 ClientSecret 到我们的配置文件中的 Github 对应的配置中。

#### 阿里云 Oss 配置
mlog-club 目前使用阿里云的 oss 来处理图片上传，所以这里需要配置一下阿里云的 oss，阿里云的 oss 目前需要付费开通，后期考虑支持更多的图片上传服务商。

#### Smtp 邮件服务器配置
TODO 因为目前没有应用场景，所以先不用配置，后面会加上邮箱验证等功能，到时候就需要改配置了。

### 启动项目

```shell
go run main.go
```

## 问题反馈
- 欢迎交流：[https://mlog.club/topics](https://mlog.club/topics)
- 提交建议：[https://mlog.club/topic/609](https://mlog.club/topic/609)
