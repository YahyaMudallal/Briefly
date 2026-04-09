import { useState } from "react";
import type { Article, User } from "./types";
import styles from './HomePage.module.css';


// --- INDIVIDUAL NEWS CARD COMPONENT ---
export default function NewsCard({ article, user }: { article: Article, user: User | null }) {
  const [showTldr, setShowTldr] = useState(false);
  const [showComments, setShowComments] = useState(false);
  const [newComment, setNewComment] = useState("");

  const handleCommentSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!newComment.trim()) return;
    console.log("Submitting comment:", newComment);
    setNewComment("");
    // Here you will eventually make your API call to Go backend
  };

  return (
    <article className={styles.card}>
      <div className={styles.cardHeader}>
        <span className={styles.sourceTag}>{article.source}</span>
      </div>
      
      <h2 className={styles.articleTitle}>{article.title}</h2>
      <p className={styles.articleDesc}>{article.description}</p>

      {/* Action Bar */}
      <div className={styles.actionBar}>
        <div className={styles.voteGroup}>
          <button className={styles.actionBtn}>▲ {article.upvotes}</button>
          <button className={styles.actionBtn}>▼ {article.downvotes}</button>
        </div>
        
        <div className={styles.featureGroup}>
          <button 
            className={`${styles.tldrBtn} ${showTldr ? styles.activeTldr : ''}`}
            onClick={() => setShowTldr(!showTldr)}
          >
            ✨ AI TL;DR
          </button>
          <button 
            className={styles.actionBtn}
            onClick={() => setShowComments(!showComments)}
          >
            💬 {article.comments.length} Comments
          </button>
        </div>
      </div>

      {/* Expandable TL;DR Section */}
      {showTldr && (
        <div className={styles.tldrBox}>
          <strong>Summary:</strong> {article.tldr}
        </div>
      )}

      {/* Expandable Comments Section */}
      {showComments && (
        <div className={styles.commentsSection}>
          <div className={styles.commentList}>
            {article.comments.length === 0 ? (
              <p className={styles.noComments}>No comments yet. Be the first!</p>
            ) : (
              article.comments.map(c => (
                <div key={c.id} className={styles.comment}>
                  <span className={styles.commentAuthor}>{c.author}</span>
                  <span className={styles.commentText}>{c.text}</span>
                </div>
              ))
            )}
          </div>

          {/* Comment Input Form (Only if logged in) */}
          {user ? (
            <form onSubmit={handleCommentSubmit} className={styles.commentForm}>
              <input 
                type="text" 
                placeholder="Add a comment..." 
                className={styles.commentInput}
                value={newComment}
                onChange={e => setNewComment(e.target.value)}
              />
              <button type="submit" className={styles.commentSubmitBtn}>Post</button>
            </form>
          ) : (
            <div className={styles.loginPrompt}>
              <a href="/login">Sign in</a> to join the conversation.
            </div>
          )}
        </div>
      )}
    </article>
  );
}
