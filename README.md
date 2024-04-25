# AWSLambdaClone
Kubernetes 의 이해를 위해 aws lambda 의 내부 시스템을 추정, 구성한 프로젝트   
A side project to create AWS lambda in a self maintained k8s cluster   

## Website
https://dev-doe.com/lambda/list   
Currently down due to server costs XD   
Will migrate to a localized server when time comes...   

## Architecture
![lambda-clone.drawio.svg](readme%2Flambda-clone.drawio.svg)   
https://drive.google.com/file/d/1f3qbuWipvJloTq3oyVkGIveh6414z1mT/view?usp=sharing   

## API
#### lambda-clone-api.dev-doe.com
1. [GET] /lambda/list
2. [GET] /lambda/default
3. [GET] /lambda/runtimes
4. [GET] /lambda/setup
5. [GET] /lambda
6. [DELETE] /lambda
7. [POST] /lambda
8. [PATCH] /lambda

#### lambda-clone-endpoint.dev-doe.com
1. [GET] /endpoint/{lambdaKey}

## Logic

**[POST] /lambda**
1. name, file, runtime, version 이 전송된다.
2. 해당 내용은 mongodb pod 에 저장된다.
3. 생성된 mongodb hash key 를 파일명으로 s3 안에 code 를 저장한다.
4. Service account 를 통해 lambda api 는 해당 k8s 권한이 있다.
5. code 와 runtime 에 맞는 이미지를 mongodb lambda.runtime 에서 조회한다.
6. Deployment 를 생성한다.
7. NodePort Service 를 생성한다.
8. Ingress 의 /endpoint/{lambdaKey} 를 Service 와 연결한다.
9. aws-load-balancer-controller 가 aws ALB 를 Ingress spec 에 맞춘다.

**[DELETE] /lambda**
1. 해당 mongodb data 를 disabled 처리한다.
2. 해당 Service, Deployment 를 제거한다.
3. Ingress 에서 해당 rule 을 제거한다.

**[PATCH] /lambda**
1. name, file, runtime, version 이 전송된다.
2. 해당 내용은 mongodb pod 에 수정된다.
3. s3 안에 code 가 수정된다.
4. Deployment 를 update 한다.

**[GET] /lambda, [GET] /lambda/list**
1. 저장된 lambda 를 반환한다.

**[GET] /lambda/default**
1. 생성 기본 구조를 lambda.default 에서 반환합니다.

**[GET] /lambda/runtimes**
1. lambda 생성시 선택할 수 있는 runtime 과 version 을 나열합니다.
2. 해당 version 과 runtime 에 따라 default_code 를 구분합니다.

**[GET] /lambda/setup**
1. 초기 생성시 lambda.default, lambda.runtime 의 기본 데이터를 저장합니다.

## Database
기본 입력이 필요한 데이터   
[GET] /lambda/setup 을 통해 해당 데이터 주입 가능
### lambda.default
```js
db.lambda.default.insertOne({
    "runtime": "node",
    "version": "20",
    "disabled":false,
    "reg_date": Date(),
    "update_date": Date(),
});
```
### lambda.runtime
```js
db.lambda.runtime.insertOne({
    "runtime": "golang",
    "version": "1.22",
    "image": "dohyung97022/aws-lambda-clone-golang:1.22.1",
    "default_code": "package main\n\nimport (\n\t\"encoding/json\"\n\t\"net/http\"\n\t\"net/url\"\n)\n\nfunc handler(params url.Values, w *http.ResponseWriter) {\n\tjson, _ := json.Marshal(map[string]any{\"message\": \"Hello World\", \"params\": params})\n\t(*w).Write(json)\n\t(*w).WriteHeader(200)\n}\n",
    "run_command": "aws s3api get-object --bucket aws-lambda-clone --key %s ./handler.go && go run *.go",
    "disabled": false,
    "reg_date": Date(),
    "update_date": Date(),
});
db.lambda.runtime.insertOne({
    "runtime": "node",
    "version": "20",
    "image": "dohyung97022/aws-lambda-clone-node:20",
    "default_code": "function handler(req, res) {\n    res.status(200).json({message: 'Hello, world!', params: req.query});\n}\n\nexport default handler\n",
    "run_command": "aws s3api get-object --bucket aws-lambda-clone --key %s ./handler.mjs && node app.mjs",
    "disabled":false,
    "reg_date": Date(),
    "update_date": Date(),
});
db.lambda.runtime.insertOne({
    "runtime": "python",
    "version":"3.12",
    "image": "dohyung97022/aws-lambda-clone-python:3.12",
    "default_code": "from flask import request, jsonify\n\ndef handler():\n    return jsonify(message=\"Hello world!\", params=request.args), 200\n",
    "run_command": "aws s3api get-object --bucket aws-lambda-clone --key %s ./handler.py && python main.py",
    "disabled":false,
    "reg_date": Date(),
    "update_date": Date(),
});
```





## Frontend
https://github.com/dohyung97022/dev-doe.com

### lambda/list   
![lambda-list.png](readme%2Flambda-list.png)      

[GET] /lambda/list   

1. Info -> [REDIRECT] /lambda/Info
2. Endpoint -> [GET] /endpoint/{lambdaKey}
3. Delete -> [DELETE] /lambda
4. Create -> [REDIRECT] /lambda/add

### lambda/add   
![lambda-add.png](readme%2Flambda-add.png)   

[GET] /lambda/default   
[GET] /lambda/runtime   

1. Create -> [POST] /lambda

### lambda/info   
![lambda-info.png](readme%2Flambda-info.png)   

[GET] /lambda   
[GET] /lambda/runtime   

1. Save -> [PATCH] /lambda
2. Endpoint -> [GET] /endpoint/{lambdaKey}


## 참고문서, 오류, 고통
https://plant-bottle-f9b.notion.site/c915c15cd1744805af3364fc0dee9f4f
