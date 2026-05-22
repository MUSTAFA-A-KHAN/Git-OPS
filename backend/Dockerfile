FROM golang:1.23-alpine AS build
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o server ./cmd/server
FROM alpine:3.20
WORKDIR /app
COPY --from=build /app/server .
EXPOSE 8080
CMD ["./server"]
