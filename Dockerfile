FROM golang:latest

LABEL Authors="alihjmm, 7abib04, Mohamed-Alasfoor, Hujafaar"
LABEL Description="Container"
LABEL Version="Latest"

EXPOSE 8080

WORKDIR /app

RUN apt-get update && apt-get install -y sqlite3 libsqlite3-dev

COPY . .

RUN go build -o forum .

CMD ["./forum"]