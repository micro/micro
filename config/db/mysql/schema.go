package mysql

var (
	changeSchema = `
CREATE TABLE IF NOT EXISTS configs (
id varchar(255) not null primary key,
author varchar(255),
comment text,
timestamp int(11),
changeset_timestamp int(11),
changeset_checksum varchar(255),
changeset_data blob,
changeset_format varchar(255),
changeset_source varchar(255),
index(id, author),
index(timestamp));
`

	changeLogSchema = `
CREATE TABLE IF NOT EXISTS change_log (
pid bigint not null primary key auto_increment,
action varchar(6),
id varchar(255) not null,
path text,
author varchar(255),
comment text,
timestamp int(11),
changeset_timestamp int(11),
changeset_checksum varchar(255),
changeset_data blob,
changeset_format varchar(255),
changeset_source varchar(255),
index(timestamp));
`
)
