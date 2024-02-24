USE bbsgo_db;
SET NAMES utf8mb4;
-- 初始化用户表
CREATE TABLE `t_user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `username` varchar(32) DEFAULT NULL,
  `email` varchar(128) DEFAULT NULL,
  `email_verified` tinyint(1) NOT NULL DEFAULT '0',
  `nickname` varchar(16) DEFAULT NULL,
  `avatar` text,
  `gender` varchar(16) DEFAULT '',
  `birthday` datetime(3) DEFAULT NULL,
  `background_image` text,
  `password` varchar(512) DEFAULT NULL,
  `home_page` varchar(1024) DEFAULT NULL,
  `description` text,
  `score` int(11) NOT NULL,
  `status` int(11) NOT NULL,
  `topic_count` int(11) NOT NULL,
  `comment_count` int(11) NOT NULL,
  `follow_count` int(11) NOT NULL,
  `fans_count` int(11) NOT NULL,
  `roles` text,
  `forbidden_end_time` bigint(20) NOT NULL DEFAULT '0',
  `create_time` bigint(20) DEFAULT NULL,
  `update_time` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`),
  UNIQUE KEY `email` (`email`),
  KEY `idx_user_score` (`score`),
  KEY `idx_user_status` (`status`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;
-- 初始化用户数据（用户名：admin、密码：123456）
INSERT INTO t_user (
    `id`,
    `username`,
    `nickname`,
    `avatar`,
    `email`,
    `password`,
    `status`,
    `create_time`,
    `update_time`,
    `roles`,
    `description`,
    `topic_count`,
    `comment_count`,
    `score`,
    `follow_count`,
    `fans_count`
  )
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
  0,
  0,
  0
FROM DUAL
WHERE NOT EXISTS(
    SELECT *
    FROM `t_user`
    WHERE `id` = 1
  );
-- 初始化话题节点
CREATE TABLE `t_topic_node` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(32) DEFAULT NULL,
  `description` varchar(1024) DEFAULT NULL,
  `logo` varchar(1024) DEFAULT NULL,
  `sort_no` int(11) DEFAULT NULL,
  `status` int(11) NOT NULL,
  `create_time` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`),
  KEY `idx_sort_no` (`sort_no`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;
INSERT INTO `t_topic_node` (
    `id`,
    `name`,
    `description`,
    `sort_no`,
    `status`,
    `create_time`
  )
SELECT 1,
  '默认节点',
  '',
  0,
  0,
  (UNIX_TIMESTAMP(now()) * 1000)
FROM DUAL
WHERE NOT EXISTS(
    SELECT *
    FROM `t_topic_node`
    WHERE `id` = 1
  );
-- 初始化系统配置表
CREATE TABLE `t_sys_config` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `key` varchar(128) NOT NULL,
  `value` text,
  `name` varchar(32) NOT NULL,
  `description` varchar(128) DEFAULT NULL,
  `create_time` bigint(20) NOT NULL,
  `update_time` bigint(20) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `key` (`key`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;
-- 初始化系统配置数据
INSERT INTO t_sys_config(
    `key`,
    `value`,
    `name`,
    `description`,
    `create_time`,
    `update_time`
  )
SELECT 'siteTitle',
  'bbs-go演示站',
  '站点标题',
  '站点标题',
  (UNIX_TIMESTAMP(now()) * 1000),
  (UNIX_TIMESTAMP(now()) * 1000)
FROM DUAL
WHERE NOT EXISTS(
    SELECT *
    FROM `t_sys_config`
    WHERE `key` = 'siteTitle'
  );
INSERT INTO t_sys_config (
    `key`,
    `value`,
    `name`,
    `description`,
    `create_time`,
    `update_time`
  )
SELECT 'siteDescription',
  'bbs-go，基于Go语言的开源社区系统',
  '站点描述',
  '站点描述',
  (UNIX_TIMESTAMP(now()) * 1000),
  (UNIX_TIMESTAMP(now()) * 1000)
FROM DUAL
WHERE NOT EXISTS(
    SELECT *
    FROM `t_sys_config`
    WHERE `key` = 'siteDescription'
  );
INSERT INTO t_sys_config (
    `key`,
    `value`,
    `name`,
    `description`,
    `create_time`,
    `update_time`
  )
SELECT 'siteKeywords',
  '["bbs-go"]',
  '站点关键字',
  '站点关键字',
  (UNIX_TIMESTAMP(now()) * 1000),
  (UNIX_TIMESTAMP(now()) * 1000)
FROM DUAL
WHERE NOT EXISTS(
    SELECT *
    FROM `t_sys_config`
    WHERE `key` = 'siteKeywords'
  );
INSERT INTO t_sys_config (
    `key`,
    `value`,
    `name`,
    `description`,
    `create_time`,
    `update_time`
  )
SELECT 'siteNavs',
  '[{\"title\":\"首页\",\"url\":\"/\"},{\"title\":\"话题\",\"url\":\"/topics\"},{\"title\":\"文章\",\"url\":\"/articles\"}]',
  '站点导航',
  '站点导航',
  (UNIX_TIMESTAMP(now()) * 1000),
  (UNIX_TIMESTAMP(now()) * 1000)
FROM DUAL
WHERE NOT EXISTS(
    SELECT *
    FROM `t_sys_config`
    WHERE `key` = 'siteNavs'
  );
INSERT INTO t_sys_config (
    `key`,
    `value`,
    `name`,
    `description`,
    `create_time`,
    `update_time`
  )
SELECT 'defaultNodeId',
  '1',
  '默认节点',
  '默认节点',
  (UNIX_TIMESTAMP(now()) * 1000),
  (UNIX_TIMESTAMP(now()) * 1000)
FROM DUAL
WHERE NOT EXISTS(
    SELECT *
    FROM `t_sys_config`
    WHERE `key` = 'defaultNodeId'
  );
INSERT INTO t_sys_config (
    `key`,
    `value`,
    `name`,
    `description`,
    `create_time`,
    `update_time`
  )
SELECT 'tokenExpireDays',
  '365',
  '用户登录有效期(天)',
  '用户登录有效期(天)',
  (UNIX_TIMESTAMP(now()) * 1000),
  (UNIX_TIMESTAMP(now()) * 1000)
FROM DUAL
WHERE NOT EXISTS(
    SELECT *
    FROM `t_sys_config`
    WHERE `key` = 'tokenExpireDays'
  );
INSERT INTO t_sys_config (
    `key`,
    `value`,
    `name`,
    `description`,
    `create_time`,
    `update_time`
  )
SELECT 'scoreConfig',
  '{"postTopicScore":1,"postCommentScore":1,"checkInScore":1}',
  '积分配置',
  '积分配置',
  (UNIX_TIMESTAMP(now()) * 1000),
  (UNIX_TIMESTAMP(now()) * 1000)
FROM DUAL
WHERE NOT EXISTS(
    SELECT *
    FROM `t_sys_config`
    WHERE `key` = 'scoreConfig'
  );

-- 菜单配置
DROP TABLE IF EXISTS `t_menu`;
CREATE TABLE `t_menu` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `parent_id` bigint(20) DEFAULT NULL,
  `name` varchar(256) DEFAULT NULL,
  `title` varchar(64) DEFAULT NULL,
  `icon` varchar(1024) DEFAULT NULL,
  `path` varchar(1024) DEFAULT NULL,
  `sort_no` bigint(20) NOT NULL DEFAULT '0',
  `status` bigint(20) DEFAULT NULL,
  `create_time` bigint(20) NOT NULL DEFAULT '0',
  `update_time` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

INSERT INTO `t_menu` (`id`, `parent_id`, `title`, `name`, `icon`, `path`, `sort_no`, `status`, `create_time`, `update_time`) VALUES
(1, 0, '仪表盘', 'Dashboard', 'icon-dashboard', '/dashboard', 0, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(2, 0, '用户管理', 'User', 'icon-user', '/user', 1, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(3, 0, '帖子管理', '', 'icon-file', '', 2, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(4, 3, '节点管理', 'TopicNode', '', '/topic/topic-node', 3, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(5, 3, '帖子管理', 'Topic', '', '/topic/index', 4, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(6, 0, '文章管理', 'Article', 'icon-nav', '/article', 5, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(7, 0, '违禁词', 'ForbiddenWord', 'icon-stop', '/forbidden-word', 6, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(8, 0, '友情链接', 'Link', 'icon-link', '/link', 7, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(9, 0, '系统设置', 'Settings', 'icon-settings', '/settings', 8, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(10, 0, '权限管理', '', 'icon-lock', '', 9, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(11, 10, '角色管理', 'Role', '', '/permission/role', 10, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(12, 10, '菜单管理', 'Menu', '', '/permission/menu', 11, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(13, 10, '权限分配', 'Permission', '', '/permission/index', 12, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000));


-- 角色配置
DROP TABLE IF EXISTS `t_role`;
CREATE TABLE `t_role` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `type` bigint NOT NULL DEFAULT '1',
  `name` varchar(64) DEFAULT NULL,
  `code` varchar(64) DEFAULT NULL,
  `sort_no` bigint DEFAULT NULL,
  `remark` varchar(256) DEFAULT NULL,
  `status` bigint DEFAULT NULL,
  `create_time` bigint(20) NOT NULL DEFAULT '0',
  `update_time` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `code` (`code`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

INSERT INTO `t_role` (`id`, `type`, `name`, `code`, `sort_no`, `remark`, `status`, `create_time`, `update_time`) VALUES
(1, 0, '超级管理员', 'owner', 0, '超级管理员', 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(2, 0, '管理员', 'admin', 1, '管理员', 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000));

-- 用户角色
DROP TABLE IF EXISTS `t_user_role`;
CREATE TABLE `t_user_role` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint DEFAULT NULL,
  `role_id` bigint DEFAULT NULL,
  `create_time` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_role` (`user_id`,`role_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

INSERT INTO `t_user_role` (`id`, `user_id`, `role_id`, `create_time`) VALUES
(1, 1, 1, (UNIX_TIMESTAMP(now()) * 1000));


-- 角色菜单
DROP TABLE IF EXISTS `t_role_menu`;
CREATE TABLE `t_role_menu` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `role_id` bigint DEFAULT NULL,
  `menu_id` bigint DEFAULT NULL,
  `create_time` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_role_menu` (`role_id`,`menu_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

INSERT INTO `t_role_menu` (`id`, `role_id`, `menu_id`, `create_time`) VALUES
(1, 1, 1, (UNIX_TIMESTAMP(now()) * 1000)),
(2, 1, 2, (UNIX_TIMESTAMP(now()) * 1000)),
(3, 1, 3, (UNIX_TIMESTAMP(now()) * 1000)),
(4, 1, 6, (UNIX_TIMESTAMP(now()) * 1000)),
(5, 1, 7, (UNIX_TIMESTAMP(now()) * 1000)),
(6, 1, 8, (UNIX_TIMESTAMP(now()) * 1000)),
(7, 1, 9, (UNIX_TIMESTAMP(now()) * 1000)),
(8, 1, 4, (UNIX_TIMESTAMP(now()) * 1000)),
(9, 1, 5, (UNIX_TIMESTAMP(now()) * 1000)),
(10, 1, 10, (UNIX_TIMESTAMP(now()) * 1000)),
(11, 1, 11, (UNIX_TIMESTAMP(now()) * 1000)),
(12, 1, 12, (UNIX_TIMESTAMP(now()) * 1000)),
(13, 1, 13, (UNIX_TIMESTAMP(now()) * 1000)),
(14, 2, 1, (UNIX_TIMESTAMP(now()) * 1000)),
(15, 2, 2, (UNIX_TIMESTAMP(now()) * 1000)),
(16, 2, 3, (UNIX_TIMESTAMP(now()) * 1000)),
(17, 2, 6, (UNIX_TIMESTAMP(now()) * 1000)),
(18, 2, 7, (UNIX_TIMESTAMP(now()) * 1000)),
(19, 2, 8, (UNIX_TIMESTAMP(now()) * 1000)),
(20, 2, 4, (UNIX_TIMESTAMP(now()) * 1000)),
(21, 2, 5, (UNIX_TIMESTAMP(now()) * 1000)),
(22, 2, 9, (UNIX_TIMESTAMP(now()) * 1000));
