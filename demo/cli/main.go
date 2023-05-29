// Copyright (c) 2023 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package main

import (
	"fmt"
	"log"

	"github.com/AccelByte/accelbyte-go-sdk/iam-sdk/pkg/iamclient/users"
	"github.com/AccelByte/accelbyte-go-sdk/iam-sdk/pkg/iamclientmodels"
	"github.com/AccelByte/accelbyte-go-sdk/services-api/pkg/factory"
	"github.com/AccelByte/accelbyte-go-sdk/services-api/pkg/repository"
	"github.com/AccelByte/accelbyte-go-sdk/services-api/pkg/service/iam"
	"github.com/AccelByte/accelbyte-go-sdk/services-api/pkg/utils/auth"

	revocationdemo "revocation-grpc-plugin-server-go-cli/pkg"
)

func main() {
	config, err := revocationdemo.GetConfig()
	if err != nil {
		log.Fatalf("Can't retrieve config: %s\n", err)
	}

	configRepo := auth.DefaultConfigRepositoryImpl()
	tokenRepo := auth.DefaultTokenRepositoryImpl()

	oauthService := &iam.OAuth20Service{
		Client:           factory.NewIamClient(configRepo),
		ConfigRepository: configRepo,
		TokenRepository:  tokenRepo,
	}

	fmt.Print("Login to AccelByte... ")
	err = oauthService.Login(config.ABUsername, config.ABPassword)
	if err != nil {
		log.Fatalf("Accelbyte account login failed: %s\n", err)
	}
	fmt.Println("[OK]")

	usersService := &iam.UsersService{
		Client:           factory.NewIamClient(configRepo),
		ConfigRepository: configRepo,
		TokenRepository:  tokenRepo,
	}
	userInfo, err := usersService.PublicGetMyUserV3Short(&users.PublicGetMyUserV3Params{})
	if err != nil {
		log.Fatalf("Get user info failed: %s\n", err)
	}
	fmt.Printf("User: %s\n", userInfo.UserName)

	// Start testing
	err = startTesting(userInfo, config, configRepo, tokenRepo)
	if err != nil {
		fmt.Println("\n[FAILED]")
		log.Fatal(err)
	}
}

func startTesting(
	userInfo *iamclientmodels.ModelUserResponseV3,
	config *revocationdemo.Config,
	configRepo repository.ConfigRepository,
	tokenRepo repository.TokenRepository) error {
	categoryPath := "/goRevocationPluginDemo"
	pdu := revocationdemo.PlatformDataUnit{
		CLIConfig:    config,
		ConfigRepo:   configRepo,
		TokenRepo:    tokenRepo,
		CurrencyCode: "VCA",
	}

	// clean up
	defer func() {
		fmt.Println("\nCleaning up...")
		fmt.Print("Deleting currency[VCA]... ")
		err := pdu.DeleteCurrency()
		if err != nil {
			return
		}
		fmt.Println("[OK]")

		fmt.Print("Deleting store... ")
		err = pdu.DeleteStore()
		if err != nil {
			return
		}
		fmt.Println("[OK]")

		err = pdu.UnsetPlatformServiceGrpcTarget()
		if err != nil {
			fmt.Printf("failed to unset platform service grpc plugin url")

			return
		}
	}()

	// 1.
	fmt.Print("Configuring platform service grpc target... ")
	err := pdu.SetPlatformServiceGrpcTarget()
	if err != nil {
		return err
	}
	fmt.Println("[OK]")

	// 2.
	fmt.Print("Creating store... ")
	err = pdu.CreateStore(true)
	if err != nil {
		return err
	}
	fmt.Println("[OK]")

	// 3.
	fmt.Print("Creating category... ")
	err = pdu.CreateCategory(categoryPath, true)
	if err != nil {
		return err
	}
	fmt.Println("[OK]")

	// 4.
	fmt.Print("Updating Revocation config... ")
	err = pdu.UpdateRevocationConfig()
	if err != nil {
		return err
	}
	fmt.Println("[OK]")

	// 5.
	fmt.Print("Setting up virtual currency [VCA]... ")
	err = pdu.CreateCurrency()
	if err != nil {
		return err
	}
	fmt.Println("[OK]")

	// 6.
	fmt.Print("Creating items...")
	itemInfos, err := pdu.CreateItems(1, categoryPath, true)
	if err != nil {
		return err
	}
	fmt.Println("[OK]")

	// 7.
	fmt.Print("Creating order...")
	orderInfo, err := pdu.CreateOrder(revocationdemo.Val(userInfo.UserID), itemInfos[0])
	if err != nil {
		return err
	}
	fmt.Println("[OK]")

	// 8.
	fmt.Print("Revoking order...")
	revocationResult, err := pdu.Revoke(revocationdemo.Val(userInfo.UserID), revocationdemo.Val(orderInfo.OrderNo), revocationdemo.Val(orderInfo.ItemID))
	if err != nil {
		return err
	}
	fmt.Println("[OK]")

	fmt.Println("Revocation Result: ")
	fmt.Printf("Revocation history id: %s\n", revocationdemo.Val(revocationResult.ID))
	fmt.Printf("Revocation status: %s\n", revocationdemo.Val(revocationResult.Status))

	for _, r := range revocationResult.ItemRevocations {
		fmt.Printf("item Id: %s\n", r.ItemID)
		fmt.Printf("item sku: %s\n", r.ItemSku)
		fmt.Printf("item type: %s\n", r.ItemType)
		fmt.Printf("quantity: %d\n", r.Quantity)
		fmt.Printf("revocation strategy: %s\n", r.Strategy)
		fmt.Printf("skipped: %t\n", r.Skipped)
		fmt.Printf("reason: %s\n", r.Reason)
		fmt.Printf("custom revocation: %s\n", r.CustomRevocation)
	}

	return nil
}
