#!/usr/bin/env sh
set -eux

tag=
while [ "$tag" = "" ]; do
  tag=$(curl -fsL "https://api.github.com/repos/kitproj/sim/releases/latest" | jq -r '.tag_name')
done

version=$(echo $tag | cut -c 2-)
url="https://github.com/kitproj/sim/releases/download/${tag}/kit_${version}_$(uname)_$(uname -m | sed 's/aarch64/arm64/').tar.gz"

while [ ! -e kit ]; do
  curl -fsL $url | tar -zxvf - sim
done

chmod +x sim
sudo mv kit /usr/local/bin/sim
