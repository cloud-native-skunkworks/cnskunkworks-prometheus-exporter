FROM golang:alpine3.14 as builder
RUN mkdir /src
ADD . /src/
WORKDIR /src
RUN go build -o cnskunkworks-prometheus-exporter
FROM alpine
COPY --from=builder /src/cnskunkworks-prometheus-exporter /app/cnskunkworks-prometheus-exporter
WORKDIR /app
ENTRYPOINT ["/app/cnskunkworks-prometheus-exporter"]