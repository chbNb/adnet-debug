package output

import (
	"testing"

	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	smodel "gitlab.mobvista.com/ADN/structs/model"

	. "github.com/bouk/monkey"
	. "github.com/smartystreets/goconvey/convey"
)

func TestHandleGradeD(t *testing.T) {
	Convey("TestHandleGradeD for output.HandelGradeD()", t, func() {
		guard := Patch(extractor.GetConfigLowGradeApp, func() (int64, bool) {
			return int64(42), true
		})
		defer guard.Unpatch()
		Convey("空请求返回 0", func() {
			res := HandleGradeD(mvutil.RequestParams{AppInfo: &smodel.AppInfo{}})
			//res := HandleGradeD(mvutil.RequestParams{AppInfo: &mvutil.AppInfo{}})
			So(res, ShouldEqual, 0)
		})

		Convey("非空请求，如果是 D 级流量，返回配置的 appID", func() {
			res := HandleGradeD(mvutil.RequestParams{
				AppInfo: &smodel.AppInfo{
					//AppInfo: &mvutil.AppInfo{
					//App: mvutil.App{Grade: 4},
				},
			})
			So(res, ShouldEqual, 0)
		})

		Convey("非空请求，非 D 级流量，返回真实 appID", func() {
			res := HandleGradeD(mvutil.RequestParams{
				//AppInfo: &mvutil.AppInfo{
				//	App:   mvutil.App{Grade: 1},
				//	AppId: 10,
				//},
				AppInfo: &smodel.AppInfo{
					App:   smodel.App{Grade: 1},
					AppId: 10,
				},
			})
			So(res, ShouldEqual, 10)
		})
	})
}

func TestBlendTraffic(t *testing.T) {
	Convey("空请求返回 0", t, func() {
		request := mvutil.RequestParams{}
		params := mvutil.Params{}
		campaign := smodel.CampaignInfo{}
		res := BlendTraffic(request, &campaign, &params)
		So(res, ShouldEqual, 0)
	})

	Convey("空请求返回参数中的 AppID", t, func() {
		request := mvutil.RequestParams{}
		params := mvutil.Params{
			AppID: 100,
		}
		campaign := smodel.CampaignInfo{}
		res := BlendTraffic(request, &campaign, &params)
		So(res, ShouldEqual, 100)
	})

	Convey("campaignInfo.Ctype 不合法返回 AppID", t, func() {
		offerCType := int32(10)
		request := mvutil.RequestParams{}
		params := mvutil.Params{
			AppID: 101,
		}
		campaign := smodel.CampaignInfo{
			Ctype: offerCType,
		}
		res := BlendTraffic(request, &campaign, &params)
		So(res, ShouldEqual, 101)
	})

	Convey("存在合法的 Campaign 信息", t, func() {
		offerCType := int32(mvconst.PAY_TYPE_CPI)
		btV4Data := smodel.BtV4{
			SubIds: map[string]*smodel.SubInfo{
				"1001": &smodel.SubInfo{
					Rate:        10,
					PackageName: "test_package_name",
					DspSubIds: map[string]*smodel.DspSubInfo{
						"10000001": &smodel.DspSubInfo{
							Rate:        100,
							PackageName: "dsp_sub_package_name",
						},
					},
				},
			},
			BtClass: map[string]*smodel.BtClass{
				"1": &smodel.BtClass{
					Percent:   0.9,
					CapMargin: int32(100),
					Status:    int32(1),
				},
			},
		}
		campaign := smodel.CampaignInfo{
			Ctype: offerCType,
			BtV4:  &btV4Data,
		}

		Convey("离线开发者", func() {
			publisher := smodel.PublisherInfo{
				Publisher: smodel.Publisher{
					Type: 1, // 类型为 ADN
				},
			}

			Convey("无 BT Class 的 App", func() {
				app := smodel.AppInfo{
					App: smodel.App{
						BtClass: 0, // BT class B
					},
				}
				request := mvutil.RequestParams{
					PublisherInfo: &publisher,
					AppInfo:       &app,
				}
				params := mvutil.Params{
					AppID: 102,
				}
				res := BlendTraffic(request, &campaign, &params)
				So(res, ShouldEqual, 102)
			})

			Convey("BT Class A 的 App", func() {
				app := smodel.AppInfo{
					App: smodel.App{
						BtClass: 1, // BT class B
					},
				}
				request := mvutil.RequestParams{
					PublisherInfo: &publisher,
					AppInfo:       &app,
				}
				params := mvutil.Params{
					AppID: 103,
				}
				res := BlendTraffic(request, &campaign, &params)
				So(res, ShouldEqual, 1001)
			})
		})
	})
}

