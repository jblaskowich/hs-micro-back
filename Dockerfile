FROM golang:1.11.5

LABEL maintainer="jblaskowichgmail.com"

WORKDIR /

RUN go get -v -d github.com/go-sql-driver/mysql
RUN go get -v -d github.com/nats-io/go-nats

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM scratch  

COPY --from=0 app /

ENTRYPOINT ["/app"]