package output

import (
	supply_mvutil "gitlab.mobvista.com/ADN/chasm/module/supply/mvutil"
	"testing"

	. "github.com/bouk/monkey"
	. "github.com/smartystreets/goconvey/convey"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	smodel "gitlab.mobvista.com/ADN/structs/model"
)

//func TestRandJumpType(t *testing.T) {
//	Convey("test RandJumpType", t, func() {
//		var res string
//		var r mvutil.RequestParams
//		var campaign smodel.CampaignInfo
//		var params mvutil.Params
//
//		Convey("RandJumpType return 0", func() {
//			r = mvutil.RequestParams{}
//			campaign = smodel.CampaignInfo{}
//			params = mvutil.Params{}
//			res = RandJumpType(&r, &campaign, &params)
//			So(res, ShouldEqual, "0")
//		})
//
//		Convey("RequestType == 10", func() {
//			r = mvutil.RequestParams{UnitInfo: &smodel.UnitInfo{}, AppInfo: &smodel.AppInfo{}, PublisherInfo: &smodel.PublisherInfo{}}
//			campaign = smodel.CampaignInfo{}
//			params.RequestType = 10
//			guard := Patch(extractor.GetCAN_CLICK_MODE_SIX_PUBLISHER, func() ([]int64, bool) {
//				return []int64{int64(12345)}, true
//			})
//			defer guard.Unpatch()
//			guard = Patch(extractor.GetJUMP_TYPE_CONFIG_ADV, func() (map[string]map[string]map[string]int32, bool) {
//				return map[string]map[string]map[string]int32{}, true
//			})
//			defer guard.Unpatch()
//			guardTp := Patch(extractor.GetJUMP_TYPE_CONFIG_THIRD_PARTY, func() (map[string]map[string]map[string]int32, bool) {
//				return map[string]map[string]map[string]int32{}, true
//			})
//			defer guardTp.Unpatch()
//			guardTp = Patch(extractor.GetOnlineEmptyDeviceNoServerJump, func() (mvutil.OnlineEmptyDeviceNoServerJump, bool) {
//				return mvutil.OnlineEmptyDeviceNoServerJump{}, true
//			})
//			defer guardTp.Unpatch()
//			guardTp = Patch(extractor.GetOnlineEmptyDeviceIPUAABTest, func() (mvutil.OnlineEmptyDeviceIPUA, bool) {
//				return mvutil.OnlineEmptyDeviceIPUA{}, true
//			})
//			defer guardTp.Unpatch()
//			res = RandJumpType(&r, &campaign, &params)
//			So(res, ShouldEqual, "0")
//		})
//
//		Convey("RequestType != 10", func() {
//			guard := Patch(getJumpTypeArr, func(r *mvutil.RequestParams, campaign *smodel.CampaignInfo, params *mvutil.Params) map[string]int32 {
//				return map[string]int32{"test_1": int32(101)}
//			})
//			defer guard.Unpatch()
//
//			guard = Patch(mvutil.RandByRate, func(rateMap map[int]int) int {
//				return 100
//			})
//			defer guard.Unpatch()
//
//			r = mvutil.RequestParams{UnitInfo: &smodel.UnitInfo{}, AppInfo: &smodel.AppInfo{}, PublisherInfo: &smodel.PublisherInfo{}}
//			campaign = smodel.CampaignInfo{}
//			params.RequestType = 1
//			res = RandJumpType(&r, &campaign, &params)
//			So(res, ShouldEqual, "0")
//		})
//	})
//}

