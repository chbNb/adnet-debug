package output

import (
	"testing"

	"github.com/bouk/monkey"
	mlogger "github.com/mae-pax/logger"
	. "github.com/smartystreets/goconvey/convey"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/watcher"
	smodel "gitlab.mobvista.com/ADN/structs/model"
)

func TestGetCampaignsInfo(t *testing.T) {
	Convey("test GetCampaignsInfo", t, func() {
		var campaignIds []int64

		Convey("返回 err", func() {
			_, err := GetCampaignsInfo(campaignIds)
			So(err, ShouldNotBeNil)
		})

		Convey("返回 true", func() {
			c := mlogger.NewFromYaml("./testdata/watch_log.yaml")
			logger := c.InitLogger("time", "", true, true)
			watcher.Init(logger)
			campaignIds = append(campaignIds, int64(1234))
			guard := monkey.Patch(extractor.GetCampaignInfo, func(int64) (camInfo *smodel.CampaignInfo, ifFind bool) {
				return &smodel.CampaignInfo{CampaignId: 1234}, true
			})
			defer guard.Unpatch()
			camInfo, ifFind := GetCampaignsInfo(campaignIds)
			So(ifFind, ShouldBeNil)
			So(camInfo, ShouldNotBeNil)
		})
	})
}

// func TestCopyMobvistaData(t *testing.T) {
// 	Convey("test CopyMobvistaData", t, func() {
// 		var result MobvistaResult
// 		var data MobvistaData

// 		Convey("测试 copy", func() {
// 			result = MobvistaResult{
// 				Status: int(1),
// 				Msg:    "test_msg",
// 				Data: MobvistaData{
// 					SessionID:         "mv_session_id",
// 					ParentSessionID:   "mv_parent_session_id",
// 					AdType:            2,
// 					Template:          3,
// 					UnitSize:          "mv_unit_size",
// 					Ads:               []Ad{},
// 					HTMLURL:           "mv_html_url",
// 					EndScreenURL:      "mv_end_screen_url",
// 					OnlyImpressionURL: "mv_impression_url",
// 				},
// 				DebugInfo: nil,
// 			}
// 			data = MobvistaData{
// 				SessionID:         "test_session_id",
// 				ParentSessionID:   "test_parent_session_id",
// 				AdType:            42,
// 				Template:          43,
// 				UnitSize:          "test_unit_size",
// 				Ads:               []Ad{},
// 				HTMLURL:           "test_html_url",
// 				EndScreenURL:      "test_end_screen_url",
// 				OnlyImpressionURL: "only_impression_url",
// 			}
// 			CopyMobvistaData(&result, &data)

// 			So(result.Data.AdType, ShouldEqual, 42)
// 			So(result.Data.EndScreenURL, ShouldEqual, "test_end_screen_url")
// 			So(result.Data.HTMLURL, ShouldEqual, "test_html_url")
// 			So(result.Data.OnlyImpressionURL, ShouldEqual, "only_impression_url")
// 			So(result.Data.ParentSessionID, ShouldEqual, "test_parent_session_id")
// 			So(result.Data.SessionID, ShouldEqual, "test_session_id")
// 			So(result.Data.UnitSize, ShouldEqual, "test_unit_size")
// 		})
// 	})
// }

// func TestMerge(t *testing.T) {
// 	Convey("test Merge", t, func() {
// 		var res []Ad
// 		var adList map[int][]Ad
// 		var backends []int

// 		Convey("空参数", func() {
// 			var exp []Ad
// 			res = Merge(adList, backends)
// 			So(res, ShouldResemble, exp)
// 		})

// 		Convey("正常参数", func() {
// 			adList = map[int][]Ad{
// 				1: []Ad{},
// 				2: []Ad{
// 					Ad{
// 						CampaignID:    int64(200),
// 						OfferID:       int(300),
// 						AppName:       "app_name",
// 						AppDesc:       "app_desc",
// 						PackageName:   "package_name",
// 						IconURL:       "icon_url",
// 						ImageURL:      "image_url",
// 						ImpressionURL: "impression_url",
// 					},
// 				},
// 			}
// 			backends = []int{1, 2, 3, 4}
// 			res = Merge(adList, backends)
// 			So(res, ShouldResemble, []Ad{
// 				Ad{
// 					CampaignID:    int64(200),
// 					OfferID:       int(300),
// 					AppName:       "app_name",
// 					AppDesc:       "app_desc",
// 					PackageName:   "package_name",
// 					IconURL:       "icon_url",
// 					ImageURL:      "image_url",
// 					ImpressionURL: "impression_url",
// 				}})
// 		})
// 	})
// }

// func TestVideoAdTypeRank(t *testing.T) {
// 	Convey("test VideoAdTypeRank", t, func() {
// 		var res []Ad
// 		var adList []Ad

// 		Convey("空参数", func() {
// 			var exp []Ad
// 			res = VideoAdTypeRank(adList)
// 			So(res, ShouldEqual, exp)
// 		})

// 		Convey("正常参数", func() {
// 			adList = []Ad{
// 				Ad{
// 					CampaignID:    int64(200),
// 					OfferID:       int(300),
// 					AppName:       "app_name_300",
// 					AppDesc:       "app_desc_300",
// 					PackageName:   "package_name_300",
// 					IconURL:       "icon_url_300",
// 					ImageURL:      "image_url_300",
// 					ImpressionURL: "impression_url_300",
// 				},
// 				Ad{
// 					CampaignID:    int64(201),
// 					OfferID:       int(301),
// 					AppName:       "app_name_301",
// 					AppDesc:       "app_desc_301",
// 					PackageName:   "package_name_301",
// 					IconURL:       "icon_url_301",
// 					ImageURL:      "image_url_301",
// 					ImpressionURL: "impression_url_301",
// 					VideoURL:      "video_url_301",
// 				},
// 			}
// 			res = VideoAdTypeRank(adList)
// 			exp := []Ad{
// 				Ad{
// 					CampaignID:    int64(201),
// 					OfferID:       int(301),
// 					AppName:       "app_name_301",
// 					AppDesc:       "app_desc_301",
// 					PackageName:   "package_name_301",
// 					IconURL:       "icon_url_301",
// 					ImageURL:      "image_url_301",
// 					ImpressionURL: "impression_url_301",
// 					VideoURL:      "video_url_301",
// 				},
// 				Ad{
// 					CampaignID:    int64(200),
// 					OfferID:       int(300),
// 					AppName:       "app_name_300",
// 					AppDesc:       "app_desc_300",
// 					PackageName:   "package_name_300",
// 					IconURL:       "icon_url_300",
// 					ImageURL:      "image_url_300",
// 					ImpressionURL: "impression_url_300",
// 				},
// 			}
// 			So(res, ShouldResemble, exp)
// 		})
// 	})
// }

