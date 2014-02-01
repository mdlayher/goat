CREATE TABLE IF NOT EXISTS api_keys (
	`id` int(11) NOT NULL AUTO_INCREMENT
	, `user_id` int(11) NOT NULL
	, `key` char(40) NOT NULL
	, `salt` char(20) NOT NULL
	, PRIMARY KEY (`id`)
	, UNIQUE KEY (`key`)
	, UNIQUE KEY (`salt`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin
