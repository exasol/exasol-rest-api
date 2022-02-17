#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

oft_version="3.4.0"
tmp_dir="/tmp/oft/"
tmp_file="$tmp_dir/openfasttrace-$oft_version.jar"

if [[ ! -f "$tmp_file" ]]; then
    mkdir -p "$tmp_dir"
    url="https://repo1.maven.org/maven2/org/itsallcode/openfasttrace/openfasttrace/$oft_version/openfasttrace-$oft_version.jar"
    echo "Downloading $url to $tmp_file"
    curl --output "$tmp_file" "$url"
fi

java -jar "$tmp_file" trace
