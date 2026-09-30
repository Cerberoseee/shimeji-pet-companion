package behavior

type ClimbState struct {
	Side         Direction
	climbSpeed   float64
	targetHeight float64
}

func (s *ClimbState) Name() string      { return "climb" }
func (s *ClimbState) Animation() string { return "climb" }

func (s *ClimbState) Enter(ctx *Context) {
	ctx.Direction = s.Side
	s.climbSpeed = -75.0
	*ctx.Vx = 0
	s.targetHeight = float64(ctx.ScreenHeight) * (0.2 + ctx.Rnd.Float64()*0.4)
}

func (s *ClimbState) Update(ctx *Context, dt float64) PetState {
	if s.Side == DirLeft {
		*ctx.X = 0
	} else {
		*ctx.X = float64(ctx.ScreenWidth - ctx.WinWidth)
	}

	*ctx.Vy = s.climbSpeed

	if *ctx.Y <= s.targetHeight || *ctx.Y <= 20 {
		*ctx.Vx = -float64(s.Side) * 120.0
		*ctx.Vy = -50.0
		*ctx.IsGrounded = false
		return &FallState{}
	}

	return nil
}

func (s *ClimbState) Exit(ctx *Context) {}