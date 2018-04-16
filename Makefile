VERSION ?= latest
image:
	docker build -t richterrettich/tmpl:${VERSION} .

build:
	go build

