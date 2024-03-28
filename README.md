# はじめかた
1. 最新のdsa_projectとサブモジュールのクローン
   ```bash
   git clone --recurse-submodules https://github.com/zakkii-k/dsa_project
   ```
   もし既にクローンしている場合は
   ```bash
   git submodule update --init --recursive
   ```
   dsa_project自体のアップデートは
   ```bash
   git pull origin master #もしくは利用中のブランチ名
   ```

2. dockerの準備
   1. インストール  
        dockerの公式サイトからdocker desktop for macをインストールし，初期設定を行う．
      ```bash
      docker --version
      ```
        などのコマンドでインストールが完了していることを確認する．

   2. 起動  
        docker desktopを起動する．  
        その後，dsa_projectをクローンしたディレクトリへ移動し，
      ```bash
      docker-compose up --build
      ```
        を入力．

# コンテナへのアクセス
1. コンテナ名の確認
  ```bash
  docker ps
  ```
  NAMESの欄からアクセスしたいコンテナの名前を確認する．
2. アクセス
  ```bash
  docker exec -it [NAMES] bash
  ```
   
