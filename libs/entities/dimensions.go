package entities

type Dimensions struct {
	Width  float32
	Height float32
	Depth  float32
}

func (d *Dimensions) ToArray() [3]float32 {
	return [3]float32{d.Width, d.Height, d.Depth}
}
