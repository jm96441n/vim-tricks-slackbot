FROM golang:1.16.2-buster as setup
WORKDIR /vim-tricks/function
# build
ADD . .
RUN go build -o /main
ENTRYPOINT [ "/aws-lambda/aws-lambda-rie" ]
CMD [ "/main" ]

FROM golang:1.16.2-buster as test
WORKDIR /vim-tricks/function
ADD . .
CMD [ "go", "test", "./..." ]
