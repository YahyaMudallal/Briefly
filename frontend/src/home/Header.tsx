import styles from '../home/HomePage.module.css';
import logoDarkFull from '../assets/logoDarkFull.png';
import type {User} from './types'

// --- HEADER COMPONENT ---
export default function Header({ user, onNavigate }: { user: User | null; onNavigate: (v: "authPage" | "homePage") => void }) {
  return (
    <header className={styles.header}>
      <div className={styles.headerContent}>
        <div className={styles.logoArea}>        
            <img className={styles.logoIcon} src={logoDarkFull} />
        </div>

        <div className={styles.authArea}>
          {user ? (
            <button className={styles.avatarBtn} title="Go to Profile">
              {user.initials}
            </button>
          ) : (
            <button className={styles.signInBtn} onClick={() => onNavigate("authPage")}>Sign In</button>
          )}
        </div>
      </div>
    </header>
  );
}