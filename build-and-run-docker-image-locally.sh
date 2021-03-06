
IMAGE_URL=oyvindsk.com:latest

# remove :latest so we don't acidentally run an old one
# don't use :latest when deploying, keep it ..
# sudo docker rmi ${IMAGE_URL}

sudo docker build -t ${IMAGE_URL} -f ./SECRET-Dockerfile .

sudo docker run -ti --rm -p 8080:8080 ${IMAGE_URL}
