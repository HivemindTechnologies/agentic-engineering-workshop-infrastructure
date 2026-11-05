# Built by `just images` from the git tree of scoring/, never from the working
# directory, so the image tag (that tree's hash) names exactly this content.
# Base images are pinned by digest: a tag can move, a digest cannot.
FROM golang:1.26.10-alpine3.24@sha256:9ce231b28998b32e116927c61aa5c406183061d4f3f2d61b292c88e7f6bdf7ca AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/score ./cmd/score

FROM alpine:3.24.2@sha256:294b683cb724975bec92580e1e685676bd4b50bda910ddb8c51d4cabeaec77e6
COPY --from=build /out/score /usr/local/bin/score
USER 65532:65532
ENTRYPOINT ["/usr/local/bin/score"]
