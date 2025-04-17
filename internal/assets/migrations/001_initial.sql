-- +migrate Up
create table transfers
(
    id bigserial primary key not null,
    "from" varchar(42) not null,
    "to" varchar(42) not null,
    "value" numeric not null
);

create index idx_transfers_from on transfers ("from");
create index idx_transfers_to on transfers ("to");

-- +migrate Down

drop table transfers;