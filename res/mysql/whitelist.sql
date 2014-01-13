CREATE TABLE IF NOT EXISTS whitelist (
	`id` int(11) NOT NULL AUTO_INCREMENT
	, `client` varchar(50) NOT NULL
	, `approved` tinyint(1) NOT NULL
	, PRIMARY KEY (`id`)
	, UNIQUE KEY (`client`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin
