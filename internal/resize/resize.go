package resize

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"video-processor/internal/filters"
)

func Resize(src image.Image, width, height int) (*image.NRGBA, error) {
	dstWidth := width
	dstHeight := height

	if src == nil {
		return nil, errors.New("source image is nil")
	}
	if dstWidth <= 0 || dstHeight <= 0 {
		return nil, fmt.Errorf("invalid dimensions: width=%d, height=%d", width, height)
	}

	srcWidth := src.Bounds().Dx()
	srcHeight := src.Bounds().Dy()

	if srcWidth != dstWidth && srcHeight != dstHeight {
		image, err := resizeHorizontal(src, dstWidth)
		if err != nil {
			return nil, err
		}
		return resizeVertical(image, dstHeight)
	}

	if srcWidth != dstWidth {
		return resizeHorizontal(src, dstWidth)
	}

	return resizeVertical(src, dstHeight)
}

func resizeVertical(src image.Image, height int) (*image.NRGBA, error) {
	srcBounds := src.Bounds()
	srcWidth := srcBounds.Dx()
	srcHeight := srcBounds.Dy()

	if srcHeight == height {
		// No vertical resize needed, just copy
		dst := image.NewNRGBA(image.Rect(0, 0, srcWidth, height))
		draw.Draw(dst, dst.Bounds(), src, srcBounds.Min, draw.Src)
		return dst, nil
	}

	dst := image.NewNRGBA(image.Rect(0, 0, srcWidth, height))
	filter := filters.NewLanczos(3)
	weights := calculateWeights(srcHeight, height, filter)

	if weights == nil {
		return nil, fmt.Errorf("failed to calculate weights for vertical resize")
	}

	// Process each column
	for x := 0; x < srcWidth; x++ {
		for dstY := 0; dstY < height; dstY++ {
			var r, g, b, a float64
			pixelWeights := weights[dstY]

			scale := float64(srcHeight) / float64(height)
			center := (float64(dstY)+0.5)*scale - 0.5
			support := float64(filter.Radius)
			if scale > 1.0 {
				support *= scale
			}

			left := int(center - support)
			right := int(center + support)

			if left < 0 {
				left = 0
			}
			if right >= srcHeight {
				right = srcHeight - 1
			}

			weightIdx := 0
			for srcY := left; srcY <= right && weightIdx < len(pixelWeights); srcY++ {
				weight := pixelWeights[weightIdx]
				if weight != 0 {
					srcColor := src.At(x+srcBounds.Min.X, srcY+srcBounds.Min.Y)
					srcR, srcG, srcB, srcA := srcColor.RGBA()

					r += float64(srcR) * weight
					g += float64(srcG) * weight
					b += float64(srcB) * weight
					a += float64(srcA) * weight
				}
				weightIdx++
			}

			// Clamp values and convert back
			if r < 0 {
				r = 0
			} else if r > 65535 {
				r = 65535
			}
			if g < 0 {
				g = 0
			} else if g > 65535 {
				g = 65535
			}
			if b < 0 {
				b = 0
			} else if b > 65535 {
				b = 65535
			}
			if a < 0 {
				a = 0
			} else if a > 65535 {
				a = 65535
			}

			dst.Set(x, dstY, color.NRGBA64{
				R: uint16(r),
				G: uint16(g),
				B: uint16(b),
				A: uint16(a),
			})
		}
	}

	return dst, nil
}

