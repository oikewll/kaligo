# ************************************************************
# Sequel Pro SQL dump
# Version 5446
#
# https://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# Host: 127.0.0.1 (MySQL 8.0.28)
# Database: kaligo
# Generation Time: 2022-05-26 12:46:01 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
SET NAMES utf8mb4;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table admin
# ------------------------------------------------------------

DROP TABLE IF EXISTS `admin`;

CREATE TABLE `admin` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `uid` char(16) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT 'UID',
  `pools` varchar(20) DEFAULT NULL COMMENT '权限池',
  `groups` varchar(1000) NOT NULL DEFAULT '' COMMENT '权限组',
  `username` varchar(20) DEFAULT NULL COMMENT '用户名',
  `password` char(32) DEFAULT NULL COMMENT '用户密码',
  `fake_password` char(32) DEFAULT NULL COMMENT '伪造密码',
  `onetime_password` char(32) DEFAULT NULL COMMENT '一次性密码',
  `realname` varchar(50) DEFAULT NULL COMMENT '真实姓名',
  `email` varchar(50) DEFAULT NULL COMMENT '邮箱',
  `safe_ips` varchar(200) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL COMMENT '登陆IP白名单',
  `first_login` tinyint(1) DEFAULT '1' COMMENT '是否首次登录',
  `date_expired` datetime DEFAULT '2088-06-06 00:00:00' COMMENT '失效日期时间',
  `otp_auth` tinyint(1) DEFAULT '0' COMMENT 'MFA认证等级 0:禁用  1:启用  2:强制启用 [未使用]',
  `otp_authcode` char(16) DEFAULT '' COMMENT 'MFA验证码',
  `need_audit` tinyint(1) DEFAULT '0' COMMENT '登陆是否需要后台进行人工审核 0: 不需要 1:需要',
  `session_id` char(26) DEFAULT '' COMMENT '登陆时session_id',
  `session_expire` int DEFAULT '1440' COMMENT 'SESSION有效期，默认24分钟',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '帐号状态 1:正常 0:禁止登陆',
  `regtime` int NOT NULL COMMENT '注册时间',
  `regip` varchar(15) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '注册IP',
  `logintime` int unsigned NOT NULL COMMENT '最后登录时间',
  `loginip` varchar(15) NOT NULL COMMENT '最后登录IP',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uid` (`uid`),
  KEY `pools` (`pools`),
  KEY `status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COMMENT='用户表';



# Dump of table admin_group
# ------------------------------------------------------------

DROP TABLE IF EXISTS `admin_group`;

