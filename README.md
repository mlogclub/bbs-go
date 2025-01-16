## 项目地址

- 演示站：[https://bbs.bbs-go.com](https://bbs.bbs-go.com)
- 文档地址：[https://bbs-go.com](https://bbs-go.com)
- 问题反馈: [https://mlog.club/topics/node/3](https://mlog.club/topics/node/3)
- Github: [https://github.com/mlogclub/bbs-go](https://github.com/mlogclub/bbs-go)
- Gitee: [https://gitee.com/mlogclub/bbs-go](https://gitee.com/mlogclub/bbs-go)

## 联系我

### 用户交流群

![BBS-GO用户交流群](docs/images/qq.png)

### 加我微信

![BBS-GO用户交流群](docs/images/wechat.png)

## 介绍

`bbs-go` 是一个基于 Go 语言开发的开源社区论坛系统。它的设计旨在提供轻量、高效的社区讨论平台，支持现代化的 Web 技术栈，并且易于扩展和部署。bbs-go 项目采用模块化架构，能够与其他服务和前端框架无缝集成，适合各种规模的在线社区。

项目的主要特点包括：

- **高性能**：基于 Go 语言的并发特性，能够在高负载下保持良好的性能表现。
- **灵活性**：支持自定义配置、插件扩展，易于适应不同需求。
- **简单易用**：提供简洁的管理后台，方便社区管理员管理论坛内容和用户。
- **支持 MySQL 数据库**：提供对常见数据库的支持，确保数据存储的可靠性和稳定性。
- **响应式设计**：前端使用现代化的技术，能够在移动设备和桌面设备上提供良好的用户体验。

项目主要面向开发者和社区管理者，适合搭建技术讨论、兴趣分享等类型的社区论坛。

![bbs-go功能简介](docs/images/features.jpg)

## 模块

### server

[![bbs-go-server](https://github.com/mlogclub/bbs-go/actions/workflows/bbs-go-server.yml/badge.svg)](https://github.com/mlogclub/bbs-go/actions/workflows/bbs-go-server.yml)

> 基于`Golang`搭建，提供接口数据支撑。

技术栈

- iris ([https://github.com/kataras/iris](https://github.com/kataras/iris)) Go语言 mvc 框架
- gorm ([http://gorm.io](http://gorm.io)) 最好用的Go语言数据库orm框架
- resty ([https://github.com/go-resty/resty](https://github.com/go-resty/resty)) Go语言好用的 http-client
- cron ([https://github.com/robfig/cron](https://github.com/robfig/cron)) 定时任务框架
- goquery ([https://github.com/PuerkitoBio/goquery](https://github.com/PuerkitoBio/goquery)) html dom 元素解析

### site

[![bbs-go-site](https://github.com/mlogclub/bbs-go/actions/workflows/bbs-go-site.yml/badge.svg)](https://github.com/mlogclub/bbs-go/actions/workflows/bbs-go-site.yml)

> 前端页面渲染服务，基于`nuxt.js`搭建。

技术栈

- vue.js ([https://vuejs.org](https://vuejs.org)) 渐进式 JavaScript 框架
- nuxt.js ([https://nuxtjs.org](https://nuxtjs.org)) 基于Vue的服务端渲染框架，效率高到爆

### admin

[![bbs-go-admin](https://github.com/mlogclub/bbs-go/actions/workflows/bbs-go-admin.yml/badge.svg)](https://github.com/mlogclub/bbs-go/actions/workflows/bbs-go-admin.yml)

> 管理后台系统，基于`vue.js + element-ui`搭建。

技术栈

- vue.js ([https://vuejs.org](https://vuejs.org)) 渐进式 JavaScript 框架
- element-ui ([https://element.eleme.cn](https://element.eleme.cn)) 饿了么开源的基于 vue.js 的前端库

## 功能预览

![首页.png](https://s2.loli.net/2022/04/12/DpvPwB9dlQ6Chef.png)
![发帖.png](https://s2.loli.net/2022/04/12/KC8eXfE6sDLq34V.png)
![发动态.png](https://s2.loli.net/2022/04/12/14pMPuGjEU6kiWV.png)
![个人中心.png](https://s2.loli.net/2022/04/12/1PVNjMh9nUAXsl8.png)
![手机版.png](https://s2.loli.net/2022/04/12/mowWb78CGIaH6T2.png)
![后台首页.png](https://s2.loli.net/2022/04/12/ErX2BLTnh7ldz8D.png)
![后台配置.png](https://s2.loli.net/2022/04/12/PwK6aC74XEZlIOL.png)

## Contributors

<a href="https://github.com/mlogclub/bbs-go/graphs/contributors"><img src="https://opencollective.com/bbs-go/contributors.svg?width=890&button=false" /></a>
