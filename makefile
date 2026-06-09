API_PATH="./api"
# Expected path relative to `./api`.
MIGRATION_PATH="./infrastructure/migration-script/migrate.go"

BUILD_MIGRATION="go-example-migration"
BUILD_DEV="go-example-dev"
BUILD_PROD="go-example-prod"
BUILD_PATH="bin"

# ===== ===== Test  ===== ===== #
.PHONY: test
test:
	@cd ${API_PATH} && \
	ENV=test go test -v ./... && \
	cd ..

# ===== ===== Database ===== ===== #
.PHONY: db-build-migration
db-build-migration:
	@cd ${API_PATH} && \
	go build -o ../${BUILD_PATH}/${BUILD_MIGRATION} ${MIGRATION_PATH} && \
	cd ..

.PHONY: db-run-migration
db-run-migration: db-build-migration
	@cd ${API_PATH} && \
	ENV=development ../${BUILD_PATH}/${BUILD_MIGRATION} up && \
	cd ..

.PHONY: db-run-migration-down
db-run-migration-down: db-build-migration
	@cd ${API_PATH} && \
	ENV=development ../${BUILD_PATH}/${BUILD_MIGRATION} down && \
	cd ..

# ===== ===== Development ===== ===== #
.PHONY: build-dev
build-dev:
	@cd ${API_PATH} && \
	go build -o ../${BUILD_PATH}/${BUILD_DEV} ./main.go && \
	cd ..

.PHONY: run-dev
run-dev: build-dev
	@cd ${API_PATH} && \
	ENV=development ../${BUILD_PATH}/${BUILD_DEV} && \
	cd ..

.PHONY: watch-dev
watch-dev:
	bash ./script/watch-dir.sh ${API_PATH} "make run-dev" 3000

.PHONY: docker-run-dev
docker-run-dev:
	@docker-compose -f ./docker/dev/docker-compose.yml up -d

# ===== ===== Production ===== ===== #
.PHONY: build-prod
build-prod:
	@cd ${API_PATH} && \
	go build -tags=production -o ../${BUILD_PATH}/${BUILD_PROD} ./main.go && \
	cd ..

.PHONY: run-prod
run-prod: build-prod
	@ENV=production ../${BUILD_PATH}/${BUILD_PROD}
