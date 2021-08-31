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
*/
import "C"
import (
	"fmt"
	gopointer "github.com/mattn/go-pointer"
	"os"
	"time"
	"unsafe"
)

// 回调函数取图
func grabImageWithCallback() {
	sdkVersion := C.MV_CC_GetSDKVersion()
	fmt.Println("SDKVersion:", sdkVersion)

	deviceList := &C.MV_CC_DEVICE_INFO_LIST{}
	var tlayerType C.uint = C.MV_GIGE_DEVICE | C.MV_USB_DEVICE

	//ch:枚举设备 | en:Enum device
	var ret C.int
	ret = C.MV_CC_EnumDevices(tlayerType, deviceList)
	if ret != 0 {
		fmt.Println("return signal is not zero:", fmt.Sprintf("0x%x", C.uint(ret)))
		return
	}
	fmt.Println("device num:", deviceList.nDeviceNum)

	handle := unsafe.Pointer(nil)
	stDeviceList := deviceList.pDeviceInfo[0]

	//  指定第一个设备并创建句柄
	ret = C.MV_CC_CreateHandle(&handle, stDeviceList)
	if ret != 0 {
		fmt.Println("MV_CC_CreateHandle failure:", fmt.Sprintf("0x%x", C.uint(ret)))
		return
	}
	defer func() {
		ret = C.MV_CC_DestroyHandle(handle)
		if ret != 0 {
			fmt.Println("MV_CC_DestroyHandle failure:", fmt.Sprintf("0x%x", C.uint(ret)))
		} else {
			fmt.Println("MV_CC_DestroyHandle success")
		}
	}()
	fmt.Println("handle:", handle)

	// 打开设备
	ret = C.MV_CC_OpenDevice(handle, C.MV_ACCESS_Exclusive, 0)
	if ret != 0 {
		fmt.Println("MV_CC_OpenDevice failure:", fmt.Sprintf("0x%x", C.uint(ret)))
		return
	}
	defer func() {
		ret = C.MV_CC_CloseDevice(handle)
		if ret != 0 {
			fmt.Println("MV_CC_CloseDevice failure:", fmt.Sprintf("0x%x", C.uint(ret)))
			return
		} else {
			fmt.Println("MV_CC_CloseDevice success")
		}
	}()

	// 探测网络最佳包大小(只对GigE相机有效)
	if stDeviceList.nTLayerType == C.MV_GIGE_DEVICE {
		var nPacketSize C.int = C.MV_CC_GetOptimalPacketSize(handle)
		if nPacketSize > 0 {
			ret = C.MV_CC_SetIntValue(handle, C.CString("GevSCPSPacketSize"), C.uint(nPacketSize))
			if ret != 0 {
				fmt.Println("MV_CC_SetIntValue failure:", fmt.Sprintf("0x%x", C.uint(ret)))
				return
			}
		} else {
			fmt.Println("MV_CC_GetOptimalPacketSize failure:", nPacketSize)
		}
	}

	// ch:设置触发模式为off | en:Set trigger mode as off
	ret = C.MV_CC_SetEnumValue(handle, C.CString("TriggerMode"), C.MV_TRIGGER_MODE_OFF)
	if ret != 0 {
		fmt.Println("MV_CC_SetEnumValue failure:", fmt.Sprintf("0x%x", ret))
		return
	}

	// ch:获取数据包大小 | en:Get payload size
	stParam := &C.MVCC_INTVALUE{}
	ret = C.MV_CC_GetIntValue(handle, C.CString("PayloadSize"), stParam)
	if ret != 0 {
		fmt.Println("MV_CC_GetIntValue failure:", fmt.Sprintf("0x%x", C.uint(ret)))
		return
	}

	// 先注册回调函数
	cb := &CallbackStruct{
		PortID:     "测试摄像头1",
		FamilyName: "测试摄像头family",
		Callback: func(frameInfo FrameOutInfo) {
			//TODO 元哥，你的回调函数写这里
			fmt.Printf("get one frame:[PortID:%v], [DataBuf len:%d], Width[%v], Height[%v], nFrameNum[%v]\n", frameInfo.PortID, len(frameInfo.DataBuf), frameInfo.Width, frameInfo.Height, frameInfo.FrameNum)
		},
	}
	p := gopointer.Save(cb)
	defer gopointer.Unref(p)

	ret = C.MV_CC_RegisterImageCallBackEx(handle, (*[0]byte)(unsafe.Pointer(C.Callback)), p)
	if ret != 0 {
		fmt.Println("MV_CC_RegisterImageCallBackEx failure:", fmt.Sprintf("0x%x", C.uint(ret)))
		return
	}

	// 开始取流
	ret = C.MV_CC_StartGrabbing(handle)
	if ret != 0 {
		fmt.Println("MV_CC_StartGrabbing failure:", fmt.Sprintf("0x%x", C.uint(ret)))
		return
	}
	defer func() {
		ret = C.MV_CC_StopGrabbing(handle)
		if ret != 0 {
			fmt.Println("MV_CC_StopGrabbing failure:", fmt.Sprintf("0x%x", C.uint(ret)))
			return
		} else {
			fmt.Println("MV_CC_StopGrabbing success")
		}
	}()

	exitSignal := make(chan os.Signal, 1)
	i := 0
	for {
		time.Sleep(time.Second)
		fmt.Println(i)
		i++
	}
	<-exitSignal

}
