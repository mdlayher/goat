BEGIN TRANSACTION;

CREATE TABLE files_users (
	file_id    int,
	user_id    int,
	ip         string,
	active     bool,
	completed  bool,
	announced  int,
	uploaded   uint64,
	downloaded uint64,
	left       uint64,
	ts         time
);

COMMIT;
