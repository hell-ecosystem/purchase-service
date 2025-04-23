CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS purchases (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    name TEXT NOT NULL,
    category TEXT NOT NULL,
    estimated_price NUMERIC(10,2),
    actual_price NUMERIC(10,2),
    created_at TIMESTAMP DEFAULT now(),
    deadline TIMESTAMP,
    is_active BOOLEAN DEFAULT true,
    is_purchased BOOLEAN DEFAULT false
);

CREATE INDEX idx_user_id ON purchases(user_id);
CREATE INDEX idx_category ON purchases(category);