// func TestVideoAdTypeFilter(t *testing.T) {
// 	Convey("test VideoAdTypeFilter", t, func() {
// 		var res []Ad
// 		var adList []Ad
// 		var VideoAdType int

// 		Convey("空参数", func() {
// 			var exp []Ad
// 			res = VideoAdTypeFilter(adList, VideoAdType)
// 			So(res, ShouldResemble, exp)
// 		})

// 		Convey("VideoAdType != 2 && VideoAdType != 3", func() {
// 			adList = []Ad{
// 				Ad{
// 					CampaignID:    int64(200),
// 					OfferID:       int(300),
// 					AppName:       "app_name_300",
// 					AppDesc:       "app_desc_300",
// 					PackageName:   "package_name_300",
// 					IconURL:       "icon_url_300",
// 					ImageURL:      "image_url_300",
// 					ImpressionURL: "impression_url_300",
// 				},
// 			}
// 			VideoAdType = 42
// 			res = VideoAdTypeFilter(adList, VideoAdType)
// 			So(res, ShouldResemble, []Ad{
// 				Ad{
// 					CampaignID:    int64(200),
// 					OfferID:       int(300),
// 					AppName:       "app_name_300",
// 					AppDesc:       "app_desc_300",
// 					PackageName:   "package_name_300",
// 					IconURL:       "icon_url_300",
// 					ImageURL:      "image_url_300",
// 					ImpressionURL: "impression_url_300",
// 				},
// 			})
// 		})

// 		Convey("VideoAdType == 2", func() {
// 			adList = []Ad{
// 				Ad{
// 					CampaignID:    int64(200),
// 					OfferID:       int(300),
// 					AppName:       "app_name_300",
// 					AppDesc:       "app_desc_300",
// 					PackageName:   "package_name_300",
// 					IconURL:       "icon_url_300",
// 					ImageURL:      "image_url_300",
// 					ImpressionURL: "impression_url_300",
// 					VideoURL:      "test_video_url",
// 				},
// 				Ad{
// 					CampaignID:    int64(201),
// 					OfferID:       int(301),
// 					AppName:       "app_name_301",
// 					AppDesc:       "app_desc_301",
// 					PackageName:   "package_name_301",
// 					IconURL:       "icon_url_301",
// 					ImageURL:      "image_url_301",
// 					ImpressionURL: "impression_url_301",
// 					VideoURL:      "",
// 				},
// 			}
// 			VideoAdType = 2
// 			res = VideoAdTypeFilter(adList, VideoAdType)
// 			So(res, ShouldResemble, []Ad{
// 				Ad{
// 					CampaignID:    int64(201),
// 					OfferID:       int(301),
// 					AppName:       "app_name_301",
// 					AppDesc:       "app_desc_301",
// 					PackageName:   "package_name_301",
// 					IconURL:       "icon_url_301",
// 					ImageURL:      "image_url_301",
// 					ImpressionURL: "impression_url_301",
// 					VideoURL:      "",
// 				},
// 			})
// 		})

// 		Convey("VideoAdType == 3", func() {
// 			adList = []Ad{
// 				Ad{
// 					CampaignID:    int64(200),
// 					OfferID:       int(300),
// 					AppName:       "app_name_300",
// 					AppDesc:       "app_desc_300",
// 					PackageName:   "package_name_300",
// 					IconURL:       "icon_url_300",
// 					ImageURL:      "image_url_300",
// 					ImpressionURL: "impression_url_300",
// 					VideoURL:      "test_video_url",
// 				},
// 				Ad{
// 					CampaignID:    int64(201),
// 					OfferID:       int(301),
// 					AppName:       "app_name_301",
// 					AppDesc:       "app_desc_301",
// 					PackageName:   "package_name_301",
// 					IconURL:       "icon_url_301",
// 					ImageURL:      "image_url_301",
// 					ImpressionURL: "impression_url_301",
// 					VideoURL:      "",
// 				},
// 			}
// 			VideoAdType = 3
// 			res = VideoAdTypeFilter(adList, VideoAdType)
// 			So(res, ShouldResemble, []Ad{
// 				Ad{
// 					CampaignID:    int64(200),
// 					OfferID:       int(300),
// 					AppName:       "app_name_300",
// 					AppDesc:       "app_desc_300",
// 					PackageName:   "package_name_300",
// 					IconURL:       "icon_url_300",
// 					ImageURL:      "image_url_300",
// 					ImpressionURL: "impression_url_300",
// 					VideoURL:      "test_video_url",
// 				},
// 			})
// 		})
// 	})
// }

func TestParseMobvistaData(t *testing.T) {
	Convey("test ParseMobvistaData", t, func() {
		// todo:
	})
}

func TestGetCampaignInfoFromMongo(t *testing.T) {
	// todo
}

