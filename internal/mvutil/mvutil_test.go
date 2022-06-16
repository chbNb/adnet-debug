package mvutil

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	smodel "gitlab.mobvista.com/ADN/structs/model"
)

func TestSliceIndex(t *testing.T) {
	Convey("SliceIndex 返回 3", t, func() {
		res := SliceIndex(10, func(i int) bool {
			return i == 3
		})
		So(res, ShouldEqual, 3)
	})

	Convey("SliceIndex 返回 -1", t, func() {
		res := SliceIndex(3, func(i int) bool {
			return false
		})
		So(res, ShouldEqual, -1)
	})
}

func TestSha1(t *testing.T) {
	Convey("Sha1", t, func() {
		res := Sha1("test_string")
		So(res, ShouldEqual, "c58efadcf9f6b303e44e4a62dd984bbc8cae6e99")
	})
}

func TestMd5(t *testing.T) {
	Convey("Md5", t, func() {
		res := Md5("test_string")
		So(res, ShouldEqual, "3474851a3410906697ec77337df7aae4")
	})
}

func TestBase64(t *testing.T) {
	Convey("Base 64", t, func() {
		byteData := []byte("test_string")
		res := Base64(byteData)
		So(res, ShouldEqual, "dGVzdF9zdHJpbmc=")
	})
}

func TestNumFormat(t *testing.T) {
	Convey("NumFormat 截取一位小数", t, func() {
		res := NumFormat(42.42, 1)
		So(res, ShouldEqual, 42.4)
	})

	Convey("NumFormat 截取两位小数", t, func() {
		res := NumFormat(0.12345, 2)
		So(res, ShouldEqual, 0.12)
	})

	Convey("NumFormat 截取两位小数", t, func() {
		res := NumFormat(42, 2)
		So(res, ShouldEqual, 42)
	})
}

func TestBase64Encode(t *testing.T) {
	Convey("转换码表后的 base 64", t, func() {
		res := Base64Encode("test_string")
		So(res, ShouldEqual, "J75AJcKAJdzuYrh=")
	})
}

func TestUrlDecodeReplace(t *testing.T) {
	data := "B30Dkxqej7rgf7ZAuYvoSxacug02kooDGzvg5IwPlvGtFLPw979JEoat8Y97uoqH5POg979LTIF%2FxIVLF8vI979JRd8tqH08xRB1kvmLlvV2GH0gxRFi51R1qAe%2F5%2BGk5HGFGsmGx8vfxIZdTY9FuA8BGLajqotdTgVOkgH%2F7HF%2F8PWL8Izx5oVqkLm7RzmxR3mj5oCGFAqzlvCRTgm5BxPWG%2B8%2FGoD1mPF68v0LkYGimvF1T%2BBoui9UEPFimowPlcR1GsZdTi9kusVRxRv7fRqgf7VFEsCnkgmdTIwcT%2B925PaO979J5sGIEA9i7PVRFxfM5gOzkLmn5AVxBov5F%2B8Nk78zFYasTLfU979Jm8V2Bxt17L92RxVGmIqzqi9P7JR1qHf6xHGOxIGGRxv2RxDWGYGOk78GR3Gjf8fcGguzEYXgRNrouH0Cf399jxBH5P9o5LW1FTR1qPC0RxGAki97GLmA5gGVR885q8v5T%2B9CR%2BRglN0OqxfNxx92lcR1qHGR7xW4FxvNfg8TfLv7jxBH5P9GGYXH5P92FLVRTNHFB783l%2B82kPCslLm2R7868cR1qPC3fxmn7Iw6kA9NfoXwlvv%2FmLBzTYW38PtH5P9O8NZdk8Bz8RCLRIF38sCjlIvO5zPzTz03kgRdkxOzFLH8lHG6FYHd8HFO5YaUTI9Ff7fHT%2Br1l%2BqgTN9AkAFP8i8njRtzBYGxqxwd7H82uxH4RxGC8sqH5PFxkPzcEHGXFxVxR3f6R3C1TPFixImUliHx979JlNf28NGiGvG2FgGUmRz3f8fRGv8A8xFURNfMRxfi8imnfzFsTAG65g8TFAmz5RZH5P9G8RCsRPzNf3VskzFiRi8VTImMq851xxGzRPwcTPuzuPVzlJR1qAV1GL97979Ll8FGmRG6FAWCksGvTs9qRg8cRPFUBRW17JR1qAFzTLI1qHP1RPv68ICx7%2B84qiqH5P97uAa8Rz0jqi9k7Lt1BHOD7705FNVRkLGC5TR1qAe%2FkvmcRHG3xHGdkIBWmzHsRzVOEPtH5PFABHOH5PFOBPVRk8Cq88HL5AV3BHfgGJR1qifdGJR1qHffkIO15zHcTcR1GHPdTdGOEgOgTLaOmICUF8qH5P9QTgfFBPtH5PFFfoaoTLffxIGs5dm5kdVt8zGifYHRlJR1qPW1f7HtGPWfTs9Gkom6TvvAf8Fq8g8AjLMzTYaA979Jm7PzkgODGvIg8LMdTIzkGACHkRFMF%2BG8FAXo979JTIOzusuH5P9xf8vQkL97uYFLl%2Bm65PCBxx9MmxHLRN0UEiRD7L9tF7RdlP96fdtMlLaARoVURRWOBdZH5P96FNRH5PF68xFBxAqzEAH1RxWAuP8J8i9iEYXgRL8xR3Fdk1R1qHv0lRFA78H85sC5kYC08Pv7RPwPTJR1GPMzTP9ARxwRls9G979L7JR1qsV8TPFq5oGsRovWEACLkxWYuHHU8IaCqxC1xxXWTYVLj3mFFzfBTLW2979LFA9TuHPH5P9CToqH5P9xuHPH5P9jTYH68PvCq89LFAW4jRmkTNm2879dlAm7udH6RPznuiRWk%2BmGqPV68L8AEAH0FAvCjxHV8H0nuHvflxW4qRmkj79NxiHJkRvt979Lk88Yj3C1RRzIuHvRkIvR73CsFPzI979JlsGF979JkRaO979JmIz5979JlHv3fzvfGN9G8LVJF38tEY8HGg0Hm8H1R8GIBPtH5PFOfdCfkI9HmzHkF3mF5L8dkI9HmotDFRzFEoV07La2979LGJR1qsZ1kN9xTdV6Ri9CBRq1GYv7BYPgfxCAB3VVk7f4GPm6RPzAkHFck88TRLF4kx9NEsqd7HGOBxF4FHfNxiHHlLWtRLF0GL8CB8FLFAWYqAH8lLeMksZdkxCABYPgfxeCBxw0FH8IRIW0kIFNBxVRGYWx5sC1mzIoEi9oxxWC5JR1qHCf5IV1770N5v9HkRvjGPmVEI9CqH9gRs9dBH98G70MGPaBTIFHR8OD770tEzOwRi0nT1R1qPMH5oqH5oqAudHgmL8aS7IAEdfpmio1fJFokxzHlAe6F7zN78qH5POD51Rgq7rD9AGCuA8imLHQEiow9AvDuveYBgoof%2B5o5iID9AvDuveDEAvaF7ziEYo6uYWVuNGzBAR6uN9HExHzETF6F3GdEd9US7PAkxzVS7GJls8Fkg8zRH8v5vOD770FkYf8RH8v5vOD770IkzrgkN9%2FmRViTd8OEzOwRH8jBPGU8HGNfoZDT%2BmIjxHLkPFCfoMDfPMH5oqH5oqAExWiS7rw9sfiuA8HEHegk3CHS7IDj%2B0M57P15JFgFLaYF39gkxe6SRzq3g56f1MD9A91BxWPS3fVE3fzEAuAmxIe7xenkxwXBTR1GiR65JXH5iV5kxWzlJRgqca0EAG1EYHP2gu65JRgqca77TzNfiIDGcaJmxHXFJR1GPWTG%2BPD7TRgqcadmcR1jTa0uN0XF8mHBPaCmJR1GiRgf1MgfcXH5iV%2FTvGf7JR1q1aXkxaH2omHBYaQ97OW2zFHusfCEYMH5PBo%2FirUqYV1EYzH979Lfiu65JMg5gPY%2FiZd2ozQBAHXFTa7BxFVuAPH5PBz5gu65gBAELv6Fd8VFYReB3OABxGCF%2BzAj78VFAqo5czV5g0V%2F7qzj%2BrajLOwfco1BAOD5i5WFA8iF%2BPAExfiS7OMfcFVuN0pmAMefJMo51MD%2Fiqo5gq157rAExePFxDeuYvaud86F1a77TzNfiIDGZ%3D%3D"
	Convey("urlDecode should be equal", t, func() {
		dataReplace := UrlDecodeReplace(data)
		dataDecode, err := url.QueryUnescape(data)
		So(err, ShouldBeNil)
		So(dataReplace, ShouldEqual, dataDecode)
	})
}

func BenchmarkUrlDecode(b *testing.B) {
	data := "B30Dkxqej7rgf7ZAuYvoSxacug02kooDGzvg5IwPlvGtFLPw979JEoat8Y97uoqH5POg979LTIF%2FxIVLF8vI979JRd8tqH08xRB1kvmLlvV2GH0gxRFi51R1qAe%2F5%2BGk5HGFGsmGx8vfxIZdTY9FuA8BGLajqotdTgVOkgH%2F7HF%2F8PWL8Izx5oVqkLm7RzmxR3mj5oCGFAqzlvCRTgm5BxPWG%2B8%2FGoD1mPF68v0LkYGimvF1T%2BBoui9UEPFimowPlcR1GsZdTi9kusVRxRv7fRqgf7VFEsCnkgmdTIwcT%2B925PaO979J5sGIEA9i7PVRFxfM5gOzkLmn5AVxBov5F%2B8Nk78zFYasTLfU979Jm8V2Bxt17L92RxVGmIqzqi9P7JR1qHf6xHGOxIGGRxv2RxDWGYGOk78GR3Gjf8fcGguzEYXgRNrouH0Cf399jxBH5P9o5LW1FTR1qPC0RxGAki97GLmA5gGVR885q8v5T%2B9CR%2BRglN0OqxfNxx92lcR1qHGR7xW4FxvNfg8TfLv7jxBH5P9GGYXH5P92FLVRTNHFB783l%2B82kPCslLm2R7868cR1qPC3fxmn7Iw6kA9NfoXwlvv%2FmLBzTYW38PtH5P9O8NZdk8Bz8RCLRIF38sCjlIvO5zPzTz03kgRdkxOzFLH8lHG6FYHd8HFO5YaUTI9Ff7fHT%2Br1l%2BqgTN9AkAFP8i8njRtzBYGxqxwd7H82uxH4RxGC8sqH5PFxkPzcEHGXFxVxR3f6R3C1TPFixImUliHx979JlNf28NGiGvG2FgGUmRz3f8fRGv8A8xFURNfMRxfi8imnfzFsTAG65g8TFAmz5RZH5P9G8RCsRPzNf3VskzFiRi8VTImMq851xxGzRPwcTPuzuPVzlJR1qAV1GL97979Ll8FGmRG6FAWCksGvTs9qRg8cRPFUBRW17JR1qAFzTLI1qHP1RPv68ICx7%2B84qiqH5P97uAa8Rz0jqi9k7Lt1BHOD7705FNVRkLGC5TR1qAe%2FkvmcRHG3xHGdkIBWmzHsRzVOEPtH5PFABHOH5PFOBPVRk8Cq88HL5AV3BHfgGJR1qifdGJR1qHffkIO15zHcTcR1GHPdTdGOEgOgTLaOmICUF8qH5P9QTgfFBPtH5PFFfoaoTLffxIGs5dm5kdVt8zGifYHRlJR1qPW1f7HtGPWfTs9Gkom6TvvAf8Fq8g8AjLMzTYaA979Jm7PzkgODGvIg8LMdTIzkGACHkRFMF%2BG8FAXo979JTIOzusuH5P9xf8vQkL97uYFLl%2Bm65PCBxx9MmxHLRN0UEiRD7L9tF7RdlP96fdtMlLaARoVURRWOBdZH5P96FNRH5PF68xFBxAqzEAH1RxWAuP8J8i9iEYXgRL8xR3Fdk1R1qHv0lRFA78H85sC5kYC08Pv7RPwPTJR1GPMzTP9ARxwRls9G979L7JR1qsV8TPFq5oGsRovWEACLkxWYuHHU8IaCqxC1xxXWTYVLj3mFFzfBTLW2979LFA9TuHPH5P9CToqH5P9xuHPH5P9jTYH68PvCq89LFAW4jRmkTNm2879dlAm7udH6RPznuiRWk%2BmGqPV68L8AEAH0FAvCjxHV8H0nuHvflxW4qRmkj79NxiHJkRvt979Lk88Yj3C1RRzIuHvRkIvR73CsFPzI979JlsGF979JkRaO979JmIz5979JlHv3fzvfGN9G8LVJF38tEY8HGg0Hm8H1R8GIBPtH5PFOfdCfkI9HmzHkF3mF5L8dkI9HmotDFRzFEoV07La2979LGJR1qsZ1kN9xTdV6Ri9CBRq1GYv7BYPgfxCAB3VVk7f4GPm6RPzAkHFck88TRLF4kx9NEsqd7HGOBxF4FHfNxiHHlLWtRLF0GL8CB8FLFAWYqAH8lLeMksZdkxCABYPgfxeCBxw0FH8IRIW0kIFNBxVRGYWx5sC1mzIoEi9oxxWC5JR1qHCf5IV1770N5v9HkRvjGPmVEI9CqH9gRs9dBH98G70MGPaBTIFHR8OD770tEzOwRi0nT1R1qPMH5oqH5oqAudHgmL8aS7IAEdfpmio1fJFokxzHlAe6F7zN78qH5POD51Rgq7rD9AGCuA8imLHQEiow9AvDuveYBgoof%2B5o5iID9AvDuveDEAvaF7ziEYo6uYWVuNGzBAR6uN9HExHzETF6F3GdEd9US7PAkxzVS7GJls8Fkg8zRH8v5vOD770FkYf8RH8v5vOD770IkzrgkN9%2FmRViTd8OEzOwRH8jBPGU8HGNfoZDT%2BmIjxHLkPFCfoMDfPMH5oqH5oqAExWiS7rw9sfiuA8HEHegk3CHS7IDj%2B0M57P15JFgFLaYF39gkxe6SRzq3g56f1MD9A91BxWPS3fVE3fzEAuAmxIe7xenkxwXBTR1GiR65JXH5iV5kxWzlJRgqca0EAG1EYHP2gu65JRgqca77TzNfiIDGcaJmxHXFJR1GPWTG%2BPD7TRgqcadmcR1jTa0uN0XF8mHBPaCmJR1GiRgf1MgfcXH5iV%2FTvGf7JR1q1aXkxaH2omHBYaQ97OW2zFHusfCEYMH5PBo%2FirUqYV1EYzH979Lfiu65JMg5gPY%2FiZd2ozQBAHXFTa7BxFVuAPH5PBz5gu65gBAELv6Fd8VFYReB3OABxGCF%2BzAj78VFAqo5czV5g0V%2F7qzj%2BrajLOwfco1BAOD5i5WFA8iF%2BPAExfiS7OMfcFVuN0pmAMefJMo51MD%2Fiqo5gq157rAExePFxDeuYvaud86F1a77TzNfiIDGZ%3D%3D"
	b.Run("Replace-Decode", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			data1 := strings.Replace(data, "%3D", "=", -1)
			data1 = strings.Replace(data1, "%2B", "+", -1)
			_ = strings.Replace(data1, "%2F", "/", -1)
		}
	})

	b.Run("Url-Decode", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			url.QueryUnescape(data)
		}
	})
}

