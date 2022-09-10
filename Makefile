TEST?=$$(go list ./... | grep -v 'vendor')
# HOSTNAME=local
# NAMESPACE=tr
# NAME=warpgate
HOSTNAME=registry.terraform.io
NAMESPACE=andreee94
NAME=warpgate
BINARY=terraform-provider-${NAME}
VERSION=0.0.1
OS_ARCH=linux_amd64

GOBIN=~/go/bin
CLIENTGENFILE=warpgate/client.gen.go
CLIENTCONFIG=warpgate/config.yaml
WARPGATEOPENAPI=https://raw.githubusercontent.com/warp-tech/warpgate/main/warpgate-web/src/admin/lib/openapi-schema.json

default: install

build:
	go build -o ${BINARY}

release:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_darwin_amd64
	GOOS=freebsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_freebsd_386
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_freebsd_amd64
	GOOS=freebsd GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_freebsd_arm
	GOOS=linux GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_linux_386
	GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_linux_amd64
	GOOS=linux GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_linux_arm
	GOOS=openbsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_openbsd_386
	GOOS=openbsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_openbsd_amd64
	GOOS=solaris GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_solaris_amd64
	GOOS=windows GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_windows_386
	GOOS=windows GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_windows_amd64

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test: 
	go test -v ./...
	# go test -i $(TEST) || exit 1                                                   
	# echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=60s -parallel=4                    

testacc: 
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m   

testcov:
	go test -v ./... -cover -coverprofile=coverage.out
	go tool cover -html=coverage.out

install-oapi-codegen:
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

gen-warpgate: install-oapi-codegen
	$(GOBIN)/oapi-codegen -config $(CLIENTCONFIG) $(WARPGATEOPENAPI) > $(CLIENTGENFILE)

install-tfplugindocs:
	go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

gen-doc: install-tfplugindocs
	$(GOBIN)/tfplugindocs generate