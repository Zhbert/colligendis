package date_service

import (
	"time"
)

func GetDaysBefore(dateOfArticle time.Time, todayDate time.Time) int {

	diff := todayDate.Sub(dateOfArticle).Hours() / 24

	return int(diff)
}