func TestGetJumpTypeArr(t *testing.T) {
	Convey("getJumpTypeArr 空", t, func() {
		r := mvutil.RequestParams{UnitInfo: &smodel.UnitInfo{}, AppInfo: &smodel.AppInfo{}, PublisherInfo: &smodel.PublisherInfo{}}
		campaign := smodel.CampaignInfo{}
		params := mvutil.Params{}
		res := getJumpTypeArr(&r, &campaign, &params)
		So(res, ShouldResemble, map[string]int32{mvconst.JUMP_TYPE_NORMAL: int32(100)})
	})

	Convey("getJumpTypeArr android", t, func() {
		r := mvutil.RequestParams{
			Param: mvutil.Params{Platform: 1},
			UnitInfo: &smodel.UnitInfo{
				UnitId: 100,
				JumptypeConfig: map[string]int32{
					"0": 300,
					"1": 301,
					"2": 302,
					"3": 303,
					"4": 304,
					"5": 305,
				},
			},
		}
		campaign := smodel.CampaignInfo{}
		params := mvutil.Params{}
		res := getJumpTypeArr(&r, &campaign, &params)
		So(res, ShouldResemble, map[string]int32{mvconst.JUMP_TYPE_NORMAL: int32(300)})
	})

	Convey("getJumpTypeArr ios", t, func() {
		r := mvutil.RequestParams{
			Param: mvutil.Params{Platform: 2},
			UnitInfo: &smodel.UnitInfo{
				UnitId: 100,
				JumptypeConfig: map[string]int32{
					"0": 400,
					"1": 401,
					"2": 402,
					"3": 403,
					"4": 404,
					"5": 405,
				},
			},
		}
		campaign := smodel.CampaignInfo{}
		params := mvutil.Params{}
		res := getJumpTypeArr(&r, &campaign, &params)
		So(res, ShouldResemble, map[string]int32{mvconst.JUMP_TYPE_NORMAL: int32(400)})
	})
}

