-- 创建库
create database if not exists gininit;

-- 切换库
use gininit;

-- 用户表
create table if not exists `user`
(
  `id`          bigint auto_increment comment 'id' primary key,
  `account`     varchar(256)                not null comment '账号',
  `password`    varchar(512)                not null comment '密码',
  `union_id`    varchar(256)                null comment '微信开放平台id',
  `mp_open_id`  varchar(256)                null comment '公众号openId',
  `nick_name`   varchar(256)                null comment '用户昵称',
  `avatar`      varchar(1024)               null comment '用户头像',
  `profile`     varchar(512)                null comment '用户简介',
  `role`        varchar(256) default 'user' not null comment '用户角色：user/admin/ban',
  `create_time` bigint                      not null comment '创建时间',
  `update_time` bigint                      not null comment '更新时间',
  `is_delete`   tinyint      default 0      not null comment '是否删除',
  index idx_unionId (`union_id`)
) comment '用户' collate = utf8mb4_general_ci;
