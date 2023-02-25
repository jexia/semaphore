module github.com/jexia/semaphore/v2

go 1.14

require (
	github.com/Knetic/govaluate v3.0.0+incompatible
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/francoispqt/gojay v1.2.13
	github.com/getkin/kin-openapi v0.22.1
	github.com/go-test/deep v1.0.7
	github.com/golang/protobuf v1.4.1
	github.com/graphql-go/graphql v0.7.9
	github.com/hashicorp/consul/api v1.7.0
	github.com/hashicorp/go-immutable-radix v1.2.0
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/hashicorp/hcl/v2 v2.3.0
	github.com/jhump/protoreflect v1.7.0
	github.com/julienschmidt/httprouter v1.3.0
	github.com/miekg/dns v1.1.27 // indirect
	github.com/mitchellh/go-wordwrap v1.0.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.2.1
	github.com/rs/cors v1.7.0
	github.com/spf13/cobra v0.0.6
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/testify v1.6.1
	github.com/zclconf/go-cty v1.3.1
	go.uber.org/zap v1.13.0
	golang.org/x/lint v0.0.0-20191125180803-fdd1cda4f05f // indirect
	golang.org/x/net v0.7.0 // indirect
	google.golang.org/grpc v1.27.0
	google.golang.org/protobuf v1.25.0
	gopkg.in/yaml.v2 v2.3.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
)

replace github.com/francoispqt/gojay v1.2.13 => github.com/Alma-media/gojay v1.2.14
