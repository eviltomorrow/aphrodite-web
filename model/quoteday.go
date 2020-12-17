package model

import (
	"context"

	"github.com/eviltomorrow/aphrodite-web/db"
	jsoniter "github.com/json-iterator/go"
)

// SelectQuoteDay select quote day
func SelectQuoteDay(db db.ExecMySQL, date string) ([]*QuoteDay, error) {
	ctx, cannel := context.WithTimeout(context.Background(), SelectTimeout)
	defer cannel()

	var _sql = `SELECT 
	s.code as code, 
	s.name as name, 
	d.open as open, 
	d.close as close, 
	d.high as high, 
	d.low as low,  
	d.volume as volume, 
	d.account as account,
	concat(format((d.close - d.yesterday_closed) / d.yesterday_closed * 100, 2), '%') as percent
FROM stock s
	LEFT JOIN (
		SELECT *
		FROM quote_day
		WHERE date = ?
	) d
	ON s.code = d.code
WHERE d.open IS NOT NULL;
`
	rows, err := db.QueryContext(ctx, _sql, date)
	if err != nil {
		return nil, err
	}

	var quotes = make([]*QuoteDay, 0, 5000)
	for rows.Next() {
		var quote = QuoteDay{}
		if err := rows.Scan(
			&quote.Code,
			&quote.Name,
			&quote.Open,
			&quote.Close,
			&quote.High,
			&quote.Low,
			&quote.Volume,
			&quote.Account,
			&quote.Percent,
		); err != nil {
			return nil, err
		}
		quote.Date = date
		quotes = append(quotes, &quote)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return quotes, nil
}

// QuoteDay quote day
type QuoteDay struct {
	Code    string  `json:"code"`
	Name    string  `json:"name"`
	Open    float64 `json:"open"`
	Close   float64 `json:"close"`
	High    float64 `json:"high"`
	Low     float64 `json:"low"`
	Volume  int64   `json:"volume"`
	Account string  `json:"account"`
	Date    string  `json:"date"`
	Color   string  `json:"color"`
	Percent string  `json:"percent"`
}

func (q *QuoteDay) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(q)
	return string(buf)
}
