import { useState } from 'react'
import { useParams, useNavigate } from '@tanstack/react-router'
import {
  useCookbook,
  useAddRecipeToCookbook,
  useRemoveRecipeFromCookbook,
  useUpdateCookbook,
} from '../hooks/useCookbooks'
import { useRecipes } from '../hooks/useRecipes'
import { Recipe, UpdateCookbookRequest } from '../lib/api'

export function CookbookDetailPage() {
  const { id } = useParams({ from: '/cookbooks/$id' })
  const navigate = useNavigate()

  const [showAddRecipeModal, setShowAddRecipeModal] = useState(false)
  const [showEditModal, setShowEditModal] = useState(false)

  const { data, isLoading, error } = useCookbook(id)
  const addRecipe = useAddRecipeToCookbook()
  const removeRecipe = useRemoveRecipeFromCookbook()
  const updateCookbook = useUpdateCookbook()

  const handleAddRecipe = (recipeId: string) => {
    addRecipe.mutate(
      { cookbookId: id, recipeId },
      {
        onSuccess: () => {
          setShowAddRecipeModal(false)
        },
      }
    )
  }

  const handleRemoveRecipe = (recipeId: string) => {
    if (confirm('Remove this recipe from the cookbook?')) {
      removeRecipe.mutate({ cookbookId: id, recipeId })
    }
  }

  const handleUpdateCookbook = (data: UpdateCookbookRequest) => {
    updateCookbook.mutate(
      { id, data },
      {
        onSuccess: () => {
          setShowEditModal(false)
        },
      }
    )
  }

  const handleRecipeClick = (recipeId: string) => {
    navigate({ to: `/recipes/${recipeId}` })
  }

  if (isLoading) {
    return <div style={{ padding: '20px' }}>Loading cookbook...</div>
  }

  if (error || !data) {
    return (
      <div style={{ padding: '20px', color: 'red' }}>
        Error loading cookbook: {error?.message || 'Not found'}
      </div>
    )
  }

  const cookbook = data.data

  return (
    <div style={{ padding: '20px', maxWidth: '1200px', margin: '0 auto' }}>
      <div style={{ marginBottom: '32px' }}>
        <div
          style={{
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'flex-start',
            marginBottom: '16px',
          }}
        >
          <div style={{ flex: 1 }}>
            <h1
              style={{
                margin: '0 0 8px 0',
                fontFamily: '"Playfair Display", serif',
                fontSize: '2rem',
                color: '#2a5d34',
              }}
            >
              {cookbook.name}
            </h1>
            {cookbook.description && (
              <p style={{ margin: '0 0 12px 0', color: '#666', fontSize: '1.125rem' }}>
                {cookbook.description}
              </p>
            )}
            <div style={{ color: '#3A7D44', fontSize: '0.9375rem', fontWeight: 500 }}>
              {cookbook.recipe_count} recipe{cookbook.recipe_count !== 1 ? 's' : ''}
            </div>
          </div>

          <div style={{ display: 'flex', gap: '8px' }}>
            <button
              onClick={() => setShowEditModal(true)}
              style={{
                padding: '10px 20px',
                backgroundColor: '#f0f0f0',
                color: '#333',
                border: 'none',
                borderRadius: '6px',
                cursor: 'pointer',
                fontSize: '0.9375rem',
              }}
            >
              Edit
            </button>
            <button
              onClick={() => setShowAddRecipeModal(true)}
              style={{
                padding: '10px 20px',
                backgroundColor: '#3A7D44',
                color: 'white',
                border: 'none',
                borderRadius: '6px',
                cursor: 'pointer',
                fontSize: '0.9375rem',
                fontWeight: 600,
              }}
            >
              + Add Recipe
            </button>
            <button
              onClick={() => navigate({ to: '/cookbooks' })}
              style={{
                padding: '10px 20px',
                backgroundColor: '#f0f0f0',
                color: '#333',
                border: 'none',
                borderRadius: '6px',
                cursor: 'pointer',
                fontSize: '0.9375rem',
              }}
            >
              Back
            </button>
          </div>
        </div>
      </div>

      {cookbook.recipes && cookbook.recipes.length === 0 && (
        <div
          style={{
            padding: '60px 20px',
            textAlign: 'center',
            backgroundColor: '#f8f9fa',
            borderRadius: '8px',
          }}
        >
          <p style={{ fontSize: '1.125rem', color: '#666', margin: 0 }}>
            No recipes in this cookbook yet. Add your first recipe to get started!
          </p>
        </div>
      )}

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(280px, 1fr))', gap: '20px' }}>
        {cookbook.recipes && cookbook.recipes.map((recipe: Recipe) => (
          <div
            key={recipe.id}
            style={{
              border: '1px solid #ddd',
              borderRadius: '8px',
              padding: '16px',
              backgroundColor: 'white',
              cursor: 'pointer',
              transition: 'transform 0.2s, box-shadow 0.2s',
            }}
            onClick={() => handleRecipeClick(recipe.id)}
            onMouseEnter={(e) => {
              e.currentTarget.style.transform = 'translateY(-2px)'
              e.currentTarget.style.boxShadow = '0 4px 8px rgba(0, 0, 0, 0.1)'
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.transform = 'translateY(0)'
              e.currentTarget.style.boxShadow = 'none'
            }}
          >
            {recipe.image_url && (
              <img
                src={recipe.image_url}
                alt={recipe.title}
                style={{
                  width: '100%',
                  height: '180px',
                  objectFit: 'cover',
                  borderRadius: '6px',
                  marginBottom: '12px',
                }}
              />
            )}
            <h3
              style={{
                margin: '0 0 8px 0',
                fontSize: '1.125rem',
                color: '#2a5d34',
                fontFamily: '"Playfair Display", serif',
              }}
            >
              {recipe.title}
            </h3>
            <div style={{ fontSize: '0.875rem', color: '#666', marginBottom: '12px' }}>
              {recipe.prep_time_minutes && recipe.cook_time_minutes && (
                <span>
                  {recipe.prep_time_minutes + recipe.cook_time_minutes} min · {recipe.servings} servings
                </span>
              )}
            </div>
            <button
              onClick={(e) => {
                e.stopPropagation()
                handleRemoveRecipe(recipe.id)
              }}
              style={{
                width: '100%',
                padding: '8px',
                backgroundColor: '#fee',
                color: '#c33',
                border: '1px solid #fcc',
                borderRadius: '4px',
                cursor: 'pointer',
                fontSize: '0.875rem',
              }}
            >
              Remove from Cookbook
            </button>
          </div>
        ))}
      </div>

      {showAddRecipeModal && (
        <AddRecipeModal
          onClose={() => setShowAddRecipeModal(false)}
          onAdd={handleAddRecipe}
          isLoading={addRecipe.isPending}
        />
      )}

      {showEditModal && (
        <EditCookbookModal
          cookbook={cookbook}
          onClose={() => setShowEditModal(false)}
          onSubmit={handleUpdateCookbook}
          isLoading={updateCookbook.isPending}
        />
      )}
    </div>
  )
}

