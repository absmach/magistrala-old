package main

import (
	"context"
	"fmt"

	"github.com/absmach/magistrala"
	authclient "github.com/absmach/magistrala/internal/clients/grpc/auth"
)

func main() {
	auth, authHandler, err := authclient.Setup("test")
	if err != nil {
		panic(err)
	}
	defer authHandler.Close()

	res, err := auth.ListPermissions(context.Background(), &magistrala.ListPermissionsReq{
		SubjectType: "user",
		Subject:     "a1c9b979-d501-4b46-82cb-7aa6b645367c",
		Object:      "a35d58ec-0d3a-446d-80ff-796c0fe17f8a",
		ObjectType:  "thing",
	})

	fmt.Println(err)
	fmt.Println(res)
}
