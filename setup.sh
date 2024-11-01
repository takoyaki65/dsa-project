#!/bin/bash

set -e

# sandboxコンテナのビルド
./dsa-judge/langs/build.sh

if [[ "$(uname)" == "Linux" ]]; then
    ./dsa-judge/prepare-sandbox.sh
fi

# resource/ディレクトリ内の全てのファイルについてsandbox内でパーミッションエラーが起きないために、
#   1. ファイルのuid:gidを1000:1000に変更
#   2. .shで終わるファイルには、ファイルのパーミッションに+xを付与
for file in $(find resource); do
    sudo chown 1000:1000 $file
    if [[ $file == *.sh ]]; then
        sudo chmod +rx $file
    fi
done
