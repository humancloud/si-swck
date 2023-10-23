#!/usr/bin/env bash

#
#
# Copyright (c) 2019 Stackinsights to present
# All rights reserved
#

set -u
set -ex

MOD=${MOD:-operator}
DIR=${DIR:-default}
IMG_PATH=${IMG_PATH:-manager}
IMG=${IMG:-controller}
NEW_IMG=${NEW_IMG:-controller}

OUT_DIR=$(mktemp -d -t ${MOD}-deploy.XXXXXXXXXX) || { echo "Failed to create temp file"; exit 1; }

SCRIPTPATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
TOOLBIN=${SCRIPTPATH}/../bin
ROOTDIR="${SCRIPTPATH}/../${MOD}"

main() {
    [[ $1 -eq 0 ]] && frag="apply" || frag="delete --ignore-not-found=true"
    cp -Rvf "${ROOTDIR}"/config/* "${OUT_DIR}"/.
    cd "${OUT_DIR}"/${IMG_PATH} && ${TOOLBIN}/kustomize edit set image ${IMG}=${NEW_IMG}
    ${TOOLBIN}/kustomize build "${OUT_DIR}"/${DIR} | kubectl ${frag} -f -
}

usage() {
cat <<EOF
Usage:
    ${0} -[duh]

Parameters:
    -d  Deploy ${MOD}
    -u  Undeploy ${MOD}
    -h  Show this help.
EOF
exit 1
}

parseCmdLine(){
    ARG=$1
    if [ $# -eq 0 ]; then
        echo "Exactly one argument required."
        usage
    fi
    case "${ARG}" in
        d) main 0;;
        u) main 1;;
        h) usage ;;
        \?) usage ;;
    esac
	  return 0
}

#
# main
#

ret=0

parseCmdLine "$@"
ret=$?
[ $ret -ne 0 ] && exit $ret
echo "Done deploy [$NEW_IMG] (exit $ret)"