/*
func TestGetCampaignDecodeInfo(t *testing.T) {
	Convey("test GetCampaignDecodeInfo", t, func() {
		// Mock GetCampaignInfoFromMongo()
		guard := Patch(GetCampaignInfoFromMongo, func(campaignIds []int64) (map[int64]smodel.CampaignInfo, error) {
			res := make(map[int64]smodel.CampaignInfo)
			for _, cID := range campaignIds {
				res[cID] = smodel.CampaignInfo{
					CampaignId: cID,
				}
			}
			return res, nil
		})
		defer guard.Unpatch()

		var campaignIds []int64

		Convey("campaignIds 为空", func() {
			campaignIds = []int64{}
			result, error := GetCampaignInfo(campaignIds)
			So(result, ShouldBeNil)
			So(error, ShouldBeError)
		})

		Convey("redis 有数据，且 json 合法", func() {
			// Mock redis.RedisClientHMGet()
			guard = Patch(redis.RedisClientHMGet, func(tablename string, args []interface{}) ([]string, error) {
				endate := []byte("\x1b\xfd\n\x00\xe4\xb7\xa6\xf6\xa7\xab\xf5\xa4\xdc\xc60\x8d*\xa1\x95S\xdb\xb6\x82\x06\x06;\xfc\x98\xb2\x80\xfd\xe3\xad\xbd\xde\x0e\x1b\xbc\xc9i\xfd\xaf\xa5C\xa6t\xe1p\b\x8bw\xcc\xfb\xff\xeefv\xef\xd2\xda\x86R\xfb\xe5\xe8.\x83\n\xc2\x81L\x14\xb2UU\xaa\xad\xc6\xb3\b\xa7\xc3\xe5\"\xc2\xa1\xecv\x18\xbclD\x82M\x0c\x832\xc6M\r\x81H\xe7Wg\xcd\x98\xbe\xd7VU_\xc6J.\x9aQ\x12\x0b\x99\x8b|\xa7H\xe7U\x19\xa5*1U\xac\xab\x98^\xbc\x95\xeb\xaa\xe6\x11bq\xd3\xe0\x94\b\xa1\x88t\x1e\x10\x10\x16\x19\xa5\xfb,\xcc\x10\xc4\"\n\xc0\x86(\xf1\x02\xbf\x93\xb3\xdc=\x9e\x11F\xaeh$R\xb8s\xa0:\xefc`\x13\x04\"\x97s\xa6/LW\x9aA\xbc-Zx$\x0b\xcd^\x11yn\xdat\x92\xc4[6\xb50\xa6z\xe4\xb9q\xe1\xf8\x1e\xa1\x1a3\x1e\xa0\x83\xeaP:\xd4m\xcbWf\x81(\xc2v\xe0\x86\xd2\xd9>\xd1>\"wnRB\xba\xa9\xa2\xb8]H\x8b\xdb\x16\xf5o\xe8i\xd0\x06\xe4\xdc=;\x01\x81\xab\xda\xe4@:\"\xda\x8fU\x928h\x86\xf1\xbc\x12\x05*\xd1964M\xb3\xd8\n\xf6t=\xf5\xb8a\\PQ\xa2\xeeJ*\\\x1f\x82-\xdb\x1d\"W\x95L\x1d\xac\xcb{X\xad\x87\x0e\xf6\xd0\xc1o\xe6\xe1\x1d\xb0=\x04\x9aI\x0c\x9d\x12\xce\x17\xfe\xac\x10\rG2KX\xa6\xe0U\x9d\xd5y\x95SnT-\xd3\xa3Uii\xd8\x92\x02\x0bR\xa7D\x80\xc9\xf1\xa2\xcdw&8\xcd\r\x8f\xed>.\xcd\n'\x8c\x92Hpx]\xde\xf2\x88\xb3=\xf2V\xb6<m\xf2\xd2\x1e\x8a&\xe4\x0e\x1d n\xcd\x84\x10\xfc\x7f\x92\xbcp\xa3D\x06M\xa7oz\xe4\xac\xfd\x1e\x8c:\xfc\x8e\xb3ze\xff-9\xfe\xa8C\x15n\xf2\x0b\xf5\xc1\xac]Q#\x94\x86\xa1hJ\x82\xbd\xab\x06\x82$\x0b\xa7\xe7(\x82\x10\xbd\x0f\x83M4lX\x1a\xe5\xc4@\x10\xc4r\xe1\x04\xca,\x93q\xd3\xc4\xf8`\xfcL\xfd\xd3;}\x89\x0c\xc0\xde\b\xc3c\xb0\xf9\x90\r\x8b\x1b\xb9\x19\xa5*Z\xb9Wk\x91Ic3\r^\a\xe7\x18I\x9b\xa2\xa3\xa4T\xf2n\xe3\xcf\xc4\x9d\xa5\x8a4\xcd\xd5\x00F\xb6\xee\xf1\x93Hq#\x8e\xf6\xa5\x8a\x9f\xc4y+\n2\xe9\xd8\x05\xcd\xc9\x80\xdb\xf4n=\x15\xc0\x97G\xc1\xa6\b\xb2;\xcf%7bad1\xb6q_G_(\xc1\x86\xb8\xe5\\&&\xec\x045\x14E\xb3y\xa9\x01\xd9\x11f\x92<\x15b\xe7r\xac`@@\x14\xf2\x96Ws\xeb\x92\xaeyi\xcb\xd1K\xddL\xb2\xa3\x97y\x99\xe4\xb3^\xee\xe7^\xe6\xe6R\xdb\xe1\xd5\xa4>\x96\xc6\xe9\x88e)\x18\x18\xaa1\x11\xc6`\xe3\xbb\x90X@o\xe2\xb0q\xff5\xf4\x1a\xd3\xca\xb2\xf0\x0b\xa4P+U\xbc\xa8\xe9\xe7\xcb\xca\x1b\xc7\xcb\x83\x8f\xcb\a\x97\x96\xef\xbf.\xbf]\xee\x1e\xbb\xd2=z\x7f\xf1\xe7\xe9\xf2\xd5\xef\xa5'\xfb\xbaG~\xff=\xfa\xb5\xbb\xff\x0e\x007\xf6\a\xb9\xcf\x8b\xa1r\xff\x8a;?\xba\xc7>.\xef=\xd5H\xaa\x7f?\xf6\xe2\x1d\xaa\xd4\xd2\xb42\x9d\xfeu\xb07\xbb\xb4<s\xa0\xbc\xf1f\xf1\xcb\xc9\xf2\xcc\t\xfc\xcd\x0e\xbf\xef\x1e\xfbX\x1ex\xf6\xec\x7f?N\xfc\xdd\xf7\xad{\xe6F\xf9\xe6xy\xfa\xd8\xf2\x9d3KO\xf6\xfd\xbdv\xa0g}\xee\xe2\xd7\xe3\xcb{\xbf\xfc\xd9\xb3\xafZY:{\xb9{\xf1K\xf9\xe0\xda\xbf\x1f\xd7\x1a\x15\xa6=c\xf1\xfb\xc3\xf2\xc1%\x84\xa5?{\xf6\xf1\xc1]\x0f3:e\xf1\xf7\xcd\xa5\x97Gi\xa5.<a\x90\xe9hR\x0c_\xff\xd0=\xf6\xf1iJ\xf7\xe0{@P\xfc\x1c\xab\xcc\x85U!c%\x13\x9d\xffZQ\x9a\xcfy^\xaa\xf5\xf8\xf3y'\xf6Ub\xa8\x98\xa9\x84\xa9DW\tW5\xd3\xd75\x1e0\xd3\xa8\x0b\xb3\x88\xb2\x1a\x8b+\xc3m\xfa\xcfhl\n\xad\x80\x02\x9b_\x15G\xe2\xbcYa\xd3\xb2\xc583s\xa9I\x89\xc6\xdbi!;\a*\x8d\b\x91iE\xd0H\xb2Nx\xd8\xbeU\x0f\x1b$\x0e`*\xeb\x16\xf4?\xc6\xe4\x9aYQ\x04`o\x87n\xa5\xdb\xde\x88\x11\xe5\x8a\xcd;7\x8e\xf4\t\x7f.\xc8'J\xc7<\x97X\xf5\xff\xdeoZ\xecw\xfaDJ\x85\x99\xdc;\x1d\xd1u\x03\x82\\\xce\xcf\xa3\x0f(\x85q#\xdc\x04\x8b\\S\xfcn\xa5\xac\r\xc2\x8f (F\"\xdcJ\xc8\xf6\x10\xd9i\xa5\x90\xab\xb3\xe6x\x98\x17S\xb0\x83\x81hz\x98O&\xd3\xa2\xf3=\xc1\xa4\x85\xf68\xf8?\x84\x87AW(\x0e\xaf\xea\xb6'~\x17M*n\xec\xedm{bU\x99F\xb3\xc2\\\xda~\x8bLzVluanH.\x03{#C\xdaf\x04Y\xf0+\xbf\xba\xc3\xb8\xd1\x9fR8)\x95\x82XN\xa5C\xa9\\A\xf4u(^\xa8\xc9\xb6\n\x96\x17\x85\xdfWP\xb0\x892\xbb\x84lO\xb9OY<.\xa9T+\xc3c\x95j\xa5\xc0\x14\xc3D\xba#zov\x9bS6\x02\xa1\x18/\xe8\xd48\x0f3\x8a\x17\xb4l\xd3\x83\xf1\x02\x1d\xb9h\x96\x9b\xf3\x11\x86\x13\xdc<z\xe8\xd0\x9a\x91\xfe\x81\xa9\xa4b`\x85\x97<\xbc\xd1@\x141\xa4!\x1dqdnF\xd0\x84\xc1\xb9\xce\xd5\xa8\x899\xe5\x16A\xe0&\xd4!\xcd\n\x88\xd9\x88=Ap\xe2\xf9\xc4\x92\xd47\xea\x84\xfb\x80\xe0[?u\xdb\xba\xb0\xab\x9d\x1cm\xd9O\x92\xb1V\x90l\n$\xd9i\x16=\x9f\x01\xda\x96\xb7\xbc.\x9b\x15\xf3\x84\x8d\xb0j\xb6\xb6f`\xdc\xadM\xf6\xbb\xe3S}\xb5q!2y_Q\xb0\xfaz\x1f\"\x06\xc7\x18\x9a\x16\x99\x88<[y\xb4@\xdb\xeb\xd1bee\x83{\x98\x89\xa1(x\xc8H\xc6\xa4\xbe3Iv:B[^\xac\xaa)2\x0f\xb7&\xb05\x8e\x1f3\xbazb\xda]\xb5~z\xc0\xed\x9b\x9a\x1c\x1c\x19\x02{;\x8c9\xd7\xa5\xa2\x7f;\x15A1\x17frZdE\xe78\xbd\x96\xa6y\xbd\xd9)e]\x85\xbf\xb5\xce_e\xe1\xb3U\xb9\xb5>!\xf2\xe6C\xe4\x88\x93;w\x02" )

				return []string{string(endate)}, nil
			guard = Patch(redis.LocalRedisHMGet, func(tablename string, args []string) ([]string, error) {
				return []string{"{test_client_1}", "test_client_2", "test_client_3"}, nil
			})
			defer guard.Unpatch()
			var camID int64 = 177233485
			campaignIds = []int64{camID}
			result, err := GetCampaignDecodeInfo(campaignIds)
			So(err, ShouldBeNil)
			So(len(result), ShouldEqual, 1)
			info, ok := result[camID]
			So(ok, ShouldBeTrue)
			So(info, ShouldNotBeNil)
			So(info.CampaignId, ShouldEqual, camID)

			//campaignIds = []int64{int64(123), int64(234)}
			//mvutil.Config.HttpConfig.ToMongo = true
			//res, _ = GetCampaignInfo(campaignIds)
			//So(res, ShouldResemble, map[int64]smodel.CampaignInfo{
			//	int64(123): smodel.CampaignInfo{
			//		CampaignId: int64(123),
			//	},
			//	int64(234): smodel.CampaignInfo{
			//		CampaignId: int64(234),
			//	},
			//})
		})
	})
}
*/