interface AddRecipeModalProps {
  onClose: () => void
  onAdd: (recipeId: string) => void
  isLoading: boolean
}

function AddRecipeModal({ onClose, onAdd, isLoading }: AddRecipeModalProps) {
  const [search, setSearch] = useState('')
  const { data: recipesData } = useRecipes({ search, limit: 50 })

  return (
    <div
      style={{
        position: 'fixed',
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,
        backgroundColor: 'rgba(0, 0, 0, 0.5)',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        zIndex: 1000,
      }}
      onClick={onClose}
    >
      <div
        style={{
          backgroundColor: 'white',
          borderRadius: '8px',
          padding: '24px',
          maxWidth: '600px',
          width: '100%',
          maxHeight: '80vh',
          overflow: 'auto',
          margin: '20px',
        }}
        onClick={(e) => e.stopPropagation()}
      >
        <h2 style={{ margin: '0 0 20px 0' }}>Add Recipe to Cookbook</h2>

        <input
          type="text"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          placeholder="Search recipes..."
          style={{
            width: '100%',
            padding: '10px',
            borderRadius: '4px',
            border: '1px solid #ddd',
            marginBottom: '16px',
          }}
        />

        <div style={{ display: 'flex', flexDirection: 'column', gap: '8px', marginBottom: '16px' }}>
          {recipesData?.data && recipesData.data.map((recipe: Recipe) => (
            <div
              key={recipe.id}
              style={{
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
                padding: '12px',
                border: '1px solid #eee',
                borderRadius: '6px',
              }}
            >
              <div>
                <div style={{ fontWeight: 500 }}>{recipe.title}</div>
                <div style={{ fontSize: '0.875rem', color: '#666' }}>
                  {recipe.servings} servings
                </div>
              </div>
              <button
                onClick={() => onAdd(recipe.id)}
                disabled={isLoading}
                style={{
                  padding: '6px 12px',
                  backgroundColor: '#3A7D44',
                  color: 'white',
                  border: 'none',
                  borderRadius: '4px',
                  cursor: 'pointer',
                  fontSize: '0.875rem',
                }}
              >
                {isLoading ? 'Adding...' : 'Add'}
              </button>
            </div>
          ))}

          {recipesData?.data && recipesData.data.length === 0 && (
            <div style={{ textAlign: 'center', padding: '20px', color: '#666' }}>
              No recipes found. Try a different search term.
            </div>
          )}
        </div>

        <div style={{ display: 'flex', justifyContent: 'flex-end' }}>
          <button
            onClick={onClose}
            style={{
              padding: '8px 16px',
              backgroundColor: '#f0f0f0',
              border: 'none',
              borderRadius: '4px',
              cursor: 'pointer',
            }}
          >
            Close
          </button>
        </div>
      </div>
    </div>
  )
}

