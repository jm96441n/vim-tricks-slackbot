FROM alpine as builder
# build
ADD https://github.com/aws/aws-lambda-runtime-interface-emulator/releases/latest/download/aws-lambda-rie /usr/bin/aws-lambda-rie
RUN chmod 755 /usr/bin/aws-lambda-rie
COPY entry.sh /
RUN chmod 755 /entry.sh

FROM golang:latest as compiler
WORKDIR /function
COPY . .
RUN go build -o /main

FROM alpine as prd
WORKDIR /
COPY --from=compiler /main /main
COPY --from=builder /entry.sh /entry.sh
COPY --from=builder /usr/bin/aws-lambda-rie /usr/bin/aws-lambda-rie
ARG SLACK_TOKEN
ENV SLACK_TOKEN=$SLACK_TOKEN

ARG VIM_CHANNEL_ID
ENV VIM_CHANNEL_ID=$VIM_CHANNEL_ID
ENTRYPOINT [ "/entry.sh" ]
CMD [ "/main" ]


FROM golang:latest as test
WORKDIR /vim-tricks/function
ADD . .
CMD [ "go", "test", "./..." ]
