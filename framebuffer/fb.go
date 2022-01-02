package framebuffer

import (
	"golang.org/x/sys/unix"
	"image"
	"image/color"
	"image/draw"
	"os"
	"syscall"
	"unsafe"
)

const (
	FBIOGET_VSCREENINFO = 0x4600
	FBIOGET_FSCREENINFO = 0x4602
	FBIOBLANK           = 0x4611
)

var _ draw.Image = (*SimpleRGBA)(nil)

// Open expects a framebuffer device as its argument (such as "/dev/fb0"). The
// device will be memory-mapped to a local buffer. Writing to the device changes
// the screen output.
// The returned Device implements the draw.Image interface. This means that you
// can use it to copy to and from other images.
// After you are done using the Device, call Close on it to unmap the memory and
// close the framebuffer file.
func Open(device string) (*Device, error) {
	file, err := os.OpenFile(device, os.O_RDWR, os.ModeDevice)
	if err != nil {
		return nil, err
	}

	_ = unix.IoctlSetInt(int(file.Fd()), FBIOBLANK, 0)
	_ = unix.IoctlSetInt(int(file.Fd()), FBIOBLANK, 4)
	_ = unix.IoctlSetInt(int(file.Fd()), FBIOBLANK, 0)

	fixInfo, _ := getFixScreenInfo(file.Fd())
	varInfo, _ := getVarScreenInfo(file.Fd())

	pixels, err := syscall.Mmap(
		int(file.Fd()),
		0,
		int(fixInfo.smemLen),
		//int(varInfo.Xres*varInfo.Yres*varInfo.bitsPerPixel/8),
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_SHARED,
	)
	if err != nil {
		file.Close()
		return nil, err
	}

	return &Device{
		file: file,
		SimpleRGBA: &SimpleRGBA{
			Pixels: pixels,
			Stride: int(fixInfo.lineLength),
			Xres:   int(varInfo.xres),
			Yres:   int(varInfo.yres),
		},
	}, nil
}

// Device represents the frame buffer. It implements the draw.Image interface.
type Device struct {
	*SimpleRGBA
	file *os.File
}

type SimpleRGBA struct {
	Pixels []byte
	Stride int
	Xres   int
	Yres   int
}

func (s *SimpleRGBA) ColorModel() color.Model {
	return color.RGBAModel
}

func (s *SimpleRGBA) Bounds() image.Rectangle {
	return image.Rect(0, 0, s.Xres, s.Yres)
}

func (s *SimpleRGBA) At(x, y int) color.Color {
	if x < 0 || x > s.Xres || y < 0 || y > s.Yres {
		return color.RGBA{}
	}
	i := y*s.Stride + x*4
	n := s.Pixels[i : i+4 : i+4] // Small cap improves performance, see https://golang.org/issue/27857
	return color.RGBA{R: n[0], G: n[1], B: n[2], A: n[3]}
}

func (s *SimpleRGBA) Black(rect image.Rectangle) {
	start := rect.Min.Y*s.Stride + rect.Min.X*4
	end := rect.Max.Y*s.Stride + rect.Max.X*4 + 4
	for i := start; i < end; i++ {
		s.Pixels[i] = 0
	}
}

func (s *SimpleRGBA) Set(x, y int, c color.Color) {
	r, g, b, a := c.RGBA()
	i := y*s.Stride + x*4
	n := s.Pixels[i : i+4 : i+4] // Small cap improves performance, see https://golang.org/issue/27857
	n[0] = byte(r)
	n[1] = byte(g)
	n[2] = byte(b)
	n[3] = byte(a)
}

// Close unmaps the framebuffer memory and closes the device file. Call this
// function when you are done using the frame buffer.
func (d *Device) Close() {
	syscall.Munmap(d.Pixels)
	d.file.Close()
}

func ioctlPtr(fd uintptr, req uint, arg unsafe.Pointer) error {
	_, _, err := unix.Syscall(unix.SYS_IOCTL, fd, uintptr(req), uintptr(arg))
	if err != 0 {
		return err
	}

	return nil
}

func getFixScreenInfo(fd uintptr) (*fixScreenInfo, error) {
	var value fixScreenInfo
	err := ioctlPtr(fd, FBIOGET_FSCREENINFO, unsafe.Pointer(&value))
	return &value, err
}

func getVarScreenInfo(fd uintptr) (*varScreenInfo, error) {
	var value varScreenInfo
	err := ioctlPtr(fd, FBIOGET_VSCREENINFO, unsafe.Pointer(&value))
	return &value, err
}

type fixScreenInfo struct {
	id           [16]byte
	smemStart    uint32
	smemLen      uint32
	fbType       uint32
	typeAux      uint32
	visual       uint32
	xPanStep     uint16
	yPanStep     uint16
	yWrapStep    uint16
	lineLength   uint32
	mmioStart    uint32
	mmioLen      uint32
	accel        uint32
	capabilities uint16
	reserved     [2]uint16
}

type bitField struct {
	offset   uint32
	length   uint32
	msbRight uint32
}

type varScreenInfo struct {
	xres        uint32
	yres        uint32
	xresVirtual uint32
	yresVirtual uint32
	xoffset     uint32
	yoffset     uint32

	bitsPerPixel uint32
	grayscale    uint32

	red    bitField
	green  bitField
	blue   bitField
	transp bitField
	nonstd uint32

	activate uint32

	height uint32
	width  uint32

	accelFlags uint32

	pixclock    uint32
	leftMargin  uint32
	rightMargin uint32
	upperMargin uint32
	lowerMargin uint32
	hsyncLen    uint32
	vsyncLen    uint32
	sync        uint32
	vmode       uint32
	rotate      uint32
	colorspace  uint32
	reserved    [4]uint32
}
