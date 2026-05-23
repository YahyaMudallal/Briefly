# Briefly

A full-stack web application whose goal is to build a smart news community where users can access the latest headlines, read AI-generated TL;DR summaries for quick insights, engage in discussions, and cast upvotes or downvotes.

The project is currently deployed. The client is accessible [here](https://briefly-frontend-5rtf.onrender.com/) and the REST API server can be reached at the following link: https://briefly-8t35.onrender.com

This README provides a general overview of the project. For further technical details, please refer to the [frontend README](./frontend/README.md) and the [backend README](./backend/README.md).

## Authors

Abdullah AL MAMUN \
Yahya MUDALLAL

## Technical Stack

**Frontend**: Built with TypeScript using React.js and Vite, hosted as a static website on [Render](https://render.com/).

**Backend**: Built with Go (Golang) using the native `net/http` package without any high-level framework. It is hosted as a web service on [Render](https://render.com/). The backend interacts with external APIs, including [NewsData.io](https://newsdata.io/) to fetch news articles and Google AI Studio to leverage the Gemini API for article summarization.

**Database**: Utilizes MongoDB for data persistence, hosted on MongoDB Atlas.

## Description of the Project

**Briefly** is a news aggregator that allows users to scroll through a feed of news articles. By creating an account, users can read AI-generated summaries, cast upvotes or downvotes, and participate in discussions through nested comments. New articles are automatically fetched and added every day at 2:00 AM via a scheduled cron job.