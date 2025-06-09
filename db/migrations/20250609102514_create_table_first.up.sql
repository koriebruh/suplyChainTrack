-- Stakeholders (manufacturers, distributors, retailers)
CREATE TABLE stakeholders
(
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name           VARCHAR(255)        NOT NULL,
    type           VARCHAR(50)         NOT NULL, -- 'manufacturer', 'distributor', 'retailer'
    wallet_address VARCHAR(42) UNIQUE,
    email          VARCHAR(255) UNIQUE NOT NULL,
    phone          VARCHAR(20),
    address        TEXT,
    is_verified    BOOLEAN          DEFAULT false,
    created_at     TIMESTAMP        DEFAULT NOW(),
    updated_at     TIMESTAMP        DEFAULT NOW()
);

-- Products
CREATE TABLE products
(
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sku             VARCHAR(100) UNIQUE NOT NULL,
    name            VARCHAR(255)        NOT NULL,
    description     TEXT,
    category        VARCHAR(100),
    manufacturer_id UUID REFERENCES stakeholders (id),
    metadata        JSONB,
    created_at      TIMESTAMP        DEFAULT NOW(),
    updated_at      TIMESTAMP        DEFAULT NOW()
);

-- Supply chain events
CREATE TABLE supply_chain_events
(
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id      UUID REFERENCES products (id),
    stakeholder_id  UUID REFERENCES stakeholders (id),
    event_type      VARCHAR(50) NOT NULL, -- 'manufactured', 'shipped', 'received', 'sold'
    location        VARCHAR(255),
    timestamp       TIMESTAMP   NOT NULL,
    metadata        JSONB,
    blockchain_hash VARCHAR(66),
    is_verified     BOOLEAN          DEFAULT false,
    created_at      TIMESTAMP        DEFAULT NOW()
);

-- Blockchain transactions
CREATE TABLE blockchain_transactions
(
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id         UUID REFERENCES supply_chain_events (id),
    transaction_hash VARCHAR(66) UNIQUE NOT NULL,
    block_number     BIGINT,
    gas_used         BIGINT,
    status           VARCHAR(20)      DEFAULT 'pending', -- 'pending', 'confirmed', 'failed'
    created_at       TIMESTAMP        DEFAULT NOW()
);