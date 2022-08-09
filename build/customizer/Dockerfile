FROM node:14-stretch as build-front

WORKDIR /app

COPY . .

WORKDIR /app/web

RUN rm -rf node_modules \
    && yarn install \
    && yarn build

FROM golang:1.18 as build-server

WORKDIR /app

COPY --from=build-front /app ./

RUN go mod download \
    && go mod verify \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /output/customizer ./cmd/customizer

FROM alpine:latest

WORKDIR /app

COPY --from=build-server /app/configs/config.toml.example /app/config.toml
COPY --from=build-server /app/build/customizer/entrypoint.sh /entrypoint.sh
COPY --from=build-server /output/customizer ./customizer
COPY --from=build-front /app/web/dist ./web/dist

RUN chmod 755 ./customizer
RUN chmod +x ./customizer

RUN chmod 755 /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]

CMD ["/app/customizer", "-c", "/app/config.toml"]