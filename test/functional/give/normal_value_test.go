package give

import (
	"testing"

	"balance/test/functional"
)

func TestGiveNormalValue(t *testing.T) {
	userID := int64(1)

	test := &test{carcass: functional.NewCarcass(t)}
	test.carcass.ResetModels(userID)

	amount := float64(100100100)
	test.giveMoneyAndCheckAmountInternal(userID, amount, amount)
	test.checkGivenMoneyExternal(userID, amount)
	test.carcass.ResetModels(userID)
}
