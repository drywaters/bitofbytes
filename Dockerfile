FROM node:latest as tailwind-builder
WORKDIR /templates
RUN npm init -y && \
    npm install tailwindcss @tailwindcss/typography
COPY ./templates ./templates
COPY ./tailwind/styles.css src/styles.css
RUN npx tailwindcss -i src/styles.css -o /styles.css --minify

FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -v -o ./bob ./cmd/bob/

FROM alpine
WORKDIR /bob
COPY ./static ./static
COPY .env .env
COPY --from=builder /app/bob ./bob
COPY --from=tailwind-builder /styles.css ./static/styles.css

EXPOSE 3000
CMD ./bob
