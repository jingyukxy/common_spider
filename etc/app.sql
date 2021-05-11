-- 主库
create database if not exists cluster_app default charset utf8mb4;
use cluster_app;

drop table if exists t_app_type;
create table t_app_type
(
    id          bigint(20)   not null auto_increment primary key comment 'id',
    name        varchar(20)  not null default '' comment '名称',
    pid         bigint(20)   not null default 0 comment '父id',
    status      tinyint(1)   not null default 0 comment '状态1,启用 0 禁用',
    description varchar(255) not null default '' comment '描述',
    keywords    varchar(255) not null default '' comment '关键字',
    index idx_type_name (name)
) charset = utf8mb4,
  engine = InnoDB, comment ='app类型表';

drop table if exists t_app;
create table t_app
(
    id          bigint(20)   not null auto_increment primary key comment 'id',
    name        varchar(50)  not null default '' comment '名称',
    title       varchar(255) not null default '' comment '标题',
    type_id     bigint(20)   not null default 0 comment 'app类型id',
    keywords    varchar(255) not null default '' comment '应用关键字',
    has_mobile  tinyint(1)   not null default 0 comment '是否有移动端',
    has_search  tinyint(1)   not null default 0 comment '是否有搜索',
    status      tinyint(1)   not null default 0 comment '状态',
    update_time timestamp    not null default now() comment '更新时间',
    create_time timestamp    not null comment '创建时间',
    description varchar(255) not null default '' comment '描述'
) charset = utf8mb4,
  engine = InnoDB, comment ='应用表';

drop table if exists t_api_type;
create table t_api_type
(
    id          bigint(20)   not null auto_increment primary key comment 'id',
    pid         bigint(20)   not null default 0 comment '父id',
    type_name   varchar(20)  not null default '' comment '类型名称',
    description varchar(255) not null default '' comment '描述'
) charset = utf8mb4,
  engine = InnoDB, comment ='api类型表';

drop table if exists t_app_url;
create table t_app_url
(
    id          bigint(20)   not null auto_increment primary key comment 'id',
    app_id      bigint(20)   not null default 0 comment 'app id',
    name        varchar(50)  not null default '' comment '名称',
    url         varchar(255) not null default '' comment '链接',
    api_type_id bigint(20)   not null default 0 comment 'api类型',
    device_type tinyint(1)   not null default 0 comment '设备类型1.pc 2.h5 3.android 4.ios',
    status      tinyint(1)   not null default 0 comment '状态',
    theme       varchar(50)  not null default '' comment '模板主题',
    create_time timestamp    not null default now() comment '创建时间'
) charset = utf8mb4,
  engine = InnoDB, comment ='应用链接表';

drop table if exists t_ads;
create table t_ads
(
    id             bigint(20)   not null auto_increment primary key comment 'id',
    merchant_id    bigint(20)   not null default 0 comment '厂商id',
    position_id    bigint(20)   not null default 0 comment '版位id',
    template_id    bigint(20)   not null default 0 comment '模板id',
    cost_type_id   bigint(20)   not null default 0 comment '费用id',
    ads_type       tinyint(1)   not null default 0 comment '广告类型1.cpc 2.cpm 3.cpt 4 cps 5 cpa',
    image          varchar(255) not null default '' comment '广告图片',
    img_height     int(10)      not null default 0 comment '图片高度',
    img_width      int(10)      not null default 0 comment '图片宽度',
    is_flash       tinyint(1)   not null default 0 comment '是否是flash',
    merchant_desc  varchar(255) not null default '' comment '厂商描述',
    open_mode      tinyint(1)   not null default 0 comment '打开方式1.轮播2.右下角3.悬浮4.弹出',
    weight         smallint(5)  not null default 0 comment '权重',
    other_settings text         not null comment '其它设置',
    status         tinyint(1)   not null default 0 comment '状态',
    expiration     timestamp    not null comment '过期时间',
    start_time     timestamp    not null comment '开始时间',
    create_time    timestamp    not null default now() comment '创建时间'
) charset = utf8mb4,
  engine = InnoDB, comment ='广告表';

