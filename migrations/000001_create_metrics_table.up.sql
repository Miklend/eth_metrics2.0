CREATE TABLE IF NOT EXISTS block_metrics (
    LastBlock BIGINT NOT NULL,
    Transaction_count BIGINT NOT NULL,
    Fees NUMERIC NOT NULL,
    Created_at TIMESTAMPTZ,
    PRIMARY KEY (LastBlock)
);
