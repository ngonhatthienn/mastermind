main_run:
	go run main.go

auth_run:
	go run cmd/auth/auth.go
  
gen:
	protoc -I ./proto \
  --go_out ./pb/game --go_opt=paths=source_relative \
  --go-grpc_out ./pb/game --go-grpc_opt paths=source_relative \
  --grpc-gateway_out ./pb/game --grpc-gateway_opt paths=source_relative \
  ./proto/gameservice.proto
	protoc -I ./proto \
  --go_out ./pb/auth --go_opt=paths=source_relative \
  --go-grpc_out ./pb/auth --go-grpc_opt paths=source_relative \
  --grpc-gateway_out ./pb/auth --grpc-gateway_opt paths=source_relative \
  ./proto/auth.proto

clean:
	rm pb/game/*.go
	rm pb/auth/*.go