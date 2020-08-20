CREATE DATABASE IF NOT EXISTS `bbsgo_db`;

USE bbsgo_db;
SET NAMES UTF8;

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
  KEY `idx_user_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

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
    `type`,
    `description`
) SELECT
    1,
    'admin',
    'bbs-go-owner',
    '',
    'a@example.com',
    '$2a$10$ofA39bAFMpYpIX/Xiz7jtOMH9JnPvYfPRlzHXqAtLPFpbE/cLdjmS',
    0,
    (UNIX_TIMESTAMP(now()) * 1000),
    (UNIX_TIMESTAMP(now()) * 1000),
    'owner',
    0,
    '轻轻地我走了，正如我轻轻的来。'
FROM
    DUAL
WHERE
    NOT EXISTS (
        SELECT
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
            `type`,
            `description`
        FROM t_user
        WHERE id = 1
    );

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
INSERT INTO t_sys_config(
    `key`,
    `value`,
    `name`,
    `description`,
    `create_time`,
    `update_time`
) SELECT
    'siteTitle',
    'bbs-go演示站',
    '站点标题',
    '站点标题',
    (UNIX_TIMESTAMP(now()) * 1000),
    (UNIX_TIMESTAMP(now()) * 1000)
FROM
    DUAL
WHERE
    NOT EXISTS (
        SELECT
            `key`,
            `value`,
            `name`,
            `description`,
            `create_time`,
            `update_time`
        FROM
            t_sys_config
        WHERE
            `key` = 'siteTitle'
    );

INSERT INTO t_sys_config (
    `key`,
    `value`,
    `name`,
    `description`,
    `create_time`,
    `update_time`
) SELECT
    'siteDescription',
    'bbs-go，基于Go语言的开源社区系统',
    '站点描述',
    '站点描述',
    (UNIX_TIMESTAMP(now()) * 1000),
    (UNIX_TIMESTAMP(now()) * 1000)
FROM
    DUAL
WHERE
    NOT EXISTS (
        SELECT
            `key`,
            `value`,
            `name`,
            `description`,
            `create_time`,
            `update_time`
        FROM
            t_sys_config
        WHERE
            `key` = 'siteDescription'
    );

INSERT INTO t_sys_config (
    `key`,
    `value`,
    `name`,
    `description`,
    `create_time`,
    `update_time`
) SELECT
    'siteKeywords',
    'bbs-go',
    '站点关键字',
    '站点关键字',
    (UNIX_TIMESTAMP(now()) * 1000),
    (UNIX_TIMESTAMP(now()) * 1000)
FROM
    DUAL
WHERE
    NOT EXISTS (
        SELECT
            `key`,
            `value`,
            `name`,
            `description`,
            `create_time`,
            `update_time`
        FROM
            t_sys_config
        WHERE
            `key` = 'siteKeywords'
    );


INSERT INTO t_sys_config (
    `key`,
    `value`,
    `name`,
    `description`,
    `create_time`,
    `update_time`
  )
SELECT
  'siteNavs',
  '[{\"title\":\"首页\",\"url\":\"/\"},{\"title\":\"话题\",\"url\":\"/topics\"},{\"title\":\"动态\",\"url\":\"/tweets\"},{\"title\":\"文章\",\"url\":\"/articles\"}]',
  '站点导航',
  '站点导航',
  (UNIX_TIMESTAMP(now()) * 1000),
  (UNIX_TIMESTAMP(now()) * 1000)
FROM DUAL
WHERE
  NOT EXISTS (
    SELECT
      `key`,
      `value`,
      `name`,
      `description`,
      `create_time`,
      `update_time`
    FROM t_sys_config
    WHERE
      `key` = 'siteNavs'
  );