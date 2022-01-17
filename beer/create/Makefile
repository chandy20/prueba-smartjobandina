.PHONY: npmi build production staging squad

npmi:
	npm ci --prefer-offline --no-audit

build:
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/v1 v1/*.go

production: build npmi
	node_modules/.bin/serverless --stage production create_domain
	node_modules/.bin/serverless --stage production deploy

squad: build npmi
	node_modules/.bin/serverless --stage squad create_domain
	node_modules/.bin/serverless --stage squad deploy
