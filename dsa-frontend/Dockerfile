# ベースイメージ
FROM node:23.9.0-bookworm-slim

# アプリケーションディレクトリを作成
WORKDIR /app

# アプリケーションの依存関係をインストール
# COPY package*.json ./
# RUN npm install

# アプリケーションのソースをバンドル (.venv, .vnodeenv, node_modulesフォルダを除く)
# COPY . .
# RUN rm -rf .venv .vnodeenv node_modules

EXPOSE 3000

CMD ["npm", "start"]