package main

import (
	"fmt"
	"time"
	"math"
	"math/rand"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_ttf"
	"github.com/veandco/go-sdl2/sdl_image"
)

var Rnd *sdl.Renderer

func main() {
	fmt.Println( "yolo" )

	bRun := true

	sdl.Init( sdl.INIT_EVERYTHING )
	ttf.Init()
	img.Init( img.INIT_PNG | img.INIT_JPG )

	rand.Seed( time.Now().UnixNano() )

	win, err := sdl.CreateWindow( "-untitled-",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		640, 360, sdl.WINDOW_SHOWN | sdl.WINDOW_RESIZABLE )

	if err != nil { panic( err ) }
	defer win.Destroy()

	Rnd, err = sdl.CreateRenderer( win, -1, sdl.RENDERER_ACCELERATED |
		sdl.RENDERER_TARGETTEXTURE )//sdl.RENDERER_PRESENTVSYNC )

	if err != nil { panic( err ) }
	defer Rnd.Destroy()

	sdl.SetHint( sdl.HINT_RENDER_SCALE_QUALITY, "linear" )
	// logical size is the larger of the two rt's
	// which should be the 2d rt and widescreen
	Rnd.SetLogicalSize( VID2D_W, VID2D_H )

	// 2d rt
	rndTex, err := Rnd.CreateTexture( sdl.PIXELFORMAT_ARGB8888,
		sdl.TEXTUREACCESS_TARGET, VID2D_W, VID2D_H )
	rndTex.SetBlendMode( sdl.BLENDMODE_BLEND )

	if err != nil { panic( err ) }
	defer rndTex.Destroy()

	// clear backbuffer
	Rnd.SetDrawColor( 0, 0, 0, 255 )
	Rnd.Clear()

	// loadings
	font0, _ := ttf.OpenFont( "fonts/NotoSans.ttf", 16 )

	// fps text
	fpsTxt := NewText( font0, "hello der" )
	defer fpsTxt.texture.Destroy()
	fpsTxt.color = sdl.Color{ 0, 255, 0, 255 }

	// create and initialize 3d render context
	// uses const VID dimensions
	ctx := CreateRenderContext()

	// texture.update takes BGRA Bitmap as RGB
	bmtex, _ := Rnd.CreateTexture( sdl.PIXELFORMAT_RGB888,
		sdl.TEXTUREACCESS_STREAMING, ctx.Bm.Width, ctx.Bm.Height )
	//bmtex.SetBlendMode( sdl.BLENDMODE_BLEND )
	defer bmtex.Destroy()

	//stars := NewStars3D( 3, 64.0, 10.0 )

	// game view pos
	//gv := Vec2{0, 0}

	theTime := float32(0.0)
	var evt sdl.Event
	ftmr := time.Now()
	fpsTmr := time.Now()
	bmdt := float32(0.0)

	var trot float32 = 0
	// texture bitmap
	tex := NewBitmapFromFile( "./hdr.png" )

	// perspective projection matrix
	// aspect should be of the logical size
	// which should be the 2d rt and also widescreen
	// 3d rt can now be any arbritary size and
	// automatically be corrected
	var ProjMat Mat4
	ProjMat.InitPerspective( math.Pi / 180.0 * 66.0,
		float32(VID2D_W) / VID2D_H, 0.1, 1000.0 )

	// main loops
	for bRun == true {
		dt := float32(time.Since( ftmr ).Seconds())
		ftmr = time.Now()
		if dt > 0.1 { dt = 0.1 }
		theTime += dt

		for evt = sdl.PollEvent(); evt != nil; evt = sdl.PollEvent() {
			switch e := evt.(type) {
				case *sdl.QuitEvent:
					bRun = false; println( e )
				case *sdl.KeyDownEvent:
					println( "KEY!" )
			}
		}

		if time.Since( fpsTmr ).Seconds() > 0.1 {
			fpsTxt.SetString( fmt.Sprintf( "%.4g - %.4gms",
				dt * 1000, bmdt ) )
			fpsTmr = time.Now()
		}

		// do drawing to 3d render context bitmap buffer
		// and then update bitmap texture with comp array
		bmdrtmr := time.Now()

		ctx.Bm.Clear( 0x10 )

		// mesh vertices and transform
		v1 := Vertex{ Vec4{ 0, 1, 0, 1 }, Vec4{ 0.5, 0, 0, 1 } }
		v2 := Vertex{ Vec4{ -1, -1, 0, 1 }, Vec4{ 0, 1, 0, 1 } }
		v3 := Vertex{ Vec4{ 1, -1, 0, 1 }, Vec4{ 1, 1, 0, 1 } }

		trot += 0.5 * dt
		var transMat Mat4
		transMat.InitTranslation( 0, 0, 3 )
		var rotMat Mat4
		rotMat.InitRotation( 0, trot, 0 )

		tform := ProjMat.Mul( transMat.Mul( rotMat ) )

		ctx.FillTriangle( v1.Transform( tform ),
			v2.Transform( tform ), v3.Transform( tform ), tex )

		//stars.UpdateAndRender( ctx, dt )

		bmtex.Update( nil, unsafe.Pointer(&ctx.Bm.Comp[0]), ctx.Bm.Width * 4 )
		bmdt = float32(time.Since( bmdrtmr ).Seconds() * 1000)

		// pure sdl 2d drawing
		// set 2d rt for drawing to
		Rnd.SetRenderTarget( rndTex )
		Rnd.SetDrawColor( 0, 0, 0, 0 )
		Rnd.Clear()

		Rnd.SetDrawColor( 255, 127, 0, 127 )
		Rnd.FillRect( &sdl.Rect{ 4, 100, 32, 32 } )

		// hud/ui view
		//Rnd.SetViewport( &sdl.Rect{ 0, 0, VID_W, VID_H } )
		fpsTxt.Draw( Rnd, 8, 8 )

		// copy final render textures to backbuffer
		Rnd.SetRenderTarget( nil )
		Rnd.SetDrawColor( 0, 0, 0, 255 )
		Rnd.Clear()
		// 3d rt
		Rnd.Copy( bmtex, nil, nil )
		// 2d rt
		Rnd.Copy( rndTex, nil, nil )
		Rnd.Present()
	}

	font0.Close()

	img.Quit()
	ttf.Quit()
	sdl.Quit()
}
