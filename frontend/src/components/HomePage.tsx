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

  // Pagination state
  const [page, setPage] = useState<number>(1);
  const [hasMore, setHasMore] = useState<boolean>(true);
  const [isLoadingMore, setIsLoadingMore] = useState<boolean>(false);

  // Sorting state
  const [sortBy, setSortBy] = useState<string>('date');
  const [order, setOrder] = useState<string>('desc');

  const LIMIT = 5; // Number of articles to fetch per page

  //  reset to page 1 when token changes
  useEffect(() => {
    setPage(1);
    setHasMore(true);
  }, [token]);

  // Fetch articles from backend on component mount
  useEffect(() => {
    if (loading) {  // Wait until loading is complete before fetching articles
      console.log("User loading state is true, waiting to fetch articles...");
      return;
    }

    const fetchArticles = async () => {
      // if we're loading more articles (page > 1), set the loading state to true
      if (page > 1) setIsLoadingMore(true);

      const headers: HeadersInit = {
        "Content-Type": "application/json",
      };
      if (token) {
        headers["Authorization"] = `Bearer ${token}`;
      }

      // Construct the URL with pagination and sorting parameters
      const url = `${API_BASE_URL}/v1/articles?page=${page}&limit=${LIMIT}&sortBy=${sortBy}&order=${order}`;
      try {
        const res = await fetch(url, { method: "GET", headers });
        const data = await res.json();

        // if the api returns 0 or less that LIMIT articles, we know there are no more articles to fetch
        const newArticles = data || [];
        if (newArticles.length < LIMIT) {
          setHasMore(false);
        }

        setArticles(prev => {
          if (page === 1) return newArticles;

          // keep from adding duplicate articles
          const existingIds = new Set(prev.map(a => a.id));
          const filteredNew = newArticles.filter((a: Article) => !existingIds.has(a.id));
          
          return [...prev, ...filteredNew];
        });

      } catch (err) {
        console.error("Failed to fetch articles:", err);
      } finally {
        setIsLoadingMore(false);
      }
    };

    fetchArticles();
  }, [token, loading, page, sortBy, order]); // Re-run effect if token changes (e.g., user logs in/out) or loading state changes or page changes

  // function the pass to the next page
  const loadMore = () => {
    if (!isLoadingMore && hasMore) {
      setPage(prevPage => prevPage + 1);
    }
  };

  return (
    <div className={styles.page}>
      <Header />
      <main className={styles.mainContainer}>
        <div className={styles.feedHeaderWrapper}>
              <div className={styles.feedHeader}>
                <h1 className={styles.feedTitle}>Top Stories</h1>
                <p className={styles.feedSubtitle}>Curated news, summarized by AI.</p>
              </div>
              
              <div className={styles.sortControls}>
                <select 
                  className={styles.sortSelect} 
                  value={sortBy} 
                  onChange={e => setSortBy(e.target.value)}
                >
                  <option value="date">Publication date</option>
                  <option value="hotness">Hotness</option>
                </select>

                <select 
                  className={styles.sortSelect} 
                  value={order} 
                  onChange={e => setOrder(e.target.value)}
                >
                  <option value="desc">Descending (Max → Min)</option>
                  <option value="asc">Ascending (Min → Max)</option>
                </select>
              </div>
            </div>


        <div className={styles.newsFeed}>
          {articles.map(article => (
            <NewsCard key={article.id} article={article} />
          ))}
        </div>

        {/* Load More Button */}
        {hasMore && (
          <div className={styles.loadMoreContainer}>
            <button 
              className={styles.loadMoreBtn} 
              onClick={loadMore} 
              disabled={isLoadingMore}
            >
              {isLoadingMore ? "Loading..." : "Load More"}
            </button>
          </div>
        )}
      </main>
    </div>
  );
}

