package helpers

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/constant"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/params"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/req_context"
	smodel "gitlab.mobvista.com/ADN/structs/model"
	"gitlab.mobvista.com/mae/go-kit/decimal"
)

var (
	appsflyerReg, _ = regexp.Compile("{af_aes_{\\S*?}}")
	matReg, _       = regexp.Compile("{mat_aes_{\\S*?}}")

	appsflyerEncrypt = NewAESECBEncrypt([]byte("wMOs8JrslMwnbK44"), 16)

	matPrivateKey = "15d09e2d839ae599992bcfd48a20e3d0"[:16]
	matKey        = "b4aa405de27c5c8dd0fea55d482ddcca"
	matEncrypt    = NewAESCBCEncrypt([]byte(matKey), []byte(matPrivateKey))
)

func GenMarketUrl(loadReq *params.LoadReqData, packageName string) string {
	marketUrl := ""
	adStacking, ifFind := extractor.GetADSTACKING()
	if !ifFind {
		return marketUrl
	}
	if loadReq.Platform == constant.Android {
		marketUrl = adStacking.Android
	} else if loadReq.Platform == constant.IOS {
		marketUrl = adStacking.IOS
	}
	if len(marketUrl) <= 0 {
		return ""
	}
	if len(packageName) == 0 {
		return ""
	}
	marketUrl = strings.Replace(marketUrl, "{package_name}", packageName, -1)
	return marketUrl
}

func upperMd5(str string) string {
	if len(str) <= 0 {
		return ""
	}
	str = strings.ToUpper(str)
	str = Md5(str)
	str = strings.ToUpper(str)
	return str
}

func enSha1(str string) string {
	if len(str) <= 0 {
		return ""
	}
	return Sha1(str)
}

func lowerMd5(str string) string {
	if len(str) <= 0 {
		return ""
	}
	str = strings.ToLower(str)
	str = Md5(str)
	return str
}

func lowerUrlEncode(str string) string {
	return url.QueryEscape(strings.ToLower(str))
}

func subString(str string, begin, length int) string {
	rs := []rune(str)
	llen := len(rs)
	if begin < 0 {
		begin = 0
	}
	if begin >= llen {
		begin = llen
	}
	end := begin + length
	if end > llen {
		end = llen
	}
	return string(rs[begin:end])
}

func advSubid(appID int64, chnID int) string {
	str := strconv.Itoa(chnID) + "_" + strconv.FormatInt(appID, 10)
	str = Md5(str)
	str = subString(str, 0, 16)
	return str
}

func replaceCDN2Https(oriUrl string, publisherId string) string {
	oriUrl = strings.Replace(oriUrl, "http://cdn-adn.rayjump.com", "https://cdn-adn-https.rayjump.com", -1)
	oriUrl = strings.Replace(oriUrl, "http://res.rayjump.com", "https://res-https.rayjump.com", -1)
	oriUrl = strings.Replace(oriUrl, "http://d11kdtiohse1a9.cloudfront.net", "https://res-https.rayjump.com", -1)
	// ABTest
	if abTestCDNDomain, ok := extractor.GetHBCDNDomainABTestPubs(publisherId); ok && len(abTestCDNDomain) > 0 {
		oriUrl = strings.Replace(oriUrl, "cdn-adn-https.rayjump.com", abTestCDNDomain, -1)
	}
	oriUrl = strings.Replace(oriUrl, "cdn-adn-https.rayjump.com", "cdn-adn-ws.rayjump.com", -1)
	return oriUrl
}

func FormatLinks(ad *params.Ad, httpReq int, publisherId string) {
	if ad == nil || httpReq != 2 {
		return
	}
	if len(ad.IconURL) > 0 {
		ad.IconURL = replaceCDN2Https(ad.IconURL, publisherId)
	}
	if len(ad.ImageURL) > 0 {
		ad.ImageURL = replaceCDN2Https(ad.ImageURL, publisherId)
	}
	if len(ad.VideoURL) > 0 {
		ad.VideoURL = replaceCDN2Https(ad.VideoURL, publisherId)
	}
	if len(ad.ExtImg) > 0 {
		ad.ExtImg = replaceCDN2Https(ad.ExtImg, publisherId)
	}
}

func UrlReplace3(str string) string {
	str = strings.Replace(str, "=", "%3D", -1)
	str = strings.Replace(str, "+", "%2B", -1)
	str = strings.Replace(str, "/", "%2F", -1)
	return str
}