func TestGetJumpTypeConf(t *testing.T) {
	campaignJConfig := map[string]int32{
		"campaign_jump_type_key_1": 1,
		"campaign_jump_type_key_2": 2,
		"campaign_jump_type_key_3": 3,
	}

	unitJConfig := map[string]int32{
		"unit_jump_type_key_1": 1,
		"unit_jump_type_key_2": 2,
		"unit_jump_type_key_3": 3,
	}

	appJConfig := map[string]int32{
		"app_jump_type_key_1": 1,
		"app_jump_type_key_2": 2,
		"app_jump_type_key_3": 3,
	}

	pubJConfig := map[string]int32{
		"app_jump_type_key_1": 1,
		"app_jump_type_key_2": 2,
		"app_jump_type_key_3": 3,
	}

	advConfig := map[string]int32{
		"adv_jump_type_key_1": 1,
		"adv_jump_type_key_2": 2,
		"adv_jump_type_key_3": 3,
	}

	Convey("空参数", t, func() {
		r := mvutil.RequestParams{UnitInfo: &smodel.UnitInfo{}, AppInfo: &smodel.AppInfo{}, PublisherInfo: &smodel.PublisherInfo{}}
		campaign := smodel.CampaignInfo{}
		guard := Patch(extractor.GetJUMP_TYPE_CONFIG_ADV, func() (map[string]map[string]map[string]int32, bool) {
			return map[string]map[string]map[string]int32{}, true
		})
		defer guard.Unpatch()
		guardTp := Patch(extractor.GetJUMP_TYPE_CONFIG_THIRD_PARTY, func() (map[string]map[string]map[string]int32, bool) {
			return map[string]map[string]map[string]int32{}, true
		})
		defer guardTp.Unpatch()
		res := getJumpTypeConf(&r, &campaign)
		var resultConf map[string]int32
		So(res, ShouldResemble, resultConf)
	})

	Convey("优先 campaign 配置", t, func() {
		campaign := smodel.CampaignInfo{
			JumpTypeConfig: campaignJConfig,
		}

		r := mvutil.RequestParams{
			UnitInfo:      &smodel.UnitInfo{JumptypeConfig: unitJConfig},
			AppInfo:       &smodel.AppInfo{JumptypeConfig: appJConfig},
			PublisherInfo: &smodel.PublisherInfo{JumptypeConfig: pubJConfig},
			Param:         mvutil.Params{Platform: 1},
		}

		res := getJumpTypeConf(&r, &campaign)
		So(res, ShouldResemble, campaignJConfig)
	})

	Convey("unit 维度", t, func() {
		r := mvutil.RequestParams{
			UnitInfo:      &smodel.UnitInfo{JumptypeConfig: unitJConfig},
			AppInfo:       &smodel.AppInfo{JumptypeConfig: appJConfig},
			PublisherInfo: &smodel.PublisherInfo{JumptypeConfig: pubJConfig},
			Param:         mvutil.Params{Platform: 1},
		}
		campaign := smodel.CampaignInfo{}
		res := getJumpTypeConf(&r, &campaign)
		So(res, ShouldResemble, unitJConfig)
	})

	Convey("app 维度", t, func() {
		r := mvutil.RequestParams{
			AppInfo:       &smodel.AppInfo{JumptypeConfig: appJConfig},
			PublisherInfo: &smodel.PublisherInfo{JumptypeConfig: pubJConfig},
			Param:         mvutil.Params{Platform: 1},
			UnitInfo:      &smodel.UnitInfo{},
		}
		campaign := smodel.CampaignInfo{}
		res := getJumpTypeConf(&r, &campaign)
		So(res, ShouldResemble, appJConfig)
	})

	Convey("Publisher 维度", t, func() {
		r := mvutil.RequestParams{
			PublisherInfo: &smodel.PublisherInfo{JumptypeConfig: pubJConfig},
			Param:         mvutil.Params{Platform: 1},
			UnitInfo:      &smodel.UnitInfo{},
			AppInfo:       &smodel.AppInfo{},
		}
		campaign := smodel.CampaignInfo{}
		res := getJumpTypeConf(&r, &campaign)
		So(res, ShouldResemble, pubJConfig)
	})

	Convey("adv 维度", t, func() {
		r := mvutil.RequestParams{
			Param:         mvutil.Params{Platform: 1},
			UnitInfo:      &smodel.UnitInfo{},
			AppInfo:       &smodel.AppInfo{},
			PublisherInfo: &smodel.PublisherInfo{},
		}
		adv := int32(123)
		campaign := smodel.CampaignInfo{
			AdvertiserId: adv,
		}
		guard := Patch(extractor.GetJUMP_TYPE_CONFIG_ADV, func() (map[string]map[string]map[string]int32, bool) {
			return map[string]map[string]map[string]int32{"1": {"123": advConfig}}, true
		})
		defer guard.Unpatch()
		guardTp := Patch(extractor.GetJUMP_TYPE_CONFIG_THIRD_PARTY, func() (map[string]map[string]map[string]int32, bool) {
			return map[string]map[string]map[string]int32{}, true
		})
		defer guardTp.Unpatch()
		res := getJumpTypeConf(&r, &campaign)
		So(res, ShouldResemble, advConfig)
	})

	Convey("兜底", t, func() {
		guard := Patch(extractor.GetJUMP_TYPE_CONFIG, func() (map[string]int32, bool) {
			return map[string]int32{"test": int32(1)}, true
		})
		defer guard.Unpatch()

		campaign := smodel.CampaignInfo{}
		// iOSDefaultJConf := map[string]int32{"ios_default_j_config_1": int32(42)}
		// androidDefaultJConf := map[string]int32{"android_default_j_config_1": int32(42)}

		// var mvConfigObj *extractor.MVConfigObj

		// mvConfigObj.JumpTypeConfig = androidDefaultJConf
		// mvConfigObj.JumpTypeIOS = iOSDefaultJConf

		// Convey("Android Platform", func() {
		// 	r := mvutil.RequestParams{Param: mvutil.Params{Platform: 1}}
		// 	res := getJumpTypeConf(r, campaign)
		// 	So(res, ShouldResemble, androidDefaultJConf)
		// })

		// Convey("iOS Platform", func() {
		// 	r := mvutil.RequestParams{Param: mvutil.Params{Platform: 2}}
		// 	res := getJumpTypeConf(r, campaign)
		// 	So(res, ShouldResemble, iOSDefaultJConf)
		// })

		Convey("Other Platform", func() {
			r := mvutil.RequestParams{Param: mvutil.Params{Platform: 4}, AppInfo: &smodel.AppInfo{}, UnitInfo: &smodel.UnitInfo{}, PublisherInfo: &smodel.PublisherInfo{}}
			guard := Patch(extractor.GetJUMP_TYPE_CONFIG_ADV, func() (map[string]map[string]map[string]int32, bool) {
				return map[string]map[string]map[string]int32{}, true
			})
			defer guard.Unpatch()
			guardTp := Patch(extractor.GetJUMP_TYPE_CONFIG_THIRD_PARTY, func() (map[string]map[string]map[string]int32, bool) {
				return map[string]map[string]map[string]int32{}, true
			})
			defer guardTp.Unpatch()
			res := getJumpTypeConf(&r, &campaign)
			So(res, ShouldBeNil)
		})
	})
}

