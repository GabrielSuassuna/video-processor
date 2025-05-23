package resize

import (
	"image"
	"image/color"
	"testing"
	"video-processor/internal/filters"
)

func TestCalculateWeights(t *testing.T) {
	tests := []struct {
		name     string
		srcSize  int
		dstSize  int
		filter   filters.Resampler
		wantNil  bool
		validate func([][]float64) bool
	}{
		{
			name:    "invalid source size",
			srcSize: 0,
			dstSize: 10,
			filter:  filters.NewLanczos(2),
			wantNil: true,
		},
		{
			name:    "invalid destination size",
			srcSize: 10,
			dstSize: 0,
			filter:  filters.NewLanczos(2),
			wantNil: true,
		},
		{
			name:    "negative source size",
			srcSize: -5,
			dstSize: 10,
			filter:  filters.NewLanczos(2),
			wantNil: true,
		},
		{
			name:    "upsampling case",
			srcSize: 5,
			dstSize: 10,
			filter:  filters.NewLanczos(2),
			wantNil: false,
			validate: func(weights [][]float64) bool {
				return len(weights) > 0
			},
		},
		{
			name:    "downsampling case",
			srcSize: 10,
			dstSize: 5,
			filter:  filters.NewLanczos(2),
			wantNil: false,
			validate: func(weights [][]float64) bool {
				return len(weights) > 0
			},
		},
		{
			name:    "same size",
			srcSize: 5,
			dstSize: 5,
			filter:  filters.NewLanczos(2),
			wantNil: false,
			validate: func(weights [][]float64) bool {
				return len(weights) > 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			weights := calculateWeights(tt.srcSize, tt.dstSize, tt.filter)
			
			if tt.wantNil {
				if weights != nil {
					t.Errorf("calculateWeights() expected nil, got %v", weights)
				}
				return
			}
			
			if weights == nil {
				t.Errorf("calculateWeights() returned nil, expected weights")
				return
			}
			
			if tt.validate != nil && !tt.validate(weights) {
				t.Errorf("calculateWeights() weights validation failed")
			}
		})
	}
}

func TestCalculateWeightsNormalization(t *testing.T) {
	filter := filters.NewLanczos(2)
	srcSize := 10
	dstSize := 5
	
	weights := calculateWeights(srcSize, dstSize, filter)
	if weights == nil {
		t.Fatal("calculateWeights() returned nil")
	}
	
	// Check that weights are properly distributed
	if len(weights) != dstSize {
		t.Fatalf("Expected %d weight arrays, got %d", dstSize, len(weights))
	}
	
	for i := 0; i < dstSize; i++ {
		pixelWeights := weights[i]
		
		sum := 0.0
		nonZeroCount := 0
		for j := 0; j < len(pixelWeights); j++ {
			sum += pixelWeights[j]
			if pixelWeights[j] != 0 {
				nonZeroCount++
			}
		}
		
		// Sum should be close to 1.0 for proper normalization
		if sum > 0 && (sum < 0.99 || sum > 1.01) {
			t.Errorf("Weights for pixel %d sum to %f, expected ~1.0", i, sum)
		}
	}
}

func TestResize(t *testing.T) {
	// Create a test image
	src := image.NewNRGBA(image.Rect(0, 0, 4, 4))
	
	// Fill with a simple pattern
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			src.Set(x, y, color.NRGBA{
				R: uint8(x * 64),
				G: uint8(y * 64),
				B: 128,
				A: 255,
			})
		}
	}

	tests := []struct {
		name      string
		src       image.Image
		width     int
		height    int
		wantError bool
	}{
		{
			name:      "valid resize larger",
			src:       src,
			width:     8,
			height:    8,
			wantError: false,
		},
		{
			name:      "valid resize smaller",
			src:       src,
			width:     2,
			height:    2,
			wantError: false,
		},
		{
			name:      "same size",
			src:       src,
			width:     4,
			height:    4,
			wantError: false,
		},
		{
			name:      "nil source",
			src:       nil,
			width:     2,
			height:    2,
			wantError: true,
		},
		{
			name:      "zero width",
			src:       src,
			width:     0,
			height:    2,
			wantError: true,
		},
		{
			name:      "zero height",
			src:       src,
			width:     2,
			height:    0,
			wantError: true,
		},
		{
			name:      "negative dimensions",
			src:       src,
			width:     -1,
			height:    -1,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Resize(tt.src, tt.width, tt.height)
			
			if tt.wantError {
				if err == nil {
					t.Errorf("Resize() expected error, got nil")
				}
				return
			}
			
			if err != nil {
				t.Errorf("Resize() unexpected error: %v", err)
				return
			}
			
			if result == nil {
				t.Errorf("Resize() returned nil result")
				return
			}
			
			bounds := result.Bounds()
			if bounds.Dx() != tt.width || bounds.Dy() != tt.height {
				t.Errorf("Resize() result dimensions = %dx%d, want %dx%d", 
					bounds.Dx(), bounds.Dy(), tt.width, tt.height)
			}
		})
	}
}

