
build:
	./gradlew compileJava
	./gradlew compileGroovy
	./gradlew compileTestJava
	./gradlew compileTestGroovy

service-docker:
ifeq ($(build), all)
	docker login -u ${DOCKER_USER} -p ${DOCKER_PASSWORD} ${DOCKER_REGISTRY}
	cd .. \
		&& make dep
	cd ../gftn-models \
		&& make
	cd ../api-service \
		&& make docker
	cd ../quotes-service \
		&& make docker
	cd ../global-whitelist-service \
		&& make docker
	cd ../crypto-service \
		&& make docker \
    	&& make docker-prod
	cd ../payment-listener \
		&& make docker
	cd ../participant-registry \
		&& make docker
	cd ../administration-service \
		&& make docker
	cd ../auth-service \
		&& make docker
	cd ../anchor-service \
		&& make docker
	cd ../send-service \
		&& make docker
	cd ../payout-service \
		&& make docker
	cd ../fee-service \
		&& make docker
	cd ../ww-gateway \
		&& make docker			
else ifeq ($(build), gas-service)
	docker login -u ${DOCKER_USER} -p ${DOCKER_PASSWORD} ${DOCKER_REGISTRY}
	cd ../gas-service \
		&& make docker
else ifeq ($(build), auth-service)
	docker login -u ${DOCKER_USER} -p ${DOCKER_PASSWORD} ${DOCKER_REGISTRY}
	cd ../auth-service \
		&& make docker
else ifeq ($(build), crypto-service)
	docker login -u ${DOCKER_USER} -p ${DOCKER_PASSWORD} ${DOCKER_REGISTRY}
	cd ../crypto-service \
		&& make docker \
		&& make docker-prod
else ifeq ($(build), msk-cli)
	docker login -u ${DOCKER_USER} -p ${DOCKER_PASSWORD} ${DOCKER_REGISTRY}
	cd ../automation-service/msk \
		&& make build-msk-cli
else ifeq ($(build), msk-server)
	docker login -u ${DOCKER_USER} -p ${DOCKER_PASSWORD} ${DOCKER_REGISTRY}
	cd ../automation-service/msk \
		&& make build-msk-burrow
else ifeq ($(build), msk-ui)
	docker login -u ${DOCKER_USER} -p ${DOCKER_PASSWORD} ${DOCKER_REGISTRY}
	cd ../automation-service/msk \
		&& make build-msk-burrow-ui
else
		cd .. \
			&& make dep
		cd ../gftn-models \
			&& make
		cd ../$(build) \
		&& make docker || exit 1
endif



service-docs:
	cd ../gftn-models \
		&& make
	cd ../api-service \
		&& make swaggergen && make docs
	cd ../anchor-service \
		&& make swaggergen && make docs
	cd ../quotes-service \
		&& make swaggergen && make docs
	cd ../participant-registry \
		&& make swaggergen && make docs
	cd ../administration-service \
		&& make swaggergen && make docs
	cd ../crypto-service \
		&& make swaggergen && make docs
	cd ../payout-service \
		&& make swaggergen && make docs
	cd ../global-whitelist-service \
		&& make swaggergen && make docs

build-go:
	cd ../gftn-models \
		&& make
	cd ../api-service \
		&& go build
	cd ../participant-registry \
		&& go build
	cd ../payment-listener \
		&& go build
	cd ../administration-service \
		&& go build
	cd ../crypto-service \
		&& go build
	cd ../fee-service \
		&& go build
	cd ../anchor-service \
		&& go build
	cd ../quotes-service \
		&& go build
	cd ../send-service \
		&& go build
	cd ../payout-service \
		&& go build
	cd ../global-whitelist-service/whitelistserver \
    	&& go build
	cd ../ww-gateway \
    	&& go build
docker:
	docker build -f src/main/docker/Dockerfile -t gftn/integration-tests .