func TestBase64Decode(t *testing.T) {
	Convey("转换码表后的 base64decode", t, func() {
		res := Base64Decode("J75AJcKAJdzuYrh=")
		data := "B30Dkxqej7rgf7ZAuYvoSxacug02kooDGzvg5IwPlvGtFLPw979JEoat8Y97uoqH5POg979LTIF%2FxIVLF8vI979JRd8tqH08xRB1kvmLlvV2GH0gxRFi51R1qAe%2F5%2BGk5HGFGsmGx8vfxIZdTY9FuA8BGLajqotdTgVOkgH%2F7HF%2F8PWL8Izx5oVqkLm7RzmxR3mj5oCGFAqzlvCRTgm5BxPWG%2B8%2FGoD1mPF68v0LkYGimvF1T%2BBoui9UEPFimowPlcR1GsZdTi9kusVRxRv7fRqgf7VFEsCnkgmdTIwcT%2B925PaO979J5sGIEA9i7PVRFxfM5gOzkLmn5AVxBov5F%2B8Nk78zFYasTLfU979Jm8V2Bxt17L92RxVGmIqzqi9P7JR1qHf6xHGOxIGGRxv2RxDWGYGOk78GR3Gjf8fcGguzEYXgRNrouH0Cf399jxBH5P9o5LW1FTR1qPC0RxGAki97GLmA5gGVR885q8v5T%2B9CR%2BRglN0OqxfNxx92lcR1qHGR7xW4FxvNfg8TfLv7jxBH5P9GGYXH5P92FLVRTNHFB783l%2B82kPCslLm2R7868cR1qPC3fxmn7Iw6kA9NfoXwlvv%2FmLBzTYW38PtH5P9O8NZdk8Bz8RCLRIF38sCjlIvO5zPzTz03kgRdkxOzFLH8lHG6FYHd8HFO5YaUTI9Ff7fHT%2Br1l%2BqgTN9AkAFP8i8njRtzBYGxqxwd7H82uxH4RxGC8sqH5PFxkPzcEHGXFxVxR3f6R3C1TPFixImUliHx979JlNf28NGiGvG2FgGUmRz3f8fRGv8A8xFURNfMRxfi8imnfzFsTAG65g8TFAmz5RZH5P9G8RCsRPzNf3VskzFiRi8VTImMq851xxGzRPwcTPuzuPVzlJR1qAV1GL97979Ll8FGmRG6FAWCksGvTs9qRg8cRPFUBRW17JR1qAFzTLI1qHP1RPv68ICx7%2B84qiqH5P97uAa8Rz0jqi9k7Lt1BHOD7705FNVRkLGC5TR1qAe%2FkvmcRHG3xHGdkIBWmzHsRzVOEPtH5PFABHOH5PFOBPVRk8Cq88HL5AV3BHfgGJR1qifdGJR1qHffkIO15zHcTcR1GHPdTdGOEgOgTLaOmICUF8qH5P9QTgfFBPtH5PFFfoaoTLffxIGs5dm5kdVt8zGifYHRlJR1qPW1f7HtGPWfTs9Gkom6TvvAf8Fq8g8AjLMzTYaA979Jm7PzkgODGvIg8LMdTIzkGACHkRFMF%2BG8FAXo979JTIOzusuH5P9xf8vQkL97uYFLl%2Bm65PCBxx9MmxHLRN0UEiRD7L9tF7RdlP96fdtMlLaARoVURRWOBdZH5P96FNRH5PF68xFBxAqzEAH1RxWAuP8J8i9iEYXgRL8xR3Fdk1R1qHv0lRFA78H85sC5kYC08Pv7RPwPTJR1GPMzTP9ARxwRls9G979L7JR1qsV8TPFq5oGsRovWEACLkxWYuHHU8IaCqxC1xxXWTYVLj3mFFzfBTLW2979LFA9TuHPH5P9CToqH5P9xuHPH5P9jTYH68PvCq89LFAW4jRmkTNm2879dlAm7udH6RPznuiRWk%2BmGqPV68L8AEAH0FAvCjxHV8H0nuHvflxW4qRmkj79NxiHJkRvt979Lk88Yj3C1RRzIuHvRkIvR73CsFPzI979JlsGF979JkRaO979JmIz5979JlHv3fzvfGN9G8LVJF38tEY8HGg0Hm8H1R8GIBPtH5PFOfdCfkI9HmzHkF3mF5L8dkI9HmotDFRzFEoV07La2979LGJR1qsZ1kN9xTdV6Ri9CBRq1GYv7BYPgfxCAB3VVk7f4GPm6RPzAkHFck88TRLF4kx9NEsqd7HGOBxF4FHfNxiHHlLWtRLF0GL8CB8FLFAWYqAH8lLeMksZdkxCABYPgfxeCBxw0FH8IRIW0kIFNBxVRGYWx5sC1mzIoEi9oxxWC5JR1qHCf5IV1770N5v9HkRvjGPmVEI9CqH9gRs9dBH98G70MGPaBTIFHR8OD770tEzOwRi0nT1R1qPMH5oqH5oqAudHgmL8aS7IAEdfpmio1fJFokxzHlAe6F7zN78qH5POD51Rgq7rD9AGCuA8imLHQEiow9AvDuveYBgoof%2B5o5iID9AvDuveDEAvaF7ziEYo6uYWVuNGzBAR6uN9HExHzETF6F3GdEd9US7PAkxzVS7GJls8Fkg8zRH8v5vOD770FkYf8RH8v5vOD770IkzrgkN9%2FmRViTd8OEzOwRH8jBPGU8HGNfoZDT%2BmIjxHLkPFCfoMDfPMH5oqH5oqAExWiS7rw9sfiuA8HEHegk3CHS7IDj%2B0M57P15JFgFLaYF39gkxe6SRzq3g56f1MD9A91BxWPS3fVE3fzEAuAmxIe7xenkxwXBTR1GiR65JXH5iV5kxWzlJRgqca0EAG1EYHP2gu65JRgqca77TzNfiIDGcaJmxHXFJR1GPWTG%2BPD7TRgqcadmcR1jTa0uN0XF8mHBPaCmJR1GiRgf1MgfcXH5iV%2FTvGf7JR1q1aXkxaH2omHBYaQ97OW2zFHusfCEYMH5PBo%2FirUqYV1EYzH979Lfiu65JMg5gPY%2FiZd2ozQBAHXFTa7BxFVuAPH5PBz5gu65gBAELv6Fd8VFYReB3OABxGCF%2BzAj78VFAqo5czV5g0V%2F7qzj%2BrajLOwfco1BAOD5i5WFA8iF%2BPAExfiS7OMfcFVuN0pmAMefJMo51MD%2Fiqo5gq157rAExePFxDeuYvaud86F1a77TzNfiIDGZ%3D%3D"
		//res2 := Base64DecodeMP("B30Dkxqej7rgf7ZAuYvoSxacug02kooDGzvg5IwPlvGtFLPw979JEoat8Y97uoqH5POg979LTIF%2FxIVLF8vI979JRd8tqH08xRB1kvmLlvV2GH0gxRFi51R1qAe%2F5%2BGk5HGFGsmGx8vfxIZdTY9FuA8BGLajqotdTgVOkgH%2F7HF%2F8PWL8Izx5oVqkLm7RzmxR3mj5oCGFAqzlvCRTgm5BxPWG%2B8%2FGoD1mPF68v0LkYGimvF1T%2BBoui9UEPFimowPlcR1GsZdTi9kusVRxRv7fRqgf7VFEsCnkgmdTIwcT%2B925PaO979J5sGIEA9i7PVRFxfM5gOzkLmn5AVxBov5F%2B8Nk78zFYasTLfU979Jm8V2Bxt17L92RxVGmIqzqi9P7JR1qHf6xHGOxIGGRxv2RxDWGYGOk78GR3Gjf8fcGguzEYXgRNrouH0Cf399jxBH5P9o5LW1FTR1qPC0RxGAki97GLmA5gGVR885q8v5T%2B9CR%2BRglN0OqxfNxx92lcR1qHGR7xW4FxvNfg8TfLv7jxBH5P9GGYXH5P92FLVRTNHFB783l%2B82kPCslLm2R7868cR1qPC3fxmn7Iw6kA9NfoXwlvv%2FmLBzTYW38PtH5P9O8NZdk8Bz8RCLRIF38sCjlIvO5zPzTz03kgRdkxOzFLH8lHG6FYHd8HFO5YaUTI9Ff7fHT%2Br1l%2BqgTN9AkAFP8i8njRtzBYGxqxwd7H82uxH4RxGC8sqH5PFxkPzcEHGXFxVxR3f6R3C1TPFixImUliHx979JlNf28NGiGvG2FgGUmRz3f8fRGv8A8xFURNfMRxfi8imnfzFsTAG65g8TFAmz5RZH5P9G8RCsRPzNf3VskzFiRi8VTImMq851xxGzRPwcTPuzuPVzlJR1qAV1GL97979Ll8FGmRG6FAWCksGvTs9qRg8cRPFUBRW17JR1qAFzTLI1qHP1RPv68ICx7%2B84qiqH5P97uAa8Rz0jqi9k7Lt1BHOD7705FNVRkLGC5TR1qAe%2FkvmcRHG3xHGdkIBWmzHsRzVOEPtH5PFABHOH5PFOBPVRk8Cq88HL5AV3BHfgGJR1qifdGJR1qHffkIO15zHcTcR1GHPdTdGOEgOgTLaOmICUF8qH5P9QTgfFBPtH5PFFfoaoTLffxIGs5dm5kdVt8zGifYHRlJR1qPW1f7HtGPWfTs9Gkom6TvvAf8Fq8g8AjLMzTYaA979Jm7PzkgODGvIg8LMdTIzkGACHkRFMF%2BG8FAXo979JTIOzusuH5P9xf8vQkL97uYFLl%2Bm65PCBxx9MmxHLRN0UEiRD7L9tF7RdlP96fdtMlLaARoVURRWOBdZH5P96FNRH5PF68xFBxAqzEAH1RxWAuP8J8i9iEYXgRL8xR3Fdk1R1qHv0lRFA78H85sC5kYC08Pv7RPwPTJR1GPMzTP9ARxwRls9G979L7JR1qsV8TPFq5oGsRovWEACLkxWYuHHU8IaCqxC1xxXWTYVLj3mFFzfBTLW2979LFA9TuHPH5P9CToqH5P9xuHPH5P9jTYH68PvCq89LFAW4jRmkTNm2879dlAm7udH6RPznuiRWk%2BmGqPV68L8AEAH0FAvCjxHV8H0nuHvflxW4qRmkj79NxiHJkRvt979Lk88Yj3C1RRzIuHvRkIvR73CsFPzI979JlsGF979JkRaO979JmIz5979JlHv3fzvfGN9G8LVJF38tEY8HGg0Hm8H1R8GIBPtH5PFOfdCfkI9HmzHkF3mF5L8dkI9HmotDFRzFEoV07La2979LGJR1qsZ1kN9xTdV6Ri9CBRq1GYv7BYPgfxCAB3VVk7f4GPm6RPzAkHFck88TRLF4kx9NEsqd7HGOBxF4FHfNxiHHlLWtRLF0GL8CB8FLFAWYqAH8lLeMksZdkxCABYPgfxeCBxw0FH8IRIW0kIFNBxVRGYWx5sC1mzIoEi9oxxWC5JR1qHCf5IV1770N5v9HkRvjGPmVEI9CqH9gRs9dBH98G70MGPaBTIFHR8OD770tEzOwRi0nT1R1qPMH5oqH5oqAudHgmL8aS7IAEdfpmio1fJFokxzHlAe6F7zN78qH5POD51Rgq7rD9AGCuA8imLHQEiow9AvDuveYBgoof%2B5o5iID9AvDuveDEAvaF7ziEYo6uYWVuNGzBAR6uN9HExHzETF6F3GdEd9US7PAkxzVS7GJls8Fkg8zRH8v5vOD770FkYf8RH8v5vOD770IkzrgkN9%2FmRViTd8OEzOwRH8jBPGU8HGNfoZDT%2BmIjxHLkPFCfoMDfPMH5oqH5oqAExWiS7rw9sfiuA8HEHegk3CHS7IDj%2B0M57P15JFgFLaYF39gkxe6SRzq3g56f1MD9A91BxWPS3fVE3fzEAuAmxIe7xenkxwXBTR1GiR65JXH5iV5kxWzlJRgqca0EAG1EYHP2gu65JRgqca77TzNfiIDGcaJmxHXFJR1GPWTG%2BPD7TRgqcadmcR1jTa0uN0XF8mHBPaCmJR1GiRgf1MgfcXH5iV%2FTvGf7JR1q1aXkxaH2omHBYaQ97OW2zFHusfCEYMH5PBo%2FirUqYV1EYzH979Lfiu65JMg5gPY%2FiZd2ozQBAHXFTa7BxFVuAPH5PBz5gu65gBAELv6Fd8VFYReB3OABxGCF%2BzAj78VFAqo5czV5g0V%2F7qzj%2BrajLOwfco1BAOD5i5WFA8iF%2BPAExfiS7OMfcFVuN0pmAMefJMo51MD%2Fiqo5gq157rAExePFxDeuYvaud86F1a77TzNfiIDGZ%3D%3D")
		data, err := url.QueryUnescape(data)
		if err != nil {
			t.Log(err)
		}
		//t.Log(data)
		res3 := Base64DecodeMP(data)
		//test := Base64DecodeMP("B30Dkxqej7rgf7ZAuYvoSxacug02kooDGzvg5IwPlvGtFLPw979JEoat8Y97uoqH5POg979LTIF")
		//res3 := Base64DecodeMPOld("B30Dkxqej7rgf7ZAuYvoSxacug02kooDGzvg5IwPlvGtFLPw979JEoat8Y97uoqH5POg979LTIF%2FxIVLF8vI979JRd8tqH08xRB1kvmLlvV2GH0gxRFi51R1qAe%2F5%2BGk5HGFGsmGx8vfxIZdTY9FuA8BGLajqotdTgVOkgH%2F7HF%2F8PWL8Izx5oVqkLm7RzmxR3mj5oCGFAqzlvCRTgm5BxPWG%2B8%2FGoD1mPF68v0LkYGimvF1T%2BBoui9UEPFimowPlcR1GsZdTi9kusVRxRv7fRqgf7VFEsCnkgmdTIwcT%2B925PaO979J5sGIEA9i7PVRFxfM5gOzkLmn5AVxBov5F%2B8Nk78zFYasTLfU979Jm8V2Bxt17L92RxVGmIqzqi9P7JR1qHf6xHGOxIGGRxv2RxDWGYGOk78GR3Gjf8fcGguzEYXgRNrouH0Cf399jxBH5P9o5LW1FTR1qPC0RxGAki97GLmA5gGVR885q8v5T%2B9CR%2BRglN0OqxfNxx92lcR1qHGR7xW4FxvNfg8TfLv7jxBH5P9GGYXH5P92FLVRTNHFB783l%2B82kPCslLm2R7868cR1qPC3fxmn7Iw6kA9NfoXwlvv%2FmLBzTYW38PtH5P9O8NZdk8Bz8RCLRIF38sCjlIvO5zPzTz03kgRdkxOzFLH8lHG6FYHd8HFO5YaUTI9Ff7fHT%2Br1l%2BqgTN9AkAFP8i8njRtzBYGxqxwd7H82uxH4RxGC8sqH5PFxkPzcEHGXFxVxR3f6R3C1TPFixImUliHx979JlNf28NGiGvG2FgGUmRz3f8fRGv8A8xFURNfMRxfi8imnfzFsTAG65g8TFAmz5RZH5P9G8RCsRPzNf3VskzFiRi8VTImMq851xxGzRPwcTPuzuPVzlJR1qAV1GL97979Ll8FGmRG6FAWCksGvTs9qRg8cRPFUBRW17JR1qAFzTLI1qHP1RPv68ICx7%2B84qiqH5P97uAa8Rz0jqi9k7Lt1BHOD7705FNVRkLGC5TR1qAe%2FkvmcRHG3xHGdkIBWmzHsRzVOEPtH5PFABHOH5PFOBPVRk8Cq88HL5AV3BHfgGJR1qifdGJR1qHffkIO15zHcTcR1GHPdTdGOEgOgTLaOmICUF8qH5P9QTgfFBPtH5PFFfoaoTLffxIGs5dm5kdVt8zGifYHRlJR1qPW1f7HtGPWfTs9Gkom6TvvAf8Fq8g8AjLMzTYaA979Jm7PzkgODGvIg8LMdTIzkGACHkRFMF%2BG8FAXo979JTIOzusuH5P9xf8vQkL97uYFLl%2Bm65PCBxx9MmxHLRN0UEiRD7L9tF7RdlP96fdtMlLaARoVURRWOBdZH5P96FNRH5PF68xFBxAqzEAH1RxWAuP8J8i9iEYXgRL8xR3Fdk1R1qHv0lRFA78H85sC5kYC08Pv7RPwPTJR1GPMzTP9ARxwRls9G979L7JR1qsV8TPFq5oGsRovWEACLkxWYuHHU8IaCqxC1xxXWTYVLj3mFFzfBTLW2979LFA9TuHPH5P9CToqH5P9xuHPH5P9jTYH68PvCq89LFAW4jRmkTNm2879dlAm7udH6RPznuiRWk%2BmGqPV68L8AEAH0FAvCjxHV8H0nuHvflxW4qRmkj79NxiHJkRvt979Lk88Yj3C1RRzIuHvRkIvR73CsFPzI979JlsGF979JkRaO979JmIz5979JlHv3fzvfGN9G8LVJF38tEY8HGg0Hm8H1R8GIBPtH5PFOfdCfkI9HmzHkF3mF5L8dkI9HmotDFRzFEoV07La2979LGJR1qsZ1kN9xTdV6Ri9CBRq1GYv7BYPgfxCAB3VVk7f4GPm6RPzAkHFck88TRLF4kx9NEsqd7HGOBxF4FHfNxiHHlLWtRLF0GL8CB8FLFAWYqAH8lLeMksZdkxCABYPgfxeCBxw0FH8IRIW0kIFNBxVRGYWx5sC1mzIoEi9oxxWC5JR1qHCf5IV1770N5v9HkRvjGPmVEI9CqH9gRs9dBH98G70MGPaBTIFHR8OD770tEzOwRi0nT1R1qPMH5oqH5oqAudHgmL8aS7IAEdfpmio1fJFokxzHlAe6F7zN78qH5POD51Rgq7rD9AGCuA8imLHQEiow9AvDuveYBgoof%2B5o5iID9AvDuveDEAvaF7ziEYo6uYWVuNGzBAR6uN9HExHzETF6F3GdEd9US7PAkxzVS7GJls8Fkg8zRH8v5vOD770FkYf8RH8v5vOD770IkzrgkN9%2FmRViTd8OEzOwRH8jBPGU8HGNfoZDT%2BmIjxHLkPFCfoMDfPMH5oqH5oqAExWiS7rw9sfiuA8HEHegk3CHS7IDj%2B0M57P15JFgFLaYF39gkxe6SRzq3g56f1MD9A91BxWPS3fVE3fzEAuAmxIe7xenkxwXBTR1GiR65JXH5iV5kxWzlJRgqca0EAG1EYHP2gu65JRgqca77TzNfiIDGcaJmxHXFJR1GPWTG%2BPD7TRgqcadmcR1jTa0uN0XF8mHBPaCmJR1GiRgf1MgfcXH5iV%2FTvGf7JR1q1aXkxaH2omHBYaQ97OW2zFHusfCEYMH5PBo%2FirUqYV1EYzH979Lfiu65JMg5gPY%2FiZd2ozQBAHXFTa7BxFVuAPH5PBz5gu65gBAELv6Fd8VFYReB3OABxGCF%2BzAj78VFAqo5czV5g0V%2F7qzj%2BrajLOwfco1BAOD5i5WFA8iF%2BPAExfiS7OMfcFVuN0pmAMefJMo51MD%2Fiqo5gq157rAExePFxDeuYvaud86F1a77TzNfiIDGZ%3D%3D")
		//t.Error(test)
		So(res, ShouldEqual, "test_string")
		//So(res2, ShouldEqual, res)
		So(res3, ShouldEqual, "appid=90358&sat=kbs0JkM0GQs0LdxThdi1%2BoKhWbSsD%2B3%2FHFKXHFeQD%2BSuhBPUYF2hWFxXJFPsYFc3%2BoK04Z2TYFwQYQMXH7KbYreXDkNCJ7K8Hk9KNVKVNFTMV3HPhgSSWVQwN3JQfd5xZTK7Lai9D5KGL2vFnTPFkdctVrH64r2knFcwLdz%2Fx7J2ZrxTYAS5D358Ynzzk7wHLbH2J2KH%2B2tDnbcNHTecx325hgz2hVcALd5Gi5ugkgHck%2BuXJaj2LbJQhQtD5B2dL%2BSnZTHXDQQaJQl9GdHi5QQtN5SbG75ok3Pp4rPi5rI9f%2Bt0nre%2BJAQdfj2SDgf34aQULAQLH2iP53xpHAcGYbJz%2BTTMnjeaG75R4aS9f%2BQGk%2BJdhTHyYa5Wx5JjJgxgJQ5nV%2BJW5gzLLnjbG7K1xQKtf5KnWVJ%2BHTx7iV5UJFPFWVzNxAH3Y5KPWk57ib5diUzTngiwVVH3kkHBY53eH02x43HrfjfdV5z9J5cdVAlwNUJqijQdiVt%2FVjMbnTlehVQsnQzrJFcXGkz9V%2BxsJTtcDTJg4kuMW5STDUfUfkPsxQccV7z7VgJdn35Rfgu1H%2BQUJgRMG5xgkVcR5aHGxAS2YduRLbJG5rHux%2BhrDbS%2FyVQuDnfnijtEJrPS5bRFkaNrL%2BfuHa2BY2RAnTJVL5jB4%2BSrkUSPNB2ZLj2bR0M0LdxThdi1%2BoKhWbRTWZTwhF9wYgSXHnJ%2FfbR%2FHbHTiZPUYF2hWbSsD%2B3wD%2BSMhB23YbJ%2FY7KtHo23HkHtJkeT%2BoK3YbJ%2FY7KtHcMXDg3wLkxhWTc7iTx%2BNr59hFNMJrQkGnHQf5VPW5f8n5Kkf%2Bu95k20DQ3Tn7HMZFjeiFxd4Ufk4%2BHB5rw%2BV5QohbSsfFx7n2JXYbxuiFPpkn50Lbhe57zBn7z8xkfSHkQNHcx%2Bndu%2FnUfXZd5nirQnfrEBV2cok3PeVQvwk%2BQAyFfMYU2zLkjAVASRLdH%2FN5JBfQlTzrQ%2FL%2BxUJFP3DgSAynjFinvrYkTKiAjrYk9KhF9wYgSXHnJ%2FfbRrY%2BiKD%2BVrY%2BNKinVAiARFfnj9GZHwJU2wzgSsynRMzr59h7QBHnTefniAfai9iaVPzrQMynjAGZ92GZ9BiAh%2FiUv9zrQMDrQThATMzgfMD%2BztY%2BiKH%2BtML%2BzQW7QMDrQThBeuhoeeG0euYrQTDbJ%2FH7zMhBewYZewY0ewhBewJ0eMYoHALkJ%2FD%2Bx2hrVKxnR2iaD2GaSci35jfaxai3jFGnRMfjVbiURPfjibGnt7NTHafjfSGZ9exnhPfADeiaVFfnvBiUxoxjx7ijfci35oialAfUDPNAhFGahTGnV2zrwQ4n2tYni0%2BZM0HrM0G0ReiANFGalBiBRsRrwbRUE0xFKXHFeQR0M0hoR1R0zK%2BN%3D%3D&system=1&os_v=24&timezone=GMT%2B03%3A00&direction=1&app_vc=4434210&app_pname=com.snaptube.premium&network=9&ima=4BzuYk5uRUE0R0M0YkcURUE0R0M0DkP3hrKuHcKuHoR1RUNbDkVTG7H0H7D9iFjFi7N06N%3D%3D&mnc=01&screen_size=1080x1920&sdkversion=MP_3.7.0&brand=samsung&ua=Mozilla%2F5.0+%28Linux%3B+Android+7.0%3B+SM-G610F+Build%2FNRD90M%3B+wv%29+AppleWebKit%2F537.36+%28KHTML%2C+like+Gecko%29+Version%2F4.0+Chrome%2F67.0.3396.87+Mobile+Safari%2F537.36&language=ar&adid=f95afd42-a30a-4580-8b16-2bb0239fecd9&mcc=286&app_vn=4.43.0.4434210&model=samsung+SM-G610F")
	})
}

