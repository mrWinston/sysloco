FROM golang:1.12 AS dev

WORKDIR /receiver
RUN ["go", "get", "github.com/githubnemo/CompileDaemon"]
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN GOOS=linux go build -a -o /app .

ENTRYPOINT CompileDaemon -log-prefix=false -build="go build -a -o /app" -command="/app -db-engine sqlite -db-loc /db/db.txt -v 2"

# ENTRYPOINT ["go", "run"]

FROM golang:1.12-alpine3.9 as prod
RUN apk add sqlite pcre
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
COPY --from=dev /app /app
ENTRYPOINT ["/app"]
