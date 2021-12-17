create table `users` (
    `id`       bigint(20)                             not null auto_increment,
    `user_id`  bigint(20)                             not null,
    `username` varchar(32) collate utf8mb4_general_ci not null,
    `password` varchar(128) collate utf8mb4_general_ci not null,
    `avatar`   varchar(32) collate utf8mb4_general_ci not null,
    primary key (`id`)
) engine = InnoDB
  default charset = utf8mb4
  collate = utf8mb4_general_ci;

create table `tetris` (
    `id` bigint(20) not null auto_increment,
    `user_id` bigint(20) not null ,
    `scores` varchar(500) collate utf8mb4_general_ci,
    primary key (`id`)
) engine = InnoDB default charset = utf8mb4 collate = utf8mb4_general_ci;

create table `mineSweep` (
    `id` bigint(20) not null auto_increment,
    `user_id` bigint(20) not null ,
    `scores` varchar(500) collate utf8mb4_general_ci,
    primary key (`id`)
) engine = InnoDB default charset = utf8mb4 collate = utf8mb4_general_ci;

create table `community` (
    `id` bigint(20) not null auto_increment,
    `page_id` bigint(20) not null ,
    `user_id` bigint(20) not null ,
    `views` int not null,
    `title` varchar(32) collate utf8mb4_general_ci not null ,
    `preContent` varchar(210) collate utf8mb4_general_ci not null ,
    `content` varchar(5000) collate utf8mb4_general_ci not null ,
    `approve` varchar(8000) collate utf8mb4_general_ci default '',
    `create_time` bigint(20) not null comment '创建时间',
    `update_time` bigint(20) comment '更新时间',
    primary key (`id`)
) engine = InnoDB default charset = utf8mb4 collate = utf8mb4_general_ci;

create table `page_id` (
    `id` bigint(20) not null auto_increment default 2,
    `user_id` bigint(20) not null ,
    `content` varchar(1000) collate utf8mb4_general_ci not null ,
    `approve` varchar(8000) collate utf8mb4_general_ci default '',
    `reply` varchar(8000) collate utf8mb4_general_ci default '',
    `create_time` bigint(20) not null comment '创建时间',
    `replyNum` int default 0 not null ,
    primary key (`id`)
) engine = InnoDB default charset = utf8mb4 collate = utf8mb4_general_ci;