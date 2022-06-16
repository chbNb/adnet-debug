package mvutil

import (
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
)

// func GetBackendId(dspId int64) int {
// 	switch dspId {
// 	case mvconst.FakeToutiao:
// 		return mvconst.TouTiao
// 	case mvconst.FakeGDT:
// 		return mvconst.Gdt
// 	case mvconst.FakeBaidu:
// 		return mvconst.Baidu
// 	default:
// 		return 0
// 	}
// }

// IsThirdDsp 是第三方DSP， 不包括 as, mvdsp, mas
func IsThirdDsp(dspId int64) bool {
	return !(dspId == mvconst.FakeAdserverDsp || dspId == mvconst.MVDSP || dspId == mvconst.MAS || dspId == mvconst.MVDSP_Retarget || dspId <= 0)
}

// 是第三方DSP， 不包括 as, mvdsp, mas。但包括rtdsp
func IsThirdDspWithRtdsp(dspId int64) bool {
	return !(dspId == mvconst.FakeAdserverDsp || dspId == mvconst.MVDSP || dspId == mvconst.MAS || dspId <= 0)
}
