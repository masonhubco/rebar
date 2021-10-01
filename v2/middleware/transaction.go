package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/masonhubco/rebar/v2"
)

type TxWrapper interface {
	WithTx(tx *sqlx.Tx, fn func(tx *sqlx.Tx) error) error
}

// Transaction returns a middleware that starts and injects database transaction for every
// http request, automatically rollback when any errors returned from request handler
func Transaction(database TxWrapper) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := database.WithTx(nil, func(tx *sqlx.Tx) error {
			// add the transaction to the context
			c.Set(rebar.TxKey, tx)

			// call the next handler or middleware
			c.Next()

			// if received any errors in context then stop and return
			if len(c.Errors) > 0 {
				return rebar.NewContextErrors([]*gin.Error(c.Errors))
			}

			// check the response status code. if the code is NOT 200..399
			// then it is considered "NOT SUCCESSFUL" and an error will be returned
			statusCode := c.Writer.Status()
			if statusCode < 200 || statusCode >= 400 {
				return errNonSuccess
			}
			return nil
		})
		// err could be one of possible values:
		// - nil - everything went well, if so, return
		// - an error returned from your application, middleware, etc...
		// - a database error - this is returned if there were problems committing the transaction
		// - a errNonSuccess - this is returned if the response status code is not between 200..399
		if err != nil && !errors.Is(err, errNonSuccess) {
			ctxErrs, ok := err.(*rebar.ContextErrors)
			if !ok || !errors.Is(ctxErrs, err) {
				// this is likely a database commit error, because it has not been added
				// to context, database tx wrapper does not have access to gin context,
				// now we need to add it to context so it can be logged by logger middleware
				c.Status(http.StatusInternalServerError)
				c.Error(err)
			}
			c.Abort()
		}
	}
}

var errNonSuccess = errors.New("non success status code")
