// Package ascii can convert a image pixel to a raw char
// base on it's RGBA value, in another word, input a image pixel
// output a raw char ascii.
package ascii

import (
	"github.com/aybabtme/rgbterm"
	"image/color"
	"math"
	"reflect"
)

// Options convert pixel to raw char
type Options struct {
	Pixels   []byte
	Reversed bool
	Colored  bool
}

// DefaultOptions that contains the default pixels
var DefaultOptions = Options{
	Pixels:   []byte(" .,:;i1tfLCG08@"),
	Reversed: false,
	Colored:  true,
}

// NewOptions create a new convert option
func NewOptions() Options {
	newOptions := Options{}
	newOptions.mergeOptions(&DefaultOptions)
	return newOptions
}

// mergeOptions merge two options
func (options *Options) mergeOptions(newOptions *Options) {
	options.Pixels = append([]byte{}, newOptions.Pixels...)
	options.Reversed = newOptions.Reversed
	options.Colored = newOptions.Colored
}

// NewPixelConverter create a new pixel converter
func NewPixelConverter() PixelConverter {
	return PixelASCIIConverter{}
}

// PixelConverter define the convert pixel operation
type PixelConverter interface {
	ConvertPixelToASCII(pixel color.Color, options *Options) string
}

// PixelASCIIConverter responsible for pixel ascii conversion
type PixelASCIIConverter struct {
}

// ConvertPixelToASCII converts a pixel to a ASCII char string
func (converter PixelASCIIConverter) ConvertPixelToASCII(pixel color.Color, options *Options) string {
	convertOptions := NewOptions()
	convertOptions.mergeOptions(options)

	if convertOptions.Reversed {
		convertOptions.Pixels = converter.reverse(convertOptions.Pixels)
	}

	r := reflect.ValueOf(pixel).FieldByName("R").Uint()
	g := reflect.ValueOf(pixel).FieldByName("G").Uint()
	b := reflect.ValueOf(pixel).FieldByName("B").Uint()
	a := reflect.ValueOf(pixel).FieldByName("A").Uint()
	value := converter.intensity(r, g, b, a)

	// Choose the char
	precision := float64(255 * 3 / (len(convertOptions.Pixels) - 1))
	rawChar := convertOptions.Pixels[converter.roundValue(float64(value)/precision)]
	if convertOptions.Colored {
		return converter.decorateWithColor(r, g, b, rawChar)
	}
	return string([]byte{rawChar})
}

func (converter PixelASCIIConverter) roundValue(value float64) int {
	return int(math.Floor(value + 0.5))
}

func (converter PixelASCIIConverter) reverse(numbers []byte) []byte {
	for i := 0; i < len(numbers)/2; i++ {
		j := len(numbers) - i - 1
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
	return numbers
}

func (converter PixelASCIIConverter) intensity(r, g, b, a uint64) uint64 {
	return (r + g + b) * a / 255
}

// decorateWithColor decorate the raw char with the color base on r,g,b value
func (converter PixelASCIIConverter) decorateWithColor(r, g, b uint64, rawChar byte) string {
	coloredChar := rgbterm.FgString(string([]byte{rawChar}), uint8(r), uint8(g), uint8(b))
	return coloredChar
}
