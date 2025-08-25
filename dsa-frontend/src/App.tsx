import { BrowserRouter, Navigate, Route, Routes } from 'react-router'
import LoginPage from './pages/LoginPage'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import ProtectedRoute from './ProtectedRoute';
import DashBoardPage from './pages/DashBoardPage';
import ProblemStatementPage from './pages/ProblemStatementPage';

const queryClient = new QueryClient();

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route path="/dashboard" element={
            <ProtectedRoute>
              <DashBoardPage />
            </ProtectedRoute>
          } />
          <Route path="/problem/:lectureid/:problemid" element={
            <ProtectedRoute>
              <ProblemStatementPage />
            </ProtectedRoute>
          } />
          <Route path="*" element={<Navigate to="/login" replace />} />
        </Routes>
      </BrowserRouter>
    </QueryClientProvider>
  )
}

export default App