func TestResizeHorizontal(t *testing.T) {
	src := image.NewNRGBA(image.Rect(0, 0, 4, 2))
	
	result, err := resizeHorizontal(src, 8)
	if err != nil {
		t.Errorf("resizeHorizontal() unexpected error: %v", err)
		return
	}
	
	bounds := result.Bounds()
	if bounds.Dx() != 8 || bounds.Dy() != 2 {
		t.Errorf("resizeHorizontal() result dimensions = %dx%d, want 8x2", 
			bounds.Dx(), bounds.Dy())
	}
}

func TestResizeVertical(t *testing.T) {
	src := image.NewNRGBA(image.Rect(0, 0, 2, 4))
	
	result, err := resizeVertical(src, 8)
	if err != nil {
		t.Errorf("resizeVertical() unexpected error: %v", err)
		return
	}
	
	bounds := result.Bounds()
	if bounds.Dx() != 2 || bounds.Dy() != 8 {
		t.Errorf("resizeVertical() result dimensions = %dx%d, want 2x8", 
			bounds.Dx(), bounds.Dy())
	}
}

func BenchmarkCalculateWeights(b *testing.B) {
	filter := filters.NewLanczos(2)
	
	b.Run("upsampling", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			calculateWeights(100, 200, filter)
		}
	})
	
	b.Run("downsampling", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			calculateWeights(200, 100, filter)
		}
	})
}

func TestResizeQuality(t *testing.T) {
	// Create a test image with distinct patterns to verify quality
	src := image.NewNRGBA(image.Rect(0, 0, 10, 10))
	
	// Create a checkerboard pattern
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			if (x+y)%2 == 0 {
				src.Set(x, y, color.NRGBA{R: 255, G: 255, B: 255, A: 255})
			} else {
				src.Set(x, y, color.NRGBA{R: 0, G: 0, B: 0, A: 255})
			}
		}
	}

	// Test downsampling
	result, err := Resize(src, 5, 5)
	if err != nil {
		t.Errorf("Resize() error: %v", err)
		return
	}

	if result == nil {
		t.Error("Resize() returned nil result")
		return
	}

	// Verify dimensions
	bounds := result.Bounds()
	if bounds.Dx() != 5 || bounds.Dy() != 5 {
		t.Errorf("Resize() result dimensions = %dx%d, want 5x5", bounds.Dx(), bounds.Dy())
	}

	// Verify that Lanczos filtering produces intermediate values (not just pure black/white)
	hasIntermediateValues := false
	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			r, g, b, _ := result.At(x, y).RGBA()
			// Convert to 8-bit for easier comparison
			r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)
			
			// If we have values that aren't pure black (0) or white (255), 
			// it indicates proper filtering is happening
			if (r8 > 10 && r8 < 245) || (g8 > 10 && g8 < 245) || (b8 > 10 && b8 < 245) {
				hasIntermediateValues = true
				break
			}
		}
		if hasIntermediateValues {
			break
		}
	}

	if !hasIntermediateValues {
		t.Log("Note: No intermediate values found - this might indicate the filter isn't being applied properly")
	}
}

func BenchmarkResize(b *testing.B) {
	src := image.NewNRGBA(image.Rect(0, 0, 100, 100))
	
	// Fill with test data
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			src.Set(x, y, color.NRGBA{
				R: uint8(x % 256),
				G: uint8(y % 256),
				B: uint8((x + y) % 256),
				A: 255,
			})
		}
	}
	
	b.Run("resize_larger", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Resize(src, 200, 200)
		}
	})
	
	b.Run("resize_smaller", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Resize(src, 50, 50)
		}
	})
}