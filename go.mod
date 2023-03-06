module github.com/dp0h/srp

go 1.16

replace golang.org/x/crypto => github.com/golang/crypto v0.0.0-20191227163750-53104e6ec876

require (
	github.com/jessevdk/go-flags v1.5.0
	github.com/rs/zerolog v1.23.0
	github.com/stretchr/testify v1.4.0
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
	golang.org/x/net v0.7.0 // indirect
	gopkg.in/yaml.v2 v2.4.0
)
