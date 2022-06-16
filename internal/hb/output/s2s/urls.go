package s2s

import "fmt"

// GetWinUrl func
func GetWinUrl(sd, td string) *string {
	url := fmt.Sprintf("%s/win?td=%s", sd, td)
	return &url
}

// GetLossUrl func
func GetLossUrl(sd, td string) *string {
	url := fmt.Sprintf("%s/loss?td=%s", sd, td)
	return &url
}
