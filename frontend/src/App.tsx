import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { useState } from 'react'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
    },
  },
})

function App() {
  const [isLoading, setIsLoading] = useState(false)

  return (
    <QueryClientProvider client={queryClient}>
      <div className="app">
        <header className="header">
          <h1>Pinecone</h1>
          <p>Recipe Management & Meal Planning</p>
        </header>

        <main className="main">
          <div className="container">
            <h2>Welcome to Pinecone</h2>
            <p>Your household recipe management system is being built...</p>

            {isLoading ? (
              <p>Loading...</p>
            ) : (
              <div>
                <p>Status: Development in progress</p>
                <p>Epic 1: Foundation & Infrastructure ✓</p>
              </div>
            )}
          </div>
        </main>
      </div>
    </QueryClientProvider>
  )
}

export default App