drop table if exists t_ads_merchant;
create table t_ads_merchant
(
    id          bigint(20)   not null auto_increment primary key comment 'id',
    name        varchar(50)  not null default '' comment '名称',
    company     varchar(255) not null default '' comment '公司名称',
    phone       varchar(20)  not null default '' comment '电话',
    qq          varchar(20)  not null default '' comment 'qq',
    description varchar(255) not null default '' comment '描述',
    link_name   varchar(255) not null default '' comment '联系人',
    create_time timestamp    not null default now() comment '创建时间'
) charset = utf8mb4,
  engine = InnoDB, comment ='广告商户表';

drop table if exists t_ads_position;
create table t_ads_position
(
    id            bigint(20)   not null auto_increment primary key comment 'id',
    name          varchar(30)  not null default '' comment '版位名称',
    description   varchar(255) not null default '' comment '版位描述',
    app_id        bigint(20)   not null default 0 comment 'app id',
    app_url_id    bigint(20)   not null default 0 comment 'app url id',
    height        smallint(5)  not null default 0 comment '高度',
    width         smallint(5)  not null default 0 comment '宽度',
    status        tinyint(1)   not null default 0 comment '状态',
    display_mode  tinyint(1)   not null default 0 comment '展示类型',
    carousel_time smallint(4)  not null default 0 comment '轮播时间，以秒为单位',
    create_time   timestamp    not null default now() comment '创建时间'
) charset = utf8mb4,
  engine = InnoDB, comment ='广告版位表';

drop table if exists t_ads_template;
create table t_ads_template
(
    id          bigint(20)  not null auto_increment primary key comment 'id',
    name        varchar(50) not null default '' comment '名称',
    js_code     text        not null comment 'js代码块',
    create_time timestamp   not null default now() comment '创建时间',
    ads_mode    tinyint(1)  not null default 0 comment '广告模式'
) charset = utf8mb4,
  engine = InnoDB, comment ='广告模板表';

drop table if exists t_ads_cost;
create table t_ads_cost
(
    id          bigint(20)   not null auto_increment primary key comment 'id',
    name        varchar(30)  not null default '' comment '计费名称',
    cost_type   tinyint(2)   not null default 0 comment '计费模式1.cpc 2cpm 3cpt 4cps 5cpa',
    price       float        not null default 0 comment '单价',
    description varchar(255) not null default '' comment '备注'
) charset = utf8mb4,
  engine = InnoDB, comment ='广告计费表';

drop table if exists t_external_link;
create table t_external_link
(
    id          bigint(20)   not null auto_increment primary key comment 'id',
    link_name   varchar(100) not null default '' comment '链接名称',
    app_id      bigint(20)   not null default 0 comment 'app id',
    app_url_id  bigint(20)   not null default 0 comment 'app url id',
    link_url    varchar(255) not null default '' comment '链接url',
    link_type   tinyint(1)   not null default 0 comment '1.文字 2.图片 3.文字加图片',
    logo        varchar(255) not null default '' comment 'logo 地址',
    start_time  timestamp    not null comment '开始时间',
    expiration  timestamp    not null comment '过期时间',
    create_time timestamp    not null default now() comment '创建时间'
) charset = utf8mb4,
  engine = InnoDB, comment ='广告计费表';

drop table if exists t_app_user;
create table t_app_user
(
    id              bigint(20)   not null auto_increment primary key comment 'id',
    username        varchar(50)  not null default '' comment '注册用户名',
    app_username    varchar(100) not null default '' comment 'app用户名加id前缀',
    phone           varchar(20)  not null default '' comment '手机',
    email           varchar(255) not null default '' comment '邮箱',
    avatar          varchar(255) not null default '' comment '头像',
    app_id          bigint(20)   not null default 0 comment 'app id',
    salt            varchar(20)  not null default '' comment '盐',
    login_count     int(10)      not null default 0 comment '登录次数',
    is_third_party  tinyint(1)   not null default 0 comment '是否为第三方用户',
    open_type       tinyint(1)   not null default 0 comment 'qq weixin weibo.....',
    last_login_time timestamp    not null comment '最后登录时间',
    status          tinyint(1)   not null default 0 comment '用户状态',
    create_time     timestamp    not null default now() comment '创建时间',
    update_time     timestamp    not null comment '更新时间',
    index idx_app_user_name (app_username)
) charset = utf8mb4,
  engine = InnoDB, comment ='app用户表';

