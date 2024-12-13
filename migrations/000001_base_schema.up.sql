BEGIN;

CREATE TABLE nfts
(
    id          SERIAL PRIMARY KEY,
    unique_hash VARCHAR(20)   NOT NULL UNIQUE, -- unique hash of the NFT, length 20 from test task
    tx_hash     VARCHAR(66)   NOT NULL, -- ETH hash, length 66 max for ETH
    media_url   VARCHAR(2048) NOT NULL, -- URL to the NFT media, length 2048 for URL
    owner       VARCHAR(42)   NOT NULL, -- ETH address, length 42 max for ETH
    created_at  TIMESTAMP DEFAULT NOW()
);

COMMIT;
