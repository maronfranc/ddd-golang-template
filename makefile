##### ##### Test ##### #####
.PHONY: test
test:
	@cd ./api && \
	go test -env=test -v ./... && \
	cd ..

##### ##### Database ##### #####
.PHONY: db-build-migration
build-migration:
	@cd ./api && \
	go build -o ../bin/go-example-migration ./infrastructure/migration-script/migrate.go && \
	cd ..
.PHONY: db-run-migration
run-migration: build-migration
	@cd ./api && \
	./bin/go-example-migration -env=development && \
	cd ..
.PHONY: docker-run-dev
docker-run-dev:
	@docker-compose -f ./docker/dev/docker-compose.yml up -d

##### ##### Development ##### #####
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
