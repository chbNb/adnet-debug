package hot_data

import (
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"

	mlogger "github.com/mae-pax/logger"
	"gitlab.mobvista.com/ADN/adnet/internal/redis"
	"gitlab.mobvista.com/ADN/adnet/internal/watcher"
)

const (
	campaignSetPreKey = "active_data_"
	maxHotDataNumber  = 20000

	CAMPAIGN  = "campaign"
	UNIT      = "unit"
	APP       = "app"
	PUBLISHER = "publisher"
)

type activeDataUnit struct {
	mutex sync.Mutex
	ids   map[int64]bool
}

var logger *mlogger.Log

var activeDataCollecter map[string]*activeDataUnit
var activeDataSwitch map[string]bool

func splitStr(str string) []string {
	if str == "" {
		return []string{}
	}
	str = strings.Replace(str, " ", "", -1)
	arr := strings.Split(str, ",")
	return arr
}

func InitActiveDataCollecter(str []string, log *mlogger.Log) {
	logger = log
	activeDataCollecter = make(map[string]*activeDataUnit)
	activeDataSwitch = make(map[string]bool)
	// arr := splitStr(str)
	for _, key := range str {
		if key != "" {
			activeDataSwitch[key] = true
			activeDataCollecter[key] = &activeDataUnit{
				ids: make(map[int64]bool),
			}
		}
	}

	go func() {
		tick := time.NewTicker(time.Minute * 1)
		for range tick.C {
			WriteToRedis()
		}
	}()
}

func canUse(key string) bool {
	if value, ok := activeDataSwitch[key]; ok && value {
		return true
	}
	return false
}

func AddToActiveDataCollecter(key string, id int64) {
	if !canUse(key) {
		return
	}
	activeDataCollecter[key].mutex.Lock()
	defer activeDataCollecter[key].mutex.Unlock()

	activeDataCollecter[key].ids[id] = true
}

func getOldRedisKeyForGetData(key string, t time.Time) string {
	var redisKey string
	if t.Hour() < 2 { // 凌晨2点前使用前一天的数据
		redisKey = campaignSetPreKey + key + "_" + t.Add(-24*time.Hour).Format("2006-01-02")
	} else {
		redisKey = campaignSetPreKey + key + "_" + t.Format("2006-01-02")
	}
	return redisKey
}

func getRedisKeyByTime(key string, t time.Time) string {
	hourStr := strconv.Itoa(t.Hour())
	minuteStr := strconv.Itoa(t.Minute() / 10)
	return campaignSetPreKey + key + "_" + t.Format("2006-01-02") + "-" + hourStr + "-" + minuteStr
}

func getRedisKey(key string) string {
	return getRedisKeyByTime(key, time.Now())
}

func getRedisKeyArrForGetData(key string) []string {
	maxLength := 6
	redisKeyArr := make([]string, maxLength+1)

	now := time.Now()
	for i := 0; i < maxLength; i++ {
		t := now.Add(time.Duration(-1*(i*10+1)) * time.Minute)
		redisKeyArr[i] = getRedisKeyByTime(key, t)
	}
	redisKeyArr[maxLength] = getOldRedisKeyForGetData(key, now) // 将旧key也加进来，防止切换的时候出现无数据的case

	return redisKeyArr
}

func GetActiveDatas(key string) ([]int64, error) {
	if !canUse(key) {
		return nil, errors.New(key + " is not in config!")
	}
	redisKeyArr := getRedisKeyArrForGetData(key)

	idsMap := make(map[int64]bool)
	for _, redisKey := range redisKeyArr {
		result, err := redis.RedisSMembers(redisKey)
		if err != nil {
			return nil, err
		}
		for _, reStr := range result {
			if k, err := strconv.ParseInt(reStr, 10, 64); err == nil {
				idsMap[k] = true
			}
			if len(idsMap) > maxHotDataNumber {
				goto MaxHotDataTo
			}
		}
	}
MaxHotDataTo:
	ids := make([]int64, len(idsMap))
	i := 0
	for k := range idsMap {
		ids[i] = k
		i++
	}
	return ids, nil
}

func (ad *activeDataUnit) getKeys() []interface{} {
	ad.mutex.Lock()
	defer ad.mutex.Unlock()

	if len(ad.ids) == 0 {
		return nil
	}
	ids := make([]interface{}, len(ad.ids))
	i := 0
	for k := range ad.ids {
		ids[i] = k
		i++
	}
	ad.ids = make(map[int64]bool) // reset

	return ids
}

func WriteToRedis() {
	for key, activeDataObj := range activeDataCollecter {
		ids := activeDataObj.getKeys()
		watchKey := key + "_to_redis"
		if ids == nil {
			watcher.AddWatchValue(watchKey, float64(0))
			continue
		}
		redisKey := getRedisKey(key)

		if err := redis.RedisSAdd(redisKey, ids...); err != nil {
			logger.Errorf("Write %s active data to redis err: %s", key, err.Error())
		}
		if _, err := redis.RedisExpire(redisKey, time.Duration(3600*6)*time.Second); err != nil {
			logger.Errorf("Expire ke %s err: %s", key, err.Error())
		}

		watcher.AddWatchValue(watchKey, float64(len(ids)))
	}
}
