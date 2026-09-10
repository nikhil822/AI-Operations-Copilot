package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type Registry struct {
	OrderTool    *OrderTool
	PaymentTool  *PaymentTool
	DeliveryTool *DeliveryTool
	SummaryTool  *SummaryTool
}

func NewRegistry(
	orderTool *OrderTool,
	paymentTool *PaymentTool,
	deliveryTool *DeliveryTool,
	summaryTool *SummaryTool,
) *Registry {
	return &Registry{
		OrderTool:    orderTool,
		PaymentTool:  paymentTool,
		DeliveryTool: deliveryTool,
		SummaryTool:  summaryTool,
	}
}

type ToolCallResult struct {
	Name   string
	Result any
	Error  error
}

type toolArguments struct {
	OrderID uint `json:"order_id"`
}

func parseToolArguments(rawArguments string) (toolArguments, error) {
	var args toolArguments

	decoder := json.NewDecoder(
		strings.NewReader(rawArguments),
	)

	// Reject fields that aren't part of our schema.
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&args); err != nil {
		return args, fmt.Errorf(
			"invalid tool arguments: %w",
			err,
		)
	}

	if args.OrderID == 0 {
		return args, fmt.Errorf(
			"order_id must be greater than zero",
		)
	}

	return args, nil
}

func (r *Registry) Execute(
	ctx context.Context,
	name string,
	rawArguments string,
) ToolCallResult {

	args, err := parseToolArguments(rawArguments)

	if err != nil {
		return ToolCallResult{
			Name:  name,
			Error: err,
		}
	}

	switch name {

	case "get_order_status":
		result, err := r.OrderTool.GetOrderStatus(
			ctx,
			args.OrderID,
		)

		return ToolCallResult{
			Name:   name,
			Result: result,
			Error:  err,
		}

	case "get_payment_status":
		result, err := r.PaymentTool.GetPaymentStatus(
			ctx,
			args.OrderID,
		)

		return ToolCallResult{
			Name:   name,
			Result: result,
			Error:  err,
		}

	case "get_delivery_status":
		result, err := r.DeliveryTool.GetDeliveryStatus(
			ctx,
			args.OrderID,
		)

		return ToolCallResult{
			Name:   name,
			Result: result,
			Error:  err,
		}

	case "get_full_order_summary":
		result, err := r.SummaryTool.GetFullOrderSummary(
			ctx,
			args.OrderID,
		)

		return ToolCallResult{
			Name:   name,
			Result: result,
			Error:  err,
		}

	default:
		return ToolCallResult{
			Name: name,
			Error: fmt.Errorf(
				"unknown tool: %s",
				name,
			),
		}
	}
}
