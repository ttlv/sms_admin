SHELL = /bin/sh

build:
	@GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -tags 'bindatafs' -o sms-admin main.go
	@docker build -t gopherlv/sms-admin .
	@rm sms-admin

push: build
    @$(eval REV := $(shell git rev-parse HEAD|cut -c 1-8))
    @docker tag gopherlv/sms-admin gopherlv/sms_service:sms-admin-$(REV)
    @docker push gopherlv/sms_service:sms-admin-$(REV)
