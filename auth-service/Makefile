FILE=../VERSION
VARIABLE=`cat $(FILE)`
export VERSION=${VARIABLE} created on: `date +'%y-%m-%d::%H:%M:%S'`

docker:
	
	# install dependencies
	bash install_dependencies.sh --gcloud

	# run unit tests
	# sh run-tests.sh

	# create images
	bash ./deployment/docker/create-images.sh --version ${VARIABLE}