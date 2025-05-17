module github.com/fmgo

go 1.24

require (
	github.com/stretchr/testify v1.8.4
	github.com/go-redis/redis/v8 v8.11.5
	github.com/beevik/etree v1.5.1
	github.com/golang-jwt/jwt/v5 v5.2.2
	golang.org/x/crypto v0.11.0
	software.sslmate.com/src/go-pkcs12 v0.5.0
)

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/fmgo/core/firma/models => ./core/firma/models
	github.com/fmgo/core/sii => ./core/sii
	github.com/fmgo/pkg/dte => ./pkg/dte
	github.com/fmgo/pkg/sii => ./pkg/sii
)
