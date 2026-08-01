package assets

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func CreateHTTPServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/validate", func(w http.ResponseWriter, r *http.Request) {
		var req CryptoRequestJSON
		_ = json.NewDecoder(r.Body).Decode(&req)

		resp := CryptoResponseJSON{Valid: true, Message: "JSON signature processed successfuly"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	fmt.Println("[HTTP/1.1] Listening http://127.0.0.1:8081")
	_ = http.ListenAndServe(":8081", mux)
}

func CreategRPCServer() {
	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer(grpc.ForceServerCodec(RawCodec{}))
	reflection.Register(grpcServer)

	grpcServer.RegisterService(&grpc.ServiceDesc{
		ServiceName: "crypto.CryptoService",
		HandlerType: (*any)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: "ValidateBinarySgnature",
				Handler: func(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
					var rawInput []byte
					// Recibimos los bytes crudos directamente de la red HTTP/2
					if err := dec(&rawInput); err != nil {
						return nil, err
					}
					// Respuesta binaria corta de Protobuf
					return []byte{0x01, 0x4F, 0x4B}, nil // 1=true, 'O', 'K'
				},
			},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "crypto.proto",
	}, nil)

	fmt.Println("[gRPC HTTP/2] Listening 127.0.0.1:8082")
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
