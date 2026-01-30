package zai

// Pricing constants for z.ai models (prices per 1M tokens in USD)
const (
	// GLM-4.7
	glm47InputPrice  = 0.6 // $0.6 / 1M tokens
	glm47OutputPrice = 2.2 // $2.2 / 1M tokens

	// GLM-4.7-FlashX
	glm47FlashXInputPrice  = 0.07 // $0.07 / 1M tokens
	glm47FlashXOutputPrice = 0.4  // $0.4 / 1M tokens

	// GLM-4.6
	glm46InputPrice  = 0.6 // $0.6 / 1M tokens
	glm46OutputPrice = 2.2 // $2.2 / 1M tokens

	// GLM-4.5
	glm45InputPrice  = 0.6 // $0.6 / 1M tokens
	glm45OutputPrice = 2.2 // $2.2 / 1M tokens

	// GLM-4.5-X
	glm45XInputPrice  = 2.2 // $2.2 / 1M tokens
	glm45XOutputPrice = 8.9 // $8.9 / 1M tokens

	// GLM-4.5-Air
	glm45AirInputPrice  = 0.2 // $0.2 / 1M tokens
	glm45AirOutputPrice = 1.1 // $1.1 / 1M tokens

	// GLM-4-32B-0414-128K
	glm32BInputPrice  = 0.1 // $0.1 / 1M tokens
	glm32BOutputPrice = 0.1 // $0.1 / 1M tokens
)

// EstimateCost calculates the cost of a request based on model and token usage
func (p *Provider) EstimateCost(model string, inputTokens, outputTokens int) float64 {
	var inputPrice, outputPrice float64

	switch model {
	case "glm-4.7":
		inputPrice, outputPrice = glm47InputPrice, glm47OutputPrice
	case "glm-4.7-flashx":
		inputPrice, outputPrice = glm47FlashXInputPrice, glm47FlashXOutputPrice
	case "glm-4.6":
		inputPrice, outputPrice = glm46InputPrice, glm46OutputPrice
	case "glm-4.5":
		inputPrice, outputPrice = glm45InputPrice, glm45OutputPrice
	case "glm-4.5-x":
		inputPrice, outputPrice = glm45XInputPrice, glm45XOutputPrice
	case "glm-4.5-air":
		inputPrice, outputPrice = glm45AirInputPrice, glm45AirOutputPrice
	case "glm-4-32b-0414-128k":
		inputPrice, outputPrice = glm32BInputPrice, glm32BOutputPrice
	default:
		// Default to glm-4.7 pricing
		inputPrice, outputPrice = glm47InputPrice, glm47OutputPrice
	}

	// Calculate cost: (input_tokens / 1M) * input_price + (output_tokens / 1M) * output_price
	inputCost := float64(inputTokens) / 1_000_000 * inputPrice
	outputCost := float64(outputTokens) / 1_000_000 * outputPrice

	return inputCost + outputCost
}
