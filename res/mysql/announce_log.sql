CREATE TABLE IF NOT EXISTS announce_log (
	`id` int(11) NOT NULL AUTO_INCREMENT
	, `info_hash` varchar(40) NOT NULL
	, `passkey` char(40) NOT NULL
	, `key` char(8) NOT NULL
	, `ip` varchar(15) NOT NULL
	, `port` int(11) NOT NULL
	, `udp` tinyint(1) NOT NULL
	, `uploaded` bigint unsigned NOT NULL
	, `downloaded` bigint unsigned NOT NULL
	, `left` bigint unsigned NOT NULL
	, `event` varchar(10) NOT NULL
	, `client` varchar(50) NOT NULL
	, `time` int(11) NOT NULL
	, PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin
