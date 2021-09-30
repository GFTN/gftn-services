swaggergen:
	rm -rf model/
	go get -u github.com/go-swagger/go-swagger/cmd/swagger
	swagger generate model --with-flatten=full --spec=../api-definition/internal-api.yaml --model-package=model --target=.