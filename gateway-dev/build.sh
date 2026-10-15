IMAGE_VERSION=$(git rev-parse --short HEAD)
TIME=$(date +%Y%m%d%H%M%S)
TAG=${IMAGE_VERSION}_${TIME}
APP_NAME="gateway"
IMAGE_REMOTE="registry-intl.ap-southeast-1.aliyuncs.com/relax_many/"
IMAGE_PUSH=${IMAGE_REMOTE}${APP_NAME}:${TAG}

echo "start docker build ${IMAGE_PUSH}"

# 构建
docker build -t "${IMAGE_PUSH}" .
# 利用 buildx 构建 amd64 环境 的镜像
#docker buildx build --platform linux/amd64 -t "${IMAGE_PUSH}" .

# push image
docker push "${IMAGE_PUSH}"
echo "push success:  ${IMAGE_PUSH}"
