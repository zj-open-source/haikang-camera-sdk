package sdk_camera

/*
#cgo CFLAGS: -I./include
#cgo LDFLAGS: -L./lib/64 -lMvCameraControl -ldl
#include "CameraParams.h"
#include "MvCameraControl.h"
#include <stdlib.h>
#include <string.h>
#include <stdio.h>
*/
import "C"
import (
	"fmt"
	"os"
	"time"
	"unsafe"
)

// 取图并保存
func grabImageWithSave() {
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
	nDataSize := C.uint(stParam.nCurValue)
	dataBuf := make([]byte, nDataSize)
	pFrameInfo := C.MV_FRAME_OUT_INFO_EX{}

	for {
		ret = C.MV_CC_GetOneFrameTimeout(handle, (*C.uchar)(unsafe.Pointer(&dataBuf[0])), nDataSize, &pFrameInfo, C.uint(1000))
		if ret == 0 {
			fmt.Printf("get one frame: Width[%v], Height[%v], PixelType[0x%v], nFrameNum[%v]\n", pFrameInfo.nWidth, pFrameInfo.nHeight, pFrameInfo.enPixelType, pFrameInfo.nFrameNum)
		} else {
			fmt.Println("no data:", fmt.Sprintf("0x%x", C.uint(ret)))
			continue
		}

		// 使用MV_CC_SaveImageEx2保存图片
		pSaveParam := &C.MV_SAVE_IMAGE_PARAM_EX{}
		db := dataBuf
		pSaveParam.pData = (*C.uchar)(unsafe.Pointer(&db[0]))

		fmt.Println(pSaveParam.pData)

		pSaveParam.nWidth = pFrameInfo.nWidth
		pSaveParam.nHeight = pFrameInfo.nHeight
		pSaveParam.nDataLen = pFrameInfo.nFrameLen
		pSaveParam.enPixelType = pFrameInfo.enPixelType
		//filePath := "/home/rock/save.bmp"
		pSaveParam.enImageType = C.MV_Image_Bmp
		//jpeg
		//filePath := "/home/rock/save.jpeg"
		//pSaveParam.enImageType = C.MV_Image_Jpeg
		//pSaveParam.nJpgQuality = C.uint(75)

		var bmpSize C.uint = C.uint(pFrameInfo.nWidth)*C.uint(pFrameInfo.nHeight)*3 + 54
		pSaveParam.nBufferSize = bmpSize
		bmpBuf := make([]byte, bmpSize)
		pSaveParam.pImageBuffer = (*C.uchar)(unsafe.Pointer(&bmpBuf[0]))

		//ret = C.MV_CC_SaveImageEx2(handle, pSaveParam)
		//if ret != 0 {
		//	fmt.Println("MV_CC_SaveImageEx2 failure:", fmt.Sprintf("0x%x", C.uint(ret)))
		//	return
		//}
		//
		//ib := C.GoBytes(unsafe.Pointer(pSaveParam.pImageBuffer), C.int(pSaveParam.nDataLen))
		//
		//err := ioutil.WriteFile(filePath, ib, 0666)
		//if err != nil {
		//	fmt.Println("Write file failure，err:", err)
		//	break
		//}
		time.Sleep(5 * time.Second)
		dataBuf = make([]byte, nDataSize)
	}
	<-exitSignal

}

// 无保存取图片
func grabImage() {
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
	ret = C.MV_CC_CreateHandle(&handle, deviceList.pDeviceInfo[0])
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
	nDataSize := C.uint(stParam.nCurValue)
	//var dataBuf C.uchar
	dataBuf := make([]C.uchar, nDataSize)
	pFrameInfo := C.MV_FRAME_OUT_INFO_EX{}
	for {
		ret = C.MV_CC_GetOneFrameTimeout(handle, &dataBuf[0], nDataSize, &pFrameInfo, C.uint(1000))
		if ret == 0 {
			fmt.Printf("get one frame: Width[%v], Height[%v], PixelType[0x%v], nFrameNum[%v]\n", pFrameInfo.nWidth, pFrameInfo.nHeight, pFrameInfo.enPixelType, pFrameInfo.nFrameNum)
		} else {
			fmt.Println("no data")
			continue
		}

		time.Sleep(time.Second)
	}
	<-exitSignal

}
