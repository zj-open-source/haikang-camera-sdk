package sdk_camera

/*
#include "include/CameraParams.h"
//#include <>
*/
import "C"
import (
	"fmt"
	gopointer "github.com/mattn/go-pointer"
	"unsafe"
)

//export Hello
func Hello() {
	fmt.Println("hello,cgo from GO")
}

//export Callback
func Callback(data *C.uint8_t, pFrameInfo *C.MV_FRAME_OUT_INFO_EX, user unsafe.Pointer) {
	dataBytes := C.GoBytes(unsafe.Pointer(data), C.int(pFrameInfo.nFrameLen)) // 读取C返回的图片buf

	// 从map中将unsafe.pointer对应的内容取出来作处理,因为C中不能长期持有GO的指针对象,只能通过全局变量来映射
	userData := gopointer.Restore(user).(*CallbackStruct)
	// 执行golang的回调
	userData.Callback(FrameOutInfo{
		Width:      uint16(pFrameInfo.nWidth),
		Height:     uint16(pFrameInfo.nHeight),
		FrameNum:   uint32(pFrameInfo.nFrameNum),
		FamilyName: userData.FamilyName,
		PortID:     userData.PortID,
		DataBuf:    dataBytes,
	})
}
