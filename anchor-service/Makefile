FILE=../VERSION
VARIABLE=`cat $(FILE)`
export VERSION=${VARIABLE} created on: `date +'%y-%m-%d::%H:%M:%S'`
all: swaggergen docker docs

docker:
	rsync -av --exclude='node_modules' --exclude='checkmarx' --exclude='integration-tests' --exclude='*/*.txt'--exclude='api-definitions' \
	 --exclude='anchor-service'  --exclude='automation-service' --exclude='*/*.toml' --exclude='*/*.yml'  --exclude='*/*.github'  --exclude='*/*.ts' --exclude='*/*.yaml'\
	  --exclude='*/*.md' --exclude='*/*.gitignore' --exclude='rdo-service' --exclude='rdo-client' --exclude='payout-service' --exclude='participant-registry' \
	  --exclude='send-service' --exclude='quotes-service'\
	  --exclude='crypto-service' --exclude='fee-service' --exclude='gas-service' --exclude='payment-listener' --exclude='gftn-web'\
	 --exclude='.idea'  --exclude='.git' $(subst /anchor-service,,$(shell pwd)) deployment/docker/vendor >> /dev/null
	docker build --build-arg BUILD_VERSION="${VERSION}" -f deployment/docker/Dockerfile -t gftn/anchor-service .
	rm -rf deployment/docker/vendor

swaggergen:
	go get -u github.com/go-swagger/go-swagger/cmd/swagger
	swagger validate api-definitions/anchor-api.yaml
	swagger validate api-definitions/anchor-onboarding-api.yaml
	swagger flatten api-definitions/anchor-api.yaml  > api-definitions/anchor-api.json
	swagger flatten api-definitions/anchor-onboarding-api.yaml  > api-definitions/anchor-onboarding-api.json

docs:
	cd api-definitions \
		&& npm install \
		&& npm run md-docs-anchor-api \
		&& npm run md-docs-anchor-onboarding-api

dep:
	go get github.com/golang/dep/cmd/dep
	export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
	dep ensure -vendor-only

lint:
	gometalinter --vendor ./... --disable-all --enable=golint

test:
	go test -run TestAuthForExternalEndpoint

build:
	go build

install:
	go install
