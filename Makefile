CONTAINER_NAME=vim-tricks
.PHONY: build
build:
	docker build --target=prd -t jmaguire5588/vim-tricks-bot-lambda:latest ./

.PHONY: run
run:
	$(MAKE) build
	docker run --rm \
		--name $(CONTAINER_NAME) \
		-p 9000:8080 \
		jmaguire5588/vim-tricks-bot-lambda:latest

.PHONY: test
test:
	docker build --target=test -t jmaguire5588/vim-tricks-bot-lambda:test ./
	docker run --rm jmaguire5588/vim-tricks-bot-lambda:test

.PHONY: curl_test
curl_test:
	curl -XPOST "http://localhost:9000/2015-03-31/functions/function/invocations" -d '{}'

.PHONY: setup
setup:
	mkdir -p ~/.aws-lambda-rie && curl -Lo ~/.aws-lambda-rie/aws-lambda-rie \
		https://github.com/aws/aws-lambda-runtime-interface-emulator/releases/latest/download/aws-lambda-rie \
		&& chmod +x ~/.aws-lambda-rie/aws-lambda-rie

.PHONY: stop
stop:
	docker stop $(CONTAINER_NAME)
