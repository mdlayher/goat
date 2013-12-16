CREATE TABLE IF NOT EXISTS files (
	`id` int(11) NOT NULL AUTO_INCREMENT
	, `info_hash` varchar(40) NOT NULL
	, `leechers` int(11) NOT NULL
	, `seeders` int(11) NOT NULL
	, `completed` int(11) NOT NULL
	, `create_time` int(11) NOT NULL
	, `update_time` int(11) NOT NULL
	, PRIMARY KEY (`id`)
	, UNIQUE KEY (`info_hash`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin
