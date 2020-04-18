rename table t_topic_like to t_user_like;

alter table t_user_like
    add column entity_type varchar(32) not null after user_id,
    change column topic_id entity_id bigint not null after entity_type,
    drop index idx_topic_like_user_id,
    drop index idx_topic_like_topic_id;

update t_user_like set entity_type = 'topic';

alter table t_topic
    drop column type,
    drop column image_list;
