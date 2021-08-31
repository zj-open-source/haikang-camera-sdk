package constants

//go:generate tools gen enum TriggerMode
type TriggerMode uint8

// TriggerMode 触发模式
const (
	TRIGGER_MODE__OFF TriggerMode = iota // 关闭触发模式
	TRIGGER_MODE__ON                     // 打开触发模式
)
