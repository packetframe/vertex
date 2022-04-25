FROM debian:bullseye
RUN apt update && apt install -y ca-certificates
COPY api /usr/bin/vertex
COPY vertexctl /usr/bin/vertexctl
ENTRYPOINT ["/usr/bin/vertex"]