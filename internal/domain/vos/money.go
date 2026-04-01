/*
* money é um value object que representa valor monetário com valor e moeda.
 */
package vos

import (
	"fmt"
	"strings"

	"github.com/guiflauzino18/economizze/internal/domain/errors"
)

type Money struct {
	cents    int64
	currency string
}

// NewMoney cria um Money com valor e moeda
func NewMoney(cents int64, currency string) (Money, error) {
	// Padroniza moeda em uppercase e remove espaços
	currency = strings.ToUpper(strings.TrimSpace(currency))

	if len(currency) != 3 {
		return Money{}, errors.NewValidationError("currency", "must be a 3 letter ISO 4217")
	}

	return Money{cents: cents, currency: currency}, nil

}

// Add acrescenta valor ao Money
func (m Money) Add(other Money) (Money, error) {

	// Checa se é a mesma moeda
	if err := m.checkSameCurrency(other); err != nil {
		return Money{}, err
	}

	return Money{cents: m.cents + other.cents, currency: m.currency}, nil
}

// Sub substrai valor de Money
func (m Money) Sub(other Money) (Money, error) {

	// Checa se é a mesma moeda
	if err := m.checkSameCurrency(other); err != nil {
		return Money{}, err
	}

	return Money{cents: m.cents - other.cents, currency: m.currency}, nil
}

// Abs retorna valor positivo
func (m Money) Abs() Money {
	if m.cents < 0 {
		return Money{cents: -m.cents, currency: m.currency}
	}

	return m
}

func (m Money) IsPositive() bool { return m.cents > 0 }
func (m Money) IsNegative() bool { return m.cents < 0 }
func (m Money) IsZero() bool     { return m.cents == 0 }
func (m Money) Cents() int64     { return m.cents }
func (m Money) Currency() string { return m.currency }

// Equals verifica se Money são iguais
func (m Money) Equals(other Money) bool {
	return m.cents == other.cents && m.currency == other.currency
}

// GreaterThan verifica se Money é maior que money passado no parâmetro
func (m Money) GreaterThan(other Money) bool {
	return m.currency == other.currency && m.cents > other.cents
}

// CheckSameCurrency verifica se moeda são iguais
func (m Money) checkSameCurrency(other Money) error {
	if m.currency != other.currency {
		return errors.NewValidationError("currency", fmt.Sprintf("cannot operate on different currencies: %s and %s", m.currency, other.currency))
	}

	return nil
}