func BenchmarkBase64Decode(b *testing.B) {
	data := "J75AJcKAJdzuYrh="
	b.Run("base64Decode", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Base64Decode(data)
		}
	})
	//b.Run("base64DecodeOld", func(b *testing.B) {
	//	for i := 0; i < b.N; i++ {
	//		Base64DecodeOld(data)
	//	}
	//})
}

func TestBase64EncodeMPNew(t *testing.T) {
	Convey("Base64 EncodeMPNew", t, func() {
		res := Base64EncodeMPNew("hello,world!")
		So(res, ShouldEqual, "OJUq9JgqzPoZ9JNW")
	})
}

func TestBase64DecodeMPNew(t *testing.T) {
	Convey("Base64 DecodeMPNew", t, func() {
		res := Base64DecodeMPNew("OJUq9JgqzPoZ9JNW")
		So(res, ShouldEqual, "hello,world!")
	})
}

func TestBase64EncodeMP(t *testing.T) {
	Convey("转换码表后的 base64EncodeMP", t, func() {
		res := Base64EncodeMP("hello,world!")
		So(res, ShouldEqual, "kL8XELbXmYe1ELqV")
	})
}

func TestBase64EncodeMV(t *testing.T) {
	Convey("转换码表后的 base64EncodeMV", t, func() {
		res := Base64EncodeMV("hello,world!")
		So(res, ShouldEqual, "Q3HzY3Wzfsx2Y3NO")
	})
}

