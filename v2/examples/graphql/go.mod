module github.com/masonhubco/rebar/v2/examples/graphql

go 1.16

require (
	github.com/99designs/gqlgen v0.13.0
	github.com/gin-gonic/gin v1.7.4
	github.com/masonhubco/rebar/samples/graphql v0.0.0-20210118195811-5c9ee77bf413
	github.com/masonhubco/rebar/v2 v2.0.0-20210829193135-c177e6f0065b
	github.com/rs/cors v1.6.0
	github.com/vektah/gqlparser/v2 v2.1.0
)

replace github.com/masonhubco/rebar/v2 v2.0.0-20210829193135-c177e6f0065b => ../../
