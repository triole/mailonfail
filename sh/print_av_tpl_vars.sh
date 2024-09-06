#!/bin/bash
scriptdir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
basedir="${scriptdir%/*}"

src="${basedir}/src/mail.go"

cat "${src}" |
  grep -Po '"[a-z_]+": .*[a-zA-Z\.()],$' |
  grep -Po '(?<=")[a-z_]+'
