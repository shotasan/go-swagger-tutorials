package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	"go-swagger-tutorials/custom-server/gen/restapi"
	"go-swagger-tutorials/custom-server/gen/restapi/operations"
)

var portFlag = flag.Int("port", 3000, "Port to run this service on")

func main() {
	// swagger仕様の読み込み
	// load embedded swagger file
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatal(err)
	}

	// APIサーバーの作成
	// create new service API
	api := operations.NewGreeterAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer func() {
		_ = server.Shutdown()
	}()

	// デフォルトのポート(ランダム)をportFlagを使って上書きする
	// parse flags
	flag.Parse()
	// set the port this service will be run on
	server.Port = *portFlag

	// GetGreetingHandler greets the given name,
	// in case the name is not giben, it will default to World
	api.GetGreetingHandler = operations.GetGreetingHandlerFunc(
		func(params operations.GetGreetingParams) middleware.Responder {
			// ポインタで渡されるparams.Nameをstringに変換する
			name := swag.StringValue(params.Name)
			if name == "" {
				name = "World"
			}

			greeting := fmt.Sprintf("Hello, %s", name)
			return operations.NewGetGreetingOK().WithPayload(greeting)
		},
	)

	// サーバーの開始
	// serve API
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
