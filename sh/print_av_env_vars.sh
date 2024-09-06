#!/bin/bash
scriptdir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
basedir="${scriptdir%/*}"

src="${basedir}/src/conf.go"

cat "${src}" | grep -Po '(?<=case ").*(?=")'
