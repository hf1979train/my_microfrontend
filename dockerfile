FROM golang:1.23.1-alpine
COPY src/ /home/myhomesv/
WORKDIR /home/myhomesv/
# 必要なパッケージインストール（例：Goモジュール）
RUN pwd
RUN ls -la
RUN go mod tidy

# ビルド
RUN go build -o myapp ./cmd/app/main.go
# ポート8080を公開(明示用)
 EXPOSE 8080
# ビルドしたバイナリを実行
CMD ["./myapp"]

# CMD go run /cmd/app/main.go
