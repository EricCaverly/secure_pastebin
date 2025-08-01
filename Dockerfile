FROM golang:1.24.5

LABEL org.opencontainers.image.authors="eric@ericc.ninja"

WORKDIR /app

COPY ./app ./

RUN go mod download

RUN go build -o spb

ENTRYPOINT [ "./spb" ]