func TestGetJumpTypeAndroid(t *testing.T) {
	r := mvutil.RequestParams{UnitInfo: &smodel.UnitInfo{}, AppInfo: &smodel.AppInfo{}, Param: mvutil.Params{}, PublisherInfo: &smodel.PublisherInfo{}}
	campaign := smodel.CampaignInfo{}
	params := mvutil.Params{}
	campaignPackageName := "test_campaign_package_name_1"
	directPackageName := "test_direct_package_name_1"
	paramsPackageName := "test_params_package_name_1"
	guard := Patch(extractor.GetJUMP_TYPE_CONFIG_ADV, func() (map[string]map[string]map[string]int32, bool) {
		return map[string]map[string]map[string]int32{}, true
	})
	defer guard.Unpatch()
	guardTp := Patch(extractor.GetJUMP_TYPE_CONFIG_THIRD_PARTY, func() (map[string]map[string]map[string]int32, bool) {
		return map[string]map[string]map[string]int32{}, true
	})
	defer guardTp.Unpatch()

	Convey("empty request", t, func() {
		res := getJumpTypeAndroid(&r, &campaign, &params)
		So(res, ShouldResemble, map[string]int32{"0": 0})
	})

	Convey("canClickMode1", t, func() {
		r = mvutil.RequestParams{
			UnitInfo: &smodel.UnitInfo{
				UnitId: 100,
				JumptypeConfig: map[string]int32{
					"1": 301,
					"2": 302,
					"3": 303,
					"4": 304,
					"5": 305,
				},
			},
		}
		jType := int32(1)
		campaign = smodel.CampaignInfo{
			PackageName:       campaignPackageName,
			DirectPackageName: directPackageName,
			JumpType:          jType,
		}
		params = mvutil.Params{
			Platform:    2,
			GAID:        "test_gaid",
			IDFA:        "test_idfa",
			LinkType:    1,
			PackageName: paramsPackageName,
		}

		// todo: 依赖 extrcator 包，暂跳过
		// res := getJumpTypeAndroid(r, campaign, params)
		// expRes := map[string]int32{
		// 	"1": 301,
		// 	"4": 304,
		// }
		// So(res, ShouldResemble, expRes)
	})
}

func TestGetJumpTypeIOS(t *testing.T) {
	r := mvutil.RequestParams{UnitInfo: &smodel.UnitInfo{}, AppInfo: &smodel.AppInfo{}, PublisherInfo: &smodel.PublisherInfo{}}
	campaign := smodel.CampaignInfo{}
	params := mvutil.Params{}
	guard := Patch(extractor.GetJUMP_TYPE_CONFIG_ADV, func() (map[string]map[string]map[string]int32, bool) {
		return map[string]map[string]map[string]int32{}, true
	})
	defer guard.Unpatch()
	guardTp := Patch(extractor.GetJUMP_TYPE_CONFIG_THIRD_PARTY, func() (map[string]map[string]map[string]int32, bool) {
		return map[string]map[string]map[string]int32{}, true
	})
	defer guardTp.Unpatch()
	Convey("empty request", t, func() {
		res := getJumpTypeIOS(&r, &campaign, &params)
		So(res, ShouldResemble, map[string]int32{"0": 0})
	})
}

func TestIsClickMode0(t *testing.T) {
	campaign := smodel.CampaignInfo{}
	params := mvutil.Params{}
	Convey("空 campaign 返回 true", t, func() {
		res := isClickMode0(&campaign, &params)
		So(res, ShouldBeTrue)
	})

	Convey("空 campaign 返回 true", t, func() {
		res := isClickMode0(&campaign, &params)
		So(res, ShouldBeTrue)
	})

	Convey("campaign packageName 空字符串返回 true", t, func() {
		packageName := ""
		res := isClickMode0(&smodel.CampaignInfo{PackageName: packageName}, &params)
		So(res, ShouldBeTrue)
	})

	// Convey("campaign link type 非 1、2 返回 true", t, func() {
	// 	packageName := "test_package_name"
	// 	res := isClickMode0(smodel.CampaignInfo{PackageName: &packageName}, mvutil.Params{LinkType: 3})
	// 	So(res, ShouldBeTrue)
	// })

	Convey("其他返回 false", t, func() {
		packageName := "test_package_name"
		res := isClickMode0(&smodel.CampaignInfo{PackageName: packageName}, &mvutil.Params{LinkType: 1})
		So(res, ShouldBeFalse)
	})
}