func resizeHorizontal(src image.Image, width int) (*image.NRGBA, error) {
	srcBounds := src.Bounds()
	srcWidth := srcBounds.Dx()
	srcHeight := srcBounds.Dy()

	if srcWidth == width {
		// No horizontal resize needed, just copy
		dst := image.NewNRGBA(image.Rect(0, 0, width, srcHeight))
		draw.Draw(dst, dst.Bounds(), src, srcBounds.Min, draw.Src)
		return dst, nil
	}

	dst := image.NewNRGBA(image.Rect(0, 0, width, srcHeight))
	filter := filters.NewLanczos(3)
	weights := calculateWeights(srcWidth, width, filter)

	if weights == nil {
		return nil, fmt.Errorf("failed to calculate weights for horizontal resize")
	}

	// Process each row
	for y := 0; y < srcHeight; y++ {
		for dstX := 0; dstX < width; dstX++ {
			var r, g, b, a float64
			pixelWeights := weights[dstX]

			scale := float64(srcWidth) / float64(width)
			center := (float64(dstX)+0.5)*scale - 0.5
			support := float64(filter.Radius)
			if scale > 1.0 {
				support *= scale
			}

			left := int(center - support)
			right := int(center + support)

			if left < 0 {
				left = 0
			}
			if right >= srcWidth {
				right = srcWidth - 1
			}

			weightIdx := 0
			for srcX := left; srcX <= right && weightIdx < len(pixelWeights); srcX++ {
				weight := pixelWeights[weightIdx]
				if weight != 0 {
					srcColor := src.At(srcX+srcBounds.Min.X, y+srcBounds.Min.Y)
					srcR, srcG, srcB, srcA := srcColor.RGBA()

					r += float64(srcR) * weight
					g += float64(srcG) * weight
					b += float64(srcB) * weight
					a += float64(srcA) * weight
				}
				weightIdx++
			}

			// Clamp values and convert back
			if r < 0 {
				r = 0
			} else if r > 65535 {
				r = 65535
			}
			if g < 0 {
				g = 0
			} else if g > 65535 {
				g = 65535
			}
			if b < 0 {
				b = 0
			} else if b > 65535 {
				b = 65535
			}
			if a < 0 {
				a = 0
			} else if a > 65535 {
				a = 65535
			}

			dst.Set(dstX, y, color.NRGBA64{
				R: uint16(r),
				G: uint16(g),
				B: uint16(b),
				A: uint16(a),
			})
		}
	}

	return dst, nil
}

func calculateWeights(srcSize, dstSize int, filter filters.Resampler) [][]float64 {
	if srcSize <= 0 || dstSize <= 0 {
		return nil
	}

	// Calculate the scaling factor
	scale := float64(srcSize) / float64(dstSize)

	// For downsampling, we need to expand the filter support
	filterRadius := 1.0
	if lanczos, ok := filter.(*filters.Lanczos); ok {
		filterRadius = float64(lanczos.Radius)
	}

	// Support radius should be at least as large as the scaling factor for downsampling
	support := filterRadius
	if scale > 1.0 {
		support *= scale
	}

	// Total number of weights needed per pixel
	weightsPerPixel := int(2*support) + 1
	weights := make([][]float64, dstSize)

	for dstIdx := 0; dstIdx < dstSize; dstIdx++ {
		// Calculate the center position in source coordinates
		center := (float64(dstIdx)+0.5)*scale - 0.5

		// Calculate the range of source pixels that contribute to this destination pixel
		left := int(center - support)
		right := int(center + support)

		// Ensure we stay within bounds
		if left < 0 {
			left = 0
		}
		if right >= srcSize {
			right = srcSize - 1
		}

		// Initialize weights slice for this destination pixel
		weights[dstIdx] = make([]float64, weightsPerPixel)

		// Calculate weights for this destination pixel
		weightSum := 0.0
		weightIdx := 0

		for srcIdx := left; srcIdx <= right; srcIdx++ {
			distance := float64(srcIdx) - center

			// Calculate weight using the filter
			var weight float64
			if scale > 1.0 {
				// Downsampling: scale the filter
				weight = filter.Kernel(distance / scale)
			} else {
				// Upsampling: use filter as-is
				weight = filter.Kernel(distance)
			}

			if weight != 0 {
				weights[dstIdx][weightIdx] = weight
				weightSum += weight
			}
			weightIdx++
		}

		// Normalize weights so they sum to 1
		if weightSum > 0 {
			for i := 0; i < weightsPerPixel; i++ {
				weights[dstIdx][i] /= weightSum
			}
		}
	}

	return weights
}
