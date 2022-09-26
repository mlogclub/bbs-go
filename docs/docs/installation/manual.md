# 手动安装

## 概要

本系统为前后端分离设计，共分为三个模块：

- server：go语言开发后端接口服务，为系统提供数据支撑
- site：基于Nuxt.js开发的社区前台UI服务
- admin：基于Vue.js、element-ui开发的运营后台

在安装时需要分别安装以上三个模块，接下来我们分别详细介绍下三个模块的安装方法。

## server 端

### Go语言环境安装

首先确保你已经成功在本地安装好了go语言开发环境，如果环境不会安装的，请参照下面文章：

- [Go简易教程](https://mlog.club/topic/231)
- [go mod 使用帮助](https://mlog.club/topic/617)
- [Go语言配置go mod代理](https://mlog.club/topic/618)

### 数据库初始化

bbs-go数据库使用的是MySQL，MySQL最低版本要求为5.7（或者MariaDB 10及以上版本）。数据库的安装这里就不做赘述了，如果不会可自行Google、百度。

数据库成功安装之后，请通过工具（工具自行选择，我喜欢用mysql-cli，也可以用Navicat、TablePlus、MySQLWorkbench...）连接数据库，然后使用以下Sql脚本初始化`bbs-go`必须需的数据库表：

<div style="background-color: #F7EABC; padding: 20px;">
 <div><strong>⚠️注意⚠️</strong></div>
 <div>一、Sql脚本会初始化默认管理员用户，用户名：<span style="color:red;">admin</span>，密码：<span style="color:red;">123456</span></div>
 <div>二、该Sql脚本只会创建启动时必须的表，bbs-go系统使用的其他表会在系统正确启动后自动创建；</div>
</div>

```sql
CREATE DATABASE IF NOT EXISTS `bbsgo_db` DEFAULT CHARACTER SET utf8mb4;

USE bbsgo_db;
SET NAMES utf8mb4;

-- 初始化用户表
CREATE TABLE `t_user`
(
    `id`                 bigint(20) NOT NULL AUTO_INCREMENT,
    `username`           varchar(32)         DEFAULT NULL,
    `email`              varchar(128)        DEFAULT NULL,
    `email_verified`     tinyint(1) NOT NULL DEFAULT '0',
    `nickname`           varchar(16)         DEFAULT NULL,
    `avatar`             text,
    `background_image`   text,
    `password`           varchar(512)        DEFAULT NULL,
    `home_page`          varchar(1024)       DEFAULT NULL,
    `description`        text,
    `score`              bigint(20) NOT NULL,
    `status`             bigint(20) NOT NULL,
    `topic_count`        bigint(20) NOT NULL,
    `comment_count`      bigint(20) NOT NULL,
    `roles`              text,
    `forbidden_end_time` bigint(20) NOT NULL DEFAULT '0',
    `create_time`        bigint(20)          DEFAULT NULL,
    `update_time`        bigint(20)          DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `username` (`username`),
    UNIQUE KEY `email` (`email`),
    KEY `idx_user_score` (`score`),
    KEY `idx_user_status` (`status`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

-- 初始化用户数据（用户名：admin、密码：123456）
INSERT INTO t_user (`id`, `username`, `nickname`, `avatar`, `email`, `password`, `status`, `create_time`, `update_time`,
                    `roles`, `description`, `topic_count`, `comment_count`, `score`)
SELECT 1,
       'admin',
       'bbsgo站长',
       '',
       'a@example.com',
       '$2a$10$ofA39bAFMpYpIX/Xiz7jtOMH9JnPvYfPRlzHXqAtLPFpbE/cLdjmS',
       0,
       (UNIX_TIMESTAMP(now()) * 1000),
       (UNIX_TIMESTAMP(now()) * 1000),
       'owner',
       '轻轻地我走了，正如我轻轻的来。',
       0,
       0,
       0
FROM DUAL
WHERE NOT EXISTS(SELECT * FROM `t_user` WHERE `id` = 1);


-- 初始化话题节点
CREATE TABLE `t_topic_node`
(
    `id`          bigint(20) NOT NULL AUTO_INCREMENT,
    `name`        varchar(32) DEFAULT NULL,
    `description` longtext,
    `sort_no`     bigint(20)  DEFAULT NULL,
    `status`      bigint(20) NOT NULL,
    `create_time` bigint(20)  DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `name` (`name`),
    KEY `idx_sort_no` (`sort_no`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

INSERT INTO `t_topic_node` (`id`, `name`, `description`, `sort_no`, `status`, `create_time`)
SELECT 1, '默认节点', '', 0, 0, (UNIX_TIMESTAMP(now()) * 1000)
FROM DUAL
WHERE NOT EXISTS(SELECT * FROM `t_topic_node` WHERE `id` = 1);

-- 初始化系统配置表
CREATE TABLE `t_sys_config`
(
    `id`          bigint(20)   NOT NULL AUTO_INCREMENT,
    `key`         varchar(128) NOT NULL,
    `value`       text,
    `name`        varchar(32)  NOT NULL,
    `description` varchar(128) DEFAULT NULL,
    `create_time` bigint(20)   NOT NULL,
    `update_time` bigint(20)   NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `key` (`key`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

-- 初始化系统配置数据
INSERT INTO t_sys_config(`key`, `value`, `name`, `description`, `create_time`, `update_time`)
SELECT 'siteTitle',
       'bbs-go演示站',
       '站点标题',
       '站点标题',
       (UNIX_TIMESTAMP(now()) * 1000),
       (UNIX_TIMESTAMP(now()) * 1000)
FROM DUAL
WHERE NOT EXISTS(SELECT * FROM `t_sys_config` WHERE `key` = 'siteTitle');

INSERT INTO t_sys_config (`key`, `value`, `name`, `description`, `create_time`, `update_time`)
SELECT 'siteDescription',
       'bbs-go，基于Go语言的开源社区系统',
       '站点描述',
       '站点描述',
       (UNIX_TIMESTAMP(now()) * 1000),
       (UNIX_TIMESTAMP(now()) * 1000)
FROM DUAL
WHERE NOT EXISTS(SELECT * FROM `t_sys_config` WHERE `key` = 'siteDescription');

INSERT INTO t_sys_config (`key`, `value`, `name`, `description`, `create_time`, `update_time`)
SELECT 'siteKeywords',
       'bbs-go',
       '站点关键字',
       '站点关键字',
       (UNIX_TIMESTAMP(now()) * 1000),
       (UNIX_TIMESTAMP(now()) * 1000)
FROM DUAL
WHERE NOT EXISTS(SELECT * FROM `t_sys_config` WHERE `key` = 'siteKeywords');

INSERT INTO t_sys_config (`key`, `value`, `name`, `description`, `create_time`, `update_time`)
SELECT 'siteNavs',
       '[{\"title\":\"首页\",\"url\":\"/\"},{\"title\":\"话题\",\"url\":\"/topics\"},{\"title\":\"文章\",\"url\":\"/articles\"}]',
       '站点导航',
       '站点导航',
       (UNIX_TIMESTAMP(now()) * 1000),
       (UNIX_TIMESTAMP(now()) * 1000)
FROM DUAL
WHERE NOT EXISTS(SELECT * FROM `t_sys_config` WHERE `key` = 'siteNavs');

INSERT INTO t_sys_config (`key`, `value`, `name`, `description`, `create_time`, `update_time`)
SELECT 'defaultNodeId',
       '1',
       '默认节点',
       '默认节点',
       (UNIX_TIMESTAMP(now()) * 1000),
       (UNIX_TIMESTAMP(now()) * 1000)
FROM DUAL
WHERE NOT EXISTS(SELECT * FROM `t_sys_config` WHERE `key` = 'defaultNodeId');

INSERT INTO t_sys_config (`key`, `value`, `name`, `description`, `create_time`, `update_time`)
SELECT 'tokenExpireDays',
       '365',
       '用户登录有效期(天)',
       '用户登录有效期(天)',
       (UNIX_TIMESTAMP(now()) * 1000),
       (UNIX_TIMESTAMP(now()) * 1000)
FROM DUAL
WHERE NOT EXISTS(SELECT * FROM `t_sys_config` WHERE `key` = 'tokenExpireDays');

INSERT INTO t_sys_config (`key`, `value`, `name`, `description`, `create_time`, `update_time`)
SELECT 'scoreConfig',
       '{"postTopicScore":1,"postCommentScore":1,"checkInScore":1}',
       '积分配置',
       '积分配置',
       (UNIX_TIMESTAMP(now()) * 1000),
       (UNIX_TIMESTAMP(now()) * 1000)
FROM DUAL
WHERE NOT EXISTS(SELECT * FROM `t_sys_config` WHERE `key` = 'scoreConfig');
```

### 编译

Go语言是支持交叉编译的，你可以在Windows系统中将他编译成可以在Linux系统上可执行的二进制文件，其他平台同理。

编译时请打开命令行，并cd进入`bbs-go/server`目录，然后按照下面步骤进行编译。

- 如果只想编译当前机器可执行文件，那么只需要执行以下命令即可:

```bash
go build
```

- 如果想编译其他类型系统可执行二进制文件（例如在Windows系统中编译Linux系统上的可执行文件），可以执行以下命令：

```bash
# 1. 设置系统类型
SET GOOS=linux
# 2. 设置CPU架构
SET GOARCH=amd64
# 3. 编译
go build
```

注意在其他类Unix系统（如Linux、MacOS）中设置系统类型和CPU架构的方式和在Windows系统中不太一样，在类Unix系统中实现上面同样效果，命令如下：

```bash
GOOS=linux GOARCH=amd64 go build
```

编译后会在`bbs-go/server`目录下生成一个`bbs-go`文件（Windows下是`bbs-go.exe`），这个就是我们编译后的可执行文件。

### 配置

系统提供了配置示例，首先将配置示例文件`bbs-go/server/bbs-go.example.yaml`复制一份，重命名为：`bbs-go/server/bbs-go.yaml`。 示例文件中对于每个配置项都做了详细注释说明，请按照注释配置即可。

这里着重讲解一下数据库的配置，数据库为必要配置，成功配置数据库之后即可正常启动bbs-go server端。在上面步骤中我们已经将bbs-go的数据库初始化完成，数据库配置有以下几点：

1. 数据库用户名，这里假设为：root（这里只是假设，请根据自己的实际情况设置）
2. 数据库密码，这里假设为：123456（这里只是假设，请根据自己的实际情况设置）
3. 数据库服务器地址，这里假设为：localhost（这里只是假设，请根据自己的实际情况设置）
4. 数据库名，我们上面脚本创建的数据库名称为：bbsgo_db（这里只是假设，请根据自己的实际情况设置）

那么我们需要将`bbs-go.yaml`配置文件中`MySqlUrl`配置修改如下：

```yaml
MySqlUrl: root:123456@tcp(localhost:3306)/bbsgo_db?charset=utf8mb4&parseTime=True&loc=Local
```

### 运行

bbs-go的运行只需要两个文件，也就是我们上面步骤涉及到的两个文件：

1. 编译好的可执行文件：`bbs-go`（Windows下是`bbs-go.exe`）
2. 配置文件：`bbs-go.yaml`

将这两个文件放到同一个目录下，然后使用使用命令行进入到该目录中，在该目录下执行命令：

```bash
./bbs-go
```

即可启动`bbs-go-server`服务，启动成功后，控制台会输入如下日志：

```
Now listening on: http://localhost:8082
Application started. Press CMD+C to shut down.
```

bbs-go-server默认的端口是：8082（你也可以在bbs-go.yaml配置文件中自行修改），在浏览器中访问：`http://localhost:8082` 就可以看到效果了，浏览器会显示：Powered by bbs-go 。

## site

### nodejs环境安装

> **推荐nodejs 版本：v16.xx**

site 模块是基于nodejs开发的，所以编译他首先要安装nodejs环境。

- Windows下安装请参照这篇文章：https://www.cnblogs.com/liuqiyun/p/8133904.html
- Linux/MacOS下安装可以使用nvm：https://github.com/nvm-sh/nvm
- 或参照官网提供的《通过包管理器方式安装 Node.js》教程：https://nodejs.org/zh-cn/download/package-manager/

使用以下命令验证nodejs是否安装成功

```bash
➜  ~ node -v
v16.15.0
➜  ~ npm -v
8.1.0
```

### 安装依赖

在环境安装好后，进入到`bbs-go/site`目录，然后执行以下命令安装bbs-go-admin模块所需依赖：

```bash
npm install
```

npm的软件源是在国外服务器的，安装起来可能比较慢，你也可以使用`cnpm`来安装依赖。首先要安装`cnpm`，安装命令如下：

```bash
npm install cnpm -g --registry=https://r.npm.taobao.org
```

`cnpm`安装完后，在`bbs-go/site`目录下执行以下命令安装依赖：

```bash
cnpm install
```

### 配置

site模块的配置很简单，只需要配置：server端服务地址即可。这里的server端就是指bbs-go-server运行的地址，也就是上面讲到的：http://localhost:8082（这里要根据你的具体情况而定）。配置方式是，打开`bbs-go/site/nuxt.config.js`，找到`proxy`配置项，将他修改为下面的配置即可：

```js
proxy: {
    '/api/': 'http://localhost:8082'
},
```

### 运行

确认正确安装依赖，配置修改成功后，在site目录使用以下命令已开发模式启动bbs-go-site项目：

```bash
npm run dev
```

bbs-go-site服务默认端口为`3000`，启动成功后你就可以在浏览器通过：[http://localhost:3000](http://localhost:3000)访问和体验整个`bbs-go`的功能啦。

### 打包

我们线上部署是不能使用：`npm run dev`方式的，该方式为开发者模式启动，系统限制最多只能同时打开三个网页。

线上部署方式如下：

1. 使用`npm run build`编译site模块，编译成功后会在目录中生成`.nuxt`目录。
2. 使用`npm run start`启动服务，服务同样启动在：3000 端口。

## admin

### nodejs环境安装

> admin端和site端都是基于nodejs进行开发的，nodejs环境安装参见site端相关文档描述。

### 安装依赖

在环境安装好后，进入到`bbs-go/admin`目录，然后执行以下命令安装bbs-go-admin模块所需依赖：

```bash
npm install
```

npm的软件源是在国外服务器的，安装起来可能比较慢，你也可以使用`cnpm`来安装依赖。首先要安装`cnpm`，安装命令如下：

```bash
npm install cnpm -g --registry=https://r.npm.taobao.org
```

`cnpm`安装完后，在`bbs-go/admin`目录下执行以下命令安装依赖：

```bash
cnpm install
```

### 配置

admin 模块的配置文件分为开发环境和生产环境，文件分别为：

```
/admin/.env.development
/admin/.env.production
```

在使用开发模式运行服务时使用`/admin/.env.development`配置文件，在打包时使用`/admin/.env.production`配置文件。

配置文件中主要有两个配置项：

```
# 接口请求地址HOST，用于admin模块请求服务端接口
# 该配置的值一般设置为server端的HOST，或者site端的HOST（因为site端代理了server端的所有接口）
VUE_APP_BASE_API = 'http://localhost:8082'

# site模块访问根目录，作用：例如后台点击帖子标题时，能够正确跳转到帖子site端的访问路径
VUE_APP_BASE_URL = 'http://localhost:3000'
```

### 运行

确保依赖安装成功、配置正确后，可以使用命令：

```bash
npm run serve
```

来启动admin服务，使用该命令启动服务，会使用配置：`/admin/.env.development`

### 打包

确保依赖安装成功、配置正确后，可以使用命令：

```bash
npm run build
```

来打包`bbs-go-admin`模块，打包时会使用配置：`/admin/.env.production`

打包的成果为：`/admin/dist/`文件夹，将该文件夹部署到nginx或者其他web容器中即可正常访问。
