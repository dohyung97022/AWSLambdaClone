FROM --platform=linux/amd64 golang:1.21.5
COPY "./src" "~/src"
WORKDIR "~/src"
RUN ["env", "GOOS=linux", "GOARCH=amd64", "go", "build", "-o", "app", "."]
CMD ["./app"]