func TestGetJumpTypeVal(t *testing.T) {
	conf := map[string]int32{
		"key1":  1,
		"key2":  2,
		"key42": 42,
	}
	jumpType := ""

	Convey("不存在的 key 返回 0", t, func() {
		res := getJumpTypeVal(conf, jumpType)
		So(res, ShouldEqual, int32(0))

		res = getJumpTypeVal(conf, "error")
		So(res, ShouldEqual, int32(0))
	})

	Convey("存在的 key 返回对应值", t, func() {
		jumpType = "key1"
		res := getJumpTypeVal(conf, jumpType)
		So(res, ShouldEqual, int32(1))

		jumpType = "key42"
		res = getJumpTypeVal(conf, jumpType)
		So(res, ShouldEqual, int32(42))
	})
}

func TestCanClickMode1(t *testing.T) {
	campaign := smodel.CampaignInfo{}
	params := mvutil.Params{}

	Convey("返回 false", t, func() {
		res := canClickMode1(&campaign, &params)
		So(res, ShouldBeFalse)
	})

	Convey("返回 false", t, func() {
		params = mvutil.Params{
			Platform: 1,
		}
		res := canClickMode1(&campaign, &params)
		So(res, ShouldBeFalse)
	})

	jType := int32(3)
	Convey("android 返回 true", t, func() {
		packageName := "package"
		campaign = smodel.CampaignInfo{
			DirectPackageName: packageName,
			JumpType:          jType,
		}
		params = mvutil.Params{
			Platform: 1,
			GAID:     "test_gaid",
		}
		res := canClickMode1(&campaign, &params)
		So(res, ShouldBeTrue)
	})

	Convey("ios 返回 true", t, func() {
		packageName := "package"

		campaign = smodel.CampaignInfo{
			DirectPackageName: packageName,
			JumpType:          jType,
		}
		params = mvutil.Params{
			Platform: 2,
			GAID:     "test_gaid",
			IDFA:     "test_idfa",
		}
		res := canClickMode1(&campaign, &params)
		So(res, ShouldBeTrue)
	})
}

func TestCanVTA(t *testing.T) {
	Convey("jumptype 空返回 false", t, func() {
		campaign := smodel.CampaignInfo{}
		res := canVTA(&campaign)
		So(res, ShouldBeFalse)
	})

	Convey("derectPackaName 非空 & jumptype = 1 返回 true", t, func() {
		packageName := "test_direct_package"
		jumpType := int32(1)
		campaign := smodel.CampaignInfo{
			JumpType:          jumpType,
			DirectPackageName: packageName,
		}
		res := canVTA(&campaign)
		So(res, ShouldBeTrue)
	})

	Convey("jumptype = 3、7 返回 true", t, func() {
		jumpType := int32(3)
		campaign := smodel.CampaignInfo{
			JumpType: jumpType,
		}
		res := canVTA(&campaign)
		So(res, ShouldBeTrue)

		jumpType = int32(7)
		campaign = smodel.CampaignInfo{
			JumpType: jumpType,
		}
		res = canVTA(&campaign)
		So(res, ShouldBeTrue)
	})

	Convey("jumptype 为其他返回 false", t, func() {
		jumpType := int32(42)
		campaign := smodel.CampaignInfo{
			JumpType: jumpType,
		}
		res := canVTA(&campaign)
		So(res, ShouldBeFalse)
	})
}