push-dockers:
	@echo $(version)
	docker login -u ${DOCKER_USER} -p ${DOCKER_PASSWORD} ${DOCKER_REGISTRY}
	docker tag gftn/api-service:latest ${DOCKER_REGISTRY}/gftn/api-service:$(version)
	docker push ${DOCKER_REGISTRY}/gftn/api-service:$(version)
	docker tag gftn/crypto-service:latest ${DOCKER_REGISTRY}/gftn/crypto-service:$(version)
	docker push ${DOCKER_REGISTRY}/gftn/crypto-service:$(version)
	docker tag gftn/crypto-service-prod:latest ${DOCKER_REGISTRY}/gftn/crypto-service-prod:$(version)
	docker push ${DOCKER_REGISTRY}/gftn/crypto-service-prod:$(version)
	docker tag gftn/participant-registry:latest ${DOCKER_REGISTRY}/gftn/participant-registry:$(version)
	docker push ${DOCKER_REGISTRY}/gftn/participant-registry:$(version)
	docker tag gftn/quotes-service:latest ${DOCKER_REGISTRY}/gftn/quotes-service:$(version)
	docker push ${DOCKER_REGISTRY}/gftn/quotes-service:$(version)
	docker tag gftn/administration-service:latest ${DOCKER_REGISTRY}/gftn/administration-service:$(version)
	docker push ${DOCKER_REGISTRY}/gftn/administration-service:$(version)
	docker tag gftn/payment-listener:latest ${DOCKER_REGISTRY}/gftn/payment-listener:$(version)
	docker push ${DOCKER_REGISTRY}/gftn/payment-listener:$(version)
	docker tag gftn/anchor-service:latest ${DOCKER_REGISTRY}/gftn/anchor-service:$(version)
	docker push ${DOCKER_REGISTRY}/gftn/anchor-service:$(version)
	docker tag gftn/gas-service:latest ${DOCKER_REGISTRY}/gftn/gas-service:$(version)
	docker push ${DOCKER_REGISTRY}/gftn/gas-service:$(version)
	docker tag gftn/payout-service:latest ${DOCKER_REGISTRY}/gftn/payout-service:$(version)
	docker push ${DOCKER_REGISTRY}/gftn/payout-service:$(version)
	docker tag gftn/send-service:latest ${DOCKER_REGISTRY}/gftn/send-service:$(version)
	docker push ${DOCKER_REGISTRY}/gftn/send-service:$(version)
	docker tag gftn/global-whitelist-service:latest ${DOCKER_REGISTRY}/gftn/global-whitelist-service:$(version)
	docker push ${DOCKER_REGISTRY}/gftn/global-whitelist-service:$(version)
	docker tag gftn/fee-service:latest ${DOCKER_REGISTRY}/gftn/fee-service:$(version)
	docker push ${DOCKER_REGISTRY}/gftn/fee-service:$(version)
	docker tag gftn/ww-gateway:latest ${DOCKER_REGISTRY}/gftn/ww-gateway:$(version)
	docker push ${DOCKER_REGISTRY}/gftn/ww-gateway:$(version)
	docker tag gftn/automation-service:latest ${DOCKER_REGISTRY}/gftn/automation-service:$(version)
	docker push ${DOCKER_REGISTRY}/gftn/automation-service:$(version)
	docker tag gftn/auth-service:latest ${DOCKER_REGISTRY}/gftn/auth-service:$(version)
	docker push ${DOCKER_REGISTRY}/gftn/auth-service:$(version)
	
push-administration-service-dockers:
	@echo $(version)
	docker login -u ${DOCKER_USER} -p ${DOCKER_PASSWORD} ${DOCKER_REGISTRY}
	docker tag gftn/administration-service:latest ${DOCKER_REGISTRY}/gftn/administration-service:$(version)
	docker push ${DOCKER_REGISTRY}/gftn/administration-service:$(version)

push-auth-service-dockers:
	@echo $(version)
	docker login -u ${DOCKER_USER} -p ${DOCKER_PASSWORD} ${DOCKER_REGISTRY}
	docker tag gftn/auth-service:latest ${DOCKER_REGISTRY}/gftn/auth-service:$(version)
	docker push ${DOCKER_REGISTRY}/gftn/auth-service:$(version)

push-anchor-service-dockers:
	@echo $(version)
	docker login -u ${DOCKER_USER} -p ${DOCKER_PASSWORD} ${DOCKER_REGISTRY}
	docker tag gftn/anchor-service:latest ${DOCKER_REGISTRY}/gftn/anchor-service:$(version)
	docker push ${DOCKER_REGISTRY}/gftn/anchor-service:$(version)

push-api-service-dockers:
	@echo $(version)
	docker login -u ${DOCKER_USER} -p ${DOCKER_PASSWORD} ${DOCKER_REGISTRY}
	docker tag gftn/api-service:latest ${DOCKER_REGISTRY}/gftn/api-service:$(version)
	docker push ${DOCKER_REGISTRY}/gftn/api-service:$(version)

push-crypto-service-dockers:
	@echo $(version)
	docker login -u ${DOCKER_USER} -p ${DOCKER_PASSWORD} ${DOCKER_REGISTRY}
	docker tag gftn/crypto-service:latest ${DOCKER_REGISTRY}/gftn/crypto-service:$(version)
	docker push ${DOCKER_REGISTRY}/gftn/crypto-service:$(version)
	docker tag gftn/crypto-service-prod:latest ${DOCKER_REGISTRY}/gftn/crypto-service-prod:$(version)
	docker push ${DOCKER_REGISTRY}/gftn/crypto-service-prod:$(version)

push-fee-service-dockers:
	@echo $(version)
	docker login -u ${DOCKER_USER} -p ${DOCKER_PASSWORD} ${DOCKER_REGISTRY}
	docker tag gftn/fee-service:latest ${DOCKER_REGISTRY}/gftn/fee-service:$(version)
	docker push ${DOCKER_REGISTRY}/gftn/fee-service:$(version)

