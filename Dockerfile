FROM golang:1.16.2-alpine3.13 as builder
WORKDIR /function
# build
ADD https://github.com/aws/aws-lambda-runtime-interface-emulator/releases/latest/download/aws-lambda-rie /usr/bin/aws-lambda-rie
COPY . .
RUN chmod 755 /usr/bin/aws-lambda-rie
COPY entry.sh /
RUN chmod 755 /entry.sh
RUN go build -o /main

FROM alpine as prd
COPY --from=builder /main /main
COPY --from=builder /entry.sh /entry.sh
COPY --from=builder /usr/bin/aws-lambda-rie /usr/bin/aws-lambda-rie
ENTRYPOINT [ "/entry.sh" ]
CMD [ "/main" ]


#FROM golang:1.16.2-buster as test
#WORKDIR /vim-tricks/function
#ADD . .
#CMD [ "go", "test", "./..." ]
