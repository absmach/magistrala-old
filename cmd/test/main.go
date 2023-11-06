package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/absmach/magistrala/auth"
	v1 "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/authzed/authzed-go/v1"
	"github.com/authzed/grpcutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	fileBytes, err := os.ReadFile("docker/spicedb/schema.zed")
	if err != nil {
		panic(err)
	}

	spicedbSchema := string(fileBytes)

	// grpcOpts, err := grpcutil.WithSystemCerts(grpcutil.SkipVerifyCA)
	// if err != nil {
	// 	panic(err)
	// }
	client, err := authzed.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpcutil.WithInsecureBearerToken("12345678"),
	)
	if err != nil {
		log.Fatalf("unable to initialize client: %s", err)
	}

	if _, err := client.SchemaServiceClient.WriteSchema(context.Background(), &v1.WriteSchemaRequest{Schema: spicedbSchema}); err != nil {
		panic(err)
	}

	client.PermissionsServiceClient.DeleteRelationships(context.Background(), &v1.DeleteRelationshipsRequest{})

	client.PermissionsServiceClient.ReadRelationships(context.Background(), &v1.ReadRelationshipsRequest{})
	if _, err := client.PermissionsServiceClient.WriteRelationships(context.Background(), &v1.WriteRelationshipsRequest{
		OptionalPreconditions: []*v1.Precondition{
			{
				Operation: v1.Precondition_OPERATION_MUST_MATCH,
				Filter: &v1.RelationshipFilter{
					ResourceType:       auth.DomainType,
					OptionalResourceId: "domain_2",
					OptionalRelation:   "member",
					OptionalSubjectFilter: &v1.SubjectFilter{
						SubjectType: auth.UserType,
					},
				},
			},
			{
				Operation: v1.Precondition_OPERATION_MUST_MATCH,
				Filter: &v1.RelationshipFilter{
					ResourceType:     auth.GroupType,
					OptionalRelation: auth.DomainRelation,
					OptionalSubjectFilter: &v1.SubjectFilter{
						SubjectType:       auth.DomainType,
						OptionalSubjectId: "domain_2",
					},
				},
			},
			{
				Operation: v1.Precondition_OPERATION_MUST_MATCH,
				Filter: &v1.RelationshipFilter{
					ResourceType:     auth.ThingType,
					OptionalRelation: auth.DomainRelation,
					OptionalSubjectFilter: &v1.SubjectFilter{
						SubjectType:       auth.DomainType,
						OptionalSubjectId: "domain_2",
					},
				},
			},

			{
				Operation: v1.Precondition_OPERATION_MUST_NOT_MATCH,
				Filter: &v1.RelationshipFilter{
					ResourceType:     auth.ThingType,
					OptionalRelation: auth.DomainRelation,
					OptionalSubjectFilter: &v1.SubjectFilter{
						SubjectType: auth.DomainType,
					},
				},
			},
			{
				Operation: v1.Precondition_OPERATION_MUST_NOT_MATCH,
				Filter: &v1.RelationshipFilter{
					ResourceType:     auth.GroupType,
					OptionalRelation: auth.DomainRelation,
					OptionalSubjectFilter: &v1.SubjectFilter{
						SubjectType: auth.DomainType,
					},
				},
			},
		},
		Updates: []*v1.RelationshipUpdate{
			{
				Operation: v1.RelationshipUpdate_OPERATION_CREATE,
				Relationship: &v1.Relationship{
					Resource: &v1.ObjectReference{
						ObjectType: auth.DomainType,
						ObjectId:   "domain_1",
					},
					Relation: auth.AdministratorRelation,
					Subject: &v1.SubjectReference{
						Object: &v1.ObjectReference{
							ObjectType: auth.UserType,
							ObjectId:   "user_100",
						},
					},
				},
			},
			// {
			// 	Operation: v1.RelationshipUpdate_OPERATION_CREATE,
			// 	Relationship: &v1.Relationship{
			// 		Resource: &v1.ObjectReference{
			// 			ObjectType: auth.PlatformType,
			// 			ObjectId:   auth.MagistralaObject,
			// 		},
			// 		Relation: auth.AdministratorRelation,
			// 		Subject: &v1.SubjectReference{
			// 			Object: &v1.ObjectReference{
			// 				ObjectType: auth.UserType,
			// 				ObjectId:   "user_0",
			// 			},
			// 		},
			// 	},
			// },

			// {
			// 	Operation: v1.RelationshipUpdate_OPERATION_CREATE,
			// 	Relationship: &v1.Relationship{
			// 		Resource: &v1.ObjectReference{
			// 			ObjectType: auth.DomainType,
			// 			ObjectId:   "domain_1",
			// 		},
			// 		Relation: auth.AdministratorRelation,
			// 		Subject: &v1.SubjectReference{
			// 			Object: &v1.ObjectReference{
			// 				ObjectType: auth.UserType,
			// 				ObjectId:   "user_1",
			// 			},
			// 		},
			// 	},
			// },

			// {
			// 	Operation: v1.RelationshipUpdate_OPERATION_CREATE,
			// 	Relationship: &v1.Relationship{
			// 		Resource: &v1.ObjectReference{
			// 			ObjectType: auth.DomainType,
			// 			ObjectId:   "domain_2",
			// 		},
			// 		Relation: auth.AdministratorRelation,
			// 		Subject: &v1.SubjectReference{
			// 			Object: &v1.ObjectReference{
			// 				ObjectType: auth.UserType,
			// 				ObjectId:   "user_2",
			// 			},
			// 		},
			// 	},
			// },

			// {
			// 	Operation: v1.RelationshipUpdate_OPERATION_CREATE,
			// 	Relationship: &v1.Relationship{
			// 		Resource: &v1.ObjectReference{
			// 			ObjectType: auth.DomainType,
			// 			ObjectId:   "domain_3",
			// 		},
			// 		Relation: auth.AdministratorRelation,
			// 		Subject: &v1.SubjectReference{
			// 			Object: &v1.ObjectReference{
			// 				ObjectType: auth.UserType,
			// 				ObjectId:   "user_3",
			// 			},
			// 		},
			// 	},
			// },

			// {
			// 	Operation: v1.RelationshipUpdate_OPERATION_CREATE,
			// 	Relationship: &v1.Relationship{
			// 		Resource: &v1.ObjectReference{
			// 			ObjectType: auth.GroupType,
			// 			ObjectId:   "group_1",
			// 		},
			// 		Relation: auth.DomainRelation,
			// 		Subject: &v1.SubjectReference{
			// 			Object: &v1.ObjectReference{
			// 				ObjectType: auth.DomainType,
			// 				ObjectId:   "domain_1",
			// 			},
			// 		},
			// 	},
			// },

			// {
			// 	Operation: v1.RelationshipUpdate_OPERATION_CREATE,
			// 	Relationship: &v1.Relationship{
			// 		Resource: &v1.ObjectReference{
			// 			ObjectType: auth.GroupType,
			// 			ObjectId:   "group_2",
			// 		},
			// 		Relation: auth.DomainRelation,
			// 		Subject: &v1.SubjectReference{
			// 			Object: &v1.ObjectReference{
			// 				ObjectType: auth.DomainType,
			// 				ObjectId:   "domain_2",
			// 			},
			// 		},
			// 	},
			// },

			// {
			// 	Operation: v1.RelationshipUpdate_OPERATION_CREATE,
			// 	Relationship: &v1.Relationship{
			// 		Resource: &v1.ObjectReference{
			// 			ObjectType: auth.GroupType,
			// 			ObjectId:   "group_3",
			// 		},
			// 		Relation: auth.DomainRelation,
			// 		Subject: &v1.SubjectReference{
			// 			Object: &v1.ObjectReference{
			// 				ObjectType: auth.DomainType,
			// 				ObjectId:   "domain_4",
			// 			},
			// 		},
			// 	},
			// },

			// {
			// 	Operation: v1.RelationshipUpdate_OPERATION_CREATE,
			// 	Relationship: &v1.Relationship{
			// 		Resource: &v1.ObjectReference{
			// 			ObjectType: auth.DomainType,
			// 			ObjectId:   "domain_1",
			// 		},
			// 		Relation: auth.ViewerRelation,
			// 		Subject: &v1.SubjectReference{
			// 			Object: &v1.ObjectReference{
			// 				ObjectType: auth.UserType,
			// 				ObjectId:   "user_11",
			// 			},
			// 		},
			// 	},
			// },
			// {
			// 	Operation: v1.RelationshipUpdate_OPERATION_CREATE,
			// 	Relationship: &v1.Relationship{
			// 		Resource: &v1.ObjectReference{
			// 			ObjectType: auth.DomainType,
			// 			ObjectId:   "domain_2",
			// 		},
			// 		Relation: auth.ViewerRelation,
			// 		Subject: &v1.SubjectReference{
			// 			Object: &v1.ObjectReference{
			// 				ObjectType: auth.UserType,
			// 				ObjectId:   "user_11",
			// 			},
			// 		},
			// 	},
			// },
			// {
			// 	Operation: v1.RelationshipUpdate_OPERATION_CREATE,
			// 	Relationship: &v1.Relationship{
			// 		Resource: &v1.ObjectReference{
			// 			ObjectType: auth.DomainType,
			// 			ObjectId:   "domain_3",
			// 		},
			// 		Relation: auth.ViewerRelation,
			// 		Subject: &v1.SubjectReference{
			// 			Object: &v1.ObjectReference{
			// 				ObjectType: auth.UserType,
			// 				ObjectId:   "user_11",
			// 			},
			// 		},
			// 	},
			// },
		},
	}); err != nil {
		fmt.Println(err)
	}

	resp, err := client.CheckPermission(context.Background(), &v1.CheckPermissionRequest{
		Resource: &v1.ObjectReference{
			ObjectType: auth.DomainType,
			ObjectId:   "domain_1",
		},
		Permission: auth.AdminPermission,
		Subject: &v1.SubjectReference{
			Object: &v1.ObjectReference{
				ObjectType: auth.UserType,
				ObjectId:   "user_100",
			},
		},
	})
	fmt.Println(err)
	fmt.Println(resp)
	// resp, err := client.PermissionsServiceClient.ExpandPermissionTree(context.Background(), &v1.ExpandPermissionTreeRequest{
	// 	Resource: &v1.ObjectReference{
	// 		ObjectType: auth.GroupType,
	// 		ObjectId:   "group_1",
	// 	},
	// 	Permission: auth.ViewPermission,
	// })
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// _ = resp

	// by, err := json.MarshalIndent(resp, "", "  ")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(string(by))

	//	resp1, err := client.PermissionsServiceClient.ReadRelationships(context.Background(), &v1.ReadRelationshipsRequest{
	//		RelationshipFilter: &v1.RelationshipFilter{
	//			ResourceType:       auth.GroupType,
	//			OptionalResourceId: "group_1",
	//			OptionalRelation:   auth.ViewerRelation,
	//			OptionalSubjectFilter: &v1.SubjectFilter{
	//				SubjectType:      auth.UserType,
	//				OptionalRelation: &v1.SubjectFilter_RelationFilter{},
	//			},
	//		},
	//	})
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//
	// loop:
	//
	//	for {
	//		readRelationResp, err := resp1.Recv()
	//		switch err {
	//		case nil:
	//			fmt.Println(readRelationResp)
	//		case io.EOF:
	//			fmt.Println("got EOF while watch streaming")
	//			break loop
	//		default:
	//			fmt.Printf("got error while watch streaming : %s\n", err.Error())
	//			break loop
	//		}
	//	}
}
