FROM node:20-alpine AS frontend-build

WORKDIR /app/frontend
COPY frontend/package.json frontend/package-lock.json* ./
RUN npm install
COPY frontend/ ./
RUN npm run build

FROM golang:1.24-alpine AS backend-build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=frontend-build /app/frontend/dist /app/frontend/dist

RUN go build -o server .

FROM alpine:latest

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

RUN apk add --no-cache file
RUN apk --update add imagemagick
RUN apk add --no-cache imagemagick imagemagick-libs libjpeg-turbo

WORKDIR /app

COPY --from=backend-build /app/server /app/
COPY --from=frontend-build /app/frontend/dist /app/frontend/dist

EXPOSE 8080

USER appuser

CMD ["./server"]
    