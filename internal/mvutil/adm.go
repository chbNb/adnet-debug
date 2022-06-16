package mvutil

import "strings"

const MinVastLength = 15 // 长度过小的vast不合法， 15是经验值

func IsAdmVast(adm string) bool {
	if len(adm) < MinVastLength {
		return false
	}
	// 如果不包含VAST/xml, 说明是banner
	admPrefix := strings.ToLower(adm[:MinVastLength])
	return strings.Contains(admPrefix, "vast") || strings.Contains(admPrefix, "xml")
}
