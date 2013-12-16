CREATE TABLE IF NOT EXISTS files_users (
	`file_id` int(11) NOT NULL
	, `user_id` int(11) NOT NULL
	, `active` tinyint(1) NOT NULL
	, `completed` tinyint(1) NOT NULL
	, `announced` int(11) NOT NULL
	, `uploaded` bigint unsigned NOT NULL
	, `downloaded` bigint unsigned NOT NULL
	, `left` bigint unsigned NOT NULL
	, `time` int(11) NOT NULL
	, UNIQUE KEY (`file_id`, `user_id`)
	, KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin
