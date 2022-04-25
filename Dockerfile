FROM debian:bullseye
RUN apt update && apt install -y ca-certificates
COPY api /usr/bin/vertex
ENTRYPOINT ["/usr/bin/vertex"]