package process_pipeline

import (
	"errors"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/output"
	"gitlab.mobvista.com/ADN/adnet/internal/uuid"
)

type MappingServerFilter struct {
}

func (taf *MappingServerFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*mvutil.RequestParams)
	if !ok {
		return nil, errors.New("MappingServerFilter input type should be *mvutil.ReqCtx")
	}
	// bkup_id为空，生成bkup_id
	NewBkupId(in)

	output.RenderNewEncryptedBkupId(in)
	// 查询生成ruid				带上参数发送http请求
	output.RenderRuid(in, "")
	// 这里的sysid和bkupid可能更新，需要更新ExtSysId
	in.Param.ExtSysId = in.Param.SysId + "," + in.Param.BkupId
	return in, nil
}

func NewBkupId(r *mvutil.RequestParams) {
	// 新版本才处理
	if r.Param.ApiVersion < mvconst.API_VERSION_2_2 {
		return
	}
	if !mvutil.IsHbOrV3OrV5Request(r.Param.RequestPath) {
		return
	}
	if len(r.Param.BkupId) == 0 {
		if r.Param.PlatformName == mvconst.PlatformNameIOS {
			// 请求mapping server生成加密的ruid
			output.RenderRuid(r, "encryption")
			adnetSwitchConf, _ := extractor.GetADNET_SWITCHS()
			newV4Backup, ok := adnetSwitchConf["newV4Backup"]
			if len(r.Param.EncryptedRuid) == 0 && ok && newV4Backup == 1 {
				// 没有切量或切量mapping server生成bkupid失败的情况下，都使用setting逻辑生成bkupid
				v4, _ := uuid.NewV4()
				r.Param.BkupId = v4.String()
			} else {
				r.Param.BkupId = r.Param.EncryptedRuid
			}
		} else {
			// 对于安卓新版本，目前还是使用setting 旧逻辑生成bkupid，等mapping server支持安卓的时候，再请求mapping server
			v4, _ := uuid.NewV4()
			r.Param.BkupId = v4.String()
		}

		r.Param.NewEncryptedBkupId = output.NewEncryptDevId(r, r.Param.BkupId)
		return
	}
}
