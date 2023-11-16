package main

import (
	"fmt"
	"math/rand"
	"time"

	mgSDK "github.com/absmach/magistrala/pkg/sdk/go"
)

func main() {

	runID := rand.Intn(10000)
	sdk := mgSDK.NewSDK(mgSDK.Config{
		ThingsURL:  "http://localhost",
		UsersURL:   "http://localhost",
		DomainsURL: "http://localhost:8189",
		HostURL:    "http://localhost",
	})

	token, err := sdk.CreateToken(mgSDK.Login{
		Identity: "user1@example.com",
		Secret:   "12345678",
	})
	if err != nil {
		panic(err)
	}

	d := mgSDK.Domain{
		Name:  fmt.Sprintf("domain_%d", runID),
		Tags:  []string{fmt.Sprintf("tag_%d", runID)},
		Alias: fmt.Sprintf("domain_%d", runID),
	}
	od, err := sdk.CreateDomain(d, token.AccessToken)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", od)

	token, err = sdk.RefreshToken(mgSDK.Login{DomainID: od.ID}, token.RefreshToken)
	if err != nil {
		panic(err)
	}

	od, err = sdk.UpdateDomain(od, token.AccessToken)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", od)

	od, err = sdk.Domain(od.ID, token.AccessToken)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", od)

	if err := sdk.DisableDomain(od.ID, token.AccessToken); err != nil {
		panic(err)
	}

	od, err = sdk.Domain(od.ID, token.AccessToken)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", od)

	if err := sdk.EnableDomain(od.ID, token.AccessToken); err != nil {
		panic(err)
	}

	dp, err := sdk.Domains(mgSDK.PageMetadata{}, token.AccessToken)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", dp)

	user2ID := "ee488271-3450-4471-9c83-53f5280c28de"

	if err := sdk.AddUserToDomain(od.ID, mgSDK.UsersRelationRequest{Relation: "viewer", UserIDs: []string{user2ID}}, token.AccessToken); err != nil {
		panic(err)
	}

	time.Sleep(2 * time.Second)
	users, err := sdk.ListDomainUsers(od.ID, mgSDK.PageMetadata{}, token.AccessToken)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", users)

	// if err := sdk.RemoveUserFromDomain(od.ID, mgSDK.UsersRelationRequest{Relation: "viewer", UserIDs: []string{user2ID}}, token.AccessToken); err != nil {
	// 	panic(err)
	// }
	// time.Sleep(4 * time.Second)

	// users, err = sdk.ListDomainUsers(od.ID, mgSDK.PageMetadata{}, token.AccessToken)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("%+v\n", users)

	token, err = sdk.CreateToken(mgSDK.Login{Identity: "admin@example.com", Secret: "12345678"})
	if err != nil {
		panic(err)
	}

	do, err := sdk.ListUserDomains(user2ID, mgSDK.PageMetadata{}, token.AccessToken)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", do)

}