// func TestGetCampaignInfo(t *testing.T) {
// 	Convey("test GetCampaignInfo", t, func() {
// 		Convey("开关关掉", func() {
// 			//Mock IsDeconde
// 			isDecode := Patch(IsDeconde, func() (isDecode bool) {
// 				return false
// 			})
// 			defer isDecode.Unpatch()

// 			// Mock GetCampaignInfoFromMongo()
// 			guard := Patch(GetCampaignInfoFromMongo, func(campaignIds []int64) (map[int64]*smodel.CampaignInfo, error) {
// 				res := make(map[int64]*smodel.CampaignInfo)
// 				for _, cID := range campaignIds {
// 					res[cID] = &smodel.CampaignInfo{
// 						CampaignId: cID,
// 					}
// 				}
// 				return res, nil
// 			})
// 			defer guard.Unpatch()

// 			var res map[int64]*smodel.CampaignInfo
// 			var campaignIds []int64

// 			var areaConfig mvutil.AreaConfig
// 			areaConfig.HttpConfig.ToMongo = true
// 			mvutil.Config.AreaConfig = &areaConfig

// 			Convey("campaignIds 为空", func() {
// 				campaignIds = []int64{}
// 				result, error := GetCampaignInfo(campaignIds)
// 				So(result, ShouldBeNil)
// 				So(error, ShouldBeError)
// 			})

