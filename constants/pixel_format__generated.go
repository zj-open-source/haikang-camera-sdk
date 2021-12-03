package constants

import (
	bytes "bytes"
	database_sql_driver "database/sql/driver"
	errors "errors"

	github_com_go_courier_enumeration "github.com/go-courier/enumeration"
)

var InvalidPixelFormat = errors.New("invalid PixelFormat type")

func ParsePixelFormatFromLabelString(s string) (PixelFormat, error) {
	switch s {
	case "":
		return PIXEL_FORMAT_UNKNOWN, nil
	case "mono8":
		return PIXEL_FORMAT__MONO8, nil
	case "bayer8":
		return PIXEL_FORMAT__BAYER8, nil
	case "yuv":
		return PIXEL_FORMAT__YUV, nil
	case "rgb":
		return PIXEL_FORMAT__RGB, nil
	case "bgr":
		return PIXEL_FORMAT__BGR, nil
	}
	return PIXEL_FORMAT_UNKNOWN, InvalidPixelFormat
}

func (v PixelFormat) String() string {
	switch v {
	case PIXEL_FORMAT_UNKNOWN:
		return ""
	case PIXEL_FORMAT__MONO8:
		return "MONO8"
	case PIXEL_FORMAT__BAYER8:
		return "BAYER8"
	case PIXEL_FORMAT__YUV:
		return "YUV"
	case PIXEL_FORMAT__RGB:
		return "RGB"
	case PIXEL_FORMAT__BGR:
		return "BGR"
	}
	return "UNKNOWN"
}

func ParsePixelFormatFromString(s string) (PixelFormat, error) {
	switch s {
	case "":
		return PIXEL_FORMAT_UNKNOWN, nil
	case "MONO8":
		return PIXEL_FORMAT__MONO8, nil
	case "BAYER8":
		return PIXEL_FORMAT__BAYER8, nil
	case "YUV":
		return PIXEL_FORMAT__YUV, nil
	case "RGB":
		return PIXEL_FORMAT__RGB, nil
	case "BGR":
		return PIXEL_FORMAT__BGR, nil
	}
	return PIXEL_FORMAT_UNKNOWN, InvalidPixelFormat
}

func (v PixelFormat) Label() string {
	switch v {
	case PIXEL_FORMAT_UNKNOWN:
		return ""
	case PIXEL_FORMAT__MONO8:
		return "mono8"
	case PIXEL_FORMAT__BAYER8:
		return "bayer8"
	case PIXEL_FORMAT__YUV:
		return "yuv"
	case PIXEL_FORMAT__RGB:
		return "rgb"
	case PIXEL_FORMAT__BGR:
		return "bgr"
	}
	return "UNKNOWN"
}

func (v PixelFormat) Int() int {
	return int(v)
}

func (PixelFormat) TypeName() string {
	return "github.com/zjzjzjzj1874/haikang-camera-sdk/constants.PixelFormat"
}

func (PixelFormat) ConstValues() []github_com_go_courier_enumeration.IntStringerEnum {
	return []github_com_go_courier_enumeration.IntStringerEnum{PIXEL_FORMAT__MONO8, PIXEL_FORMAT__BAYER8, PIXEL_FORMAT__YUV, PIXEL_FORMAT__RGB, PIXEL_FORMAT__BGR}
}

func (v PixelFormat) MarshalText() ([]byte, error) {
	str := v.String()
	if str == "UNKNOWN" {
		return nil, InvalidPixelFormat
	}
	return []byte(str), nil
}

func (v *PixelFormat) UnmarshalText(data []byte) (err error) {
	*v, err = ParsePixelFormatFromString(string(bytes.ToUpper(data)))
	return
}

func (v PixelFormat) Value() (database_sql_driver.Value, error) {
	offset := 0
	if o, ok := (interface{})(v).(github_com_go_courier_enumeration.DriverValueOffset); ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}

func (v *PixelFormat) Scan(src interface{}) error {
	offset := 0
	if o, ok := (interface{})(v).(github_com_go_courier_enumeration.DriverValueOffset); ok {
		offset = o.Offset()
	}

	i, err := github_com_go_courier_enumeration.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*v = PixelFormat(i)
	return nil
}
