all: swaggergen

swaggergen:
	mv \
	model/asset.go \
  model/participant.go \
	model/sweep.go \
	model/quote.go \
	model/quote_request.go \
	model/exchange.go \
	model/exchange_receipt.go \
	./ || true
	rm -rf model/
	go get -u github.com/go-swagger/go-swagger/cmd/swagger
	swagger generate model --with-flatten=full --spec=model-definitions/gftn-models.yaml --model-package=model --target=.
	mv \
	asset.go \
  participant.go \
	sweep.go \
	quote.go \
	quote_request.go \
	exchange.go \
	exchange_receipt.go \
	model/

swaggergenasset:
	swagger generate model --with-flatten=full --spec=model-definitions/gftn-asset-model.yaml --model-package=model --target=.

swaggergenparticipant:
	swagger generate model --with-flatten=full --spec=model-definitions/gftn-participant-model.yaml --model-package=model --target=.


glide:
	curl https://glide.sh/get | sh
	glide install

build:
	go build

install:
	go install
