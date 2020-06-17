CREATE TABLE audit.mysql_audit_log (
 `time` DateTime,
 `date` Date DEFAULT toDate(time),
 `name` String,
 `record` String,
 `command_class` String,
 `connection_id` String,
 `status` UInt32,
 `sqltext` String,
 `user` LowCardinality(String),
 `host` LowCardinality(String),
 `os_user` String,
 `ip` String,
 `db` LowCardinality(String),
 `dbserver` LowCardinality(String)
 )ENGINE = MergeTree() PARTITION BY toDate(time) ORDER BY time SETTINGS index_granularity = 8192; 