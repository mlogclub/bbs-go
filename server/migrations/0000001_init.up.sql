-- 初始化角色数据
INSERT INTO `t_role` (`id`, `type`, `name`, `code`, `sort_no`, `remark`, `status`, `create_time`, `update_time`) VALUES
(1, 0, '超级管理员', 'owner', 0, '超级管理员', 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(2, 0, '管理员', 'admin', 1, '管理员', 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000));

-- 初始化节点数据
INSERT INTO `t_topic_node` (`id`, `name`, `description`, `logo`, `sort_no`, `status`, `create_time`) VALUES
(1, '默认节点', '', NULL, 0, 0, (UNIX_TIMESTAMP(now()) * 1000));

-- 初始化系统配置数据
INSERT INTO `t_sys_config` (`id`, `key`, `value`, `name`, `description`, `create_time`, `update_time`) VALUES
(1, 'siteTitle', 'bbs-go演示站', '站点标题', '站点标题', (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(2, 'siteDescription', 'bbs-go，基于Go语言的开源社区系统', '站点描述', '站点描述', (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(3, 'siteKeywords', '[\"bbs-go\"]', '站点关键字', '站点关键字', (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(4, 'siteNavs', '[{\"title\":\"首页\",\"url\":\"/\"},{\"title\":\"话题\",\"url\":\"/topics\"},{\"title\":\"文章\",\"url\":\"/articles\"}]', '站点导航', '站点导航', (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(5, 'defaultNodeId', '1', '默认节点', '默认节点', (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(6, 'tokenExpireDays', '365', '用户登录有效期(天)', '用户登录有效期(天)', (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(7, 'scoreConfig', '{\"postTopicScore\":1,\"postCommentScore\":1,\"checkInScore\":1}', '积分配置', '积分配置', (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(8, 'urlRedirect', 'false', '', '', (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(9, 'enableHideContent', 'false', '', '', (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(10, 'siteLogo', '', '', '', (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(11, 'siteNotification', '', '', '', (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(12, 'recommendTags', '', '', '', (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(13, 'modules', '{\"tweet\":true,\"topic\":true,\"article\":true}', '', '', (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000));

-- 初始化菜单数据
INSERT INTO `t_menu` (`id`, `parent_id`, `type`, `name`, `title`, `icon`, `path`, `component`, `sort_no`, `status`, `create_time`, `update_time`) VALUES
(1, 0, 'menu', 'Dashboard', '仪表盘', 'icon-dashboard', '/dashboard', 'dashboard/index', 0, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(2, 0, 'menu', 'User', '用户管理', 'icon-user', '/user', 'user/index', 1, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(4, 0, 'menu', 'Permission', '权限管理', 'icon-lock', '', NULL, 9, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(5, 4, 'menu', 'Role', '角色管理', '', '/permission/role', 'system/role/index', 10, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(6, 4, 'menu', 'Menu', '菜单管理', '', '/permission/menu', 'system/menu/index', 16, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(7, 4, 'menu', 'Api', '接口管理', '', '/permission/api', 'system/api/index', 13, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(8, 4, 'menu', 'Permission', '权限分配', '', '/permission/index', 'system/permission/index', 20, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(9, 2, 'func', '', '编辑', '', '', NULL, 2, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(10, 5, 'func', '', '新增', '', '', NULL, 12, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(11, 5, 'func', '', '编辑', '', '', NULL, 11, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(12, 6, 'func', '', '新增', '', '', NULL, 17, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(13, 6, 'func', '', '编辑', '', '', NULL, 18, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(14, 6, 'func', '', '排序', '', '', NULL, 19, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(15, 8, 'func', '', '保存', '', '', NULL, 21, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(16, 7, 'func', '', '新增', 'icon-settings', '', NULL, 15, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(17, 7, 'func', '', '编辑', 'icon-settings', '', NULL, 14, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(18, 0, 'menu', 'System', '系统管理', 'icon-settings', '', NULL, 22, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(19, 18, 'menu', 'Settings', '系统设置', '', '/system/settings', 'settings/index', 23, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(20, 18, 'menu', 'Dict', '字典管理', '', '/system/dict', 'system/dict/index', 24, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(21, 0, 'menu', 'Article', '文章管理', 'icon-file', '/article', 'article/index', 6, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(22, 0, 'menu', 'Forbidden-word', '违禁词', 'icon-safe', '/forbidden-word', 'forbidden-word/index', 7, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(23, 0, 'menu', 'Link', '友情链接', 'icon-link', '/link', 'link/index', 8, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(24, 0, 'menu', '', '帖子管理', 'icon-share-alt', '', '', 3, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(25, 24, 'menu', 'Topic', '帖子管理', '', '/topic', 'topic/index', 4, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(26, 24, 'menu', 'TopicNode', '节点管理', '', '/topic-node', 'topic-node/index', 5, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000));

-- 初始化角色菜单数据
INSERT INTO `t_role_menu` (`id`, `role_id`, `menu_id`, `create_time`) VALUES
(1, 1, 1, (UNIX_TIMESTAMP(now()) * 1000)),
(2, 1, 2, (UNIX_TIMESTAMP(now()) * 1000)),
(3, 1, 4, (UNIX_TIMESTAMP(now()) * 1000)),
(4, 1, 9, (UNIX_TIMESTAMP(now()) * 1000)),
(5, 1, 5, (UNIX_TIMESTAMP(now()) * 1000)),
(6, 1, 11, (UNIX_TIMESTAMP(now()) * 1000)),
(7, 1, 10, (UNIX_TIMESTAMP(now()) * 1000)),
(8, 1, 7, (UNIX_TIMESTAMP(now()) * 1000)),
(9, 1, 17, (UNIX_TIMESTAMP(now()) * 1000)),
(10, 1, 16, (UNIX_TIMESTAMP(now()) * 1000)),
(11, 1, 6, (UNIX_TIMESTAMP(now()) * 1000)),
(12, 1, 12, (UNIX_TIMESTAMP(now()) * 1000)),
(13, 1, 13, (UNIX_TIMESTAMP(now()) * 1000)),
(14, 1, 14, (UNIX_TIMESTAMP(now()) * 1000)),
(15, 1, 8, (UNIX_TIMESTAMP(now()) * 1000)),
(16, 1, 15, (UNIX_TIMESTAMP(now()) * 1000)),
(17, 1, 18, (UNIX_TIMESTAMP(now()) * 1000)),
(18, 1, 19, (UNIX_TIMESTAMP(now()) * 1000)),
(19, 1, 20, (UNIX_TIMESTAMP(now()) * 1000)),
(20, 1, 21, (UNIX_TIMESTAMP(now()) * 1000)),
(21, 1, 24, (UNIX_TIMESTAMP(now()) * 1000)),
(22, 1, 22, (UNIX_TIMESTAMP(now()) * 1000)),
(23, 1, 23, (UNIX_TIMESTAMP(now()) * 1000)),
(24, 1, 25, (UNIX_TIMESTAMP(now()) * 1000)),
(25, 1, 26, (UNIX_TIMESTAMP(now()) * 1000));