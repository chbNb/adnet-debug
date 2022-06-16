package config

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestParseConfig(t *testing.T) {
	convey.Convey("ParseConfig ok", t, func() {
		cfg, err := ParseConfig("testdata", "cloud", "test")
		convey.So(err, convey.ShouldBeNil)
		convey.So(cfg, convey.ShouldNotBeNil)
	})
}

func TestDefaultConfigDir(t *testing.T) {
	convey.Convey("DefaultConfigDir ok", t, func() {
		dftConfig := DefaultConfigDir()
		convey.So(dftConfig, convey.ShouldEqual, DefaultConfigDir())
	})
}

func TestDefaultSrvAddr(t *testing.T) {
	convey.Convey("DefaultSrvAddr ok", t, func() {
		dftAddr := DefaultSrvAddr()
		convey.So(dftAddr, convey.ShouldEqual, ":9102")
	})
}

func TestConfigData(t *testing.T) {
	convey.Convey("SetConfigValue ok", t, func() {
		cfg, err := ParseConfig("testdata", "cloud", "test")
		convey.So(err, convey.ShouldBeNil)
		SetConfig(cfg)
		cfgData := GetConfig()
		convey.So(cfgData, convey.ShouldNotBeNil)
		convey.So(cfgData.ServerCfg.HTTPAddr, convey.ShouldEqual, ":9102")
		convey.So(cfgData.TBConfigPath, convey.ShouldEqual, "testdata/cloud/test/treasure_box.yaml")
		convey.So(cfgData.ConsulCfg.Address, convey.ShouldEqual, "127.0.0.1:8500")
		convey.So(cfgData.ConsulCfg.Aerospike.ServiceName, convey.ShouldEqual, "hb-aerospike")
		convey.So(cfgData.ConsulCfg.Adx.ServiceName, convey.ShouldEqual, "adx")
		convey.So(cfgData.ConsulCfg.AdnetAerospike.ServiceName, convey.ShouldEqual, "adnet-aerospike")

		// convey.So(cfgData.NetSvrCfg.GRPCAddress, convey.ShouldEqual, "a8644b6f65ba7404b93e678c19a1c62e-114010529.us-east-1.elb.amazonaws.com:9000")
	})
}

func TestBidWhiteListConfigData(t *testing.T) {
	convey.Convey("SetConfigValue ok", t, func() {
		cfg, err := ParseConfig("testdata", "cloud", "test")
		convey.So(err, convey.ShouldBeNil)
		SetConfig(cfg)
		cfgData := GetConfig()
		convey.So(cfgData, convey.ShouldNotBeNil)
	})
}
