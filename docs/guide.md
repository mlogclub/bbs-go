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

- iris ([https://github.com/kataras/iris](https://github.com/kataras/iris)) Go语言 mvc 框架
- gorm ([http://gorm.io/](http://gorm.io/)) Go语言 orm 框架
- resty ([https://github.com/go-resty/resty](https://github.com/go-resty/resty)) Go语言好用的 http-client
- cron ([https://github.com/robfig/cron](https://github.com/robfig/cron)) 定时任务
- goquery ([https://github.com/PuerkitoBio/goquery](https://github.com/PuerkitoBio/goquery)) html dom 元素解析
- nuxt.js ([https://nuxtjs.org](https://nuxtjs.org)) 基于Vue的服务端渲染框架
- element-UI ([https://element.eleme.cn](https://element.eleme.cn)) 饿了么开源的基于 vue.js 的前端库
- vditor ([https://github.com/b3log/vditor](https://github.com/b3log/vditor)) Markdown 编辑器

## 项目结构

bbs-go采用前后端分离技术，网站和后台均使用`http api`进行数据通信。所以bbs-go包含三个模块：server、site、admin。

### server模块

server模块是基于Go语言搭建的，为bbs-go提供数据接口支撑的服务。

### site模块

site模块使用`nuxt.js`进行搭建，该模块是bbs-go的用户前端网页。nuxt.js相关知识可以去它的官网查看：[https://nuxtjs.org](https://nuxtjs.org)

### admin模块

admin模块是bbs-go的管理后台，他基于element-ui搭建，element-ui相关知识可以去它的官网查看：[https://element.eleme.cn](https://element.eleme.cn/)

## 本地快速安装

> 说明：适用于本地开发和体验，各端运行后需要保持前台窗口进程。

### 安装依赖

```shell
# 第一步 clone 代码
git clone https://github.com/mlogclub/mlog.git

# 第二步 安装依赖
cd mlog
go mod tidy
```
> 说明  :bbs-go 的依赖是使用go mod来进行管理的，go mod使用帮助看这里：[https://mlog.club/topic/9](https://mlog.club/topic/9)

### 配置文件

在server目录中新建bbs-go.yaml配置文件（或者将bbs-go.example.yaml重命名)，配置内容请参考bbs-go.example.yaml中的说明。

> **注意：运行项目前先配置好数据库，否则程序无法运行。**

### 启动服务
在server目录中运行命令：
```shell
go run main.go
```

### 数据初始化

```sql
-- 初始化用户（用户名：admin、密码：123456）
INSERT INTO `t_user`(`id`, `username`, `nickname`, `avatar`, `email`, `password`, `status`, `create_time`, `update_time`, `roles`, `type`, `description`) VALUES (1, 'admin', '管理员', '', '', '$2a$10$ofA39bAFMpYpIX/Xiz7jtOMH9JnPvYfPRlzHXqAtLPFpbE/cLdjmS', 0, 1555419028975, 1555419028975, '管理员', 0, '轻轻地我走了，正如我轻轻的来。');

-- 初始化系统配置
insert into t_sys_config(`key`, `value`, `name`, `description`, `create_time`, `update_time`) values
    ('site.title', 'M-LOG', '站点标题', '站点标题', 1555419028975, 1555419028975),
    ('site.description', 'M-LOG社区，基于Go语言的开源社区系统', '站点描述', '站点描述', 1555419028975, 1555419028975),
    ('site.keywords', 'M-LOG,Go语言', '站点关键字', '站点关键字', 1555419028975, 1555419028975),
    ('recommend.tags', '', '推荐标签', '推荐标签，多个标签之间用英文逗号分隔', 1555419028975, 1555419028975);
```

### 启动网站前端
在site目录中运行命令：
```shell
npm install
npm run dev
```
正常启动后，打开 http://127.0.0.1:3000 访问网站。

### 启动管理后台

```shell
npm install
npm run serve
```
正常启动后，打开 http://127.0.0.1:8080 访问管理后台。

## 生产环境编译部署

编译安装

> TODO

Docker安装

> TODO

## 问题反馈

- 欢迎交流：[https://mlog.club/topics](https://mlog.club/topics)
- 提交建议：[https://mlog.club/topic/609](https://mlog.club/topic/609)