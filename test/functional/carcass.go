package functional

import (
	"flag"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/assert"

	"balance/internal/application"
	"balance/internal/models"
)

type Carcass struct {
	T            *testing.T
	App          *application.Application
	Expectations *httpexpect.Expect
	API          *API
}

func IsDebug() bool {
	debug := flag.Lookup("debug")
	if debug == nil {
		return false
	}
	return debug.Value.(flag.Getter).Get().(bool)
}

func NewCarcass(t *testing.T) *Carcass {
	debug := IsDebug()
	expectations, app := Expectations(t, debug)
	api := &API{expectations: expectations}
	return &Carcass{
		T:            t,
		App:          app,
		Expectations: expectations,
		API:          api,
	}
}

func (c *Carcass) GetOneBalanceModel(userID int64) *models.Balance {
	balance := new(models.Balance)
	err := c.App.Container.Postgres.Model(balance).Where("user_id = ?", userID).Select()
	assert.Nil(c.T, err, "error on get one balance model for user ID %d", userID)
	if err != nil {
		c.T.Fatal(err)
	}
	return balance
}

func (c *Carcass) GetOneTransactionModel(userID int64) *models.Transaction {
	transaction := new(models.Transaction)
	err := c.App.Container.Postgres.Model(transaction).Where("user_id = ?", userID).Select()
	assert.Nil(c.T, err, "error on get one transaction model for user ID %d", userID)
	if err != nil {
		c.T.Fatal(err)
	}
	return transaction
}

func (c *Carcass) ResetModels(userID int64) {
	var err error
	_, err = app.Container.Postgres.Model((*models.Balance)(nil)).Where("user_id = ?", userID).Delete()
	assert.Nil(c.T, err, "error on reset user balance models with ID %d", userID)
	_, err = app.Container.Postgres.Model((*models.Transaction)(nil)).Where("user_id = ?", userID).Delete()
	assert.Nil(c.T, err, "error on reset user transaction models with ID %d", userID)
	if err != nil {
		c.T.Fatal(err)
	}
}
