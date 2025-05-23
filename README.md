# Video Processor - Image Resizer

A high-performance command-line tool written in Go for resizing images with professional-quality filters.

## Features

- **High-Quality Resizing**: Uses Lanczos-2 filter for superior image quality
- **Anti-Aliasing**: Reduces jagged edges and artifacts during scaling
- **Detail Preservation**: Maintains fine details during downsampling
- **Smooth Gradients**: Creates smoother color transitions
- **Format Support**: JPEG, PNG, and other common image formats
- **Efficient Processing**: Optimized algorithms for fast resizing
- **Comprehensive Testing**: Full test suite with benchmarks

## Requirements

- Go 1.16 or higher

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/video-processor.git
   cd video-processor
   ```

2. Build the project:
   ```
   go build -o resizer ./cmd/main.go
   ```

## Usage

Basic usage:

```
./resizer -input input.jpg -output resized.jpg -width 800 -height 600
```

### Command Line Options

| Flag | Description |
|------|-------------|
| `-input` | Path to input image file (required) |
| `-output` | Path to output image file (default: input file with _resized suffix) |
| `-width` | Target width in pixels (required) |
| `-height` | Target height in pixels (required) |
| `-verbose` | Enable verbose output |

### Examples

Resize an image to 800x600:
```
./resizer -input image.jpg -width 800 -height 600
```

Resize with a specific output filename:
```
./resizer -input image.png -output thumbnail.png -width 300 -height 200
```

Enable verbose output to see processing details:
```
./resizer -input image.jpg -width 1024 -height 768 -verbose
```

## How It Works

The Image Resizer uses advanced Lanczos interpolation for high-quality resizing:

1. **Weight Calculation**: Computes optimal filter weights for each destination pixel
2. **Lanczos Filtering**: Applies Lanczos-2 filter for superior quality
3. **Color Interpolation**: Uses full-precision (16-bit) color processing
4. **Normalization**: Ensures weights sum to 1.0 for accurate color reproduction
5. **Two-Pass Resizing**: Handles complex resizes with horizontal then vertical passes

### Algorithm Features

- **Lanczos-2 Filter**: Provides excellent balance between sharpness and artifacts
- **Adaptive Support**: Automatically adjusts filter support for downsampling
- **Bounds Checking**: Prevents memory access errors during processing
- **Weight Normalization**: Maintains color accuracy across all scaling ratios

### Quality Benefits

- **Reduced Aliasing**: Minimizes stair-stepping and jagged edges
- **Better Detail Preservation**: Maintains fine image details during scaling
- **Smoother Gradients**: Creates natural color transitions
- **Minimal Ringing**: Lanczos-2 reduces common filter artifacts

## Testing

Run the complete test suite:

```bash
cd video-processor
go test -v ./internal/resize
```

Run benchmarks to measure performance:

```bash
go test -bench=. ./internal/resize
```

Run the quality demonstration example:

```bash
go run examples/lanczos_resize_example.go
```

### Test Coverage

- **Unit Tests**: Weight calculation, resize functions, error handling
- **Quality Tests**: Verifies Lanczos filter application and intermediate values
- **Edge Cases**: Nil inputs, invalid dimensions, boundary conditions
- **Performance Tests**: Benchmarks for upsampling and downsampling operations
- **Integration Tests**: End-to-end resize operations with real image data

### Example Output

The example creates demonstration images showing:
- `original_400x400.jpg` - High-resolution source with fine details
- `lanczos_resized_100x100.jpg` - Downsampled with Lanczos quality
- `lanczos_upsampled_200x200.jpg` - Upsampled result

## Performance

Benchmark results on Apple M4 Pro:
- Weight calculation: ~7.4μs for 100→200 pixel scaling
- Downsampling: ~760μs for 100x100→50x50 resize
- Upsampling: ~3.9ms for 100x100→200x200 resize

## API Usage

The resize package can be used programmatically in your Go applications:

```go
package main

import (
    "image"
    "image/jpeg"
    "os"
    "video-processor/internal/resize"
)

func main() {
    // Load an image
    file, err := os.Open("input.jpg")
    if err != nil {
        panic(err)
    }
    defer file.Close()
    
    src, err := jpeg.Decode(file)
    if err != nil {
        panic(err)
    }
    
    // Resize with Lanczos filter
    resized, err := resize.Resize(src, 800, 600)
    if err != nil {
        panic(err)
    }
    
    // Save the result
    output, err := os.Create("output.jpg")
    if err != nil {
        panic(err)
    }
    defer output.Close()
    
    jpeg.Encode(output, resized, &jpeg.Options{Quality: 90})
}
```

### API Functions

#### `resize.Resize(src image.Image, width, height int) (*image.NRGBA, error)`

Main resize function that handles both upsampling and downsampling with Lanczos-2 filter.

**Parameters:**
- `src`: Source image implementing the `image.Image` interface
- `width`: Target width in pixels (must be > 0)
- `height`: Target height in pixels (must be > 0)

**Returns:**
- `*image.NRGBA`: Resized image with full color precision
- `error`: Error if invalid parameters or processing fails

**Features:**
- Automatic optimization for single-dimension resizes
- Two-pass resizing for complex scaling operations
- High-precision color processing (16-bit per channel)
- Lanczos-2 filter for optimal quality/performance balance

## Package Structure

```
video-processor/
├── cmd/main.go              # CLI application
├── internal/
│   ├── filters/filter.go    # Lanczos and other filters
│   └── resize/              
│       ├── resize.go        # Main resize functions
│       └── resize_test.go   # Comprehensive tests
└── examples/
    └── lanczos_resize_example.go  # Quality demonstration
```

## License

MIT
