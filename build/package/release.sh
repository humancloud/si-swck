#!/usr/bin/env bash

#
# Copyright (c) 2019 Stackinsights to present
# All rights reserved
#
#

set -ex
SCRIPTDIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
BUILDDIR=${SCRIPTDIR}/..
ROOTDIR=${BUILDDIR}/..

RELEASE_TAG=$(git describe --tags $(git rev-list --tags --max-count=1))
RELEASE_VERSION=${RELEASE_TAG#"v"}

binary(){
    bindir=${BUILDDIR}/release/binary
    rm -rf ${bindir}
    mkdir -p ${bindir}/config
    # Copy relevant files
    cp -Rfv ${BUILDDIR}/bin ${bindir}
    cp -Rfv ${ROOTDIR}/CHANGES.md ${bindir}
    cp -Rfv ${ROOTDIR}/docs/binary-readme.md ${bindir}/README.md
    cp -Rfv ${ROOTDIR}/dist/* ${bindir}
    cp -Rfv ${ROOTDIR}/operator/dist/* ${bindir}
    mkdir -p ${bindir}/licenses/adapter-licenses
    cp -Rfv ${ROOTDIR}/adapter/dist/licenses/* ${bindir}/licenses/adapter-licenses
	cat ${ROOTDIR}/adapter/dist/LICENSE >> ${bindir}/LICENSE
    # Docker
    cp -Rfv ${ROOTDIR}/build/images/Dockerfile.release-bin ${bindir}/Dockerfile
    echo -e "build:" > ${bindir}/Makefile
    echo -e "\tdocker build . -t apache/stackinsights-swck:${RELEASE_TAG}" >> ${bindir}/Makefile
    # Generates CRDs and deployment manifests
    pushd ${ROOTDIR}/operator/config/manager
    kustomize edit set image controller=apache/stackinsights-swck:${RELEASE_TAG}
    popd
    kustomize build operator/config/default > ${bindir}/config/operator-bundle.yaml
    pushd ${ROOTDIR}/adapter/config/namespaced/adapter
    kustomize edit set image metrics-adapter=apache/stackinsights-swck:${RELEASE_TAG}
    popd
    kustomize build adapter/config > ${bindir}/config/adapter-bundle.yaml
    # Package
    tar -czf ${BUILDDIR}/release/stackinsights-swck-${RELEASE_VERSION}-bin.tgz -C ${bindir} .
    rm -rf ${bindir}
}

source(){
    # Package
    rm -rf ${BUILDDIR}/release/stackinsights-swck-${RELEASE_VERSION}-src.tgz
    pushd ${ROOTDIR}
    tar \
        --exclude=".DS_Store" \
        --exclude=".git" \
        --exclude=".github" \
        --exclude=".gitignore" \
        --exclude=".asf.yaml" \
        --exclude=".idea" \
        --exclude="bin"  \
        --exclude="operator/bin"  \
        --exclude="adapter/bin"  \
        --exclude="build/release"  \
        --exclude="*.test"  \
        --exclude="*.out"  \
        -czf ./build/release/stackinsights-swck-${RELEASE_VERSION}-src.tgz \
        .
    popd
}

sign(){
    type=$1
    pushd ${BUILDDIR}/release/
    gpg --batch --yes --armor --detach-sig stackinsights-swck-${RELEASE_VERSION}-${type}.tgz
	shasum -a 512 stackinsights-swck-${RELEASE_VERSION}-${type}.tgz > stackinsights-swck-${RELEASE_VERSION}-${type}.tgz.sha512
	popd
}

parseCmdLine(){
    ARGS=$1
    if [ $# -eq 0 ]; then
        echo "Exactly one argument required."
        usage
    fi
    while getopts  "bsk:h" FLAG; do
        case "${FLAG}" in
            b) binary ;;
            s) source ;;
            k) sign ${OPTARG} ;;
            h) usage ;;
            \?) usage ;;
        esac
    done
	return 0
}



usage() {
cat <<EOF
Usage:
    ${0} -[bsh]

Parameters:
    -b  Build and assemble the binary package
    -s  Assemble the source package
    -h  Show this help.
EOF
exit 1
}

#
# main
#

ret=0

parseCmdLine "$@"
ret=$?
[ $ret -ne 0 ] && exit $ret
echo "Done release [RELEASE_VERSION] (exit $ret)"
