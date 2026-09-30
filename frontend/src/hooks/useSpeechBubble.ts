import { useState, useRef, useCallback, useEffect } from 'react'

export function useSpeechBubble(defaultDurationMs = 4500) {
  const [isOpen, setIsOpen] = useState(false)
  const [text, setText] = useState('')
  const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null)

  const showBubble = useCallback((message: string, durationMs = defaultDurationMs) => {
    if (timerRef.current) {
      clearTimeout(timerRef.current)
    }

    setText(message)
    setIsOpen(true)

    timerRef.current = setTimeout(() => {
      setIsOpen(false)
      timerRef.current = null
    }, durationMs)
  }, [defaultDurationMs])

  const hideBubble = useCallback(() => {
    if (timerRef.current) {
      clearTimeout(timerRef.current)
    }
    setIsOpen(false)
  }, [])

  useEffect(() => {
    return () => {
      if (timerRef.current) clearTimeout(timerRef.current)
    }
  }, [])

  return { isOpen, text, showBubble, hideBubble }
}