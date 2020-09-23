package core

import (
	"time"

	"github.com/go-pg/pg/v10"

	"balance/internal/models"
	"balance/internal/postgres"
)

func (s *Service) Move(
	fromUserID int64,
	toUserID int64,
	amount float64,
	reason string,
	sleep int64,
) (*models.Transaction, *models.Transaction, error) {
	balanceFrom := new(models.Balance)
	balanceTo := new(models.Balance)
	transactionFrom := new(models.Transaction)
	transactionTo := new(models.Transaction)

	err := s.postgres.RunInTransaction(s.context, func(tx *pg.Tx) error {
		err := postgres.SetIsolationLevelReadCommitted(tx)

		if err != nil {
			s.logger.WithContext(s.context).Error(err)
			return err
		}

		err = tx.Model(balanceFrom).
			Where("user_id = ?", fromUserID).
			For("UPDATE").
			Select()

		if err != nil && err != pg.ErrNoRows {
			s.logger.WithContext(s.context).WithField("from_user_id", fromUserID).Error(err)
			return err
		}

		if err == pg.ErrNoRows {
			return ErrorUserNotFound
		}

		if sub(balanceFrom.Balance, amount) < 0 {
			return ErrorInsufficientFunds
		}

		err = tx.Model(balanceTo).
			Where("user_id = ?", toUserID).
			For("UPDATE").
			Select()

		if err != nil && err != pg.ErrNoRows {
			s.logger.WithContext(s.context).WithField("from_user_id", fromUserID).Error(err)
			return err
		}

		balanceTo.UpdatedAt = time.Now()
		balanceTo.UserId = toUserID
		if err == pg.ErrNoRows {
			balanceTo.Balance = round(amount)

			_, err = tx.Model(balanceTo).Insert()
		} else {
			balanceTo.Balance = add(balanceTo.Balance, amount)

			_, err = tx.Model(balanceTo).
				Set("balance = ?balance").
				Where("id = ?id").
				Update()
		}

		if err != nil {
			s.logger.WithContext(s.context).WithField("balance_to", balanceTo).Error(err)
			return err
		}

		time.Sleep(time.Duration(sleep) * time.Millisecond) // For testing concurrent

		balanceFrom.UpdatedAt = time.Now()
		balanceFrom.UserId = fromUserID
		balanceFrom.Balance = sub(balanceFrom.Balance, amount)
		_, err = tx.Model(balanceFrom).
			Set("balance = ?balance").
			Where("id = ?id").
			Update()

		if err != nil {
			s.logger.WithContext(s.context).WithField("balance_from", balanceFrom).Error(err)
			return err
		}

		transactionFrom.UserId = fromUserID
		transactionFrom.Amount = -round(amount)
		transactionFrom.InitiatorId = &toUserID
		transactionFrom.Reason = reason
		transactionFrom.CreatedAt = time.Now()

		_, err = tx.Model(transactionFrom).Insert()
		if err != nil {
			s.logger.WithContext(s.context).WithField("transaction_from", transactionFrom).Error(err)
			return err
		}

		transactionTo.UserId = toUserID
		transactionTo.Amount = round(amount)
		transactionTo.InitiatorId = &fromUserID
		transactionTo.Reason = reason
		transactionTo.CreatedAt = time.Now()

		_, err = tx.Model(transactionTo).Insert()
		if err != nil {
			s.logger.WithContext(s.context).WithField("transaction_to", transactionTo).Error(err)
			return err
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return transactionFrom, transactionTo, nil
}
