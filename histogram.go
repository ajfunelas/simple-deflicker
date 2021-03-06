package main

import (
	"image"
	"image/color"

	"github.com/disintegration/imaging"
)

func generateHistogramFromImage(input image.Image) histogram {
	var histogram histogram
	var pixels uint64
	for y := input.Bounds().Min.Y; y < input.Bounds().Max.Y; y++ {
		for x := input.Bounds().Min.X; x < input.Bounds().Max.X; x++ {
			r, g, b, _ := input.At(x, y).RGBA()
			intensity := float64(0.2126*float64(r)+0.7152*float64(g)+0.0722*float64(b)) / 256.0
			histogram[int(intensity)]++
			pixels++
		}
	}
	return histogram
}

func convertToCumulativeHistogram(input histogram) histogram {
	var targetHistogram histogram
	targetHistogram[0] = input[0]
	for i := 1; i < 256; i++ {
		targetHistogram[i] = targetHistogram[i-1] + input[i]
	}
	return targetHistogram
}

func generateLutFromHistograms(current histogram, target histogram) lut {
	currentCumulativeHistogram := convertToCumulativeHistogram(current)
	targetCumulativeHistogram := convertToCumulativeHistogram(target)

	ratio := float64(currentCumulativeHistogram[255]) / float64(targetCumulativeHistogram[255])
	for i := 0; i < 256; i++ {
		targetCumulativeHistogram[i] = uint32(0.5 + float64(targetCumulativeHistogram[i])*ratio)
	}

	//Generate LUT
	var lut lut
	var p uint8 = 0
	for i := 0; i < 256; i++ {
		for targetCumulativeHistogram[p] < currentCumulativeHistogram[i] {
			p++
		}
		lut[i] = p
	}
	return lut
}
func applyLutToImage(input image.Image, lut lut) image.Image {
	result := imaging.AdjustFunc(input, func(c color.NRGBA) color.NRGBA {
		c.R = uint8(lut[c.R])
		c.G = uint8(lut[c.G])
		c.B = uint8(lut[c.B])
		return c
	})
	return result
}
