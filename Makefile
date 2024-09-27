.DEFAULT_GOAL := help

# Change these variables as necessary.
main_package_path = ./cmd/serve
binary_name = binary

# Can use `go run honnef.co/go/tools/cmd/staticcheck@latest`
static_check = staticcheck

# Can use `go run golang.org/x/vuln/cmd/govulncheck@latest`
govulncheck = govulncheck

# Can use `go run github.com/golangci/golangci-lint@latest`
golangci = golangci-lint

# Can use `go fmt`
fmt = gofumpt -l -w

# Can use `go run github.com/cosmtrek/air@v1.43.0`
air = air

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

.PHONY: no-dirty
no-dirty:
	@test -z "$(shell git status --porcelain)"

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: run quality control checks
.PHONY: audit
audit: test
	@go mod tidy -diff
	@go mod verify
	@test -z "$(shell gofmt -l .)"
	@go vet ./...
	@${govulncheck} ./...
	@${golangci} run --fix --issues-exit-code 0 ./...
	@echo Audit completed.

# Add this to the command above to use staticcheck
# ${static_check} -checks=all,-ST1000,-U1000 ./...

## test: run all tests
.PHONY: test
test:
	@go test -v -race -buildvcs ./...
	@echo Tests completed.

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	@go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	@go tool cover -html=/tmp/coverage.out
	@echo Tests completed.

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## tidy: tidy modfiles and format .go files
.PHONY: tidy
tidy:
	@go mod tidy -v
	@${fmt} .
	@echo Tidy completed.

## build: build the application
.PHONY: build
build:
	@# Include additional build steps, like TypeScript, SCSS or Tailwind compilation here...
	@go build -o=/tmp/bin/${binary_name} ${main_package_path}
	@echo Build completed.

## run: run the  application
.PHONY: run
run: build
	@/tmp/bin/${binary_name}
	@echo Application completed.

## run/live: run the application with reloading on file changes
.PHONY: run/live
run/live:
	@${air} \
	--build.cmd "make build" --build.bin "/tmp/bin/${binary_name}" --build.delay "100" \
	--build.exclude_dir "" \
	--build.include_ext "go, tpl, tmpl, html, css, scss, js, ts, sql, jpeg, jpg, gif, png, bmp, svg, webp, ico" \
	--misc.clean_on_exit "true"
	@echo Air exited.

## pprof: run a profiler on the running server at port 3000
.PHONY:pprof
pprof:
	@go tool pprof -http=localhost:3000 default.pgo
	@echo Pprof completed.

# ==================================================================================== #
# OPERATIONS
# ==================================================================================== #

## push: push changes to the remote Git repository
.PHONY: push
push: confirm audit no-dirty
	@git push
	@echo Pushed.

## production/deploy: deploy the application to production
.PHONY: production/deploy
production/deploy: confirm audit no-dirty
	@GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=/tmp/bin/linux_amd64/${binary_name} ${main_package_path}
	@upx -5 /tmp/bin/linux_amd64/${binary_name}
	@# Include additional deployment steps here...
	@echo Deploy completed.