func UrlReplace1(str string) string {
	return strings.Replace(str, "=", "%3D", -1)
}

func GetJumpUrl(loadReq *params.LoadReqData, campaign *smodel.CampaignInfo, queryQ string, creative *params.ParamCreative, ad *params.Ad) string {
	// if loadReq.Platform == 0 {
	// 	loadReq.Platform = constant.Android
	// }
	jumpUrl := ""
	switch loadReq.Extra10 {
	case constant.SDK_TO_MARKET:
		jumpUrl = GenMarketUrl(loadReq, campaign.PackageName)
	case constant.CLIENT_SEND_DEVID_PING_SERVER:
		jumpUrl = createDirectUrl(campaign.DirectUrl, loadReq, campaign, creative)
	case constant.CLIENT_DO_ALL_PING_SERVER:
		jumpUrl = createDirectUrl(campaign.DirectUrl, loadReq, campaign, creative)
	default:
		jumpUrl = createTrackUrl(loadReq, campaign, "", queryQ, creative, ad)
	}
	return jumpUrl
}

func createTrackUrl(loadReq *params.LoadReqData, campaign *smodel.CampaignInfo, trackUrl string, queryQ string, creative *params.ParamCreative, ad *params.Ad) string {
	if len(trackUrl) == 0 {
		if len(campaign.TrackingUrl) == 0 {
			return ""
		}
		trackUrl = campaign.TrackingUrl

		// android https
		if loadReq.Platform == constant.Android && (loadReq.HttpReq == 2) && len(campaign.TrackingUrlHttps) > 0 {
			trackUrl = campaign.TrackingUrlHttps
		}
	}

	if campaign.IsSSPlatform() {
		trackUrl = createDirectUrl(trackUrl, loadReq, campaign, creative)
		return trackUrl
	}

	// append url
	if campaign.Network != 0 {
		confs, ifFind := extractor.GetTRACK_URL_CONFIG_NEW()
		if ifFind {
			conf, ok := confs[campaign.Network]
			if ok {
				urlAppend := ""
				if loadReq.Platform == constant.Android {
					urlAppend = conf.Android
				} else if loadReq.Platform == constant.IOS {
					urlAppend = conf.IOS
				}
				// 针对小程序做处理
				// if network == mvconst.NETWORK_SMALL_ROUTINE {
				// 	// 小程序相关替换规则
				// 	trackUrl = renderSmallRoutineUrl(trackUrl, params, campaign.CampaignId, urlAppend)
				// } else {
				trackUrl = trackUrl + urlAppend
				// }
			}
		}
	}

	// if loadReq.ExtdeleteDevid == constant.DELETE_DIVICEID_TRUE || loadReq.ExtdeleteDevid == constant.DELETE_DIVICEID_BUT_NOT_IMPRESSION {
	// 	loadReq.Gaid = ""
	// 	loadReq.Idfa = ""
	// 	loadReq.AndroidId = ""
	// }

	// gaid
	trackUrl = strings.Replace(trackUrl, "{gaid}", loadReq.Gaid, -1)
	// idfa
	trackUrl = strings.Replace(trackUrl, "{idfa}", loadReq.Idfa, -1)
	// devid
	devId := ""
	if loadReq.AppDevinfoEncrypt == constant.DEVINFO_ENCRYPT_DONT {
		devId = loadReq.AndroidId
	}
	trackUrl = strings.Replace(trackUrl, "{devId}", devId, -1)
	// imei
	trackUrl = strings.Replace(trackUrl, "{imei}", loadReq.Imei, -1)
	// oaid
	trackUrl = strings.Replace(trackUrl, "{oaid}", url.QueryEscape(loadReq.Oaid), -1)
	// mac
	trackUrl = strings.Replace(trackUrl, "{mac}", loadReq.Mac, -1)
	// packageName
	trackUrl = strings.Replace(trackUrl, "{package_name}", loadReq.ExtfinalPackageName, -1)
	trackUrl = strings.Replace(trackUrl, "{network}", loadReq.NetworkTypeName, -1)
	// subId
	subId := loadReq.AppId
	subIdStr := strconv.FormatInt(subId, 10)
	trackUrl = strings.Replace(trackUrl, "{subId}", subIdStr, -1)
	// clickId
	if len(queryQ) <= 0 {
		queryQ = GenParamQ(loadReq, campaign, ad)
	}
	queryQ = url.QueryEscape(queryQ)
	trackUrl = strings.Replace(trackUrl, "{clickId}", queryQ, -1)
	// 素材独享逻辑 TODO
	crat := strconv.Itoa(loadReq.CreativeAdType)
	advCrid := strconv.Itoa(loadReq.AdvCreativeID)
	trackUrl = strings.Replace(trackUrl, "{creativeId}", advCrid, -1)
	trackUrl = strings.Replace(trackUrl, "{adType}", crat, -1)
	trackUrl = strings.Replace(trackUrl, "{creativeName}", loadReq.CreativeName, -1)

	// 系统是SA，替换域名
	// system := extractor.GetSYSTEM()
	// if system == mvconst.SERVER_SYSTEM_SA {
	// 	domain := extractor.GetDOMAIN()
	// 	domainTrack := extractor.GetDOMAIN_TRACK()
	// 	trackUrl = strings.Replace(trackUrl, "http://net.rayjump.com", "http://"+domain, -1)
	// 	trackUrl = strings.Replace(trackUrl, "http://tknet.rayjump.com", "http://"+domainTrack, -1)
	// }
	// 替换agent/click的https
	if loadReq.HttpReq == 2 {
		trackUrl = strings.Replace(trackUrl, "http://net.rayjump.com", "https://net.rayjump.com", -1)
		trackUrl = strings.Replace(trackUrl, "http://tknet.rayjump.com", "https://tknet.rayjump.com", -1)
	}

	// 3s CN domain
	trackUrl = change3SCNDomain(trackUrl, loadReq.CountryCode)
	// mtgId
	mtgIdStr := "mtg" + strconv.FormatInt(loadReq.ExtMtgId, 10)
	trackUrl = strings.Replace(trackUrl, "{mtgId}", mtgIdStr, -1)
	// dsp domain
	// changeDspDomain

	return trackUrl
}

