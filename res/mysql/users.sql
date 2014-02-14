CREATE TABLE IF NOT EXISTS users (
	`id` int(11) NOT NULL AUTO_INCREMENT
	, `username` varchar(20) NOT NULL
	, `password` char(60) NOT NULL
	, `passkey` char(40) NOT NULL
	, `torrent_limit` int(11) NOT NULL
	, PRIMARY KEY (`id`)
	, UNIQUE KEY (`username`)
	, UNIQUE KEY (`password`)
	, UNIQUE KEY (`passkey`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin
