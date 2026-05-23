# Briefly: Frontend

This folder contains the frontend part of the project, which is built with TypeScript using React and Vite. The client uses React Router to assign the route "/" for the HomePage, "/auth" for the authentication page, and "/profile" for the profile page.

## Running the frontend

To run the client, first install the node dependencies with the command:
```bash
npm install
```

Then, start the frontend for the development workflow with:
```bash
npm run dev
```

## Design of the frontend

The folder structure of the React project is the following:
```plain text
├── README.md
├── eslint.config.js
├── index.html
├── package-lock.json
├── package.json
├── public
│   ├── _redirects
│   ├── favicon.svg
│   └── icons.svg
├── src
│   ├── App.tsx
│   ├── assets
│   │   ├── logoBrightFull.png
│   │   ├── logoBrightHalf.png
│   │   ├── logoDarkFull.png
│   │   └── logoDarkHalf.png
│   ├── components
│   │   ├── AuthPage.tsx
│   │   ├── FormatTimeAgo.tsx
│   │   ├── Header.tsx
│   │   ├── HomePage.tsx
│   │   ├── NewsCard.tsx
│   │   ├── ProfilePage.tsx
│   │   └── ProtectedRoute.tsx
│   ├── config.ts
│   ├── context
│   │   └── UserContext.tsx
│   ├── css
│   │   ├── AuthPage.module.css
│   │   └── HomePage.module.css
│   ├── main.tsx
│   └── types
│       └── types.tsx
├── tsconfig.app.json
├── tsconfig.json
├── tsconfig.node.json
└── vite.config.ts
```