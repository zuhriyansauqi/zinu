package zinu

import (
	"encoding/json"
	"image"
	"image/color"
	"io/ioutil"
	"math"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
	"github.com/zuhriyan/zinu/utils"
)

type (

	// Zinu zinu object
	Zinu struct {
		Name       string        `json:"name"`
		Version    int64         `json:"version"`
		Background string        `json:"background"`
		Objects    []*zinuObject `json:"objects"`
		FileDir    string
	}

	zinuObject struct {
		Type        string  `json:"type"`
		Width       float64 `json:"width"`
		Height      float64 `json:"height"`
		X           float64 `json:"x"`
		Y           float64 `json:"y"`
		InverseX    bool
		InverseY    bool
		Value       string  `json:"value"`
		Font        string  `json:"font"`
		FontSize    float64 `json:"size"`
		LineSpacing float64 `json:"lineSpacing"`
		WordWrap    int64   `json:"wordWrap"`
	}
)

var (
	igWidth  int = 1080
	igHeight int = 1920
)

// Generate Generate InstaStrory bitmap
func (z *Zinu) Generate(output string) error {
	bg, err := imaging.Open(filepath.Join(z.FileDir, z.Background))
	if err != nil {
		return err
	}

	dst := imaging.New(igWidth, igHeight, color.NRGBA{0, 0, 0, 0})
	z.fillBackground(bg, &dst)

	for _, obj := range z.Objects {
		if err := z.draw(obj, &dst); err != nil {
			return err
		}
	}

	if err := z.save(output, dst); err != nil {
		return err
	}

	return nil
}

func (z *Zinu) fillBackground(bg image.Image, dst **image.NRGBA) {
	bgFill := imaging.Fill(bg, igWidth, igHeight, imaging.Center, imaging.Lanczos)
	*dst = imaging.Paste(*dst, bgFill, image.Pt(0, 0))
}

func (z *Zinu) draw(obj *zinuObject, dst **image.NRGBA) error {
	if obj.Type == "text" {
		return z.drawText(obj, dst)
	}
	return z.drawImage(obj, dst)
}

func (z *Zinu) drawText(obj *zinuObject, dst **image.NRGBA) error {
	var posX, posY float64

	fontBytes, err := ioutil.ReadFile(filepath.Join(z.FileDir, obj.Font))
	if err != nil {
		return err
	}

	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return err
	}

	var text = utils.WordWrap(obj.Value, int(obj.WordWrap))
	var lines = utils.Tokenize(text, "\n")
	var count = len(lines)

	obj.Height = (obj.FontSize * float64(count)) + (float64(count-1) * obj.FontSize * (obj.LineSpacing / obj.FontSize))

	if posX = obj.X; obj.InverseX {
		posX = math.Abs(posX) - obj.Width
	}

	if posY = obj.Y; obj.InverseY {
		posY = float64(igHeight) - math.Abs(posY) - obj.Height
	}

	canvas := imaging.New(int(obj.Width), int(obj.Height), color.NRGBA64{0, 0, 0, 0})

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(f)
	c.SetFontSize(obj.FontSize)
	c.SetClip(canvas.Bounds())
	c.SetDst(canvas)
	c.SetSrc(image.White)

	pt := freetype.Pt(0, int(c.PointToFixed(obj.FontSize)>>6))
	for _, s := range lines {
		_, err := c.DrawString(s, pt)
		if err != nil {
			return err
		}

		pt.Y += c.PointToFixed((obj.LineSpacing / obj.FontSize) * obj.FontSize)
	}

	*dst = imaging.Overlay(*dst, canvas, image.Pt(int(posX), int(posY)), 1)

	return nil
}

func (z *Zinu) drawImage(obj *zinuObject, dst **image.NRGBA) error {
	var posX, posY float64

	if posX = obj.X; obj.InverseX {
		posX = math.Abs(posX) - obj.Width
	}

	if posY = obj.Y; obj.InverseY {
		posY = math.Abs(posY) - obj.Height
	}

	img, err := imaging.Open(filepath.Join(z.FileDir, obj.Value))
	if err != nil {
		return err
	}

	*dst = imaging.Overlay(*dst, img, image.Pt(int(posX), int(posY)), 1)

	return nil
}

func (z *Zinu) save(output string, dst *image.NRGBA) error {
	if err := imaging.Save(dst, output+".jpg", imaging.JPEGQuality(80)); err != nil {
		return err
	}

	return nil
}

// Load load zinu template from json file
func Load(filename string) (*Zinu, error) {
	path := filepath.Dir(filename)

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	z := new(Zinu)
	if err := json.Unmarshal(b, z); err != nil {
		return nil, err
	}

	for _, obj := range z.Objects {
		if obj.InverseX = false; obj.X < 0 {
			obj.InverseX = true
		}
		if obj.InverseY = false; obj.Y < 0 {
			obj.InverseY = true
		}
	}

	z.FileDir = path

	return z, nil
}
