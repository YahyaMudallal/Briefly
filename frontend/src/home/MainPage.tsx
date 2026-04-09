import styles from './HomePage.module.css';

import type { Article, User, Comment } from './types';
import Header from './Header'
import NewsCard from './NewsCard'

// --- MOCK DATA ---
const MOCK_USER: User | null = null;//{ username: "jane_doe", initials: "JD" }; // Set to null to see logged-out state
const MOCK_ARTICLES: Article[] = [
  {
    id: "1",
    title: "Global Tech Summit Announces Groundbreaking AI Regulations",
    description: "Leaders from top tech firms and government agencies have agreed on a new framework for artificial intelligence development, focusing on safety and transparency...",
    tldr: "Tech giants and governments agreed on new AI safety and transparency rules to prevent misuse.",
    source: "TechInsider",
    upvotes: 142,
    downvotes: 12,
    comments: [
      { id: "c1", author: "code_ninja", text: "Finally some clear guidelines!", timestamp: "2h ago" }
    ]
  },
  {
    id: "2",
    title: "Breakthrough in Solid-State Battery Technology",
    description: "Researchers at MIT have developed a new solid-state battery architecture that promises to double the range of electric vehicles while eliminating fire risks...",
    tldr: "MIT researchers created a safer solid-state battery that could double EV range.",
    source: "ScienceDaily",
    upvotes: 89,
    downvotes: 3,
    comments: []
  }
];


// --- MAIN PAGE COMPONENT ---
export default function MainPage({onNavigate}:{onNavigate: (v: "authPage" | "homePage") => void }) {
  return (
    <div className={styles.page}>
      <Header user={MOCK_USER} onNavigate={onNavigate} />

      <main className={styles.mainContainer}>
        <div className={styles.feedHeader}>
          <h1 className={styles.feedTitle}>Top Stories</h1>
          <p className={styles.feedSubtitle}>Curated news, summarized by AI.</p>
        </div>

        <div className={styles.newsFeed}>
          {MOCK_ARTICLES.map(article => (
            <NewsCard key={article.id} article={article} user={MOCK_USER} />
          ))}
        </div>
      </main>
    </div>
  );
}

