# petrender
docker build -t <コンテナ名> . \
docker run -it -v ${PWD}:/work <コンテナ名> \
go build parse_log.go
zbar/ean.cをこのファイルにあるean.cに変更