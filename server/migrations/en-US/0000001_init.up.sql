-- Initialize roles data
INSERT INTO `t_role` (`id`, `type`, `name`, `code`, `sort_no`, `remark`, `status`, `create_time`, `update_time`) VALUES
(1, 0, 'Owner', 'owner', 0, 'Owner with highest privileges', 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(2, 0, 'Admin', 'admin', 1, 'Admin with management privileges', 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000));

-- Initialize topic nodes data
INSERT INTO `t_topic_node` (`id`, `name`, `description`, `logo`, `sort_no`, `status`, `create_time`) VALUES
(1, 'Default', '', NULL, 0, 0, (UNIX_TIMESTAMP(now()) * 1000));

-- Initialize system configuration data
INSERT INTO `t_sys_config` (`id`, `key`, `value`, `name`, `description`, `create_time`, `update_time`) VALUES
(1, 'siteTitle', 'BBS-GO Demo Site', 'Site Title', 'Site Title', (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(2, 'siteDescription', 'BBS-GO, an open source community system based on Go language', 'Site Description', 'Site Description', (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(3, 'siteKeywords', '[\"bbs-go\"]', 'Site Keywords', 'Site Keywords', (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(4, 'siteNavs', '[{\"title\":\"Topics\",\"url\":\"/topics\"},{\"title\":\"Articles\",\"url\":\"/articles\"}]', 'Site Navigation', 'Site Navigation', (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(5, 'defaultNodeId', '1', 'Default Node', 'Default Node', (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(6, 'tokenExpireDays', '365', 'User Login Validity Period (Days)', 'User Login Validity Period (Days)', (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(7, 'scoreConfig', '{\"postTopicScore\":1,\"postCommentScore\":1,\"checkInScore\":1}', 'Score Configuration', 'Score Configuration', (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(8, 'urlRedirect', 'false', '', '', (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(9, 'enableHideContent', 'false', '', '', (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(10, 'siteLogo', '', '', '', (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(11, 'siteNotification', '', '', '', (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(12, 'recommendTags', '', '', '', (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(13, 'modules', '{\"tweet\":true,\"topic\":true,\"article\":true}', '', '', (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000));

-- Initialize menu data
INSERT INTO `t_menu` (`id`, `parent_id`, `type`, `name`, `title`, `icon`, `path`, `component`, `sort_no`, `status`, `create_time`, `update_time`) VALUES
(1, 0, 'menu', 'Dashboard', 'Dashboard', 'icon-dashboard', '/dashboard', 'dashboard/index', 0, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(2, 0, 'menu', 'User', 'User', 'icon-user', '/user', 'user/index', 1, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(4, 0, 'menu', 'Permission', 'Permission', 'icon-lock', '', NULL, 9, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(5, 4, 'menu', 'Role', 'Role', '', '/permission/role', 'system/role/index', 10, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(6, 4, 'menu', 'Menu', 'Menu', '', '/permission/menu', 'system/menu/index', 16, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(7, 4, 'menu', 'Api', 'API', '', '/permission/api', 'system/api/index', 13, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(8, 4, 'menu', 'Permission', 'Permission', '', '/permission/index', 'system/permission/index', 20, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(9, 2, 'func', '', 'Edit', '', '', NULL, 2, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(10, 5, 'func', '', 'Add', '', '', NULL, 12, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(11, 5, 'func', '', 'Edit', '', '', NULL, 11, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(12, 6, 'func', '', 'Add', '', '', NULL, 17, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(13, 6, 'func', '', 'Edit', '', '', NULL, 18, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(14, 6, 'func', '', 'Sort', '', '', NULL, 19, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(15, 8, 'func', '', 'Save', '', '', NULL, 21, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(16, 7, 'func', '', 'Add', 'icon-settings', '', NULL, 15, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(17, 7, 'func', '', 'Edit', 'icon-settings', '', NULL, 14, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(18, 0, 'menu', 'System', 'System', 'icon-settings', '', NULL, 22, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(19, 18, 'menu', 'Settings', 'Settings', '', '/system/settings', 'settings/index', 23, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(20, 18, 'menu', 'Dict', 'Dictionary', '', '/system/dict', 'system/dict/index', 24, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(21, 0, 'menu', 'Article', 'Article', 'icon-file', '/article', 'article/index', 6, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(22, 0, 'menu', 'Forbidden-word', 'Forbidden Words', 'icon-safe', '/forbidden-word', 'forbidden-word/index', 7, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(23, 0, 'menu', 'Link', 'Friend Links', 'icon-link', '/link', 'link/index', 8, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(24, 0, 'menu', '', 'Topic', 'icon-share-alt', '', '', 3, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(25, 24, 'menu', 'Topic', 'Topic', '', '/topic', 'topic/index', 4, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000)),
(26, 24, 'menu', 'TopicNode', 'Node', '', '/topic-node', 'topic-node/index', 5, 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000));

-- Initialize role menu data
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