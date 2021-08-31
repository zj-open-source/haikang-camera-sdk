package constants

import (
	bytes "bytes"
	database_sql_driver "database/sql/driver"
	errors "errors"

	github_com_go_courier_enumeration "github.com/go-courier/enumeration"
)

var InvalidTriggerMode = errors.New("invalid TriggerMode type")

func ParseTriggerModeFromLabelString(s string) (TriggerMode, error) {
	switch s {
	case "":
		return TRIGGER_MODE__OFF, nil
	case "关闭触发模式":
		return TRIGGER_MODE__OFF, nil
	case "打开触发模式":
		return TRIGGER_MODE__ON, nil
	}
	return TRIGGER_MODE__OFF, InvalidTriggerMode
}

func (v TriggerMode) String() string {
	switch v {
	case TRIGGER_MODE__OFF:
		return "OFF"
	case TRIGGER_MODE__ON:
		return "ON"
	}
	return "UNKNOWN"
}

func ParseTriggerModeFromString(s string) (TriggerMode, error) {
	switch s {
	case "":
		return TRIGGER_MODE__OFF, nil
	case "OFF":
		return TRIGGER_MODE__OFF, nil
	case "ON":
		return TRIGGER_MODE__ON, nil
	}
	return TRIGGER_MODE__OFF, InvalidTriggerMode
}

func (v TriggerMode) Label() string {
	switch v {
	case TRIGGER_MODE__OFF:
		return "关闭触发模式"
	case TRIGGER_MODE__ON:
		return "打开触发模式"
	}
	return "UNKNOWN"
}

func (v TriggerMode) Int() int {
	return int(v)
}

func (TriggerMode) TypeName() string {
	return "git.querycap.com/aia/env-sdk-camera/constants.TriggerMode"
}

func (TriggerMode) ConstValues() []github_com_go_courier_enumeration.IntStringerEnum {
	return []github_com_go_courier_enumeration.IntStringerEnum{TRIGGER_MODE__OFF, TRIGGER_MODE__ON}
}

func (v TriggerMode) MarshalText() ([]byte, error) {
	str := v.String()
	if str == "UNKNOWN" {
		return nil, InvalidTriggerMode
	}
	return []byte(str), nil
}

func (v *TriggerMode) UnmarshalText(data []byte) (err error) {
	*v, err = ParseTriggerModeFromString(string(bytes.ToUpper(data)))
	return
}

func (v TriggerMode) Value() (database_sql_driver.Value, error) {
	offset := 0
	if o, ok := (interface{})(v).(github_com_go_courier_enumeration.DriverValueOffset); ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}

func (v *TriggerMode) Scan(src interface{}) error {
	offset := 0
	if o, ok := (interface{})(v).(github_com_go_courier_enumeration.DriverValueOffset); ok {
		offset = o.Offset()
	}

	i, err := github_com_go_courier_enumeration.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*v = TriggerMode(i)
	return nil
}
