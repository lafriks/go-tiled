package tiled

type aliasGroup Group
type aliasImageLayer ImageLayer
type internalLayer Layer
type aliasLayer struct {
	internalLayer
	// Layer data in raw format
	Data *Data `xml:"data"`
}
type aliasMap Map
type aliasObject Object
type aliasText Text

// SetDefaults provides default values for Group.
func (a *aliasGroup) SetDefaults() {
	a.Opacity = 1
	a.Visible = true
}

// SetDefaults provides default values for ImageLayer.
func (a *aliasImageLayer) SetDefaults() {
	a.Opacity = 1
	a.Visible = true
}

// SetDefaults provides default values for Layer.
func (a *aliasLayer) SetDefaults() {
	a.internalLayer.Opacity = 1
	a.internalLayer.Visible = true
}

// SetDefaults provides default values for Map.
func (a *aliasMap) SetDefaults() {
	a.RenderOrder = "right-down"
}

// SetDefaults provides default values for Object.
func (a *aliasObject) SetDefaults() {
	a.Visible = true
}

// SetDefaults provides default values for Text.
func (a *aliasText) SetDefaults() {
	a.FontFamily = "sans-serif"
	a.Size = 16
	a.Kerning = true
	a.HAlign = "left"
	a.VAlign = "top"
	a.Color = &HexColor{}
}
