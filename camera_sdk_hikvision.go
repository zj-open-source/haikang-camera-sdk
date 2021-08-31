// Package sdk_camera 海康摄像头SDK集成
package sdk_camera

/*
#cgo CFLAGS: -I./include/
#cgo LDFLAGS: -L./lib/64/ -lMvCameraControl -ldl
#include "CameraParams.h"
#include "MvCameraControl.h"
#include "callback.h"
#include <stdlib.h>
#include <string.h>
#include <stdio.h>

struct _MV_CamL_DEV_INFO_ MV_CC_DEIVCE_INFO_GET_stCamLInFo(struct _MV_CC_DEVICE_INFO_* s) {
	return s->SpecialInfo.stCamLInfo;
}

struct _MV_USB3_DEVICE_INFO_ MV_CC_DEIVCE_INFO_GET_stUsb3VInfo(struct _MV_CC_DEVICE_INFO_* s) {
	return s->SpecialInfo.stUsb3VInfo;
}

struct _MV_GIGE_DEVICE_INFO_ MV_CC_DEIVCE_INFO_GET_stGigEInfo(struct _MV_CC_DEVICE_INFO_* s) {
	return s->SpecialInfo.stGigEInfo;
}
*/
import "C"
import (
	"errors"
	"fmt"
	"git.querycap.com/aia/env-sdk-camera/constants"
	"unsafe"
)

func init() {
	DefaultCameraSDK.Register(TypeofHikvisionCameraSDK, &hikvisionCameraSDK{})
}

type (
	hikvisionCameraSDK struct{}
	HikvisionCameraSDK struct {
		handle     unsafe.Pointer
		deviceInfo *C.MV_CC_DEVICE_INFO
		deviceList C.MV_CC_DEVICE_INFO_LIST
		nDataSize  C.uint
	}
)

func (h *HikvisionCameraSDK) SetDeviceList(old *HikvisionCameraSDK) {
	h.deviceList = old.deviceList
}

// GetOneFrameFroRGBWithCallback 入参：传入回调函数
func (h *HikvisionCameraSDK) GetOneFrameFroRGBWithCallback(addr unsafe.Pointer) error {
	ret := C.MV_CC_RegisterImageCallBackForRGB(h.handle, (*[0]byte)(unsafe.Pointer(C.Callback)), addr)
	if ret != 0 {
		fmt.Println("MV_CC_RegisterImageCallBackEx failure:", fmt.Sprintf("0x%x", C.uint(ret)))
		return errors.New(fmt.Sprintf("MV_CC_RegisterImageCallBackEx failure:0x%x", C.uint(ret)))
	}
	return nil
}

// GetOneFrameFroBGRWithCallback 入参：传入回调函数
func (h *HikvisionCameraSDK) GetOneFrameFroBGRWithCallback(addr unsafe.Pointer) error {
	ret := C.MV_CC_RegisterImageCallBackForBGR(h.handle, (*[0]byte)(unsafe.Pointer(C.Callback)), addr)
	if ret != 0 {
		fmt.Println("MV_CC_RegisterImageCallBackEx failure:", fmt.Sprintf("0x%x", C.uint(ret)))
		return errors.New(fmt.Sprintf("MV_CC_RegisterImageCallBackEx failure:0x%x", C.uint(ret)))
	}
	return nil
}

// GetOneFrameWithCallback 入参：传入回调函数
func (h *HikvisionCameraSDK) GetOneFrameWithCallback(addr unsafe.Pointer) error {
	ret := C.MV_CC_RegisterImageCallBackEx(h.handle, (*[0]byte)(unsafe.Pointer(C.Callback)), addr)
	if ret != 0 {
		fmt.Println("MV_CC_RegisterImageCallBackEx failure:", fmt.Sprintf("0x%x", C.uint(ret)))
		return errors.New(fmt.Sprintf("MV_CC_RegisterImageCallBackEx failure:0x%x", C.uint(ret)))
	}
	return nil
}