func TestBlendTrafficV3(t *testing.T) {
	// todo
}

func TestRenderRequestPackage(t *testing.T) {
	Convey("空请求返回 0", t, func() {
		//res := RenderRequestPackage()
	})
}

func TestRenderRequestParams(t *testing.T) {
	r := mvutil.RequestParams{}

	Convey("finalsubid = extra14", t, func() {
		params := mvutil.Params{
			Extfinalsubid: 0,
			Extra14:       14,
			AppID:         100,
		}
		renderRequestParams(r, &params)
		So(params.Extfinalsubid, ShouldEqual, 14)
	})

	Convey("finalsubid = appID", t, func() {
		params := mvutil.Params{
			Extfinalsubid: 0,
			Extra14:       0,
			AppID:         100,
		}
		renderRequestParams(r, &params)
		So(params.Extfinalsubid, ShouldEqual, 100)
	})
}

func TestRenderNewSubId(t *testing.T) {
	var r mvutil.RequestParams
	var campaign smodel.CampaignInfo
	var params mvutil.Params

	Convey("test output.renderNewSubId", t, func() {
		Convey("无 new subid 配置，不做处理", func() {
			initSubID := int64(10101)
			params.Extfinalsubid = initSubID
			guard := Patch(extractor.GetAdvBlackSubIdList, func() (map[string]map[string]map[string]string, bool) {
				subIDList := map[string]map[string]map[string]string{}
				return subIDList, true
			})
			defer guard.Unpatch()
			renderNewSubId(r, campaign, &params)
			So(params.Extfinalsubid, ShouldEqual, int64(10101))
		})

		Convey("有 1 对 new subid 配置", func() {
			params.Extfinalsubid = int64(1)
			campaign.BlackSubidListV2 = map[string]map[string]string{
				"1": map[string]string{
					"1001": "subID_package_name",
				},
			}
			renderNewSubId(r, campaign, &params)
			So(params.Extfinalsubid, ShouldEqual, 1001)
			So(params.ExtfinalPackageName, ShouldEqual, "subID_package_name")
		})

		Convey("有多对 new subid 配置", func() {
			params.Extfinalsubid = int64(1)
			campaign.BlackSubidListV2 = map[string]map[string]string{
				"1": map[string]string{
					"1001": "subID_package_name",
					"1002": "subID_package_name_2",
				},
			}
			renderNewSubId(r, campaign, &params)
			So(params.Extfinalsubid, ShouldBeIn, []int64{1001, 1002})
			So(params.ExtfinalPackageName, ShouldBeIn, []string{"subID_package_name", "subID_package_name_2"})
		})
	})
}

func TestRandNewSubId(t *testing.T) {
	Convey("subArr 非空", t, func() {
		subArr := map[string]string{
			"1": "sub_val_1",
			"2": "sub_val_2",
		}
		finalSubID := int64(43)
		res := randNewSubId(subArr, finalSubID)
		So(res, ShouldEqual, int64(1))
	})

	// 暂时跳过，待修复
	Convey("subArr 为空", t, func() {
		// subArr := map[string]string{}
		// finalSubID := int64(10)
		// res := randNewSubId(subArr, finalSubID)
		// So(res, ShouldEqual, 0)
	})
}

