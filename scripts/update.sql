alter table twitter_oauth add column `redirect_uri`  varchar(255) NOT NULL DEFAULT '' after twitter_screen_name;
alter table discord_oauth add column `redirect_uri`  varchar(255) NOT NULL DEFAULT '' after discord_user_id;


----
alter table mission add column `sort_id`  int(4) NOT NULL DEFAULT 0 after title_cn;



alter table users add column  `from_kol_ref_code` varchar(64) NOT NULL DEFAULT '';
alter table users add column  `from_kol_user_id` varchar(64) NOT NULL DEFAULT '';