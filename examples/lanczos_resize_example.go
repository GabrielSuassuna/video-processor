package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"
	"video-processor/internal/resize"
)

func main() {
	// Create a high-resolution test image with fine details
	src := createTestImage(400, 400)
	
	// Save original
	saveImage(src, "original_400x400.jpg")
	fmt.Println("Created original_400x400.jpg")
	
	// Resize to smaller size using Lanczos
	resized, err := resize.Resize(src, 100, 100)
	if err != nil {
		log.Fatalf("Error resizing image: %v", err)
	}
	
	saveImage(resized, "lanczos_resized_100x100.jpg")
	fmt.Println("Created lanczos_resized_100x100.jpg using Lanczos filter")
	
	// Resize back up to demonstrate upsampling quality
	upsampled, err := resize.Resize(resized, 200, 200)
	if err != nil {
		log.Fatalf("Error upsampling image: %v", err)
	}
	
	saveImage(upsampled, "lanczos_upsampled_200x200.jpg")
	fmt.Println("Created lanczos_upsampled_200x200.jpg")
	
	fmt.Println("\nLanczos resize example completed!")
	fmt.Println("The Lanczos filter provides high-quality resizing with:")
	fmt.Println("- Reduced aliasing artifacts")
	fmt.Println("- Better preservation of fine details")
	fmt.Println("- Smoother gradients")
	fmt.Println("- Minimal ringing artifacts")
}

func createTestImage(width, height int) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Create a complex pattern with:
			// 1. Diagonal stripes
			// 2. Radial gradient
			// 3. High frequency details
			
			// Diagonal stripes
			stripe := (x + y) % 20 < 10
			
			// Radial gradient from center
			centerX, centerY := float64(width)/2, float64(height)/2
			dx, dy := float64(x)-centerX, float64(y)-centerY
			distance := (dx*dx + dy*dy) / (centerX * centerY)
			
			// High frequency checkerboard in corner
			highFreq := (x/2+y/2)%2 == 0 && x < width/4 && y < height/4
			
			var r, g, b uint8
			
			if highFreq {
				// High frequency area - alternating black/white
				if (x+y)%2 == 0 {
					r, g, b = 255, 255, 255
				} else {
					r, g, b = 0, 0, 0
				}
			} else if stripe {
				// Diagonal stripes with gradient
				intensity := uint8(255 * (1 - distance))
				r, g, b = intensity, intensity/2, intensity/4
			} else {
				// Background gradient
				intensity := uint8(128 * (1 - distance))
				r, g, b = intensity/4, intensity/2, intensity
			}
			
			img.Set(x, y, color.NRGBA{R: r, G: g, B: b, A: 255})
		}
	}
	
	return img
}

func saveImage(img image.Image, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Error creating file %s: %v", filename, err)
	}
	defer file.Close()
	
	err = jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
	if err != nil {
		log.Fatalf("Error encoding image %s: %v", filename, err)
	}
}