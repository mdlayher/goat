CREATE TABLE IF NOT EXISTS api_keys (
	`id` int(11) NOT NULL AUTO_INCREMENT
	, `user_id` int(11) NOT NULL
	, `pubkey` char(40) NOT NULL
	, `secret` char(40) NOT NULL
	, `expire` int(11) NOT NULL
	, PRIMARY KEY (`id`)
	, UNIQUE KEY (`pubkey`)
	, UNIQUE KEY (`secret`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin
