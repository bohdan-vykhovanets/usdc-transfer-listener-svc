-- +migrate Up
alter table transfers add column block_number bigint not null default 0;
alter table transfers add column tx_hash varchar(66) not null default '0x0000000000000000000000000000000000000000000000000000000000000000';
alter table transfers add column log_index integer not null default 0;

DELETE FROM transfers
WHERE block_number = 0
  AND tx_hash = '0x0000000000000000000000000000000000000000000000000000000000000000'
  AND log_index = 0;

alter table transfers alter column block_number drop default;
alter table transfers alter column tx_hash drop default;
alter table transfers alter column log_index drop default;

alter table transfers add constraint transfers_unique_event unique (block_number, tx_hash, log_index);
create index idx_transfers_block_number on transfers (block_number);

-- +migrate Down

drop index idx_transfers_block_number;
alter table transfers drop constraint transfers_unique_event;

alter table transfers drop column log_index;
alter table transfers drop column tx_hash;
alter table transfers drop column block_number;
