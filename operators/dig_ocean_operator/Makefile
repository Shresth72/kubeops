# Default kubectl command
KUBECTL := sudo k3s kubectl

# Variables
MODULE_NAME := kluster
GROUP_NAME := shresth72
API_VERSION := v1alpha1
SECRET_NAME := dosecret
GITHUB_PKG := github.com/${GROUP_NAME}
DO_TOKEN := ""

# Paths
CODEGEN_PKG := /home/shrestha/go/pkg/mod/k8s.io/code-generator@v0.29.2/kube_codegen.sh

HEADER_FILE := ~/go/pkg/mod/k8s.io/code-generator@v0.29.2/hack/boilerplate.go.txt

CRD_MANIFEST := manifests/${GROUP_NAME}.dev_${MODULE_NAME}s.yaml

MANIFEST_DEPLOY_PATH := manifests/deploy





# Generator
generate-all: generate-code generate-crd

generate-code:
	${CODEGEN_PKG} all ${GITHUB_PKG}/${MODULE_NAME}/pkg/client \
		${GITHUB_PKG}/${MODULE_NAME}/pkg/apis ${GROUP_NAME}.dev:${API_VERSION} \
		--go-header-file ${HEADER_FILE}

generate-crd:
	controller-gen paths="${GITHUB_PKG}/${MODULE_NAME}/pkg/apis/${GROUP_NAME}.dev/${API_VERSION}" \
		crd:trivialVersion=true \
		crd:crdVersion=v1 \
		out:crd:artifacts:config=manifests




# CRD 
create-crd:
	${KUBECTL} create -f ${CRD_MANIFEST}

delete-crd:
	${KUBECTL} delete -f ${CRD_MANIFEST}

list-crds:
	${KUBECTL} get crds




# Secrets
create-secret:
	${KUBECTL} create secret ${SECRET_NAME} --from-literal token=${DO_TOKEN}

list-secret:
	${KUBECTL} get secrets ${SECRET_NAME} -oyaml




# Kluster
create-kluster:
	${KUBECTL} create -f manifests/klusterone.yaml

delete-kluster:
	${KUBECTL} delete -f manifests/klusterone.yaml

list-klusters:
	${KUBECTL} get kluster




# Docker
docker-build:
	docker build -t shresth72/kluster:0.1.0 .

docker-push:
	docker push shresth72/kluster:0.1.0

docker-run:
	docker run --rm -it shresth72/kluster:0.1.0 




# Deployments
create-all-deploy:
	${KUBECTL} create -f manifests/deploy

create-deploy-yaml:
	${KUBECTL} create deployment kluster --image shresth72/kluster:0.1.0 --dry-run=client -oyaml > ${MANIFEST_DEPLOY_PATH}/deploy.yaml

create-deploy:
	${KUBECTL} create -f ${MANIFEST_DEPLOY_PATH}/deploy.yaml




# RBAC
create-service-account:
	${KUBECTL} create serviceaccount kluster-sa --dry-run=client -oyaml > ${MANIFEST_DEPLOY_PATH}/sa.yaml

create-clusterrole:
	${KUBECTL} create clusterrole kluster-cr --resource Kluster --verb list,watch --dry-run=client -oyaml > ${MANIFEST_DEPLOY_PATH}/cr.yaml

create-clusterrole-binding:
	${KUBECTL} create clusterrolebinding kluster-crb --clusterrole kluster-cr --serviceaccount default:kluster-sa --dry-run=client -oyaml > ${MANIFEST_DEPLOY_PATH}/crb.yaml

create-role:
	${KUBECTL} create role k-role --resource secret --verb get --dry-run=client -oyaml > ${MANIFEST_DEPLOY_PATH}/krole.yaml

create-role-binding:
	${KUBECTL} create rolebinding k-rb --serviceaccount default:kluster-sa --role k-role --dry-run=client -oyaml > ${MANIFEST_DEPLOY_PATH}/kroleb.yaml




clean: delete-crd delete-kluster

.PHONY: all generate-code generate-crd create-crd delete-crd list-crds create-kluster delete-kluster

