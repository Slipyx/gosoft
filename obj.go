package main

import (
	"os"
	"strings"
	"bufio"
	"strconv"
	"fmt"
)

type OBJMesh struct {
	v []Vec4
	vt []Vec4

	f [][3]int
}

// load an obj file and convert to indexed Mesh

func LoadOBJMesh( file string ) *Mesh {
	var om *OBJMesh = &OBJMesh{}

	om.v = make( []Vec4, 0 )
	om.vt = make( []Vec4, 0 )
	om.f = make( [][3]int, 0 )

	objFile, err := os.Open( file )
	if err != nil { panic( err ) }
	defer objFile.Close()

	objReader := bufio.NewReader( objFile )

	// load objmesh structure

	// assumptions:
	// len of om.f should be multiple of 3
	// len of om.v and om.vt should be equal
	// no mtl, multitexture

	for true {
		var objLine string
		objLine, err := objReader.ReadString( '\n' )
		if err != nil { break }
		objFlds := strings.Fields( objLine )

		if len( objFlds ) == 0 { continue }

		if objFlds[0] == "v" {
			var v Vec4
			var ffloat float64

			ffloat, _ = strconv.ParseFloat( objFlds[1], 32 )
			v.X = float32(ffloat)
			ffloat, _ = strconv.ParseFloat( objFlds[2], 32 )
			v.Y = float32(ffloat)
			ffloat, _ = strconv.ParseFloat( objFlds[3], 32 )
			v.Z = float32(ffloat)
			v.W = 1.0

			om.v = append( om.v, v )
		}

		if objFlds[0] == "vt" {
			var vtx, vty float64

			vtx, _ = strconv.ParseFloat( objFlds[1], 32 )
			vty, _ = strconv.ParseFloat( objFlds[2], 32 )
			/*for vtx < 0 { vtx += 1 }
			for vtx > 1 { vtx -= 1 }
			for vty < 0 { vty += 1 }
			for vty > 1 { vty -= 1 }*/
			vty = 1 - vty
			//vtx = 1 - vtx

			om.vt = append( om.vt, Vec4{ float32(vtx), float32(vty), 0, 1 } )
		}

		// 1 based index
		if objFlds[0] == "f" {
			var f [3][3]int

			fsplit := strings.Split( objFlds[1], "/" )
			f[0][0], _ = strconv.Atoi( fsplit[0] )
			f[0][1], f[0][2] = f[0][0], f[0][0]
			if len( fsplit ) > 1 {
				f[0][1], _ = strconv.Atoi( fsplit[1] )
			}
			if len( fsplit ) > 2 {
				f[0][2], _ = strconv.Atoi( fsplit[2] )
			}

			fsplit = strings.Split( objFlds[2], "/" )
			f[1][0], _ = strconv.Atoi( fsplit[0] )
			f[1][1], f[1][2] = f[1][0], f[1][0]
			if len( fsplit ) > 1 {
				f[1][1], _ = strconv.Atoi( fsplit[1] )
			}
			if len( fsplit ) > 2 {
				f[1][2], _ = strconv.Atoi( fsplit[2] )
			}

			fsplit = strings.Split( objFlds[3], "/" )
			f[2][0], _ = strconv.Atoi( fsplit[0] )
			f[2][1], f[2][2] = f[2][0], f[2][0]
			if len( fsplit ) > 1 {
				f[2][1], _ = strconv.Atoi( fsplit[1] )
			}
			if len( fsplit ) > 2 {
				f[2][2], _ = strconv.Atoi( fsplit[2] )
			}

			om.f = append( om.f, f[0], f[1], f[2] )
		}
	}

	// convert to indexed Mesh
	m := &Mesh{}

	m.Vertices = make( []Vertex, len( om.v ) )
	m.Indices = make( []int, len( om.f ) )

	// TRY THIS BRO
	// go through each om.f first

	// assign vertex positions
	for i := 0; i < len( m.Vertices ); i++ {
		m.Vertices[i] = Vertex{ om.v[i], Vec4{0,0,0,1} }
	}

	// assign tc's
	for i := 0; i < len( om.f ); i += 3 {
		m.Vertices[om.f[i][0] - 1].TexCoord = om.vt[om.f[i][1] - 1]
		m.Vertices[om.f[i + 1][0] - 1].TexCoord = om.vt[om.f[i + 1][1] - 1]
		m.Vertices[om.f[i + 2][0] - 1].TexCoord = om.vt[om.f[i + 2][1] - 1]
	}

	// copy indices
	for i := 0; i < len( m.Indices ); i += 3 {
		m.Indices[i] = om.f[i][0] - 1
		m.Indices[i + 1] = om.f[i + 1][0] - 1
		m.Indices[i + 2] = om.f[i + 2][0] - 1
	}

	fmt.Printf( "%v\n", m.Vertices[m.Indices[0]] )

	return m
}
