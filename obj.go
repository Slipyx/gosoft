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

	// v/vt/vn v/vt/vn v/vt/vn
	f [][3][3]int
}

func OBJVHash( v [3]int ) int {
	res := 17
	res = 31 * res + v[0]
	res = 31 * res + v[1]
	res = 31 * res + v[2]
	return res
}

// what even
type Vinf struct {
	vi int
	vert *Vertex
}

// load an obj file and convert to indexed Mesh

func LoadOBJMesh( file string ) *Mesh {
	var om *OBJMesh = &OBJMesh{}

	om.v = make( []Vec4, 0 )
	om.vt = make( []Vec4, 0 )
	om.f = make( [][3][3]int, 0 )

	objFile, err := os.Open( file )
	if err != nil { panic( err ) }
	defer objFile.Close()

	objReader := bufio.NewReader( objFile )

	// load objmesh structure

	// assumptions:
	// len of om.f should be multiple of 3
	// no mtl, multitexture, NO QUADS

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

			om.f = append( om.f, f )
		}
	}

	// convert to indexed Mesh
	m := &Mesh{}

	m.Vertices = make( []Vertex, 0 )
	m.Indices = make( []int, 0 )
	vi := 0

	objVMap := make( map[int]*Vinf )

	for _, f := range om.f {
		for _, v := range f {
			hash := OBJVHash( v )
			if objVMap[hash] == nil {
				vert := Vertex{ om.v[ v[0] - 1 ], om.vt[ v[1] - 1 ] }
				m.Vertices = append( m.Vertices, vert )
				objVMap[hash] = &Vinf{ vi, &vert }
				m.Indices = append( m.Indices, vi )
				vi += 1
			} else {
				m.Indices = append( m.Indices, objVMap[hash].vi )
			}
		}
	}

	fmt.Printf( "%v\n", m.Vertices[m.Indices[0]] )

	return m
}

