REPOSITORY ?= ip-team-worldwire-docker-local.artifactory.swg-devops.com

push-msk-cli:
	docker push ${REPOSITORY}/gftn/msk-cli:latest

build-msk-cli:
	docker build ./msk-cli/ -t gftn/msk-cli:latest
# 	docker tag gftn/msk-cli:latest ${REPOSITORY}/gftn/msk-cli:latest  || exit 1

push-msk-burrow:
	docker push ${REPOSITORY}/gftn/msk-burrow:latest

build-msk-burrow:
	docker build ./msk-monitoring/ -t gftn/msk-burrow:latest
# 	docker tag gftn/msk-burrow:latest ${REPOSITORY}/gftn/msk-burrow:latest  || exit 1

push-msk-burrow-ui:
	docker push ${REPOSITORY}/gftn/msk-burrow-ui:latest

build-msk-burrow-ui:
	docker build ./msk-monitoring/dashboard -t gftn/msk-burrow-ui:latest
# 	docker tag gftn/msk-burrow-ui:latest ${REPOSITORY}/gftn/msk-burrow-ui:latest  || exit 1