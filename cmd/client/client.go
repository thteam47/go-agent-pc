package clienthttp

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"path"
	"strings"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	gw "github.com/thteam47/common/api/agent-pc"
	"google.golang.org/grpc"
)

func Run(lis net.Listener, grpc_port string, http_port string) error {
	// grpcServerEndpoint := flag.String("grpc-server-endpoint", grpc_port, "gRPC server endpoint")
	// log.Printf("dial server %s", *grpcServerEndpoint)
	transportOption := grpc.WithInsecure()
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	dialOpts := []grpc.DialOption{transportOption}
	gwmux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.HTTPBodyMarshaler{
		Marshaler: &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   true,
				EmitUnpopulated: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		},
	}))

	err := gw.RegisterAgentpcServiceHandlerFromEndpoint(ctx, gwmux, grpc_port, dialOpts)
	if err != nil {
		return err
	}
	mux := http.NewServeMux()
	mux.Handle("/", gwmux)
	mux.HandleFunc("/agent-pc/", serveSwaggerFile)
	log.Println("REST server ready...")
	// s := &http.Server{Handler: allowCORS(mux)}
	return http.ListenAndServe(http_port, allowCORS(mux))
	// return s.Serve(lis)
}
func serveSwaggerFile(w http.ResponseWriter, r *http.Request) {
	if !strings.HasSuffix(r.URL.Path, "swagger.json") {
		fmt.Printf("Not Found: %s\r\n", r.URL.Path)
		http.NotFound(w, r)
		return
	}
	p := strings.TrimPrefix(r.URL.Path, "/agent-pc")
	p = path.Join("./pkg/client/api/", p)
	fmt.Printf("Serving swagger-file: %s\r\n", p)
	http.ServeFile(w, r, p)
}
func preflightHandler(w http.ResponseWriter, r *http.Request) {
	headers := []string{"Content-Type", "Accept", "Authorization"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
	glog.Infof("preflight request for %s", r.URL.Path)
}
func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        // Xử lý preflight request (nếu có)
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        // Gọi handler ban đầu
        h.ServeHTTP(w, r)
	})
}