CREATE TABLE `admin_group` (
  `id` char(32) NOT NULL DEFAULT '' COMMENT 'ID',
  `name` varchar(20) DEFAULT NULL COMMENT '用户组名称',
  `pools` varchar(20) DEFAULT NULL COMMENT '权限池',
  `purviews` text NOT NULL COMMENT '用户组权限',
  `uptime` int DEFAULT NULL COMMENT '修改时间',
  `addtime` int DEFAULT NULL COMMENT '添加时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COMMENT='用户组表';

LOCK TABLES `admin_group` WRITE;
/*!40000 ALTER TABLE `admin_group` DISABLE KEYS */;

INSERT INTO `admin_group` (`id`, `name`, `pools`, `purviews`, `uptime`, `addtime`)
VALUES
	('1','超级管理员','admin','*',1504839424,1504839424),
	('2','普通管理员','admin','GET-/api/user,GET-/api/user/:id,POST-/api/user,PUT-/api/user/:id,DELETE-/api/user/:id,POST-admin/editpwd,GET-admin/mypurview',1523269932,1504839539),
	('3','服务站','admin','GET-/api/user,GET-/api/user/:id,POST-/api/user,PUT-/api/user/:id,DELETE-/api/user/:id,POST-admin/editpwd,GET-admin/mypurview',1533791312,1529983318);

/*!40000 ALTER TABLE `admin_group` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table admin_loginlog
# ------------------------------------------------------------

DROP TABLE IF EXISTS `admin_loginlog`;

CREATE TABLE `admin_loginlog` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `pools` varchar(20) NOT NULL DEFAULT 'admin' COMMENT '应用池',
  `uid` char(32) NOT NULL DEFAULT '' COMMENT '用户ID',
  `username` varchar(60) NOT NULL DEFAULT '' COMMENT '用户名',
  `session_id` char(26) DEFAULT NULL COMMENT 'SESSION ID',
  `agent` varchar(500) DEFAULT NULL COMMENT '浏览器信息',
  `logintime` int unsigned NOT NULL COMMENT '登录时间',
  `loginip` varchar(15) NOT NULL DEFAULT '' COMMENT '登录IP',
  `logincountry` char(2) DEFAULT NULL COMMENT '登陆国家',
  `loginsta` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '登录时状态 1=成功，0=失败',
  `cli_hash` varchar(32) NOT NULL COMMENT '用户登录名和ip的hash',
  PRIMARY KEY (`id`),
  KEY `logintime` (`logintime`),
  KEY `cli_hash` (`cli_hash`,`loginsta`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COMMENT='用户登陆记录表';



# Dump of table admin_oplog
# ------------------------------------------------------------

DROP TABLE IF EXISTS `admin_oplog`;

CREATE TABLE `admin_oplog` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '操作日志ID',
  `pools` varchar(20) DEFAULT 'admin' COMMENT '应用池',
  `uid` char(32) DEFAULT NULL COMMENT '用户ID',
  `username` varchar(20) NOT NULL DEFAULT '' COMMENT '管理员用户名',
  `session_id` char(26) DEFAULT NULL COMMENT 'SESSION ID',
  `msg` varchar(250) NOT NULL COMMENT '消息内容',
  `do_time` int unsigned NOT NULL COMMENT '发生时间',
  `do_ip` varchar(15) NOT NULL COMMENT '客户端IP',
  `do_country` char(2) NOT NULL DEFAULT '' COMMENT '国家',
  `do_url` varchar(100) NOT NULL COMMENT '操作网址',
  PRIMARY KEY (`id`),
  KEY `user_name` (`username`),
  KEY `do_time` (`do_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COMMENT='用户操作日志';



# Dump of table admin_purview
# ------------------------------------------------------------

DROP TABLE IF EXISTS `admin_purview`;

CREATE TABLE `admin_purview` (
  `pools` varchar(20) DEFAULT 'admin',
  `uid` char(32) NOT NULL DEFAULT '' COMMENT '管理员ID',
  `purviews` text NOT NULL COMMENT '配置字符',
  PRIMARY KEY (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COMMENT='用户权限表';



# Dump of table category
# ------------------------------------------------------------

DROP TABLE IF EXISTS `category`;

CREATE TABLE `category` (
  `id` smallint unsigned NOT NULL AUTO_INCREMENT COMMENT '分类表',
  `name` varchar(50) DEFAULT NULL COMMENT '名称',
  `sort` int DEFAULT '100' COMMENT '排序',
  `create_user` int DEFAULT NULL COMMENT '创建用户',
  `create_time` int DEFAULT NULL COMMENT '创建时间',
  `update_user` int DEFAULT NULL COMMENT '修改用户',
  `update_time` int DEFAULT NULL COMMENT '修改时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COMMENT='分类表';

LOCK TABLES `category` WRITE;
/*!40000 ALTER TABLE `category` DISABLE KEYS */;

INSERT INTO `category` (`id`, `name`, `sort`, `create_user`, `create_time`, `update_user`, `update_time`)
VALUES
	(1,'视频',2,1,1511258578,1,1537002795),
	(2,'音乐',3,1,1511258584,NULL,NULL),
	(3,'小说',4,1,1511258589,1,1537001412);

/*!40000 ALTER TABLE `category` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table config
# ------------------------------------------------------------

DROP TABLE IF EXISTS `config`;

CREATE TABLE `config` (
  `sort` smallint NOT NULL DEFAULT '0' COMMENT '排序id',
  `name` char(100) NOT NULL DEFAULT '' COMMENT '变量名',
  `value` text COMMENT '变量值',
  `title` varchar(50) NOT NULL DEFAULT '' COMMENT '说明标题',
  `info` varchar(200) NOT NULL COMMENT '备注',
  `group` smallint unsigned NOT NULL DEFAULT '1' COMMENT '分组',
  `type` varchar(10) NOT NULL DEFAULT 'string' COMMENT '变量类型',
  PRIMARY KEY (`name`),
  KEY `sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COMMENT='系统配置变量表';

LOCK TABLES `config` WRITE;
/*!40000 ALTER TABLE `config` DISABLE KEYS */;

INSERT INTO `config` (`sort`, `name`, `value`, `title`, `info`, `group`, `type`)
VALUES
	(4,'attachment_image','jpg|png|gif|bmp|ico','图片文件类型','',2,'string'),
	(5,'attachment_media','mp3|avi|mpg|mp4|3gp|flv|rm|rmvb|wmv|swf','多媒体文件类型','',2,'string'),
	(7,'attachment_size','16','最大附件大小(Mb)','',2,'number'),
	(6,'attachment_soft','zip|7z|rar|gz|bz2|tar|iso|exe|dll|doc|xls|ppt|docx|xlsx|pptx|wps|pdf|psd','其它文件件类型','',2,'string'),
	(6,'authorized_time','10','登录授权时间','用户登录多长时间会被踢出',1,'number'),
	(2,'doc_auto_des','1','自动提取摘要','',3,'bool'),
	(6,'doc_auto_des_len','150','自动摘要长度','',3,'number'),
	(1,'doc_auto_keywords','1','自动获取关键字','',3,'bool'),
	(3,'doc_auto_thumb','0','自动提取缩略图','',3,'bool'),
	(7,'doc_down_remove','0','抓取远程资源','',3,'bool'),
	(5,'doc_thumb_h','200','缩略图默认高度','',3,'number'),
	(4,'doc_thumb_w','200','缩略图默认宽度','',3,'number'),
	(0,'ip_limit','','后台登录IP限制','',0,'string'),
	(1,'open_upload','1','是否允许上传文件','',2,'bool'),
	(4,'site_description','PHPCALL开发框架','当前站点摘要信息','',1,'text'),
	(3,'site_keyword','PHPCALL开发框架','当前站点关键字','',1,'string'),
	(1,'site_name','PHPCALL开发框架','当前站点名称','',1,'string'),
	(5,'site_tj','','当前站点统计代码','',1,'text'),
	(2,'site_upload_path','/uploads','附件上传目录','',2,'string'),
	(3,'site_upload_url','http://uploads.phpcall.org','附件目录网址','如果不使用二级域名，此项留空',2,'string'),
	(2,'site_url','http://www.phpcall.org','当前站点URL','',1,'string'),
	(7,'user_guide_url','','用户向导URL','',1,'string');

/*!40000 ALTER TABLE `config` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table content
# ------------------------------------------------------------

DROP TABLE IF EXISTS `content`;

CREATE TABLE `content` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '内容表',
  `catid` smallint DEFAULT NULL COMMENT '分类ID',
  `name` varchar(50) DEFAULT NULL COMMENT '名称',
  `image` varchar(50) DEFAULT NULL COMMENT '封面图',
  `images` varchar(2000) DEFAULT NULL COMMENT '套图',
  `content` text COMMENT '内容',
  `create_user` int DEFAULT NULL COMMENT '创建用户',
  `create_time` int DEFAULT NULL COMMENT '创建时间',
  `update_user` int DEFAULT NULL COMMENT '修改用户',
  `update_time` int DEFAULT NULL COMMENT '修改时间',
  PRIMARY KEY (`id`),
  KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COMMENT='内容表';



# Dump of table crond
# ------------------------------------------------------------

DROP TABLE IF EXISTS `crond`;

CREATE TABLE `crond` (
  `id` smallint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `sort` smallint NOT NULL COMMENT '排序',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '任务名',
  `filename` varchar(248) NOT NULL DEFAULT '' COMMENT '执行脚本',
  `runtime_format` varchar(20) NOT NULL DEFAULT '' COMMENT '执行时间',
  `lasttime` int unsigned NOT NULL DEFAULT '0' COMMENT '最后执行时间',
  `runtime` varchar(30) NOT NULL DEFAULT '0' COMMENT '运行时间',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态：1=启动 0=停止',
  `uptime` int DEFAULT NULL COMMENT '更新时间',
  `addtime` int DEFAULT NULL COMMENT '添加时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 ROW_FORMAT=DYNAMIC COMMENT='计划任务表';



# Dump of table oauth_access_tokens
# ------------------------------------------------------------

DROP TABLE IF EXISTS `oauth_access_tokens`;

CREATE TABLE `oauth_access_tokens` (
  `access_token` char(32) NOT NULL DEFAULT '',
  `client_id` char(32) NOT NULL DEFAULT '',
  `user_id` varchar(80) DEFAULT NULL,
  `openid` char(32) DEFAULT NULL,
  `expires` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `scope` varchar(2000) DEFAULT NULL,
  PRIMARY KEY (`access_token`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

LOCK TABLES `oauth_access_tokens` WRITE;
/*!40000 ALTER TABLE `oauth_access_tokens` DISABLE KEYS */;

INSERT INTO `oauth_access_tokens` (`access_token`, `client_id`, `user_id`, `openid`, `expires`, `scope`)
VALUES
	('4b74f0b5770533bde997c8eca52ab91e','testclient','test888','9d9e0f4a9998c4272b99cd381ec23461','2018-06-10 18:01:05',NULL);

/*!40000 ALTER TABLE `oauth_access_tokens` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table oauth_authorization_codes
# ------------------------------------------------------------

DROP TABLE IF EXISTS `oauth_authorization_codes`;

CREATE TABLE `oauth_authorization_codes` (
  `authorization_code` char(32) NOT NULL DEFAULT '',
  `client_id` char(32) NOT NULL DEFAULT '',
  `user_id` varchar(80) DEFAULT NULL,
  `redirect_uri` varchar(2000) DEFAULT NULL,
  `expires` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `scope` varchar(2000) DEFAULT NULL,
  `id_token` varchar(1000) DEFAULT NULL,
  PRIMARY KEY (`authorization_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;



# Dump of table oauth_clients
# ------------------------------------------------------------

DROP TABLE IF EXISTS `oauth_clients`;

CREATE TABLE `oauth_clients` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(20) DEFAULT NULL COMMENT '应用名称',
  `website` varchar(200) DEFAULT NULL COMMENT '应用网站',
  `cate` tinyint(1) DEFAULT '1' COMMENT '应用分类 1、网页应用 2、客户端',
  `desc` varchar(200) DEFAULT NULL COMMENT '应用介绍',
  `domain` varchar(200) DEFAULT NULL COMMENT '域名绑定，绑定后的域名才可访问client_id',
  `ip` varchar(200) DEFAULT NULL COMMENT '信任IP，以逗号分隔，信任IP才可访问OpenAPI',
  `client_id` char(32) NOT NULL DEFAULT '' COMMENT 'App Key，只生成一次',
  `client_secret` char(32) NOT NULL DEFAULT '' COMMENT 'App Secret，后台可以重置',
  `redirect_uri` varchar(2000) NOT NULL DEFAULT '' COMMENT '授权回调页',
  `cancel_uri` varchar(2000) DEFAULT NULL COMMENT '取消授权回调页',
  `grant_types` varchar(80) DEFAULT NULL COMMENT '授权方式',
  `scope` varchar(2000) DEFAULT NULL COMMENT '授权作用域',
  `user_id` varchar(80) DEFAULT NULL COMMENT '用户ID',
  `create_user` int DEFAULT NULL COMMENT '添加用户',
  `create_time` int DEFAULT NULL COMMENT '添加时间',
  `update_user` int DEFAULT NULL COMMENT '修改用户',
  `update_time` int DEFAULT NULL COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `client_id` (`client_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

LOCK TABLES `oauth_clients` WRITE;
/*!40000 ALTER TABLE `oauth_clients` DISABLE KEYS */;

INSERT INTO `oauth_clients` (`id`, `name`, `website`, `cate`, `desc`, `domain`, `ip`, `client_id`, `client_secret`, `redirect_uri`, `cancel_uri`, `grant_types`, `scope`, `user_id`, `create_user`, `create_time`, `update_user`, `update_time`)
VALUES
	(1,'测试应用','http://www1.phpcall.org',1,'这是一个测试应用','www1.phpcall.org','127.0.0.1,192.168.0.46','testclient','testpass','http://www1.phpcall.org/oauth2_sdk/callback.php',NULL,'authorization_code,refresh_token,password,client_credentials','basic','user',1,1526151992,1,1528379432),
	(7,'fesfes',NULL,1,'fes','fes','fs','4396ff23eb80bbe11fabd578458b27c8','0ef64bd272146903d624a852158b9be0','http://fes','http://fes','authorization_code','basic','user',1,1526151992,1,1526152004);

/*!40000 ALTER TABLE `oauth_clients` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table oauth_jwt
# ------------------------------------------------------------

DROP TABLE IF EXISTS `oauth_jwt`;

CREATE TABLE `oauth_jwt` (
  `client_id` char(32) NOT NULL DEFAULT '',
  `subject` varchar(80) DEFAULT NULL,
  `public_key` varchar(2000) DEFAULT NULL,
  PRIMARY KEY (`client_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;



# Dump of table oauth_permissions
# ------------------------------------------------------------

DROP TABLE IF EXISTS `oauth_permissions`;

CREATE TABLE `oauth_permissions` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(50) DEFAULT NULL COMMENT '权限名称',
  `info` varchar(200) DEFAULT NULL COMMENT '权限说明',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

LOCK TABLES `oauth_permissions` WRITE;
/*!40000 ALTER TABLE `oauth_permissions` DISABLE KEYS */;

INSERT INTO `oauth_permissions` (`id`, `name`, `info`)
VALUES
	(1,'edit',NULL),
	(2,'admin',NULL);

/*!40000 ALTER TABLE `oauth_permissions` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table oauth_refresh_tokens
# ------------------------------------------------------------

DROP TABLE IF EXISTS `oauth_refresh_tokens`;

CREATE TABLE `oauth_refresh_tokens` (
  `refresh_token` char(32) NOT NULL DEFAULT '',
  `client_id` char(32) NOT NULL DEFAULT '',
  `user_id` varchar(80) DEFAULT NULL,
  `openid` char(32) DEFAULT NULL,
  `expires` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `scope` varchar(2000) DEFAULT NULL,
  PRIMARY KEY (`refresh_token`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;



# Dump of table oauth_scopes
# ------------------------------------------------------------

DROP TABLE IF EXISTS `oauth_scopes`;

CREATE TABLE `oauth_scopes` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(100) DEFAULT NULL COMMENT '授权名称',
  `scope` varchar(80) DEFAULT NULL COMMENT '授权',
  `is_default` tinyint(1) DEFAULT '0' COMMENT '是否默认',
  `desc` varchar(200) DEFAULT NULL COMMENT '授权说明',
  `create_user` int DEFAULT NULL,
  `create_time` int DEFAULT NULL,
  `update_user` int DEFAULT NULL,
  `update_time` int DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `scope` (`scope`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

LOCK TABLES `oauth_scopes` WRITE;
/*!40000 ALTER TABLE `oauth_scopes` DISABLE KEYS */;

INSERT INTO `oauth_scopes` (`id`, `name`, `scope`, `is_default`, `desc`, `create_user`, `create_time`, `update_user`, `update_time`)
VALUES
	(1,'基础信息','basic',1,'登陆即可获取：包含userid、userinfo_basic，用户ID、姓名、头像、性别',1,1526146038,1,1528271472),
	(2,'用户信息','userinfo',0,'姓名、头像、性别、省市、Email等',1,1526146038,1,1526146038),
	(3,'用户权限','user_permissions',0,NULL,1,1526146038,1,1526146038),
	(4,'查看下级信息','child_userinfo',0,NULL,1,1526146038,1,1526146038),
	(5,'查看下级详细信息','child_userinfo_all',0,NULL,1,1526146038,1,1526146038),
	(7,'通过关键词搜索用户','search_users_keywords',0,NULL,1,1526146038,1,1526146038),
	(8,'搜索用户时的联想搜索建议','search_users',0,NULL,1,1526146038,1,1526146038),
	(10,'更改头像','update_avatar',0,NULL,1,1526146038,1,1526146038),
	(11,'更改用户资料','update_userinfo',0,NULL,1,1526146038,1,1526146038);

/*!40000 ALTER TABLE `oauth_scopes` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table oauth_users
# ------------------------------------------------------------

DROP TABLE IF EXISTS `oauth_users`;

CREATE TABLE `oauth_users` (
  `userid` char(32) DEFAULT NULL COMMENT '用户ID',
  `username` varchar(80) NOT NULL DEFAULT '' COMMENT '账号',
  `password` varchar(32) DEFAULT NULL COMMENT '密码',
  `realname` varchar(80) DEFAULT NULL COMMENT '真实姓名',
  `avatar` varchar(255) DEFAULT NULL COMMENT '头像',
  `email` varchar(80) DEFAULT NULL COMMENT '邮箱',
  `email_verified` tinyint(1) DEFAULT '0' COMMENT '是否邮箱验证',
  `scope` varchar(2000) DEFAULT 'base' COMMENT '授权范围',
  PRIMARY KEY (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

LOCK TABLES `oauth_users` WRITE;
/*!40000 ALTER TABLE `oauth_users` DISABLE KEYS */;

INSERT INTO `oauth_users` (`userid`, `username`, `password`, `realname`, `avatar`, `email`, `email_verified`, `scope`)
VALUES
	('1a1dc91c907325c69271ddf0c944bc72','user','1a1dc91c907325c69271ddf0c944bc72',NULL,'http://www.dahouduan.com/wp-content/uploads/2017/11/avatar.jpg',NULL,NULL,'basic');

/*!40000 ALTER TABLE `oauth_users` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table user
# ------------------------------------------------------------

DROP TABLE IF EXISTS `user`;

CREATE TABLE `user` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(50) DEFAULT NULL,
  `password` varchar(50) DEFAULT NULL,
  `realname` varchar(50) DEFAULT NULL,
  `email` varchar(50) DEFAULT NULL,
  `status` tinyint(1) DEFAULT '0',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `creator_id` int DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `updator_id` int DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `deletor_id` int DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;

INSERT INTO `user` (`id`, `username`, `password`, `realname`, `email`, `status`, `created_at`, `creator_id`, `updated_at`, `updator_id`, `deleted_at`, `deletor_id`)
VALUES
	(1,'jjj',NULL,NULL,NULL,0,'2022-05-26 17:01:11',NULL,'2022-05-26 17:01:11',NULL,NULL,NULL);

/*!40000 ALTER TABLE `user` ENABLE KEYS */;
UNLOCK TABLES;



/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
