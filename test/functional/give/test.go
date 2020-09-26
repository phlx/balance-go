package give

import (
	"net/http"

	"github.com/stretchr/testify/assert"

	"balance/test/functional"
)

type carcassDecorator struct { //nolint
	carcass *functional.Carcass
}

func (t carcassDecorator) giveMoneyAndCheckAmountInternal(userID int64, originalAmount, expectedAmount float64) {
	var err error

	t.carcass.API.GiveExpect(userID, originalAmount, nil, "").Status(http.StatusOK)

	balance := t.carcass.GetOneBalanceModel(userID)
	assert.Nil(t.carcass.T, err, "error on selecting balance for user ID %d", userID)
	assert.Equal(t.carcass.T, balance.Balance, expectedAmount)

	transaction := t.carcass.GetOneTransactionModel(userID)
	assert.Nil(t.carcass.T, err, "error on selecting transaction for user ID %d", userID)
	assert.Equal(t.carcass.T, transaction.Amount, expectedAmount)
}

func (t carcassDecorator) checkGivenMoneyExternal(userID int64, expectedAmount float64) {
	t.carcass.API.BalanceExpect(userID).
		Status(http.StatusOK).
		JSON().
		Object().
		ValueEqual("balance", expectedAmount)

	transactions := t.carcass.Expectations.GET("/v1/transactions").
		WithQuery("user_id", userID).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		Value("transactions").
		Array()

	transactions.Length().Equal(1)

	transactions.Element(0).Object().ValueEqual("amount", expectedAmount)
}
