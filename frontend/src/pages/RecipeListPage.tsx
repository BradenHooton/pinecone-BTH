import { useState } from 'react'
import { useNavigate } from '@tanstack/react-router'
import { useRecipes } from '../hooks/useRecipes'
import { RecipeCard } from '../components/recipe/RecipeCard'

export function RecipeListPage() {
  const navigate = useNavigate()
  const [search, setSearch] = useState('')
  const { data, isLoading, error } = useRecipes({ search: search || undefined })

  if (isLoading) {
    return <div style={{ padding: '20px' }}>Loading recipes...</div>
  }

  if (error) {
    return (
      <div style={{ padding: '20px', color: 'red' }}>
        Error loading recipes: {error.message}
      </div>
    )
  }

  return (
    <div style={{ padding: '20px' }}>
      <div
        style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          marginBottom: '24px',
        }}
      >
        <h1 style={{ margin: 0 }}>Recipes</h1>
        <button
          onClick={() => navigate({ to: '/recipes/new' })}
          style={{
            padding: '10px 20px',
            backgroundColor: '#3A7D44',
            color: 'white',
            border: 'none',
            borderRadius: '4px',
            cursor: 'pointer',
            fontSize: '1rem',
          }}
        >
          Create Recipe
        </button>
      </div>

      <div style={{ marginBottom: '20px' }}>
        <input
          type="text"
          placeholder="Search recipes..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          style={{
            width: '100%',
            padding: '10px',
            fontSize: '1rem',
            border: '1px solid #ddd',
            borderRadius: '4px',
          }}
        />
      </div>

      {data && data.data.length === 0 ? (
        <div style={{ textAlign: 'center', padding: '40px', color: '#666' }}>
          <p>No recipes found. Create your first recipe!</p>
        </div>
      ) : (
        <div
          style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(auto-fill, minmax(300px, 1fr))',
            gap: '20px',
          }}
        >
          {data?.data.map((recipe) => (
            <RecipeCard
              key={recipe.id}
              recipe={recipe}
              onClick={() => navigate({ to: `/recipes/${recipe.id}` })}
            />
          ))}
        </div>
      )}

      {data && data.meta.total > data.meta.limit && (
        <div style={{ marginTop: '20px', textAlign: 'center', color: '#666' }}>
          Showing {data.data.length} of {data.meta.total} recipes
        </div>
      )}
    </div>
  )
}
