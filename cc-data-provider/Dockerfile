FROM golang:1.16-alpine AS build
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN go build -o /cc-data-provider

FROM golang:1.16-alpine AS run
WORKDIR /app

COPY --from=build /cc-data-provider /cc-data-provider
COPY .env .env

EXPOSE 3000

CMD [ "/cc-data-provider" ]