func TestBase64DecodeMV(t *testing.T) {
	Convey("转换码表后的 base64DecodeMV", t, func() {
		res := Base64DecodeMV("Q3HzY3Wzfsx2Y3NO")
		So(res, ShouldEqual, "hello,world!")
	})
}

func TestBase64DecodeMP(t *testing.T) {
	Convey("转换码表后的 base64DecodeMP", t, func() {
		res := Base64DecodeMP("kL8XELbXmYe1ELqV")
		So(res, ShouldEqual, "hello,world!")
	})
}

func TestOriBase64Decode(t *testing.T) {
	Convey("标准 base64decode", t, func() {
		dataStr := "eyJjcnRfY2lkIjoiMzc2NzI0ODEzIiwiY3J0X3JpZCI6IjYyNTgyYWM3ZDgyZjQ2MTI5YzNjNTUzeSIsInJ2X3RpZCI6IjEwMiIsImVjX2lkIjoiNzA0IiwidHBsZ3AiOiIxIiwidl9mbWQ1IjoiIiwiaV9mbWQ1IjoiYTBhYTA0Mzc3NTdmMDMzZjBmMTYzYjAyY2M2Nzc5NjIiLCJoNV90IjoxLCJtb2ZfdCI6MX0="
		res := OriBase64Decode(dataStr)
		So(string(res), ShouldEqual, `{"crt_cid":"376724813","crt_rid":"62582ac7d82f46129c3c553y","rv_tid":"102","ec_id":"704","tplgp":"1","v_fmd5":"","i_fmd5":"a0aa0437757f033f0f163b02cc677962","h5_t":1,"mof_t":1}`)
	})
}

func TestDeBase64(t *testing.T) {
	Convey("转换码表后的 base64decode", t, func() {
		byteData := []byte("dGVzdF9zdHJpbmc=")
		res, _ := DeBase64(byteData)
		So(string(res), ShouldEqual, "test_string")
	})
}

func TestRequestHeader(t *testing.T) {
	httpReq := http.Request{
		Method: "post",
		Proto:  "HTTP/1.0",
		Header: map[string][]string{
			"header_1": []string{"header_1_val_1", "header_1_val_2"},
			"header_2": []string{"header_2_val_1", "header_2_val_2"},
		},
	}

	Convey("请求头不存在的 key", t, func() {
		res := RequestHeader(&httpReq, "test_key")
		So(res, ShouldEqual, "")
	})

	Convey("请求头的 Method = post", t, func() {
		res := RequestHeader(&httpReq, "header_1")
		So(res, ShouldEqual, "header_1_val_1")
	})
}

func TestGetRequestID(t *testing.T) {
	Convey("GetRequestID返回不能重复", t, func() {
		reqID := GetRequestID()
		if reqID == GetRequestID() {
			t.Fatalf("getRequestID error")
		}
		So(reqID, ShouldNotEqual, GetRequestID())
	})
}

func BenchmarkGetTkClickID(b *testing.B) {
	b.Run("GetGoTkClickID", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			GetGoTkClickID()
		}
	})
	b.Run("GetGoTkClickIDNew", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			GetGoTkClickIDNew(GetRequestID())
		}
	})

}

func TestGetGoTkClickID(t *testing.T) {
	Convey("GetGoTkClickID 正常返回", t, func() {
		tkId := GetGoTkClickID()
		tkIdNew := GetGoTkClickIDNew(tkId)
		//t.Error(tkId)
		So(tkId, ShouldNotBeNil)
		So(tkIdNew, ShouldEqual, tkId)
	})
}

// func TestGen32BitsIDFromRequestID(t *testing.T) {
// 	Convey("返回空字符串", t, func() {
// 		res := Gen32BitsIDFromRequestID("test_string")
// 		So(res, ShouldEqual, "")
// 	})

// 	Convey("返回 32 字符长度 requestID", t, func() {
// 		res := Gen32BitsIDFromRequestID("abcd1234ABCD0123456789+-")
// 		So(res, ShouldEqual, "916176beabcd1234ABCD0123456789+-")
// 	})
// }

// func TestClientIP(t *testing.T) {
// 	Convey("Header 存在 X-Real-Ip 则从 X-Real-Ip 获取客户端 IP", t, func() {
// 		httpReq := http.Request{
// 			Method: "post",
// 			Proto:  "HTTP/1.0",
// 			Header: map[string][]string{
// 				"header_1":  []string{"header_1_val_1", "header_1_val_2"},
// 				"header_2":  []string{"header_2_val_1", "header_2_val_2"},
// 				"X-Real-Ip": []string{"100.0.0.1"},
// 			},
// 		}
// 		res := ClientIP(&httpReq)
// 		So(res, ShouldEqual, "100.0.0.1")
// 	})

// 	Convey("Header 不存在 X-Real-Ip，存在 X-Forwarded-For，则从 X-Forwarded-For 获取客户端 IP", t, func() {
// 		httpReq := http.Request{
// 			Method: "post",
// 			Proto:  "HTTP/1.0",
// 			Header: map[string][]string{
// 				"header_1":        []string{"header_1_val_1", "header_1_val_2"},
// 				"header_2":        []string{"header_2_val_1", "header_2_val_2"},
// 				"X-Forwarded-For": []string{"100.0.0.2,100.0.0.3"},
// 			},
// 		}
// 		res := ClientIP(&httpReq)
// 		So(res, ShouldEqual, "100.0.0.2")
// 	})

// 	Convey("否则，则从 RemoteAddr 获取客户端 IP", t, func() {
// 		httpReq := http.Request{
// 			Method:     "post",
// 			Proto:      "HTTP/1.0",
// 			RemoteAddr: "100.0.0.4:80",
// 			Header: map[string][]string{
// 				"header_1": []string{"header_1_val_1", "header_1_val_2"},
// 				"header_2": []string{"header_2_val_1", "header_2_val_2"},
// 			},
// 		}
// 		res := ClientIP(&httpReq)
// 		So(res, ShouldEqual, "100.0.0.4")
// 	})

// 	Convey("RemoteAddr 格式非法，返回空", t, func() {
// 		httpReq := http.Request{
// 			Method:     "post",
// 			Proto:      "HTTP/1.0",
// 			RemoteAddr: "100.0.0.5",
// 			Header: map[string][]string{
// 				"header_1": []string{"header_1_val_1", "header_1_val_2"},
// 				"header_2": []string{"header_2_val_1", "header_2_val_2"},
// 			},
// 		}
// 		res := ClientIP(&httpReq)
// 		So(res, ShouldEqual, "")
// 	})

// 	Convey("没有 IP 信息", t, func() {
// 		httpReq := http.Request{
// 			Method: "post",
// 			Proto:  "HTTP/1.0",
// 			Header: map[string][]string{
// 				"header_1": []string{"header_1_val_1", "header_1_val_2"},
// 				"header_2": []string{"header_2_val_1", "header_2_val_2"},
// 			},
// 		}
// 		res := ClientIP(&httpReq)
// 		So(res, ShouldEqual, "")
// 	})
// }

func TestVerCampare(t *testing.T) {
	Convey("版本相等", t, func() {
		res, err := VerCampare("1.0.0", "1.0.0")
		So(res, ShouldEqual, 0)
		So(err, ShouldBeNil)
	})

	Convey("版本 1 > 版本 2", t, func() {
		res, err := VerCampare("1.0.42", "1.0.1")
		So(res, ShouldEqual, 1)
		So(err, ShouldBeNil)
	})

	Convey("版本 1 < 版本 2", t, func() {
		res, err := VerCampare("1.0.1", "1.0.42")
		So(res, ShouldEqual, -1)
		So(err, ShouldBeNil)
	})

	Convey("版本 1 非法", t, func() {
		res, err := VerCampare("test.version", "1.0.1")
		So(res, ShouldEqual, -2)
		So(err, ShouldNotBeNil)
	})

	Convey("版本 2 非法", t, func() {
		res, err := VerCampare("1.0.1", "100_")
		So(res, ShouldEqual, -2)
		So(err, ShouldNotBeNil)
	})
}

func BenchmarkIntVer(b *testing.B) {
	b.ResetTimer()
	osVersion := "10.3.4.2"
	for i := 0; i < b.N; i++ {
		IntVer(osVersion)
	}
}

// func BenchmarkIntVerOld(b *testing.B) {
// 	b.ResetTimer()
// 	osVersion := "10.3.4.2"
// 	for i := 0; i < b.N; i++ {
// 		IntVerOld(osVersion)
// 	}
// }

func TestInVers(t *testing.T) {
	Convey("10.2.3.4", t, func() {
		osVersion := "10.2.3.4"
		res1, err1 := IntVer(osVersion)
		//res2, err2 := IntVerOld(osVersion)
		//So(res1, ShouldEqual, res2)
		//So(err1, ShouldEqual, err2)
		So(err1, ShouldBeNil)
		So(res1, ShouldEqual, 10020304)
		//osVersion = "3.2.4"
		//res1, err1 = IntVer(osVersion)
		//res2, err2 = IntVerOld(osVersion)
		//So(res1, ShouldEqual, res2)
		//So(err1, ShouldEqual, err2)
		res2, err := IntVer("1.4")
		t.Log(res2, err)
	})
}

