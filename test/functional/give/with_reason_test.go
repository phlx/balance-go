package give

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"balance/test/functional"
)

func TestGiveWithReason(t *testing.T) {
	userID := int64(1)
	amount := 5000.00
	reason := "Feel free to идти в McDonald's, чувак 🤑"

	test := &carcassDecorator{carcass: functional.NewCarcass(t)}
	test.carcass.ResetModels(userID)

	test.carcass.API.GiveExpect(userID, amount, nil, reason).Status(http.StatusOK)

	transaction := test.carcass.GetOneTransactionModel(userID)
	assert.Equal(t, reason, transaction.Reason)

	test.carcass.ResetModels(userID)
}
