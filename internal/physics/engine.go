package physics

import (
	"math"
	"sync"
	"syscall"
	"time"
	"math/rand"

	behavior "shimeji-pet-companion/internal/behavior"

	"github.com/wailsapp/wails/v3/pkg/application"
)

var (
	user32              = syscall.NewLazyDLL("user32.dll")
	procGetAsyncKeyState = user32.NewProc("GetAsyncKeyState")
)

const VK_LBUTTON = 0x01

func isLButtonPressed() bool {
	val, _, _ := procGetAsyncKeyState.Call(uintptr(VK_LBUTTON))
	return (val & 0x8000) != 0
}

type Config struct {
	TargetFPS     int     // Change this to whatever FPS you want (e.g., 15, 30, 60)
	Gravity       float64 // Downward acceleration in pixels/sec² (e.g. 1800)
	Friction      float64 // Ground friction factor per second (e.g. 0.05)
	AirResistance float64 // Air drag factor per second (e.g. 0.90)
	Bounce        float64 // Restitution factor (0.0 to 1.0)
	FloorOffset   int     // Taskbar offset padding
}

// Framerate-independent defaults (values are expressed per second)
var DefaultConfig = Config{
	TargetFPS:     20,     // Default target FPS
	Gravity:       1800.0, // Pixels per second²
	Friction:      0.02,   // Velocity multiplier when sliding on floor
	AirResistance: 0.85,   // Velocity multiplier in air
	Bounce:        0.28,   // Rebound velocity ratio
	FloorOffset:   48,     // Taskbar padding
}

type Engine struct {
	mu           sync.Mutex
	window       *application.WebviewWindow
	config       Config
	screenWidth  int
	screenHeight int
	winWidth     int
	winHeight    int

	x, y         float64
	vx, vy       float64
	lastSetX     int
	lastSetY     int
	isDragging   bool
	isGrounded   bool
	lastMoveTime time.Time

	fpsChan  chan int
	stopChan chan struct{}

	sm *behavior.StateMachine
}

func NewEngine(app *application.App, win *application.WebviewWindow, screenW, screenH, winW, winH int, cfg *Config) *Engine {
	c := DefaultConfig
	if cfg != nil {
		c = *cfg
	}
	if c.TargetFPS <= 0 {
		c.TargetFPS = 20
	}

	startX := float64(screenW - winW - 60)
	startY := float64(100)

	e := &Engine{
		window:       win,
		config:       c,
		screenWidth:  screenW,
		screenHeight: screenH,
		winWidth:     winW,
		winHeight:    winH,
		x:            startX,
		y:            startY,
		lastSetX:     int(startX),
		lastSetY:     int(startY),
		fpsChan:      make(chan int, 1),
		stopChan:     make(chan struct{}),
	}

	ctx := &behavior.Context{
		App:          app,
		Window:       win,
		ScreenWidth:  screenW,
		ScreenHeight: screenH,
		WinWidth:     winW,
		WinHeight:    winH,
		FloorY:       float64(screenH - c.FloorOffset - winH),
		X:            &e.x,
		Y:            &e.y,
		Vx:           &e.vx,
		Vy:           &e.vy,
		IsGrounded:   &e.isGrounded,
		IsDragging:   &e.isDragging,
		Direction:    behavior.DirRight,
		Rnd:          rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	if win != nil {
		win.SetPosition(int(startX), int(startY))
	}

	e.sm = behavior.NewStateMachine(ctx, &behavior.SitState{})

	return e
}

func (e *Engine) Start() {
	go e.loop()
}

func (e *Engine) Stop() {
	close(e.stopChan)
}

// SetFPS allows changing the loop rate at runtime
func (e *Engine) SetFPS(fps int) {
	if fps < 1 {
		fps = 1
	} else if fps > 120 {
		fps = 120
	}
	e.fpsChan <- fps
}

func (e *Engine) loop() {
	currentFPS := e.config.TargetFPS
	interval := time.Second / time.Duration(currentFPS)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	lastTick := time.Now()

	for {
		select {
		case <-e.stopChan:
			return
		case newFPS := <-e.fpsChan:
			e.mu.Lock()
			e.config.TargetFPS = newFPS
			e.mu.Unlock()
			ticker.Reset(time.Second / time.Duration(newFPS))
		case now := <-ticker.C:
			dt := now.Sub(lastTick).Seconds()
			lastTick = now

			// Clamp dt to avoid huge position jumps if thread pauses
			if dt > 0.1 {
				dt = 0.1
			}
			e.tick(dt)
		}
	}
}

func (e *Engine) StartCustomDrag() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.isDragging = true
	e.isGrounded = false
	e.vx = 0
	e.vy = 0
}

