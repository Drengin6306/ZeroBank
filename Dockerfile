FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download

COPY . .

ARG SERVICE_PATH
ARG SERVICE_NAME

RUN cd ${SERVICE_PATH} && go build -ldflags="-s -w" -o /app/service .

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

WORKDIR /app

ARG SERVICE_PATH
ARG SERVICE_PORT

COPY --from=builder /app/service /app/service
COPY ${SERVICE_PATH}/etc /app/etc

EXPOSE ${SERVICE_PORT}

ENTRYPOINT ["/app/service"]