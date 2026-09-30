package behavior

// DragState activates when the user clicks and drags the pet window.
type DragState struct{}

func (s *DragState) Name() string {
	return "drag"
}

func (s *DragState) Animation() string {
	// Triggers the dangling/pinched sprite pose in the frontend
	return "drag"
}

func (s *DragState) Enter(ctx *Context) {
	// Kill autonomous walking velocities so it doesn't fight the cursor
	*ctx.Vx = 0
	*ctx.Vy = 0
}

func (s *DragState) Update(ctx *Context, dt float64) PetState {
	// When the user releases the left mouse button
	if !*ctx.IsDragging {
		if *ctx.IsGrounded {
			// Placed directly on the floor/taskbar
			return &SitState{}
		}
		// Released in mid-air -> let gravity & throw momentum take over
		return &FallState{}
	}

	return nil
}

func (s *DragState) Exit(ctx *Context) {
	// Clean up or trigger sound effects on drop if desired
}
