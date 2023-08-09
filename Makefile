# Include variables from the .envrc file
include .envrc

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #
## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	go run ./cmd/api

# ==================================================================================== #
# BUILD
# ==================================================================================== #

## build/api: build the cmd/api application using ldflags on to strip debug info and decrease size
.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	go build -o=./bin/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/api ./cmd/api

## generates a swagger.yaml file
# install this first: go install github.com/go-swagger/go-swagger/cmd/swagger@latest
.PHONY: swagger
swagger: vendor
	swagger generate spec -o ./swagger.yaml --scan-models

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: tidy and vendor dependencies and format, vet and test all code
.PHONY: audit
audit: vendor
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'


# ==================================================================================== #
# PRODUCTION
# ==================================================================================== #
production_host_ip = '146.190.76.183'
## production/connect: connect to the production server
.PHONY: production/connect
production/connect:
	ssh rpg_manager@${production_host_ip}

## production/deploy/api: deploy the api to production
.PHONY: production/deploy/api
production/deploy/api:
	rsync -P ./bin/linux_amd64/api rpg_manager@${production_host_ip}:~
	rsync -P ./remote/production/api.service rpg_manager@${production_host_ip}:~
	rsync -P ./remote/production/Caddyfile rpg_manager@${production_host_ip}:~
	ssh -t rpg_manager@${production_host_ip} '\
		sudo mv ~/api.service /etc/systemd/system/ \
		&& sudo systemctl enable api \
		&& sudo systemctl restart api \
		&& sudo mv ~/Caddyfile /etc/caddy/ \
		&& sudo systemctl reload caddy \
	'
