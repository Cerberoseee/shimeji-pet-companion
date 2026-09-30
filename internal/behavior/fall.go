package behavior

type FallState struct{}

func (s *FallState) Name() string      { return "fall" }
func (s *FallState) Animation() string { return "fall" }
func (s *FallState) Enter(ctx *Context) {}
func (s *FallState) Update(ctx *Context, dt float64) PetState {
	if *ctx.IsGrounded {
		return &SitState{}
	}
	return nil
}
func (s *FallState) Exit(ctx *Context) {}