#!/usr/bin/env bash

#
#
# Copyright (c) 2019 Stackinsights to present
# All rights reserved
#

set -ex

SCRIPTDIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

[[ -z "$1" ]] && { echo  "\nPlease call '$0 <file path>' to run this command!\n"; exit 1; }
TARGET=$1
OUT=$(mktemp) || { echo "Failed to create temp file"; exit 1; }

cat "${SCRIPTDIR}/boilerplate.go.txt" > "${OUT}"
echo "" >> "${OUT}"
cat "${TARGET}" >> "${OUT}"
mv "${OUT}" "${TARGET}"
