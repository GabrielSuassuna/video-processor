package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"video-processor/internal/resize"
)

func main() {
	// Define command-line flags
	inputFile := flag.String("input", "", "Path to input image file (required)")
	outputFile := flag.String("output", "", "Path to output image file (default: input file with _resized suffix)")
	width := flag.Int("width", 0, "Target width in pixels (required)")
	height := flag.Int("height", 0, "Target height in pixels (required)")
	verbose := flag.Bool("verbose", false, "Enable verbose output")

	// Parse command-line flags
	flag.Parse()

	// Validate input file
	if *inputFile == "" {
		fmt.Println("Error: Input file is required")
		flag.Usage()
		os.Exit(1)
	}

	// Check if input file exists
	if _, err := os.Stat(*inputFile); os.IsNotExist(err) {
		fmt.Printf("Error: Input file does not exist: %s\n", *inputFile)
		os.Exit(1)
	}

	// Validate dimensions
	if *width <= 0 || *height <= 0 {
		fmt.Println("Error: Both width and height must be greater than 0")
		flag.Usage()
		os.Exit(1)
	}

	// Generate default output file name if not specified
	if *outputFile == "" {
		ext := filepath.Ext(*inputFile)
		baseName := strings.TrimSuffix(*inputFile, ext)
		*outputFile = fmt.Sprintf("%s_resized%s", baseName, ext)
	}

	if *verbose {
		fmt.Println("Starting image resizing...")
		fmt.Printf("Input: %s\n", *inputFile)
		fmt.Printf("Output: %s\n", *outputFile)
		fmt.Printf("Dimensions: %d x %d\n", *width, *height)
	}

	// Load the input image
	inputImg, format, err := loadImage(*inputFile)
	if err != nil {
		fmt.Printf("Error loading image: %v\n", err)
		os.Exit(1)
	}

	// Resize the image
	resizedImg, err := resize.Resize(inputImg, *width, *height)
	if err != nil {
		fmt.Printf("Error resizing image: %v\n", err)
		os.Exit(1)
	}

	// Save the resized image
	err = saveImage(*outputFile, resizedImg, format)
	if err != nil {
		fmt.Printf("Error saving image: %v\n", err)
		os.Exit(1)
	}

	if *verbose {
		fmt.Println("Image resizing completed successfully")
	}
}

// loadImage loads an image from the given file path
func loadImage(filePath string) (image.Image, string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode image: %w", err)
	}

	return img, format, nil
}

// saveImage saves an image to the given file path
func saveImage(filePath string, img *image.NRGBA, format string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(filePath))
	
	switch {
	case format == "jpeg" || ext == ".jpg" || ext == ".jpeg":
		err = jpeg.Encode(file, img, &jpeg.Options{Quality: 95})
	case format == "png" || ext == ".png":
		err = png.Encode(file, img)
	default:
		// Default to JPEG if format is unknown
		err = jpeg.Encode(file, img, &jpeg.Options{Quality: 95})
	}

	if err != nil {
		return fmt.Errorf("failed to encode image: %w", err)
	}

	return nil
}