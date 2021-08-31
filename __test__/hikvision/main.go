package main

import (
	sdk_camera "git.querycap.com/aia/env-sdk-camera"
	"git.querycap.com/aia/env-sdk-camera/__test__/sdk"
)

func main() {
	sdk.SDKExampleWithCallback(sdk_camera.TypeofHikvisionCameraSDK)
}
