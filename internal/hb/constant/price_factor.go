package constant

const (
	PriceFactorTag_A  = 1
	PriceFactorTag_B  = 2
	PriceFactorTag_BB = 3
)

const (
	PriceFactorHit_TRUE  = 1
	PriceFactorHit_FALSE = 2
)
const (
	FreqControlTimeWindowModeByHour = 1
	FreqControlTimeWindowModeByDate = 2
)

const PriceFactor_AerospikePrefix = "dd_"

const PriceFactor_DefaultValue = 1
const PriceFactor_MAXValue = 10
const PriceFactor_MINValue = 0

const (
	SALT_PRICE_FACTOR_GROUP_ABTEST   = "salt_price_factor_group_abtest"
	SALT_PRICE_FACTOR_RATE_ABTEST    = "salt_price_factor_rate_abtest"
	SALT_PRICE_FACTOR_SUBRATE_ABTEST = "salt_price_factor_subrate_abtest"
)
