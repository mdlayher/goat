BEGIN TRANSACTION;

CREATE TABLE files_users (
	file_id    int,
	user_id    int,
	ip         string,
	active     bool,
	completed  bool,
	announced  int,
	uploaded   int64,
	downloaded int64,
	left       int64,
	ts         time
);

COMMIT;
