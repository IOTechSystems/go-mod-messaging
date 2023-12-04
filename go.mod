module github.com/edgexfoundry/go-mod-messaging/v2

go 1.18

require (
	github.com/eclipse/paho.mqtt.golang v1.4.2
	github.com/edgexfoundry/go-mod-core-contracts/v2 v2.3.0
	github.com/go-redis/redis/v7 v7.3.0
	github.com/google/uuid v1.3.0
	github.com/nats-io/nats-server/v2 v2.10.6
	github.com/nats-io/nats.go v1.31.0
	github.com/pebbe/zmq4 v1.2.7
	github.com/stretchr/testify v1.8.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fxamacker/cbor/v2 v2.4.0 // indirect
	github.com/go-kit/log v0.2.0 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-playground/validator/v10 v10.11.1 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/klauspost/compress v1.17.3 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/minio/highwayhash v1.0.2 // indirect
	github.com/nats-io/jwt/v2 v2.5.3 // indirect
	github.com/nats-io/nkeys v0.4.6 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	golang.org/x/crypto v0.15.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.14.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/time v0.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/edgexfoundry/go-mod-core-contracts/v2 => github.com/IOTechSystems/go-mod-core-contracts/v2 v2.3.0
