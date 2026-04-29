import { useUser } from "../context/UserContext";
import Header from "./Header";
import styles from "../css/HomePage.module.css";

export default function ProfilePage() {
  const { user } = useUser();
  return (
    <div className={styles.page}>
      <Header />
      <div style={{ padding: '20px', maxWidth: '900px', margin: '0 auto' }}>
        <h1>Profile Page</h1>
        {user && (
          <div>
            <p>Name: {user.firstName} {user.lastName}</p>
            <p>Email: {user.email}</p>
            <p>Bookmarked Articles : </p>
          </div>
        )}
      </div>
    </div>
  );
}