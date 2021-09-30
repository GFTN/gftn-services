FILE=../VERSION
VARIABLE=`cat $(FILE)`
export VERSION=${VARIABLE} created on: `date +'%y-%m-%d::%H:%M:%S'`

docker:
	rsync -av --exclude='node_modules' --exclude='checkmarx' --exclude='integration-tests' --exclude='*/*.txt'--exclude='api-definitions' \
	 --exclude='anchor-service'  --exclude='automation-service' --exclude='*/*.toml' --exclude='*/*.yml'  --exclude='*/*.github'  --exclude='*/*.ts'  --exclude='*/*.json' --exclude='*/*.yaml'\
	  --exclude='*/*.md' --exclude='*/*.gitignore' --exclude='rdo-service' --exclude='rdo-client' --exclude='payout-service' --exclude='participant-registry' \
	  --exclude='send-service' --exclude='administration-service' --exclude='api-service' --exclude='quotes-service'\
	  --exclude='crypto-service' --exclude='fee-service' --exclude='gas-service' --exclude='payment-listener' --exclude='gftn-web' --exclude='kubernetes' \
	 --exclude='.idea'  --exclude='.git' $(subst /fee-service,,$(shell pwd)) deployment/docker/vendor >> /dev/null 
	docker build --build-arg BUILD_VERSION="${VERSION}" -f deployment/docker/Dockerfile -t gftn/fee-service .
	rm -rf deployment/docker/vendor

docs:
	cd api-definitions \
		&& npm install \
		&& npm run md-docs-client-api \


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
