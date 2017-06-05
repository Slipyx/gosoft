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

func (m *Mesh) Draw( ctx *RenderContext, transform Mat4, texture *Bitmap ) {
	for i := 0; i < len( m.Indices ); i += 3 {
		ctx.DrawTriangle( m.Vertices[m.Indices[i]].Transform( transform ),
			m.Vertices[m.Indices[i + 1]].Transform( transform ),
			m.Vertices[m.Indices[i + 2]].Transform( transform ), texture )
	}
}

