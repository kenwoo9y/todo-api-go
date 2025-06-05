FROM golang:1.23-bullseye as deploy-builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -trimpath -ldflags "-w -s" -o app ./api/cmd

# ---------------------------------------------------

FROM debian:bullseye-slim as deploy

RUN apt-get update

COPY --from=deploy-builder /app/app .

EXPOSE 8080
CMD ["./app"]

# ---------------------------------------------------

FROM golang:1.23 as dev
WORKDIR /app
ENV GOPATH=/root/go
ENV PATH=/usr/local/go/bin:$GOPATH/bin:$PATH
RUN go install github.com/sqldef/sqldef/cmd/mysqldef@latest
RUN go install github.com/sqldef/sqldef/cmd/psqldef@latest
WORKDIR /app/api
CMD ["go", "run", "./cmd"]