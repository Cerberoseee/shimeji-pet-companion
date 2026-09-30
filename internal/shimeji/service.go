package shimeji

import (
	config "shimeji-pet-companion/internal/config"
	physics "shimeji-pet-companion/internal/physics"
)

type ShimejiService struct {
	cfg             *config.Config
	promptGenerator *config.PromptGenerator
	physics         *physics.Engine
}

func NewShimejiService(cfg *config.Config, physics *physics.Engine) *ShimejiService {
	if cfg == nil {
		cfg = config.LoadConfig("config.yml")
	}

	return &ShimejiService{
		cfg:             cfg,
		promptGenerator: config.NewPromptGenerator(""),
		physics:         physics,
	}
}

func (s *ShimejiService) SetPhysics(phy *physics.Engine) {
	s.physics = phy
}

func (s *ShimejiService) StartCustomDrag() {
	if s.physics != nil {
		s.physics.StartCustomDrag()
	}
}

func (s *ShimejiService) MoveCustomDrag(x, y int, vx, vy float64) {
	if s.physics != nil {
		s.physics.MoveCustomDrag(x, y, vx, vy)
	}
}

func (s *ShimejiService) EndCustomDrag(vx, vy float64) {
	if s.physics != nil {
		s.physics.EndCustomDrag(vx, vy)
	}
}