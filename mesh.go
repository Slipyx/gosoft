package main

// CW is front face

type Mesh struct {
	Vertices []Vertex
	Indices []int
}

func NewMesh() *Mesh {
	m := &Mesh{}

	m.Vertices = make( []Vertex, 0 )
	m.Indices = make( []int, 0 )

	return m
}

