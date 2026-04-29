export interface Comment {
  id:           string;
  authorId:     string;
  authorName:   string;
  articleId:    string;
  content:      string;
  parentId?:    string;
  createdAt:    Date;
  updatedAt:    Date;
}

export interface Article {
  id:           string;
  title:        string;
  description:  string;
  content:      string;
  originalUrl:  string;
  ImageUrl:     string;
  source:       string;
  upvotes:      number;
  downvotes:    number;
  nbComments:   number;
  summary:      string; // The AI summary
  userVote:     number; // 1 (up), -1 (down), or 0/undefined (none)
  publishedAt:  Date;
  createdAt:    Date;
  updatedAt:    Date;
}

export interface User {
  id:           string;
  email:        string;
  firstName:    string;
  lastName:     string;
  initials:     string;
  isAdmin:      boolean;
}

