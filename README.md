# myhomesv
* golangのバージョンを切り替える
https://qiita.com/snyt45/items/2425a849db8947001587

## 開発用
```script
$ docker compose -f docker-compose.dev.yml up --build -d

# 解析する時
$ docker compose -f docker-compose.dev.yml build --progress=plain --no-cache


```
## 本番用