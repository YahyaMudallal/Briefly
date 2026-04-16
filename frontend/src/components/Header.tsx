import styles from '../css/HomePage.module.css';
import logoDarkFull from '../assets/logoDarkFull.png';
import { useUser } from '../context/UserContext';
import { Link, useNavigate } from 'react-router-dom';

  
// --- HEADER COMPONENT ---
export default function Header() {
  const { user, logout } = useUser();
  const navigate = useNavigate();
  return (
    <header className={styles.header}>
      <div className={styles.headerContent}>
        <div className={styles.logoArea}>
          <Link to="/">
            <img className={styles.logoIcon} src={logoDarkFull} />
          </Link>       
        </div>

        <div className={styles.authArea}>
          {user ? (
            <>
              <button className={styles.avatarBtn} onClick={() => navigate("/profile")} title="Go to Profile">
                {user.initials}
              </button>
              <button className={styles.signOutBtn} onClick={()=> logout()}>
                Sign Out
              </button>
            </>
          ) : (
            <button className={styles.signInBtn} onClick={() => navigate("auth")}>Sign In</button>
          )}
        </div>
      </div>
    </header>
  );
}