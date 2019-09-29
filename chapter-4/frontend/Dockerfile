FROM golang:1.12 AS gobuilder

WORKDIR /app
COPY . .

ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOARCH amd64
RUN go build -o main -mod vendor ./main.go

FROM gcr.io/distroless/base
WORKDIR /app
COPY --from=gobuilder /app/main /app/main
COPY --from=gobuilder /app/web /app/web

ENTRYPOINT ["./main"]

