run:
	go run main.go

gen:
	protoc -I ./proto \
  --go_out ./pb --go_opt paths=source_relative \
  --go-grpc_out ./pb --go-grpc_opt paths=source_relative \
    --grpc-gateway_opt logtostderr=true \
  --grpc-gateway_out ./pb --grpc-gateway_opt paths=source_relative \
   --grpc-gateway_opt generate_unbound_methods=true \
  ./proto/*.proto
clean:
	rm pb/*.go