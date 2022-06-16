package backend

import (
	"testing"

	. "github.com/bouk/monkey"
	. "github.com/smartystreets/goconvey/convey"
	"gitlab.mobvista.com/ADN/adnet/internal/corsair_proto"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/mtg_openrtb/pkg/mtgrtb"
	smodel "gitlab.mobvista.com/ADN/structs/model"
)

func TestConstructAsInfo(t *testing.T) {
	Convey("ConstructMasInfo", t, func() {
		patch := Patch(extractor.GetCLOUD_NAME, func() string {
			return "awssg"
		})
		defer patch.Unpatch()
		reqCtx := &mvutil.ReqCtx{
			ReqParams: &mvutil.RequestParams{
				Param: mvutil.Params{
					UnitSize: "unitsize",
				},
				UnitInfo: &smodel.UnitInfo{
					Unit: smodel.Unit{
						PlacementId: int64(12345),
					},
				},
			},
		}
		guard := Patch(extractor.GetPassthroughData, func() []string {
			return []string{}
		})
		defer guard.Unpatch()
		asInfo, err := ConstructMasInfo(reqCtx)
		So(err, ShouldBeNil)
		So(asInfo.UnitSize, ShouldNotBeNil)
		So(*asInfo.UnitSize, ShouldEqual, "unitsize")
		So(*asInfo.Extra, ShouldEqual, "awssg_pioneer")

	})
}

func Test_fillMasAd(t *testing.T) {
	Convey("fillMasAd", t, func() {
		extra2 := "extra2"
		advId := int32(2)
		z := "z"
		resp := &mtgrtb.BidResponse{
			AsResp: &mtgrtb.BidResponse_AsResp{
				Extra2: &extra2,
				SdkParam: []*mtgrtb.BidResponse_SdkParam{
					{
						AdvId:       &advId,
						ImageSizeId: &advId,
					},
				},
				UrlParam: []*mtgrtb.BidResponse_UrlParam{
					{
						Z: &z,
					},
				},
			},
		}
		reqCtx := &mvutil.ReqCtx{
			ReqParams: &mvutil.RequestParams{
				Param: mvutil.Params{},
				AsResp: &mtgrtb.BidResponse_AsResp{
					UrlParam: []*mtgrtb.BidResponse_UrlParam{
						{Z: &z},
					},
				},
			},
		}
		ad := new(corsair_proto.Campaign)
		fillMasAd(resp, reqCtx, nil, ad, 0)
		//So(reqCtx.ReqParams.Param.Extra2, ShouldEqual, extra2)
		So(*ad.AdvertiserId, ShouldEqual, advId)
		So(*(reqCtx.ReqParams.AsResp.UrlParam[0].Z), ShouldEqual, z)
	})
}
