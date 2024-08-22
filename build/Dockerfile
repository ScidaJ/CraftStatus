# syntax=docker/dockerfile:1
FROM golang:latest AS build
WORKDIR /app
COPY go.mod go.sum main.go /app/
COPY internal/ /app/internal
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/bot ./main.go

FROM alpine:latest AS certs
RUN apk --update add ca-certificates

FROM scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /bin/bot /bin/bot
EXPOSE 8080
CMD [ "/bin/bot" ]