package sdk_camera

import "testing"

func TestGrabImage(t *testing.T) {
	t.Run("#取图并保存", func(t *testing.T) {
		grabImageWithSave()
	})

	t.Run("#回调函数取图", func(t *testing.T) {
		grabImageWithCallback()
	})
}
