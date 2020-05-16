FROM golang:1.14.3

ENV APP_NAME feedthembot
ENV PORT 8080

COPY ./src /go/src/${APP_NAME}
ADD botsettings.json /go/src/${APP_NAME}
WORKDIR /go/src/${APP_NAME}

RUN go get .
RUN go build -o ${APP_NAME}

CMD ./${APP_NAME}

EXPOSE ${PORT}