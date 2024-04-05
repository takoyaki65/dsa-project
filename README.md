# はじめかた
1. homebrewのインストール  
   [homebrewのページ](https://brew.sh/ja/)を開いてインストールコマンドをターミナルにコピペ&実行．  
   パス等の設定ができていない場合はそれも行う．

2. gitのインストール
   ```sh
   # インストール
   brew install git

   # パスを通す
   # zshrcの部分は設定が書いてあるファイル
   echo 'export PATH="/opt/homebrew/bin:$PATH"' >> ~/.zshrc

   # ターミナルを再起動するか以下で設定リロード
   source ~/.zshrc
   ```

3. 最新のdsa_projectとサブモジュールのクローン
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

4. サブモジュールの環境設定  
   [dsa_backのREADME](https://github.com/zakkii-k/dsa_back/blob/main/README.md )と，[dsa_frontのREADME](https://github.com/zakkii-k/dsa_front/blob/main/README.md )を見ながら必要なパッケージ等をインストール

5. dockerの準備
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
      
