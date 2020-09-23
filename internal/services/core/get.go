package core

import (
	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"

	"balance/internal/exchangerates"
	"balance/internal/models"
)

func (s *Service) Get(UserID int64, currency exchangerates.Currency) (float64, error) {
	balance := new(models.Balance)
	err := s.postgres.Model(balance).Where("user_id = ?", UserID).Select()
	if err == pg.ErrNoRows {
		return 0, ErrorUserNotFound
	}

	rate, err := s.getRate(currency)
	if err != nil {
		s.logger.WithContext(s.context).WithField("currency", currency).Error(err)
		return 0, errors.Wrap(err, "failed to get rate for currency")
	}

	multiplied := div(balance.Balance, rate)

	return round(multiplied), nil
}

func (s *Service) getRate(currency exchangerates.Currency) (float64, error) {
	rate := float64(1)
	var err error

	if currency != exchangerates.RUB {
		rate, err = s.exchangeRatesService.GetRate(s.context, currency, exchangerates.RUB)
		if err != nil {
			return 0, err
		}
	}

	return rate, nil
}
