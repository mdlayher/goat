package goat

var mysql_schemas = []string{`
CREATE TABLE IF NOT EXISTS announce_log (
	id int(11) NOT NULL AUTO_INCREMENT
	, info_hash varchar(40) NOT NULL
	, passkey char(40) NOT NULL
	, ` + "`key` " + `char(8) NOT NULL
	, ip varchar(15) NOT NULL
	, port int(11) NOT NULL
	, udp tinyint(1) NOT NULL
	, uploaded bigint unsigned NOT NULL
	, downloaded bigint unsigned NOT NULL
	, ` + "`left` " + `bigint unsigned NOT NULL
	, event varchar(10) NOT NULL
	, client varchar(50) NOT NULL
	, time int(11) NOT NULL
	, PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin`,
	`CREATE TABLE IF NOT EXISTS api_keys (
	id int(11) NOT NULL AUTO_INCREMENT
	, ` + "`user_id` " + `int(11) NOT NULL
	, ` + "`key` " + `char(40) NOT NULL
	, ` + "`salt` " + `char(20) NOT NULL
	, PRIMARY KEY (id)
	, UNIQUE KEY (user_id)
	, UNIQUE KEY (` + "`key`" + `)
	, UNIQUE KEY (salt)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin`,
	`CREATE TABLE IF NOT EXISTS files (
	id int(11) NOT NULL AUTO_INCREMENT
	, info_hash varchar(40) NOT NULL
	, verified tinyint(1) NOT NULL
	, create_time int(11) NOT NULL
	, update_time int(11) NOT NULL
	, PRIMARY KEY (id)
	, UNIQUE KEY (info_hash)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin`,
	`CREATE TABLE IF NOT EXISTS files_users (
	file_id int(11) NOT NULL
	, ` + "`user_id` " + `int(11) NOT NULL
	, ip varchar(15) NOT NULL
	, active tinyint(1) NOT NULL
	, completed tinyint(1) NOT NULL
	, announced int(11) NOT NULL
	, uploaded bigint unsigned NOT NULL
	, downloaded bigint unsigned NOT NULL
	, ` + "`left` " + `bigint unsigned NOT NULL
	, time int(11) NOT NULL
	, UNIQUE KEY (file_id, user_id, ip)
	, KEY (file_id)
	, KEY (file_id)
	, KEY (ip)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin`,
	`CREATE TABLE IF NOT EXISTS scrape_log (
	id int(11) NOT NULL AUTO_INCREMENT
	, info_hash char(40) NOT NULL
	, passkey char(40) NOT NULL
	, ip varchar(15) NOT NULL
	, time int(11) NOT NULL
	, PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin`,
	`CREATE TABLE IF NOT EXISTS users (
	id int(11) NOT NULL AUTO_INCREMENT
	, username varchar(20) NOT NULL
	, passkey char(40) NOT NULL
	, torrent_limit int(11) NOT NULL
	, PRIMARY KEY (id)
	, UNIQUE KEY (username)
	, UNIQUE KEY (passkey)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin`,
	`CREATE TABLE IF NOT EXISTS whitelist (
	id int(11) NOT NULL AUTO_INCREMENT
	, client varchar(50) NOT NULL
	, approved tinyint(1) NOT NULL
	, PRIMARY KEY (id)
	, UNIQUE KEY (client)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin
`}
