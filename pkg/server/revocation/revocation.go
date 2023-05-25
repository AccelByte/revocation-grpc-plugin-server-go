// Copyright (c) 2023 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package revocation

import (
	"fmt"

	pb "lootbox-roll-function-grpc-plugin-server-go/pkg/pb"
)

type RevokeEntryType string

const (
	StatusFail    = "FAIL"
	StatusSuccess = "SUCCESS"
)

const (
	RevokeEntryTypeItem        RevokeEntryType = "ITEM"
	RevokeEntryTypeCurrency    RevokeEntryType = "CURRENCY"
	RevokeEntryTypeEntitlement RevokeEntryType = "ENTITLEMENT"
)

type Revocation interface {
	Revoke(namespace string, userId string, quantity int32, request *pb.RevokeRequest) (*pb.RevokeResponse, error)
}

type ItemRevocation struct {
}

func (r *ItemRevocation) Revoke(namespace string, userId string, quantity int32, request *pb.RevokeRequest) (*pb.RevokeResponse, error) {
	item := request.GetItem()
	customRevocation := map[string]string{}
	customRevocation["namespace"] = namespace
	customRevocation["userId"] = userId
	customRevocation["quantity"] = fmt.Sprintf("%d", quantity)
	customRevocation["itemId"] = item.GetItemId()
	customRevocation["sku"] = item.GetItemSku()
	customRevocation["itemType"] = item.GetItemType()
	customRevocation["useCount"] = fmt.Sprintf("%d", item.GetUseCount())
	customRevocation["entitlementType"] = item.GetEntitlementType()

	return &pb.RevokeResponse{
		Status:           StatusSuccess,
		CustomRevocation: customRevocation,
	}, nil
}

type CurrencyRevocation struct {
}

func (r *CurrencyRevocation) Revoke(namespace string, userId string, quantity int32, request *pb.RevokeRequest) (*pb.RevokeResponse, error) {
	currency := request.GetCurrency()
	customRevocation := map[string]string{}
	customRevocation["namespace"] = namespace
	customRevocation["userId"] = userId
	customRevocation["quantity"] = fmt.Sprintf("%d", quantity)
	customRevocation["currencyNamespace"] = currency.GetNamespace()
	customRevocation["currencyCode"] = currency.GetCurrencyCode()
	customRevocation["balanceOrigin"] = currency.GetBalanceOrigin()

	return &pb.RevokeResponse{
		Status:           StatusSuccess,
		CustomRevocation: customRevocation,
	}, nil
}

type EntitlementRevocation struct {
}

func (r *EntitlementRevocation) Revoke(namespace string, userId string, quantity int32, request *pb.RevokeRequest) (*pb.RevokeResponse, error) {
	entitlement := request.GetEntitlement()
	customRevocation := map[string]string{}
	customRevocation["namespace"] = namespace
	customRevocation["userId"] = userId
	customRevocation["quantity"] = fmt.Sprintf("%d", quantity)
	customRevocation["entitlementId"] = entitlement.GetEntitlementId()
	customRevocation["itemId"] = entitlement.GetItemId()
	customRevocation["sku"] = entitlement.GetSku()

	return &pb.RevokeResponse{
		Status:           StatusSuccess,
		CustomRevocation: customRevocation,
	}, nil
}

var revocations = map[RevokeEntryType]Revocation{
	RevokeEntryTypeItem:        &ItemRevocation{},
	RevokeEntryTypeCurrency:    &CurrencyRevocation{},
	RevokeEntryTypeEntitlement: &EntitlementRevocation{},
}

func GetRevocation(revocationType RevokeEntryType) (Revocation, error) {
	if revocation, ok := revocations[revocationType]; ok {
		return revocation, nil
	}

	return nil, fmt.Errorf("revocation type '%s' is not supported", revocationType)
}
