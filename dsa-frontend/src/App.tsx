import { BrowserRouter, Navigate, Route, Routes } from 'react-router'
import LoginPage from './pages/LoginPage'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import ProtectedRoute from './ProtectedRoute';
import DashBoardPage from './pages/DashBoardPage';
import ProblemStatementPage from './pages/ProblemStatementPage';
import NotFoundPage from './pages/NotFound';
import ValidationResultsListing from './pages/ValidationResultsListing';
import ValidationDetail from './pages/ValidationDetail';

const queryClient = new QueryClient();

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route path="/dashboard" element={<ProtectedRoute><DashBoardPage /></ProtectedRoute>} />
          <Route path="/problem/:lectureid/:problemid" element={<ProtectedRoute><ProblemStatementPage /></ProtectedRoute>} />
          <Route path="/validation/results" element={<ProtectedRoute><ValidationResultsListing /></ProtectedRoute>} />
          <Route path="/validation/detail/:idParam" element={<ProtectedRoute><ValidationDetail /></ProtectedRoute>} />
          <Route path="/" element={<Navigate to="/dashboard" replace />} />
          <Route path="*" element={<NotFoundPage />} />
        </Routes>
      </BrowserRouter>
    </QueryClientProvider>
  )
}

export default App
