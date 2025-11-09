import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { RouterProvider, createRouter, createRootRoute, createRoute, Link, Outlet, useNavigate } from '@tanstack/react-router'
import { RegisterForm } from './components/auth/RegisterForm'
import { LoginForm } from './components/auth/LoginForm'
import { ProtectedRoute } from './components/auth/ProtectedRoute'
import { useAuth } from './hooks/useAuth'
import { RecipeListPage } from './pages/RecipeListPage'
import { RecipeDetailPage } from './pages/RecipeDetailPage'
import { RecipeFormPage } from './pages/RecipeFormPage'
import { MealPlanPage } from './pages/MealPlanPage'
import { GroceryListPage } from './pages/GroceryListPage'
import { MenuPage } from './pages/MenuPage'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
    },
  },
})

// Root route
const rootRoute = createRootRoute({
  component: RootComponent,
})

function RootComponent() {
  return (
    <QueryClientProvider client={queryClient}>
      <div className="app">
        <Outlet />
      </div>
    </QueryClientProvider>
  )
}

// Home route (protected)
const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/',
  component: HomePage,
})

function HomePage() {
  const { user, logout, isLoading } = useAuth()
  const navigate = useNavigate()

  return (
    <ProtectedRoute>
      <header className="header">
        <div style={{ maxWidth: '1200px', margin: '0 auto', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div>
            <h1>Pinecone</h1>
            <p>Recipe Management & Meal Planning</p>
          </div>
          {user && (
            <div>
              <span style={{ marginRight: '1rem', color: 'white' }}>Welcome, {user.name}!</span>
              <button
                onClick={() => logout()}
                disabled={isLoading}
                style={{
                  padding: '0.5rem 1rem',
                  backgroundColor: 'white',
                  color: 'var(--color-primary)',
                  border: 'none',
                  borderRadius: 'var(--radius-base)',
                  cursor: 'pointer',
                  fontWeight: 600
                }}
              >
                {isLoading ? 'Logging out...' : 'Logout'}
              </button>
            </div>
          )}
        </div>
      </header>

      <main className="main">
        <div className="container">
          <h2>Welcome to Pinecone</h2>
          <p>Your household recipe management system</p>

          <div style={{ marginTop: '2rem' }}>
            <button
              onClick={() => navigate({ to: '/recipes' })}
              style={{
                padding: '1rem 2rem',
                backgroundColor: 'var(--color-primary)',
                color: 'white',
                border: 'none',
                borderRadius: 'var(--radius-base)',
                cursor: 'pointer',
                fontSize: '1.125rem',
                fontWeight: 600,
              }}
            >
              View Recipes
            </button>
          </div>

          <div style={{ marginTop: '2rem' }}>
            <h3>Progress Status</h3>
            <ul>
              <li>✅ Epic 1: Foundation & Infrastructure</li>
              <li>✅ Epic 2: User Authentication</li>
              <li>🚧 Epic 3: Recipe Management (In Progress)</li>
            </ul>
          </div>
        </div>
      </main>
    </ProtectedRoute>
  )
}

// Register route
const registerRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/register',
  component: RegisterPage,
})

function RegisterPage() {
  return (
    <div style={{ minHeight: '100vh', backgroundColor: 'var(--color-background)' }}>
      <div style={{ padding: '2rem 1rem', textAlign: 'center' }}>
        <Link to="/">
          <h1 style={{ color: 'var(--color-primary)' }}>Pinecone</h1>
        </Link>
      </div>
      <RegisterForm />
    </div>
  )
}

// Login route
const loginRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/login',
  component: LoginPage,
})

function LoginPage() {
  return (
    <div style={{ minHeight: '100vh', backgroundColor: 'var(--color-background)' }}>
      <div style={{ padding: '2rem 1rem', textAlign: 'center' }}>
        <Link to="/">
          <h1 style={{ color: 'var(--color-primary)' }}>Pinecone</h1>
        </Link>
      </div>
      <LoginForm />
    </div>
  )
}

// Recipe routes
const recipesRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/recipes',
  component: RecipesPage,
})

function RecipesPage() {
  return (
    <ProtectedRoute>
      <RecipeListPage />
    </ProtectedRoute>
  )
}

const recipeNewRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/recipes/new',
  component: RecipeNewPage,
})

function RecipeNewPage() {
  return (
    <ProtectedRoute>
      <RecipeFormPage />
    </ProtectedRoute>
  )
}

const recipeDetailRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/recipes/$id',
  component: RecipeDetailPageWrapper,
})

function RecipeDetailPageWrapper() {
  return (
    <ProtectedRoute>
      <RecipeDetailPage />
    </ProtectedRoute>
  )
}

const recipeEditRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/recipes/$id/edit',
  component: RecipeEditPage,
})

function RecipeEditPage() {
  return (
    <ProtectedRoute>
      <RecipeFormPage />
    </ProtectedRoute>
  )
}

// Meal plan route
const mealPlanRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/mealplan',
  component: MealPlanPageWrapper,
})

function MealPlanPageWrapper() {
  return (
    <ProtectedRoute>
      <MealPlanPage />
    </ProtectedRoute>
  )
}

// Grocery list route
const groceryListRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/grocery',
  component: GroceryListPageWrapper,
})

function GroceryListPageWrapper() {
  return (
    <ProtectedRoute>
      <GroceryListPage />
    </ProtectedRoute>
  )
}

// Menu recommendation route
const menuRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/menu',
  component: MenuPageWrapper,
})

function MenuPageWrapper() {
  return (
    <ProtectedRoute>
      <MenuPage />
    </ProtectedRoute>
  )
}

// Create router
const routeTree = rootRoute.addChildren([
  indexRoute,
  registerRoute,
  loginRoute,
  recipesRoute,
  recipeNewRoute,
  recipeDetailRoute,
  recipeEditRoute,
  mealPlanRoute,
  groceryListRoute,
  menuRoute,
])

const router = createRouter({ routeTree })

// Declare router type for TypeScript
declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router
  }
}

function App() {
  return <RouterProvider router={router} />
}

export default App
