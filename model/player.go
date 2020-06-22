package model

type Player struct {
	XPos, YPos    float64
	Speed         float64
	Shoot         bool
	FireCount     int
	MissilesFired []Missile
}

type Missile struct {
	XPos, YPos float64
	Visible    bool
}

type Alien struct {
	XPos, YPos float64
}
