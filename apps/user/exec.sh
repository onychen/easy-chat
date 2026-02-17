easy-chat&gt; go mod init
easy-chat&gt; cd /user/rpc
easy-chat/user/rpc&gt; goctl rpc protoc ./apps/user/rpc/user.proto --go_out=./apps/user/rpc/ --go-grpc_out=./apps/user/rpc/ --zrpc_out=./apps/user/rpc/
easy-chat/user/rpc&gt; go mod tidy

goctl model mysql ddl -src="./deploy/sql/user.sql" -dir="./apps/user/models/" -c
goctl api go -api apps/user/api/user.api -dir apps/user/api -style gozero