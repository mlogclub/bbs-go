alter table t_topic
add column ip_location varchar(64) null comment 'IP属地'
after ip;

alter table t_comment
add column ip_location varchar(64) null comment 'IP属地'
after ip;