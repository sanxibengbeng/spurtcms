run:
	docker-compose down    
	docker-compose build
	docker-compose up -d

buildwithview:
	mkdir build
	cp -r lang build
	cp -r locales build
	cp -r public build
	cp -r view build
	cp -r cms.sql build
	CGO_ENABLED=0  GOOS=linux GOARCH=amd64 go build -o build/spurtcms-admin

build:
	CGO_ENABLED=0  GOOS=linux GOARCH=amd64 go build -o build/spurtcms-admin

permission:
	chmod -R 777 spurtcms-admin

start:
	sudo systemctl start spurtcms-admin

test:
	go test ./... -v

test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

test-short:
	go test ./... -short

clean:
	rm -f spurtcms-admin
	rm -f coverage.out

.PHONY: run buildwithview build permission start test test-coverage test-short clean