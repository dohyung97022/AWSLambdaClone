FROM --platform=linux/amd64 python:3.12-alpine3.19
COPY "./runtimeImages/python3.12/src" "~/src"
WORKDIR "~/src"
RUN ["sh", "-c", "pip install -r requirements.txt && apk add --no-cache aws-cli && rm ./handler.py"]
