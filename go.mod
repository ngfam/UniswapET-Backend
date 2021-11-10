module github.com/ngfam/uniswap69420

go 1.17

require (
	github.com/99designs/gqlgen v0.14.0
	github.com/vektah/gqlparser/v2 v2.2.0
)

require (
	github.com/Microsoft/go-winio v0.4.15-0.20190919025122-fc70bd9a86b5 // indirect
	github.com/agnivade/levenshtein v1.1.0 // indirect
	github.com/containerd/containerd v1.4.1 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/golang-migrate/migrate v3.5.4+incompatible // indirect
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/hashicorp/golang-lru v0.5.1 // indirect
	github.com/mitchellh/mapstructure v0.0.0-20180220230111-00c29f56e238 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sirupsen/logrus v1.7.0 // indirect
	github.com/stretchr/testify v1.5.1 // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	google.golang.org/genproto v0.0.0-20201030142918-24207fddd1c3 // indirect
	google.golang.org/grpc v1.33.1 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
)

require internal/pkg/db/mysql v1.0.0

replace internal/pkg/db/mysql => ./internal/pkg/db/mysql

require (
	github.com/docker/docker v17.12.0-ce-rc1.0.20200618181300-9dc6525e6118+incompatible
	github.com/go-chi/chi v1.5.4
	github.com/rs/cors v1.8.0
	internal/links v0.0.0-00010101000000-000000000000
	internal/pairs v0.0.0-00010101000000-000000000000
	internal/tokens v0.0.0-00010101000000-000000000000
	internal/users v1.0.0
	pkg/jwt v0.0.0-00010101000000-000000000000
)

replace internal/users => ./internal/users

replace internal/links => ./internal/links

replace internal/tokens => ./internal/tokens

replace internal/pairs => ./internal/pairs

replace pkg/jwt => ./pkg/jwt
