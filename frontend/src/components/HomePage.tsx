import { useUser } from '../context/UserContext';
import styles from '../css/HomePage.module.css';
import type { Article } from '../types/types';
import Header from './Header'
import NewsCard from './NewsCard'
import { useEffect, useState} from 'react';
import { API_BASE_URL } from '../config';

// --- MAIN PAGE COMPONENT ---
export default function HomePage() {

  //get articles from backend and display them in news cards
  const [articles, setArticles] = useState<Article[]>([]);
  const { token, loading } = useUser(); 
  
  const url = `${API_BASE_URL}/v1/articles`;
  
  // Fetch articles from backend on component mount
  useEffect(() => {
    if (loading) {  // Wait until loading is complete before fetching articles
      console.log("User loading state is true, waiting to fetch articles...");
      return;
    }
    const headers: HeadersInit = {
      "Content-Type": "application/json",
    };
    if (token) {
      headers["Authorization"] = `Bearer ${token}`;
    }

    fetch(url, {
      method: "GET",
      headers: headers,
      //body is not needed for GET request
    })
    .then(res => res.json())
    .then(data => {
      setArticles(data);
    })
    .catch(err => console.error("Failed to fetch articles:", err));
  }, [token, loading]); // Re-run effect if token changes (e.g., user logs in/out) or loading state changes

  return (
    <div className={styles.page}>
      <Header />
      <main className={styles.mainContainer}>
        <div className={styles.feedHeader}>
          <h1 className={styles.feedTitle}>Top Stories</h1>
          <p className={styles.feedSubtitle}>Curated news, summarized by AI.</p>
        </div>

        <div className={styles.newsFeed}>
          {articles.map(article => (
            <NewsCard key={article.id} article={article} />
          ))}
        </div>
      </main>
    </div>
  );
}

