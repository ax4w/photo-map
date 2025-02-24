
    FROM golang:1.24-alpine AS build

    WORKDIR /app
    
    COPY go.mod ./
    RUN go mod download

    COPY . .

    RUN go build -o server main.go

    FROM alpine:latest

    RUN addgroup -S appgroup && adduser -S appuser -G appgroup
    
    WORKDIR /app
    
    COPY --from=build /app/server /app/
    
    COPY --from=build /app/index.html /app/
    
    EXPOSE 8080
    
    USER appuser

    CMD ["./server"]
    