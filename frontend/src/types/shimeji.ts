export type ShimejiAnimation = 'idle' | 'pout' | 'shocked' | 'happy'

export interface OllamaResponse {
  subtitle_en: string
  animation: ShimejiAnimation
}