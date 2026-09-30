import { useRef, useCallback } from 'react'
import { ShimejiService } from '../../bindings/shimeji-pet-companion/internal/shimeji'

interface DragConfig {
  fps?: number // Desired drag update rate (e.g., 15, 20, 30)
}

export function useThrottledDrag({ fps = 20 }: DragConfig = {}) {
  const isDragging = useRef(false)
  const clickOrigin = useRef<{ x: number; y: number }>({ x: 0, y: 0 })
  const lastScreenPos = useRef<{ x: number; y: number; time: number }>({ x: 0, y: 0, time: 0 })
  const lastThrottleTime = useRef(0)
  const currentVelocity = useRef<{ vx: number; vy: number }>({ vx: 0, vy: 0 })

  const frameInterval = 1000 / fps // e.g. 50ms for 20 FPS

  const onPointerDown = useCallback((e: React.PointerEvent<HTMLDivElement>) => {
    if (e.button !== 0) return

    e.currentTarget.setPointerCapture(e.pointerId)
    isDragging.current = true

    clickOrigin.current = {
      x: e.clientX,
      y: e.clientY,
    }

    lastScreenPos.current = {
      x: e.screenX,
      y: e.screenY,
      time: performance.now(),
    }

    lastThrottleTime.current = performance.now()
    ShimejiService.StartCustomDrag()
  }, [])

  const onPointerMove = useCallback(
    (e: React.PointerEvent<HTMLDivElement>) => {
      if (!isDragging.current) return

      const now = performance.now()

      // Calculate throw velocity over desktop screen coordinates
      const dt = Math.max((now - lastScreenPos.current.time) / 1000, 0.016)
      currentVelocity.current = {
        vx: (e.screenX - lastScreenPos.current.x) / dt,
        vy: (e.screenY - lastScreenPos.current.y) / dt,
      }
      lastScreenPos.current = { x: e.screenX, y: e.screenY, time: now }

      // THROTTLE: Only update window position once every `frameInterval` ms
      if (now - lastThrottleTime.current >= frameInterval) {
        lastThrottleTime.current = now

        const targetX = Math.round(e.screenX - clickOrigin.current.x)
        const targetY = Math.round(e.screenY - clickOrigin.current.y)

        ShimejiService.MoveCustomDrag(
          targetX,
          targetY,
          currentVelocity.current.vx,
          currentVelocity.current.vy
        )
      }
    },
    [frameInterval]
  )

  const onPointerUp = useCallback((e: React.PointerEvent<HTMLDivElement>) => {
    if (!isDragging.current) return

    if (e.currentTarget.hasPointerCapture(e.pointerId)) {
      e.currentTarget.releasePointerCapture(e.pointerId)
    }

    isDragging.current = false

    // Clamp throw velocity
    const maxThrow = 1500
    const vx = Math.max(Math.min(currentVelocity.current.vx, maxThrow), -maxThrow)
    const vy = Math.max(Math.min(currentVelocity.current.vy, maxThrow), -maxThrow)

    ShimejiService.EndCustomDrag(vx, vy)
  }, [])

  return {
    dragProps: {
      onPointerDown,
      onPointerMove,
      onPointerUp,
      onPointerCancel: onPointerUp,
    },
  }
}