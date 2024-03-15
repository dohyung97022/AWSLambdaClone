FROM --platform=linux/amd64 node:20-alpine3.19
COPY "./runtimeImages/node20/src" "~/src"
WORKDIR "~/src"
RUN ["sh", "-c", "npm install && apk add --no-cache aws-cli && rm ./handler.mjs"]
