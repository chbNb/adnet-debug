package mvutil

import (
	"testing"

	smodel "gitlab.mobvista.com/ADN/structs/model"

	. "github.com/smartystreets/goconvey/convey"
)

// func TestGetBtV3(t *testing.T) {
// 	var btv3_subInfo = map[int64]SubInfo{
// 		1: SubInfo{10, "pkg", map[int64]DspSubInfo{
// 			1: DspSubInfo{20, "dsp-pkg"},
// 		}},
// 	}
//
// 	var btv3_subInfoe = map[int64]SubInfoe{
// 		1: SubInfoe{10, "pkg", []int32{1, 2, 3}},
// 	}
//
// 	Convey("Given GetBtV3 input type assertion", t, func() {
// 		res := GetBtV3("some string")
// 		Convey("Then GetBtV3 type assertion failed", func() {
// 			So(res, ShouldBeNil)
// 		})
//
// 		res = GetBtV3(btv3_subInfo)
// 		Convey("Then GetBtV3 type assertion success and return btv3", func() {
// 			So(res, ShouldResemble, btv3_subInfo)
// 		})
//
// 		res = GetBtV3(btv3_subInfoe)
// 		Convey("Then GetBtV3 type assertion success and return btv3 with []DspSubIds", func() {
// 			So(res, ShouldResemble, map[int64]SubInfo{
// 				1: SubInfo{10, "pkg", map[int64]DspSubInfo{}},
// 			})
// 		})
// 	})
// }

func TestGetCityCode(t *testing.T) {
	var cityCode = map[string]int64{"BJ": 1, "SH": 2}

	Convey("Given GetCityCode type assertion", t, func() {
		res := GetCityCode("some string")
		Convey("Then GetCityCode type assertion failed", func() {
			So(res, ShouldBeNil)
		})

		res = GetCityCode(cityCode)
		Convey("Then GetCityCode type assertion success and return cityCode", func() {
			So(res, ShouldResemble, cityCode)
		})
	})
}

func TestGetFCA(t *testing.T) {
	request := &RequestParams{AppInfo: &smodel.AppInfo{}}
	camIds := []int64{}
	newDefaultFca := 0
	Convey("获取VBA campaign的fca(FrequencyCap)", t, func() {
		campaign := &smodel.CampaignInfo{}
		Convey("fca从campaign的VBA配置中(CampaignInfo.ConfigVBA.FrequencyCap)获取", func() {
			campaign.ConfigVBA = &smodel.ConfigVBA{UseVBA: 1, FrequencyCap: 2, Status: 1}
			fca := GetFCA(request, campaign, camIds, newDefaultFca)
			So(fca, ShouldEqual, campaign.ConfigVBA.FrequencyCap)
		})
		Convey("当campaign的VBA配置中fca小于等于0，则返回fca等于1", func() {
			campaign.ConfigVBA = &smodel.ConfigVBA{UseVBA: 1, FrequencyCap: 0, Status: 1}
			fca := GetFCA(request, campaign, camIds, newDefaultFca)
			So(fca, ShouldEqual, 1)
			campaign.ConfigVBA.FrequencyCap = -1
			fca = GetFCA(request, campaign, camIds, newDefaultFca)
			So(fca, ShouldEqual, 1)
		})
	})

	Convey("获取非VBA campaign的fca(FrequencyCap)", t, func() {
		campaign := &smodel.CampaignInfo{}
		frequencyCap := int32(3)
		campaign.FrequencyCap = frequencyCap
		Convey("fca从campaign的配置中(CampaignInfo.FrequencyCap)获取,要求配置大于0", func() {
			fca := GetFCA(request, campaign, camIds, newDefaultFca)
			So(fca, ShouldEqual, 3)
		})
	})

	Convey("获取非VBA campaign的fca(FrequencyCap),且campaign的配置(CampaignInfo.FrequencyCap)小于等于0", t, func() {
		campaign := &smodel.CampaignInfo{}
		request.AppInfo.App.FrequencyCap = 6
		Convey("fca配置从app的配置中(AppInfo.App.FrequencyCap)获取", func() {
			fca := GetFCA(request, campaign, camIds, newDefaultFca)
			So(fca, ShouldEqual, 6)
		})
	})

	Convey(`获取非VBA campaign的fca(FrequencyCap),且campaign的配置(CampaignInfo.FrequencyCap)小于等于0
	且app的配置(AppInfo.App.FrequencyCap)小于0`, t, func() {
		campaign := &smodel.CampaignInfo{}
		request.AppInfo.App.FrequencyCap = -1
		Convey("fca返回默认值5", func() {
			fca := GetFCA(request, campaign, camIds, newDefaultFca)
			So(fca, ShouldEqual, 5)
		})
	})
}

func TestGetFCB(t *testing.T) {
	Convey("fcb 原本用于控制展示cap,目前用于标记素材来源", t, func() {
		Convey("素材Desc的source为ADN,则返回1", func() {
			fcb := GetFCB(int32(1))
			So(fcb, ShouldEqual, 1)
		})
		Convey("素材Desc的source为非ADN,则返回2", func() {
			fcb := GetFCB(int32(2))
			So(fcb, ShouldEqual, 2)
			fcb = GetFCB(int32(3))
			So(fcb, ShouldEqual, 2)
		})
	})
}

