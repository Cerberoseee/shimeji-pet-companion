package behavior

import (
	"math/rand"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type Direction int

const (
	DirLeft  Direction = -1
	DirRight Direction = 1
)

type PetState interface {
	Name() string
	Animation() string
	Enter(ctx *Context)
	Update(ctx *Context, dt float64) PetState
	Exit(ctx *Context)
}

type Context struct {
	App          *application.App
	Window       *application.WebviewWindow
	ScreenWidth  int
	ScreenHeight int
	WinWidth     int
	WinHeight    int
	FloorY       float64

	X, Y        *float64
	Vx, Vy      *float64
	IsGrounded  *bool
	IsDragging  *bool
	Direction   Direction

	Rnd *rand.Rand
}

type StateMachine struct {
	current PetState
	ctx     *Context
}

func NewStateMachine(ctx *Context, initial PetState) *StateMachine {
	sm := &StateMachine{ctx: ctx}
	sm.TransitionTo(initial)
	return sm
}

func (sm *StateMachine) CurrentState() PetState {
	return sm.current
}

func (sm *StateMachine) TransitionTo(next PetState) {
	if next == nil {
		return
	}
	if sm.current != nil {
		sm.current.Exit(sm.ctx)
	}
	sm.current = next
	sm.current.Enter(sm.ctx)

	// Emit directly to the pet's WebviewWindow
	if sm.ctx.Window != nil {
		sm.ctx.Window.EmitEvent("pet:state_changed", map[string]any{
			"state":     next.Name(),
			"animation": next.Animation(),
			"direction": sm.ctx.Direction,
		})
	}
}

func (sm *StateMachine) Update(dt float64) {
	if sm.current == nil {
		return
	}

	if *sm.ctx.IsDragging {
		if sm.current.Name() != "drag" {
			sm.TransitionTo(&DragState{})
		}
		return
	}

	if !*sm.ctx.IsGrounded && sm.current.Name() != "climb" && sm.current.Name() != "fall" && sm.current.Name() != "drag" {
		sm.TransitionTo(&FallState{})
		return
	}

	if next := sm.current.Update(sm.ctx, dt); next != nil {
		sm.TransitionTo(next)
	}
}