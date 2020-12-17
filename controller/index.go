package controller

import (
	"fmt"
	"time"

	"github.com/eviltomorrow/aphrodite-web/db"
	"github.com/eviltomorrow/aphrodite-web/model"

	"github.com/gin-gonic/gin"
)

// Index index
func Index(c *gin.Context) {
	var yesterday = time.Now().AddDate(0, 0, -2)
	var date string
	switch yesterday.Weekday() {
	case time.Saturday:
		date = yesterday.AddDate(0, 0, -1).Format("2006-01-02")

	case time.Sunday:
		date = yesterday.AddDate(0, 0, -2).Format("2006-01-02")

	default:
		date = yesterday.Format("2006-01-02")

	}

	quotes, err := model.SelectQuoteDay(db.MySQL, date)
	if err != nil {
		c.String(200, fmt.Sprintf("Server internal error, nest error: %v", err))
		return
	}

	c.HTML(200, "index.html", gin.H{
		"date":   date,
		"quotes": quotes,
	})
}
