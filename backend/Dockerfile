FROM golang:1.19 as build
WORKDIR /go/src
COPY go.mod go.sum ./
RUN go mod download

ADD . /backend
WORKDIR /backend
RUN go build -o bin ../backend/app/main.go

FROM gcr.io/distroless/base-debian11
COPY --from=build /backend/bin /app
EXPOSE 8080
CMD ["/app"]