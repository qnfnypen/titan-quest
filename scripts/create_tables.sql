CREATE DATABASE IF NOT EXISTS `titan_quest` DEFAULT CHARACTER SET 'utf8mb4' COLLATE 'utf8mb4_unicode_ci';

USE `titan_quest`;

-- ----------------------------
-- Table structure for users
-- ----------------------------
CREATE TABLE `users` (
`id` bigint(20) NOT NULL AUTO_INCREMENT,
`uuid` varchar(255) NOT NULL DEFAULT '',
`avatar` varchar(255) NOT NULL DEFAULT '',
`username` varchar(255) NOT NULL DEFAULT '',
`pass_hash` varchar(255) NOT NULL DEFAULT '',
`user_email` varchar(255) NOT NULL DEFAULT '',
`wallet_address` varchar(255) NOT NULL DEFAULT '',
`role` tinyint(4) NOT NULL DEFAULT '0',
`allocate_storage` int(1) NOT NULL DEFAULT '0',
`created_at` datetime(3) NOT NULL DEFAULT '0000-00-00 00:00:00.000',
`updated_at` datetime(3) NOT NULL DEFAULT '0000-00-00 00:00:00.000',
`deleted_at` datetime(3) NOT NULL DEFAULT '0000-00-00 00:00:00.000',
`project_id` int(20) NOT NULL DEFAULT '0',
`referral_code` varchar(64) NOT NULL DEFAULT '',
`referrer` varchar(64) NOT NULL DEFAULT '',
`referrer_user_id` varchar(64) NOT NULL DEFAULT '',
`credits` bigint(20) NOT NULL DEFAULT 0,
`from_kol_ref_code` varchar(64) NOT NULL DEFAULT '',
`from_kol_user_id` varchar(64) NOT NULL DEFAULT '',
PRIMARY KEY (`id`),
UNIQUE KEY `uniq_username` (`username`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE `login_log`  (
`id` bigint(20) NOT NULL AUTO_INCREMENT,
`login_username` varchar(50) NOT NULL DEFAULT '',
`ipaddr` varchar(50) NOT NULL DEFAULT '',
`login_location` varchar(255) NOT NULL DEFAULT '',
`browser` varchar(50) NOT NULL DEFAULT '',
`os` varchar(50) NOT NULL DEFAULT '',
`status` tinyint(4) NOT NULL DEFAULT 0,
`msg` varchar(255) NOT NULL DEFAULT '',
`created_at` datetime(3) NOT NULL DEFAULT 0,
PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4;


CREATE TABLE `operation_log`  (
 `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
 `title` varchar(50) NOT NULL DEFAULT '',
 `business_type` int(2) NOT NULL DEFAULT 0 ,
 `method` varchar(100) NOT NULL DEFAULT '',
 `request_method` varchar(10) NOT NULL DEFAULT '',
 `operator_type` int(1) NOT NULL DEFAULT 0,
 `operator_username` varchar(50) NOT NULL DEFAULT '',
 `operator_url` varchar(500)  NOT NULL DEFAULT '',
 `operator_ip` varchar(50) NOT NULL DEFAULT '',
 `operator_location` varchar(255) NOT NULL DEFAULT '',
 `operator_param` text NOT NULL DEFAULT '',
 `json_result` text NOT NULL DEFAULT '',
 `status` int(1) NOT NULL DEFAULT 0,
 `error_msg` varchar(2000) NOT NULL DEFAULT '',
 `created_at` datetime(3) NOT NULL DEFAULT 0,
 `updated_at` datetime(3) NOT NULL DEFAULT 0,
 PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4;


-- ----------------------------
-- Table structure for location_cn
-- ----------------------------

CREATE TABLE `location_cn` (
`id` bigint(20) NOT NULL AUTO_INCREMENT,
`ip` varchar(28) NOT NULL DEFAULT '',
`continent` varchar(28) NOT NULL DEFAULT '',
`country` varchar(128) NOT NULL DEFAULT '',
`province` varchar(128) NOT NULL DEFAULT '',
`city` varchar(128) NOT NULL DEFAULT '',
`longitude` varchar(28) NOT NULL DEFAULT '',
`area_code` varchar(28) NOT NULL DEFAULT '',
`latitude` varchar(28) NOT NULL DEFAULT '',
`isp` varchar(256) NOT NULL DEFAULT '',
`zip_code` varchar(28) NOT NULL DEFAULT '',
`elevation` varchar(28) NOT NULL DEFAULT '',
`created_at` datetime NOT NULL DEFAULT 0,
`updated_at` datetime NOT NULL DEFAULT 0,
PRIMARY KEY (`id`),
UNIQUE KEY `uniq_uuid` (`ip`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for location_en
-- ----------------------------

CREATE TABLE `location_en` (
`id` bigint(20) NOT NULL AUTO_INCREMENT,
`ip` varchar(28) NOT NULL DEFAULT '',
`continent` varchar(28) NOT NULL DEFAULT '',
`country` varchar(128) NOT NULL DEFAULT '',
`province` varchar(128) NOT NULL DEFAULT '',
`city` varchar(128) NOT NULL DEFAULT '',
`longitude` varchar(28) NOT NULL DEFAULT '',
`area_code` varchar(28) NOT NULL DEFAULT '',
`latitude` varchar(28) NOT NULL DEFAULT '',
`isp` varchar(256) NOT NULL DEFAULT '',
`zip_code` varchar(28) NOT NULL DEFAULT '',
`elevation` varchar(28) NOT NULL DEFAULT '',
`created_at` datetime NOT NULL DEFAULT 0,
`updated_at` datetime NOT NULL DEFAULT 0,
PRIMARY KEY (`id`),
UNIQUE KEY `uniq_uuid` (`ip`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;



CREATE TABLE `mission` (
`id` bigint(20) NOT NULL AUTO_INCREMENT,
`title` varchar(128) NOT NULL DEFAULT '',
`title_cn` varchar(128) NOT NULL DEFAULT '',
`channel` varchar(28) NOT NULL DEFAULT '',
`logo` varchar(128) NOT NULL DEFAULT '',
`credit` bigint(20) NOT NULL DEFAULT 0,
`status` int(1) NOT NULL DEFAULT 0,
`open_url` text NOT NULL,
`start_time` datetime NOT NULL DEFAULT 0,
`end_time` datetime NOT NULL DEFAULT 0,
`type` int(4) NOT NULL DEFAULT 0,
`sort_id` int(4) NOT NULL DEFAULT 0,
`parent_id` bigint(20) NOT NULL DEFAULT 0,
`created_at` datetime NOT NULL DEFAULT 0,
`updated_at` datetime NOT NULL DEFAULT 0,
PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


create table sub_mission like mission;


CREATE TABLE `user_mission` (
`id` bigint(20) NOT NULL AUTO_INCREMENT,
`username` varchar(128) NOT NULL DEFAULT '',
`mission_id` bigint(20) NOT NULL DEFAULT 0,
`sub_mission_id` bigint(20) NOT NULL DEFAULT 0,
`type` int(4) NOT NULL DEFAULT 0,
`credit` bigint(20) NOT NULL DEFAULT 0,
`content` varchar(128) NOT NULL DEFAULT '',
`created_at` datetime NOT NULL DEFAULT 0,
`updated_at` datetime NOT NULL DEFAULT 0,
PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;



CREATE TABLE `twitter_oauth` (
`id` bigint(20) NOT NULL AUTO_INCREMENT,
`request_token` varchar(128) NOT NULL DEFAULT '',
`username` varchar(128) NOT NULL DEFAULT '',
`twitter_user_id` varchar(128) NOT NULL DEFAULT '',
`twitter_screen_name` varchar(128) NOT NULL DEFAULT '',
`redirect_uri`  varchar(255) NOT NULL DEFAULT '',
`created_at` datetime NOT NULL DEFAULT 0,
`updated_at` datetime NOT NULL DEFAULT 0,
PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE `discord_oauth` (
`id` bigint(20) NOT NULL AUTO_INCREMENT,
`state` varchar(128) NOT NULL DEFAULT '',
`username` varchar(128) NOT NULL DEFAULT '',
`discord_user_id` varchar(128) NOT NULL DEFAULT '',
`redirect_uri`  varchar(255) NOT NULL DEFAULT '',
`email` varchar(128) NOT NULL DEFAULT '',
`created_at` datetime NOT NULL DEFAULT 0,
`updated_at` datetime NOT NULL DEFAULT 0,
PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE `user_twitter_link` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `username` varchar(128) NOT NULL DEFAULT '',
 `link` varchar(128) NOT NULL DEFAULT '',
 `mission_id` bigint(20) NOT NULL DEFAULT 0,
 `created_at` datetime NOT NULL DEFAULT 0,
 `updated_at` datetime NOT NULL DEFAULT 0,
 PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
