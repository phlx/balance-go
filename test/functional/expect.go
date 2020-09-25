package functional

import (
	"context"
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/go-pg/pg/v10"

	"balance/internal/application"
)

var (
	expect *httpexpect.Expect
	app    *application.Application
)

func Expectations(t *testing.T, debug bool) (*httpexpect.Expect, *application.Application) {
	ctx := context.Background()

	if expect == nil {
		app = application.Engine(ctx, debug, true)

		expect = httpexpect.WithConfig(httpexpect.Config{
			Client: &http.Client{
				Transport: httpexpect.NewBinder(app.Router),
				Jar:       httpexpect.NewJar(),
			},
			Reporter: httpexpect.NewAssertReporter(t),
			Printers: []httpexpect.Printer{
				httpexpect.NewDebugPrinter(t, true),
			},
		})

		app.Container.Logger.Debug("Loaded application engine")

		var db string
		_, _ = app.Container.Postgres.Query(pg.Scan(&db), "select current_database()")
		app.Container.Logger.Debugf("Used database name: [%s]", db)
	}

	return expect, app
}
