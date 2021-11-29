proxy:
	grpcwebproxy --backend_addr=localhost:9999 --backend_tls_noverify --run_http_server --run_tls_server=false --use_websockets --allow_all_origins

proto:
	protoc -I=proto --go_out=. --twirp_out=. stream.proto