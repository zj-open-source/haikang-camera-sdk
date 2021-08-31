package sdk

import (
	"bytes"
	"fmt"
	sdk_camera "git.querycap.com/aia/env-sdk-camera"
	"git.querycap.com/aia/env-sdk-camera/constants"
	gopointer "github.com/mattn/go-pointer"
	"io/ioutil"
	"sync"
	"time"
	//"io/ioutil"
)

func SDKExampleWithCallback(sdkType sdk_camera.CameraSDKType) {
	sdk, err := sdk_camera.GetCameraSDK(sdkType)
	if err != nil {
		fmt.Println("get camera sdk failure:", err)
		return
	}

	sdkVersion := sdk.GetSDKVersion()
	fmt.Println("SDKVersion:", sdkVersion)

	deviceList, err := sdk.GetDeviceList()
	if err != nil {
		fmt.Println("get device list failure:", err)
		return
	}
	if deviceList.(*sdk_camera.MvCcDeviceInfoList).DeviceNum == 0 {
		fmt.Println("no device was found")
		return
	}

	wg := &sync.WaitGroup{}
	for i, val := range *deviceList.(*sdk_camera.MvCcDeviceInfoList).MvCcDeviceInfo {
		wg.Add(1)
		go func(idx int, wgSync *sync.WaitGroup, cInfo sdk_camera.MvCcDeviceInfo) {
			sdkNew, err := sdk_camera.GetCameraSDK(sdkType)
			if err != nil {
				fmt.Println("get new camera sdk failure:", err)
				return
			}

			sdkNew.(*sdk_camera.HikvisionCameraSDK).SetDeviceList(sdk.(*sdk_camera.HikvisionCameraSDK))
			err = sdkNew.CreateHandle(idx)
			if err != nil {
				fmt.Println("create handle failure:", err)
				return
			}
			defer func() {
				_ = sdkNew.DestroyHandle()
			}()

			err = sdkNew.OpenDevice()
			if err != nil {
				fmt.Println("open device failure:", err)
				return
			}
			defer func() {
				_ = sdkNew.CloseDevice()
			}()

			cb := &sdk_camera.CallbackStruct{
				FamilyName: cInfo.FamilyName,
				PortID:     cInfo.PortID,
				Callback: func(frameInfo sdk_camera.FrameOutInfo) {
					// 回调函数写这里
					fmt.Printf("get one frame:[PortID:%v], [DataBuf len:%d], Width[%v], Height[%v], nFrameNum[%v]\n", frameInfo.PortID, len(frameInfo.DataBuf), frameInfo.Width, frameInfo.Height, frameInfo.FrameNum)

					buf := &bytes.Buffer{}
					err := sdk_camera.BGRToJpeg(buf, frameInfo.DataBuf, int(frameInfo.Width), int(frameInfo.Height), nil)
					if err != nil {
						fmt.Println("BGRTOJpeg failure:", err)
					}
					err = ioutil.WriteFile("/home/rock/test.jpeg", buf.Bytes(), 0666)
					if err != nil {
						fmt.Println("WriteFile failure:", err)
					}
				},
			}
			p := gopointer.Save(cb)
			defer gopointer.Unref(p)
			if err := sdkNew.GetOneFrameFroBGRWithCallback(p); err != nil {
				//if err := sdkNew.GetOneFrameFroRGBWithCallback(p); err != nil {
				//if err := sdkNew.GetOneFrameWithCallback(p); err != nil {
				fmt.Println("GetOneFrameWithCallback === ", err)
				return
			}

			err = sdkNew.StartGrabbing(constants.TRIGGER_MODE__ON)
			defer func() {
				_ = sdkNew.StopGrabbing()
			}()

			j := 0
			for {
				time.Sleep(time.Second)
				fmt.Println(j)
				j++
			}
		}(i, wg, val)
	}
	wg.Wait()
	fmt.Println("10s later process will exit.")
	time.Sleep(time.Second * 10)
	fmt.Println("bye bye!")
}
