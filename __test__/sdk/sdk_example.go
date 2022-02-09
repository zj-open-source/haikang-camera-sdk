package sdk

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"sync"
	"time"
	
	sdk_camera "github.com/zjzjzjzj1874/haikang-camera-sdk"
	"github.com/zjzjzjzj1874/haikang-camera-sdk/constants"
)

func SDKExample(sdkType sdk_camera.CameraSDKType) {
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
	for i := range *deviceList.(*sdk_camera.MvCcDeviceInfoList).MvCcDeviceInfo {
		wg.Add(1)
		go func(idx int, wgSync *sync.WaitGroup) {
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

			err = sdkNew.StartGrabbing(constants.TRIGGER_MODE__OFF)
			defer func() {
				_ = sdkNew.StopGrabbing()
			}()

			for {
				// dataBuf:测试用例中，不需要处理，所以无需接受
				frameInfo, _, err := sdkNew.GetImageForRGB(1000)
				//frameInfo, _, err := sdkNew.GetOneFrameTimeout(1000)
				//frameInfo, _, err := (*cameraSDK).GetOneFrameTimeout(1000)
				//frameInfo, dataBuf, err := (*cameraSDK).GetOneFrameTimeout(1000)
				if err != nil {
					fmt.Println("GetOneFrameTimeout failure:", err)
					continue
				}
				if frameInfo != nil {
					fmt.Printf("idx:%d,图片高度:%d;当前帧数：%d \n", idx, frameInfo.(*sdk_camera.MvFrameOutInfoEx).Height, frameInfo.(*sdk_camera.MvFrameOutInfoEx).FrameNum)
					if frameInfo.(*sdk_camera.MvFrameOutInfoEx).FrameNum == 10000 {
						wgSync.Done()
						break
					}
				}

				//mat, err := gocv.NewMatFromBytes(int(frameInfo.Height), int(frameInfo.Width), gocv.MatTypeCV8UC3, dataBuf)
				//if err != nil {
				//	fmt.Printf("new mat failure:[err:%v]\n", err)
				//	continue
				//}
				//img, err := mat.ToImage()
				//if err != nil {
				//	fmt.Printf("mat to image failure:[err:%v]\n", err)
				//	continue
				//}
				//
				//imageByte, err := ImageToJpeg(img)
				//if err != nil {
				//	continue
				//}
				//
				//err = ioutil.WriteFile("/tmp/save.jpeg", imageByte, 0666)
				//if err != nil {
				//	fmt.Println("write file error:", err)
				//}
			}
			time.Sleep(time.Second)
		}(i, wg)
	}
	wg.Wait()
	fmt.Println("10s later process will exit.")
	time.Sleep(time.Second * 10)
	fmt.Println("bye bye!")
}

func ImageToJpeg(img image.Image) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := jpeg.Encode(buf, img, &jpeg.Options{Quality: 75})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
