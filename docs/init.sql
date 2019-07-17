-- mlog-club 使用的gorm打开了AutoMigrate功能系统会在启动的时候自动根据我们定义的实体类来初始化表结构，
-- 所以我们要做的就是正确创建和配置数据库，建表、建索引功能交个gorm即可。

-- 初始化oauth-client
INSERT INTO `t_oauth_client`(`id`, `client_id`, `client_secret`, `domain`, `callback_url`, `status`, `create_time`) VALUES (1, 'clientid请自定义', 'clientsecret请自定义', 'http://mlog.club(根据实际情况修改)', 'http://mlog.club(根据实际情况修改)', 0, '2019-06-21 14:02:02');
-- 初始化用户（用户名：admin、密码：123456）
INSERT INTO `t_user`(`id`, `username`, `nickname`, `avatar`, `email`, `password`, `status`, `create_time`, `update_time`, `roles`, `type`, `description`) VALUES (1, 'admin', '管理员', '', '', '$2a$10$ofA39bAFMpYpIX/Xiz7jtOMH9JnPvYfPRlzHXqAtLPFpbE/cLdjmS', 0, 1555419028975, 1555419028975, '管理员', 0, '轻轻地我走了，正如我轻轻的来。');

-- 初始化系统配置
insert into t_sys_config(`key`, `value`, `name`, `description`, `create_time`, `update_time`) values
    ('site.title', 'M-LOG', '站点标题', '站点标题', 1555419028975, 1555419028975),
    ('site.description', 'M-LOG社区，基于Go语言的开源社区系统', '站点描述', '站点描述', 1555419028975, 1555419028975),
    ('site.keywords', 'M-LOG,Go语言', '站点关键字', '站点关键字', 1555419028975, 1555419028975);
