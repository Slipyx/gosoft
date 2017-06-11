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
		960, 540, sdl.WINDOW_SHOWN | sdl.WINDOW_RESIZABLE | sdl.WINDOW_OPENGL )

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
	tex := NewBitmapFromFile( "./tex.png" )
	//toptex := NewBitmapFromFile( "./top.png" )
	// obj mesh, front face is CW
	objmesh := NewMesh()//LoadOBJMesh( "./mesh.obj" )
	objmesh.Vertices = append( objmesh.Vertices,
		Vertex{ Vec4{ -0.5, 0, 0.5, 1 }, Vec4{ 0, 0, 0, 1 } },
		Vertex{ Vec4{ 0.5, 0, 0.5, 1 }, Vec4{ 100, 0, 0, 1 } },
		Vertex{ Vec4{ 0.5, 0, -0.5, 1 }, Vec4{ 100, 100, 0, 1 } },
		Vertex{ Vec4{ -0.5, 0, -0.5, 1 }, Vec4{ 0, 100, 0, 1 } }, )
	objmesh.Indices = append( objmesh.Indices,
		0, 1, 2, 0, 2, 3 )

	// transform
	tform := NewTransform( Vec3{ 0, 0, 0 } )
	tform.Scale = Vec3{ 100, 1, 100 }

	// perspective projection matrix
	// aspect should be of the logical size
	// which should be the 2d rt and also widescreen
	// 3d rt can now be any arbritary size and
	// automatically be corrected
	aspect := float32(VID_W) / float32(VID_H)
	// fovh in degrees
	cam := NewCamera( Vec3{ 0, 0.5, 10 }, 90, aspect, 0.1, 1024.0 )

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
		ctx.ClearDepthBuffer()

		trot += dt

		//tform.Rot.Y = trot/4
		cam.Pos.Z += dt*2
		cam.Pos.X += dt

		objmesh.Draw( ctx, cam.GetViewProj().Mul( tform.GetModel() ), tex )

		//stars.UpdateAndRender( ctx, dt )

		bmtex.Update( nil, unsafe.Pointer(&ctx.Bm.Comp[0]), ctx.Bm.Width * 4 )
		bmdt = float32(time.Since( bmdrtmr ).Seconds() * 1000)

		// set 2d rt for drawing to
		Rnd.SetRenderTarget( rndTex )
		//Rnd.SetDrawColor( 255, 0, 255, 0 )
		//Rnd.Clear()

		// copy 3d framebuffer first to entire render texture
		Rnd.Copy( bmtex, nil, nil )

		// now do sdl's 2d drawing over the 3d buffer
		Rnd.SetDrawColor( 255, 127, 0, 64 )
		Rnd.FillRect( &sdl.Rect{ 4, 100, 32, 32 } )

		// hud/ui view
		//Rnd.SetViewport( &sdl.Rect{ 0, 0, VID_W, VID_H } )
		fpsTxt.Draw( Rnd, 8, 8 )

		// copy final render texture to display
		Rnd.SetRenderTarget( nil )
		Rnd.SetDrawColor( 0, 0, 0, 255 )
		Rnd.Clear()
		// rt
		Rnd.Copy( rndTex, nil, nil )
		Rnd.Present()
	}

	fmt.Println( math.Pi )

	font0.Close()

	img.Quit()
	ttf.Quit()
	sdl.Quit()
}
