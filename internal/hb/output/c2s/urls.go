package c2s

func GetWinWithMacors() string {
	return "{sd}/win?td={td}"
}

func GetLossWithMacros() string {
	return "{sd}/loss?td={td}"
}

func GetWinWithMacorsWithoutTd() string {
	return "{sd}/win?token={token}"
}

func GetLossWithMacrosWithoutTd() string {
	return "{sd}/loss?token={token}"
}
