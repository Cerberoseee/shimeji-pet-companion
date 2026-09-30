package behavior

type WalkState struct {
	duration float64
	elapsed  float64
	speed    float64
}

func (s *WalkState) Name() string      { return "walk" }
func (s *WalkState) Animation() string { return "walk" }

func (s *WalkState) Enter(ctx *Context) {
	s.speed = 65.0
	s.duration = 4.0 + ctx.Rnd.Float64()*5.0
	s.elapsed = 0

	if ctx.Direction == 0 {
		ctx.Direction = DirRight
	}
	*ctx.Vx = float64(ctx.Direction) * s.speed
}

func (s *WalkState) Update(ctx *Context, dt float64) PetState {
	s.elapsed += dt
	*ctx.Vx = float64(ctx.Direction) * s.speed

	rightWall := float64(ctx.ScreenWidth - ctx.WinWidth)

	if *ctx.X <= 0 {
		*ctx.X = 0
		if ctx.Rnd.Float64() < 0.6 {
			return &ClimbState{Side: DirLeft}
		}
		ctx.Direction = DirRight
	} else if *ctx.X >= rightWall {
		*ctx.X = rightWall
		if ctx.Rnd.Float64() < 0.6 {
			return &ClimbState{Side: DirRight}
		}
		ctx.Direction = DirLeft
	}

	if s.elapsed >= s.duration {
		return &SitState{}
	}
	return nil
}

func (s *WalkState) Exit(ctx *Context) {}