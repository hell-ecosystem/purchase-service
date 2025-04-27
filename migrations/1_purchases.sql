CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS purchases (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    name TEXT NOT NULL,
    category TEXT NOT NULL,
    estimated_price NUMERIC(10,2),
    actual_price NUMERIC(10,2),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    deadline TIMESTAMP,
    is_active BOOLEAN NOT NULL DEFAULT true,
    is_purchased BOOLEAN NOT NULL DEFAULT false
);

CREATE INDEX idx_purchases_user_id ON purchases(user_id);
CREATE INDEX idx_purchases_category ON purchases(category);
CREATE INDEX idx_purchases_is_active ON purchases(is_active);