func TestIsVTA(t *testing.T) {
	Convey("lost VbaConnecting is not VTA", t, func() {
		isVTA := IsVTA(&smodel.CampaignInfo{})
		So(isVTA, ShouldBeFalse)
	})

	Convey("lost VbaTracking is not VTA", t, func() {
		var connecting int32 = 1
		var trackingURL = ""

		isVTA := IsVTA(&smodel.CampaignInfo{
			VbaConnecting:   connecting,
			VbaTrackingLink: trackingURL,
		})
		So(isVTA, ShouldBeFalse)
	})

	Convey("is VTA", t, func() {
		var connecting int32 = 1
		var trackingURL = "http://www.mobvista.com"

		isVTA := IsVTA(&smodel.CampaignInfo{
			VbaConnecting:   connecting,
			VbaTrackingLink: trackingURL,
		})
		So(isVTA, ShouldBeTrue)
	})
}

func TestIsVBA(t *testing.T) {
	Convey("VTA 单子不是 VBA", t, func() {
		var connecting int32 = 1
		var trackingURL = "http://www.mobvista.com"

		vbaCampaign := IsVBA(&smodel.CampaignInfo{
			VbaConnecting:   connecting,
			VbaTrackingLink: trackingURL,
		})
		So(vbaCampaign, ShouldBeFalse)
	})

	Convey("ConfigVBA 不存在不是 VBA", t, func() {
		vbaCampaign := IsVBA(&smodel.CampaignInfo{})
		So(vbaCampaign, ShouldBeFalse)
	})

	Convey("ConfigVBA 存在", t, func() {
		vbaConfig := smodel.ConfigVBA{
			UseVBA:       1,
			FrequencyCap: 10,
			Status:       1,
		}

		Convey("UseVBA 不为 1，不是 VBA 单子", func() {
			vbaConfig.UseVBA = 2
			campaign2 := smodel.CampaignInfo{
				ConfigVBA: &vbaConfig,
			}
			vbaCampaign2 := IsVBA(&campaign2)
			So(vbaCampaign2, ShouldBeFalse)
		})
		Convey("Status 不为 1，不是 VBA 单子", func() {
			vbaConfig.Status = 2
			campaign := smodel.CampaignInfo{
				ConfigVBA: &vbaConfig,
			}
			vbaCampaign := IsVBA(&campaign)
			So(vbaCampaign, ShouldBeFalse)
		})
		Convey("Status == 1 && UseVBA == 1，是 VBA 单子", func() {
			campaign := smodel.CampaignInfo{
				ConfigVBA: &vbaConfig,
			}
			vbaCampaign := IsVBA(&campaign)
			So(vbaCampaign, ShouldBeTrue)
		})
	})
}

func TestIsBrandOffer(t *testing.T) {
	Convey("Tag 不存在不是 Brand Offer", t, func() {
		res := IsBrandOffer(&smodel.CampaignInfo{})
		So(res, ShouldBeFalse)
	})

	Convey("Tag != 4 不是 Brand Offer", t, func() {
		tagVal := int32(1)
		res := IsBrandOffer(&smodel.CampaignInfo{
			Tag: tagVal,
		})
		So(res, ShouldBeFalse)
	})

	Convey("Tag == 4 是 Brand Offer", t, func() {
		tagVal := int32(4)
		res := IsBrandOffer(&smodel.CampaignInfo{
			Tag: tagVal,
		})
		So(res, ShouldBeTrue)
	})
}

func TestIsCityOffer(t *testing.T) {
	Convey("CityCode 不存在不是 City Offer", t, func() {
		res := IsCityOffer(&smodel.CampaignInfo{})
		So(res, ShouldBeFalse)
	})

	Convey("CityCode 不为 map 不是 City Offer", t, func() {
		res := IsCityOffer(&smodel.CampaignInfo{
			CityCodeV2: map[string][]int32{},
		})
		So(res, ShouldBeFalse)
	})

	Convey("CityCode 存在且为 map 是 City Offer", t, func() {
		res := IsCityOffer(&smodel.CampaignInfo{
			CityCodeV2: map[string][]int32{"BJ": []int32{1}, "SH": []int32{2}},
		})
		So(res, ShouldBeTrue)
	})
}

func TestGetVTALink(t *testing.T) {
	Convey("VbaTrackingLink 不存在返回 vta link 为空", t, func() {
		res := GetVTALink(&smodel.CampaignInfo{})
		So(res, ShouldEqual, "")
	})

	Convey("VbaTrackingLink 存在返回 VbaTrackingLink 的值", t, func() {
		vbaLink := "http://www.mobvista.com"
		res := GetVTALink(&smodel.CampaignInfo{
			VbaTrackingLink: vbaLink,
		})
		So(res, ShouldEqual, vbaLink)
	})
}

func TestGetVTATag(t *testing.T) {
	trackingURL := "http://www.mobvista.com"

	Convey("VTA 单子返回 1", t, func() {
		connecting := int32(1)
		res := GetVTATag(&smodel.CampaignInfo{
			VbaConnecting:   connecting,
			VbaTrackingLink: trackingURL,
		})
		So(res, ShouldEqual, 1)
	})

	Convey("VBA 单子返回 2", t, func() {
		res := GetVTATag(&smodel.CampaignInfo{
			ConfigVBA: &smodel.ConfigVBA{
				UseVBA:       1,
				FrequencyCap: 10,
				Status:       1,
			},
		})
		So(res, ShouldEqual, 2)
	})

	Convey("vbaLink 单子返回 3", t, func() {
		res := GetVTATag(&smodel.CampaignInfo{
			VbaTrackingLink: trackingURL,
		})
		So(res, ShouldEqual, 3)
	})

	Convey("其他单子返回 0", t, func() {
		res := GetVTATag(&smodel.CampaignInfo{})
		So(res, ShouldEqual, 0)
	})
}