// 			Convey("redis 有数据，但 json 不合法", func() {
// 				// Mock redis.RedisClientHMGet()
// 				guard = Patch(redis.LocalRedisHMGet, func(tablename string, args []string) ([]string, error) {
// 					return []string{"{test_client_1}", "test_client_2", "test_client_3"}, nil
// 				})
// 				defer guard.Unpatch()

// 				campaignIds = []int64{int64(123), int64(234)}

// 				res, _ = GetCampaignInfo(campaignIds)
// 				So(res, ShouldResemble, map[int64]*smodel.CampaignInfo{
// 					int64(123): &smodel.CampaignInfo{
// 						CampaignId: int64(123),
// 					},
// 					int64(234): &smodel.CampaignInfo{
// 						CampaignId: int64(234),
// 					},
// 				})
// 			})

// 			Convey("redis 有数据，json 合法", func() {
// 				// todo:
// 				// Mock redis.RedisClientHMGet()
// 				// guard = Patch(redis.RedisClientHMGet, func(tablename string, args []interface{}) ([]string, error) {
// 				// 	return []string{"{\"CampaignId\":111123,\"ThirdParty\":\"test_third_party\"}", "test_client_2", "test_client_3"}, nil
// 				// })
// 				// defer guard.Unpatch()

// 				// campaignIds = []int64{int64(123), int64(234)}
// 				// res, _ = GetCampaignInfo(campaignIds)
// 				// So(res, ShouldResemble, map[int64]smodel.CampaignInfo{
// 				// 	int64(111123): smodel.CampaignInfo{
// 				// 		CampaignId: int64(111123),
// 				// 		ThirdParty: "test_third_party",
// 				// 	},
// 				// })
// 			})
// 		})
// 		Convey("开关打开", func() {
// 			//Mock IsDeconde
// 			isDecode := Patch(IsDeconde, func() (isDecode bool) {
// 				return true
// 			})
// 			defer isDecode.Unpatch()

// 			var areaConfig mvutil.AreaConfig
// 			areaConfig.HttpConfig.RedisDecode = true
// 			mvutil.Config.AreaConfig = &areaConfig

// 			// Mock GetCampaignInfoFromMongo()
// 			guard := Patch(GetCampaignInfoFromMongo, func(campaignIds []int64) (map[int64]*smodel.CampaignInfo, error) {
// 				res := make(map[int64]*smodel.CampaignInfo)
// 				for _, cID := range campaignIds {
// 					res[cID] = &smodel.CampaignInfo{
// 						CampaignId: cID,
// 					}
// 				}
// 				return res, nil
// 			})
// 			defer guard.Unpatch()

// 			var campaignIds []int64

// 			Convey("campaignIds 为空", func() {
// 				campaignIds = []int64{}
// 				result, error := GetCampaignInfo(campaignIds)
// 				So(result, ShouldBeNil)
// 				So(error, ShouldBeError)
// 			})

