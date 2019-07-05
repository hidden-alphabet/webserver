FROM golang:alpine

WORKDIR /app

COPY server/ server/
COPY main.go ./

RUN apk add git

RUN go get -u -v github.com/gorilla/handlers && \
      go get -u -v github.com/gorilla/mux && \
      go get -u -v github.com/joho/godotenv && \
      go get -u -v github.com/joho/godotenv/autoload && \
      go get -u -v github.com/lib/pq && \
      go get -u -v github.com/lib/pq/oid && \
      go get -u -v github.com/lib/pq/scram && \
      go get -u -v github.com/satori/go.uuid && \
      go get -u -v golang.org/x/crypto/argon2 && \
      go get -u -v golang.org/x/crypto/blake2b && \
      go get -u -v golang.org/x/sys/cpu

RUN apk del git

RUN mkdir -p bin && \
    go build -o bin/main main.go

ENTRYPOINT ["./bin/main"]
