FROM golang:1.25.1-alpine AS build
RUN apk add --no-cache curl libstdc++ libgcc

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go install github.com/a-h/templ/cmd/templ@latest && \
    templ generate && \
    curl -sL https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64-musl -o tailwindcss && \
    chmod +x tailwindcss && \
    ./tailwindcss -i internal/server/styles/input.css -o internal/server/assets/css/output.css

RUN go build -o main cmd/api/main.go

FROM alpine:3.22.1 AS prod
LABEL org.opencontainers.image.source=https://github.com/jacobsky/meishi
WORKDIR /app
COPY --from=build /app/main /app/main
EXPOSE ${PORT}
CMD ["./main"]


