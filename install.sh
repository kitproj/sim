#!/usr/bin/env sh
set -eux

tag=$(curl --retry 3 -s "https://api.github.com/repos/kitproj/sim/releases/latest" | jq -r '.tag_name')
version=$(echo $tag | cut -c 2-)
url="https://github.com/kitproj/sim/releases/download/${tag}/sim_${version}_$(uname)_$(uname -m | sed 's/aarch64/arm64/').tar.gz"
curl --retry 3 -L $url | tar -zxvf - sim
chmod +x sim
sudo mv sim /usr/local/bin/sim
