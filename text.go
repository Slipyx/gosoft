package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_ttf"
)

// text type
type Text struct {
    texture *sdl.Texture
    str string
    font *ttf.Font
    color sdl.Color
}

func NewText( font *ttf.Font, str string ) *Text {
    ntxt := &Text{
        nil, "", font, sdl.Color{255,255,255,255},
    }
    ntxt.SetString( str )
    return ntxt
}

func (self *Text) SetString( str string ) {
    self.str = str
    self.texture.Destroy()
    txsfc, _ := self.font.RenderUTF8_Blended( self.str, self.color )
    self.texture, _ = Rnd.CreateTextureFromSurface( txsfc )
    txsfc.Free()
}

func (self *Text) Draw( rnd *sdl.Renderer, x, y int32 ) {
    _, _, tw, th, _ := self.texture.Query()
    rnd.Copy( self.texture, nil, &sdl.Rect{ x, y, tw, th } )
}

