FROM node:20-alpine AS tailwind-builder
ARG TAILWIND_VERSION=4.1.13
WORKDIR /app
RUN npm init -y \
    && npm install tailwindcss@${TAILWIND_VERSION} @tailwindcss/typography
COPY ./templates ./templates
COPY ./tailwind/styles.css src/styles.css
RUN npx @tailwindcss/cli -i src/styles.css -o /styles.css --minify

FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -v -o ./bob ./cmd/bob/

FROM alpine
ARG LOG_LEVEL=warn
WORKDIR /bob
COPY ./static ./static
COPY --from=builder /app/bob ./bob
COPY --from=tailwind-builder /styles.css ./static/styles.css
RUN addgroup -S bob \
    && adduser -S -G bob bob \
    && chown -R bob:bob /bob

ENV CSRF_KEY_FILE=/run/secrets/csrf_key
ENV SERVER_ADDRESS=:3000
ENV CSRF_SECURE=true
ENV LOG_LEVEL=${LOG_LEVEL}

USER bob

EXPOSE 3000
CMD ["./bob"]
