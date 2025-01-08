package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Yaska1706/rates-improved/pkg/models"
	"github.com/jmoiron/sqlx"
)

type RateRepo struct {
	db     *sqlx.DB
	logger *slog.Logger
	ctx    context.Context
}

func NewRateRepo(ctx context.Context, db *sqlx.DB, logger *slog.Logger) *RateRepo {
	return &RateRepo{
		db:     db,
		logger: logger,
		ctx:    ctx,
	}
}

func (r *RateRepo) Insert(currencyRates []models.CurrencyRate) error {
	query := `INSERT INTO rates(rate_id,currency,date,rate,created_at) VALUES($1,$2,$3,$4,$5) ON CONFLICT (rate_id) DO UPDATE SET
	rate = EXCLUDED.rate
	;`

	tx, err := r.db.BeginTx(r.ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	//nolint
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(r.ctx, query)
	if err != nil {
		return fmt.Errorf("prepare stmt: %w", err)
	}
	defer stmt.Close()

	r.logger.InfoContext(r.ctx, "Insert process starting")

	for _, currencyRate := range currencyRates {
		_, err := stmt.ExecContext(r.ctx, currencyRate.RateID, currencyRate.Currency, currencyRate.Date, currencyRate.Rate, time.Now())
		if err != nil {
			return fmt.Errorf("DBInsert :%w", err)
		}

	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	r.logger.InfoContext(r.ctx, "Insert/Update successful")

	return nil
}

func (r *RateRepo) GetLatestRates() ([]models.CurrencyRate, error) {
	query := `SELECT r.currency, r.rate
				FROM rates r
				JOIN (
					SELECT currency, MAX(date) AS latest_date
					FROM rates
					GROUP BY currency
				) latest
				ON r.currency = latest.currency AND r.date = latest.latest_date
			;`

	rows, err := r.db.QueryContext(r.ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query : %w", err)
	}

	defer rows.Close()

	var currencyRates []models.CurrencyRate
	for rows.Next() {
		var currencyRate models.CurrencyRate
		if err := rows.Scan(&currencyRate.Currency, &currencyRate.Rate); err != nil {
			return nil, fmt.Errorf("scan(): %w", err)
		}

		currencyRates = append(currencyRates, currencyRate)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return currencyRates, nil
}

func (r *RateRepo) GetCurrencyStats() ([]models.CurrencyStat, error) {
	query := `SELECT
				currency,
				MIN(rate) AS min_rate,
				MAX(rate) AS max_rate,
				ROUND(AVG(rate), 4) AS avg_rate
			FROM
				rates
			GROUP BY
				currency;`

	rows, err := r.db.QueryContext(r.ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query : %w", err)
	}

	defer rows.Close()
	var currencyStats []models.CurrencyStat
	for rows.Next() {
		var currencyStat models.CurrencyStat
		if err := rows.Scan(&currencyStat.Currency,
			&currencyStat.MinRate,
			&currencyStat.MaxRate,
			&currencyStat.AvgRate,
		); err != nil {
			return nil, fmt.Errorf("scan(): %w", err)
		}

		currencyStats = append(currencyStats, currencyStat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return currencyStats, nil

}

func (r *RateRepo) GetRatesByDate(date time.Time) ([]models.CurrencyRate, error) {
	query := `SELECT r.currency , r.rate FROM rates r WHERE r.date = $1`

	rows, err := r.db.QueryContext(r.ctx, query, date)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query : %w", err)
	}
	defer rows.Close()

	var currencyRates []models.CurrencyRate
	for rows.Next() {
		var currencyRate models.CurrencyRate
		if err := rows.Scan(&currencyRate.Currency, &currencyRate.Rate); err != nil {
			return nil, fmt.Errorf("scan(): %w", err)
		}

		currencyRates = append(currencyRates, currencyRate)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return currencyRates, nil
}
