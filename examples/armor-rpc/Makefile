
default:
	@godep go build
	@ls -ltrh

setup:
	@echo Installing build tooling...
	go get github.com/tools/godep
	go get github.com/laher/goxc

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
	go test

bump:
	@goxc bump

release:
	godep save
	goxc -env GOPATH=`godep path` -bc="linux,amd64" -d . xc # we only use basic xc for now, see github.com/laher/goxc for more

docker: release
	@docker build -t armor-rpc .
	@echo Container [armor-rpc] built. Run with: make docker-run

docker-run:
	docker run -p 80:6060 armor-rpc

.PHONY: heroku build test setup release docker docker-run bump
