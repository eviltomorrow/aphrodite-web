package model

import (
	"context"
	"fmt"
	"time"

	"github.com/eviltomorrow/aphrodite-web/db"
)

// SelectQuoteDay select quote day
func SelectQuoteDay(db db.ExecMySQL, date string) ([]*QuoteDay, error) {
	ctx, cannel := context.WithTimeout(context.Background(), SelectTimeout)
	defer cannel()

	var _sql = fmt.Sprintf("select id, code, open, close, high, low, yesterday_closed, volume, account, date, day_of_year, create_timestamp, modify_timestamp from quote_day where date = ?")

	rows, err := db.QueryContext(ctx, _sql, date)
	if err != nil {
		return nil, err
	}

	var quotes = make([]*QuoteDay, 0, 16)
	for rows.Next() {
		var quote = QuoteDay{}
		if err := rows.Scan(
			&quote.Code,
			&quote.Open,
			&quote.Close,
			&quote.High,
			&quote.Low,
			&quote.Volume,
			&quote.Account,
			&quote.Date,
		); err != nil {
			return nil, err
		}
		quotes = append(quotes, &quote)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return quotes, nil
}

// QuoteDay quote day
type QuoteDay struct {
	Code    string    `json:"code"`
	Open    float64   `json:"open"`
	Close   float64   `json:"close"`
	High    float64   `json:"high"`
	Low     float64   `json:"low"`
	Volume  int64     `json:"volume"`
	Account float64   `json:"account"`
	Date    time.Time `json:"date"`
}
