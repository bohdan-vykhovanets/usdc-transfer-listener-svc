-- +migrate Up
create table transfers
(
    id bigserial primary key not null,
    "from" varchar(42) not null,
    "to" varchar(42) not null,
    "value" numeric not null,
    "block_number" bigint not null,
    "tx_hash" varchar(66) not null,
    "log_index" integer not null
);

alter table transfers add constraint transfers_unique_event unique (block_number, tx_hash, log_index);

create index idx_transfers_from on transfers ("from");
create index idx_transfers_to on transfers ("to");
create index idx_transfers_block_number on transfers (block_number);

-- +migrate Down
drop index idx_transfers_block_number;
drop index idx_transfers_to;
drop index idx_transfers_from;

drop table transfers;

