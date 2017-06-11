package main

type Transform struct {
	Pos, Rot, Scale Vec3
}

func NewTransform( pos Vec3 ) Transform {
	t := Transform{
		pos, Vec3{ 0, 0, 0 }, Vec3{ 1, 1, 1 },
	}

	return t
}

func (t Transform) GetModel() Mat4 {
	var posMat Mat4; posMat.InitTranslation( t.Pos.X, t.Pos.Y, t.Pos.Z )
	var rotMat Mat4; rotMat.InitRotation( t.Rot.X, t.Rot.Y, t.Rot.Z )
	var scaleMat Mat4; scaleMat.InitScale( t.Scale.X, t.Scale.Y, t.Scale.Z )

	return posMat.Mul( rotMat ).Mul( scaleMat )
}

