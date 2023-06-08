// Copyright (c) 2023 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package platformservice

import (
	"context"
	"net/url"
	"time"

	"github.com/AccelByte/accelbyte-go-sdk/services-api/pkg/repository"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	"revocation-grpc-plugin-server-go-cli/pkg/client/platformservice/openapi2/client"
	"revocation-grpc-plugin-server-go-cli/pkg/client/platformservice/openapi2/client/service_plugin_config"
	"revocation-grpc-plugin-server-go-cli/pkg/client/platformservice/openapi2/models"
)

type Client struct {
	tokenRepo         repository.TokenRepository
	platformSvcClient *client.JusticePlatformService
}

func NewClient(serviceAddress string, tokenRepo repository.TokenRepository) (*Client, error) {
	u, err := url.Parse(serviceAddress)
	if err != nil {
		return nil, err
	}
	platformSvcClient := client.New(httptransport.New(u.Host, "platform", []string{u.Scheme}), strfmt.Default)

	return &Client{
		tokenRepo:         tokenRepo,
		platformSvcClient: platformSvcClient,
	}, nil
}

func (c *Client) UpdateRevocationPluginConfig(namespace string, config *models.RevocationPluginConfigUpdate) error {
	token, err := c.tokenRepo.GetToken()
	if err != nil {
		return err
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()

	bearerToken := httptransport.BearerToken(*token.AccessToken)
	_, err = c.platformSvcClient.ServicePluginConfig.UpdateLootBoxPluginConfig1(&service_plugin_config.UpdateLootBoxPluginConfig1Params{
		Namespace: namespace,
		Body:      config,
		Context:   ctx,
	}, bearerToken)

	return err
}

func (c *Client) DeleteRevocationPluginConfig(namespace string) error {
	token, err := c.tokenRepo.GetToken()
	if err != nil {
		return err
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()

	bearerToken := httptransport.BearerToken(*token.AccessToken)
	_, err = c.platformSvcClient.ServicePluginConfig.DeleteLootBoxPluginConfig1(&service_plugin_config.DeleteLootBoxPluginConfig1Params{
		Namespace: namespace,
		Context:   ctx,
	}, bearerToken)

	return err
}
