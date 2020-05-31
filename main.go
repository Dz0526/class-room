package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
)

var weekdays = map[string]int{
	"Sunday":    0,
	"Monday":    1,
	"Tuesday":   2,
	"Wednesday": 3,
	"Thursday":  4,
	"Friday":    5,
	"Saturday":  6,
}

var match_judge = [8]string{"明日", "Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}

func main() {
	file, err := os.Open("class.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	record, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)

	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	router.POST("/post", func(c *gin.Context) {
		events, err := bot.ParseRequest(req)
		if err != nil {
			log.Fatal(err)
		}

		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TexstMessage:
					message_d := event.Message.(*linebot.TextMessage)
					text := message_d.Text
					// judge text
					judge := match_text(text, match_judge)
					if judge {
						if text == "明日" {
							t := time.Now()
							weekday := t.Weekday()
							// get number by using weekday
							weekday_data := weekdays[weekday]

						} else {
							// get number by using text
							weekday_data := weekdays[text] - 1
						}
						// get class by using csv data
						class_room := record[weekday_data][1]
						// space to \n
						format_class := strings.Replace(class_room, " ", "\n", -1)
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(format_class)).Do; err != nil {
							log.Fatal(err)
						}
					} else {
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("Nothing")).Do; err != nil {
							log.Fatal(err)
						}
					}

				}
			}
		}
	})
	router.Run(":" + port)

}

func match_text(s string, lis [8]string) bool {
	for _, value := range lis {
		if value == s {
			return true
		}
	}
	return false
}