drop table if exists t_user_profile;
create table t_user_profile
(
    id         bigint(20) not null auto_increment primary key comment 'id',
    uid        bigint(20) not null default 0 comment 'uid',
    birth      datetime   not null default '1970-01-01 00:00:00' comment '生日',
    user_level tinyint(1) not null default 0 comment '用户级别',
    gender     tinyint(1) not null default 0 comment '性别',
    score      int(10)    not null default 0 comment '积分',
    is_vip     tinyint(1) not null default 0 comment '是否是vip'
) charset = utf8mb4,
  engine = InnoDB, comment ='app用户表';

drop table if exists t_dns_merchant;
create table t_dns_merchant
(
    id            bigint(20)   not null auto_increment primary key comment 'id',
    name          varchar(50)  not null default '' comment '厂商名称',
    merchant_type tinyint(1)   not null default 0 comment '厂商类型',
    description   varchar(255) not null default '' comment '厂商描述',
    has_api       tinyint(1)   not null default 0 comment '是否有api'
) charset = utf8mb4,
  engine = InnoDB, comment ='dns厂商';

drop table if exists t_dns_domain;
create table t_dns_domain
(
    id          bigint(20)   not null auto_increment primary key comment 'id',
    dns_id      bigint(20)   not null default 0 comment 'dns id',
    domain_name varchar(255) not null default '' comment 'domain name',
    expiration  timestamp    not null comment '过期时间',
    is_used     tinyint(1)   not null default 0 comment '是否已被使用'
) charset = utf8mb4,
  engine = InnoDB, comment ='dns域名';

drop table if exists t_dns_sub_domain;
create table t_dns_sub_domain
(
    id          bigint(20)   not null auto_increment primary key comment 'id',
    domain_id   bigint(20)   not null default 0 comment '域名id',
    sub_name    varchar(50)  not null default '子域名',
    domain_type tinyint(1)   not null default 0 comment '域名类型 txt cname ....',
    sub_dest    varchar(255) not null default '' comment '子域名源地址',
    ttl         int(10)      not null default 0 comment 'ttl时间'
) charset = utf8mb4,
  engine = InnoDB, comment ='dns子域名';

drop table if exists t_spider;
create table t_spider
(
    id             bigint(20)   not null auto_increment primary key comment 'id',
    app_id         bigint(20)   not null default 0 comment 'app id',
    name           varchar(50)  not null default '' comment '名称',
    description    varchar(255) not null default '' comment '描述',
    app_type       bigint(20)   not null default 0 comment 'app type',
    source_ip      varchar(100) not null default 0 comment '源ip',
    queue_name     varchar(50)  not null default '' comment '队列名',
    rule           varchar(100) not null default '' comment '规则',
    status         tinyint(1)   not null default 0 comment '状态',
    need_download  tinyint(1)   not null default 0 comment '是否需要下载',
    need_translate tinyint(1)   not null default 0 comment '是否需要翻译',
    need_tag       tinyint(1)   not null default 0 comment '是否需要打标签',
    need_proxy     tinyint(1)   not null default 0 comment '是否需要代理',
    source_code    int          not null default 0 comment '来源id',
    create_time    timestamp    not null default now() comment '创建时间'
) charset = utf8mb4,
  engine = InnoDB, comment ='采集器';

drop table if exists t_spider_task;
create table t_spider_task
(
    id            bigint(20)   not null auto_increment primary key comment 'id',
    task_name     varchar(50)  not null default '' comment '任务名称',
    description   varchar(255) not null default '' comment '任务说明',
    status        tinyint(1)   not null default 0 comment '状态',
    task_url      varchar(255) not null default '' comment '任务url',
    task_schedule varchar(255) not null default 0 comment '调度设定',
    create_time   timestamp    not null default now() comment '创建时间'
) charset = utf8mb4,
  engine = InnoDB, comment ='采集任务';

