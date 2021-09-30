FILE=../VERSION
VARIABLE=`cat $(FILE)`
export VERSION=${VARIABLE} created on: `date +'%y-%m-%d::%H:%M:%S'`

all: swaggergen docker docs

docker:
	rsync -av --exclude='node_modules' --exclude='checkmarx' --exclude='integration-tests' --exclude='*/*.txt'--exclude='api-definitions' \
	 --exclude='anchor-service'  --exclude='automation-service' --exclude='*/*.toml' --exclude='*/*.yml'  --exclude='*/*.github'  --exclude='*/*.ts'  --exclude='*/*.json' --exclude='*/*.yaml'\
	  --exclude='*/*.md' --exclude='*/*.gitignore' --exclude='rdo-service' --exclude='rdo-client' --exclude='payout-service' --exclude='participant-registry' \
	  --exclude='send-service' --exclude='administration-service' --exclude='api-service' --exclude='quotes-service'\
	  --exclude='crypto-service' --exclude='fee-service' --exclude='gas-service' --exclude='payment-listener' --exclude='gftn-web' --exclude='kubernetes' \
	 --exclude='.idea'  --exclude='.git' $(subst /crypto-service,,$(shell pwd)) deployment/docker/vendor >> /dev/null 
	docker build --build-arg BUILD_VERSION="${VERSION}" -f deployment/docker/dev/Dockerfile -t gftn/crypto-service --build-arg VAULT_CERT="$(cat ../utility/vault/certs/certificate.crt)" --build-arg VAULT_KEY="$(cat ../utility/vault/certs/privateKey.key)" .
	rm -rf deployment/docker/vendor

docker-prod:
	rsync -av --exclude='node_modules' --exclude='checkmarx' --exclude='integration-tests' --exclude='*/*.txt'--exclude='api-definitions' \
	 --exclude='anchor-service'  --exclude='automation-service' --exclude='*/*.toml' --exclude='*/*.yml'  --exclude='*/*.github'  --exclude='*/*.ts'  --exclude='*/*.json' --exclude='*/*.yaml'\
	  --exclude='*/*.md' --exclude='*/*.gitignore' --exclude='rdo-service' --exclude='rdo-client' --exclude='payout-service' --exclude='participant-registry' \
	  --exclude='send-service' --exclude='administration-service' --exclude='api-service' --exclude='quotes-service'\
	  --exclude='crypto-service' --exclude='fee-service' --exclude='gas-service' --exclude='payment-listener' --exclude='gftn-web' --exclude='kubernetes' \
	 --exclude='.idea'  --exclude='.git' $(subst /crypto-service,,$(shell pwd)) deployment/docker/vendor >> /dev/null 
	docker build --build-arg BUILD_VERSION="${VERSION}" -f deployment/docker/prod/Dockerfile -t gftn/crypto-service-prod --build-arg VAULT_CERT="$(cat ../utility/vault/certs/certificate.crt)" --build-arg VAULT_KEY="$(cat ../utility/vault/certs/privateKey.key)" .
	rm -rf deployment/docker/vendor

swaggergen:
	go get -u github.com/go-swagger/go-swagger/cmd/swagger
	swagger validate api-definitions/internal-api.yaml
	swagger flatten api-definitions/internal-api.yaml  > api-definitions/internal-api.json

docs:
	cd api-definitions \
		&& npm install \
		&& npm run md-docs-internal-api

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

gitsub:
	# Setting up your project in vscode may require that you add the utility repo as a submodule
	# Reminder: Make sure you have the correct branch compatabile with your api-service branch 
	# eg: PS C:\Users\your.user\go\src\github.ibm.com\gftn\api-service> git submodule add --force  ../utility
	git submodule add --force  ../utility