func change3SCNDomain(trackUrl string, countryCode string) string {
	conf, ifFind := extractor.Get3S_CHINA_DOMAIN()
	if !ifFind {
		return trackUrl
	}
	countrys := conf.Countrys
	if InStrArray(countryCode, countrys) {
		return trackUrl
	}
	domains := conf.Domains
	cnDomain := conf.CNDomain
	if len(domains) <= 0 || len(cnDomain) <= 0 {
		return trackUrl
	}
	for _, v := range domains {
		if strings.Contains(trackUrl, v) {
			// if params.CNDomainTest == 2 {
			cnDomain = conf.CNLineDo
			// }
			trackUrl = strings.Replace(trackUrl, v, cnDomain, -1)
			trackUrl = trackUrl + "&mb_trackingcn=1"
			return trackUrl
		}
	}
	return trackUrl
}

func createDirectUrl(dUrl string, loadReq *params.LoadReqData, campaign *smodel.CampaignInfo, creative *params.ParamCreative) string {
	// if len(campaign.DirectUrl) == 0 {
	// 	return ""
	// }
	// dUrl := campaign.DirectUrl
	// if loadReq.ExtdeleteDevid == constant.DELETE_DIVICEID_TRUE || loadReq.ExtdeleteDevid == constant.DELETE_DIVICEID_BUT_NOT_IMPRESSION {
	// 	loadReq.Gaid = ""
	// 	loadReq.Idfa = ""
	// 	loadReq.AndroidId = ""
	// }

	gaid := loadReq.Gaid
	idfa := loadReq.Idfa
	devId := ""
	if loadReq.AppDevinfoEncrypt == constant.DEVINFO_ENCRYPT_DONT {
		devId = loadReq.AndroidId
	}
	idfaMd5 := upperMd5(idfa)
	idfaSha1 := enSha1(idfa)
	gaidMd5 := upperMd5(gaid)
	gaidSha1 := enSha1(gaid)
	devIdMd5 := upperMd5(devId)
	devIdSha1 := enSha1(devId)
	imei := loadReq.Imei
	imeiMd5 := lowerMd5(imei)
	imeiSha1 := enSha1(imei)
	mac := loadReq.Mac
	macMd5 := upperMd5(mac)
	macSha1 := enSha1(mac)
	gaidDevId := gaid
	if len(gaid) <= 0 {
		gaidDevId = devId
	}
	ip := loadReq.ClientIp
	countryCode := strings.ToLower(loadReq.CountryCode)
	ua := loadReq.UserAgent
	uaOsPlatform := lowerUrlEncode(loadReq.Os)
	uaOsVersion := lowerUrlEncode(loadReq.OsVersion)
	uaDeviceModel := lowerUrlEncode(loadReq.Model)
	uaInfo := req_context.GetInstance().UaParser.Parse(ua)
	uaOs := url.QueryEscape(uaInfo.Os.Family + " " + uaInfo.Os.ToVersionString())
	ua = url.QueryEscape(ua)
	timestamp := time.Now().Unix()
	microtime := time.Now().UnixNano() / 1e6
	cityString := url.QueryEscape(loadReq.CityName)
	mgid := "mtg" + loadReq.BidId
	subId := loadReq.AppId
	// todo
	// subId := getSubId(params)
	chnId := campaign.ChnID
	advSubId := advSubid(subId, chnId)
	mbSubId := "mob" + advSubId
	cMbSubId := "c_" + advSubId
	// todo
	packageName := url.QueryEscape(loadReq.ExtfinalPackageName)
	packageMbSubId := packageName
	if len(packageMbSubId) <= 0 {
		packageMbSubId = url.QueryEscape(mbSubId)
	}
	// 大写的gaid或idfa
	gaidOrIdfa := ""
	if loadReq.Platform == constant.Android {
		gaidOrIdfa = gaid
	} else if loadReq.Platform == constant.IOS {
		gaidOrIdfa = idfa
	}
	gaidOrIdfa = strings.ToUpper(gaidOrIdfa)
	// 大写的countrycode
	upperCountryCode := strings.ToUpper(loadReq.CountryCode)
	// 素材独享逻辑 TODO
	crat := strconv.Itoa(loadReq.CreativeAdType)
	// advCreativeid
	advCrid := strconv.Itoa(loadReq.AdvCreativeID)
	priceIn, err := formatFloat64(loadReq.LocalPriceIn, 6)
	if err != nil {
		req_context.GetInstance().MLogs.Runtime.Warnf("Campaign OriPrice float64 to string error: %s", err.Error())
	}

	dUrl = strings.Replace(dUrl, "{mgid}", mgid, -1)
	dUrl = strings.Replace(dUrl, "{mbSubId}", mbSubId, -1)
	dUrl = strings.Replace(dUrl, "{ip}", ip, -1)
	dUrl = strings.Replace(dUrl, "{package_name}", packageName, -1)
	dUrl = strings.Replace(dUrl, "{gaid_devId}", gaidDevId, -1)
	dUrl = strings.Replace(dUrl, "{ua}", ua, -1)
	dUrl = strings.Replace(dUrl, "{microtime}", strconv.FormatInt(microtime, 10), -1)
	dUrl = strings.Replace(dUrl, "{countryCode}", countryCode, -1)
	dUrl = strings.Replace(dUrl, "{package_mbSubId}", packageMbSubId, -1)
	dUrl = strings.Replace(dUrl, "{c_mbSubId}", cMbSubId, -1)
	dUrl = strings.Replace(dUrl, "{idfa}", idfa, -1)
	dUrl = strings.Replace(dUrl, "{idfaMd5}", idfaMd5, -1)
	dUrl = strings.Replace(dUrl, "{idfaSha1}", idfaSha1, -1)
	dUrl = strings.Replace(dUrl, "{gaid}", gaid, -1)
	dUrl = strings.Replace(dUrl, "{gaidMd5}", gaidMd5, -1)
	dUrl = strings.Replace(dUrl, "{gaidSha1}", gaidSha1, -1)
	dUrl = strings.Replace(dUrl, "{devId}", devId, -1)
	dUrl = strings.Replace(dUrl, "{devIdMd5}", devIdMd5, -1)
	dUrl = strings.Replace(dUrl, "{devIdSha1}", devIdSha1, -1)
	dUrl = strings.Replace(dUrl, "{imei}", imei, -1)
	dUrl = strings.Replace(dUrl, "{oaid}", url.QueryEscape(loadReq.Oaid), -1)
	dUrl = strings.Replace(dUrl, "{imeiMd5}", imeiMd5, -1)
	dUrl = strings.Replace(dUrl, "{imeiSha1}", imeiSha1, -1)
	dUrl = strings.Replace(dUrl, "{mac}", mac, -1)
	dUrl = strings.Replace(dUrl, "{macMd5}", macMd5, -1)
	dUrl = strings.Replace(dUrl, "{macSha1}", macSha1, -1)
	dUrl = strings.Replace(dUrl, "{uaDevice}", uaDeviceModel, -1)
	dUrl = strings.Replace(dUrl, "{uaOsPlatform}", uaOsPlatform, -1)
	dUrl = strings.Replace(dUrl, "{uaOsVersion}", uaOsVersion, -1)
	dUrl = strings.Replace(dUrl, "{timestamp}", strconv.FormatInt(timestamp, 10), -1)
	dUrl = strings.Replace(dUrl, "{city}", cityString, -1)
	dUrl = strings.Replace(dUrl, "{uaOs}", uaOs, -1)
	dUrl = strings.Replace(dUrl, "{gaid_idfa}", gaidOrIdfa, -1)
	dUrl = strings.Replace(dUrl, "{upperCountryCode}", upperCountryCode, -1)
	dUrl = strings.Replace(dUrl, "{creativeId}", advCrid, -1)
	dUrl = strings.Replace(dUrl, "{adType}", crat, -1)
	dUrl = strings.Replace(dUrl, "{creativeName}", loadReq.CreativeName, -1)
	dUrl = strings.Replace(dUrl, "{price}", priceIn, -1)
	dUrl = strings.Replace(dUrl, "{network}", loadReq.NetworkTypeName, -1)
	dUrl = strings.Replace(dUrl, "{lang}", url.QueryEscape(loadReq.Language), -1)
	// mtgId
	mtgIdStr := "mtg" + strconv.FormatInt(loadReq.ExtMtgId, 10)
	dUrl = strings.Replace(dUrl, "{mtgId}", mtgIdStr, -1)

	// 顶级宏替换
	dUrl = renderAppsflyerAESEncrypt(dUrl)
	dUrl = renderMatAESEncrypt(dUrl)
	//dUrl = removeAdjustCallback(dUrl, loadReq, campaign)
	return dUrl
}