// func TestIntVerOld(t *testing.T) {
// 	Convey("版本号 1.0.1 转换为数字", t, func() {
// 		res, err := IntVerOld("1.0.1")
// 		So(res, ShouldEqual, int32(1000100))
// 		So(err, ShouldBeNil)
// 	})

// 	Convey("版本号 42 转换为数字", t, func() {
// 		res, err := IntVerOld("42")
// 		So(res, ShouldEqual, int32(42000000))
// 		So(err, ShouldBeNil)
// 	})

// 	Convey("非法版本号", t, func() {
// 		res, err := IntVerOld("42_1")
// 		So(res, ShouldEqual, int32(0))
// 		So(err, ShouldNotBeNil)
// 	})
// }

func TestIntVer(t *testing.T) {
	Convey("版本号 1.0.1 转换为数字", t, func() {
		res, err := IntVer("1.0.1")
		So(res, ShouldEqual, int32(1000100))
		So(err, ShouldBeNil)
	})

	Convey("版本号 42 转换为数字", t, func() {
		res, err := IntVer("42")
		So(res, ShouldEqual, int32(42000000))
		So(err, ShouldBeNil)
	})

	Convey("非法版本号", t, func() {
		res, err := IntVer("42_1")
		So(res, ShouldEqual, int32(0))
		So(err, ShouldNotBeNil)
	})
}

func TestInStrArray(t *testing.T) {
	arr := []string{"test1", "test2", "test3"}
	Convey("字符串在数组中", t, func() {
		res := InStrArray("test1", arr)
		So(res, ShouldBeTrue)
	})

	Convey("字符串不在数组中", t, func() {
		res := InStrArray("test4", arr)
		So(res, ShouldBeFalse)
	})
}

func TestInArray(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}
	Convey("数字在数组中", t, func() {
		res := InArray(1, arr)
		So(res, ShouldBeTrue)
	})

	Convey("数字不在数组中", t, func() {
		res := InArray(42, arr)
		So(res, ShouldBeFalse)
	})
}

func TestInInt64Arr(t *testing.T) {
	arr := []int64{1, 2, 3, 4, 5}
	Convey("数字在数组中", t, func() {
		res := InInt64Arr(1, arr)
		So(res, ShouldBeTrue)
	})

	Convey("数字不在数组中", t, func() {
		res := InInt64Arr(42, arr)
		So(res, ShouldBeFalse)
	})
}

func TestInInt32Arr(t *testing.T) {
	arr := []int32{1, 2, 3, 4, 5}
	Convey("数字在数组中", t, func() {
		res := InInt32Arr(1, arr)
		So(res, ShouldBeTrue)
	})

	Convey("数字不在数组中", t, func() {
		res := InInt32Arr(42, arr)
		So(res, ShouldBeFalse)
	})
}

// func TestGetInternalIP(t *testing.T) {
// 	// todo
// }

// func TestGetEndCardUrl(t *testing.T) {
// 	Convey("offerwall_urls.http.rewardvideo_end_screen is empty", t, func() {
// 		res, err := GetEndCardUrl(1, "1.0.1")
// 		So(res, ShouldEqual, "")
// 		So(err, ShouldNotBeNil)
// 	})

// 	Convey("end_screen 包含 ?", t, func() {
// 		Rewardvideo_end_screen.Value.Http.Rewardvideo_end_screen = "http://www.mobvista.com?endcard"
// 		res, err := GetEndCardUrl(1, "1.0.1")
// 		So(res, ShouldEqual, "http://www.mobvista.com?endcard&unit_id=1&sdk_version=1.0.1")
// 		So(err, ShouldBeNil)
// 	})

// 	Convey("end_screen 不包含 ?", t, func() {
// 		Rewardvideo_end_screen.Value.Http.Rewardvideo_end_screen = "http://www.mobvista.com/endcard"
// 		res, err := GetEndCardUrl(1, "1.0.1")
// 		So(res, ShouldEqual, "http://www.mobvista.com/endcard?unit_id=1&sdk_version=1.0.1")
// 		So(err, ShouldBeNil)
// 	})
// }

// func TestGetRvImagesKey(t *testing.T) {
// 	Convey("play_all_300x250", t, func() {
// 		res := GetRvImagesKey("250", "300")
// 		So(res, ShouldEqual, "play_all_300x250")
// 	})

// 	Convey("click_all_500x500", t, func() {
// 		res := GetRvImagesKey("500", "500")
// 		So(res, ShouldEqual, "click_all_500x500")
// 	})

// 	Convey("endcard_portrait_600x800", t, func() {
// 		res := GetRvImagesKey("800", "600")
// 		So(res, ShouldEqual, "endcard_portrait_600x800")
// 	})

// 	Convey("endcard_landscape_800x600", t, func() {
// 		res := GetRvImagesKey("600", "800")
// 		So(res, ShouldEqual, "endcard_landscape_800x600")
// 	})

// 	Convey("other sizes", t, func() {
// 		res := GetRvImagesKey("42", "42")
// 		So(res, ShouldEqual, "")
// 	})
// }

// func TestFormatScreenSize(t *testing.T) {
// 	Convey("screenSize 不包含 x", t, func() {
// 		width, height := FormatScreenSize("", "")
// 		So(width, ShouldEqual, 0)
// 		So(height, ShouldEqual, 0)
// 	})

// 	Convey("screenSize 包含 x", t, func() {
// 		width, height := FormatScreenSize("", "300x250")
// 		So(width, ShouldEqual, 300)
// 		So(height, ShouldEqual, 250)
// 	})

// 	Convey("screenSize 不合法", t, func() {
// 		Convey("screenSize 包含多个 x", func() {
// 			width, height := FormatScreenSize("", "300x250x123")
// 			So(width, ShouldEqual, 0)
// 			So(height, ShouldEqual, 0)
// 		})
// 	})
// }

func TestGetUrlScheme(t *testing.T) {
	Convey("参数 2 时返回 https://", t, func() {
		res := GetUrlScheme(2)
		So(res, ShouldEqual, "https://")
	})

	Convey("其他参数时返回 http://", t, func() {
		res := GetUrlScheme(0)
		So(res, ShouldEqual, "http://")
	})
}

func TestGetAdTypeStr(t *testing.T) {
	Convey("unknown adtype", t, func() {
		res := GetAdTypeStr(0)
		So(res, ShouldEqual, "unknown")
	})

	Convey("text adtype", t, func() {
		res := GetAdTypeStr(1)
		So(res, ShouldEqual, "text")
	})

	Convey("banner adtype", t, func() {
		res := GetAdTypeStr(2)
		So(res, ShouldEqual, "banner")
	})

	Convey("appwall adtype", t, func() {
		res := GetAdTypeStr(3)
		So(res, ShouldEqual, "appwall")
	})

	Convey("overlay adtype", t, func() {
		res := GetAdTypeStr(4)
		So(res, ShouldEqual, "overlay")
	})

	Convey("fullscreen adtype", t, func() {
		res := GetAdTypeStr(5)
		So(res, ShouldEqual, "full_screen")
	})

	Convey("interstitial adtype", t, func() {
		res := GetAdTypeStr(29)
		So(res, ShouldEqual, "interstitial")
	})

	Convey("native adtype", t, func() {
		res := GetAdTypeStr(42)
		So(res, ShouldEqual, "native")
	})

	Convey("rewarded video adtype", t, func() {
		res := GetAdTypeStr(94)
		So(res, ShouldEqual, "rewarded_video")
	})

	Convey("feeds video adtype", t, func() {
		res := GetAdTypeStr(95)
		So(res, ShouldEqual, "feeds_video")
	})

	Convey("offerwall adtype", t, func() {
		res := GetAdTypeStr(278)
		So(res, ShouldEqual, "offerwall")
	})

	Convey("interstitial sdk adtype", t, func() {
		res := GetAdTypeStr(279)
		So(res, ShouldEqual, "interstitial_sdk")
	})

	Convey("other", t, func() {
		res := GetAdTypeStr(99999)
		So(res, ShouldEqual, "")
	})

	// todo: online video lost
}

func TestSerializeMPV2(t *testing.T) {
	param := &Params{}
	param.RequestID = "requestId"
	param.PublisherID = 123
	param.AppID = 456
	param.UnitID = 345
	param.MWadBackend = "MWadBackend"
	param.MWadBackendData = "MWadBackendData"
	param.Scenario = "open_api"
	param.AdType = 42
	param.ImageSize = "imageSize"
	param.RequestType = 4
	param.Platform = 2
	param.OSVersion = "10.13.2"
	param.SDKVersion = "1.2"
	param.Model = "model"
	param.Orientation = 2
	param.ScreenSize = "ScreenSize"
	param.CountryCode = "HK"
	param.Language = "zh"
	param.NetworkType = 9
	param.MCC = "460"
	param.MNC = "01"
	param.Extra = "extra"
	param.Extra3 = "extra3"
	param.Extra4 = "extra4"
	param.Extra5 = "extra5"
	param.Extra7 = 7
	param.Extra8 = 8
	param.Extra9 = "extra9"
	param.Extra10 = "extra10"
	param.RequestID = "requestId"
	param.ClientIP = "192.168.1.1"
	param.IMEI = "imei"
	param.MAC = "mac"
	param.AndroidID = "androidID"
	param.ServerIP = "serverIP"
	param.GAID = "gaid"
	param.IDFA = "idfa"
	param.AppVersionName = "appversionName"
	param.Brand = "brand"
	param.RemoteIP = "remoteIp"
	param.SessionID = "sessionID"
	param.ParentSessionID = "parentSessionID"
	param.CityCode = 3256
	param.AdNum = 13
	param.TNum = 14
	param.MWRandValue = 15
	param.Extra16 = 16
	param.IDFV = "idfv"
	param.OpenIDFA = "openIdfv"
	param.MWbackendConfig = "MWbackendConfig"
	param.Extfinalsubid = 10
	param.ExtpackageName = "extpackageName"
	param.MWFlowTagID = 100
	param.Extendcard = "extendcard"
	param.ExtrushNoPre = 1
	param.ExtfinalPackageName = "extfinalPackageName"
	param.Extnativex = 2
	param.Extctype = 42
	param.Extrvtemplate = 1
	param.Extabtest1 = 1
	param.Extb2t = 1
	param.Extchannel = "extChannel"
	param.Extbp = "extbp"
	param.Extsource = 32
	param.Extalgo = "extalgo"
	param.MWplayInfo = "MWplayInfo"
	param.ExtifLowerImp = 0
	param.ExtsystemUseragent = "ExtsystemUseragent"
	param.ExtMpNormalMap = "extmpnormalMap"
	mwCreative := "mwCreative"
	Convey("SerializeMP equal", t, func() {
		res1 := SerializeMP(param, mwCreative)
		//res2 := SerializeMPOld(param, mwCreative)
		So(res1, ShouldNotBeEmpty)
	})
}

func BenchmarkSerializeMP(b *testing.B) {
	param := &Params{}
	param.RequestID = "requestId"
	param.PublisherID = 123
	param.AppID = 456
	param.UnitID = 345
	param.MWadBackend = "MWadBackend"
	param.MWadBackendData = "MWadBackendData"
	param.Scenario = "open_api"
	param.AdType = 42
	param.ImageSize = "imageSize"
	param.RequestType = 4
	param.Platform = 2
	param.OSVersion = "10.13.2"
	param.SDKVersion = "1.2"
	param.Model = "model"
	param.Orientation = 2
	param.ScreenSize = "ScreenSize"
	param.CountryCode = "HK"
	param.Language = "zh"
	param.NetworkType = 9
	param.MCC = "460"
	param.MNC = "01"
	param.Extra = "extra"
	param.Extra3 = "extra3"
	param.Extra4 = "extra4"
	param.Extra5 = "extra5"
	param.Extra7 = 7
	param.Extra8 = 8
	param.Extra9 = "extra9"
	param.Extra10 = "extra10"
	param.RequestID = "requestId"
	param.ClientIP = "192.168.1.1"
	param.IMEI = "imei"
	param.MAC = "mac"
	param.AndroidID = "androidID"
	param.ServerIP = "serverIP"
	param.GAID = "gaid"
	param.IDFA = "idfa"
	param.AppVersionName = "appversionName"
	param.Brand = "brand"
	param.RemoteIP = "remoteIp"
	param.SessionID = "sessionID"
	param.ParentSessionID = "parentSessionID"
	param.CityCode = 3256
	param.AdNum = 13
	param.TNum = 14
	param.MWRandValue = 15
	param.Extra16 = 16
	param.IDFV = "idfv"
	param.OpenIDFA = "openIdfv"
	param.MWbackendConfig = "MWbackendConfig"
	param.Extfinalsubid = 10
	param.ExtpackageName = "extpackageName"
	param.MWFlowTagID = 100
	param.Extendcard = "extendcard"
	param.ExtrushNoPre = 1
	param.ExtfinalPackageName = "extfinalPackageName"
	param.Extnativex = 2
	param.Extctype = 42
	param.Extrvtemplate = 1
	param.Extabtest1 = 1
	param.Extb2t = 1
	param.Extchannel = "extChannel"
	param.Extbp = "extbp"
	param.Extsource = 32
	param.Extalgo = "extalgo"
	param.MWplayInfo = "MWplayInfo"
	param.ExtifLowerImp = 0
	param.ExtsystemUseragent = "ExtsystemUseragent"
	param.ExtMpNormalMap = "extmpnormalMap"
	mwCreative := "mwCreative"
	b.Run("SerializeMP", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			SerializeMP(param, mwCreative)
		}
	})
	// b.Run("SerializeMPOld", func(b *testing.B) {
	// 	for i := 0; i < b.N; i++ {
	// 		SerializeMPOld(param, mwCreative)
	// 	}
	// })
}

