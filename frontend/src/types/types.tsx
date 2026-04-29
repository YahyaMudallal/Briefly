export interface Comment {
  id: string;
  authorId: string;
  articleId: string;
  content: string;
  parentId?: string;
  createdAt: Date;
  updatedAt: Date;
  authorName: string;
}

export interface Article {
  id: string;
  title: string;
  description: string;
  summary: string; // The AI summary
  source: string;
  content: string;
  upvotes: number;
  downvotes: number;
  userVote: number; // 1 (up), -1 (down), or 0/undefined (none)
  nbComments: number;
  createdAt: Date;
  updatedAt: Date;
}

export interface User {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  initials: string;
  isAdmin: boolean;
}

