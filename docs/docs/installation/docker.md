# Docker 安装

## 概要

在使用docker安装之前，要求你能熟练使用docker，详见官方文档：[https://docs.docker.com](https://docs.docker.com)

本系统为前后端分离设计，共分为三个模块：

- server：go语言开发后端接口服务，为系统提供数据支撑
- site：基于Nuxt.js开发的社区前台UI服务
- admin：基于Vue.js、element-ui开发的运营后台

只有server、site模块提供Docker安装功能，因为admin模块不依赖与服务，成功变异之后可以直接运行，admin模块安装请参考本文档：[手动安装 -> admin端](/installation/manual.html#admin)

## 目录结构

在下载bbs-go源码后，在源码中提供了`docker-compose.yml`和各模块的`Dockerfile`，目录结构如下：

```text
.
├── .docker-compose
│   ├── mysql
│   │   ├── docker-entrypoint-initdb.d
│   │   │   └── init.sql   (数据库初始化脚本)
├── server
│   ├── Dockerfile   (server 模块Dockerfile)
│   ├── bbs-go.docker.yaml   (server 模块用于docker环境中的配置文件)
├── site
│   ├── Dockerfile   (site 模块Dockerfile)
│   ├── nuxt.config.docker.js   (site 模块用于docker环境中的配置文件)
└── ...
```

## 安装

## 配置

> TODO


<!--
### 构建镜像

docker服务成功安装且启动后，在项目根目录执行以下命令构建镜像：

> 构建时，请保证你的网速良好，因为会下载各种依赖

```bash
docker compose build
```
-->

### 启动服务

```bash
docker-compose pull
docker-compose up -d --no-build
```

启动成功后即可通过`3000`端口访问到你的服务了。

### 停止服务

```bash
docker compose stop
```
