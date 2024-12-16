BEGIN;

CREATE TABLE transfers (
    id SERIAL PRIMARY KEY,
    from_address VARCHAR(42) NOT NULL,
    to_address VARCHAR(42) NOT NULL,
    tx_hash VARCHAR(66) NOT NULL,
    status VARCHAR(10) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX index_transfers_tx_hash ON transfers (tx_hash);
CREATE INDEX index_nfts_tx_hash ON nfts (tx_hash);

COMMIT;
