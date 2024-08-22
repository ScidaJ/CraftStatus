# syntax=docker/dockerfile:1
FROM golang:latest AS build
WORKDIR /app
COPY go.mod go.sum main.go /app/
COPY internal/ /app/internal

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/bot ./main.go

FROM scratch
EXPOSE 8080
COPY --from=build /bin/bot /bin/bot
CMD [ "/bin/bot" ]