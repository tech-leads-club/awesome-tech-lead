FROM golang:1.24

WORKDIR /app

COPY . /app

RUN make setup

CMD ["make", "generate-readme"]