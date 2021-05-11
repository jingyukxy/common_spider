create database if not exists cluster_vod default charset utf8mb4;
use cluster_vod;

drop table if exists t_vod_staff;
create table t_vod_staff
(
    id          bigint(20)    not null auto_increment primary key comment 'id',
    name        varchar(50)   not null default '' comment '名称',
    staff_type  tinyint(1)    not null default 0 comment '类型 演员 导演 编剧',
    gender      tinyint(1)    not null default 0 comment '性别',
    image       varchar(255)  not null default '' comment '图片',
    height      tinyint(3)    not null default 0 comment '身高',
    country     varchar(255)  not null default '' comment '国家',
    description varchar(1000) not null default '' comment '描述',
    update_time timestamp     not null default now() comment '更新时间'
) charset = utf8mb4,
  engine = InnoDB, comment ='演职人员';

drop table if exists t_vod_episodes;
create table t_vod_episodes
(
    id          bigint(20)   not null auto_increment primary key comment 'id',
    vod_id      bigint(20)   not null default 0 comment 'vod id',
    title       varchar(255) not null default '' comment '标题',
    description text         not null comment '内容',
    keywords    varchar(255) not null default '' comment '关键词',
    create_time timestamp    not null default now() comment '创建时间'
) charset = utf8mb4,
  engine = InnoDB, comment ='剧情';

drop table if exists t_vod;
create table t_vod
(
    id              bigint(20)    not null auto_increment primary key comment 'id',
    name            varchar(255)  not null default '' comment '名称',
    category_id     bigint(20)    not null default 0 comment '类型',
    title           varchar(255)  not null default '' comment '副标题',
    keywords        varchar(255)  not null default '' comment '关键词',
    image           varchar(255)  not null default '' comment '图片',
    area            varchar(255)  not null default '' comment '国家',
    language        varchar(255)  not null default '' comment '语言',
    show_time       timestamp     not null comment '上映时间',
    total_episode   int(10)       not null default 0 comment '总集数',
    new_episode     int(10)       not null default 0 comment '更新至',
    is_end          tinyint(1)    not null default 0 comment '是否已经完结',
    show_year       smallint(5)   not null default 0 comment '上映年',
    add_time        timestamp     not null default now() comment '添加时间',
    hit_day_count   bigint(20)    not null default 0 comment '日点击',
    hit_week_count  bigint(20)    not null default 0 comment '周点击',
    hit_month_count bigint(20)    not null default 0 comment '月点击',
    stars           smallint(4)   not null default 0 comment '星级',
    status          tinyint(1)    not null default 0 comment '状态',
    editor          varchar(20)   not null default '' comment '编辑',
    letter          char(1)       not null default '' comment '首字母',
    score           float         not null default 0 comment '分数',
    douban_score    float         not null default 0 comment '豆瓣评分',
    is_film         tinyint(1)    not null default 0 comment '是否为电影',
    vod_length      int(10)       not null default 0 comment '长度 秒',
    description     varchar(2000) not null default '' comment '描述',
    update_time     timestamp     not null comment '更新时间'
) charset = utf8mb4,
  engine = InnoDB, comment ='视频';

drop table if exists t_vod_tag;
create table t_vod_tag
(
    id          bigint(20)   not null auto_increment primary key comment 'id',
    name        varchar(255) not null default '' comment '关键词',
    category_id bigint(20)   not null default 0 comment '分类'
) charset = utf8mb4,
  engine = InnoDB, comment ='视频词库';

drop table if exists t_vod_content;
create table t_vod_content
(
    id        bigint(20)   not null auto_increment primary key comment 'id',
    title     varchar(255) not null default '' comment '标题',
    play_type varchar(20)  not null default '' comment '播放类型',
    play_url  varchar(255) not null default '' comment '播放地址',
    vod_id    bigint(20)   not null default 0 comment '视频id'
) charset = utf8mb4,
  engine = InnoDB, comment ='视频内容';

drop table if exists t_vod_carousel;
create table t_vod_carousel
(
    id           bigint(20)   not null auto_increment primary key comment 'id',
    name         varchar(255) not null default '' comment 'name',
    description  varchar(200) not null default '' comment 'description',
    vod_id       bigint(20)   not null default 0 comment 'vod id',
    category_id  bigint(20)   not null default 0 comment '分类id',
    status       tinyint(1)   not null default 0 comment '状态',
    image        varchar(255) not null default '' comment '图片',
    external_url varchar(255) not null default '' comment '外部链接'
) charset = utf8mb4,
  engine = InnoDB, comment ='视频轮播';

drop table if exists t_vod_staff_associate;
create table t_vod_staff_associate
(
    vod_id   bigint(20) not null comment 'vod id',
    staff_id bigint(20) not null comment 'staff id',
    primary key idx_vsa_index (vod_id, staff_id)
) charset = utf8mb4,
  engine = InnoDB, comment ='角色关联表';

drop table if exists t_vod_category_binder;
create table t_vod_category_binder
(
    id           bigint(20)   not null auto_increment primary key comment 'id',
    spider_id    bigint(20)   not null default 0 comment '采集id',
    source_type  varchar(200) not null default '' comment '来源类型',
    category_ids varchar(200) not null default '' comment '绑定id,不同的app_id下的'
) charset = utf8mb4,
  engine = InnoDB, comment ='来源绑定表';

drop table if exists t_vod_category;
create table t_vod_category
(
    vod_id      bigint(20) not null default 0 comment 'vod id',
    category_id bigint(20) not null default 0 comment 'category id',
    app_id      bigint(20) not null default 0 comment 'app_id',
    primary key idx_vci (vod_id, category_id, app_id)
) charset = utf8mb4,
  engine = InnoDB, comment ='vod绑定表';
