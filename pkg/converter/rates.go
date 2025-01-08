package converter

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/Yaska1706/rates-improved/pkg/models"
)

func XMLtoCurrencyRate(envelope models.Envelope) []models.CurrencyRate {
	var currencyRates []models.CurrencyRate

	for _, dateCube := range envelope.Cube.Cubes {

		for _, currencyCube := range dateCube.Cubes {
			currencyRates = append(currencyRates, models.CurrencyRate{
				RateID:   generateRateID(currencyCube.Currency, FormatDate(dateCube.Time)),
				Date:     FormatDate(dateCube.Time),
				Currency: currencyCube.Currency,
				Rate:     formatFloat(currencyCube.Rate),
			})
		}
	}

	return currencyRates

}

func FormatDate(dateString string) time.Time {
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		slog.Error("failed time conversion", "error", err.Error())
		return time.Now()
	}

	return date
}

func formatFloat(val string) float64 {
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		slog.Error("failed float conversion", "error", err.Error())
	}

	return f
}
func generateRateID(currency string, date time.Time) int64 {

	data := fmt.Sprintf("%s-%s", currency, date.Format("2006-01-02"))

	hash := sha256.Sum256([]byte(data))

	rateID := int64(binary.BigEndian.Uint64(hash[:8]))
	if rateID < 0 {
		rateID = -rateID
	}

	return rateID
}
