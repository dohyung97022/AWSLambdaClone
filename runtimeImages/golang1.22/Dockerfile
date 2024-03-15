FROM --platform=linux/amd64 golang:1.22.1-alpine3.19
COPY "./runtimeImages/golang1.22/src" "~/src"
WORKDIR "~/src"
ENV GOOS linux
ENV GOARCH amd64
RUN ["sh", "-c", "apk add --no-cache aws-cli && rm ./handler.go"]
