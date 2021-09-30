
all: swaggergen docker docs

docker:
	rsync -av --exclude='node_modules' --exclude='checkmarx' --exclude='integration-tests' --exclude='*/*.txt'--exclude='api-definitions' \
	 --exclude='anchor-service'  --exclude='automation-service' --exclude='*/*.toml' --exclude='*/*.yml'  --exclude='*/*.github'  --exclude='*/*.ts'  --exclude='*/*.json' --exclude='*/*.yaml'\
	  --exclude='*/*.md' --exclude='*/*.gitignore' --exclude='rdo-service' --exclude='rdo-client' --exclude='payout-service' --exclude='participant-registry' \
	  --exclude='administration-service' --exclude='api-service' --exclude='quotes-service' --exclude='send-service'\
	  --exclude='crypto-service' --exclude='fee-service' --exclude='gas-service' --exclude='payment-listener' --exclude='ww-gateway' --exclude='gftn-web' --exclude='kubernetes' \
	 --exclude='.idea'  --exclude='.git' $(subst /ww-gateway,,$(shell pwd)) deployment/docker/vendor >> /dev/null
	docker build -f deployment/docker/Dockerfile -t gftn/ww-gateway .
	rm -rf deployment/docker/vendor

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
	git submodule add --force  ../utility
