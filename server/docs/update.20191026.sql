alter table t_user drop index username;

CREATE TABLE `t_third_account` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) DEFAULT NULL,
  `avatar` varchar(1024) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `nickname` varchar(32) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `third_type` varchar(32) COLLATE utf8mb4_general_ci NOT NULL,
  `third_id` varchar(32) COLLATE utf8mb4_general_ci NOT NULL,
  `extra_data` longtext COLLATE utf8mb4_general_ci,
  `create_time` bigint(20) DEFAULT NULL,
  `update_time` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

alter table t_third_account add unique index idx_third(third_type, third_id);
alter table t_third_account add unique index idx_user_id_third_type(user_id, third_type);
