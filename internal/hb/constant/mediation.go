package constant

const (
	MaxMediation        = 8 // ref wiki[http://confluence.mobvista.com/pages/viewpage.action?pageId=5364205]
	IronsourceMediation = 7
	TapjoyMediation     = 13
	FyberMediation      = 5
	MopubMediation      = 3
	ToponMediation      = 10
)

const (
	MaxMediationIdStr        = "8" // ref wiki[http://confluence.mobvista.com/pages/viewpage.action?pageId=5364205]
	IronsourceMediationIdStr = "7"
	TapjoyMediationIdStr     = "13"
	FyberMediationIdStr      = "5"
	MopubMediationIdStr      = "3"
	ToponMediationIdStr      = "10"
)

const (
	MaxMediationName        = "APPLOVIN_MEDIATION"
	IronsourceMediationName = "IRONSOURCE"
	TapjoyMediationName     = "TAPDAQ"
	FyberMediationName      = "FYBER"
	MopubMediationName      = "MOPUB"
	ToponMediationName      = "TOPON"
)

func GetMediationChnId(mediationName string) int {
	switch mediationName {
	case MaxMediationName:
		return MaxMediation
	case IronsourceMediationName:
		return IronsourceMediation
	case TapjoyMediationName:
		return TapjoyMediation
	case FyberMediationName:
		return FyberMediation
	case MopubMediationName:
		return MopubMediation
	case ToponMediationName:
		return ToponMediation
	default:
		return 0
	}
}
