BEGIN TRANSACTION;

CREATE TABLE files (
	id          int,
	info_hash   string,
	verified    bool,
	create_time time,
	update_time time
);

COMMIT;
