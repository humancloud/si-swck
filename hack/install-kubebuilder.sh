#!/usr/bin/env bash

#
#
# Copyright (c) 2019 Stackinsights to present
# All rights reserved
#
set -ex

os=$(go env GOOS)
arch=$(go env GOARCH)
export K8S_VERSION=1.19.2
export PATH=$PATH:/usr/local/kubebuilder/bin
sudo mkdir -p /usr/local/kubebuilder/bin
curl -sSLo kubebuilder https://github.com/kubernetes-sigs/kubebuilder/releases/download/v3.2.0/kubebuilder_${os}_${arch}
sudo mv ./kubebuilder /usr/local/bin/

curl -sSLo envtest-bins.tar.gz "https://go.kubebuilder.io/test-tools/${K8S_VERSION}/${os}/${arch}"
sudo  tar -C /usr/local/kubebuilder --strip-components=1 -zvxf envtest-bins.tar.gz
