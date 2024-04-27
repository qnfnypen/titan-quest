alter table twitter_oauth add column `redirect_uri`  varchar(255) NOT NULL DEFAULT '' after twitter_screen_name;
alter table discord_oauth add column `redirect_uri`  varchar(255) NOT NULL DEFAULT '' after discord_user_id;
