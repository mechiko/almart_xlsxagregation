package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/mechiko/utility"
)

type Record struct {
	Cis      *utility.CisInfo
	Korob    string
	Palet    string
	Produced time.Time
	Expired  time.Time
	Order    int64
	Serial   string
}

func NewRecord(row []string) (*Record, error) {
	if len(row) < 8 {
		return nil, fmt.Errorf("записей %d должно быть больше или равно 8", len(row))
	}
	// fakeCode := strings.Replace(row[3], "_x001D_", "\x1D", 1)
	cis, err := utility.ParseCisInfo(row[3])
	if err != nil {
		return nil, fmt.Errorf("получение КМ %w", err)
	}
	// produced, err := parseDate(row[3])
	// if err != nil {
	// 	return nil, fmt.Errorf("ошибка даты производства %s%w", row[3], err)
	// }
	// expired, err := parseDate(row[4])
	// if err != nil {
	// 	return nil, fmt.Errorf("ошибка даты срока годности %s%w", row[4], err)
	// }

	r := &Record{
		Cis:   cis,
		Korob: row[7],
		Palet: row[8],
		// Produced: produced,
		// Expired:  expired,
	}
	return r, nil
}

// полная строка 11 ячеек
func IsRecord(row []string) bool {
	return len(row) >= 6 && strings.HasPrefix(row[3], "01")
}

func parseDate(s string) (time.Time, error) {
	layout := "02.01.2006" // Corresponds to DD.MM.YYYY
	parsedTime, err := time.Parse(layout, s)
	if err != nil {
		return time.Time{}, err
	}
	return parsedTime, nil
}
