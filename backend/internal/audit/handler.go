package audit

import (
	"context"

	"github.com/Paukku/ajanvarausjarjestelma/backend/pb/common"
)

type AuditHandler struct {
	service *Service
}

func (h *AuditHandler) GetAuditLogs(
	ctx context.Context,
	req *common.GetAuditLogsRequest,
) (*common.GetAuditLogsResponse, error) {

	limit := int(req.Limit)
	offset := int(req.Offset)

	logs, err := h.service.GetLogs(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	var pbLogs []*common.AuditLog
	for _, l := range logs {
		pbLogs = append(pbLogs, mapLogToProto(l))
	}

	return &common.GetAuditLogsResponse{
		Logs: pbLogs,
	}, nil
}
