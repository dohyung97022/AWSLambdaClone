FROM --platform=linux/amd64 golang:1.21.5
COPY "./src" "~/src"
WORKDIR "~/src"
RUN ["go", "build", "-o", "app", "."]
CMD ["./app"]