// 			// Convey("redis 有数据，且 json 合法", func() {
// 			// 	// Mock redis.RedisClientHMGet()
// 			// 	guard = Patch(redis.LocalRedisHMGet, func(tablename string, args []string) ([]string, error) {
// 			// 		endate := []byte("\x1b\xfd\n\x00\xe4\xb7\xa6\xf6\xa7\xab\xf5\xa4\xdc\xc60\x8d*\xa1\x95S\xdb\xb6\x82\x06\x06;\xfc\x98\xb2\x80\xfd\xe3\xad\xbd\xde\x0e\x1b\xbc\xc9i\xfd\xaf\xa5C\xa6t\xe1p\b\x8bw\xcc\xfb\xff\xeefv\xef\xd2\xda\x86R\xfb\xe5\xe8.\x83\n\xc2\x81L\x14\xb2UU\xaa\xad\xc6\xb3\b\xa7\xc3\xe5\"\xc2\xa1\xecv\x18\xbclD\x82M\x0c\x832\xc6M\r\x81H\xe7Wg\xcd\x98\xbe\xd7VU_\xc6J.\x9aQ\x12\x0b\x99\x8b|\xa7H\xe7U\x19\xa5*1U\xac\xab\x98^\xbc\x95\xeb\xaa\xe6\x11bq\xd3\xe0\x94\b\xa1\x88t\x1e\x10\x10\x16\x19\xa5\xfb,\xcc\x10\xc4\"\n\xc0\x86(\xf1\x02\xbf\x93\xb3\xdc=\x9e\x11F\xaeh$R\xb8s\xa0:\xefc`\x13\x04\"\x97s\xa6/LW\x9aA\xbc-Zx$\x0b\xcd^\x11yn\xdat\x92\xc4[6\xb50\xa6z\xe4\xb9q\xe1\xf8\x1e\xa1\x1a3\x1e\xa0\x83\xeaP:\xd4m\xcbWf\x81(\xc2v\xe0\x86\xd2\xd9>\xd1>\"wnRB\xba\xa9\xa2\xb8]H\x8b\xdb\x16\xf5o\xe8i\xd0\x06\xe4\xdc=;\x01\x81\xab\xda\xe4@:\"\xda\x8fU\x928h\x86\xf1\xbc\x12\x05*\xd1964M\xb3\xd8\n\xf6t=\xf5\xb8a\\PQ\xa2\xeeJ*\\\x1f\x82-\xdb\x1d\"W\x95L\x1d\xac\xcb{X\xad\x87\x0e\xf6\xd0\xc1o\xe6\xe1\x1d\xb0=\x04\x9aI\x0c\x9d\x12\xce\x17\xfe\xac\x10\rG2KX\xa6\xe0U\x9d\xd5y\x95SnT-\xd3\xa3Uii\xd8\x92\x02\x0bR\xa7D\x80\xc9\xf1\xa2\xcdw&8\xcd\r\x8f\xed>.\xcd\n'\x8c\x92Hpx]\xde\xf2\x88\xb3=\xf2V\xb6<m\xf2\xd2\x1e\x8a&\xe4\x0e\x1d n\xcd\x84\x10\xfc\x7f\x92\xbcp\xa3D\x06M\xa7oz\xe4\xac\xfd\x1e\x8c:\xfc\x8e\xb3ze\xff-9\xfe\xa8C\x15n\xf2\x0b\xf5\xc1\xac]Q#\x94\x86\xa1hJ\x82\xbd\xab\x06\x82$\x0b\xa7\xe7(\x82\x10\xbd\x0f\x83M4lX\x1a\xe5\xc4@\x10\xc4r\xe1\x04\xca,\x93q\xd3\xc4\xf8`\xfcL\xfd\xd3;}\x89\x0c\xc0\xde\b\xc3c\xb0\xf9\x90\r\x8b\x1b\xb9\x19\xa5*Z\xb9Wk\x91Ic3\r^\a\xe7\x18I\x9b\xa2\xa3\xa4T\xf2n\xe3\xcf\xc4\x9d\xa5\x8a4\xcd\xd5\x00F\xb6\xee\xf1\x93Hq#\x8e\xf6\xa5\x8a\x9f\xc4y+\n2\xe9\xd8\x05\xcd\xc9\x80\xdb\xf4n=\x15\xc0\x97G\xc1\xa6\b\xb2;\xcf%7bad1\xb6q_G_(\xc1\x86\xb8\xe5\\&&\xec\x045\x14E\xb3y\xa9\x01\xd9\x11f\x92<\x15b\xe7r\xac`@@\x14\xf2\x96Ws\xeb\x92\xaeyi\xcb\xd1K\xddL\xb2\xa3\x97y\x99\xe4\xb3^\xee\xe7^\xe6\xe6R\xdb\xe1\xd5\xa4>\x96\xc6\xe9\x88e)\x18\x18\xaa1\x11\xc6`\xe3\xbb\x90X@o\xe2\xb0q\xff5\xf4\x1a\xd3\xca\xb2\xf0\x0b\xa4P+U\xbc\xa8\xe9\xe7\xcb\xca\x1b\xc7\xcb\x83\x8f\xcb\a\x97\x96\xef\xbf.\xbf]\xee\x1e\xbb\xd2=z\x7f\xf1\xe7\xe9\xf2\xd5\xef\xa5'\xfb\xbaG~\xff=\xfa\xb5\xbb\xff\x0e\x007\xf6\a\xb9\xcf\x8b\xa1r\xff\x8a;?\xba\xc7>.\xef=\xd5H\xaa\x7f?\xf6\xe2\x1d\xaa\xd4\xd2\xb42\x9d\xfeu\xb07\xbb\xb4<s\xa0\xbc\xf1f\xf1\xcb\xc9\xf2\xcc\t\xfc\xcd\x0e\xbf\xef\x1e\xfbX\x1ex\xf6\xec\x7f?N\xfc\xdd\xf7\xad{\xe6F\xf9\xe6xy\xfa\xd8\xf2\x9d3KO\xf6\xfd\xbdv\xa0g}\xee\xe2\xd7\xe3\xcb{\xbf\xfc\xd9\xb3\xafZY:{\xb9{\xf1K\xf9\xe0\xda\xbf\x1f\xd7\x1a\x15\xa6=c\xf1\xfb\xc3\xf2\xc1%\x84\xa5?{\xf6\xf1\xc1]\x0f3:e\xf1\xf7\xcd\xa5\x97Gi\xa5.<a\x90\xe9hR\x0c_\xff\xd0=\xf6\xf1iJ\xf7\xe0{@P\xfc\x1c\xab\xcc\x85U!c%\x13\x9d\xffZQ\x9a\xcfy^\xaa\xf5\xf8\xf3y'\xf6Ub\xa8\x98\xa9\x84\xa9DW\tW5\xd3\xd75\x1e0\xd3\xa8\x0b\xb3\x88\xb2\x1a\x8b+\xc3m\xfa\xcfhl\n\xad\x80\x02\x9b_\x15G\xe2\xbcYa\xd3\xb2\xc583s\xa9I\x89\xc6\xdbi!;\a*\x8d\b\x91iE\xd0H\xb2Nx\xd8\xbeU\x0f\x1b$\x0e`*\xeb\x16\xf4?\xc6\xe4\x9aYQ\x04`o\x87n\xa5\xdb\xde\x88\x11\xe5\x8a\xcd;7\x8e\xf4\t\x7f.\xc8'J\xc7<\x97X\xf5\xff\xdeoZ\xecw\xfaDJ\x85\x99\xdc;\x1d\xd1u\x03\x82\\\xce\xcf\xa3\x0f(\x85q#\xdc\x04\x8b\\S\xfcn\xa5\xac\r\xc2\x8f (F\"\xdcJ\xc8\xf6\x10\xd9i\xa5\x90\xab\xb3\xe6x\x98\x17S\xb0\x83\x81hz\x98O&\xd3\xa2\xf3=\xc1\xa4\x85\xf68\xf8?\x84\x87AW(\x0e\xaf\xea\xb6'~\x17M*n\xec\xedm{bU\x99F\xb3\xc2\\\xda~\x8bLzVluanH.\x03{#C\xdaf\x04Y\xf0+\xbf\xba\xc3\xb8\xd1\x9fR8)\x95\x82XN\xa5C\xa9\\A\xf4u(^\xa8\xc9\xb6\n\x96\x17\x85\xdfWP\xb0\x892\xbb\x84lO\xb9OY<.\xa9T+\xc3c\x95j\xa5\xc0\x14\xc3D\xba#zov\x9bS6\x02\xa1\x18/\xe8\xd48\x0f3\x8a\x17\xb4l\xd3\x83\xf1\x02\x1d\xb9h\x96\x9b\xf3\x11\x86\x13\xdc<z\xe8\xd0\x9a\x91\xfe\x81\xa9\xa4b`\x85\x97<\xbc\xd1@\x141\xa4!\x1dqdnF\xd0\x84\xc1\xb9\xce\xd5\xa8\x899\xe5\x16A\xe0&\xd4!\xcd\n\x88\xd9\x88=Ap\xe2\xf9\xc4\x92\xd47\xea\x84\xfb\x80\xe0[?u\xdb\xba\xb0\xab\x9d\x1cm\xd9O\x92\xb1V\x90l\n$\xd9i\x16=\x9f\x01\xda\x96\xb7\xbc.\x9b\x15\xf3\x84\x8d\xb0j\xb6\xb6f`\xdc\xadM\xf6\xbb\xe3S}\xb5q!2y_Q\xb0\xfaz\x1f\"\x06\xc7\x18\x9a\x16\x99\x88<[y\xb4@\xdb\xeb\xd1bee\x83{\x98\x89\xa1(x\xc8H\xc6\xa4\xbe3Iv:B[^\xac\xaa)2\x0f\xb7&\xb05\x8e\x1f3\xbazb\xda]\xb5~z\xc0\xed\x9b\x9a\x1c\x1c\x19\x02{;\x8c9\xd7\xa5\xa2\x7f;\x15A1\x17frZdE\xe78\xbd\x96\xa6y\xbd\xd9)e]\x85\xbf\xb5\xce_e\xe1\xb3U\xb9\xb5>!\xf2\xe6C\xe4\x88\x93;w\x02")

