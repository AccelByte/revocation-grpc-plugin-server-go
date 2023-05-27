// Copyright (c) 2023 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package revocationdemo

import (
	"fmt"

	"github.com/AccelByte/accelbyte-go-sdk/platform-sdk/pkg/platformclient/catalog_changes"
	"github.com/AccelByte/accelbyte-go-sdk/platform-sdk/pkg/platformclient/category"
	"github.com/AccelByte/accelbyte-go-sdk/platform-sdk/pkg/platformclient/currency"
	"github.com/AccelByte/accelbyte-go-sdk/platform-sdk/pkg/platformclient/item"
	"github.com/AccelByte/accelbyte-go-sdk/platform-sdk/pkg/platformclient/order"
	"github.com/AccelByte/accelbyte-go-sdk/platform-sdk/pkg/platformclient/revocation"
	"github.com/AccelByte/accelbyte-go-sdk/platform-sdk/pkg/platformclient/service_plugin_config"
	"github.com/AccelByte/accelbyte-go-sdk/platform-sdk/pkg/platformclient/store"
	"github.com/AccelByte/accelbyte-go-sdk/platform-sdk/pkg/platformclientmodels"
	"github.com/AccelByte/accelbyte-go-sdk/services-api/pkg/factory"
	"github.com/AccelByte/accelbyte-go-sdk/services-api/pkg/repository"
	"github.com/AccelByte/accelbyte-go-sdk/services-api/pkg/service/platform"
	"github.com/pkg/errors"
)

var (
	abStoreName = "GO Revocation Plugin Demo Store"
	abStoreDesc = "GO Description for revocation grpc plugin demo store"
)

var errEmptyStoreID = errors.New("error empty store id, createStore first")

type PlatformDataUnit struct {
	CLIConfig    *Config
	ConfigRepo   repository.ConfigRepository
	TokenRepo    repository.TokenRepository
	storeID      string
	CurrencyCode string
}

func (p *PlatformDataUnit) SetPlatformServiceGrpcTarget() error {
	grpcServerUrl := p.CLIConfig.GRPCServerURL
	if grpcServerUrl == "" {
		return errors.New("gRPC server url can't be empty")
	}

	wrapper := platform.ServicePluginConfigService{
		Client:           factory.NewPlatformClient(p.ConfigRepo),
		ConfigRepository: p.ConfigRepo,
		TokenRepository:  p.TokenRepo,
	}
	// call https://demo.accelbyte.io/platform/apidocs/#/ServicePluginConfig/updateServicePluginConfig
	_, err := wrapper.UpdateServicePluginConfigShort(&service_plugin_config.UpdateServicePluginConfigParams{
		Body: &platformclientmodels.ServicePluginConfigUpdate{
			GrpcServerAddress: grpcServerUrl,
		},
		Namespace: p.CLIConfig.ABNamespace,
	})
	return err
}

