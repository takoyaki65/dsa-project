# 準備
1. パッケージのインストール  
   注意: このセクションは，ローカル環境でインテリセンスが利くようにするための手順について
   記述している．実際にアプリを動かす際には，venvを作成しなくても良い．

   Pythonパッケージマネージャーの[uv](https://docs.astral.sh/uv/)をインストールする．
   ```bash
   # macOSまたはLinuxの場合
   curl -LsSf https://astral.sh/uv/install.sh | sh
   # 実行後、環境変数の設定が即座に反映されない場合は、シェルを開きなおしてください
   # バージョン確認
   uv --version
   ```

   venvを作成する．
   ```bash
   .../dsa-back$ uv sync
   ```

   venvをactivateする
   ```bash
   .../dsa-back$ . .venv/bin/activate
   ```
