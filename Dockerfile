FROM alpine:latest

ENV APP_NAME feedthembot
RUN apk add --no-cache git make musl-dev go

ENV GOPATH /go
ENV PATH /go/bin:$PATH

COPY ./src ${GOPATH}/src/${APP_NAME}
WORKDIR ${GOPATH}/src/${APP_NAME}

RUN apk add --no-cache ca-certificates &&\
    chmod +x .
RUN go get .
RUN go build -o ${APP_NAME}

EXPOSE 8080
CMD ./${APP_NAME}