drop table if exists t_spider_task_log;
create table t_spider_task_log
(
    id          bigint(20)  not null auto_increment primary key comment 'id',
    task_name   varchar(50) not null default '' comment '任务名称',
    task_id     bigint(20)  not null default 0 comment 'task id',
    exec_time   timestamp   not null comment '执行时间',
    exec_status tinyint(1)  not null default 0 comment '执行状态 执行中，执行成功 执行失败',
    exec_log    text        not null comment '执行任务'
) charset = utf8mb4,
  engine = InnoDB, comment ='采集任务';

drop table if exists t_proxy_config;
create table t_proxy_config
(
    id          bigint(20)   not null auto_increment primary key comment 'id',
    addr        varchar(255) not null default '' comment '代理地址',
    protocol    varchar(50)  not null default '' comment '代理协议 sock5 http..',
    create_time timestamp    not null default now() comment '创建时间',
    status      tinyint(1)   not null default 0 comment '状态'
) charset = utf8mb4,
  engine = InnoDB, comment ='采集任务';

drop table if exists t_third_api_config;
create table t_third_api_config
(
    id             bigint(20)   not null auto_increment primary key comment 'id',
    name           varchar(50)  not null default '' comment '第三方协议名称',
    api_addr       varchar(255) not null default '' comment 'api地址',
    api_auth_token varchar(255) not null default '' comment '第三方认证token',
    create_time    timestamp    not null default now() comment '创建时间',
    status         tinyint(1)   not null default 0 comment '状态',
    other_settings text         not null comment '其它设置'
) charset = utf8mb4,
  engine = InnoDB, comment ='采集任务';

drop table if exists t_tag;
create table t_tag
(
    id          bigint(20)   not null auto_increment primary key comment 'id',
    app_type_id bigint(20)   not null default 0 comment 'app type',
    weight      tinyint(2)   not null default 0 comment '权重',
    name        varchar(255) not null default '' comment '标签',
    part        varchar(50)  not null default '' comment '词性',
    index idx_tag_name (name)
) charset = utf8mb4,
  engine = InnoDB, comment ='采集任务';

drop table if exists t_synonym;
create table t_synonym
(
    id          bigint(20)   not null auto_increment primary key comment 'id',
    source_word varchar(200) not null default '' comment '原词',
    dest_word   varchar(200) not null default '' comment '目标词'
) charset = utf8mb4,
  engine = InnoDB, comment ='同义词库';

drop table if exists t_menu;
create table t_menu
(
    id          bigint(20)   not null auto_increment primary key comment 'id',
    app_id      bigint(20)   not null default 0 comment 'app id',
    app_type_id bigint(20)   not null default 0 comment 'app type id',
    app_url_id  bigint(20)   not null default 0 comment 'app url',
    name        varchar(50)  not null default '' comment '名称',
    pid         bigint(20)   not null default 0 comment 'pid',
    menu_level  tinyint(2)   not null default 0 comment '菜单层级',
    status      tinyint(1)   not null default 0 comment '状态',
    component   varchar(200) not null default '' comment '组件',
    title       varchar(200) not null default '' comment '标题',
    keywords    varchar(200) not null default '' comment '关键字',
    description varchar(255) not null default '' comment '描述',
    create_time timestamp    not null default now() comment '创建时间'
) charset = utf8mb4,
  engine = InnoDB, comment ='菜单';

drop table if exists t_comment;
create table t_comment
(
    id          bigint(20)   not null auto_increment primary key comment 'id',
    app_id      bigint(20)   not null default 0 comment 'app id',
    app_type_id bigint(20)   not null default 0 comment 'app type id',
    category_id bigint(20)   not null default 0 comment 'category id',
    title       varchar(255) not null default '' comment '标题',
    username    varchar(50)  not null default '' comment '用户',
    ip          varchar(20)  not null default '' comment '评论ip',
    create_time timestamp    not null default now() comment '创建时间',
    audit_user  varchar(30)  not null default '' comment '审核用户',
    status      tinyint(1)   not null default 0 comment '状态',
    content     text         not null comment '内容'
) charset = utf8mb4,
  engine = InnoDB, comment ='评论';
