.PHONY: proto data build

proto:
	for d in api srv; do \
		for f in $$d/**/proto/*.proto; do \
			protoc --go_out=plugins=micro:. $$f; \
			echo compiled: $$f; \
		done \
	done

lint:
	./bin/lint.sh

build:
	./bin/build.sh

data:
	go-bindata -o data/bindata.go -pkg data data/*.json

run:
	docker-compose build
	docker-compose up
