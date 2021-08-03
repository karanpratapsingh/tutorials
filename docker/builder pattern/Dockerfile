FROM golang:1.16.5
WORKDIR /app
COPY go.mod .
COPY main.go .
RUN go build -o app
EXPOSE 8080
CMD [ "./app" ]
