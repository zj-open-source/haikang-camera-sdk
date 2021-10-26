package sdk_camera

import (
	"errors"
	"github.com/zjzjzjzj1874/haikang-camera-sdk/constants"
	"unsafe"
)

type CameraSDKType int

const (
	TypeofHikvisionCameraSDK CameraSDKType = iota
)

var DefaultCameraSDK = NewCameraSDK()

type (
	newSDK interface {
		New() CameraSDK
	}

	CameraSDK interface {
		GetSDKVersion() uint32                                 // 获取SDK版本
		GetDeviceList() (interface{}, error)                   // 获取摄像头列表
		CreateHandle(idx int) error                            // 创建摄像头句柄
		DestroyHandle() error                                  // 销毁句柄
		OpenDevice() error                                     // 打开摄像头
		CloseDevice() error                                    // 关闭摄像头
		StartGrabbing(triggerMode constants.TriggerMode) error // 开始取流
		StopGrabbing() error                                   // 停止取流
		// GetOneFrameTimeout 获取一帧数据:bmp
		GetOneFrameTimeout(timeout uint32) (interface{}, []byte, error)
		// GetImageForRGB 获取一帧RGB数据
		GetImageForRGB(timeout uint32) (res interface{}, dataBuf []byte, err error)
		// GetOneFrameWithCallback 回调取图
		GetOneFrameWithCallback(argsAddr unsafe.Pointer) error
		// GetOneFrameFroRGBWithCallback 回调取图(RGB)
		GetOneFrameFroRGBWithCallback(argsAddr unsafe.Pointer) error
		// GetOneFrameFroBGRWithCallback 回调取图(BGR)
		GetOneFrameFroBGRWithCallback(argsAddr unsafe.Pointer) error
		// FeatureSave 将相机属性导出到文件中
		FeatureSave(path string) error
		// FeatureLoad 从文件中导入相机属性
		FeatureLoad(path string) error
		// GetHeartBeatTimeout 相机心跳检测
		GetHeartBeatTimeout() error
		// IsDeviceConnected 相机是否连接
		IsDeviceConnected() bool
	}
)

type SDKSet map[CameraSDKType]newSDK

func NewCameraSDK() SDKSet {
	return map[CameraSDKType]newSDK{}
}

func (s SDKSet) Register(sdkType CameraSDKType, sdk newSDK) {
	if s[sdkType] != nil {
		panic("sdk is existing")
	}
	s[sdkType] = sdk
}

func GetCameraSDK(sdkType CameraSDKType) (CameraSDK, error) {
	if DefaultCameraSDK[sdkType] == nil {
		return nil, errors.New("camera sdk is not exist")
	}
	return DefaultCameraSDK[sdkType].New(), nil
}