func (p *PlatformDataUnit) CreateStore(doPublish bool) error {
	storeWrapper := platform.StoreService{
		Client:           factory.NewPlatformClient(p.ConfigRepo),
		ConfigRepository: p.ConfigRepo,
		TokenRepository:  p.TokenRepo,
	}

	// Clean up existing stores
	storeInfo, err := storeWrapper.ListStoresShort(&store.ListStoresParams{
		Namespace: p.CLIConfig.ABNamespace,
	})
	if err != nil {
		return err
	}
	for _, s := range storeInfo {
		if Val(s.Published) == false {
			_, _ = storeWrapper.DeleteStoreShort(&store.DeleteStoreParams{
				Namespace: p.CLIConfig.ABNamespace,
				StoreID:   Val(s.StoreID),
			})
		}
	}

	// Create and publish new store
	newStore, err := storeWrapper.CreateStoreShort(&store.CreateStoreParams{
		Namespace: p.CLIConfig.ABNamespace,
		Body: &platformclientmodels.StoreCreate{
			DefaultLanguage:    "en",
			DefaultRegion:      "US",
			Description:        abStoreDesc,
			SupportedLanguages: []string{"en"},
			SupportedRegions:   []string{"US"},
			Title:              &abStoreName,
		},
	})
	if err != nil {
		return fmt.Errorf("could not create new store: %w", err)
	}

	p.storeID = Val(newStore.StoreID)
	if doPublish {
		err = p.publishStoreChange()
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PlatformDataUnit) publishStoreChange() error {
	catalogWrapper := platform.CatalogChangesService{
		Client:           factory.NewPlatformClient(p.ConfigRepo),
		ConfigRepository: p.ConfigRepo,
		TokenRepository:  p.TokenRepo,
	}
	_, err := catalogWrapper.PublishAllShort(&catalog_changes.PublishAllParams{
		Namespace: p.CLIConfig.ABNamespace,
		StoreID:   p.storeID,
	})
	if err != nil {
		return fmt.Errorf("could not publish store: %w", err)
	}
	return nil
}

func (p *PlatformDataUnit) CreateCategory(categoryPath string, doPublish bool) error {
	if p.storeID == "" {
		return errEmptyStoreID
	}

	categoryWrapper := platform.CategoryService{
		Client:           factory.NewPlatformClient(p.ConfigRepo),
		ConfigRepository: p.ConfigRepo,
		TokenRepository:  p.TokenRepo,
	}
	_, err := categoryWrapper.CreateCategoryShort(&category.CreateCategoryParams{
		Namespace: p.CLIConfig.ABNamespace,
		StoreID:   p.storeID,
		Body: &platformclientmodels.CategoryCreate{
			CategoryPath: &categoryPath,
			LocalizationDisplayNames: map[string]string{
				"en": categoryPath,
			},
		},
	})
	if err != nil {
		fmt.Errorf("could not create new category: %w", err)
	}
	return nil
}

func (p *PlatformDataUnit) UnsetPlatformServiceGrpcTarget() error {
	wrapper := platform.ServicePluginConfigService{
		Client:           factory.NewPlatformClient(p.ConfigRepo),
		ConfigRepository: p.ConfigRepo,
		TokenRepository:  p.TokenRepo,
	}
	return wrapper.DeleteServicePluginConfigShort(&service_plugin_config.DeleteServicePluginConfigParams{
		Namespace: p.CLIConfig.ABNamespace,
	})
}

func (p *PlatformDataUnit) DeleteCurrency() error {
	currencyWrapper := platform.CurrencyService{
		Client:           factory.NewPlatformClient(p.ConfigRepo),
		ConfigRepository: p.ConfigRepo,
		TokenRepository:  p.TokenRepo,
	}
	_, err := currencyWrapper.DeleteCurrencyShort(&currency.DeleteCurrencyParams{
		Namespace:    p.CLIConfig.ABNamespace,
		CurrencyCode: p.CurrencyCode,
	})
	return err
}

func (p *PlatformDataUnit) DeleteStore() error {
	storeWrapper := platform.StoreService{
		Client:           factory.NewPlatformClient(p.ConfigRepo),
		ConfigRepository: p.ConfigRepo,
		TokenRepository:  p.TokenRepo,
	}
	_, err := storeWrapper.DeleteStoreShort(&store.DeleteStoreParams{
		Namespace: p.CLIConfig.ABNamespace,
		StoreID:   p.storeID,
	})
	return err
}

func (p *PlatformDataUnit) UpdateRevocationConfig() error {
	revocationWrapper := platform.RevocationService{
		Client:           factory.NewPlatformClient(p.ConfigRepo),
		ConfigRepository: p.ConfigRepo,
		TokenRepository:  p.TokenRepo,
	}
	_, err := revocationWrapper.UpdateRevocationConfigShort(&revocation.UpdateRevocationConfigParams{
		Namespace: p.CLIConfig.ABNamespace,
		Body: &platformclientmodels.RevocationConfigUpdate{
			Entitlement: &platformclientmodels.EntitlementRevocationConfig{
				Consumable: nil,
				Durable: &platformclientmodels.DurableEntitlementRevocationConfig{
					Enabled:  false,
					Strategy: platformclientmodels.DurableEntitlementRevocationConfigStrategyCUSTOM,
				},
			},
			Wallet: &platformclientmodels.WalletRevocationConfig{
				Enabled:  true,
				Strategy: platformclientmodels.WalletRevocationConfigStrategyCUSTOM,
			},
		},
	})
	return err
}

func (p *PlatformDataUnit) CreateCurrency() error {
	currencyWrapper := platform.CurrencyService{
		Client:           factory.NewPlatformClient(p.ConfigRepo),
		ConfigRepository: p.ConfigRepo,
		TokenRepository:  p.TokenRepo,
	}
	_, err := currencyWrapper.CreateCurrencyShort(&currency.CreateCurrencyParams{
		Namespace: p.CLIConfig.ABNamespace,
		Body: &platformclientmodels.CurrencyCreate{
			CurrencyCode:   Ptr(p.CurrencyCode),
			CurrencySymbol: "$V",
			CurrencyType:   platformclientmodels.CurrencyCreateCurrencyTypeREAL,
			Decimals:       2,
		},
	})

	return err
}

func (p *PlatformDataUnit) CreateItems(itemCount int, categoryPath string, doPublish bool) ([]SimpleItemInfo, error) {
	if p.storeID == "" {
		return nil, errEmptyStoreID
	}

	itemWrapper := platform.ItemService{
		Client:           factory.NewPlatformClient(p.ConfigRepo),
		ConfigRepository: p.ConfigRepo,
		TokenRepository:  p.TokenRepo,
	}
	var items []SimpleItemInfo
	itemDiff := RandomString("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 6)
	for i := 0; i < itemCount; i++ {
		itemInfo := SimpleItemInfo{
			Title: fmt.Sprintf("Item %s Titled %d", itemDiff, i+1),
			SKU:   fmt.Sprintf("SKU_%s_%d", itemDiff, i+1),
		}

		newItem, err := itemWrapper.CreateItemShort(&item.CreateItemParams{
			Namespace: p.CLIConfig.ABNamespace,
			StoreID:   p.storeID,
			Body: &platformclientmodels.ItemCreate{
				Name:               &itemInfo.Title,
				ItemType:           Ptr(platformclientmodels.ItemCreateItemTypeCOINS),
				CategoryPath:       &categoryPath,
				EntitlementType:    Ptr(platformclientmodels.ItemCreateEntitlementTypeCONSUMABLE),
				SeasonType:         platformclientmodels.ItemCreateSeasonTypeTIER,
				Status:             Ptr(platformclientmodels.ItemCreateStatusACTIVE),
				UseCount:           1,
				TargetCurrencyCode: p.CurrencyCode,
				Listable:           true,
				Purchasable:        true,
				Sku:                itemInfo.SKU,
				Localizations: map[string]platformclientmodels.Localization{
					"en": {
						Title: Ptr(itemInfo.Title),
					},
				},
				RegionData: map[string][]platformclientmodels.RegionDataItemDTO{
					"US": {
						{
							CurrencyCode:      Ptr(p.CurrencyCode),
							CurrencyNamespace: Ptr(p.CLIConfig.ABNamespace),
							CurrencyType:      Ptr(platformclientmodels.RegionDataItemDTOCurrencyTypeVIRTUAL),
							Price:             Ptr(int32(0)),
						},
					},
				},
			},
		})
		if err != nil {
			return nil, fmt.Errorf("could not create new store item: %w", err)
		}

		itemInfo.ID = *newItem.ItemID
		items = append(items, itemInfo)
	}

	if doPublish {
		_ = p.publishStoreChange()
	}

	return items, nil
}

func (p *PlatformDataUnit) CreateOrder(userID string, itemInfo SimpleItemInfo) (*platformclientmodels.OrderInfo, error) {
	orderWrapper := platform.OrderService{
		Client:           factory.NewPlatformClient(p.ConfigRepo),
		ConfigRepository: p.ConfigRepo,
		TokenRepository:  p.TokenRepo,
	}
	return orderWrapper.PublicCreateUserOrderShort(&order.PublicCreateUserOrderParams{
		Namespace: p.CLIConfig.ABNamespace,
		UserID:    userID,
		Body: &platformclientmodels.OrderCreate{
			CurrencyCode:    Ptr(p.CurrencyCode),
			ItemID:          Ptr(itemInfo.ID),
			Price:           Ptr(int32(0)),
			Quantity:        Ptr(int32(1)),
			DiscountedPrice: Ptr(int32(0)),
		},
	})
}

func (p *PlatformDataUnit) Revoke(userID string, orderNo string, itemID string) (*platformclientmodels.RevocationResult, error) {
	revocationWrapper := platform.RevocationService{
		Client:           factory.NewPlatformClient(p.ConfigRepo),
		ConfigRepository: p.ConfigRepo,
		TokenRepository:  p.TokenRepo,
	}
	return revocationWrapper.DoRevocationShort(&revocation.DoRevocationParams{
		Namespace: p.CLIConfig.ABNamespace,
		UserID:    userID,
		Body: &platformclientmodels.RevocationRequest{
			Source:        platformclientmodels.RevocationRequestSourceORDER,
			TransactionID: orderNo,
			RevokeEntries: []*platformclientmodels.RevokeEntry{
				{
					Quantity: 1,
					Type:     platformclientmodels.RevokeEntryTypeITEM,
					Item: &platformclientmodels.RevokeItem{
						ItemIdentityType: platformclientmodels.RevokeItemItemIdentityTypeITEMID,
						ItemIdentity:     itemID,
					},
				},
			},
		},
	})
}
