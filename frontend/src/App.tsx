
import { useState } from 'react';
import AuthPage from './auth/AuthPage'
import MainPage from './home/MainPage';

type CurrentPage = "authPage" | "homePage";

function App() {
  const [currentPage, setCurrentPage] = useState<CurrentPage>("homePage");
  return currentPage==="authPage" ? 
    <AuthPage onNavigate={setCurrentPage}/> : <MainPage onNavigate={setCurrentPage}/>
}

export default App
