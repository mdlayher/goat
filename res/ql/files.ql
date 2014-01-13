BEGIN TRANSACTION;

CREATE TABLE files (
	info_hash   string,
	verified    bool,
	create_time time,
	update_time time
);

COMMIT;
