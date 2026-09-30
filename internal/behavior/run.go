package behavior

type RunState struct {
	duration float64
	elapsed  float64
	speed    float64
}

func (s *RunState) Name() string      { return "run" }
func (s *RunState) Animation() string { return "run" }

func (s *RunState) Enter(ctx *Context) {
	s.speed = 180.0 // Fast burst
	s.duration = 2.0 + ctx.Rnd.Float64()*3.0
	s.elapsed = 0
	*ctx.Vx = float64(ctx.Direction) * s.speed
}

func (s *RunState) Update(ctx *Context, dt float64) PetState {
	s.elapsed += dt
	*ctx.Vx = float64(ctx.Direction) * s.speed

	rightWall := float64(ctx.ScreenWidth - ctx.WinWidth)

	// Wall hit while running -> slide or climb rapidly
	if *ctx.X <= 0 {
		*ctx.X = 0
		return &ClimbState{Side: DirLeft}
	} else if *ctx.X >= rightWall {
		*ctx.X = rightWall
		return &ClimbState{Side: DirRight}
	}

	if s.elapsed >= s.duration {
		return &SitState{} // Slide to sit after running
	}
	return nil
}

func (s *RunState) Exit(ctx *Context) {}