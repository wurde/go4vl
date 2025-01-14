package v4l2

// #include <linux/videodev2.h>
import "C"

import (
	"fmt"
	"unsafe"
)

// Rect (v4l2_rect)
// https://www.kernel.org/doc/html/v4.14/media/uapi/v4l/dev-overlay.html?highlight=v4l2_rect#c.v4l2_rect
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L412
type Rect struct {
	Left   int32
	Top    int32
	Width  uint32
	Height uint32
}

// Fract (v4l2_fract)
// https://www.kernel.org/doc/html/v4.14/media/uapi/v4l/vidioc-enumstd.html#c.v4l2_fract
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L419
type Fract struct {
	Numerator   uint32
	Denominator uint32
}

// CropCapability (v4l2_cropcap)
// https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-cropcap.html#c.v4l2_cropcap
// https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/videodev2.h#L1221
type CropCapability struct {
	StreamType  uint32
	Bounds      Rect
	DefaultRect Rect
	PixelAspect Fract
	_           [4]uint32
}

// GetCropCapability  retrieves cropping info for specified device
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-cropcap.html#ioctl-vidioc-cropcap
func GetCropCapability(fd uintptr) (CropCapability, error) {
	var cap C.struct_v4l2_cropcap
	cap._type = C.uint(BufTypeVideoCapture)

	if err := send(fd, C.VIDIOC_CROPCAP, uintptr(unsafe.Pointer(&cap))); err != nil {
		return CropCapability{}, fmt.Errorf("crop capability: %w", err)
	}

	return *(*CropCapability)(unsafe.Pointer(&cap)), nil
}

// SetCropRect sets the cropping dimension for specified device
// See https://www.kernel.org/doc/html/latest/userspace-api/media/v4l/vidioc-g-crop.html#ioctl-vidioc-g-crop-vidioc-s-crop
func SetCropRect(fd uintptr, r Rect) error {
	var crop C.struct_v4l2_crop
	crop._type = C.uint(BufTypeVideoCapture)
	crop.c = *(*C.struct_v4l2_rect)(unsafe.Pointer(&r))

	if err := send(fd, C.VIDIOC_S_CROP, uintptr(unsafe.Pointer(&crop))); err != nil {
		return fmt.Errorf("set crop: %w", err)
	}
	return nil
}

func (c CropCapability) String() string {
	return fmt.Sprintf("default:{top=%d, left=%d, width=%d,height=%d};  bounds:{top=%d, left=%d, width=%d,height=%d}; pixel-aspect{%d:%d}",
		c.DefaultRect.Top,
		c.DefaultRect.Left,
		c.DefaultRect.Width,
		c.DefaultRect.Height,

		c.Bounds.Top,
		c.Bounds.Left,
		c.Bounds.Width,
		c.Bounds.Height,

		c.PixelAspect.Numerator,
		c.PixelAspect.Denominator,
	)
}
