package core

import (
	"time"

	"github.com/go-pg/pg/v10"

	"balance/internal/models"
	"balance/internal/postgres"
)

func (s *Service) Take(
	userID int64,
	amount float64,
	initiatorID *int64,
	reason string,
	sleep int64,
) (*models.Balance, *models.Transaction, error) {
	balance := new(models.Balance)
	transaction := new(models.Transaction)

	err := s.postgres.RunInTransaction(s.context, func(tx *pg.Tx) error {
		err := postgres.SetIsolationLevelReadCommitted(tx)

		if err != nil {
			s.logger.WithContext(s.context).Error(err)
			return err
		}

		err = tx.Model(balance).
			Where("user_id = ?", userID).
			For("UPDATE").
			Select()

		if err != nil && err != pg.ErrNoRows {
			s.logger.WithContext(s.context).WithField("user_id", userID).Error(err)
			return err
		}

		if err == pg.ErrNoRows {
			return ErrorUserNotFound
		}

		balance.Balance = sub(balance.Balance, amount)
		balance.UpdatedAt = time.Now()

		if balance.Balance < 0 {
			return ErrorInsufficientFunds
		}

		time.Sleep(time.Duration(sleep) * time.Millisecond) // For testing concurrent

		_, err = tx.Model(balance).
			Set("balance = ?balance").
			Where("id = ?id").
			Update()

		if err != nil {
			s.logger.WithContext(s.context).WithField("balance", balance).Error(err)
			return err
		}

		transaction.UserId = userID
		transaction.Amount = -round(amount)
		transaction.InitiatorId = initiatorID
		transaction.Reason = reason
		transaction.CreatedAt = time.Now()

		_, err = tx.Model(transaction).Insert()
		if err != nil {
			s.logger.WithContext(s.context).WithField("transaction", transaction).Error(err)
			return err
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return balance, transaction, nil
}
