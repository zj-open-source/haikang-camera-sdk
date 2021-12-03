package constants

//go:generate tools gen enum PixelFormat
type PixelFormat uint8

// PixelFormat 图片像素格式
const (
	PIXEL_FORMAT_UNKNOWN PixelFormat = iota // 未知格式
	PIXEL_FORMAT__MONO8                     // mono8
	PIXEL_FORMAT__BAYER8                    // bayer8
	PIXEL_FORMAT__YUV                       // yuv
	PIXEL_FORMAT__RGB                       // rgb
	PIXEL_FORMAT__BGR                       // bgr
)
