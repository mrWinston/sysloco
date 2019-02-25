FROM golang:1.11 AS dev

# Copy the code from the host and compile it
WORKDIR /receiver
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN GOOS=linux go build -a -o /app .

ENTRYPOINT ["go", "run"]

FROM golang:1.11.4-alpine3.8
RUN apk add sqlite pcre
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
COPY --from=dev /app /app 
ENTRYPOINT ["/app"]
