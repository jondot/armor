
default:
	@go build
	@ls -ltrh

setup:
	@echo Installing build tooling...
	go get github.com/laher/goxc
	go get github.com/tools/godep

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

release:
	godep save
	goxc -bc="linux,amd64" -d . xc # we only use basic xc for now, see github.com/laher/goxc for more

.PHONY: heroku build test setup