// 			// 		return []string{string(endate)}, nil
// 			// 	})

// 			// 	defer guard.Unpatch()
// 			// 	var camID int64 = 177233485
// 			// 	campaignIds = []int64{camID}
// 			// 	result, err := GetCampaignDecodeInfo(campaignIds)
// 			// 	So(err, ShouldBeNil)
// 			// 	So(len(result), ShouldEqual, 1)
// 			// 	info, ok := result[camID]
// 			// 	So(ok, ShouldBeTrue)
// 			// 	So(info, ShouldNotBeNil)
// 			// 	So(info.CampaignId, ShouldEqual, camID)
// 			// })
// 		})

// 	})
// }

func TestGetCreativeInfoFromMongo(t *testing.T) {
	// todo
}

// func TestGetCreativeInfo(t *testing.T) {
// 	Convey("test GetCreativeInfo", t, func() {
// 		Convey("开关关掉", func() {
// 			//Mock IsDeconde
// 			isDecode := Patch(IsDeconde, func() (isDecode bool) {
// 				return false
// 			})
// 			defer isDecode.Unpatch()

// 			var res map[int64]mvutil.Content
// 			var creativeID string
// 			var creativeMap map[enum.CreativeType]int64

// 			// Mock GetCreativeInfoFromMongo()
// 			guard := Patch(GetCreativeInfoFromMongo, func(creativeID string, creativeMap map[enum.CreativeType]int64) (map[int64]mvutil.Content, error) {
// 				result := map[int64]mvutil.Content{
// 					int64(101010): mvutil.Content{
// 						CreativeId: int64(42),
// 					},
// 				}
// 				return result, nil
// 			})
// 			defer guard.Unpatch()

// 			Convey("creativeID 为空", func() {
// 				creativeID = ""
// 				result, err := GetCreativeInfo(creativeID, creativeMap)
// 				So(result, ShouldBeNil)
// 				So(err, ShouldBeError)
// 			})

// 			Convey("creativeMap 为空", func() {
// 				result, err := GetCreativeInfo(creativeID, creativeMap)
// 				So(result, ShouldBeNil)
// 				So(err, ShouldBeError)
// 			})

// 			Convey("redis 数据", func() {
// 				// Mock redis.RedisClientHMGet()
// 				guard = Patch(redis.LocalRedisHMGet, func(tablename string, args []string) ([]string, error) {
// 					return []string{"{test_client_1}", "", "test_client_3"}, nil
// 				})
// 				defer guard.Unpatch()

// 				var areaConfig mvutil.AreaConfig
// 				areaConfig.HttpConfig.ToMongo = true
// 				mvutil.Config.AreaConfig = &areaConfig
// 				creativeID = "1234"
// 				creativeMap = map[enum.CreativeType]int64{
// 					enum.CreativeType(10001): int64(10002),
// 					enum.CreativeType(20001): int64(20002),
// 				}

// 				res, _ = GetCreativeInfo(creativeID, creativeMap)
// 				So(res, ShouldResemble, map[int64]mvutil.Content{
// 					int64(101010): mvutil.Content{
// 						CreativeId: int64(42),
// 					},
// 				})
// 			})
// 		})
// 		Convey("开关打开", func() {
// 			//Mock IsDeconde
// 			isDecode := Patch(IsDeconde, func() (isDecode bool) {
// 				return true
// 			})
// 			defer isDecode.Unpatch()

