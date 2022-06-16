package link

import (
	"strconv"

	"gitlab.mobvista.com/ADN/adnet/internal/hb/constant"
)

var tempMidway map[string]string
var tempMidwayAndAdn map[string]string

func init() {
	tempMidway = make(map[string]string, 3)
	tempMidway["notice"] = "{sh}://{do}/click?k={k}&mp={mp}"
	tempMidway["impression"] = "{sh}://{do}/impression?k={k}&mp={mp}"
	tempMidway["event"] = "{sh}://{do}/ad/log/play?k={k}&mp={mp}&type={t}&key="
	tempMidway["pv"] = "{sh}://{do}/onlyImpression?k=%s&mp=%s"

	tempMidwayAndAdn = make(map[string]string, 3)
	tempMidwayAndAdn["notice"] = "{sh}://{do}/click?k={k}&z={z}&q={q}&r={r}&al={al}&csp={csp}&c={c}&notice=1"
	tempMidwayAndAdn["impression"] = "{sh}://{do}/impression?k={k}&z={z}&q={q}&x=0&r={r}&al={al}&csp={csp}&c={c}"
	tempMidwayAndAdn["event"] = "{sh}://{do}/trackv2?z={z}&q={q}&type={t}&r={r}&c={c}&csp={csp}&key="
	tempMidwayAndAdn["pv"] = "{sh}://{do}/onlyImpression?k=%s&p=%s&csp=%s&c={c}"
}

func temps(logType int) map[string]string {
	if logType == constant.LogMidway {
		return tempMidway
	} else {
		return tempMidwayAndAdn
	}
}

func template(category int, logType string) string {
	temp := temps(category)
	tempUrl, ok := temp[logType]
	if !ok {
		return ""
	}
	return tempUrl
}

func ImpTemplate(category int) string {
	return template(category, constant.Imp)
}

func NoticeTemplate(category int) string {
	return template(category, constant.Notice)
}

func PlayTemplate(category int, rate int) string {
	return EventTemplate(category, constant.Play) + "&rate=" + strconv.Itoa(rate)
}

func EventTemplate(category int, key string) string {
	return template(category, constant.Event) + key
}

func PvTemplate(category int) string {
	return template(category, constant.Pv)
}
