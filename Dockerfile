FROM golang:1.18.3-alpine3.16 AS development
ENV GO111MODULE=on \
    CGO_ENABLED=1  \
    GOARCH=amd64 \
    GOOS=linux

ARG TIMESTAMP
ARG HASH_VALUE
ENV BUILD_TIME=${TIMESTAMP}
ENV COMMIT_HASH=${HASH_VALUE}

WORKDIR /app
COPY . .
RUN go mod download && go mod tidy -go=1.18
EXPOSE 8000 9000
HEALTHCHECK --interval=5m --timeout=3s CMD curl --fail http://localhost:8000/ || exit 1
CMD ["go", "run", "./cmd"]

FROM golang:1.18.3-alpine3.16 AS build
ENV GO111MODULE=on \
    CGO_ENABLED=1  \
    GOARCH=amd64 \
    GOOS=linux

COPY --from=development /app/ /app/
WORKDIR  /app/cmd
RUN go build -o app

FROM alpine:3.16 AS production
ARG TIMESTAMP
ARG HASH_VALUE
ENV BUILD_TIME=${TIMESTAMP}
ENV COMMIT_HASH=${HASH_VALUE}

COPY --from=build /app/cmd/app /usr/local/app
EXPOSE 8000 9000
USER nobody:nobody

HEALTHCHECK --interval=5m --timeout=3s CMD curl --fail http://localhost:8000/ || exit 1
ENTRYPOINT ["/usr/local/app"]

