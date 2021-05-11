-- 文章库
create database if not exists cluster_article default charset utf8mb4;
use cluster_article;

drop table if exists t_article;
create table t_article
(
    id               bigint(20)   not null auto_increment primary key comment 'id',
    category_id      bigint(20)   not null default 0 comment 'category id',
    title            varchar(255) not null default '' comment '标题',
    sub_title        varchar(255) not null default '' comment '副标题',
    thumb_img        varchar(255) not null default '' comment '缩略图',
    source_url       varchar(255) not null default '' comment '来源地址',
    author           varchar(50)  not null default '' comment '作者',
    tags             varchar(255) not null default '' comment '标签',
    description      varchar(255) not null default '' comment '描述',
    create_time      timestamp    not null default now() comment '创建时间',
    publish_time     timestamp    not null comment '发布时间',
    update_time      timestamp    not null comment '更新时间',
    status           tinyint(1)   not null default 0 comment '状态',
    is_top           tinyint(1)   not null default 0 comment '是否置顶',
    read_day_count   bigint(20)   not null default 0 comment '日阅读量',
    read_week_count  bigint(20)   not null default 0 comment '周阅读量',
    read_month_count bigint(20)   not null default 0 comment '月阅读量'
) charset = utf8mb4,
  engine = InnoDB, comment ='文章表';

drop table if exists t_article_content;
create table t_article_content
(
    id            bigint(20) not null auto_increment primary key comment 'id',
    article_id    bigint(20) not null default 0 comment '文章id',
    content       text       not null comment '内容',
    content_type  tinyint(1) not null default 0 comment '文章类型 文字 图片 文字加图片',
    allow_comment tinyint(1) not null default 0 comment '是否支持回复',
    up_count      bigint(20) not null default 0 comment '点赞量',
    down_count    bigint(20) not null default 0 comment '反对量',
    create_time   timestamp  not null default now() comment '创建时间'
) charset = utf8mb4,
  engine = InnoDB, comment ='文章内容';

