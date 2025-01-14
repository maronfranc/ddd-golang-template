##### ##### Test ##### #####
.PHONY: test
test:
	@cd ./api && \
	go test -v ./... && \
	cd ..

##### ##### Database ##### #####
.PHONY: migration
build-migration:
	@go build -o ./bin/go-example-migration ./api/infrastructure/migration-script/migrate.go
.PHONY: run-migration
run-migration: build-migration
	@./bin/go-example-migration -env=development

##### ##### Development ##### #####
.PHONY: docker-run-dev
docker-run-dev:
	@docker-compose -f ./docker/dev/docker-compose.yml up -d
.PHONY: build-dev
build-dev:
	@cd ./api && \
	go build -o ../bin/go-example-dev ./main.go && \
	cd ..
.PHONY: run-dev
run-dev: build-dev
	@cd ./api && \
	../bin/go-example-dev -env=development && \
	cd ..
.PHONY: watch-dev
watch-dev:
	bash ./script/watch-dir.sh "./api" "make run-dev" 3000

##### ##### Production ##### #####
.PHONY: build-prod
build-prod:
	@cd ./api && \
	go build -tags=production -o ./bin/go-example-prod ./main.go && \
	cd ..
.PHONY: run-prod
run-prod: build-prod
	@./bin/go-example-prod -env=production