// GetOneFrameTimeout timeout 默认1000
func (h *HikvisionCameraSDK) GetOneFrameTimeout(timeout uint32) (interface{}, []byte, error) {
	pFrameInfo := C.MV_FRAME_OUT_INFO_EX{}
	dataBuf := make([]byte, h.nDataSize*h.deviceList.nDeviceNum)
	ret := C.MV_CC_GetOneFrameTimeout(h.handle, (*C.uchar)(unsafe.Pointer(&dataBuf[0])), h.nDataSize*h.deviceList.nDeviceNum, &pFrameInfo, C.uint(timeout))
	if ret == 0 {
		fmt.Printf("get one frame with timeout: Width[%v], Height[%v], PixelType[0x%v], nFrameNum[%v]\n", pFrameInfo.nWidth, pFrameInfo.nHeight, pFrameInfo.enPixelType, pFrameInfo.nFrameNum)
	} else {
		fmt.Println("no data:", fmt.Sprintf("0x%x", C.uint(ret)))
		return nil, nil, nil
	}

	return &MvFrameOutInfoEx{
		Width:     uint16(pFrameInfo.nWidth),
		Height:    uint16(pFrameInfo.nHeight),
		PixelType: uint16(pFrameInfo.enPixelType),
		FrameNum:  uint32(pFrameInfo.nFrameNum),
		FrameLen:  uint32(pFrameInfo.nFrameLen),
	}, dataBuf, nil
}

// GetImageForRGB ：timeout 默认1000
func (h *HikvisionCameraSDK) GetImageForRGB(timeout uint32) (interface{}, []byte, error) {
	pFrameInfo := C.MV_FRAME_OUT_INFO_EX{}
	dataBuf := make([]byte, h.nDataSize*h.deviceList.nDeviceNum)
	ret := C.MV_CC_GetImageForRGB(h.handle, (*C.uchar)(unsafe.Pointer(&dataBuf[0])), h.nDataSize*h.deviceList.nDeviceNum, &pFrameInfo, C.int(timeout))
	if ret == 0 {
		fmt.Printf("get image for RGB: Width[%v], Height[%v], PixelType[0x%v], nFrameNum[%v]\n", pFrameInfo.nWidth, pFrameInfo.nHeight, pFrameInfo.enPixelType, pFrameInfo.nFrameNum)
	} else {
		fmt.Println("no data:", fmt.Sprintf("0x%x", C.uint(ret)))
		return nil, nil, nil
	}

	return &MvFrameOutInfoEx{
		Width:     uint16(pFrameInfo.nWidth),
		Height:    uint16(pFrameInfo.nHeight),
		PixelType: uint16(pFrameInfo.enPixelType),
		FrameNum:  uint32(pFrameInfo.nFrameNum),
		FrameLen:  uint32(pFrameInfo.nFrameLen),
	}, dataBuf, nil
}

func (h *HikvisionCameraSDK) GetDeviceList() (interface{}, error) {
	var (
		tLayerType C.uint = C.MV_GIGE_DEVICE | C.MV_USB_DEVICE
		deviceList        = &C.MV_CC_DEVICE_INFO_LIST{}
	)
	//ch:枚举设备 | en:Enum device
	var ret C.int
	ret = C.MV_CC_EnumDevices(tLayerType, deviceList)
	if ret != 0 {
		fmt.Println("MV_CC_EnumDevices failure:", fmt.Sprintf("0x%x", C.uint(ret)))
		return nil, errors.New(fmt.Sprintf("MV_CC_EnumDevices failure:0x%x", C.uint(ret)))
	}

	h.deviceList = *deviceList
	deviceInfoList := &MvCcDeviceInfoList{
		DeviceNum: uint32(deviceList.nDeviceNum),
	}

	infos := make([]MvCcDeviceInfo, 0)
	for i := 0; i < int(deviceList.nDeviceNum); i++ {
		stCamLInfo := C.MV_CC_DEIVCE_INFO_GET_stCamLInFo(deviceList.pDeviceInfo[i])
		familyName := C.GoStringN((*C.char)(unsafe.Pointer(&stCamLInfo.chFamilyName[0])), 64)
		portID := C.GoStringN((*C.char)(unsafe.Pointer(&stCamLInfo.chPortID[0])), 64)
		mcdi := MvCcDeviceInfo{
			TLayerType: uint32(deviceList.pDeviceInfo[i].nTLayerType),
			FamilyName: familyName,
			PortID:     portID,
		}
		infos = append(infos, mcdi)

	}

	deviceInfoList.MvCcDeviceInfo = &infos
	return deviceInfoList, nil
}

