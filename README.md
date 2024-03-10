# AWSLambdaClone
Kubernetes 의 이해를 위해 aws lambda 의 내부 시스템을 추정, 구성한 프로젝트   
A side project to create AWS lambda in a self maintained k8s cluster

## Architecture

## API

1. [GET] /lambda/list
2. [GET] /lambda/default
3. [GET] /lambda/runtimes
4. [GET] /lambda
5. [DELETE] /lambda
6. [POST] /lambda
7. [PATCH] /lambda
8. [GET] /{lambdaHashKey}

## Logic

**[POST] /lambda**
1. /lambda/create 페이지에서 [POST] /lambda 가 호출될 때, name, file, runtime, version 이 전송된다.
2. lambda api pod 에서 해당 내용을 받은 이후, name, runtime, version 은 mongodb pod 에 저장된다. hash 가 생성되어 mongodb key 로 지정된다. file 의 경우 s3 내에 파일 형태로 저장되고 해당 주소도 mongodb 에 저장된다.
3. lambda api pod 는 그 이후 controlplane api 에 service account 의 권한을 통해 해당 runtime, version 에 해당되는 image docker hub 에서 받아 s3의 코드 파일을 주입, 실행시키는 command 의 deployment 를 만든다.
4. 해당 pod label 에 clusterIp 의 service 를 생성하고, 해당 service 를 mongo 저장시 생성된 hash key 의 path 로 ingress 에 지정한다.

**[DELETE] /lambda**
1. lambda api pod 에 key 값인 hash 가 전송되면 controlplane api 에서 service account 의 권한을 통해 deployment 를 삭제하고, ingress 에 해당 hash 의 path 를 제거한다.
2. 해당 값이 이미 존재하지 않거나, 삭제에 성공한 경우 mongodb 에서 해당 hash key 의 데이터 상태를 삭제 상태로 바꾼다.

**[PATCH] /lambda**
1. hash, name, file, runtime, version 이 전송된다. 해당 hash key 의 데이터를 mongodb에서 찾아 name, runtime, version 이 수정되고 file 의 경우 수정 내용이 있을 경우 s3에 신규 코드 파일로 대체된다.
2. lambda api pod 는 이후 controlplane api 에 deployment 를 해당 runtime, version 의 이미지로 변경하고 command 또한 신규 s3 로 변경한다.
3. ingerss 의 경우 수정사항이 없이 그대로 유지한다.

**[GET] /lambda, [GET] /lambda/list**
1. mongodb 에서 저장된 lambda 내용을 반환한다.

**[GET] /lambda/default**
1. lambda 의 생성시 기본 구조를 반환합니다.

**[GET] /lambda/runtimes**
1. lambda 생성시 선택할 수 있는 runtime 과 version 을 나열합니다.


## Database
기본 입력이 필요한 데이터   
[GET] /lambda/setup 을 통해 해당 데이터 주입 가능
### lambda.default
```js
db.lambda.default.insertOne({
    "runtime": "golang",
    "version":'1.22',
    "disabled":false,
    "reg_date": Date(),
    "update_date": Date(),
});
```
### lambda.runtime
```js
db.lambda.runtime.insertOne({
    "runtime": "golang",
    "version":'1.22',
    "disabled":false,
    "reg_date": Date(),
    "update_date": Date(),
});
db.lambda.runtime.insertOne({
    "runtime": "node",
    "version":'20',
    "disabled":false,
    "reg_date": Date(),
    "update_date": Date(),
});
db.lambda.runtime.insertOne({
    "runtime": "python",
    "version":'3.12',
    "disabled":false,
    "reg_date": Date(),
    "update_date": Date(),
});
```





## Frontend

**lambda/list**
1. delete
2. create
3. click → info

**lambda/create**
1. Lambda Name
2. Code
3. Runtime
4. Create Button

**lambda/info**
1. EndpointUrl
2. Test, Parameters
3. Result
4. edit
