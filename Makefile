VERSION=latest

# 测试版本
VERSION_TEST=$(VERSION)
# 编译的程序名称
APP_NAME_TEST=api-test
# 测试版本
K8S_NAME_TEST=api-test
# 测试下的编译文件
DOCKER_FILE_TEST=./deploy/dockerfile_dev
# 测试环境配置
DOCKER_REPO_TEST=
# 触发容器重新拉取镜像重新部署的URL
K8S_UPDATE_URL_TEST=

# 测试环境的编译发布
build-test:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build . -f ${DOCKER_FILE_TEST} --no-cache -t ${APP_NAME_TEST}

build-test-win:
	go env -w GOOS=linux GOARCH=amd64 CGO_ENABLED=0
	go build . -f ${DOCKER_FILE_TEST} --no-cache -t ${APP_NAME_TEST}
	go env -w GOOS=windows

tag-test:
	@echo 'create tag ${VERSION_TEST}'
	docker tag ${APP_NAME_TEST} ${DOCKER_REPO_TEST}:${VERSION_TEST}

publish-test:
	@echo 'publish ${VERSION_TEST} to ${DOCKER_REPO_TEST}'
	docker push ${DOCKER_REPO_TEST}:${VERSION_TEST}

redeploy-k8s-test:
	@echo 'update k8s container image'
	curl ${K8S_UPDATE_URL_TEST}

release-test: build-test tag-test publish-test redeploy-k8s-test

release-test-win: build-test-win tag-test publish-test redeploy-k8s-test

# 测试版本
VERSION_PROD=$(VERSION)
# 编译的程序名称
APP_NAME_PROD=api-prod
# 测试版本
K8S_NAME_PROD=api-prod
# 测试下的编译文件
DOCKER_FILE_PROD=./deploy/dockerfile_dev
# 测试环境配置
DOCKER_REPO_PROD=
# 触发容器重新拉取镜像重新部署的URL
K8S_UPDATE_URL_PROD=

# 测试环境的编译发布
build-prod:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build . -f ${DOCKER_FILE_PROD} --no-cache -t ${APP_NAME_PROD}

build-prod-win:
	go env -w GOOS=linux GOARCH=amd64 CGO_ENABLED=0
	go build . -f ${DOCKER_FILE_PROD} --no-cache -t ${APP_NAME_PROD}
	go env -w GOOS=windows

tag-prod:
	@echo 'create tag ${VERSION_PROD}'
	docker tag ${APP_NAME_PROD} ${DOCKER_REPO_PROD}:${VERSION_PROD}

publish-prod:
	@echo 'publish ${VERSION_PROD} to ${DOCKER_REPO_PROD}'
	docker push ${DOCKER_REPO_PROD}:${VERSION_PROD}

redeploy-k8s-prod:
	@echo 'update k8s container image'
	curl ${K8S_UPDATE_URL_PROD}

release-prod: build-prod tag-prod publish-prod redeploy-k8s-prod

release-prod-win: build-prod-win tag-prod publish-prod redeploy-k8s-prod