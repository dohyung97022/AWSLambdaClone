docker build -t dohyung97022/aws-lambda-clone:latest .
docker push dohyung97022/aws-lambda-clone:latest
ssh master rm -r "~/aws-lambda-clone"
ssh master mkdir "~/aws-lambda-clone"
scp -r "./kubernetes" "master:~/aws-lambda-clone"
ssh master "cd ~/aws-lambda-clone/kubernetes; rm lambda-clone-ingress.yaml"
ssh master "cd ~/aws-lambda-clone/kubernetes; kubectl apply -f ."