// 			var areaConfig mvutil.AreaConfig
// 			areaConfig.HttpConfig.RedisDecode = true
// 			mvutil.Config.AreaConfig = &areaConfig

// 			// var res map[int64]mvutil.Content
// 			// var creativeID string
// 			// var creativeMap map[enum.CreativeType]int64

// 			// Mock GetCreativeInfoFromMongo()
// 			guard := Patch(GetCreativeInfoFromMongo, func(creativeID string, creativeMap map[enum.CreativeType]int64) (map[int64]mvutil.Content, error) {
// 				result := map[int64]mvutil.Content{
// 					int64(101010): mvutil.Content{
// 						CreativeId: int64(42),
// 					},
// 				}
// 				return result, nil
// 			})
// 			defer guard.Unpatch()

// 			// Convey("creativeID 为空", func() {
// 			// 	creativeID = ""
// 			// 	result, err := GetCreativeDecodeInfo(creativeID, creativeMap)
// 			// 	So(result, ShouldBeNil)
// 			// 	So(err, ShouldBeError)
// 			// })

// 			// Convey("creativeMap 为空", func() {
// 			// 	result, err := GetCreativeDecodeInfo(creativeID, creativeMap)
// 			// 	So(result, ShouldBeNil)
// 			// 	So(err, ShouldBeError)
// 			// })

// 			// Convey("redis 数据", func() {
// 			// 	// Mock redis.RedisClientHMGet()
// 			// 	guard = Patch(redis.LocalRedisHMGet, func(tablename string, args []string) ([]string, error) {
// 			// 		endata := []byte("\x1bh\x00\x00\xc4\xc0m\xa5\xbe\x0c\xdc\xc4C\xa3D1\xf9\a\x87\x9c\xea\x01$i\x9eI \x9c\x05\x139\x91\xa4\x05\xe7\xe7#\x96%\xce\x89\xe7Lq'\x94\x926 \xc7\xda/\x1a\x1e\xc8\xf21L*Z\x93\xa7=\x869Y\xfcG>\xcc`\x06\x8c\x82Y0\x01\xae\xa1\\\xa1\x0bU\n%\x84,\x7f\"\xfc\a")
// 			// 		return []string{string(endata)}, nil
// 			// 	})
// 			// 	defer guard.Unpatch()

// 			// 	//mvutil.Config.HttpConfig.ToMongo = true
// 			// 	creativeID = "1234"
// 			// 	creativeMap = map[enum.CreativeType]int64{
// 			// 		enum.CreativeType(10001): int64(10002),
// 			// 		//enum.CreativeType(20001): int64(20002),
// 			// 	}

// 			// 	res, _ = GetCreativeDecodeInfo(creativeID, creativeMap)
// 			// 	So(res, ShouldNotBeNil)
// 			// 	So(len(res), ShouldEqual, 1)
// 			// 	info, ok := res[10002]
// 			// 	So(ok, ShouldBeTrue)
// 			// 	So(len(info.Url), ShouldNotBeEmpty)
// 			// })
// 		})
// 	})
// }

/*
func TestGetCreativeDecodeInfo(t *testing.T) {
	Convey("test GetCreativeDecodeInfo", t, func() {
		var res map[int64]mvutil.Content
		var creativeID string
		var creativeMap map[enum.CreativeType]int64

		// Mock GetCreativeInfoFromMongo()
		guard := Patch(GetCreativeInfoFromMongo, func(creativeID string, creativeMap map[enum.CreativeType]int64) (map[int64]mvutil.Content, error) {
			result := map[int64]mvutil.Content{
				int64(101010): mvutil.Content{
					CreativeId: int64(42),
				},
			}
			return result, nil
		})
		defer guard.Unpatch()

		Convey("creativeID 为空", func() {
			creativeID = ""
			result, err := GetCreativeDecodeInfo(creativeID, creativeMap)
			So(result, ShouldBeNil)
			So(err, ShouldBeError)
		})

		Convey("creativeMap 为空", func() {
			result, err := GetCreativeDecodeInfo(creativeID, creativeMap)
			So(result, ShouldBeNil)
			So(err, ShouldBeError)
		})

		Convey("redis 数据", func() {
			// Mock redis.RedisClientHMGet()
			guard = Patch(redis.RedisClientHMGetKey, func(tablename string, args []interface{}) ([]string, error) {
				endata := []byte("\x1bh\x00\x00\xc4\xc0m\xa5\xbe\x0c\xdc\xc4C\xa3D1\xf9\a\x87\x9c\xea\x01$i\x9eI \x9c\x05\x139\x91\xa4\x05\xe7\xe7#\x96%\xce\x89\xe7Lq'\x94\x926 \xc7\xda/\x1a\x1e\xc8\xf21L*Z\x93\xa7=\x869Y\xfcG>\xcc`\x06\x8c\x82Y0\x01\xae\xa1\\\xa1\x0bU\n%\x84,\x7f\"\xfc\a")
				return []string{string(endata)}, nil
			guard = Patch(redis.LocalRedisHMGet, func(tablename string, args []string) ([]string, error) {
				return []string{"{test_client_1}", "", "test_client_3"}, nil
			})
			defer guard.Unpatch()

			//mvutil.Config.HttpConfig.ToMongo = true
			creativeID = "1234"
			creativeMap = map[enum.CreativeType]int64{
				enum.CreativeType(10001): int64(10002),
				//enum.CreativeType(20001): int64(20002),
			}

			res, _ = GetCreativeDecodeInfo(creativeID, creativeMap)
			So(res, ShouldNotBeNil)
			So(len(res), ShouldEqual, 1)
			info, ok := res[10002]
			So(ok, ShouldBeTrue)
			So(len(info.Url), ShouldNotBeEmpty)
			//t.Errorf("%+v", res)
			//So(res, ShouldResemble, map[int64]mvutil.Content{
			//	int64(101010): mvutil.Content{
			//		CreativeId: int64(42),
			//	},
			//})
		})
	})
}
*/
