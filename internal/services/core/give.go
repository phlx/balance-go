package core

import (
	"github.com/go-pg/pg/v10"

	"balance/internal/models"
	"balance/internal/pkg/time"
	"balance/internal/postgres"
)

func (s *Service) Give(
	userID int64,
	amount float64,
	initiatorID *int64,
	reason string,
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

		balance.UpdatedAt = time.Now()
		if err == pg.ErrNoRows {
			balance.UserId = userID
			balance.Balance = round(amount)

			_, err = tx.Model(balance).Insert()
		} else {
			balance.Balance = add(balance.Balance, amount)

			_, err = tx.Model(balance).
				Set("balance = ?balance").
				Where("user_id = ?user_id").
				Update()
		}

		if err != nil {
			s.logger.WithContext(s.context).WithField("balance", balance).Error(err)
			return err
		}

		transaction.UserId = userID
		transaction.Amount = round(amount)
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
