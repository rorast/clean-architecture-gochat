FROM golang:1.21-alpine

WORKDIR /app

# 複製必要的檔案
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 編譯遷移工具
RUN go build -o migrate ./sql/migrate.go

# 設定執行權限
RUN chmod +x /app/migrate

# 設定預設命令
CMD ["/app/migrate"] 