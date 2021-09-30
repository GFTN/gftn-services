


docker:
	cd ../utility/aws/javascript/ && npm install && npm run build
	mkdir -p ./utility/aws/
	rsync -rupE ../utility/aws/javascript ./utility/aws/ >> /dev/null
	docker login -u ${DOCKER_USER} -p ${DOCKER_PASSWORD} ${DOCKER_REGISTRY}
	docker pull ${DOCKER_REGISTRY}/gftn/node-alpine:latest
	docker tag ${DOCKER_REGISTRY}/gftn/node-alpine:latest gftn/node-alpine:latest
	docker images
	docker build -f deployment/Dockerfile -t gftn/gas-service .
install:
	npm install
test:
	npm run test \
	npm run stress-test \

