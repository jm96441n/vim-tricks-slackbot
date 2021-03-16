.PHONY: build
build:
	docker build -t jmaguire5588/vim-tricks-bot-lambda:latest ./

.PHONY: run
run:
	$(MAKE) build
	docker run --rm \
		--name vim-tricks \
		-p 9000:8080 \
		jmaguire5588/vim-tricks-bot-lambda:latest

.PHONY: test
test:
	curl -XPOST "http://localhost:9000/2015-03-31/functions/function/invocations" -d '{}'

