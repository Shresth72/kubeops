MODULE_NAME := kluster
GROUP_NAME := shresth72
API_VERSION := v1alpha1
SECRET_NAME := dosecret
GITHUB_PKG := github.com/${GROUP_NAME}
DO_TOKEN := ""

# Paths
CODEGEN_PKG := /home/shrestha/go/pkg/mod/k8s.io/code-generator@v0.29.2/generate-groups.sh
HEADER_FILE := /home/shrestha/go_lang/kubectl/operators/digital_ocean_operator/hack/boilerplate.go.txt
CRD_MANIFEST := manifests/${GROUP_NAME}.dev_${MODULE_NAME}s.yaml

# Generator
generate-all: generate-code generate-crd

generate-code:
	${CODEGEN_PKG} all ${GITHUB_PKG}/${MODULE_NAME}/pkg/client \
		${GITHUB_PKG}/${MODULE_NAME}/pkg/apis ${GROUP_NAME}:${API_VERSION} \
		--go-header-file ${HEADER_FILE}

generate-crd:
	controller-gen paths="${GITHUB_PKG}/${MODULE_NAME}/pkg/apis/${GROUP_NAME}.dev/${API_VERSION}" \
		crd:trivialVersion=true \
		crd:crdVersion=v1 \
		out:crd:artifacts:config=manifests




# CRD 
create-crd:
	sudo k3s kubectl create -f ${CRD_MANIFEST}

delete-crd:
	sudo k3s kubectl delete -f ${CRD_MANIFEST}

list-crds:
	sudo k3s kubectl get crds




# Secrets
create-secret:
	sudo k3s kubectl create secret ${SECRET_NAME} --from-literal token=${DO_TOKEN}

list-secret:
	sudo k3s kubectl get secrets ${SECRET_NAME} -oyaml




# Kluster
create-kluster:
	sudo k3s kubectl create -f manifests/klusterone.yaml

delete-kluster:
	sudo k3s kubectl delete -f manifests/klusterone.yaml

list-klusters:
	sudo k3s kubectl get kluster




clean: delete-crd delete-kluster

.PHONY: all generate-code generate-crd create-crd delete-crd list-crds create-kluster delete-kluster

