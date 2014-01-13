BEGIN TRANSACTION;

CREATE TABLE announce_log (
	info_hash  string,
	passkey    string,
	key        string,
	ip         string,
	port       int32,
	udp        bool,
	uploaded   int64,
	downloaded int64,
	left       int64,
	event      string,
	client     string,
	ts         time
);

COMMIT;