func TestCompareSDKVersionForJumpType(t *testing.T) {
	Convey("test compareSDKVersionForJumpType", t, func() {
		var params mvutil.Params
		var jumpType string
		var res bool

		Convey("not get jump type version", func() {
			guard := Patch(extractor.GetJUMPTYPE_SDKVERSION, func() (map[string]map[string]string, bool) {
				return nil, false
			})
			defer guard.Unpatch()

			res = compareSDKVersionForJumpType(&params, jumpType)
			So(res, ShouldBeFalse)
		})

		Convey("get jump type version, jump type not in config", func() {
			guard := Patch(extractor.GetJUMPTYPE_SDKVERSION, func() (map[string]map[string]string, bool) {
				return map[string]map[string]string{
					"j_type_1": map[string]string{
						"j_type_1_sub_key1": "j_type_sub_1_val1",
						"j_type_1_sub_key2": "j_type_sub_1_val2",
					},
				}, true
			})
			defer guard.Unpatch()

			jumpType = "j_type_2"
			res = compareSDKVersionForJumpType(&params, jumpType)
			So(res, ShouldBeFalse)
		})

		Convey("get jump type version, jump type in config, sdktype is empty", func() {
			guard := Patch(extractor.GetJUMPTYPE_SDKVERSION, func() (map[string]map[string]string, bool) {
				return map[string]map[string]string{
					"j_type_1": map[string]string{
						"j_type_1_sub_key1": "j_type_sub_1_val1",
						"j_type_1_sub_key2": "j_type_sub_1_val2",
					},
				}, true
			})
			defer guard.Unpatch()

			guard = Patch(supply_mvutil.RenderSDKVersion, func(sdkversion string) supply_mvutil.SDKVersionItem {
				return supply_mvutil.SDKVersionItem{
					SDKType:        "",
					SDKNumber:      "s_num_1",
					SDKVersionCode: int32(42),
				}
			})
			defer guard.Unpatch()

			jumpType = "j_type_1"
			res = compareSDKVersionForJumpType(&params, jumpType)
			So(res, ShouldBeFalse)
		})

		Convey("get jump type version, jump type in config, sdktype is not in config", func() {
			guard := Patch(extractor.GetJUMPTYPE_SDKVERSION, func() (map[string]map[string]string, bool) {
				return map[string]map[string]string{
					"j_type_1": map[string]string{
						"j_type_1_sub_key1": "j_type_sub_1_val1",
						"j_type_1_sub_key2": "j_type_sub_1_val2",
					},
				}, true
			})
			defer guard.Unpatch()

			guard = Patch(supply_mvutil.RenderSDKVersion, func(sdkversion string) supply_mvutil.SDKVersionItem {
				return supply_mvutil.SDKVersionItem{
					SDKType:        "j_type_1_sub_key100",
					SDKNumber:      "s_num_1",
					SDKVersionCode: int32(42),
				}
			})
			defer guard.Unpatch()

			jumpType = "j_type_1"
			res = compareSDKVersionForJumpType(&params, jumpType)
			So(res, ShouldBeFalse)
		})

		Convey("get jump type version, jump type in config, sdktype is in config, < sdk version", func() {
			guard := Patch(extractor.GetJUMPTYPE_SDKVERSION, func() (map[string]map[string]string, bool) {
				return map[string]map[string]string{
					"j_type_1": map[string]string{
						"j_type_1_sub_key1": "j_type_sub_1_val1",
						"j_type_1_sub_key2": "j_type_sub_1_val2",
					},
				}, true
			})
			defer guard.Unpatch()

			guard = Patch(supply_mvutil.RenderSDKVersion, func(sdkversion string) supply_mvutil.SDKVersionItem {
				return supply_mvutil.SDKVersionItem{
					SDKType:        "j_type_1_sub_key1",
					SDKNumber:      "s_num_1",
					SDKVersionCode: int32(42),
				}
			})
			defer guard.Unpatch()

			guard = Patch(mvutil.GetVersionCode, func(version string) int32 {
				return int32(100)
			})
			defer guard.Unpatch()

			jumpType = "j_type_1"
			res = compareSDKVersionForJumpType(&params, jumpType)
			So(res, ShouldBeFalse)
		})

		Convey("get jump type version, jump type in config, sdktype is in config, >= sdk version", func() {
			guard := Patch(extractor.GetJUMPTYPE_SDKVERSION, func() (map[string]map[string]string, bool) {
				return map[string]map[string]string{
					"j_type_1": map[string]string{
						"j_type_1_sub_key1": "j_type_sub_1_val1",
						"j_type_1_sub_key2": "j_type_sub_1_val2",
					},
				}, true
			})
			defer guard.Unpatch()

			guard = Patch(supply_mvutil.RenderSDKVersion, func(sdkversion string) supply_mvutil.SDKVersionItem {
				return supply_mvutil.SDKVersionItem{
					SDKType:        "j_type_1_sub_key1",
					SDKNumber:      "s_num_1",
					SDKVersionCode: int32(42),
				}
			})
			defer guard.Unpatch()

			guard = Patch(mvutil.GetVersionCode, func(version string) int32 {
				return int32(1)
			})
			defer guard.Unpatch()

			jumpType = "j_type_1"
			res = compareSDKVersionForJumpType(&params, jumpType)
			So(res, ShouldBeTrue)
		})
	})
}

