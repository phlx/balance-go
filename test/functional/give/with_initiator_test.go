package give

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"balance/test/functional"
)

func TestGiveWithInitiator(t *testing.T) {
	userID := int64(1)

	test := &test{carcass: functional.NewCarcass(t)}
	test.carcass.ResetModels(userID)

	initiatorID := int64(2)
	test.carcass.API.GiveExpect(userID, 100, &initiatorID, "").Status(http.StatusOK)

	transaction := test.carcass.GetOneTransactionModel(userID)
	assert.Equal(t, initiatorID, *transaction.InitiatorId)

	test.carcass.ResetModels(userID)
}