func (e *Engine) MoveCustomDrag(newX, newY int, vx, vy float64) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.x = float64(newX)
	e.y = float64(newY)
	e.lastSetX = newX
	e.lastSetY = newY
	e.vx = vx
	e.vy = vy

	if e.window != nil {
		e.window.SetPosition(newX, newY)
	}
}

func (e *Engine) EndCustomDrag(throwVx, throwVy float64) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.vx = throwVx
	e.vy = throwVy
	e.isDragging = false // Resumes gravity loop at TargetFPS
}

func (e *Engine) tick(dt float64) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.window == nil {
		return
	}

	curX, curY := e.window.Position()
	lDown := isLButtonPressed()
	now := time.Now()

	// 1. Detect native window dragging & calculate throw velocity
	if curX != e.lastSetX || curY != e.lastSetY {
		if lDown {
			moveDt := now.Sub(e.lastMoveTime).Seconds()
			if !e.lastMoveTime.IsZero() && moveDt > 0.01 && moveDt < 0.3 {
				e.vx = (float64(curX) - e.x) / moveDt
				e.vy = (float64(curY) - e.y) / moveDt
			}
			e.x = float64(curX)
			e.y = float64(curY)
			e.lastSetX = curX
			e.lastSetY = curY
			e.lastMoveTime = now
			e.isDragging = true
			e.isGrounded = false
		}
	}

	// 2. Detect drag release & apply throw momentum clamp
	if e.isDragging {
		if !lDown {
			e.isDragging = false
			if now.Sub(e.lastMoveTime) > 120*time.Millisecond {
				e.vx = 0
				e.vy = 0
			} else {
				const maxThrow = 1500.0 // Pixels per second
				e.vx = math.Max(math.Min(e.vx, maxThrow), -maxThrow)
				e.vy = math.Max(math.Min(e.vy, maxThrow), -maxThrow)
			}
		}
	}

	// 3. Update Behavior State Machine (handles transitions, walk, climb, sit)
	if e.sm != nil {
		e.sm.Update(dt)
	}

	// 4. If user is currently holding the window, skip autonomous physics
	if e.isDragging {
		return
	}

	currentState := ""
	if e.sm != nil && e.sm.CurrentState() != nil {
		currentState = e.sm.CurrentState().Name()
	}
	isClimbing := currentState == "climb"

	// 5. Idle check: Save CPU when sitting still on the floor.
	// (Runs AFTER sm.Update so sit timer still counts down to trigger walking)
	if e.isGrounded && (currentState == "sit" || currentState == "idle") && math.Abs(e.vx) < 1.0 && math.Abs(e.vy) < 1.0 {
		return
	}

	// 6. Gravity and air drag (disabled during climb or when grounded)
	if !isClimbing && !e.isGrounded {
		e.vy += e.config.Gravity * dt
		e.vx *= math.Pow(e.config.AirResistance, dt)
		e.vy *= math.Pow(e.config.AirResistance, dt)
	}

	// 7. Integrate Position
	e.x += e.vx * dt
	e.y += e.vy * dt

	// 8. Floor Collision
	floorY := float64(e.screenHeight - e.config.FloorOffset - e.winHeight)
	if e.y >= floorY {
		e.y = floorY
		if math.Abs(e.vy) > 80.0 {
			e.vy = -e.vy * e.config.Bounce
		} else {
			e.vy = 0
			e.isGrounded = true
		}

		// Apply friction only when settling down, not while actively walking or running
		if currentState == "sit" || currentState == "idle" || currentState == "fall" {
			e.vx *= math.Pow(e.config.Friction, dt)
			if math.Abs(e.vx) < 5.0 {
				e.vx = 0
			}
		}
	} else if !isClimbing {
		e.isGrounded = false
	}

	// 9. Wall Collisions (Ignored while climbing so the pet can scale the wall)
	rightWall := float64(e.screenWidth - e.winWidth)
	if !isClimbing {
		if e.x < 0 {
			e.x = 0
			e.vx = -e.vx * e.config.Bounce
		} else if e.x > rightWall {
			e.x = rightWall
			e.vx = -e.vx * e.config.Bounce
		}
	}

	// 10. Update native window position only when integer pixel values change
	newX := int(math.Round(e.x))
	newY := int(math.Round(e.y))
	if newX != e.lastSetX || newY != e.lastSetY {
		e.lastSetX = newX
		e.lastSetY = newY
		e.window.SetPosition(newX, newY)
	}
}