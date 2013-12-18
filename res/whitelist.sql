CREATE TABLE IF NOT EXISTS whitelist (
	`client` varchar(50) NOT NULL
	, `approved` tinyint(1) NOT NULL
	, PRIMARY KEY (`client`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin
