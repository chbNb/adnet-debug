package output

import (
	"fmt"
	"sync"
	"sync/atomic"

	jsoniter "github.com/json-iterator/go"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

// 广告抓包 支持多用户
// 当抓包开关打开且adreq or adresp条件满足时抓包

var (
	_snaps      = make(map[string]*AdSnapshot)
	_snaps_copy = make(map[string]*AdSnapshot)
	mux         = new(sync.RWMutex)
)

func copysnap() {
	snaps := make(map[string]*AdSnapshot)
	for k, v := range _snaps {
		snaps[k] = v
	}
	_snaps_copy = snaps
}

// 添加一个抓包条件
// 如果条件已存在则返回已存在的条件
// 如果条件不存在，则返回新增的条件
func SetGetSnapshot(s *AdSnapshot) (*AdSnapshot, bool) {
	mux.Lock()
	defer mux.Unlock()
	tmp, ok := _snaps[s.Sign()]
	if ok {
		s.Stop()
		return tmp, ok
	}

	_snaps[s.Sign()] = s
	copysnap()
	return s, false
}

func DelSnapshot(s *AdSnapshot) {
	mux.Lock()
	delete(_snaps, s.Sign())
	copysnap()
	mux.Unlock()
}

// CaptureAdPack 获取广告包
func CaptureAdPack(req *mvutil.Params, res *MobvistaResult) {
	for _, snap := range _snaps_copy {
		snap.Snap(req, res)
	}
}

type AdPack struct {
	Req *mvutil.Params
	Res interface{}
}

type AdSnapshot struct {
	once   *sync.Once
	done   chan struct{}
	num    int32
	res    []*AdPack
	params *mvutil.Params
}

func NewAdSnapshot(params *mvutil.Params, num int) *AdSnapshot {
	return &AdSnapshot{
		once:   new(sync.Once),
		done:   make(chan struct{}),
		num:    int32(num),
		res:    make([]*AdPack, 0, num),
		params: params,
	}
}

// Sign 唯一条件
func (s *AdSnapshot) Sign() string {
	return mvutil.Md5(fmt.Sprintf("%d_%d_%d_%d_%s_%s_%s_%s_%d",
		s.params.PublisherID, s.params.AppID, s.params.UnitID, s.params.CampaignID,
		s.params.GAID, s.params.IDFA, s.params.ClientIP, s.params.Scenario, s.params.AdType))
}

func (s *AdSnapshot) Condition(params *mvutil.Params, res *MobvistaResult) (ok bool) {
	if s.params.PublisherID != 0 && s.params.PublisherID != params.PublisherID {
		return
	}

	if s.params.AppID != 0 && s.params.AppID != params.AppID {
		return
	}

	if s.params.UnitID != 0 && s.params.UnitID != params.UnitID {
		return
	}

	if s.params.GAID != "" && s.params.GAID != params.GAID {
		return
	}

	if s.params.IDFA != "" && s.params.IDFA != params.IDFA {
		return
	}

	if s.params.ClientIP != "" && s.params.ClientIP != params.ClientIP {
		return
	}

	if s.params.Scenario != "" && s.params.Scenario != params.Scenario {
		return
	}

	if s.params.AdType != 0 && s.params.AdType != params.AdType {
		return
	}

	if s.params.CampaignID != 0 {
		var hadMatchCampaign bool
		for _, ad := range res.Data.Ads {
			if ad.CampaignID == s.params.CampaignID {
				hadMatchCampaign = true
			}
		}

		if !hadMatchCampaign {
			return
		}
	}

	if atomic.AddInt32(&s.num, -1) < 0 {
		s.Stop()
		return
	}
	return true
}

func (s *AdSnapshot) Done() <-chan struct{} {
	return s.done
}

func (s *AdSnapshot) Download() ([]byte, error) {
	data, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(s.res)
	return data, err
}

func (s *AdSnapshot) Snap(req *mvutil.Params, res *MobvistaResult) {
	if !s.Condition(req, res) {
		return
	}

	s.res = append(s.res, &AdPack{
		Req: req,
		Res: res,
	})

	if s.num < 0 {
		s.Stop()
	}
}

func (s *AdSnapshot) Stop() {
	s.once.Do(func() {
		close(s.done)
	})
}
