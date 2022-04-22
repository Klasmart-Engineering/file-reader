package instrument

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	newrelic "github.com/newrelic/go-agent"
)

type key int

const (
	KeyNrID key = iota
)

func InitNewRelic() (newrelic.Application, error) {
	var err error
	nrConfig := newrelic.NewConfig("file reader", os.Getenv("NEW_RELIC_LICENSE_KEY_ID"))
	app, err := newrelic.NewApplication(nrConfig)

	if err != nil {
		panic("Failed to setup NewRelic: " + err.Error())

	}
	return app, nil
}

//populateNewRelicInContext get the request context populated
func SetNewRelicInContext() gin.HandlerFunc {

	return func(c *gin.Context) {
		//Setup context
		ctx := c.Request.Context()

		//Set newrelic context
		var txn newrelic.Transaction
		//newRelicTransaction is the key populated by nrgin Middleware
		value, exists := c.Get("newRelicTransaction")
		if exists {
			if v, ok := value.(newrelic.Transaction); ok {
				txn = v
			}
			ctx = context.WithValue(ctx, KeyNrID, txn)
		}
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
