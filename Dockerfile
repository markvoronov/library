FROM golang:1.22-alpine AS builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash git make gcc gettext musl-dev postgresql-client

COPY ./ ./

# make wait-for-postgres.sh executable
#RUN chmod +x wait-for-postgres.sh

# build go app
RUN go mod download
RUN go build -o library-app ./cmd/main.go

CMD ["./library-app"]