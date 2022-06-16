package mvutil

type DistributionInfo struct {
	FlowTagId        int               `bson:"flow_tag_id,omitempty" json:"flow_tag_id"`
	Name             string            `bson:"name,omitempty" json:"name"`
	AdBackends       []int             `bson:"ad_backends,omitempty" json:"ad_backends"`
	Ratio            int               `bson:"ratio,omitempty" json:"ratio"`
	AppIds           []int             `bson:"app_ids,omitempty" json:"app_ids"`
	UnitIds          []int             `bson:"unit_ids,omitempty" json:"unit_ids"`
	CountryCodes     []string          `bson:"country_codes,omitempty" json:"country_codes"`
	CityCodes        []uint64          `bson:"city_codes,omitempty" json:"city_codes"`
	AdTypes          []int             `bson:"ad_types,omitempty" json:"ad_types"`
	AndroidVersions  []int64           `bson:"android_versions,omitempty" json:"android_versions"`
	IosVersions      []int64           `bson:"ios_versions,omitempty" json:"ios_versions"`
	Platforms        []int             `bson:"platforms,omitempty" json:"platforms"`
	DistributionMode int               `bson:"distribution_mode,omitempty" json:"distribution_mode"`
	Weight           int               `bson:"weight,omitempty" json:"weight"`
	Mtime            int               `bson:"mtime,omitempty" json:"mtime"`
	Status           int               `bson:"status,omitempty" json:"status"`
	Updated          int64             `bson:"updated,omitempty" json:"updated"`
	AdReqKeys        map[string]string `bson:"adReqKeys,omitempty" json:"adReqKeys"`
	DeviceType       []int             `bson:"deviceType,omitempty" json:"deviceType"`
	NetworkType      []int             `bson:"networkType,omitempty" json:"networkType"`
	OrienDevIds      []string          `bson:"orienDevIds,omitempty" json:"orienDevIds"`
	SDKVersion       InfoSDKVersion    `bson:"sdk_version,omitempty" json:"sdkVersion"`
}

type InfoSDKVersion struct {
	Include []*InfoSDKVersionItem `bson:"include,omitempty" json:"include"`
	Exclude []*InfoSDKVersionItem `bson:"exclude,omitempty" json:"exclude"`
}

type InfoSDKVersionItem struct {
	Max string `bson:"max,omitempty" json:"max"`
	Min string `bson:"min,omitempty" json:"min"`
}

type DistributionData struct {
	FlowTagId        int             `bson:"flow_tag_id,omitempty" json:"flow_tag_id"`
	Name             string          `bson:"name,omitempty" json:"name"`
	AdBackends       []int           `bson:"ad_backends,omitempty" json:"ad_backends"`
	Ratio            int             `bson:"ratio,omitempty" json:"ratio"`
	AppIds           map[int]bool    `bson:"app_ids,omitempty" json:"app_ids"`
	UnitIds          map[int]bool    `bson:"unit_ids,omitempty" json:"unit_ids"`
	CountryCodes     map[string]bool `bson:"country_codes,omitempty" json:"country_codes"`
	CityCodes        map[uint64]bool `bson:"city_codes,omitempty" json:"city_codes"`
	AndroidVersions  []int64         `bson:"android_versions,omitempty" json:"android_versions"`
	IosVersions      []int64         `bson:"ios_versions,omitempty" json:"ios_versions"`
	Platforms        map[int]bool    `bson:"platforms,omitempty" json:"platforms"`
	AdTypes          map[int]bool    `bson:"ad_types,omitempty" json:"ad_types"`
	DistributionMode int             `bson:"distribution_mode,omitempty" json:"distribution_mode"`
	Weight           int             `bson:"weight,omitempty" json:"weight"`
	Mtime            int             `bson:"mtime,omitempty" json:"mtime"`
	Updated          int64           `bson:"updated,omitempty" json:"updated"`
	AdReqKeys        map[int]string  `bson:"adReqKeys,omitempty" json:"adReqKeys"`
	DeviceType       map[int]bool    `bson:"device_type,omitempty" json:"deviceType"`
	NetworkType      map[int]bool    `bson:"network_type,omitempty" json:"networkType"`
	OrienDevIds      map[string]bool `bson:"orienDevIds,omitempty" json:"orienDevIds"`
	IOSSDKVersion    *DataSDKVersion `json:"iosSDKVersion"`
	AndroidSDKVerion *DataSDKVersion `json:"androidSDKVersion"`
}

type DataSDKVersion struct {
	Include []*DataSDKVersionItem `json:"include"`
	Exclude []*DataSDKVersionItem `json:"exclude"`
}

type DataSDKVersionItem struct {
	Max int32 `json:"max"`
	Min int32 `json:"min"`
}
