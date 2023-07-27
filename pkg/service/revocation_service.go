// Copyright (c) 2023 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package service

import (
	"context"

	pb "revocation-grpc-plugin-server-go/pkg/pb"
	"revocation-grpc-plugin-server-go/pkg/service/revocation"
)

type RevocationServiceServer struct {
	pb.UnimplementedRevocationServer
}

func NewRevocationServiceServer() *RevocationServiceServer {
	return &RevocationServiceServer{}
}

func (s *RevocationServiceServer) Revoke(_ context.Context, req *pb.RevokeRequest) (*pb.RevokeResponse, error) {
	revocationEntryType := revocation.RevokeEntryType(req.GetRevokeEntryType())
	revocationObj, err := revocation.GetRevocation(revocationEntryType)
	if err != nil {
		return &pb.RevokeResponse{
			Status: revocation.StatusFail,
			Reason: err.Error(),
		}, nil
	}

	revocationResp, err := revocationObj.Revoke(req.GetNamespace(), req.GetUserId(), req.GetQuantity(), req)
	if err != nil {
		return &pb.RevokeResponse{
			Status: revocation.StatusFail,
			Reason: err.Error(),
		}, nil
	}
	return revocationResp, nil
}
