package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/eviltomorrow/aphrodite-base/zlog"
	"github.com/eviltomorrow/aphrodite-web/cache"
	"github.com/eviltomorrow/aphrodite-web/db"
	"github.com/eviltomorrow/aphrodite-web/model"
)

// Index index
func Index(c *gin.Context) {
	var date = calculateRecordDate().Format("2006-01-02")
	quotes, err := getCache(date)
	if err == nil {
		zlog.Debug("Get quotes from redis cache success", zap.Int("quotes-count", len(quotes)))
		c.HTML(200, "index.html", gin.H{
			"date":   date,
			"quotes": quotes,
		})
		return
	}

	quotes, err = model.SelectQuoteDay(db.MySQL, date)
	if err != nil {
		c.String(200, fmt.Sprintf("Server internal error, nest error: %v", err))
		return
	}
	zlog.Debug("Get quotes from db success", zap.Int("quotes-count", len(quotes)))

	if len(quotes) != 0 {
		if err = setCache(date, quotes); err == nil {
			zlog.Debug("Set quotes to redis cache success", zap.Int("quotes-count", len(quotes)))
		}
	}

	c.HTML(200, "index.html", gin.H{
		"date":   date,
		"quotes": quotes,
	})
}

func setCache(key string, quotes []*model.QuoteDay) error {
	buf, err := json.Marshal(quotes)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if status := cache.Redis.Set(ctx, key, buf, 3*time.Hour); status.Err() != nil {
		return status.Err()
	}
	return nil
}

func getCache(key string) ([]*model.QuoteDay, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	status := cache.Redis.Get(ctx, key)
	if status.Err() != nil {
		return nil, status.Err()
	}

	buf, err := status.Bytes()
	if err != nil {
		return nil, err
	}

	var quotes = make([]*model.QuoteDay, 0, 100)
	if err := json.Unmarshal(buf, &quotes); err != nil {
		return nil, err
	}
	return quotes, nil
}

func calculateRecordDate() time.Time {
	var now = time.Now()
	var point = time.Date(now.Year(), now.Month(), now.Day(), 23, 10, 0, 0, time.Local)

	switch now.Weekday() {
	case time.Monday:
		if now.After(point) {
			return now.AddDate(0, 0, -3) // 星期五
		}
		return now.AddDate(0, 0, -4) // 星期四

	case time.Tuesday:
		if now.After(point) {
			return now.AddDate(0, 0, -1) // 星期一
		}
		return now.AddDate(0, 0, -4) // 星期五

	case time.Wednesday, time.Thursday, time.Friday:
		if now.After(point) {
			return now.AddDate(0, 0, -1) // 星期二， 星期三， 星期四
		}
		return now.AddDate(0, 0, -2) // 星期一， 星期二， 星期三

	case time.Saturday:
		return now.AddDate(0, 0, -2) // 星期四

	default:
		return now.AddDate(0, 0, -3) // 星期四
	}
}
