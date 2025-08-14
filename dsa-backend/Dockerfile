FROM golang:1.24.6-bookworm

WORKDIR /app

# アプリケーションを起動
CMD ["dockerize", "-wait", "tcp://db:5432", "-timeout", "30s", "go", "run", "main.go"]
