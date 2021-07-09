# =============================================================================
# Target: base
FROM golang:1.16-buster as base

EXPOSE 8181

RUN curl -fsSL https://deb.nodesource.com/setup_16.x | bash -

RUN apt-get install -y nodejs && \
    npm install -g npm             # update to latest NPM

RUN npm install -g \
    sass \
    sass-lint

WORKDIR /go/src/app
COPY . .

RUN go mod download

# TODO: use magefile?

RUN sass --no-quiet --stop-on-error resources/css/main.css scss/main.scss

# TODO: figure out why tests hang in container
#RUN go test ./...
RUN go install

ENTRYPOINT ["adler"]
