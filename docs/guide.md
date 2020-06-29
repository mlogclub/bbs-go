---
sidebar: auto
---

# 使用帮助

## 公众号

> 欢迎关注公众号`码农俱乐部`获取更多干货资源。

![码农俱乐部](https://open.weixin.qq.com/qr/code?username=gh_950827012b8d)

## 简介

`bbs-go`是一个使用Go语言搭建的开源社区系统，采用前后端分离技术，Go语言提供api进行数据支撑，用户界面使用Nuxt.js进行渲染，后台界面基于element-ui。如果你正在学习Go语言，或者考虑转Go语言的Phper/Javaer...那么该项目对你有的学习会有很大的帮助，欢迎一起来交流。

主要功能如下：

- 用户中心
- 论坛功能
- 多人博客
- 站内消息
- 收藏功能
- 站内消息

## 项目地址

- Github：[https://github.com/mlogclub/bbs-go](https://github.com/mlogclub/bbs-go)
- 码云：[https://gitee.com/mlogclub/bbs-go](https://gitee.com/mlogclub/bbs-go)

## 演示

[https://mlog.club](https://mlog.club)

## 技术栈

- iris ([https://github.com/kataras/iris](https://github.com/kataras/iris)) Go语言 mvc 框架
- gorm ([http://gorm.io/](http://gorm.io/)) Go语言 orm 框架
- resty ([https://github.com/go-resty/resty](https://github.com/go-resty/resty)) Go语言好用的 http-client
- cron ([https://github.com/robfig/cron](https://github.com/robfig/cron)) 定时任务
- goquery ([https://github.com/PuerkitoBio/goquery](https://github.com/PuerkitoBio/goquery)) html dom 元素解析
- nuxt.js ([https://nuxtjs.org](https://nuxtjs.org)) 基于Vue的服务端渲染框架
- element-UI ([https://element.eleme.cn](https://element.eleme.cn)) 饿了么开源的基于 vue.js 的前端库
- vditor ([https://github.com/b3log/vditor](https://github.com/b3log/vditor)) Markdown 编辑器

## 获取源码

`bbs-go`的源码托管在Github：[https://github.com/mlogclub/bbs-go](https://github.com/mlogclub/bbs-go)，通过以下命令将源代码克隆到本地：

```bash
git clone https://github.com/mlogclub/mlog.git
```

## 项目结构

bbs-go采用前后端分离技术，网站和后台均使用`http api`进行数据通信。bbs-go包含三个模块：server、site，两个模块的介绍如下：

### server模块

`server`模块基于Go语言开发，他为整个项目提供接口数据支撑。`site`模块的数据都是从该模块获取的。

### site模块

`site`模块使用`nuxt.js`进行搭建，该模块是bbs-go的用户前端网页。`nuxt.js`相关知识可以去它的官网查看：[https://nuxtjs.org](https://nuxtjs.org)

## 配置详解

### server模块配置

`server`模块的示例配置文件为`server/bbs-go.example.yaml`，内容如下：

```yaml
Env: prod # 环境，线上环境：prod、测试环境：dev
BaseUrl: https://mlog.club # 网站域名
Port: '8082' # 端口
LogFile: /data/logs/bbs-go.log # 日志文件
ShowSql: false # 是否打印sql
StaticPath: /data/www  # 根路径下的静态文件目录，可配置绝对路径

# 数据库连接
MySqlUrl: username:password@tcp(localhost:3306)/bbsgo_db?charset=utf8mb4&parseTime=True&loc=Local

# github登录配置
Github:
  ClientID:
  ClientSecret:

# qq登录配置
QQConnect:
  AppId:
  AppKey:

# 阿里云oss配置
AliyunOss:
  Host: 请配置成你自己的
  Bucket: 请配置成你自己的
  Endpoint: 请配置成你自己的
  AccessId: 请配置成你自己的
  AccessSecret: 请配置成你自己的

# 邮件服务器配置，用于邮件通知
Smtp:
  Addr: smtp.qq.com
  Port: '25'
  Username: 请配置成你自己的
  Password: 请配置成你自己的

# 百度ai配置，用于自动分析文章摘要、标签
BaiduAi:
  ApiKey:
  SecretKey:
```

请复制该文件到：`server/bbs-go.yaml`，并根据配置文件中的注释将配置修改成你自己的。

### site模块配置

`site`模块是基于`nuxt.js`开发的，他的配置文件为：`site/nuxt.config.js`，我们主要关注一下两项配置即可：

1. port：site模块启动端口，默认为3000
2. proxy：`server`模块的连接地址，通过该地址可以请求`server`模块数据

## 快速启动

`bbs-go`总用有两个模块：server、site，接下来我们一步步的启动这二个模块。

### server模块启动

#### 安装依赖

server模块使用`go mod`管理依赖，如果你不清楚如何使用`go mod`，请先认真读一下下面两篇文章：

- [go mod使用帮助](https://mlog.club/topic/617)
- [配置go mod代理](https://mlog.club/topic/618)

在项目的`server`目录下执行下面命令来下载`server`模块依赖：

```bash
go mod download
```

#### 初始化数据库

新建数据库`bbsgo_db`(或者其他名字，你高兴就好)。并按照要求配置好你的数据库链接（请参见：[ server模块配置](#server模块配置)）。

配置好数据库链接后，`bbs-go`在启动的时候会自动建表，所以我们无需手动建表，但是有些数据是需要提前初始化的，例如：管理员用户，基本配置，所以我们需要执行下面sql脚本进行数据初始化：

```sql

-- 初始化用户表
CREATE TABLE IF NOT EXISTS `t_user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `username` varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  `email` varchar(128) COLLATE utf8_unicode_ci DEFAULT NULL,
  `nickname` varchar(16) COLLATE utf8_unicode_ci DEFAULT NULL,
  `avatar` text COLLATE utf8_unicode_ci,
  `password` varchar(512) COLLATE utf8_unicode_ci DEFAULT NULL,
  `status` int(11) NOT NULL,
  `roles` text COLLATE utf8_unicode_ci,
  `type` int(11) NOT NULL,
  `description` text COLLATE utf8_unicode_ci,
  `create_time` bigint(20) DEFAULT NULL,
  `update_time` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`),
  UNIQUE KEY `email` (`email`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- 初始化用户数据（用户名：admin、密码：123456）
INSERT INTO `t_user`(`id`, `username`, `nickname`, `avatar`, `email`, `password`, `status`, `create_time`, `update_time`, `roles`, `type`, `description`) VALUES (1, 'admin', '管理员', '', '', '$2a$10$ofA39bAFMpYpIX/Xiz7jtOMH9JnPvYfPRlzHXqAtLPFpbE/cLdjmS', 0, 1555419028975, 1555419028975, 'owner', 0, '轻轻地我走了，正如我轻轻的来。');

-- 初始化系统配置表
CREATE TABLE IF NOT EXISTS `t_sys_config` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `key` varchar(128) COLLATE utf8_unicode_ci NOT NULL,
  `value` text COLLATE utf8_unicode_ci,
  `name` varchar(32) COLLATE utf8_unicode_ci NOT NULL,
  `description` varchar(128) COLLATE utf8_unicode_ci DEFAULT NULL,
  `create_time` bigint(20) NOT NULL,
  `update_time` bigint(20) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `key` (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- 初始化系统配置数据
insert into t_sys_config(`key`, `value`, `name`, `description`, `create_time`, `update_time`) values
    ('siteTitle', 'bbs-go', '站点标题', '站点标题', 1555419028975, 1555419028975),
    ('siteDescription', 'bbs-go，基于Go语言的开源社区系统', '站点描述', '站点描述', 1555419028975, 1555419028975),
    ('siteKeywords', 'bbs-go', '站点关键字', '站点关键字', 1555419028975, 1555419028975),
    ('siteNavs', '[{\"title\":\"首页\",\"url\":\"/\"},{\"title\":\"话题\",\"url\":\"/topics\"},{\"title\":\"文章\",\"url\":\"/articles\"}]', '站点导航', '站点导航', 1555419028975, 1555419028975);
```

#### 配置server模块

参见：[ server模块配置](#server模块配置)

#### 启动server模块

再配置好数据库链接并初始化数据库之后，在server模块目录下执行下面脚本启动server模块：

```bash
go run main.go
```

### site模块启动

第一步：进入site模块目录，执行下面命令安装依赖：

```bash
npm install
```

第二步：打开`site/nuxt.config.js`进行相关配置，请参见：[site模块配置](#site模块配置)。

第三步：执行下面命令启动site模块服务：

```bash
npm run dev
```

正常启动后，打开 [http://127.0.0.1:8080](http://127.0.0.1:8080) 访问网站。

## Docker启动

感谢 [athom](https://github.com/athom) ，Docker启动功能由 [athom](https://github.com/athom) 提供支持，详见：[https://github.com/mlogclub/bbs-go/pull/25](https://github.com/mlogclub/bbs-go/pull/25)

使用Docker快速启动项目，首先参照[配置详解](#配置详解)，配置好你的项目，然后执行项目根目录下的：up.sh 启动服务。

## 问题反馈

- 欢迎交流：[https://mlog.club/topics](https://mlog.club/topics)
- 提交建议：[https://mlog.club/topic/609](https://mlog.club/topic/609)