interface EditCookbookModalProps {
  cookbook: any
  onClose: () => void
  onSubmit: (data: UpdateCookbookRequest) => void
  isLoading: boolean
}

function EditCookbookModal({ cookbook, onClose, onSubmit, isLoading }: EditCookbookModalProps) {
  const [name, setName] = useState(cookbook.name)
  const [description, setDescription] = useState(cookbook.description || '')

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()

    onSubmit({
      name,
      description: description || undefined,
    })
  }

  return (
    <div
      style={{
        position: 'fixed',
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,
        backgroundColor: 'rgba(0, 0, 0, 0.5)',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        zIndex: 1000,
      }}
      onClick={onClose}
    >
      <div
        style={{
          backgroundColor: 'white',
          borderRadius: '8px',
          padding: '24px',
          maxWidth: '500px',
          width: '100%',
          margin: '20px',
        }}
        onClick={(e) => e.stopPropagation()}
      >
        <h2 style={{ margin: '0 0 20px 0' }}>Edit Cookbook</h2>

        <form onSubmit={handleSubmit}>
          <div style={{ marginBottom: '16px' }}>
            <label style={{ display: 'block', marginBottom: '4px', fontWeight: 500 }}>
              Name *
            </label>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
              maxLength={200}
              style={{
                width: '100%',
                padding: '8px',
                borderRadius: '4px',
                border: '1px solid #ddd',
              }}
            />
          </div>

          <div style={{ marginBottom: '24px' }}>
            <label style={{ display: 'block', marginBottom: '4px', fontWeight: 500 }}>
              Description
            </label>
            <textarea
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              rows={4}
              style={{
                width: '100%',
                padding: '8px',
                borderRadius: '4px',
                border: '1px solid #ddd',
                fontFamily: 'inherit',
                resize: 'vertical',
              }}
            />
          </div>

          <div style={{ display: 'flex', gap: '8px', justifyContent: 'flex-end' }}>
            <button
              type="button"
              onClick={onClose}
              style={{
                padding: '8px 16px',
                backgroundColor: '#f0f0f0',
                border: 'none',
                borderRadius: '4px',
                cursor: 'pointer',
              }}
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={isLoading}
              style={{
                padding: '8px 16px',
                backgroundColor: '#3A7D44',
                color: 'white',
                border: 'none',
                borderRadius: '4px',
                cursor: 'pointer',
                fontWeight: 600,
              }}
            >
              {isLoading ? 'Saving...' : 'Save Changes'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