func TestSerializeMP(t *testing.T) {
	Convey("序列化空 params", t, func() {
		params := &Params{}
		res := SerializeMP(params, "0")
		So(res, ShouldEqual, "6aSIidMM6deI6aSIidMM6deI6deI6aSI6aSI6deI6deI6deIidMM6aSIideI6deI6aSI6aSI6deI6deIideIideI")
	})

	Convey("序列化 params", t, func() {
		params := &Params{
			RequestID:   "test_request_ID",
			PublisherID: 1000,
			AppID:       42,
			CityCode:    1,
		}
		res := SerializeMP(params, "0")
		So(res, ShouldEqual, "J75AJcKBH+c2H+fT+TQj6ajMiaSIfazIideI6dMe6aSIideI6deI6dMM6dMM6deI6deI6deI6aSIidMM6aSI6deI6dMM6dMM6deI6deI6aSI6aSI6v==")
	})
}

func TestSerializeMPPart(t *testing.T) {
	Convey("序列化参数", t, func() {
		// pubID := int64(100)
		// appID := int64(42)
		// unitID := int64(3)
		// queryParam := corsair_proto.QueryParam{
		// 	PublisherId: &pubID,
		// 	AppId:       &appID,
		// 	UnitId:      &unitID,
		// }
		reqParams := RequestParams{
			Param: Params{
				MWplayInfo: "test_mw_play_info",
				Extchannel: "1",
			},
		}
		reqCTX := ReqCtx{
			//ParamList: &queryParam,
			ReqParams: &reqParams,
		}
		res := SerializeMPPart(reqCTX.ReqParams, "ad_backend", "ad_backend_data",
			"backend_config", "", "345", "8:0:0.83")
		So(res, ShouldEqual, "|0|0|0|ad_backend|ad_backend_data||0|0|0|||||||0||0||0|0|0||||||0|0|0|0|backend_config|test_mw_play_info||||0||345||||8:0:0.83|1||0")
	})
}

func TestSerializeMPPartV2(t *testing.T) {
	adBackend := "1,8"
	adBackendData := "adBackendData"
	params := Params{}
	params.RequestID = "requestID"
	params.PublisherID = 345
	params.AppID = 456
	params.UnitID = 567
	params.CountryCode = "CN"
	params.CityCode = 9527
	params.Platform = 1
	params.AdType = 42
	params.OSVersion = "10.23"
	params.SDKVersion = "mi_2.3.0"
	params.AppVersionName = "3.2"
	params.Brand = "brand"
	params.Model = "model"
	params.ScreenSize = "ScreenSize"
	params.Orientation = 2
	params.Language = "zh"
	params.NetworkType = 9
	params.MCC = "460"
	params.MNC = "01"
	params.IMEI = "imei"
	params.MAC = "mac"
	params.GAID = "gaid"
	params.AndroidID = "androidid"
	params.IDFA = "idfa"
	params.AdNum = 10
	params.Scenario = "open_api"
	params.MWplayInfo = "mwplayInfo"
	reqParam := &RequestParams{Param: params}
	reqParam.FlowTagID = 100
	backendConfig := "backendConfig"
	Convey("SerializeMPPart equal", t, func() {
		res1 := SerializeMPPart(reqParam, adBackend, adBackendData, backendConfig, "", "345", "8:0:0.83")
		//res2 := SerializeMPPartOld(reqParam, adBackend, adBackendData, backendConfig)
		So(res1, ShouldNotBeEmpty)
	})
}

func BenchmarkSerializeMPPart(b *testing.B) {
	adBackend := "1,8"
	adBackendData := "adBackendData"
	params := Params{}
	params.RequestID = "requestID"
	params.PublisherID = 345
	params.AppID = 456
	params.UnitID = 567
	params.CountryCode = "CN"
	params.CityCode = 9527
	params.Platform = 1
	params.AdType = 42
	params.OSVersion = "10.23"
	params.SDKVersion = "mi_2.3.0"
	params.AppVersionName = "3.2"
	params.Brand = "brand"
	params.Model = "model"
	params.ScreenSize = "ScreenSize"
	params.Orientation = 2
	params.Language = "zh"
	params.NetworkType = 9
	params.MCC = "460"
	params.MNC = "01"
	params.IMEI = "imei"
	params.MAC = "mac"
	params.GAID = "gaid"
	params.AndroidID = "androidid"
	params.IDFA = "idfa"
	params.AdNum = 10
	params.Scenario = "open_api"
	params.MWplayInfo = "mwplayInfo"
	reqParam := &RequestParams{Param: params}
	reqParam.FlowTagID = 100
	backendConfig := "backendConfig"
	b.Run("SerializeMPPart", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			SerializeMPPart(reqParam, adBackend, adBackendData, backendConfig, "", "345", "8:0:0.83")
		}
	})
	// b.Run("SerializeMPPartOld", func(b *testing.B) {
	// 	for i := 0; i < b.N; i++ {
	// 		SerializeMPPartOld(reqParam, adBackend, adBackendData, backendConfig)
	// 	}
	// })
}

func TestSerializePV2(t *testing.T) {
	params := &Params{}
	params.AdType = 42
	params.ImageSize = "ImageSize"
	params.PlatformName = "ios"
	params.OSVersion = "10.11"
	params.SDKVersion = "mi_2.3.0"
	params.Model = "model"
	params.ScreenSize = "screenSize"
	params.Orientation = 2
	params.Language = "zh"
	params.NetworkTypeName = "wifi"
	params.MCC = "460"
	params.MNC = "01"
	params.Extra3 = "extra3"
	params.Extra4 = "extra4"
	params.ClientIP = "192.168.0.1"
	params.IMEI = "imei"
	params.MAC = "mac"
	params.AndroidID = "androidid"
	params.GAID = "gaid"
	params.IDFA = "idfa"
	params.Brand = "brand"
	params.RemoteIP = "192.168.0.1"
	params.SessionID = "sessionID"
	params.ParentSessionID = "parentSessionID"
	params.TNum = 10
	params.IDFV = "idfv"
	params.OpenIDFA = "openidfa"
	params.Extstats = "extstats"
	Convey("SerializeP equal", t, func() {
		//camInfo := &CampaignInfo{}
		res1 := SerializeP(params)
		//res2 := SerializePOld(params, camInfo)
		So(res1, ShouldNotBeEmpty)
	})
}

func BenchmarkSerializeP(b *testing.B) {
	params := &Params{}
	params.AdType = 42
	params.ImageSize = "ImageSize"
	params.PlatformName = "ios"
	params.OSVersion = "10.11"
	params.SDKVersion = "mi_2.3.0"
	params.Model = "model"
	params.ScreenSize = "screenSize"
	params.Orientation = 2
	params.Language = "zh"
	params.NetworkTypeName = "wifi"
	params.MCC = "460"
	params.MNC = "01"
	params.Extra3 = "extra3"
	params.Extra4 = "extra4"
	params.ClientIP = "192.168.0.1"
	params.IMEI = "imei"
	params.MAC = "mac"
	params.AndroidID = "androidid"
	params.GAID = "gaid"
	params.IDFA = "idfa"
	params.Brand = "brand"
	params.RemoteIP = "192.168.0.1"
	params.SessionID = "sessionID"
	params.ParentSessionID = "parentSessionID"
	params.TNum = 10
	params.IDFV = "idfv"
	params.OpenIDFA = "openidfa"
	params.Extstats = "extstats"
	//camInfo := &CampaignInfo{}
	b.Run("SerializeP", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			SerializeP(params)
		}
	})
	// b.Run("SerializePNew", func(b *testing.B) {
	// 	for i := 0; i < b.N; i++ {
	// 		SerializePOld(params, camInfo)
	// 	}
	// })
}

func TestSerializeP(t *testing.T) {
	params := Params{
		AdType:       42,
		ImageSize:    "300x300",
		PlatformName: "ios",
		OSVersion:    "1.1.1",
	}
	//campaign := CampaignInfo{}
	Convey("序列化参数", t, func() {
		res := SerializeP(&params)
		So(res, ShouldNotBeEmpty)
		//So(res, ShouldEqual, "fHx8fHx8fG5hdGl2ZXwzMDB4MzAwfHxpb3N8MS4xLjF8fHx8MHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8MHx8LHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8")
	})
}

func TestSerializeQ(t *testing.T) {
	params := &Params{
		AdType:       42,
		ImageSize:    "300x300",
		PlatformName: "ios",
		OSVersion:    "1.1.1",
	}

	Convey("序列化参数，campaign.OriPrice、campaign.Price 为空", t, func() {
		campaign := &smodel.CampaignInfo{}
		res := SerializeQ(params, campaign)
		So(res, ShouldNotBeEmpty)
		//So(res, ShouldEqual, "a_i09M6deI6aSIidMM6aSIidMMWUvM6av/iaSI6aSI6aSI6aSIidMM6deIideIW+MM6aSIidMM6aSIidMM6deI6aSI6aSIidMM6dMM6aSI6aSIideIideIidMM6aSI")
	})

	Convey("序列化参数，campaign.OriPrice 为空", t, func() {
		price := float64(1.23)
		campaign := &smodel.CampaignInfo{
			Price: price,
		}
		res := SerializeQ(params, campaign)
		So(res, ShouldNotBeEmpty)
		//So(res, ShouldEqual, "a_i09M6deI6aSIidMM6aSIidMMWUvM6aj/iUfI6aSI6aSI6aSIidMM6deIideIW+MM6aSIidMM6aSIidMM6deI6aSI6aSIidMM6dMM6aSI6aSIideIideIidMM6aSI")
	})

	Convey("序列化参数，campaign.Price 为空", t, func() {
		oriPrice := float64(0.01)
		campaign := &smodel.CampaignInfo{
			OriPrice: oriPrice,
		}
		res := SerializeQ(params, campaign)
		So(res, ShouldNotBeEmpty)
		//So(res, ShouldEqual, "a_i09M6deI6aSIidMM6aSIidMMWUve6av/iaSI6aSI6aSI6aSIidMM6deIideIW+MM6aSIidMM6aSIidMM6deI6aSI6aSIidMM6dMM6aSI6aSIideIideIidMM6aSI")
	})
}

func BenchmarkSerializeCSP(b *testing.B) {
	param := &Params{}
	param.MWadBackend = "8,2"
	param.MWadBackendData = "1:2333334:1"
	param.MWFlowTagID = 102
	param.MWRandValue = 3721
	param.MWbackendConfig = "hello,world"
	param.AdNum = 20
	param.TNum = 5
	b.Run("SerializeCSP", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			SerializeCSP(param, "")
		}
	})
	// b.Run("SerializeCSPOld", func(b *testing.B) {
	// 	for i := 0; i < b.N; i++ {
	// 		SerializeCSPOld(param)
	// 	}
	// })
	// b.Run("SerializeCSPOld", func(b *testing.B) {
	// 	for i := 0; i < b.N; i++ {
	// 		SerializeCSPOld(param)
	// 	}
	// })
}

func TestSerializeCSP(t *testing.T) {
	params := &Params{
		MWadBackend:     "ad_backend",
		MWadBackendData: "test",
		MWFlowTagID:     42,
		MWRandValue:     1,
	}

	Convey("序列化参数，campaign.Price 为空", t, func() {
		res := SerializeCSP(params, "")
		//res2 := SerializeCSPOld(params)
		So(res, ShouldEqual, "Dkx6DrcULF5/HdeTH+fT6aNB6acI6aSIidM=")
		//So(res, ShouldEqual, res2)
	})
}

func TestCheckParam(t *testing.T) {
	Convey("没有 | 字符", t, func() {
		res := CheckParam("test")
		So(res, ShouldEqual, "test")
	})

	Convey("包含一个 | 字符", t, func() {
		res := CheckParam("tes|t")
		So(res, ShouldEqual, "test")
	})

	Convey("包含多个 | 字符", t, func() {
		res := CheckParam("t|e|s|t")
		So(res, ShouldEqual, "test")
	})
}

