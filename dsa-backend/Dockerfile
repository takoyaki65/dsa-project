FROM ubuntu:24.04

WORKDIR /app

# アプリケーションを起動
CMD ["dockerize", "-wait", "tcp://db:5432", "-timeout", "30s", "./main"]
