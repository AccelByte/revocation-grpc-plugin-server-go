// Copyright (c) 2023 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/AccelByte/accelbyte-go-sdk/iam-sdk/pkg/iamclient/users"
	"github.com/AccelByte/accelbyte-go-sdk/iam-sdk/pkg/iamclient/users_v4"
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
	err = oauthService.LoginClient(&config.ABClientID, &config.ABClientSecret)
	if err != nil {
		log.Fatalf("Accelbyte account login failed: %s\n", err)
	}
	fmt.Println("[OK]")

	usersService := &iam.UsersV4Service{
		Client:           factory.NewIamClient(configRepo),
		ConfigRepository: configRepo,
		TokenRepository:  tokenRepo,
	}
	verified := true
	nameId := revocationdemo.RandomString("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 8)
	dName := "Extend Test User " + nameId
	username := fmt.Sprintf("extend_%s_user", nameId)
	email := username + "@dummy.net"
	country := "ID"
	dob := "1990-01-01"
	password := revocationdemo.RandomString("ABCDEFGHIJKlmnopqrstuvwxyz0123456789!@#$%^&", 16)

	var acceptedPolicies []*iamclientmodels.LegalAcceptedPoliciesRequest
	authType := iamclientmodels.AccountCreateTestUserRequestV4AuthTypeEMAILPASSWD
	userInfo, err := usersService.PublicCreateTestUserV4Short(&users_v4.PublicCreateTestUserV4Params{
		Body: &iamclientmodels.AccountCreateTestUserRequestV4{
			AcceptedPolicies:  acceptedPolicies,
			AuthType:          &authType,
			Country:           &country,
			DateOfBirth:       &dob,
			DisplayName:       &dName,
			EmailAddress:      &email,
			Password:          &password,
			UniqueDisplayName: dName,
			Username:          &username,
			Verified:          &verified,
		},
		Namespace: os.Getenv("AB_NAMESPACE"),
	})
	if err != nil {
		log.Fatalf("Get user info failed: %s\n", err)
	}
	fmt.Printf("Test User Created: %s\n", *userInfo.UserID)

	// Start testing
	err = startTesting(userInfo, config, configRepo, tokenRepo)
	if err != nil {
		fmt.Println("\n[FAILED]")
		log.Fatal(err)
	}
	fmt.Println("\n[SUCCESS]")
}

func startTesting(
	userInfo *iamclientmodels.AccountCreateUserResponseV4,
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

		fmt.Print("Deleting Test User... ")
		err = deleteUser(userInfo, configRepo, tokenRepo)
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
		fmt.Println("[ERR]")

		return err
	}
	fmt.Println("[OK]")

	// 2.
	fmt.Print("Creating store... ")
	err = pdu.CreateStore(true)
	if err != nil {
		fmt.Println("[ERR]")

		return err
	}
	fmt.Println("[OK]")

	// 3.
	fmt.Print("Creating category... ")
	err = pdu.CreateCategory(categoryPath, true)
	if err != nil {
		fmt.Println("[ERR]")

		return err
	}
	fmt.Println("[OK]")

	// 4.
	fmt.Print("Updating Revocation config... ")
	err = pdu.UpdateRevocationConfig()
	if err != nil {
		fmt.Println("[ERR]")

		return err
	}
	fmt.Println("[OK]")

	// 5.
	fmt.Print("Setting up virtual currency [VCA]... ")
	err = pdu.CreateCurrency()
	if err != nil {
		fmt.Println("[ERR]")

		return err
	}
	fmt.Println("[OK]")

	// 6.
	fmt.Print("Creating items...")
	itemInfos, err := pdu.CreateItems(1, categoryPath, true)
	if err != nil {
		fmt.Println("[ERR]")

		return err
	}
	fmt.Println("[OK]")

	// 7.
	fmt.Print("Creating order...")
	orderInfo, err := pdu.CreateOrder(revocationdemo.Val(userInfo.UserID), itemInfos[0])
	if err != nil {
		fmt.Println("[ERR]")

		return err
	}
	fmt.Println("[OK]")

	// 8.
	fmt.Print("Revoking order...")
	revocationResult, err := pdu.Revoke(revocationdemo.Val(userInfo.UserID), revocationdemo.Val(orderInfo.OrderNo), revocationdemo.Val(orderInfo.ItemID))
	if err != nil {
		fmt.Println("[ERR]")

		return err
	}
	fmt.Println("[OK]")

	revocationStatus := revocationdemo.Val(revocationResult.Status)
	fmt.Println("Revocation Result: ")
	fmt.Printf("Revocation history id: %s\n", revocationdemo.Val(revocationResult.ID))
	fmt.Printf("Revocation status: %s\n", revocationStatus)

	if strings.ToUpper(revocationStatus) == "FAIL" {
		return fmt.Errorf("revocation status is FAIL")
	}

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

func deleteUser(
	userInfo *iamclientmodels.AccountCreateUserResponseV4,
	configRepo repository.ConfigRepository,
	tokenRepo repository.TokenRepository) error {
	userService := &iam.UsersService{
		Client:           factory.NewIamClient(configRepo),
		ConfigRepository: configRepo,
		TokenRepository:  tokenRepo,
	}
	errDelete := userService.AdminDeleteUserInformationV3Short(&users.AdminDeleteUserInformationV3Params{
		Namespace: *userInfo.Namespace,
		UserID:    *userInfo.UserID,
	})
	if errDelete != nil {
		return errDelete
	}

	return nil
}
