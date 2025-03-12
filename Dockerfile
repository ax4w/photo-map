
    FROM golang:1.24-alpine AS build

    WORKDIR /app
    
    COPY go.mod go.sum ./
    RUN go mod download
    RUN go install github.com/a-h/templ/cmd/templ@latest

    COPY . .

    RUN templ generate
    RUN go build -o server .

    FROM alpine:latest

    RUN addgroup -S appgroup && adduser -S appuser -G appgroup

    RUN apk add --no-cache file
    RUN apk --update add imagemagick
    RUN apk add --no-cache imagemagick imagemagick-libs libjpeg-turbo
    
    WORKDIR /app
    
    COPY --from=build /app/server /app/
    
    COPY --from=build /app/frontend/index.html /app/
    
    EXPOSE 8080
    
    USER appuser

    CMD ["./server"]
    