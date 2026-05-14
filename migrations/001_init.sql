-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE funds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL CHECK (length(trim(name)) > 0),
    vintage_year INT NOT NULL CHECK (vintage_year BETWEEN 1900 AND 2100),
    target_size_usd NUMERIC(18, 2) NOT NULL CHECK (target_size_usd > 0),
    status TEXT NOT NULL CHECK (status IN ('Fundraising', 'Investing', 'Closed')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE investors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL CHECK (length(trim(name)) > 0),
    investor_type TEXT NOT NULL CHECK (investor_type IN ('Individual', 'Institution', 'Family Office')),
    email TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE investments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    investor_id UUID NOT NULL REFERENCES investors(id) ON DELETE RESTRICT,
    fund_id UUID NOT NULL REFERENCES funds(id) ON DELETE RESTRICT,
    amount_usd NUMERIC(18, 2) NOT NULL CHECK (amount_usd > 0),
    investment_date DATE NOT NULL
);

CREATE INDEX idx_funds_created_at ON funds (created_at DESC);
CREATE INDEX idx_investors_created_at ON investors (created_at DESC);
CREATE INDEX idx_investments_fund_id_date ON investments (fund_id, investment_date DESC);
CREATE INDEX idx_investments_investor_id ON investments (investor_id);

-- +goose Down
DROP INDEX IF EXISTS idx_investments_investor_id;
DROP INDEX IF EXISTS idx_investments_fund_id_date;
DROP INDEX IF EXISTS idx_investors_created_at;
DROP INDEX IF EXISTS idx_funds_created_at;
DROP TABLE IF EXISTS investments;
DROP TABLE IF EXISTS investors;
DROP TABLE IF EXISTS funds;
