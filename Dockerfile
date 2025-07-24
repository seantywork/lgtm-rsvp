FROM ubuntu:24.04

ARG DEBIAN_FRONTEND=noninteractive

RUN mkdir -p /workspace

WORKDIR /workspace

RUN apt-get update 

RUN apt-get install -y make build-essential ca-certificates sqlite3

COPY --from=golang:1.23.4 /usr/local/go/ /usr/local/go/

ENV PATH="/usr/local/go/bin:${PATH}"

COPY . . 

RUN make clean

RUN go clean -modcache

RUN go mod tidy

RUN	make

ENTRYPOINT ["/bin/bash", "./run.sh"]
CMD [""]