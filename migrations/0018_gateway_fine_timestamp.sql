-- +migrate Up
alter table gateway
    add column fpga_id bytea,
    add column fine_timestamp_aes_key bytea;

-- +migrate Down
alter table gateway
    drop column fpga_id,
    drop column fine_timestamp_aes_key;
