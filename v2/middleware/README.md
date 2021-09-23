# rebar middleware

### `ForceSSL`

Redirect to https on production

```go
app := rebar.New(rebar.Options{ /* configs */ })
router.Use(middleware.ForceSSL(rebar.Production))
```

### `I18n`

Detect language from client side and use that for selecting i18n language files

```go
app := rebar.New(rebar.Options{ /* configs */ })
router.Use(middleware.I18n())
```

### `Logger`

Write structured http request logs with zap

```go
logger, err := rebar.NewStandardLogger()
if err != nil {
	log.Fatal("ERROR:", err)
}
app := rebar.New(rebar.Options{ /* configs */ })
router.Use(middleware.Logger(logger))
```

### `Recovery`

Recover from any panics and writes a 500 if there was one.

```go
app := rebar.New(rebar.Options{ /* configs */ })
router.Use(middleware.Recovery())
```

### `Transaction`

Start and inject database transaction for every http request, automatically rollback
when any errors returned from request handler

```go
db, err := somepkg.NewDatabaseConnection()
if err != nil {
	log.Fatal("ERROR:", err)
}
app := rebar.New(rebar.Options{ /* configs */ })
router.Use(middleware.Transaction(db))
```

`db` in the example above should implement `TxWrapper`.

```go
type TxWrapper interface {
	WithTx(tx *sqlx.Tx, fn func(tx *sqlx.Tx) error) error
}
```

### `BaiscJWT`

Allow you to create a middleware with the SystemToken

```go
app := rebar.New(rebar.Options{ /* configs */ })
router.Use(middleware.BasicJWT("test-auth-token"))
```
