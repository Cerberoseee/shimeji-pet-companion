package behavior

type SitState struct {
	duration float64
	elapsed  float64
}

func (s *SitState) Name() string      { return "sit" }
func (s *SitState) Animation() string { return "sit" }

func (s *SitState) Enter(ctx *Context) {
	*ctx.Vx = 0
	*ctx.Vy = 0
	s.duration = 3.0 + ctx.Rnd.Float64()*4.0
	s.elapsed = 0
}

func (s *SitState) Update(ctx *Context, dt float64) PetState {
	s.elapsed += dt
	if s.elapsed >= s.duration {
		if ctx.Rnd.Float64() < 0.7 {
			return &WalkState{}
		}
		return &RunState{}
	}
	return nil
}

func (s *SitState) Exit(ctx *Context) {}