package postgres

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"

	"balance/internal/models"
)

func CreateDatabaseSchema(db *pg.DB) error {
	internalModels := []interface{}{
		(*models.Balance)(nil),
		(*models.Transaction)(nil),
	}

	for _, model := range internalModels {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp:          false,
			IfNotExists:   true,
			Varchar:       0,
			FKConstraints: false,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