func (h *HikvisionCameraSDK) CreateHandle(idx int) error {
	if h.deviceList.nDeviceNum < C.uint(idx+1) {
		return errors.New("invalid index")
	}
	handle := unsafe.Pointer(nil)
	ret := C.MV_CC_CreateHandle(&handle, h.deviceList.pDeviceInfo[idx])
	if ret != 0 {
		fmt.Println("MV_CC_CreateHandle failure:", fmt.Sprintf("0x%x", C.uint(ret)))
		return errors.New(fmt.Sprintf("MV_CC_CreateHandle failure:0x%x", C.uint(ret)))
	}
	h.deviceInfo = h.deviceList.pDeviceInfo[idx]
	h.handle = handle
	return nil
}

func (h *HikvisionCameraSDK) DestroyHandle() error {
	ret := C.MV_CC_DestroyHandle(h.handle)
	if ret != 0 {
		fmt.Println("MV_CC_DestroyHandle failure:", fmt.Sprintf("0x%x", C.uint(ret)))
		return errors.New(fmt.Sprintf("MV_CC_DestroyHandle failure:0x%x", C.uint(ret)))
	}
	return nil
}

func (h *HikvisionCameraSDK) OpenDevice() error {
	ret := C.MV_CC_OpenDevice(h.handle, C.MV_ACCESS_Exclusive, 0)
	if ret != 0 {
		fmt.Println("MV_CC_OpenDevice failure:", fmt.Sprintf("0x%x", C.uint(ret)))
		return errors.New(fmt.Sprintf("MV_CC_OpenDevice failure:0x%x", C.uint(ret)))
	}
	return nil
}

func (h *HikvisionCameraSDK) CloseDevice() error {
	ret := C.MV_CC_CloseDevice(h.handle)
	if ret != 0 {
		fmt.Println("MV_CC_CloseDevice failure:", fmt.Sprintf("0x%x", C.uint(ret)))
		return errors.New(fmt.Sprintf("MV_CC_CloseDevice failure:0x%x", C.uint(ret)))
	} else {
		fmt.Println("MV_CC_CloseDevice success!!")
	}
	return nil
}

// StartGrabbing triggerModel = 1(开触发拍照模式);while triggerMode = 0(关闭触发拍照模式)
func (h *HikvisionCameraSDK) StartGrabbing(triggerMode constants.TriggerMode) error {
	// 探测网络最佳包大小(只对GigE相机有效)
	if h.deviceInfo.nTLayerType == C.MV_GIGE_DEVICE {
		var nPacketSize C.int = C.MV_CC_GetOptimalPacketSize(h.handle)
		if nPacketSize > 0 {
			ret := C.MV_CC_SetIntValue(h.handle, C.CString("GevSCPSPacketSize"), C.uint(nPacketSize))
			if ret != 0 {
				fmt.Println("MV_CC_SetIntValue failure:", fmt.Sprintf("0x%x", C.uint(ret)))
				return errors.New(fmt.Sprintf("MV_CC_SetIntValue failure:0x%x", C.uint(ret)))
			}
		} else {
			fmt.Println("MV_CC_GetOptimalPacketSize failure:", nPacketSize)
			return errors.New(fmt.Sprintf("MV_CC_GetOptimalPacketSize failure:nPacketSize:%v", nPacketSize))
		}
	}

	// ch:设置触发模式为off | en:Set trigger mode as off  ==>  C.MV_TRIGGER_MODE_OFF == 0  || C.MV_TRIGGER_MODE_ON == 1
	triggerModeKey := C.CString("TriggerMode")
	defer C.free(unsafe.Pointer(triggerModeKey))
	ret := C.MV_CC_SetEnumValue(h.handle, triggerModeKey, C.uint(triggerMode))
	//ret := C.MV_CC_SetEnumValue(h.handle, C.CString("TriggerMode"), C.MV_TRIGGER_MODE_OFF)
	if ret != 0 {
		fmt.Println("MV_CC_SetEnumValue failure:", fmt.Sprintf("0x%x", ret))
		return errors.New(fmt.Sprintf("MV_CC_SetEnumValue failure:0x%x", C.uint(ret)))
	}

	// ch:获取数据包大小 | en:Get payload size
	stParam := &C.MVCC_INTVALUE{}
	ret = C.MV_CC_GetIntValue(h.handle, C.CString("PayloadSize"), stParam)
	if ret != 0 {
		fmt.Println("MV_CC_GetIntValue failure:", fmt.Sprintf("0x%x", C.uint(ret)))
		return errors.New(fmt.Sprintf("MV_CC_GetIntValue failure:0x%x", C.uint(ret)))
	}

	// 开始取流
	ret = C.MV_CC_StartGrabbing(h.handle)
	if ret != 0 {
		fmt.Println("MV_CC_StartGrabbing failure:", fmt.Sprintf("0x%x", C.uint(ret)))
		return errors.New(fmt.Sprintf("MV_CC_StartGrabbing failure:0x%x", C.uint(ret)))
	}

	h.nDataSize = C.uint(stParam.nCurValue)
	return nil
}

