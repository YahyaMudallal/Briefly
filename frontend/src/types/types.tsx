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
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  initials: string;
  isAdmin: boolean;
}

