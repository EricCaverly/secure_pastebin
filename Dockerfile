FROM golang:1.24.5

WORKDIR /app

COPY ./app ./

RUN go mod download

RUN go build -o spb

ENTRYPOINT [ "./spb" ]