func (h *HikvisionCameraSDK) StopGrabbing() error {
	ret := C.MV_CC_StopGrabbing(h.handle)
	if ret != 0 {
		fmt.Println("MV_CC_StopGrabbing failure:", fmt.Sprintf("0x%x", C.uint(ret)))
		return errors.New(fmt.Sprintf("MV_CC_StopGrabbing failure:0x%x", C.uint(ret)))
	}
	return nil
}

func (h *HikvisionCameraSDK) GetSDKVersion() uint32 {
	return uint32(C.MV_CC_GetSDKVersion())
}

func (h *hikvisionCameraSDK) New() CameraSDK {
	return &HikvisionCameraSDK{}
}

type (
	MvCcDeviceInfoList struct {
		DeviceNum      uint32            `json:"nDeviceNum"`  // 设备数量
		MvCcDeviceInfo *[]MvCcDeviceInfo `json:"pDeviceInfo"` // 设备列表
	}
	MvCcDeviceInfo struct {
		MajorVer    uint16 `json:"nMajorVer"`
		MinorVer    uint16 `json:"nMinorVer"`
		MacAddrHigh uint32 `json:"nMacAddrHigh"`
		MacAddrLow  uint32 `json:"nMacAddrLow"`
		TLayerType  uint32 `json:"nTLayerType"`
		FamilyName  string `json:"chFamilyName"`
		PortID      string `json:"chPortID"`
	}

	MvFrameOutInfoEx struct {
		Width             uint16      `json:"nWidth"`             // 宽
		Height            uint16      `json:"nHeight"`            // 高
		PixelType         uint16      `json:"enPixelType"`        // todo 枚举,对应什么数据结构
		FrameNum          uint32      `json:"nFrameNum"`          // 当前帧数
		DevTimeStampHigh  uint32      `json:"nDevTimeStampHigh"`  //
		DevTimeStampLow   uint32      `json:"nDevTimeStampLow"`   //
		Reserved0         uint32      `json:"nReserved0"`         //
		HostTimeStamp     uint        `json:"nHostTimeStamp"`     //
		FrameLen          uint32      `json:"nFrameLen"`          //
		SecondCount       uint32      `json:"nSecondCount"`       //
		CycleCount        uint32      `json:"nCycleCount"`        //
		CycleOffset       uint32      `json:"nCycleOffset"`       //
		Gain              float32     `json:"fGain"`              //
		ExposureTime      float32     `json:"fExposureTime"`      //
		AverageBrightness uint32      `json:"nAverageBrightness"` //
		Red               uint32      `json:"nRed"`               //
		Green             uint32      `json:"nGreen"`             //
		Blue              uint32      `json:"nBlue"`              //
		FrameCounter      uint32      `json:"nFrameCounter"`      //
		TriggerIndex      uint32      `json:"nTriggerIndex"`      //
		Input             uint32      `json:"nInput"`             //
		Output            uint32      `json:"nOutput"`            //
		OffsetX           uint16      `json:"nOffsetX"`           //
		OffsetY           uint16      `json:"nOffsetY"`           //
		ChunkWidth        uint16      `json:"nChunkWidth"`        //
		ChunkHeight       uint16      `json:"nChunkHeight"`       //
		LostPacket        uint32      `json:"nLostPacket"`        //
		UnparsedChunkNum  uint32      `json:"nUnparsedChunkNum"`  //
		UnparsedChunkList interface{} `json:"UnparsedChunkList"`  //
		Reserved          []uint32    `json:"nReserved"`          //
	}
)

type (
	CallbackStruct struct {
		FamilyName string       `json:"chFamilyName"`
		PortID     string       `json:"chPortID"`
		Callback   CallbackFunc `json:"callback"`
	}

	FrameOutInfo struct {
		// 图片宽
		Width uint16 `json:"width"`
		// 图片高
		Height uint16 `json:"height"`
		// 当前帧数
		FrameNum uint32 `json:"frameNum"`
		// 图片内容
		DataBuf    []byte `json:"dataBuf"`
		FamilyName string `json:"chFamilyName"`
		PortID     string `json:"chPortID"`
	}
)

type CallbackFunc func(frameInfo FrameOutInfo)
