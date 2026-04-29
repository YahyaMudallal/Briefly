import { useEffect, useState } from "react";
import type { Article, Comment } from "../types/types";
import styles from '../css/HomePage.module.css';
import { Link, useNavigate } from "react-router";
import { useUser } from "../context/UserContext";
import FormatTimeAgo from "./FormatTimeAgo";


// --- INDIVIDUAL NEWS CARD COMPONENT ---
export default function NewsCard({ article }: { article: Article}) {
  const { user, token } = useUser();
  const navigate = useNavigate();
  const [showTldr, setShowTldr] = useState(false);
  const [showComments, setShowComments] = useState(false);
  const [newComment, setNewComment] = useState("");
  const [comments, setComments] = useState<Comment[]>([]);
  const [isUpvoted, setIsUpvoted] = useState(article.userVote === 1);
  const [isDownvoted, setIsDownvoted] = useState(article.userVote === -1);
  const [nbUpvotes, setNbUpvotes] = useState(article.upvotes);
  const [nbDownvotes, setNbDownvotes] = useState(article.downvotes);

  useEffect(() => {
    setIsUpvoted(article.userVote === 1);
    setIsDownvoted(article.userVote === -1);
    setNbUpvotes(article.upvotes);
    setNbDownvotes(article.downvotes);
  }, [article]);

  const handleCommentSubmit = (e: React.SubmitEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!newComment.trim()) return;
    console.log("Submitting comment:", newComment);
    setNewComment("");
    const url = "http://localhost:8080/v1/comments";
    const body = {
      articleId: article.id,
      authorId: user?.id,
      authorName: `${user?.firstName} ${user?.lastName}`,
      content: newComment
    };
    fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${token}`
      },
      body: JSON.stringify(body)
    })
    .then(res => {
      if (!res.ok) {
        throw new Error("Failed to post comment");
      }
      return res.json();
    })
    .then(data => {
      console.log("Comment posted successfully:", data);
      // Optionally, we can update the local state to show the new comment immediately
      setComments(prev => [...prev, data]);
      // Or we can refetch the comments to get the updated list from the backend
      // fetchComments();
    })
    .catch(err => {
      console.error("Error posting comment:", err);
      alert("Failed to post comment. Please try again.");
    }); 

  };


  //fetch comments for this article from backend and display them in the comments section when showComments is trues
  const fetchComments = () => {
    const url = `http://localhost:8080/v1/comments/article/${article.id}`;
    fetch(url, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${token}`
      }
    })
    .then(res => {
      if (!res.ok) {
        throw new Error("Failed to fetch comments");
      } 
      return res.json();
    })
    .then(data => {
      setComments(data);  
    })
    .catch(err => {
      console.error("Error fetching comments:", err);
      alert("Failed to load comments. Please try again.");
    });
  };

  const toggleUpVote = () => {
    const url = `http://localhost:8080/v1/articles/${article.id}/upvote`;
    fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${token}`
      }
    })
    .then(res => {
      if (!res.ok) {
        throw new Error("Failed to toggle upvote");
      }
      return res.json();
    })
    .then(data => {
      console.log("Upvote toggled successfully:", data);
      setNbUpvotes(data.upvotes);
    })
    .catch(err => {
      console.error("Error toggling upvote:", err);
      alert("Failed to toggle upvote. Please try again.");
    });
  };

  const toggleDownVote = () => {
    const url = `http://localhost:8080/v1/articles/${article.id}/downvote`;
    fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${token}`
      }
    })
    .then(res => {
      if (!res.ok) {
        throw new Error("Failed to toggle downvote");
      }
      return res.json();
    })
    .then(data => {
      console.log("Downvote toggled successfully:", data);
      // Optionally, update the local state to reflect the new vote count
      setNbDownvotes(data.downvotes);
    })
    .catch(err => {
      console.error("Error toggling downvote:", err);
      alert("Failed to toggle downvote. Please try again.");
    });
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
          <button 
            className={`${styles.actionBtn} ${isUpvoted ? styles.activeVoteBtn : ''}`} 
            onClick={()=>{
              if (!token) {
                navigate("/auth");
                return;
              }
              if (isDownvoted){
                setIsDownvoted(false);
                setIsUpvoted(true);
                setNbUpvotes(prev => prev + 1);
                setNbDownvotes(prev => prev - 1);
              } else {
                setIsUpvoted(!isUpvoted);
                setNbUpvotes(prev => isUpvoted ? prev - 1 : prev + 1);
              }
              toggleUpVote();
            }}
          >
            ▲ {nbUpvotes}
          </button>

          <button 
            className={`${styles.actionBtn} ${isDownvoted ? styles.activeVoteBtn : ''}`}
            onClick={()=>{
              if (!token) {
                navigate("/auth");
                return;
              }
              if (isUpvoted){
                setIsUpvoted(false);
                setIsDownvoted(true);
                setNbDownvotes(prev => prev + 1);
                setNbUpvotes(prev => prev - 1);
              } else {
                setIsDownvoted(!isDownvoted);
                setNbDownvotes(prev => isDownvoted ? prev - 1 : prev + 1);
              }
              toggleDownVote();
            }}
          >
            ▼ {nbDownvotes}
          </button>
        </div>
        
        <div className={styles.featureGroup}>
          <button 
            className={`${styles.tldrBtn} ${showTldr ? styles.activeTldr : ''}`}
            onClick={() => setShowTldr(!showTldr)}
          >
            ✨ AI TL;DR
          </button>
          <button 
            className={`${styles.actionBtn} ${showComments ? styles.activeCommentBtn : ''}`}
            onClick={() => {
              setShowComments(!showComments); 
              !showComments && fetchComments();}}
          >
            💬 {comments.length===0 ? article.nbComments : comments.length} Comments
          </button>
        </div>
      </div>

      {/* Expandable TL;DR Section */}
      {showTldr && (
        <div className={styles.tldrBox}>
          <strong>Summary:</strong> {article.summary}
        </div>
      )}

      {/* Expandable Comments Section */}
      {showComments && (
        <div className={styles.commentsSection}>
          <div className={styles.commentList}>
            {article.nbComments === 0 ? (
              <p className={styles.noComments}>No comments yet. Be the first!</p>
            ) : (
              comments.map(c => (
                <div key={c.id} className={styles.comment}>
                  <span className={styles.commentAuthor}>{c.authorName}</span>
                  <span className={styles.commentTimestamp}>•</span>
                  <span className={styles.commentTimestamp}>{FormatTimeAgo(c.createdAt)}</span>
                  <p className={styles.commentText}>{c.content}</p>
                </div>
              ))
            )}
          </div>

          {/* Comment Input Form (Only if logged in) */}
          {user ? (
            <form onSubmit={handleCommentSubmit} className={styles.commentForm}>
              <input 
                placeholder="Add a comment..." 
                className={styles.commentInput}
                value={newComment}
                onChange={e => setNewComment(e.target.value)}
              />
              <button type="submit" className={styles.commentSubmitBtn}>Post</button>
            </form>
          ) : (
            <div className={styles.loginPrompt}>
              <Link to="/auth">Sign in</Link> to join the conversation.
            </div>
          )}
        </div>
      )}
    </article>
  );
}

