-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS rates (
	id BIGSERIAL PRIMARY KEY,
    rate_id BIGINT UNIQUE NOT NULL, 
	currency varchar(50) NOT NULL,
	date date  NOT NULL,
	rate decimal NOT NULL,
	created_at timestamp WITH time zone DEFAULT now() NOT NULL,
	updated_at timestamp WITH time zone DEFAULT now() NOT NULL
);
-- Automatically update `updated_at` on row modification
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_rates_updated_at
BEFORE UPDATE ON rates
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS rates;
-- +goose StatementEnd
