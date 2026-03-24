clear

cd ..

#(
#echo "building ui"
#cd ui && npm run generate
#)
#(
#echo "server vendor"
#cd server && go mod vendor
#)

#docker pull --platform linux/amd64 golang:1.25.0
#docker pull --platform linux/arm64 golang:1.25.0
#docker pull --platform linux/amd64 alpine:latest
#docker pull --platform linux/arm64 alpine:latest

export GOVERSION=1.25.0

export TAG=0.1.1


# local build test
#(
#docker build \
#--build-arg GOVERSION=${GOVERSION} \
#-t kk-scheduler-local:${TAG} \
#-f ./image-build/Dockerfile .
#)

# docker login

(
docker buildx build --platform linux/amd64,linux/arm64  \
            --build-arg GOVERSION=${GOVERSION} \
            -t cruvie/kk-scheduler:${TAG} \
            -t cruvie/kk-scheduler:latest \
            -f ./image-build/Dockerfile . \
            --push
)


