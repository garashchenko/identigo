package identigo

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
)

// Identicon is a struct with picture and internal pixel grid details
type Identicon struct {
	picSide        int
	picBorder      int
	squareRowCount int
	squareSide     int
	squareColor    color.Color
	fillSquare     []bool
}

// NewIdenticon creates Identicon (side x side) with (gridRowCount x gridRowCount) squares based on (key) hash
func NewIdenticon(side int, gridRowCount int, key string) *Identicon {
	squareSide := side / gridRowCount
	border := (side - squareSide*gridRowCount) / 2

	hash := hashString(key)

	// filling only half of the grid because of the symmetry
	// if gridRowCount is odd we need to fill one more column (the center one)
	fillCellCount := gridRowCount * (gridRowCount/2 + gridRowCount%2)
	squareColor, _ := getColorFromBytes(hash[:3])

	return &Identicon{
		picSide:        side,
		picBorder:      border,
		squareRowCount: gridRowCount,
		squareSide:     squareSide,
		squareColor:    squareColor,
		fillSquare:     getCellsToFill(hash[3:], fillCellCount),
	}
}

func hashString(str string) (result []byte) {
	sha := sha256.New()

	sha.Write([]byte(str))
	result = sha.Sum(nil)

	return result
}

func getColorFromBytes(hash []byte) (c color.Color, err error) {
	if len(hash) < 3 {
		return nil, fmt.Errorf("identigo: not enough bytes for color determination")
	}
	return color.RGBA{hash[0], hash[1], hash[2], 255}, nil
}

func getCellsToFill(hash []byte, cellCount int) (fillCell []bool) {
	processedCells := 0
outerLoop:
	for _, hashByte := range hash {
		for i := 0; i < 8; i++ {
			fillCell = append(fillCell, hashByte&1 == 1)

			hashByte = hashByte >> 1

			processedCells++

			if processedCells == cellCount {
				break outerLoop
			}
		}
	}
	return fillCell
}

func createBackground(side int, color color.Color) *image.RGBA {
	pic := image.NewRGBA(image.Rect(0, 0, side, side))
	draw.Draw(pic, pic.Bounds(), &image.Uniform{color}, image.ZP, draw.Src)
	return pic
}

func drawSquare(pic *image.RGBA, side int, x int, y int, color color.Color) {
	x0 := x
	y0 := y
	x1 := x0 + side
	y1 := y0 + side
	square := image.NewRGBA(image.Rect(x0, y0, x1, y1))
	draw.Draw(pic, square.Bounds(), &image.Uniform{color}, image.ZP, draw.Src)
}

func getCoord(index, picBorder, squareSide, squareRowCount int) (x int, y int) {
	var rowCount, columnCount int
	rowCount = index % squareRowCount
	columnCount = index / squareRowCount
	x = picBorder + squareSide*columnCount
	y = picBorder + squareSide*rowCount
	return x, y
}

func getSymmetricCoord(x, y, picSide int, squareSide int) (xSymm, ySymm int) {
	xSymm = picSide - x - squareSide
	ySymm = y
	return xSymm, ySymm
}

// Render method generates PNG picture
func (icon *Identicon) Render() []byte {

	squareColor := icon.squareColor

	white := color.RGBA{255, 255, 255, 255}
	image := createBackground(icon.picSide, white)

	for index, fill := range icon.fillSquare {

		if fill == false {
			continue
		}

		x0, y0 := getCoord(index, icon.picBorder, icon.squareSide, icon.squareRowCount)
		drawSquare(image, icon.squareSide, x0, y0, squareColor)

		// draw a square mirrored verically if it makes sense
		// (it may not if gridRowCount is odd and we are filling the center column)
		xSymm, ySymm := getSymmetricCoord(x0, y0, icon.picSide, icon.squareSide)
		if xSymm != x0 {
			drawSquare(image, icon.squareSide, xSymm, ySymm, icon.squareColor)
		}
	}

	var buf bytes.Buffer
	png.Encode(&buf, image)
	return buf.Bytes()

}
