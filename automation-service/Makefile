FILE=../VERSION
VARIABLE=`cat $(FILE)`
export VERSION=${VARIABLE} created on: `date +'%y-%m-%d::%H:%M:%S'`

all: dep docker

docker:
	# docker build -f deployment/docker/Dockerfile -t ip-team-worldwire-docker-local.artifactory.swg-devops.com/gftn/automation-service .
	# docker tag ip-team-worldwire-docker-local.artifactory.swg-devops.com/gftn/automation-service:latest ip-team-worldwire-docker-local.artifactory.swg-devops.com/gftn/automation-service:latest || exit 1
	rsync -av --exclude='node_modules' --exclude='checkmarx' --exclude='integration-tests' --exclude='*/*.txt'--exclude='api-definitions' \
	 --exclude='anchor-service' --exclude='*/*.toml' --exclude='*/*.yml'  --exclude='*/*.github'  --exclude='*/*.ts'  --exclude='*/*.json' --exclude='*/*.yaml'\
	  --exclude='*/*.md' --exclude='*/*.gitignore' --exclude='rdo-service' --exclude='rdo-client' --exclude='payout-service' --exclude='participant-registry' \
	  --exclude='send-service' --exclude='callback-service' --exclude='administration-service' --exclude='api-service' --exclude='quotes-service'\
	  --exclude='crypto-service' --exclude='fee-service' --exclude='gas-service' --exclude='payment-listener' --exclude='gftn-web' --exclude='kubernetes' \
	 --exclude='.idea'  --exclude='.git' $(subst /automation-service,,$(shell pwd)) deployment/docker/vendor >> /dev/null  
	docker build --build-arg BUILD_VERSION="${VERSION}" -f deployment/docker/Dockerfile -t gftn/automation-service .
	rm -rf deployment/docker/vendor/

dep:
	go get github.com/golang/dep/cmd/dep
	export PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/go/bin:/go/bin
	dep ensure -v

swaggergen:
	go get -u github.com/go-swagger/go-swagger/cmd/swagger
	swagger validate api-definition/internal-api.yaml
	swagger flatten api-definition/internal-api.yaml  > api-definition/internal-api.json

push:
	docker push ip-team-worldwire-docker-local.artifactory.swg-devops.com/gftn/automation-service:latest

build:
	go build

install:
	go install

test:
	go test -run TestAuthForExternalEndpoint