push-gas-service-dockers:
	@echo $(version)
	docker login -u ${DOCKER_USER} -p ${DOCKER_PASSWORD} ${DOCKER_REGISTRY}
	docker tag gftn/gas-service:latest ${DOCKER_REGISTRY}/gftn/gas-service:$(version)
	docker push ${DOCKER_REGISTRY}/gftn/gas-service:$(version)

push-global-whitelist-service-dockers:
	@echo $(version)
	docker login -u ${DOCKER_USER} -p ${DOCKER_PASSWORD} ${DOCKER_REGISTRY}
	docker tag gftn/global-whitelist-service:latest ${DOCKER_REGISTRY}/gftn/global-whitelist-service:$(version)
	docker push ${DOCKER_REGISTRY}/gftn/global-whitelist-service:$(version)

push-participant-registry-dockers:
	@echo $(version)
	docker login -u ${DOCKER_USER} -p ${DOCKER_PASSWORD} ${DOCKER_REGISTRY}
	docker tag gftn/participant-registry:latest ${DOCKER_REGISTRY}/gftn/participant-registry:$(version)
	docker push ${DOCKER_REGISTRY}/gftn/participant-registry:$(version)

push-payment-listener-dockers:
	@echo $(version)
	docker login -u ${DOCKER_USER} -p ${DOCKER_PASSWORD} ${DOCKER_REGISTRY}
	docker tag gftn/payment-listener:latest ${DOCKER_REGISTRY}/gftn/payment-listener:$(version)
	docker push ${DOCKER_REGISTRY}/gftn/payment-listener:$(version)

push-payout-service-dockers:
	@echo $(version)
	docker login -u ${DOCKER_USER} -p ${DOCKER_PASSWORD} ${DOCKER_REGISTRY}
	docker tag gftn/payout-service:latest ${DOCKER_REGISTRY}/gftn/payout-service:$(version)
	docker push ${DOCKER_REGISTRY}/gftn/payout-service:$(version)

push-quotes-service-dockers:
	@echo $(version)
	docker login -u ${DOCKER_USER} -p ${DOCKER_PASSWORD} ${DOCKER_REGISTRY}
	docker tag gftn/quotes-service:latest ${DOCKER_REGISTRY}/gftn/quotes-service:$(version)
	docker push ${DOCKER_REGISTRY}/gftn/quotes-service:$(version)

push-send-service-dockers:
	@echo $(version)
	docker login -u ${DOCKER_USER} -p ${DOCKER_PASSWORD} ${DOCKER_REGISTRY}
	docker tag gftn/send-service:latest ${DOCKER_REGISTRY}/gftn/send-service:$(version)
	docker push ${DOCKER_REGISTRY}/gftn/send-service:$(version)


push-ww-gateway-dockers:
	@echo $(version)
	docker login -u ${DOCKER_USER} -p ${DOCKER_PASSWORD} ${DOCKER_REGISTRY}
	docker tag gftn/ww-gateway:latest ${DOCKER_REGISTRY}/gftn/ww-gateway:$(version)
	docker push ${DOCKER_REGISTRY}/gftn/ww-gateway:$(version)

push-automation-service-dockers:
	@echo $(version)
	docker login -u ${DOCKER_USER} -p ${DOCKER_PASSWORD} ${DOCKER_REGISTRY}
	docker tag gftn/automation-service:latest ${DOCKER_REGISTRY}/gftn/automation-service:$(version)
	docker push ${DOCKER_REGISTRY}/gftn/automation-service:$(version)

push-msk-cli-dockers:
	@echo $(version)
	docker login -u ${DOCKER_USER} -p ${DOCKER_PASSWORD} ${DOCKER_REGISTRY}
	docker tag gftn/msk-cli:latest ${DOCKER_REGISTRY}/gftn/msk-cli:$(version)
	docker push ${DOCKER_REGISTRY}/gftn/msk-cli:$(version)

push-msk-monitoring-server-dockers:
	@echo $(version)
	docker login -u ${DOCKER_USER} -p ${DOCKER_PASSWORD} ${DOCKER_REGISTRY}
	docker tag gftn/msk-burrow:latest ${DOCKER_REGISTRY}/gftn/msk-burrow:$(version)
	docker push ${DOCKER_REGISTRY}/gftn/msk-burrow:$(version)

push-msk-monitoring-ui-dockers:
	@echo $(version)
	docker login -u ${DOCKER_USER} -p ${DOCKER_PASSWORD} ${DOCKER_REGISTRY}
	docker tag gftn/msk-burrow-ui:latest ${DOCKER_REGISTRY}/gftn/msk-burrow-ui:$(version)
	docker push ${DOCKER_REGISTRY}/gftn/msk-burrow-ui:$(version)