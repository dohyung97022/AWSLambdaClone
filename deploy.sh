docker build -t dohyung97022/aws-lambda-clone:latest .
docker push dohyung97022/aws-lambda-clone:latest
ssh master mkdir "~/aws-lambda-clone"
scp -r "./kubernetes" "master:~/aws-lambda-clone"

## deploy kubernetes
#ssh master
#cd ~/aws-lambda-clone/kubernetes
#kubectl apply -f .
