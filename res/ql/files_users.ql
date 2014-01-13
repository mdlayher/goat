BEGIN TRANSACTION;

CREATE TABLE files_users (
	file_id    int64,
	user_id    int64,
	ip         string,
	active     bool,
	completed  bool,
	announced  int64,
	uploaded   int64,
	downloaded int64,
	left       int64,
	ts         time
);

COMMIT;
