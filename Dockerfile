FROM golang:1.24.3
WORKDIR /app
COPY . .
RUN go mod tidy && go mod vendor
RUN go build -o main .
EXPOSE 8080
CMD ["./main"]