package filter

import (
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/smartystreets/goconvey/convey"
)

type Message struct {
	Name string
	Body string
	Time int64
}

func TestError(t *testing.T) {
	convey.Convey("filter error code", t, func() {
		cause := errors.New("something error")
		err := errors.WithMessage(New(BidAdxDoBidError), cause.Error())
		t.Logf("%s", err.Error())
		switch errors.Cause(err).(type) {
		case *PipeLineHandlerError:
			t.Log("===== PipeLineHandlerError =====")
			wrapErrMsg, rawErrMsg := FormatErrorMessage(err.Error())
			t.Log(wrapErrMsg)
			code, _ := UnmarshalMessage(rawErrMsg)
			convey.So(code, convey.ShouldEqual, BidAdxDoBidError)
			convey.So(code.String(), convey.ShouldEqual, "bid adx http do error")
		default:
			t.Log("===== default =====")
		}

		convey.So(err, convey.ShouldBeError)

		cause = errors.New("whoops")
		err = errors.Wrap(cause, "oh noes")
		wrapErrMsg, rawErrMsg := FormatErrorMessage(err.Error())
		convey.So(wrapErrMsg, convey.ShouldEqual, "oh noes")
		convey.So(rawErrMsg, convey.ShouldEqual, "whoops")

		wrapErrMsg, rawErrMsg = FormatErrorMessage(cause.Error())
		convey.So(wrapErrMsg, convey.ShouldEqual, "")
		convey.So(rawErrMsg, convey.ShouldEqual, "whoops")

		var msg Message
		b := []byte(`{"Name":"Alice","Body":123,"Time":"1294706395881547000"}`)
		err = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(b, &msg)
		if err != nil {
			err = errors.Wrap(New(BidRequestNot200), err.Error())
			wrapErrMsg, rawErrMsg = FormatErrorMessage(err.Error())
			t.Logf("wrapErrMsg: %s, rawErrMsg: %s", wrapErrMsg, rawErrMsg)
			convey.So(rawErrMsg, convey.ShouldEqual, `{"status":10903,"msg":"bid request is not 200"}`)
			convey.So(wrapErrMsg, convey.ShouldEqual, `filter.Message.Body: ReadString: expects " or n, but found 1, error found in #10 byte of ...|","Body":123,"Time":|..., bigger context ...|{"Name":"Alice","Body":123,"Time":"1294706395881547000"}|...`)
		}

		// multi error wrap
		if err != nil {
			err = errors.Wrap(err, "wrap error message")
			wrapErrMsg, rawErrMsg = FormatErrorMessage(err.Error())
			convey.So(wrapErrMsg, convey.ShouldEqual, `wrap error message: filter.Message.Body: ReadString: expects " or n, but found 1, error found in #10 byte of ...|","Body":123,"Time":|..., bigger context ...|{"Name":"Alice","Body":123,"Time":"1294706395881547000"}|...`)
			convey.So(rawErrMsg, convey.ShouldEqual, `{"status":10903,"msg":"bid request is not 200"}`)
		}

		err = AppPlatformError
		wrapErrMsg, rawErrMsg = FormatErrorMessage(err.Error())
		convey.So(wrapErrMsg, convey.ShouldEqual, "")
		convey.So(rawErrMsg, convey.ShouldEqual, `{"status":10611,"msg":"app platform error"}`)

		errMsg := `Post http://adn-adx-internal-sg.rayjump.com/hbrtb, context deadline exceeded: {"status":10905,"msg":"bid adx http do error"}`
		wrapErrMsg, rawErrMsg = FormatErrorMessage(errMsg)
		t.Logf("wrapErrMsg: %s, rawErrMsg: %s", wrapErrMsg, rawErrMsg)
		filterCode, err := UnmarshalMessage(rawErrMsg)
		convey.So(err, convey.ShouldBeEmpty)
		convey.So(filterCode.Int(), convey.ShouldEqual, 10905)
		convey.So(filterCode.String(), convey.ShouldEqual, "bid adx http do error")
	})
}
