import { useState } from 'react'
import { ShimejiService } from '../bindings/shimeji-pet-companion/internal/shimeji'
import { ShimejiSprite } from './components/Shimeji/Shimeji'
import { SpeechBubble } from './components/SpeechBubble/SpeechBubble'
import { useSpeechBubble } from './hooks/useSpeechBubble'
import { ShimejiAnimation } from './types/shimeji'
import styles from './styles/App.module.css'
import { useThrottledDrag } from './hooks/useThrottleDrag'

export function App() {
  const [loading, setLoading] = useState(false)
  const [animation, setAnimation] = useState<ShimejiAnimation>('idle')

  const { isOpen, text, showBubble, hideBubble } = useSpeechBubble(4500)

  const { dragProps } = useThrottledDrag({ fps: 24 }) // Change to 15, 20, 24, etc.

  // Inside handlePoke after receiving Ollama's response:
  const handlePoke = async () => {
    if (loading) return

    setLoading(true)
    showBubble('...', 30000)

    try {
      const res = await ShimejiService.AskOllama('The user poked you on the head!')
      if (!res) throw new Error('No response from Ollama')

      setAnimation((res.animation as ShimejiAnimation) || 'idle')
      showBubble(res.subtitle_en)

    } catch (err) {
      console.error('Inference error:', err)
      hideBubble()
      setAnimation('idle')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className={styles.container} onClick={handlePoke} {...dragProps}>
      <div className={styles.noDrag}>
        <SpeechBubble text={text} visible={isOpen} />
      </div>
      <ShimejiSprite animation={animation} loading={loading} />
    </div>
  )
}

export default App