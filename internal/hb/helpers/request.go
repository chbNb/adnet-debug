package helpers

import "gopkg.in/mgo.v2/bson"

func GetRequestID() string {
	return bson.NewObjectId().Hex()
	// return primitive.NewObjectID().Hex()
}

func GetBidId(reqId string) string {
	return SubString(reqId, 0, 18) + SubString(reqId, 19, 5) + "x"
}

func SubString(str string, begin, length int) string {
	rs := []rune(str)
	llen := len(rs)
	if begin < 0 {
		begin = 0
	}
	if begin >= llen {
		begin = llen
	}
	end := begin + length
	if end > llen {
		end = llen
	}
	return string(rs[begin:end])
}
