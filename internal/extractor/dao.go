package extractor

import (
	"errors"

	"gitlab.mobvista.com/ADN/structs/constant"
	smodel "gitlab.mobvista.com/ADN/structs/model"
)

var DemandDao = new(Dao)

type Dao struct{}

func (d *Dao) GetCampaignInfo(id int64, force ...bool) (*smodel.CampaignInfo, error) {
	camp, ifind := GetCampaignInfo(id)
	if camp != nil && ifind {
		return camp, nil
	}

	return nil, errors.New("not found campaign")
}

func (d *Dao) GetAdvOfferByUUID(uuid string, force ...bool) (*smodel.AdvOffer, error) {
	return nil, errors.New("don't support")
}

func (d *Dao) GetMVConfig(key string) (interface{}, bool) {
	return GetMVConfigValue(key)
}

func (d *Dao) GetAppPackageMtgID(appPackage string) (int64, bool) {
	mtgId := GetAppPackageMtgID(appPackage)
	if mtgId != constant.DefaultMtgId {
		return mtgId, true
	}
	return mtgId, false
}

// 先放着，避免编译不过
func (d *Dao) GetCampaignRTAutoRateInfo(key string) (*smodel.CampaignRTAutoRateInfo, bool) {
	return &smodel.CampaignRTAutoRateInfo{}, false
}
