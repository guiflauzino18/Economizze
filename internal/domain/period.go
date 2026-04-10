/*
* period é um value object que representa inicio e fim de um período.
 */

package domain

import (
	"fmt"
	"time"
)

type Period struct {
	start time.Time
	end   time.Time
}

// NewPeriod cria um novo periodo
func NewPediod(start, end time.Time) (Period, error) {

	if start.IsZero() {
		return Period{}, NewValidationError("start", "cannot be zero")
	}

	if end.IsZero() {
		return Period{}, NewValidationError("end", "cannot be zero")
	}

	if !end.After(start) {
		return Period{}, NewValidationError("end", "must be after start")
	}

	return Period{
		start: start,
		end:   end,
	}, nil

}

// NewMontPeriod é um atalho para o mes inteiro
func NewMonthPeriod(year int, month time.Month) Period {
	start := time.Date(year, month, 0, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, -1)
	p, _ := NewPediod(start, end)
	return p
}

func (p Period) Start() time.Time { return p.start }
func (p Period) End() time.Time   { return p.end }

// String retorna periodo em texto
func (p Period) String() string {

	return fmt.Sprintf("%s to %s", p.start.Format("2006-01-02"), p.end.Format("2006-01-02"))
}
