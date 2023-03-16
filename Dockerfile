FROM node:alpine as tailwind
WORKDIR /usr/src/ktn
COPY . .
RUN npm install -g tailwindcss
RUN npx tailwindcss -i ./web/templates/input.css -o ./web/static/main.css

FROM golang:1.20-alpine as builder
WORKDIR /usr/src/ktn
COPY . .
RUN go build

FROM alpine
COPY --from=builder /usr/src/ktn/ktn-go /usr/local/bin/ktn
COPY ./web/static ./web/static
COPY ./web/templates ./web/templates
COPY --from=tailwind /usr/src/ktn/web/static/main.css ./web/static/

EXPOSE 8080
EXPOSE 2525

ENV GIN_MODE release

CMD ["ktn"]
