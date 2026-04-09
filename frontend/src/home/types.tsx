export interface Comment {
  id: string;
  author: string;
  text: string;
  timestamp: string;
}

export interface Article {
  id: string;
  title: string;
  description: string;
  tldr: string; // The AI summary
  source: string;
  upvotes: number;
  downvotes: number;
  comments: Comment[];
}

export interface User {
  username: string;
  initials: string;
}