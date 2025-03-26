FROM node:18-alpine AS frontend-build

WORKDIR /app/frontend-react
COPY frontend-react/package.json frontend-react/package-lock.json* ./
RUN npm install
COPY frontend-react/ ./
RUN npm run build

FROM golang:1.24-alpine AS backend-build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=frontend-build /app/frontend-react/dist /app/frontend-react/dist

RUN go build -o server .

FROM alpine:latest

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

RUN apk add --no-cache file
RUN apk --update add imagemagick
RUN apk add --no-cache imagemagick imagemagick-libs libjpeg-turbo

WORKDIR /app

COPY --from=backend-build /app/server /app/
COPY --from=frontend-build /app/frontend-react/dist /app/frontend-react/dist

EXPOSE 8080

USER appuser

CMD ["./server"]
    