func TestGetJumpTypeConfV2(t *testing.T) {
	campaignJConfig := map[string]int32{
		"campaign_jump_type_key_1": 1,
		"campaign_jump_type_key_2": 2,
		"campaign_jump_type_key_3": 3,
	}

	unitJConfig := map[string]int32{
		"unit_jump_type_key_1": 1,
		"unit_jump_type_key_2": 2,
		"unit_jump_type_key_3": 3,
	}

	appJConfig := map[string]int32{
		"app_jump_type_key_1": 1,
		"app_jump_type_key_2": 2,
		"app_jump_type_key_3": 3,
	}

	pubJConfig := map[string]int32{
		"app_jump_type_key_1": 1,
		"app_jump_type_key_2": 2,
		"app_jump_type_key_3": 3,
	}

	advConfig := map[string]int32{
		"adv_jump_type_key_1": 1,
		"adv_jump_type_key_2": 2,
		"adv_jump_type_key_3": 3,
	}

	Convey("空参数", t, func() {
		r := &mvutil.RequestParams{
			UnitInfo:      &smodel.UnitInfo{},
			AppInfo:       &smodel.AppInfo{},
			PublisherInfo: &smodel.PublisherInfo{}}
		campaign := &smodel.CampaignInfo{}
		guard := Patch(extractor.GetJUMP_TYPE_CONFIG_ADV, func() (map[string]map[string]map[string]int32, bool) {
			return map[string]map[string]map[string]int32{}, true
		})
		defer guard.Unpatch()
		guardTp := Patch(extractor.GetJUMP_TYPE_CONFIG_THIRD_PARTY, func() (map[string]map[string]map[string]int32, bool) {
			return map[string]map[string]map[string]int32{}, true
		})
		defer guardTp.Unpatch()
		res := getJumpTypeConfV2(*r, *campaign)
		So(res, ShouldResemble, map[string]int32{"0": int32(1)})
	})

	Convey("优先 campaign 配置", t, func() {
		campaign := smodel.CampaignInfo{
			JumpTypeConfig2: campaignJConfig,
		}

		r := mvutil.RequestParams{
			UnitInfo:      &smodel.UnitInfo{JumptypeConfigV2: unitJConfig},
			AppInfo:       &smodel.AppInfo{JumptypeConfigV2: appJConfig},
			PublisherInfo: &smodel.PublisherInfo{JumptypeConfigV2: pubJConfig},
		}

		res := getJumpTypeConfV2(r, campaign)
		So(res, ShouldResemble, campaignJConfig)
	})

	Convey("unit 维度", t, func() {
		r := mvutil.RequestParams{
			UnitInfo:      &smodel.UnitInfo{JumptypeConfigV2: unitJConfig},
			AppInfo:       &smodel.AppInfo{JumptypeConfigV2: appJConfig},
			PublisherInfo: &smodel.PublisherInfo{JumptypeConfigV2: pubJConfig},
		}
		campaign := smodel.CampaignInfo{}
		guard := Patch(extractor.GetJUMP_TYPE_CONFIG_ADV, func() (map[string]map[string]map[string]int32, bool) {
			return map[string]map[string]map[string]int32{}, true
		})
		defer guard.Unpatch()
		guardTp := Patch(extractor.GetJUMP_TYPE_CONFIG_THIRD_PARTY, func() (map[string]map[string]map[string]int32, bool) {
			return map[string]map[string]map[string]int32{}, true
		})
		defer guardTp.Unpatch()
		res := getJumpTypeConfV2(r, campaign)
		So(res, ShouldResemble, unitJConfig)
	})

	Convey("app 维度", t, func() {
		r := mvutil.RequestParams{
			AppInfo:       &smodel.AppInfo{JumptypeConfigV2: appJConfig},
			PublisherInfo: &smodel.PublisherInfo{JumptypeConfigV2: pubJConfig},
			UnitInfo:      &smodel.UnitInfo{},
		}
		campaign := smodel.CampaignInfo{}
		guard := Patch(extractor.GetJUMP_TYPE_CONFIG_ADV, func() (map[string]map[string]map[string]int32, bool) {
			return map[string]map[string]map[string]int32{}, true
		})
		defer guard.Unpatch()
		guardTp := Patch(extractor.GetJUMP_TYPE_CONFIG_THIRD_PARTY, func() (map[string]map[string]map[string]int32, bool) {
			return map[string]map[string]map[string]int32{}, true
		})
		defer guardTp.Unpatch()
		res := getJumpTypeConfV2(r, campaign)
		So(res, ShouldResemble, appJConfig)
	})

	Convey("Publisher 维度", t, func() {
		r := mvutil.RequestParams{
			PublisherInfo: &smodel.PublisherInfo{JumptypeConfigV2: pubJConfig},
			UnitInfo:      &smodel.UnitInfo{},
			AppInfo:       &smodel.AppInfo{},
		}
		campaign := smodel.CampaignInfo{}
		guard := Patch(extractor.GetJUMP_TYPE_CONFIG_ADV, func() (map[string]map[string]map[string]int32, bool) {
			return map[string]map[string]map[string]int32{}, true
		})
		defer guard.Unpatch()
		guardTp := Patch(extractor.GetJUMP_TYPE_CONFIG_THIRD_PARTY, func() (map[string]map[string]map[string]int32, bool) {
			return map[string]map[string]map[string]int32{}, true
		})
		defer guardTp.Unpatch()
		res := getJumpTypeConfV2(r, campaign)
		So(res, ShouldResemble, pubJConfig)
	})

	Convey("adv 维度", t, func() {
		r := mvutil.RequestParams{
			Param:         mvutil.Params{Platform: 1},
			PublisherInfo: &smodel.PublisherInfo{},
			UnitInfo:      &smodel.UnitInfo{},
			AppInfo:       &smodel.AppInfo{},
		}
		adv := int32(123)
		campaign := smodel.CampaignInfo{
			AdvertiserId: adv,
		}
		guard := Patch(extractor.GetJUMP_TYPE_CONFIG_ADV, func() (map[string]map[string]map[string]int32, bool) {
			return map[string]map[string]map[string]int32{"androidV2": {"123": advConfig}}, true
		})
		defer guard.Unpatch()
		guardTp := Patch(extractor.GetJUMP_TYPE_CONFIG_THIRD_PARTY, func() (map[string]map[string]map[string]int32, bool) {
			return map[string]map[string]map[string]int32{}, true
		})
		defer guardTp.Unpatch()
		res := getJumpTypeConfV2(r, campaign)
		So(res, ShouldResemble, advConfig)
	})

	Convey("兜底", t, func() {
		r := mvutil.RequestParams{
			PublisherInfo: &smodel.PublisherInfo{},
			UnitInfo:      &smodel.UnitInfo{},
			AppInfo:       &smodel.AppInfo{},
		}
		campaign := smodel.CampaignInfo{}
		guard := Patch(extractor.GetJUMP_TYPE_CONFIG_ADV, func() (map[string]map[string]map[string]int32, bool) {
			return map[string]map[string]map[string]int32{}, true
		})
		defer guard.Unpatch()
		guardTp := Patch(extractor.GetJUMP_TYPE_CONFIG_THIRD_PARTY, func() (map[string]map[string]map[string]int32, bool) {
			return map[string]map[string]map[string]int32{}, true
		})
		defer guardTp.Unpatch()
		res := getJumpTypeConfV2(r, campaign)
		So(res, ShouldResemble, map[string]int32{"0": int32(1)})
	})
}
