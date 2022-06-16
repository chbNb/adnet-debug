package config

type TkCfg struct {
	TrackHost     string
	PlayTrackPath string
}

func parseTkCfg(fileName string) (*TkCfg, error) {
	viper, err := ReadConfig(fileName)
	if err != nil {
		return nil, err
	}
	tkCfg := &TkCfg{}
	tkCfg.TrackHost = viper.GetString("track_host")
	tkCfg.PlayTrackPath = viper.GetString("play_track_path")
	return tkCfg, nil
}
