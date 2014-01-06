BEGIN TRANSACTION;

CREATE TABLE announce_log (
	info_hash  string,
	passkey    string,
	key        string,
	ip         string,
	port       int,
	udp        bool,
	uploaded   uint64,
	downloaded uint64,
	left       uint64,
	event      string,
	client     string,
	ts         time
);

COMMIT;
