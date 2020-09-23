package core

import (
	"math"
	"strings"

	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"

	"balance/internal/models"
)

func (s *Service) List(
	userID int64,
	sort map[string]string,
	limit int,
	page int,
) ([]models.Transaction, int, error) {
	transactions := make([]models.Transaction, 0)
	count := 0

	order, err := order(sort)
	if err != nil {
		return nil, 0, ErrorInvalidSort
	}

	err = s.postgres.RunInTransaction(s.context, func(tx *pg.Tx) error {
		count, err = tx.Model(&transactions).
			Where("user_id = ?", userID).
			Order(order).
			Limit(limit).
			Offset(limit * (page - 1)).
			SelectAndCount()

		if err != nil {
			s.logger.
				WithContext(s.context).
				WithField("user_id", userID).
				WithField("order", order).
				Error(errors.Wrap(err, "Failed to select and count transactions"))

			return err
		}
		return nil
	})

	if err != nil {
		return nil, 0, err
	}

	pages := int(math.Ceil(float64(count) / float64(limit)))

	return transactions, pages, nil
}

func order(sort map[string]string) (string, error) {
	var field, direction string
	for field, direction = range sort {
	} // get last pair
	fieldsToColumns := map[string]string{
		"id":     "id",
		"time":   "created_at",
		"amount": "amount",
	}
	directions := []string{
		"asc",
		"desc",
	}

	isValidField := false
	for f := range fieldsToColumns {
		if f == field {
			isValidField = true
			break
		}
	}
	if !isValidField {
		return "", errors.Errorf("Invalid order field %s", field)
	}

	isValidDirection := false
	for _, d := range directions {
		if d == strings.ToLower(direction) {
			isValidDirection = true
			break
		}
	}
	if !isValidDirection {
		return "", errors.Errorf("Invalid order direction %s", field)
	}

	return fieldsToColumns[field] + " " + strings.ToUpper(direction), nil
}
