#!/bin/bash

# Tip: Output from this is easy to pipe to xargs rm

die () {
    echo >&2 "$@"
    exit 1
}

[ "$#" -eq 1 ] || die "path argument required, $# args provided"


find $1 -name "metadata.toml.hash"
find $1 -name "content.adoc.hash"
find $1 -name "content.html"
find $1 -name "full.pdf"
find $1 -name "full.xml"