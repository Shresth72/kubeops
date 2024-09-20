MODULE_NAME := kluster
GROUP_NAME := shresth72
API_VERSION := v1alpha1
GITHUB_PKG := github.com/${GROUP_NAME}

# Paths
CODEGEN_PKG := /home/shrestha/go/pkg/mod/k8s.io/code-generator@v0.29.2/generate-groups.sh
HEADER_FILE := /home/shrestha/go_lang/kubectl/operators/digital_ocean_operator/hack/boilerplate.go.txt
CRD_MANIFEST := manifests/${GROUP_NAME}.dev_${MODULE_NAME}s.yaml

# Default target
all: generate-code generate-crd create-crd

# Generate client and API code
generate-code:
	${CODEGEN_PKG} all ${GITHUB_PKG}/${MODULE_NAME}/pkg/client \
		${GITHUB_PKG}/${MODULE_NAME}/pkg/apis ${GROUP_NAME}:${API_VERSION} \
		--go-header-file ${HEADER_FILE}

# Generate CRD YAML
generate-crd:
	controller-gen paths="${GITHUB_PKG}/${MODULE_NAME}/pkg/apis/${GROUP_NAME}.dev/${API_VERSION}" \
		crd:trivialVersion=true \
		crd:crdVersion=v1 \
		out:crd:artifacts:config=manifests

# Apply the CRD to the cluster
create-crd:
	sudo k3s kubectl create -f ${CRD_MANIFEST}

# Delete the CRD from the cluster
delete-crd:
	sudo k3s kubectl delete -f ${CRD_MANIFEST}

# List all CRDs in the cluster
list-crds:
	sudo k3s kubectl get crds

.PHONY: all generate-code generate-crd create-crd delete-crd list-crds

