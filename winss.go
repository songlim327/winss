package winss

import (
	"errors"
	"image"
	"image/color"
	"image/png"
	"os"
	"unsafe"

	"github.com/gonutz/w32/v2"
)

// FindWindow returns the window handler
func FindWindow(name string) (w32.HWND, error) {
	hWnd := w32.FindWindow("", name)
	if hWnd == 0 {
		return w32.HWND(unsafe.Pointer(nil)), errors.New("windows unable to find the handler for the name specified")
	}
	return hWnd, nil
}

// Screenshot will take a screenshot of the given windows handler
func Screenshot(hWnd w32.HWND) *image.RGBA {
	// Retrieve the handle to a display device context for the client window
	hdcWindow := w32.GetWindowDC(hWnd)
	defer w32.ReleaseDC(hWnd, hdcWindow)

	hdcMemDC := w32.CreateCompatibleDC(hdcWindow)
	defer w32.DeleteDC(hdcMemDC)

	// Get client area for size calculations
	rcClient := w32.GetWindowRect(hWnd)
	width := int(rcClient.Right) - int(rcClient.Left)
	height := int(rcClient.Bottom) - int(rcClient.Top)

	// Create a compatible bitmap from the Window DC
	hbmWindow := w32.CreateCompatibleBitmap(hdcWindow, width, height)
	defer w32.DeleteObject(w32.HGDIOBJ(hbmWindow))

	// Select the compatible bitmap into the compatible memory DC
	hOld := w32.SelectObject(hdcMemDC, w32.HGDIOBJ(hbmWindow))

	// Bit block transfer into our compatible memory DC
	w32.BitBlt(hdcMemDC, 0, 0, width, height, hdcWindow, 0, 0, w32.SRCCOPY)
	w32.SelectObject(hdcMemDC, hOld)

	bmInfo := w32.BITMAPINFO{
		BmiHeader: w32.BITMAPINFOHEADER{
			BiSize:  uint32(unsafe.Sizeof(w32.BITMAPINFOHEADER{})),
			BiWidth: int32(width),
			// Pixel data is top down format
			BiHeight:      -int32(height),
			BiPlanes:      1,
			BiBitCount:    32,
			BiCompression: w32.BI_RGB,
			BiSizeImage:   uint32(width * height * 4),
		},
	}

	// Initialize array of bytes with the size of width x height x 4 (each pixel consists of 4 bytes, representing RGBA)
	pixels := make([]byte, width*height*4)
	w32.GetDIBits(hdcWindow, hbmWindow, 0, uint(height), unsafe.Pointer(&pixels[0]), &bmInfo, w32.DIB_RGB_COLORS)

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			offset := (y*int(width) + x) * 4
			img.SetRGBA(x, y, color.RGBA{R: pixels[offset+2], G: pixels[offset+1], B: pixels[offset], A: 255})
		}
	}
	return img
}

// SaveImage saves the image to the specified file name
func SaveImage(img image.Image, filename string) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		closeErr := file.Close()
		if closeErr != nil {
			err = closeErr
		}
	}(f)

	if err != nil {
		return err
	}

	if err = png.Encode(f, img); err != nil {
		return err
	}
	return nil
}