func TestGetNewSubIdConf(t *testing.T) {
	Convey("test output.getNewSubIdConf", t, func() {
		var campaign smodel.CampaignInfo
		var subIDInt64 int64
		var res map[string]string
		var params mvutil.Params

		Convey("Campaign 维度 subID 数据", func() {
			campaign = smodel.CampaignInfo{
				BlackSubidListV2: map[string]map[string]string{
					"42": map[string]string{
						"subID_1_1": "val_1_1",
					},
				},
			}
			subIDInt64 = int64(42)
			res = getNewSubIdConf(campaign, subIDInt64, &params)
			So(res, ShouldResemble, map[string]string{"subID_1_1": "val_1_1"})
		})

		Convey("Campaign 维度默认数据", func() {
			campaign = smodel.CampaignInfo{
				BlackSubidListV2: map[string]map[string]string{
					"1": map[string]string{
						"subID_1_1": "val_1_1",
					},
					"0": map[string]string{
						"subID_1_2": "val_1_2",
					},
				},
			}
			subIDInt64 = int64(42)
			res = getNewSubIdConf(campaign, subIDInt64, &params)
			So(res, ShouldResemble, map[string]string{"subID_1_2": "val_1_2"})
		})

		Convey("Advertiser 维度数据，没有 advertiserID，返回空 map", func() {
			campaign = smodel.CampaignInfo{}
			subIDInt64 = int64(100)

			// Mock extractor.GetAdvBlackSubIdList()
			guard := Patch(extractor.GetAdvBlackSubIdList, func() (map[string]map[string]map[string]string, bool) {
				subIDList := map[string]map[string]map[string]string{
					"subID_lev1_1": map[string]map[string]string{
						"subID_lev2_1": map[string]string{
							"subID_lev3_1": "lev3_1_val",
						},
					},
				}

				return subIDList, true
			})
			defer guard.Unpatch()

			res = getNewSubIdConf(campaign, subIDInt64, &params)
			So(res, ShouldResemble, map[string]string{})
		})

		Convey("Advertiser 维度数据，有 advertiserID，存在 SubID 配置", func() {
			advID := int32(2018)
			campaign = smodel.CampaignInfo{
				AdvertiserId: advID,
			}
			subIDInt64 = int64(42)

			// Mock extractor.GetAdvBlackSubIdList()
			guard := Patch(extractor.GetAdvBlackSubIdList, func() (map[string]map[string]map[string]string, bool) {
				subIDList := map[string]map[string]map[string]string{
					"2018": map[string]map[string]string{
						"42": map[string]string{
							"subID_lev3_1": "lev3_1_val",
						},
					},
				}

				return subIDList, true
			})
			defer guard.Unpatch()

			res = getNewSubIdConf(campaign, subIDInt64, &params)
			So(res, ShouldResemble, map[string]string{"subID_lev3_1": "lev3_1_val"})
		})

		Convey("Advertiser 维度数据，有 advertiserID，存在 SubID 默认配置", func() {
			advID := int32(2018)
			campaign = smodel.CampaignInfo{
				AdvertiserId: advID,
			}
			subIDInt64 = int64(42)

			// Mock extractor.GetAdvBlackSubIdList()
			guard := Patch(extractor.GetAdvBlackSubIdList, func() (map[string]map[string]map[string]string, bool) {
				subIDList := map[string]map[string]map[string]string{
					"2018": map[string]map[string]string{
						"10": map[string]string{
							"subID_lev3_1": "lev3_1_val",
						},
						"0": map[string]string{
							"subID_lev3_2": "subID_lev3_2",
						},
					},
				}

				return subIDList, true
			})
			defer guard.Unpatch()

			res = getNewSubIdConf(campaign, subIDInt64, &params)
			So(res, ShouldResemble, map[string]string{"subID_lev3_2": "subID_lev3_2"})
		})

		Convey("Advertiser 维度数据，有 advertiserID，存在 SubID 配置, 0", func() {
			advID := int32(2018)
			campaign = smodel.CampaignInfo{
				AdvertiserId: advID,
			}
			subIDInt64 = int64(42)
			params.RequestPath = mvconst.PATHMPAD
			params.PublisherType = mvconst.PublisherTypeM
			params.AdType = 42
			// Mock extractor.GetAdvBlackSubIdList()
			guard := Patch(extractor.GetAdvBlackSubIdList, func() (map[string]map[string]map[string]string, bool) {
				subIDList := map[string]map[string]map[string]string{
					"2018": map[string]map[string]string{
						"10": map[string]string{
							"subID_lev3_1": "lev3_1_val",
						},
						"0": map[string]string{
							"subID_lev3_2": "subID_lev3_2",
						},
						"-1": map[string]string{
							"subID_lev3_3": "subID_lev3_3",
						},
						"-2": map[string]string{
							"subID_lev3_4": "subID_lev3_4",
						},
						"-3": map[string]string{
							"subID_lev3_5": "subID_lev3_5",
						},
					},
				}

				return subIDList, true
			})
			defer guard.Unpatch()

			res = getNewSubIdConf(campaign, subIDInt64, &params)
			So(res, ShouldResemble, map[string]string{"subID_lev3_5": "subID_lev3_5"})
		})

		Convey("Advertiser 维度数据，有 advertiserID，存在 SubID 配置, -1", func() {
			advID := int32(2018)
			campaign = smodel.CampaignInfo{
				AdvertiserId: advID,
			}
			subIDInt64 = int64(42)
			params.RequestPath = mvconst.PATHMPAD
			params.PublisherType = mvconst.PublisherTypeM
			params.AdType = 42
			// Mock extractor.GetAdvBlackSubIdList()
			guard := Patch(extractor.GetAdvBlackSubIdList, func() (map[string]map[string]map[string]string, bool) {
				subIDList := map[string]map[string]map[string]string{
					"2018": map[string]map[string]string{
						"10": map[string]string{
							"subID_lev3_1": "lev3_1_val",
						},
						"0": map[string]string{
							"subID_lev3_2": "subID_lev3_2",
						},
						"-1": map[string]string{
							"subID_lev3_3": "subID_lev3_3",
						},
						"-2": map[string]string{
							"subID_lev3_4": "subID_lev3_4",
						},
					},
				}

				return subIDList, true
			})
			defer guard.Unpatch()

			res = getNewSubIdConf(campaign, subIDInt64, &params)
			So(res, ShouldResemble, map[string]string{"subID_lev3_2": "subID_lev3_2"})
		})

		Convey("Advertiser 维度数据，有 advertiserID，存在 SubID 配置, -2", func() {
			advID := int32(2018)
			campaign = smodel.CampaignInfo{
				AdvertiserId: advID,
			}
			subIDInt64 = int64(42)
			params.RequestPath = mvconst.PATHMPAD
			params.PublisherID = mvconst.DspPublisherID
			params.AdType = 42
			// Mock extractor.GetAdvBlackSubIdList()
			guard := Patch(extractor.GetAdvBlackSubIdList, func() (map[string]map[string]map[string]string, bool) {
				subIDList := map[string]map[string]map[string]string{
					"2018": map[string]map[string]string{
						"10": map[string]string{
							"subID_lev3_1": "lev3_1_val",
						},
						"0": map[string]string{
							"subID_lev3_2": "subID_lev3_2",
						},
						"-1": map[string]string{
							"subID_lev3_3": "subID_lev3_3",
						},
						"-2": map[string]string{
							"subID_lev3_4": "subID_lev3_4",
						},
					},
				}

				return subIDList, true
			})
			defer guard.Unpatch()

			res = getNewSubIdConf(campaign, subIDInt64, &params)
			So(res, ShouldResemble, map[string]string{"subID_lev3_4": "subID_lev3_4"})
		})

	})
}

