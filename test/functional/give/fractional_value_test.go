package give

import (
	"testing"

	"balance/test/functional"
)

func TestGiveFractionalValue(t *testing.T) {
	var originalAmount, expectedAmount float64
	userID := int64(1)

	test := &test{carcass: functional.NewCarcass(t)}
	test.carcass.ResetModels(userID)

	originalAmount = 0.125
	expectedAmount = 0.12

	test.giveMoneyAndCheckAmountInternal(userID, originalAmount, expectedAmount)
	test.checkGivenMoneyExternal(userID, expectedAmount)
	test.carcass.ResetModels(userID)

	originalAmount = 0.135
	expectedAmount = 0.14

	test.giveMoneyAndCheckAmountInternal(userID, originalAmount, expectedAmount)
	test.checkGivenMoneyExternal(userID, expectedAmount)
	test.carcass.ResetModels(userID)
}
