package main

const TEMPL_MAKE = `
default:
	@godep save
	@godep go build
	@ls -ltrh

setup: .goxc.ok
	@echo Installing dependency management tools...
	go get github.com/tools/godep

.goxc.ok:
	@echo Installing crossbuild tooling. This will take a while...
	go get github.com/laher/goxc
	goxc -t
	touch .goxc.ok

heroku:
	@echo Bootstrapping with godep
	@go get github.com/tools/godep
	@godep save
	@git add -A .
	@git commit -am "dependencies"
	@echo Creating a Heroku Go app...
	@heroku create -b https://github.com/kr/heroku-buildpack-go.git
	@git push heroku master

test:
	@godep go test

bump:
	@goxc bump

release:
	godep save
	goxc -env GOPATH=` + "`godep path`" + ` -bc="linux,amd64" -d . xc # we only use basic xc for now, see github.com/laher/goxc for more

docker: release
	@docker build -t {{.Product}} .
	@echo Container [{{.Product}}] built. Run with: make docker-run

docker-run:
	docker run -p 80:6060 {{.Product}}

.PHONY: heroku build test setup release docker docker-run bump
`