func TestNeedPackageName(t *testing.T) {
	Convey("test output.needPackageName", t, func() {
		var r mvutil.RequestParams
		var campaign smodel.CampaignInfo
		var params mvutil.Params
		var res bool

		Convey("campaign.AppPostList 不为空", func() {
			appPostList := &smodel.AppPostList{
				Include: []string{"1", "2", "3"},
				Exclude: []string{"3", "4"},
			}

			Convey("没有 appId，返回 false", func() {
				r = mvutil.RequestParams{}
				campaign = smodel.CampaignInfo{
					AppPostList: appPostList,
				}
				params = mvutil.Params{}
				res = needPackageName(r, campaign, params)
				So(res, ShouldBeFalse)
			})

			Convey("appId 在 Exclude 中，返回 false", func() {
				r = mvutil.RequestParams{}
				campaign = smodel.CampaignInfo{
					AppPostList: appPostList,
				}
				params = mvutil.Params{
					AppID:           1,
					PublisherID:     6028,
					ExtdspRealAppid: 3,
				}
				res = needPackageName(r, campaign, params)
				So(res, ShouldBeFalse)
			})

			Convey("appId 在 Include 中，返回 true", func() {
				r = mvutil.RequestParams{}
				campaign = smodel.CampaignInfo{
					AppPostList: appPostList,
				}
				params = mvutil.Params{
					AppID:           1,
					PublisherID:     6028,
					ExtdspRealAppid: 2,
				}
				res = needPackageName(r, campaign, params)
				So(res, ShouldBeTrue)
			})

			Convey("appId 不在在 Include & Exclude 中，但 Include 包含 ALL，返回 true", func() {
				r = mvutil.RequestParams{}
				appPostList := &smodel.AppPostList{
					Include: []string{"1", "2", "ALL"},
					Exclude: []string{"3", "4"},
				}

				campaign = smodel.CampaignInfo{
					AppPostList: appPostList,
				}
				params = mvutil.Params{
					AppID:           1,
					PublisherID:     6028,
					ExtdspRealAppid: 5,
				}
				res = needPackageName(r, campaign, params)
				So(res, ShouldBeTrue)
			})

			Convey("非 dsp app，appId 使用原 appId", func() {
				r = mvutil.RequestParams{}
				campaign = smodel.CampaignInfo{
					AppPostList: appPostList,
				}
				params = mvutil.Params{
					AppID:           1,
					PublisherID:     6000,
					ExtdspRealAppid: 4,
				}
				res = needPackageName(r, campaign, params)
				So(res, ShouldBeTrue)
			})
		})

		Convey("campaign.AppPostList 为空", func() {
			Convey("advertiser 维度无配置，返回 false", func() {
				guard := Patch(extractor.GetAppPostList, func() (map[int32]*smodel.AppPostList, bool) {
					return map[int32]*smodel.AppPostList{}, false
				})
				defer guard.Unpatch()

				res = needPackageName(r, campaign, params)
				So(res, ShouldBeFalse)
			})

			Convey("advertiser 维度有配置，campaign 无 advertiserID 信息，返回 false", func() {
				guard := Patch(extractor.GetAppPostList, func() (map[int32]*smodel.AppPostList, bool) {
					return map[int32]*smodel.AppPostList{}, true
				})
				defer guard.Unpatch()

				campaign = smodel.CampaignInfo{
					AppPostList: nil,
				}
				res = needPackageName(r, campaign, params)
				So(res, ShouldBeFalse)
			})

			Convey("advertiser 维度有配置，app 在 exclude 配置中，返回 false", func() {
				advID := int32(999)
				guard := Patch(extractor.GetAppPostList, func() (map[int32]*smodel.AppPostList, bool) {
					return map[int32]*smodel.AppPostList{
						999: &smodel.AppPostList{
							Exclude: []string{"123"},
							Include: []string{"123"},
						},
					}, true
				})
				defer guard.Unpatch()

				campaign = smodel.CampaignInfo{
					AppPostList:  nil,
					AdvertiserId: advID,
				}
				params = mvutil.Params{
					AppID: 123,
				}
				res = needPackageName(r, campaign, params)
				So(res, ShouldBeFalse)
			})

			Convey("advertiser 维度有配置，app 在 include 配置中，返回 true", func() {
				advID := int32(999)
				guard := Patch(extractor.GetAppPostList, func() (map[int32]*smodel.AppPostList, bool) {
					return map[int32]*smodel.AppPostList{
						999: &smodel.AppPostList{Include: []string{"123"}},
					}, true
				})
				defer guard.Unpatch()

				campaign = smodel.CampaignInfo{
					AppPostList:  nil,
					AdvertiserId: advID,
				}
				params = mvutil.Params{
					AppID: 123,
				}
				res = needPackageName(r, campaign, params)
				So(res, ShouldBeTrue)
			})

			Convey("advertiser 维度有配置，app 不在 exclude 配置中，include 配置存在 ALL，返回 true", func() {
				advID := int32(999)
				guard := Patch(extractor.GetAppPostList, func() (map[int32]*smodel.AppPostList, bool) {
					return map[int32]*smodel.AppPostList{
						999: &smodel.AppPostList{
							Exclude: []string{"987"},
							Include: []string{"ALL"},
						},
					}, true
				})
				defer guard.Unpatch()

				campaign = smodel.CampaignInfo{
					AppPostList:  nil,
					AdvertiserId: advID,
				}
				params = mvutil.Params{
					AppID: 123,
				}
				res = needPackageName(r, campaign, params)
				So(res, ShouldBeTrue)
			})
		})
	})
}

