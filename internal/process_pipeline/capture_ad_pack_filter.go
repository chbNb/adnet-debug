package process_pipeline

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/output"
)

type CaptureAdPackFilter struct {
}

func (c *CaptureAdPackFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*http.Request)
	if !ok {
		return nil, errors.New("QueryMemDataFilter input type should be *http.Request")
	}

	params := new(mvutil.Params)
	if publisherId := in.FormValue("publisher_id"); publisherId != "" {
		params.PublisherID, _ = strconv.ParseInt(publisherId, 10, 64)
	}

	if appId := in.FormValue("app_id"); appId != "" {
		params.AppID, _ = strconv.ParseInt(appId, 10, 64)
	}

	if unitId := in.FormValue("unit_id"); unitId != "" {
		params.UnitID, _ = strconv.ParseInt(unitId, 10, 64)
	}

	if campaignId := in.FormValue("campaign_id"); campaignId != "" {
		params.CampaignID, _ = strconv.ParseInt(campaignId, 10, 64)
	}

	if gaid := in.FormValue("gaid"); gaid != "" {
		params.GAID = gaid
	}

	if idfa := in.FormValue("idfa"); idfa != "" {
		params.IDFA = idfa
	}

	if clientIp := in.FormValue("client_ip"); clientIp != "" {
		params.ClientIP = clientIp
	}

	if scenario := in.FormValue("scenario"); scenario != "" {
		params.Scenario = scenario
	}

	if adType := in.FormValue("ad_type"); adType != "" {
		adTypeInt, err := strconv.Atoi(adType)
		if err == nil {
			params.AdType = int32(adTypeInt)
		}
	}

	var num int
	if value := in.FormValue("num"); value != "" {
		num, _ = strconv.Atoi(value)
	}

	if num <= 0 {
		num = 1
	}

	var timeout int
	if value := in.FormValue("timeout"); value != "" {
		timeout, _ = strconv.Atoi(value)
	}
	if timeout <= 0 {
		timeout = 10
	}

	snap := output.NewAdSnapshot(params, num)
	snap, ok = output.SetGetSnapshot(snap)

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(timeout))
	defer cancel()

	select {
	case <-snap.Done():
	case <-ctx.Done():
		snap.Stop()
	}

	if !ok {
		output.DelSnapshot(snap)
	}

	var res string
	snapData, err := snap.Download()
	if err != nil {
		res = fmt.Sprintf("download err %s", err.Error())
		return &res, nil
	}

	res = string(snapData)
	return &res, nil
}
