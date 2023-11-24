// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

// Package main contains cli main function to run the cli.
package main

import (
	"fmt"
	"log"

	"github.com/absmach/magistrala/cli"
	sdk "github.com/absmach/magistrala/pkg/sdk/go"
	"github.com/spf13/cobra"
)

const (
	defURL          string = "http://localhost"
	defUsersURL     string = defURL + ":9002"
	defThingsURL    string = defURL + ":9000"
	defBootstrapURL string = defURL + ":9013"
	defDomainsURL   string = defURL + ":8189"
	defCertsURL     string = defURL + ":9019"
)

func main() {
	msgContentType := string(sdk.CTJSONSenML)
	sdkConf := sdk.Config{
		ThingsURL:       defThingsURL,
		UsersURL:        defUsersURL,
		ReaderURL:       defURL,
		HTTPAdapterURL:  fmt.Sprintf("%s/http", defURL),
		BootstrapURL:    defBootstrapURL,
		CertsURL:        defCertsURL,
		DomainsURL:      defDomainsURL,
		MsgContentType:  sdk.ContentType(msgContentType),
		TLSVerification: false,
		HostURL:         defURL,
	}

	// Root
	rootCmd := &cobra.Command{
		Use: "magistrala-cli",
		PersistentPreRun: func(_ *cobra.Command, _ []string) {
			cliConf, err := cli.ParseConfig(sdkConf)
			if err != nil {
				log.Fatalf("Failed to parse config: %s", err)
			}
			if cliConf.MsgContentType == "" {
				cliConf.MsgContentType = sdk.ContentType(msgContentType)
			}
			s := sdk.NewSDK(cliConf)
			cli.SetSDK(s)
		},
	}
	// API commands
	healthCmd := cli.NewHealthCmd()
	usersCmd := cli.NewUsersCmd()
	domainsCmd := cli.NewDomainsCmd()
	thingsCmd := cli.NewThingsCmd()
	groupsCmd := cli.NewGroupsCmd()
	channelsCmd := cli.NewChannelsCmd()
	messagesCmd := cli.NewMessagesCmd()
	provisionCmd := cli.NewProvisionCmd()
	bootstrapCmd := cli.NewBootstrapCmd()
	certsCmd := cli.NewCertsCmd()
	subscriptionsCmd := cli.NewSubscriptionCmd()
	configCmd := cli.NewConfigCmd()

	// Root Commands
	rootCmd.AddCommand(healthCmd)
	rootCmd.AddCommand(usersCmd)
	rootCmd.AddCommand(domainsCmd)
	rootCmd.AddCommand(groupsCmd)
	rootCmd.AddCommand(thingsCmd)
	rootCmd.AddCommand(channelsCmd)
	rootCmd.AddCommand(messagesCmd)
	rootCmd.AddCommand(provisionCmd)
	rootCmd.AddCommand(bootstrapCmd)
	rootCmd.AddCommand(certsCmd)
	rootCmd.AddCommand(subscriptionsCmd)
	rootCmd.AddCommand(configCmd)

	// Root Flags
	rootCmd.PersistentFlags().StringVarP(
		&sdkConf.BootstrapURL,
		"bootstrap-url",
		"b",
		sdkConf.BootstrapURL,
		"Bootstrap service URL",
	)

	rootCmd.PersistentFlags().StringVarP(
		&sdkConf.CertsURL,
		"certs-url",
		"s",
		sdkConf.CertsURL,
		"Certs service URL",
	)

	rootCmd.PersistentFlags().StringVarP(
		&sdkConf.ThingsURL,
		"things-url",
		"t",
		sdkConf.ThingsURL,
		"Things service URL",
	)

	rootCmd.PersistentFlags().StringVarP(
		&sdkConf.UsersURL,
		"users-url",
		"u",
		sdkConf.UsersURL,
		"Users service URL",
	)

	rootCmd.PersistentFlags().StringVarP(
		&sdkConf.DomainsURL,
		"domains-url",
		"d",
		sdkConf.DomainsURL,
		"Domains service URL",
	)

	rootCmd.PersistentFlags().StringVarP(
		&sdkConf.HTTPAdapterURL,
		"http-url",
		"p",
		sdkConf.HTTPAdapterURL,
		"HTTP adapter URL",
	)

	rootCmd.PersistentFlags().StringVarP(
		&sdkConf.ReaderURL,
		"reader-url",
		"R",
		sdkConf.ReaderURL,
		"Reader URL",
	)

	rootCmd.PersistentFlags().StringVarP(
		&sdkConf.HostURL,
		"host-url",
		"H",
		sdkConf.HostURL,
		"Host URL",
	)

	rootCmd.PersistentFlags().StringVarP(
		&msgContentType,
		"content-type",
		"y",
		msgContentType,
		"Message content type",
	)

	rootCmd.PersistentFlags().BoolVarP(
		&sdkConf.TLSVerification,
		"insecure",
		"i",
		sdkConf.TLSVerification,
		"Do not check for TLS cert",
	)

	rootCmd.PersistentFlags().StringVarP(
		&cli.ConfigPath,
		"config",
		"c",
		cli.ConfigPath,
		"Config path",
	)

	rootCmd.PersistentFlags().BoolVarP(
		&cli.RawOutput,
		"raw",
		"r",
		cli.RawOutput,
		"Enables raw output mode for easier parsing of output",
	)

	// Client and Channels Flags
	rootCmd.PersistentFlags().Uint64VarP(
		&cli.Limit,
		"limit",
		"l",
		10,
		"Limit query parameter",
	)

	rootCmd.PersistentFlags().Uint64VarP(
		&cli.Offset,
		"offset",
		"o",
		0,
		"Offset query parameter",
	)

	rootCmd.PersistentFlags().StringVarP(
		&cli.Name,
		"name",
		"n",
		"",
		"Name query parameter",
	)

	rootCmd.PersistentFlags().StringVarP(
		&cli.Identity,
		"identity",
		"I",
		"",
		"User identity query parameter",
	)

	rootCmd.PersistentFlags().StringVarP(
		&cli.Metadata,
		"metadata",
		"m",
		"",
		"Metadata query parameter",
	)

	rootCmd.PersistentFlags().StringVarP(
		&cli.Status,
		"status",
		"S",
		"",
		"User status query parameter",
	)

	rootCmd.PersistentFlags().StringVarP(
		&cli.State,
		"state",
		"z",
		"",
		"Bootstrap state query parameter",
	)

	rootCmd.PersistentFlags().StringVarP(
		&cli.Topic,
		"topic",
		"T",
		"",
		"Subscription topic query parameter",
	)

	rootCmd.PersistentFlags().StringVarP(
		&cli.Contact,
		"contact",
		"C",
		"",
		"Subscription contact query parameter",
	)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
