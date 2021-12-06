cd ../pb-client
git pull
protoc --go-grpc_out=. --go_out=. *.proto
cp *.go ../go-gateway/pb