// func TestGetAdTrackURL(t *testing.T) {
// 	var commonConfig CommonConfig
// 	commonConfig.TrackConfig.TrackHost = "mtrack.rayjump.com"
// 	commonConfig.TrackConfig.PlayTrackPath = "/ad/log/play"
// 	Config.CommonConfig = &commonConfig
// 	// pubID := int64(100)
// 	// appID := int64(42)
// 	// unitID := int64(3)
// 	// httpReqNum := int32(1000)
// 	// queryParam := corsair_proto.QueryParam{
// 	// 	PublisherId: &pubID,
// 	// 	AppId:       &appID,
// 	// 	UnitId:      &unitID,
// 	// 	HttpReq:     &httpReqNum,
// 	// }
// 	var param Params
// 	param.NetworkType = mvconst.NETWORK_TYPE_4G
// 	param.Platform = mvconst.PlatformIOS
// 	param.RequestID = "5b167141c6c1e24c7cf21b9b"
// 	param.IDFA = "35951D6A-DB66-47A1-8CD0-DB1E9E0D33CD"
// 	param.ClientIP = "192.168.1.1"
// 	param.UserAgent = "Mozilla%252f5.0%2B%28iPhone%25253B%2BCPU%2BiPhone%2BOS%2B9_3_5%2Blike%2BMac%2BOS%2BX%29%2BAppleWebKit%252f601.1.46%2B%28KHTML%25252C%2Blike%2BGecko%29%2BMobile%252f13G36"
// 	param.TNum = 1
// 	param.OSVersion = "9.3.5"
// 	param.Model = "iphone6%2C2"
// 	param.UnitID = 128
// 	param.AppName = "iOSTest"
// 	param.PackageName = "com.tianye1.mvsdk"
// 	param.FormatAdType = mvconst.ADTypeNativeVideo
// 	param.AdType = mvconst.ADTypeNative
// 	param.VideoVersion = "1.1"
// 	param.VideoH = 300
// 	param.VideoW = 400
// 	param.MWplayInfo = "test_mw_play_info"
// 	param.Extchannel = "1"
// 	reqParam := &RequestParams{Param: param}
// 	reqCTX := ReqCtx{
// 		//ParamList: &queryParam,
// 		ReqParams: reqParam,
// 	}

// 	Convey("包含多个 | 字符", t, func() {
// 		res := GetAdTrackURL(reqCTX.ReqParams, 1, "offer_ID_1", "backend_config", "extra_param")
// 		So(res, ShouldEqual, "http://mtrack.rayjump.com/ad/log/play?k=5b167141c6c1e24c7cf21b9b&mp=fkRefUhefacUfrieHnRTDAJUHUReDUQ06aSIidMeiUtIi%2BMeGrKrHr5B%2BTQj%2BAj1i%2BeIidMB6aNB6a3%2FiB926deI67QML7K%2FHnDQi3iB6dMM6dMT6dMM6aSIideI6deIiAVPfncjf3jwxjRFf0TTfTjeWntaxavwxjRexnQcijNAiTfj6aSIidMe6aSIDrcULF5%2FHcKUYFPrLkJIJ75AJcKwJ2KMY7cP%2BFQ%2FHrKI6deIideIideI6deIi%2BeIiv%3D%3D&extra_param")
// 	})
// }

func TestTrimBlank(t *testing.T) {
	Convey("无换行符", t, func() {
		res := TrimBlank("test_sting_without_rn")
		So(res, ShouldEqual, "test_sting_without_rn")
	})

	Convey("去除换行符和 tab", t, func() {
		res := TrimBlank("test_sting_with_n\n_n_and_with_t\t_t")
		So(res, ShouldEqual, "test_sting_with_n_n_and_with_t_t")
	})
}

func TestRandByRate(t *testing.T) {
	Convey("只有一个元素，一定 rand 到该元素", t, func() {
		rateMap := map[int]int{
			21: 42,
		}
		res := RandByRate(rateMap)
		So(res, ShouldEqual, 21)
	})
}

func TestRandIntArr(t *testing.T) {
	Convey("只有一个元素，一定 rand 到该元素", t, func() {
		arr := []int{42}
		res := RandIntArr(arr)
		So(res, ShouldEqual, 42)
	})

	Convey("空数组返回 0", t, func() {
		arr := []int{}
		res := RandIntArr(arr)
		So(res, ShouldEqual, 0)
	})
}

func TestRandInt64Arr(t *testing.T) {
	Convey("只有一个元素，一定 rand 到该元素", t, func() {
		arr := []int64{42}
		res := RandInt64Arr(arr)
		So(res, ShouldEqual, 42)
	})

	Convey("空数组返回 0", t, func() {
		arr := []int64{}
		res := RandInt64Arr(arr)
		So(res, ShouldEqual, 0)
	})
}

func TestGetRandConsiderZero(t *testing.T) {
	Convey("没有 idfa 和 gaid 返回 -1", t, func() {
		res := GetRandConsiderZero("", "", "test_salt", 42)
		So(res, ShouldEqual, -1)
	})

	Convey("没有 idfa 和 gaid 返回 -1", t, func() {
		res := GetRandConsiderZero("gaid", "idfa", "test_salt", 42)
		So(res, ShouldEqual, -1)
	})

	Convey("没有 idfa 和 gaid 返回 -1", t, func() {
		res := GetRandConsiderZero("gaid", "00000000-0000-0000-0000-000000000000", "test_salt", 42)
		So(res, ShouldEqual, -1)
	})

	Convey("没有 idfa 和 gaid 返回 -1", t, func() {
		res := GetRandConsiderZero("test_gaid", "test_idfa", "test_salt", 42)
		So(res, ShouldEqual, 9)
	})
}

func TestGetPureRand(t *testing.T) {
	Convey("randSum 为0的情况，返回0", t, func() {
		res := GetPureRand(0)
		So(res, ShouldEqual, 0)
	})
	Convey("randSum>0的情况，返回[0,randSum)的值", t, func() {
		randSum := 2
		for i := 0; i < 5; i++ { // 随机数, 多测几次
			res := GetPureRand(randSum)
			So(res, ShouldBeGreaterThanOrEqualTo, 0)
			So(res, ShouldBeLessThan, randSum)
		}
	})
}

func TestRandArr(t *testing.T) {
	Convey("没有 idfa 和 gaid 返回 0", t, func() {
		res := RandArr(map[int]int{}, "", "", "test_salt")
		So(res, ShouldEqual, 0)
	})

	Convey("randArr 只有一个元素，返回该元素 key", t, func() {
		res := RandArr(map[int]int{1: 42}, "test", "", "test_salt")
		So(res, ShouldEqual, 1)
	})

	Convey("返回 2", t, func() {
		res := RandArr(map[int]int{1: 0, 2: 1, 3: 40}, "test_gaid", "test_idfa", "test_salt")
		So(res, ShouldEqual, 3)
	})
}

func TestIsDevidEmpty(t *testing.T) {
	Convey("devID 为空", t, func() {
		params := &Params{
			GAID: "",
			IDFA: "",
		}
		res := IsDevidEmpty(params)
		So(res, ShouldBeTrue)
	})

	Convey("gaid = gaid & idfa = idfa 为空", t, func() {
		params := &Params{
			GAID: "gaid",
			IDFA: "idfa",
		}
		res := IsDevidEmpty(params)
		So(res, ShouldBeTrue)
	})

	Convey("gaid = gaid & idfa = 00000000-0000-0000-0000-000000000000 为空", t, func() {
		params := &Params{
			GAID: "gaid",
			IDFA: "00000000-0000-0000-0000-000000000000",
		}
		res := IsDevidEmpty(params)
		So(res, ShouldBeTrue)
	})

	Convey("gaid 或 idfa 有值则不为空", t, func() {
		Convey("GAID 非空", func() {
			params := &Params{
				GAID: "test_gaid",
				IDFA: "",
			}
			res := IsDevidEmpty(params)
			So(res, ShouldBeFalse)
		})

		Convey("IDFA 非空", func() {
			params := &Params{
				GAID: "",
				IDFA: "test_idfa",
			}
			res := IsDevidEmpty(params)
			So(res, ShouldBeFalse)
		})
	})
}

func TestMax(t *testing.T) {
	Convey("返回 2", t, func() {
		res := Max(1, 2)
		So(res, ShouldEqual, 2)
	})

	Convey("返回 1", t, func() {
		res := Max(1, 1)
		So(res, ShouldEqual, 1)
	})
}

func TestSubString(t *testing.T) {
	Convey("返回 test", t, func() {
		res := SubString("test_string", 0, 4)
		So(res, ShouldEqual, "test")
	})

	Convey("返回 test", t, func() {
		res := SubString("test_string", -1, 4)
		So(res, ShouldEqual, "test")
	})

	Convey("返回 test", t, func() {
		res := SubString("test_string", 12, 4)
		So(res, ShouldEqual, "")
	})

	Convey("返回 test", t, func() {
		res := SubString("test_string", 5, 20)
		So(res, ShouldEqual, "string")
	})
}

func TestGetVersionCode(t *testing.T) {
	Convey("TestGetVersionCode", t, func() {
		Convey("空字符返回 0", func() {
			res := GetVersionCode("")
			So(res, ShouldEqual, 0)
		})

		Convey("正常返回", func() {
			res := GetVersionCode("1.2.3")
			So(res, ShouldEqual, 10203)
		})

		Convey("返回前三组数字", func() {
			res := GetVersionCode("1.2.3.4")
			So(res, ShouldEqual, 10203)
		})

		Convey("排除非法字符", func() {
			res := GetVersionCode("1.10——2.3.4")
			So(res, ShouldEqual, 10003)
		})

		// todo, now will error
		Convey(". 分隔后少于三个部分", func() {
			// res := GetVersionCode("1.1")
			// So(res, ShouldEqual, 10101)
		})
	})
}

func TestHttpBuildQuery(t *testing.T) {
	Convey("query 排序", t, func() {
		queryMap := map[string]string{
			"3":      "val_3",
			"test_2": "val_2",
			"1":      "val_1",
		}
		res := HttpBuildQuery(queryMap)
		So(res, ShouldEqual, "1=val_1&3=val_3&test_2=val_2")
	})
}

func TestUrlEncode(t *testing.T) {
	Convey("UrlEncode", t, func() {
		res := UrlEncode("http://www.test.com?test=hello world&param2=13")
		So(res, ShouldEqual, "http%3A%2F%2Fwww.test.com%3Ftest%3Dhello+world%26param2%3D13")
	})
}

func TestUrlDecode(t *testing.T) {
	Convey("UrlDecode", t, func() {
		res := UrlDecode("http%3A%2F%2Fwww.test.com%3Ftest%3Dhello+world%26param2%3D13")
		So(res, ShouldEqual, "http://www.test.com?test=hello world&param2=13")
	})
}

func TestSubUtf8Str(t *testing.T) {
	Convey("长度不用截取", t, func() {
		res := SubUtf8Str("test_string", 100)
		So(res, ShouldEqual, "test_string")
	})

	Convey("截取 3 个字符", t, func() {
		res := SubUtf8Str("test_string", 3)
		So(res, ShouldEqual, "tes......")
	})

	Convey("截取 3 个中文字符", t, func() {
		res := SubUtf8Str("测试字符串", 3)
		So(res, ShouldEqual, "测试字......")
	})

	Convey("截取 3 个其他文字字符", t, func() {
		res := SubUtf8Str("Больше", 3)
		So(res, ShouldEqual, "Бол......")
	})
}

func TestSerializeOImpV2(t *testing.T) {
	param := &Params{}
	param.PublisherID = 123
	param.AppID = 456
	param.UnitID = 345
	param.AdvCreativeID = 934
	param.CampaignID = 23456
	param.Scenario = "open_api"
	param.AdType = 42
	param.ImageSize = "imageSize"
	param.RequestType = 4
	param.PlatformName = "ios"
	param.OSVersion = "10.13.2"
	param.SDKVersion = "1.2"
	param.Model = "model"
	param.ScreenSize = "screenSize"
	param.Orientation = 2
	param.CountryCode = "HK"
	param.Language = "zh"
	param.NetworkTypeName = "wifi"
	param.MCC = "460"
	param.MNC = "01"
	param.Extra = "extra"
	param.Extra3 = "extra3"
	param.Extra4 = "extra4"
	param.Extra5 = "extra5"
	param.Extra7 = 7
	param.Extra8 = 8
	param.Extra9 = "extra9"
	param.Extra10 = "extra10"
	param.RequestID = "requestId"
	param.ClientIP = "192.168.1.1"
	param.IMEI = "imei"
	param.MAC = "mac"
	param.AndroidID = "androidID"
	param.ServerIP = "serverIP"
	param.GAID = "gaid"
	param.IDFA = "idfa"
	param.AppVersionName = "appversionName"
	param.Brand = "brand"
	param.RemoteIP = "remoteIp"
	param.SessionID = "sessionID"
	param.ParentSessionID = "parentSessionID"
	param.CityCode = 3256
	param.Extra13 = 13
	param.Extra14 = 14
	param.Extra15 = 15
	param.Extra16 = 16
	param.IDFV = "idfv"
	param.OpenIDFA = "openIdfv"
	param.Extra20 = "extra20"
	param.Extfinalsubid = 10
	param.ExtpackageName = "extpackageName"
	param.ExtflowTagId = 100
	param.Extendcard = "extendcard"
	param.ExtrushNoPre = 1
	param.ExtfinalPackageName = "extfinalPackageName"
	param.Extnativex = 2
	param.Extctype = 42
	param.Extrvtemplate = 1
	param.Extabtest1 = 1
	param.Extb2t = 1
	param.Extchannel = "extChannel"
	param.Extbp = "extbp"
	param.Extsource = 32
	param.Extalgo = "extalgo"
	param.ExtthirdCid = "extthirdCid"
	param.ExtifLowerImp = 0
	param.ExtsystemUseragent = "ExtsystemUseragent"
	param.ExtMpNormalMap = "extmpnormalMap"
	Convey("SerializeOImpP equal", t, func() {
		res1 := SerializeOImpP(param)
		//res2 := SerializeOImpPOld(param)
		So(res1, ShouldNotBeEmpty)
	})
}

