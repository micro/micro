package mysql

var (
	changeSchema = `
CREATE TABLE IF NOT EXISTS configs (
key varchar(255) not null primary key,
value blob,
expiry bigint(20)
);
`
)
