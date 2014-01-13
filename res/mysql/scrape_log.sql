CREATE TABLE IF NOT EXISTS scrape_log (
	`id` int(11) NOT NULL AUTO_INCREMENT
	, `info_hash` char(40) NOT NULL
	, `passkey` char(40) NOT NULL
	, `ip` varchar(15) NOT NULL
	, `time` int(11) NOT NULL
	, PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin
