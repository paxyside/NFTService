BEGIN;

ALTER TABLE nfts ADD COLUMN token_id BIGINT;

COMMIT;
