/*
 Navicat Premium Data Transfer

 Source Server         : kiple-test
 Source Server Type    : MySQL
 Source Server Version : 50729
 Source Host           : 218.0.49.239:3306
 Source Schema         : gkuser

 Target Server Type    : MySQL
 Target Server Version : 50729
 File Encoding         : 65001

 Date: 28/09/2020 09:09:37
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for admin_login_history
-- ----------------------------
DROP TABLE IF EXISTS `admin_login_history`;
CREATE TABLE `admin_login_history` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `user_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'users id',
  `ip` varchar(32) NOT NULL DEFAULT '' COMMENT 'login ip',
  `status` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '0:default,1:login success,2:login failed',
  `type` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT 'type(0.login 1.logout)',
  `content` varchar(100) NOT NULL DEFAULT '' COMMENT 'login failed reason',
  `created_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_user_id` (`user_id`) USING BTREE,
  KEY `idx_created_time` (`created_time`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=28 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='admin user login history';

-- ----------------------------
-- Table structure for admin_login_info
-- ----------------------------
DROP TABLE IF EXISTS `admin_login_info`;
CREATE TABLE `admin_login_info` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `user_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'users id',
  `login_history_id` int(11) NOT NULL DEFAULT '0' COMMENT 'table user_login_history:id',
  `wrong_attempt_count` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT 'wrong attempt count',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='admin login info';

-- ----------------------------
-- Table structure for admin_users
-- ----------------------------
DROP TABLE IF EXISTS `admin_users`;
CREATE TABLE `admin_users` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `mobile` varchar(20) NOT NULL COMMENT 'mobile',
  `email` varchar(50) NOT NULL COMMENT 'email',
  `password` varchar(64) NOT NULL COMMENT 'password',
  `username` varchar(50) NOT NULL COMMENT 'nickname',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '1->active 2->inactive',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `google_token` char(16) DEFAULT NULL COMMENT 'google auth code ',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `mobile_index` (`mobile`),
  UNIQUE KEY `email_index` (`email`),
  UNIQUE KEY `name_index` (`username`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=27 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for gk_auth_group
-- ----------------------------
DROP TABLE IF EXISTS `gk_auth_group`;
CREATE TABLE `gk_auth_group` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `group_name` varchar(255) NOT NULL COMMENT 'group name',
  `group_desc` varchar(255) NOT NULL COMMENT 'group description',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for gk_auth_group_role
-- ----------------------------
DROP TABLE IF EXISTS `gk_auth_group_role`;
CREATE TABLE `gk_auth_group_role` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `group_id` int(11) NOT NULL COMMENT 'group id (gk_auth_group.id)',
  `role_id` int(11) NOT NULL COMMENT 'role id(gk_auth_role.id)',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `group_role_index` (`group_id`,`role_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=34 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for gk_auth_log
-- ----------------------------
DROP TABLE IF EXISTS `gk_auth_log`;
CREATE TABLE `gk_auth_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) NOT NULL COMMENT 'portal user id',
  `user_type` tinyint(4) NOT NULL COMMENT '1 admin ,2 merchant,3 user',
  `uri` varchar(255) NOT NULL COMMENT 'request resource',
  `create_time` datetime NOT NULL,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for gk_auth_resource_uri
-- ----------------------------
DROP TABLE IF EXISTS `gk_auth_resource_uri`;
CREATE TABLE `gk_auth_resource_uri` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `parent_id` int(11) DEFAULT NULL COMMENT 'the parent id of uri',
  `resource_type` tinyint(4) NOT NULL DEFAULT '1' COMMENT '1:menu 2: btn',
  `resource_name` varchar(100) NOT NULL COMMENT 'resource name ',
  `resource_uri` varchar(255) DEFAULT NULL COMMENT 'request uri address',
  `resource_desc` varchar(100) NOT NULL COMMENT 'resource function desciption ,like create user',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `sort` int(9) NOT NULL DEFAULT '1' COMMENT 'sort 1+',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uri_index` (`resource_uri`) USING BTREE,
  KEY `parent_id_Index` (`parent_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=56 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for gk_auth_role
-- ----------------------------
DROP TABLE IF EXISTS `gk_auth_role`;
CREATE TABLE `gk_auth_role` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `role_name` varchar(255) NOT NULL COMMENT 'role name',
  `role_desc` varchar(255) NOT NULL COMMENT 'role description',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=37 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for gk_auth_role_uri
-- ----------------------------
DROP TABLE IF EXISTS `gk_auth_role_uri`;
CREATE TABLE `gk_auth_role_uri` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `role_id` int(11) NOT NULL COMMENT 'role id (gk_auth_role.id)',
  `uri_id` int(11) NOT NULL COMMENT 'resource uri id(gk_auth_resource_uri.id)',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `role_uri_key` (`role_id`,`uri_id`) USING BTREE,
  KEY `uri_key` (`uri_id`) USING BTREE,
  KEY `role_id_index` (`role_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=175 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for gk_auth_user_group
-- ----------------------------
DROP TABLE IF EXISTS `gk_auth_user_group`;
CREATE TABLE `gk_auth_user_group` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_type` tinyint(4) NOT NULL COMMENT 'type,1:admin 2:merchant,3:user,101:admin common,102:merchant common,103:user common',
  `user_id` bigint(20) NOT NULL COMMENT 'user id',
  `group_id` int(11) NOT NULL COMMENT 'group id',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_group_unique_index` (`user_id`,`group_id`) USING BTREE,
  KEY `group_index` (`group_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=43 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for gk_merchant_users
-- ----------------------------
DROP TABLE IF EXISTS `gk_merchant_users`;
CREATE TABLE `gk_merchant_users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `merchant_id` bigint(20) NOT NULL COMMENT 'merchant id',
  `name` varchar(20) NOT NULL DEFAULT '' COMMENT 'login name',
  `password` varchar(255) NOT NULL COMMENT 'login pwd',
  `google_token` char(16) NOT NULL,
  `create_time` datetime NOT NULL,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `merchant_id_index` (`merchant_id`) USING BTREE,
  UNIQUE KEY `name_index` (`name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=31 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for merchant_attestation
-- ----------------------------
DROP TABLE IF EXISTS `merchant_attestation`;
CREATE TABLE `merchant_attestation` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `merchant_id` int(11) NOT NULL COMMENT 'merchant',
  `web_site` varchar(100) NOT NULL COMMENT 'merchant web site',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_merchant_id` (`merchant_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='merchant attestation';

-- ----------------------------
-- Table structure for merchant_handle_log
-- ----------------------------
DROP TABLE IF EXISTS `merchant_handle_log`;
CREATE TABLE `merchant_handle_log` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `merchant_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'merchant id',
  `admin_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'admin id',
  `handle_type` varchar(50) NOT NULL DEFAULT '' COMMENT 'handle type(example:reset pin,reset password,update profile,status change)',
  `handle_text` varchar(800) NOT NULL DEFAULT '' COMMENT 'handle type text',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idex_merchant_id` (`merchant_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='merchant handle log';

-- ----------------------------
-- Table structure for merchant_rsa_key
-- ----------------------------
DROP TABLE IF EXISTS `merchant_rsa_key`;
CREATE TABLE `merchant_rsa_key` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'pk/主键',
  `merchant_id` int(10) unsigned NOT NULL COMMENT 'merchant id/商户id',
  `app_id` char(32) NOT NULL COMMENT 'uuid no ''-'' /32位 uuid 除开 连接符 -',
  `merchant_private_key` varchar(2000) NOT NULL COMMENT 'merchant privete key pem type/商户private key pem格式',
  `merchant_public_key` varchar(500) NOT NULL COMMENT 'merchant public key pem type/商户public key pem格式',
  `merchant_key_type` tinyint(4) NOT NULL COMMENT '1->pkcs_1 2->pkcs_8',
  `platform_private_key` varchar(2000) NOT NULL COMMENT 'platform privete key pem type/平台private key pem格式',
  `platform_public_key` varchar(500) NOT NULL COMMENT 'platform public key pem type/平台public key pem格式',
  `platform_key_type` tinyint(4) NOT NULL COMMENT '1->pkcs_1 2->pkcs_8',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `merchant_id` (`merchant_id`) USING BTREE,
  UNIQUE KEY `app_id_index` (`app_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='merchant rsa key table ,include merchant public key ,merchant private key,platform public key ,platform private key';

-- ----------------------------
-- Table structure for merchants
-- ----------------------------
DROP TABLE IF EXISTS `merchants`;
CREATE TABLE `merchants` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(10) unsigned NOT NULL COMMENT 'user id ,merchant admin id',
  `email` varchar(20) NOT NULL DEFAULT '' COMMENT 'email',
  `mobile` varchar(100) NOT NULL DEFAULT '' COMMENT 'mobile',
  `country` varchar(50) NOT NULL DEFAULT '' COMMENT 'merchant country',
  `city` varchar(100) NOT NULL DEFAULT '' COMMENT 'merchant city',
  `address` varchar(200) NOT NULL DEFAULT '' COMMENT 'merchant address',
  `company_name` varchar(100) NOT NULL DEFAULT '' COMMENT 'company name',
  `master_merchant_id` int(11) NOT NULL DEFAULT '0' COMMENT 'master merchant id',
  `sub_master_merchant_id` int(11) NOT NULL DEFAULT '0' COMMENT 'sub master merchant id',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '1->active  2->inactive  3->freeze  10->delete',
  `last_login_time` timestamp NULL DEFAULT NULL COMMENT 'last login time',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_user_id` (`user_id`) USING BTREE,
  KEY `idx_last_login_time` (`last_login_time`) USING BTREE,
  KEY `idx_create_time` (`create_time`) USING BTREE,
  KEY `idx_company_name` (`company_name`) USING BTREE,
  KEY `idx_mobile` (`mobile`) USING BTREE,
  KEY `idx_email` (`email`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='kiple merchants';

-- ----------------------------
-- Table structure for user_device_info
-- ----------------------------
DROP TABLE IF EXISTS `user_device_info`;
CREATE TABLE `user_device_info` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `user_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'users id',
  `username` varchar(80) NOT NULL DEFAULT '' COMMENT 'user name',
  `device_model` varchar(255) NOT NULL DEFAULT '' COMMENT 'login model',
  `device_id` varchar(255) NOT NULL DEFAULT '' COMMENT 'device id',
  `method` varchar(20) NOT NULL DEFAULT '' COMMENT 'request method',
  `source` varchar(255) NOT NULL DEFAULT '' COMMENT 'request source',
  `ip_address` varchar(64) NOT NULL DEFAULT '' COMMENT 'ip',
  `latitude` varchar(20) NOT NULL DEFAULT '' COMMENT 'latitude',
  `longitude` varchar(20) DEFAULT NULL COMMENT 'longitude',
  `created_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=7073 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='user device info';

-- ----------------------------
-- Table structure for user_ekyc_info
-- ----------------------------
DROP TABLE IF EXISTS `user_ekyc_info`;
CREATE TABLE `user_ekyc_info` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'user id',
  `status` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT 'status(1.pending 2.face success 3.face fail)',
  `realname` varchar(128) NOT NULL DEFAULT '' COMMENT 'real name',
  `card_type` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT 'card type(1.MyCard 2.passport)',
  `card_number` varchar(50) NOT NULL DEFAULT '' COMMENT 'card number',
  `nationality` varchar(100) NOT NULL DEFAULT '' COMMENT 'nationality',
  `passport_issue_country` varchar(100) NOT NULL DEFAULT '' COMMENT 'passport issuing country',
  `passport_expire_date` varchar(50) NOT NULL DEFAULT '' COMMENT 'passport expire date',
  `address` varchar(255) NOT NULL DEFAULT '' COMMENT 'address',
  `gender` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT 'gender(0.unknown 1.male 2.female)',
  `birthdate` varchar(20) NOT NULL DEFAULT '' COMMENT 'birthdate(YY-mm-dd)',
  `face_capture_percentage` float(6,2) NOT NULL DEFAULT '0.00' COMMENT 'face capture percentage',
  `created_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
  `updated_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_user_status` (`user_id`,`status`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='user ekyc info table';

-- ----------------------------
-- Table structure for user_extends
-- ----------------------------
DROP TABLE IF EXISTS `user_extends`;
CREATE TABLE `user_extends` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `user_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'users id',
  `id_type` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '1:ID card,2:passport',
  `id_number` varchar(255) NOT NULL DEFAULT '' COMMENT 'ID number',
  `nationality_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'user_nationality.id',
  `state_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'user_state.id',
  `occupation_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'user_occupation.id',
  `nature_business_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'user_nature_business.id',
  `is_pass_ekyc` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT 'is pass ekyc(0.no 1.yes)',
  `latitude` varchar(20) NOT NULL DEFAULT '' COMMENT 'latitude',
  `longitude` varchar(20) DEFAULT NULL COMMENT 'longitude',
  `created_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
  `updated_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uni_user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=3232 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='user extends table';

-- ----------------------------
-- Table structure for user_handle_log
-- ----------------------------
DROP TABLE IF EXISTS `user_handle_log`;
CREATE TABLE `user_handle_log` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `user_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'user id',
  `ip` varchar(32) NOT NULL DEFAULT '' COMMENT 'ip address',
  `client_type` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT 'client type(1.app 2.web)',
  `handle_type` varchar(50) NOT NULL DEFAULT '' COMMENT 'handle type(example:reset pin,reset password,update profile,status change)',
  `handle_text` varchar(800) NOT NULL DEFAULT '' COMMENT 'handle type text',
  `user_agent` varchar(200) NOT NULL DEFAULT '' COMMENT 'request header User-Agent',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idex_user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=3361 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='user handle log';

-- ----------------------------
-- Table structure for user_login_device
-- ----------------------------
DROP TABLE IF EXISTS `user_login_device`;
CREATE TABLE `user_login_device` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `user_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'users id',
  `device_id` varchar(255) NOT NULL DEFAULT '' COMMENT 'login device',
  `count` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'login count',
  `created_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
  `updated_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_uid_deviceid` (`user_id`,`device_id`(30)) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=3212 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='user login device';

-- ----------------------------
-- Table structure for user_login_history
-- ----------------------------
DROP TABLE IF EXISTS `user_login_history`;
CREATE TABLE `user_login_history` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `user_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'users id',
  `ip` varchar(32) NOT NULL DEFAULT '' COMMENT 'login ip',
  `status` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '0:default,1:login success,2:login failed',
  `type` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT 'type(0.login 1.logout)',
  `failed_reason` varchar(100) NOT NULL DEFAULT '' COMMENT 'login failed reason',
  `latitude` varchar(20) NOT NULL DEFAULT '' COMMENT 'latitude',
  `longitude` varchar(20) DEFAULT NULL COMMENT 'longitude',
  `country` varchar(40) NOT NULL DEFAULT '' COMMENT 'login country',
  `device_id` varchar(255) NOT NULL DEFAULT '' COMMENT 'login device',
  `user_agent` varchar(200) NOT NULL DEFAULT '' COMMENT 'request header User-Agent',
  `created_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_user_id` (`user_id`) USING BTREE,
  KEY `idx_created_time` (`created_time`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=3407 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='user login history';

-- ----------------------------
-- Table structure for user_login_info
-- ----------------------------
DROP TABLE IF EXISTS `user_login_info`;
CREATE TABLE `user_login_info` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `user_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'users id',
  `login_history_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT 'table user_login_history:id',
  `wrong_attempt_count` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT 'wrong attempt count',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=3011 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='user login info';

-- ----------------------------
-- Table structure for user_login_location
-- ----------------------------
DROP TABLE IF EXISTS `user_login_location`;
CREATE TABLE `user_login_location` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `user_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'users id',
  `country` varchar(40) NOT NULL DEFAULT '' COMMENT 'login country',
  `city` varchar(64) NOT NULL DEFAULT '' COMMENT 'city',
  `count` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'login count',
  `created_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
  `updated_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_uid_country` (`user_id`,`country`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=3195 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='user login address';

-- ----------------------------
-- Table structure for user_nationality
-- ----------------------------
DROP TABLE IF EXISTS `user_nationality`;
CREATE TABLE `user_nationality` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `value` varchar(80) NOT NULL DEFAULT '' COMMENT 'value',
  `created_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
  `updated_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_value` (`value`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=114 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='user nationality';

-- ----------------------------
-- Table structure for user_nature_business
-- ----------------------------
DROP TABLE IF EXISTS `user_nature_business`;
CREATE TABLE `user_nature_business` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `value` varchar(80) NOT NULL DEFAULT '' COMMENT 'value',
  `created_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
  `updated_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_value` (`value`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='user nature of business table';

-- ----------------------------
-- Table structure for user_occupation
-- ----------------------------
DROP TABLE IF EXISTS `user_occupation`;
CREATE TABLE `user_occupation` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `value` varchar(80) NOT NULL DEFAULT '' COMMENT 'value',
  `created_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
  `updated_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_value` (`value`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='user occupation';

-- ----------------------------
-- Table structure for user_state
-- ----------------------------
DROP TABLE IF EXISTS `user_state`;
CREATE TABLE `user_state` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `value` varchar(80) NOT NULL DEFAULT '' COMMENT 'value',
  `created_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
  `updated_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_value` (`value`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='user state';

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `username` varchar(80) NOT NULL DEFAULT '' COMMENT 'user name',
  `full_name` varchar(80) NOT NULL DEFAULT '' COMMENT 'full name',
  `mobile` varchar(30) NOT NULL DEFAULT '' COMMENT 'mobile number',
  `email` varchar(60) NOT NULL DEFAULT '' COMMENT 'email address',
  `gender` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT 'gender 1:male,2:female',
  `birthdate` date DEFAULT NULL COMMENT 'birth date',
  `pin` varbinary(32) NOT NULL DEFAULT '' COMMENT 'pin code',
  `wallet_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'wallet id',
  `password` varchar(128) NOT NULL DEFAULT '' COMMENT 'user password',
  `refresh_token` varchar(255) NOT NULL DEFAULT '' COMMENT 'refresh token',
  `status` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '0:Inactive,1:Activated,2:Disabled,3:Closed 10:deleted',
  `source_merchant_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'source merchant',
  `address` varchar(255) NOT NULL DEFAULT '' COMMENT 'detail address',
  `city` varchar(64) NOT NULL DEFAULT '' COMMENT 'city',
  `postcode` varchar(12) NOT NULL DEFAULT '' COMMENT 'postcode',
  `user_type` tinyint(4) NOT NULL DEFAULT '1' COMMENT '1->normal user 2->merchant user',
  `uuid` char(32) NOT NULL COMMENT 'user uuid',
  `created_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
  `updated_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_username` (`username`(20)) USING BTREE,
  KEY `idx_smid_username` (`source_merchant_id`,`username`(20)) USING BTREE,
  KEY `idx_mobile_smid` (`mobile`,`source_merchant_id`) USING BTREE,
  KEY `idx_email_smid` (`email`,`source_merchant_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=3239 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='user table';

SET FOREIGN_KEY_CHECKS = 1;