// hardcode encrypt
func renderAppsflyerAESEncrypt(durl string) string {
	if !strings.Contains(durl, "{af_aes_{") {
		return durl
	}

	count := strings.Count(durl, "{af_aes_{")
	for i := 0; i < count; i++ {
		matchStr := appsflyerReg.FindString(durl)
		if matchStr == "" {
			continue
		}

		repStr := strings.Replace(matchStr, "{af_aes_{", "", 1)
		if len(repStr) > 2 {
			repStr = repStr[:len(repStr)-2]
		}

		if repStr == "" {
			continue
		}

		encrypt, err := appsflyerEncrypt.Encrypt(repStr)
		if err != nil {
			req_context.GetInstance().MLogs.Runtime.Warnf("renderAppsflyerAESEncrypt raw:%v Encrypt:%v err:%v", matchStr, repStr, encrypt)
		}
		durl = strings.Replace(durl, matchStr, encrypt, 1)
	}
	return durl
}

// hardcode encrypt
func renderMatAESEncrypt(durl string) string {
	if !strings.Contains(durl, "{mat_aes_{") {
		return durl
	}

	count := strings.Count(durl, "{mat_aes_{")
	for i := 0; i < count; i++ {
		matchStr := matReg.FindString(durl)
		if matchStr == "" {
			continue
		}

		repStr := strings.Replace(matchStr, "{mat_aes_{", "", 1)
		if len(repStr) > 2 {
			repStr = repStr[:len(repStr)-2]
		}

		if repStr == "" {
			continue
		}

		encrypt, err := matEncrypt.Encrypt(repStr)
		if err != nil {
			req_context.GetInstance().MLogs.Runtime.Warnf("renderMatAESEncrypt raw:%v Encrypt:%v err:%v", matchStr, repStr, encrypt)
		}
		durl = strings.Replace(durl, matchStr, encrypt, 1)
	}
	return durl
}

func formatFloat64(f float64, frac int) (string, error) {
	dec1 := decimal.NewMDecimal()
	err := dec1.FromFloat64(f)
	if err != nil {
		return "", err
	}
	dec2, err := dec1.Round(frac)
	if err != nil {
		return "", err
	}
	return string(dec2.ToString()), nil
}
