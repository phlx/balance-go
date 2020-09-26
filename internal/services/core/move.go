package core

import (
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"

	"balance/internal/models"
	time2 "balance/internal/pkg/time"
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
	balances := make([]models.Balance, 0)
	transactionFrom := new(models.Transaction)
	transactionTo := new(models.Transaction)

	err := s.postgres.RunInTransaction(s.context, func(tx *pg.Tx) error {
		err := postgres.SetIsolationLevelReadCommitted(tx)

		if err != nil {
			s.logger.WithContext(s.context).Error(err)
			return err
		}

		err = tx.Model(&balances).
			WhereOr("user_id = ?", fromUserID).
			WhereOr("user_id = ?", toUserID).
			For("UPDATE").
			Select()

		if err != nil {
			s.logger.WithContext(s.context).WithField("from_user_id", fromUserID).Error(err)
			return err
		}

		if len(balances) == 0 || len(balances) == 1 && balances[0].UserId != fromUserID {
			return ErrorUserNotFound
		}

		if len(balances) > 2 {
			err = errors.Errorf("Queried balances for move must be length of 2, %d given", len(balances))
			s.logger.WithContext(s.context).
				WithField("from_user_id", fromUserID).
				WithField("to_user_id", toUserID).
				WithField("balances", balances).
				Error(err)
			return err
		}

		if len(balances) > 1 {
			if fromUserID == balances[0].UserId {
				balanceFrom = &balances[0]
				balanceTo = &balances[1]
			} else {
				balanceFrom = &balances[1]
				balanceTo = &balances[0]
			}
		} else {
			balanceFrom = &balances[0]
		}

		if sub(balanceFrom.Balance, amount) < 0 {
			return ErrorInsufficientFunds
		}

		balanceTo.UpdatedAt = time2.Now()
		if balanceTo.UserId == 0 {
			balanceTo.UserId = toUserID
			balanceTo.Balance = round(amount)

			_, err = tx.Model(balanceTo).Insert()
		} else {
			balanceTo.UserId = toUserID
			balanceTo.Balance = add(balanceTo.Balance, amount)

			_, err = tx.Model(balanceTo).
				Set("balance = ?balance").
				Where("user_id = ?user_id").
				Update()
		}

		if err != nil {
			s.logger.WithContext(s.context).WithField("balance_to", balanceTo).Error(err)
			return err
		}

		time.Sleep(time.Duration(sleep) * time.Millisecond) // For testing concurrent

		balanceFrom.UpdatedAt = time2.Now()
		balanceFrom.UserId = fromUserID
		balanceFrom.Balance = sub(balanceFrom.Balance, amount)
		_, err = tx.Model(balanceFrom).
			Set("balance = ?balance").
			Where("user_id = ?user_id").
			Update()

		if err != nil {
			s.logger.WithContext(s.context).WithField("balance_from", balanceFrom).Error(err)
			return err
		}

		transactionFrom.UserId = fromUserID
		transactionFrom.Amount = -round(amount)
		transactionFrom.InitiatorId = &toUserID
		transactionFrom.Reason = reason
		transactionFrom.CreatedAt = time2.Now()

		_, err = tx.Model(transactionFrom).Insert()
		if err != nil {
			s.logger.WithContext(s.context).WithField("transaction_from", transactionFrom).Error(err)
			return err
		}

		transactionTo.UserId = toUserID
		transactionTo.Amount = round(amount)
		transactionTo.InitiatorId = &fromUserID
		transactionTo.Reason = reason
		transactionTo.CreatedAt = time2.Now()

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
