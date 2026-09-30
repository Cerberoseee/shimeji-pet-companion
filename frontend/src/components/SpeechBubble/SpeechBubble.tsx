interface SpeechBubbleProps {
  text: string
  visible: boolean
}
import styles from './SpeechBubble.module.css'

export function SpeechBubble({ text, visible }: SpeechBubbleProps) {
  if (!visible || !text) return null

  return (
    <div className={styles.speechBubble}>
      {text}
    </div>
  )
}