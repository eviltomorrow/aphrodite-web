package model

import (
	"testing"

	"github.com/eviltomorrow/aphrodite-web/db"
	"github.com/stretchr/testify/assert"
)

func TestSelectQuoteDay(t *testing.T) {
	_assert := assert.New(t)
	var date = "2020-12-02"
	quotes, err := SelectQuoteDay(db.MySQL, date)
	_assert.Nil(err)
	t.Logf("Count: %v\r\n", len(quotes))

	date = "2020-12-22"
	quotes, err = SelectQuoteDay(db.MySQL, date)
	_assert.Nil(err)
	t.Logf("Count: %v\r\n", len(quotes))
}

func BenchmarkSelectQuoteDay(b *testing.B) {
	var date = "2020-12-02"
	for i := 0; i < b.N; i++ {
		SelectQuoteDay(db.MySQL, date)
	}
}
