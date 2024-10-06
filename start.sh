#!/bin/bash

# ジャッジサーバーで用いる、サンドボックスコンテナをビルド
./dsa-judge/langs/build.sh

docker-compose up --build

