package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/absmach/magistrala/auth"
	v1 "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/authzed/authzed-go/v1"
	"github.com/authzed/grpcutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	client, err := authzed.NewClientWithExperimentalAPIs(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpcutil.WithInsecureBearerToken("12345678"),
	)
	if err != nil {
		log.Fatalf("unable to initialize client: %s", err)
	}
	stream, err := client.PermissionsServiceClient.ReadRelationships(context.Background(), &v1.ReadRelationshipsRequest{
		RelationshipFilter: &v1.RelationshipFilter{
			ResourceType:       auth.ThingType,
			OptionalResourceId: "a35d58ec-0d3a-446d-80ff-796c0fe17f8a",
			// OptionalRelation:   auth.ViewerRelation,
			OptionalSubjectFilter: &v1.SubjectFilter{
				SubjectType:       auth.UserType,
				OptionalSubjectId: "ee857235-131e-4cb6-a55f-ccb10877cc10_1fcc26f0-0d97-477f-b63c-b361f2842ff1",
			},
		},
		OptionalLimit: 0,
	})

	if err != nil {
		panic(err)
	}

loop:
	for {
		resp, err := stream.Recv()
		switch {
		case err == io.EOF:
			fmt.Println("Stream ended with EOF")
			break loop
		case err == nil:
			b, err := json.MarshalIndent(&resp, "", "  ")
			if err != nil {
				panic(err)
			}
			fmt.Println(string(b))
		default:
			panic(err)
		}
	}

	resp, err := client.PermissionsServiceClient.CheckPermission(context.Background(), &v1.CheckPermissionRequest{
		Resource: &v1.ObjectReference{
			ObjectType: auth.ThingType,
			ObjectId:   "a35d58ec-0d3a-446d-80ff-796c0fe17f8a",
		},
		Permission: auth.AdminPermission,
		Subject: &v1.SubjectReference{
			Object: &v1.ObjectReference{
				ObjectType: auth.UserType,
				ObjectId:   "ee857235-131e-4cb6-a55f-ccb10877cc10_1fcc26f0-0d97-477f-b63c-b361f2842ff1",
			},
		},
	})

	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Permissionship)

	// resp1, err := client.PermissionsServiceClient.ExpandPermissionTree(context.Background(), &v1.ExpandPermissionTreeRequest{
	// 	Resource: &v1.ObjectReference{
	// 		ObjectType: auth.ThingType,
	// 		ObjectId:   "a35d58ec-0d3a-446d-80ff-796c0fe17f8a",
	// 	},
	// 	Permission: auth.ViewPermission,
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// bb, err := json.MarshalIndent(&resp1, "", "  ")
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(string(bb))

	resp1, err := client.ExperimentalServiceClient.BulkCheckPermission(context.Background(), &v1.BulkCheckPermissionRequest{
		Items: []*v1.BulkCheckPermissionRequestItem{
			{
				Resource: &v1.ObjectReference{
					ObjectType: auth.ThingType,
					ObjectId:   "a35d58ec-0d3a-446d-80ff-796c0fe17f8a",
				},
				Permission: auth.AdminPermission, //admin permission
				Subject: &v1.SubjectReference{
					Object: &v1.ObjectReference{
						ObjectType: auth.UserType,
						ObjectId:   "ee857235-131e-4cb6-a55f-ccb10877cc10_1fcc26f0-0d97-477f-b63c-b361f2842ff1",
					},
				},
			},
			{
				Resource: &v1.ObjectReference{
					ObjectType: auth.ThingType,
					ObjectId:   "a35d58ec-0d3a-446d-80ff-796c0fe17f8a",
				},
				Permission: auth.EditPermission, //edit permission
				Subject: &v1.SubjectReference{
					Object: &v1.ObjectReference{
						ObjectType: auth.UserType,
						ObjectId:   "ee857235-131e-4cb6-a55f-ccb10877cc10_1fcc26f0-0d97-477f-b63c-b361f2842ff1",
					},
				},
			},
			{
				Resource: &v1.ObjectReference{
					ObjectType: auth.ThingType,
					ObjectId:   "a35d58ec-0d3a-446d-80ff-796c0fe17f8a",
				},
				Permission: auth.ViewPermission, //view permission
				Subject: &v1.SubjectReference{
					Object: &v1.ObjectReference{
						ObjectType: auth.UserType,
						ObjectId:   "ee857235-131e-4cb6-a55f-ccb10877cc10_1fcc26f0-0d97-477f-b63c-b361f2842ff1",
					},
				},
			},
			{
				Resource: &v1.ObjectReference{
					ObjectType: auth.ThingType,
					ObjectId:   "a35d58ec-0d3a-446d-80ff-796c0fe17f8a",
				},
				Permission: auth.MembershipPermission, //membership permission
				Subject: &v1.SubjectReference{
					Object: &v1.ObjectReference{
						ObjectType: auth.UserType,
						ObjectId:   "ee857235-131e-4cb6-a55f-ccb10877cc10_1fcc26f0-0d97-477f-b63c-b361f2842ff1",
					},
				},
			},
		},
	})

	if err != nil {
		panic(err)
	}

	for _, pair := range resp1.Pairs {
		item := pair.GetItem()
		req := pair.GetRequest()
		err := pair.GetError()
		fmt.Println(req.Permission)
		if item != nil {
			fmt.Println(item.Permissionship)
		}
		if err != nil {
			fmt.Println(err)
		}
	}
	// bb, err := json.MarshalIndent(&resp1, "", "  ")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(string(bb))
}