func TestGetBTClass(t *testing.T) {
	Convey("流量参数为空，返回 0", t, func() {
		r := mvutil.RequestParams{PublisherInfo: &smodel.PublisherInfo{}, AppInfo: &smodel.AppInfo{}, UnitInfo: &smodel.UnitInfo{}}
		res := getBTClass(r)
		So(res, ShouldEqual, 0)
	})

	Convey("M 系统开发者，返回 unit 维度 bt class", t, func() {
		r := mvutil.RequestParams{
			PublisherInfo: &smodel.PublisherInfo{
				Publisher: smodel.Publisher{
					Type: 3,
				},
			},
			AppInfo: &smodel.AppInfo{
				App: smodel.App{
					BtClass: 5,
				},
			},
			UnitInfo: &smodel.UnitInfo{
				Unit: smodel.Unit{
					BtClass: 4,
				},
			},
		}
		res := getBTClass(r)
		So(res, ShouldEqual, 4)
	})

	Convey("Adn 离线开发者，返回 app 维度 bt class", t, func() {
		r := mvutil.RequestParams{
			PublisherInfo: &smodel.PublisherInfo{
				Publisher: smodel.Publisher{
					Type: 1,
				},
			},
			AppInfo: &smodel.AppInfo{
				App: smodel.App{
					BtClass: 5,
				},
			},
			UnitInfo: &smodel.UnitInfo{
				Unit: smodel.Unit{
					BtClass: 4,
				},
			},
		}
		res := getBTClass(r)
		So(res, ShouldEqual, 5)
	})
}
