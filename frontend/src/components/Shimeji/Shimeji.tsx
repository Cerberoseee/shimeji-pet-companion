import { ShimejiAnimation } from '../../types/shimeji'
import watermelonImg from '../../assets/pet.gif'
import styles from './Shimeji.module.scss'

interface ShimejiSpriteProps {
  animation: ShimejiAnimation
  loading: boolean
}

export function ShimejiSprite({ animation, loading }: ShimejiSpriteProps) {
  const getTransform = () => {
    if (loading) return 'scale(0.95)'
    switch (animation) {
      case 'shocked':
        return 'translateY(-6px) scale(1.08)'
      case 'pout':
        return 'rotate(-4deg)'
      case 'happy':
        return 'scale(1.05)'
      default:
        return 'none'
    }
  }

  return (
    <img
      src={watermelonImg}
      alt="Desktop Shimeji"
      draggable={false}
      className={styles.shimeji}
      style={{ transform: getTransform(), filter: loading ? 'brightness(0.9)' : 'none' }}
    />
  )
}