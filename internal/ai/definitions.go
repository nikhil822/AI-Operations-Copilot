package ai

func ToolDefinitions() []ToolDefinition {
	orderIDParameter := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"order_id": map[string]interface{}{
				"type":        "integer",
				"description": "The numeric order ID",
			},
		},
		"required": []string{"order_id"},
	}

	return []ToolDefinition{
		{
			Type: "function",
			Function: FunctionDefinition{
				Name:        "get_order_status",
				Description: "Get the current status and basic information for an order. Use this when the user asks about the order itself.",
				Parameters:  orderIDParameter,
			},
		},
		{
			Type: "function",
			Function: FunctionDefinition{
				Name:        "get_payment_status",
				Description: "Get the authoritative payment status, amount, and payment date for an order. Use this when the user asks about payment.",
				Parameters:  orderIDParameter,
			},
		},
		{
			Type: "function",
			Function: FunctionDefinition{
				Name:        "get_delivery_status",
				Description: "Get the authoritative delivery status, scheduled date, and delivery date for an order. Use this when the user asks about delivery.",
				Parameters:  orderIDParameter,
			},
		},
		{
			Type: "function",
			Function: FunctionDefinition{
				Name:        "get_full_order_summary",
				Description: "Get a complete order summary including order, customer, vehicle, payment, and delivery information. Use this when the user asks for a full status summary.",
				Parameters:  orderIDParameter,
			},
		},
	}
}