func BenchmarkSerializeOImpP(b *testing.B) {
	param := &Params{}
	param.PublisherID = 123
	param.AppID = 456
	param.UnitID = 345
	param.AdvCreativeID = 934
	param.CampaignID = 23456
	param.Scenario = "open_api"
	param.AdType = 42
	param.ImageSize = "imageSize"
	param.RequestType = 4
	param.PlatformName = "ios"
	param.OSVersion = "10.13.2"
	param.SDKVersion = "1.2"
	param.Model = "model"
	param.ScreenSize = "screenSize"
	param.Orientation = 2
	param.CountryCode = "HK"
	param.Language = "zh"
	param.NetworkTypeName = "wifi"
	param.MCC = "460"
	param.MNC = "01"
	param.Extra = "extra"
	param.Extra3 = "extra3"
	param.Extra4 = "extra4"
	param.Extra5 = "extra5"
	param.Extra7 = 7
	param.Extra8 = 8
	param.Extra9 = "extra9"
	param.Extra10 = "extra10"
	param.RequestID = "requestId"
	param.ClientIP = "192.168.1.1"
	param.IMEI = "imei"
	param.MAC = "mac"
	param.AndroidID = "androidID"
	param.ServerIP = "serverIP"
	param.GAID = "gaid"
	param.IDFA = "idfa"
	param.AppVersionName = "appversionName"
	param.Brand = "brand"
	param.RemoteIP = "remoteIp"
	param.SessionID = "sessionID"
	param.ParentSessionID = "parentSessionID"
	param.CityCode = 3256
	param.Extra13 = 13
	param.Extra14 = 14
	param.Extra15 = 15
	param.Extra16 = 16
	param.IDFV = "idfv"
	param.OpenIDFA = "openIdfv"
	param.Extra20 = "extra20"
	param.Extfinalsubid = 10
	param.ExtpackageName = "extpackageName"
	param.ExtflowTagId = 100
	param.Extendcard = "extendcard"
	param.ExtrushNoPre = 1
	param.ExtfinalPackageName = "extfinalPackageName"
	param.Extnativex = 2
	param.Extctype = 42
	param.Extrvtemplate = 1
	param.Extabtest1 = 1
	param.Extb2t = 1
	param.Extchannel = "extChannel"
	param.Extbp = "extbp"
	param.Extsource = 32
	param.Extalgo = "extalgo"
	param.ExtthirdCid = "extthirdCid"
	param.ExtifLowerImp = 0
	param.ExtsystemUseragent = "ExtsystemUseragent"
	param.ExtMpNormalMap = "extmpnormalMap"
	b.Run("SeralizeOImp", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			SerializeOImpP(param)
		}
	})
	// b.Run("SeralizeOImpOld", func(b *testing.B) {
	// 	for i := 0; i < b.N; i++ {
	// 		SerializeOImpPOld(param)
	// 	}
	// })
}

func TestSerializeOImpP(t *testing.T) {
	Convey("", t, func() {
		params := &Params{
			PublisherID: 100,
			AppID:       1000,
			UnitID:      10000,
			AdType:      42,
			RequestType: 1,
			Extendcard:  "test",
			Extra3:      "extra3",
			Extra4:      "extra4",
		}
		res := SerializeOImpP(params)
		So(res, ShouldNotBeEmpty)
		//So(res, ShouldEqual, "MTAwfDEwMDB8MTAwMDB8MHwwfHx8bmF0aXZlfHwxfHx8fHx8MHx8fHx8fHxleHRyYTN8ZXh0cmE0fHx8MHwwfHx8fHx8fHx8fHx8fHx8fHx8fHwwfDB8MHwwfDB8LHx8fHx8MHx8fHx8fDB8fHRlc3R8MHx8fDB8fHx8MHwwfDB8fHx8fHwwfHx8fHx8MHx8fHwwfHx8fHx8fHx8")
	})
}

func TestGetMPSdkVersionCompare(t *testing.T) {
	Convey("非mp 的sdkversion", t, func() {
		res := GetMPSdkVersionCompare("mal_2.4.6")
		So(res, ShouldBeFalse)
	})

	Convey("mp 的sdkversion 版本小于4.xxx", t, func() {
		res := GetMPSdkVersionCompare("mp_2.4.6")
		So(res, ShouldBeTrue)
	})

	Convey("mp 的sdkversion 版本大于4.xxx", t, func() {
		res := GetMPSdkVersionCompare("mp_4.4.6")
		So(res, ShouldBeFalse)
	})

	Convey("mp 的sdkversion 错误版本", t, func() {
		res := GetMPSdkVersionCompare("mp_")
		So(res, ShouldBeFalse)
	})

	Convey("mp 的sdkversion 错误版本无前缀", t, func() {
		res := GetMPSdkVersionCompare("1.2.3")
		So(res, ShouldBeFalse)
	})
}

func BenchmarkSerializeQ(b *testing.B) {
	jsonstr := "{\"campaignId\":230414590,\"apkUrl\":\"http://cdn.salmonads.com/apk/dmp/18/06/22/12/34/5b2c7c72e0805.apk\",\"advertiserId\":937,\"name\":\"mobvistawebeye_webeye_48207405\",\"platform\":1,\"trackingUrl\":\"http://track.clickhubs.com/v1/ad/click?h=1530587813084\\u0026pubid=10396\\u0026campid=48207405\\u0026gaid={gaid}\\u0026sub={clickId}\",\"directUrl\":\"\",\"price\":5.25,\"oriPrice\":5.25,\"startTime\":1529395261,\"endTime\":2398348800,\"countryCode\":[\"HK\"],\"cityCode\":[],\"status\":1,\"network\":121,\"previewUrl\":\"https://play.google.com/store/apps/details?id=com.pixonic.wwr\",\"packageName\":\"com.pixonic.wwr\",\"campaignType\":2,\"specialType\":[0],\"ctype\":1,\"networkCid\":\"48207405\",\"deviceModel\":[\"all\"],\"networkType\":[0],\"osVersion\":[\"4.0\",\"4.1\",\"4.2\",\"4.3\",\"4.4\",\"5.0\",\"5.1\",\"6.0\",\"7.0\",\"99.0\"],\"osVersionMin\":14,\"osVersionMax\":1000,\"osVersionMinV2\":4000000,\"osVersionMaxV2\":99000000,\"appName\":\"War Robots_HK_Android_CPI- CR 0.5- Need Transparency\",\"appDesc\":\"Win the Great Iron War! New 3D mech robot shooter in PvP mode!\",\"iconUrl\":\"http://cdn-adn.rayjump.com/cdn-adn/dmp/18/03/24/12/23/5ab5d2b426589.png\",\"appSize\":\"32\",\"appScore\":4.6,\"appInstall\":50000000,\"direct\":2,\"tag\":1,\"adSourceId\":1,\"category\":1,\"publisherId\":0,\"preClickRate\":{\"default\":[0,10000]},\"preClickCacheTime\":604800,\"frequencyCap\":0,\"directPackageName\":\"type9\",\"sdkPackageName\":\"\",\"landingType\":3,\"ctatext\":\"Install\",\"tImp\":0,\"advImp\":[],\"adUrlList\":[],\"jumpType\":9,\"isNoPayment\":2,\"newVersion\":\"4.0.0\",\"vbaConnecting\":2,\"vbaTrackingLink\":\"\",\"contentRating\":0,\"system\":[3,5],\"retargetingDevice\":2,\"sendDeviceidRate\":100,\"nxAdvName\":\"\",\"advOfferName\":\"\",\"updated\":1530588634,\"_id\":\"5b28b83e3728bc46ea412df9\",\"belongType\":0,\"vtaJump\":1,\"vtaTime\":2,\"retargetOffer\":2,\"openType\":0,\"subCategoryName\":[\"GAME_ACTION\"],\"isCampaignCreative\":0,\"costType\":1,\"source\":589,\"JUMP_TYPE_CONFIG\":{},\"chnId\":0,\"thirdParty\":\"\"}"
	campaignInfo := &smodel.CampaignInfo{}
	json.Unmarshal([]byte(jsonstr), campaignInfo)
	param := &Params{}
	param.NetworkType = mvconst.NETWORK_TYPE_4G
	param.Platform = mvconst.PlatformIOS
	param.RequestID = "5b167141c6c1e24c7cf21b9b"
	param.IDFA = "35951D6A-DB66-47A1-8CD0-DB1E9E0D33CD"
	param.ClientIP = "192.168.1.1"
	param.UserAgent = "Mozilla%252f5.0%2B%28iPhone%25253B%2BCPU%2BiPhone%2BOS%2B9_3_5%2Blike%2BMac%2BOS%2BX%29%2BAppleWebKit%252f601.1.46%2B%28KHTML%25252C%2Blike%2BGecko%29%2BMobile%252f13G36"
	param.TNum = 1
	param.OSVersion = "9.3.5"
	param.Model = "iphone6%2C2"
	param.UnitID = 128
	param.AppName = "iOSTest"
	param.PackageName = "com.tianye1.mvsdk"
	param.FormatAdType = mvconst.ADTypeNativeVideo
	param.AdType = mvconst.ADTypeNative
	param.VideoVersion = "1.1"
	param.VideoH = 300
	param.VideoW = 400

	for i := 0; i < b.N; i++ {
		_ = SerializeQ(param, campaignInfo)
	}
}

func TestSerializeCSPNew(t *testing.T) {

	param := &Params{}
	param.MWadBackend = "8,2"
	param.MWadBackendData = "1:2333334:1"
	param.MWFlowTagID = 102
	param.MWRandValue = 3721
	param.MWbackendConfig = "hello,world"
	param.AdNum = 20
	param.TNum = 5
	Convey("SerializeCSP equal", t, func() {
		res1 := SerializeCSP(param, "")
		//res2 := SerializeCSPOld(param)
		So(res1, ShouldNotBeEmpty)
	})
}

// func BenchmarkIntVer(b *testing.B) {
// 	b.Run("intver", func(b *testing.B) {
// 		for i := 0; i < b.N; i++ {
// 			IntVer("10.1.2.4")
// 		}
// 	})
// 	b.Run("old intver", func(b *testing.B) {
// 		for i := 0; i < b.N; i++ {
// 			IntVerOld("10.1.2.4")
// 		}
// 	})
// }

type IntRate int

func (i IntRate) GetRate() int {
	return int(i)
}

func TestRandByRateInMapV2(t *testing.T) {

	type args struct {
		valInMap map[string]IRate
		getRate  func(sum int) int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "check sum",
			args: args{
				valInMap: map[string]IRate{
					"abc": IntRate(100),
				},
				getRate: func(sum int) int {
					return rand.Intn(sum)
				},
			},
			want: "abc",
		},
		{
			name: "check sum 2",
			args: args{
				valInMap: map[string]IRate{
					"abc": IntRate(100),
					"bcd": IntRate(0),
				},
				getRate: func(sum int) int {
					return rand.Intn(sum)
				},
			},
			want: "abc",
		},
		{
			name: "check sum 3",
			args: args{
				valInMap: map[string]IRate{
					"abc": IntRate(0),
					"bcd": IntRate(100),
				},
				getRate: func(sum int) int {
					return rand.Intn(sum)
				},
			},
			want: "bcd",
		},
		{
			name: "check sum 4",
			args: args{
				valInMap: map[string]IRate{
					"abc": IntRate(0),
					"bcd": IntRate(0),
				},
				getRate: func(sum int) int {
					if sum <= 0 {
						return 0
					}
					return rand.Intn(sum)
				},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandByRateInMapV2(tt.args.valInMap, tt.args.getRate); got != tt.want {
				t.Errorf("RandByRateInMapV2() = %v, want %v", got, tt.want)
			}
		})
	}
}
