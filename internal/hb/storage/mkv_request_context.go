package storage

import (
	"bytes"
	"compress/gzip"
	"errors"
	jsoniter "github.com/json-iterator/go"
	"gitlab.mobvista.com/ADN/adnet/internal/corsair_proto"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/exporter/metrics"
	"io/ioutil"
	"strings"
	"time"
)

type ReqCtxKey struct {
	Token string `json:"Token,omitempty"`
}

type ReqCtxCompressor struct{}

func NewReqCtxCompressor() *ReqCtxCompressor {
	return &ReqCtxCompressor{}
}

func (c *ReqCtxCompressor) Compress(key interface{}) string {
	return key.(*ReqCtxKey).Token
}

type ReqCtxVal struct {
	ReqParams       *mvutil.RequestParams       `json:"req,omitempty"`
	Result          *corsair_proto.QueryResult_ `json:"res,omitempty"`
	MaxTimeout      int                         `json:"m_t,omitempty"`
	FlowTagID       int                         `json:"ft_id,omitempty"`
	RandValue       int                         `json:"r_v,omitempty"`
	AdsTest         bool                        `json:"a_t,omitempty"`
	Elapsed         int                         `json:"e,omitempty"`
	Finished        int                         `json:"f,omitempty"`
	Backends        map[int]*mvutil.BackendCtx  `json:"b,omitempty"`
	OrderedBackends []int                       `json:"o_b,omitempty"`
	DebugModeInfo   []interface{}               `json:"d,omitempty"`
	IsWhiteGdt      bool                        `json:"gdt,omitempty"`
	IsNativeVideo   bool                        `json:"nv,omitempty"`
}

type ReqCtxSerializer struct{}

func NewReqCtxSerializer() *ReqCtxSerializer {
	return &ReqCtxSerializer{}
}

func (s *ReqCtxSerializer) Marshal(val interface{}) (buf []byte, err error) {
	// return jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(&val)
	data, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(&val)
	if err != nil {
		return nil, err
	}

	// 类型断言查看是否开启 gzip 压缩
	vv, ok := val.(*ReqCtxVal)
	if !ok {
		/* not gzip */
		return data, nil
	}

	// 需要Gzip压缩
	if vv.ReqParams.Param.ExtDataInit.AerospikeGzipEnable == 1 {
		now := time.Now()
		/* gzip */
		gzipBuf := bytes.NewBuffer(buf)
		gzipWriter := gzip.NewWriter(gzipBuf)

		/* defer */
		defer func() {
			if gzipWriter != nil {
				_ = gzipWriter.Close()
			}
		}()
		_, err = gzipWriter.Write(data)
		if err != nil {
			return nil, errors.New(strings.Join([]string{"Gzip marshal write error: ", err.Error()}, ""))
		}
		err = gzipWriter.Flush()
		if err != nil {
			return nil, errors.New(strings.Join([]string{"Gzip marshal flush error:", err.Error()}, " "))
		}
		// 统计使用了 gzip 序列化的请求次数 & 耗时
		costTime := (float64)(time.Since(now) / time.Microsecond)
		metrics.IncCounterWithLabelValues(27, "gzip")
		metrics.SetGaugeWithLabelValues(costTime, 28, "gzip")
		//
		return gzipBuf.Bytes(), nil
	} else {
		/* not gzip */
		return data, nil
	}
}

func (s *ReqCtxSerializer) Unmarshal(buf []byte, val interface{}) (err error) {
	// return jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(buf, &val)
	now := time.Now()
	gzipBuf := bytes.NewReader(buf)
	gzipReader, err := gzip.NewReader(gzipBuf)

	/* defer */
	defer func() {
		if gzipReader != nil {
			_ = gzipReader.Close()
		}
	}()

	/* read compatible ether gzip or not giZip data*/
	if err != nil {
		/* unGzip error, treated as old data*/
		return jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(buf, &val)
	} else {
		data, err := ioutil.ReadAll(gzipReader)
		/* normal data will also throw 'unexpected EOF' error */
		if err != nil && !strings.Contains(err.Error(), "unexpected EOF") {
			return errors.New(strings.Join([]string{"Unmarshal gzip data error: ", err.Error()}, ""))
		}

		// 统计使用gzip反序列化的次数 & 耗时
		costTime := (float64)(time.Since(now) / time.Microsecond)
		metrics.IncCounterWithLabelValues(27, "unGzip")
		metrics.SetGaugeWithLabelValues(costTime, 28, "unGzip")

		return jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(data, &val)
	}
}
