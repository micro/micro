package cockroach

var (
	changeSchema = `create table configs
(
    key    varchar(255) not null
        constraint configs_key
            primary key,
    value  bytea,
    expiry bigint
);`
)
