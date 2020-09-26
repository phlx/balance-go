package metrics

const (
	ControllersHealthTiming       = "app.controllers.health.timing"
	ControllersBalanceTiming      = "app.controllers.balance.timing"
	ControllersGiveTiming         = "app.controllers.give.timing"
	ControllersTakeTiming         = "app.controllers.take.timing"
	ControllersMoveTiming         = "app.controllers.move.timing"
	ControllersTransactionsTiming = "app.controllers.transactions.timing"

	ControllersHealthCount       = "app.controllers.health.count"
	ControllersBalanceCount      = "app.controllers.balance.count"
	ControllersGiveCount         = "app.controllers.give.count"
	ControllersTakeCount         = "app.controllers.take.count"
	ControllersMoveCount         = "app.controllers.move.count"
	ControllersTransactionsCount = "app.controllers.transactions.count"

	Responses200AllCount = "app.responses.200.all.count"
	Responses400AllCount = "app.responses.200.all.count"
	Responses500AllCount = "app.responses.200.all.count"
)
