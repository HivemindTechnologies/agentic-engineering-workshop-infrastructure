# Built by `just images` from the git tree of setup/loadgen/, like scoring/Dockerfile.
FROM golang:1.26.10-alpine3.24@sha256:9ce231b28998b32e116927c61aa5c406183061d4f3f2d61b292c88e7f6bdf7ca AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -o /out/loadgen .

FROM alpine:3.24.2@sha256:294b683cb724975bec92580e1e685676bd4b50bda910ddb8c51d4cabeaec77e6
COPY --from=build /out/loadgen /usr/local/bin/loadgen
USER 65532:65532
ENTRYPOINT ["/usr/local/bin/loadgen"]
