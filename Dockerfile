# 構建階段
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安裝構建工具
RUN apk add --no-cache git gcc musl-dev

# 複製依賴文件
COPY go.mod go.sum ./
RUN go mod download

# 複製源代碼
COPY . .

# 構建應用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 運行階段
FROM alpine:latest

WORKDIR /app

# 安裝運行時依賴
RUN apk --no-cache add ca-certificates tzdata

# 設置時區
ENV TZ=Asia/Taipei

# 從構建階段複製二進制文件和配置文件
COPY --from=builder /app/main .
COPY --from=builder /app/web ./web
COPY --from=builder /app/internal/config ./internal/config

# 暴露端口
EXPOSE 8080

# 運行應用
CMD ["./main"] 