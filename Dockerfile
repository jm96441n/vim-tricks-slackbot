FROM golang:1.16.2-buster as builder
WORKDIR /vim-tricks/function
# build
ADD . .
RUN go build -o /main

FROM alpine as prd
COPY --from=builder /main /main
ENTRYPOINT [ "/aws-lambda/aws-lambda-rie" ]
CMD [ "/main" ]

FROM golang:1.16.2-buster as test
WORKDIR /vim-tricks/function
ADD . .
CMD [ "go", "test", "./..." ]
