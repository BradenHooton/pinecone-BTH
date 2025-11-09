import { useParams, useNavigate } from '@tanstack/react-router'
import { useRecipe, useDeleteRecipe } from '../hooks/useRecipes'

export function RecipeDetailPage() {
  const { id } = useParams({ strict: false }) as { id: string }
  const navigate = useNavigate()
  const { data, isLoading, error } = useRecipe(id)
  const deleteMutation = useDeleteRecipe()

  const handleDelete = async () => {
    if (confirm('Are you sure you want to delete this recipe?')) {
      try {
        await deleteMutation.mutateAsync(id)
        navigate({ to: '/recipes' })
      } catch (err) {
        alert('Failed to delete recipe')
      }
    }
  }

  if (isLoading) {
    return <div style={{ padding: '20px' }}>Loading recipe...</div>
  }

  if (error || !data) {
    return (
      <div style={{ padding: '20px', color: 'red' }}>
        Error loading recipe: {error?.message || 'Recipe not found'}
      </div>
    )
  }

  const recipe = data.data

  return (
    <div style={{ padding: '20px', maxWidth: '800px', margin: '0 auto' }}>
      <div
        style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          marginBottom: '24px',
        }}
      >
        <button
          onClick={() => navigate({ to: '/recipes' })}
          style={{
            padding: '8px 16px',
            backgroundColor: '#f0f0f0',
            border: 'none',
            borderRadius: '4px',
            cursor: 'pointer',
          }}
        >
          Back
        </button>
        <div style={{ display: 'flex', gap: '8px' }}>
          <button
            onClick={() => navigate({ to: `/recipes/${id}/edit` })}
            style={{
              padding: '8px 16px',
              backgroundColor: '#3A7D44',
              color: 'white',
              border: 'none',
              borderRadius: '4px',
              cursor: 'pointer',
            }}
          >
            Edit
          </button>
          <button
            onClick={handleDelete}
            style={{
              padding: '8px 16px',
              backgroundColor: '#dc3545',
              color: 'white',
              border: 'none',
              borderRadius: '4px',
              cursor: 'pointer',
            }}
          >
            Delete
          </button>
        </div>
      </div>

      {recipe.image_url && (
        <img
          src={recipe.image_url}
          alt={recipe.title}
          style={{
            width: '100%',
            maxHeight: '400px',
            objectFit: 'cover',
            borderRadius: '8px',
            marginBottom: '20px',
          }}
        />
      )}

      <h1 style={{ marginBottom: '16px' }}>{recipe.title}</h1>

      <div style={{ display: 'flex', gap: '24px', marginBottom: '24px', color: '#666' }}>
        <div>
          <strong>Servings:</strong> {recipe.servings}
        </div>
        {recipe.prep_time_minutes && (
          <div>
            <strong>Prep:</strong> {recipe.prep_time_minutes} min
          </div>
        )}
        {recipe.cook_time_minutes && (
          <div>
            <strong>Cook:</strong> {recipe.cook_time_minutes} min
          </div>
        )}
        <div>
          <strong>Total:</strong> {recipe.total_time_minutes} min
        </div>
      </div>

      {recipe.tags && recipe.tags.length > 0 && (
        <div style={{ marginBottom: '24px', display: 'flex', gap: '8px', flexWrap: 'wrap' }}>
          {recipe.tags.map((tag) => (
            <span
              key={tag.id}
              style={{
                padding: '4px 12px',
                backgroundColor: '#3A7D44',
                color: 'white',
                borderRadius: '4px',
                fontSize: '0.875rem',
              }}
            >
              {tag.tag_name}
            </span>
          ))}
        </div>
      )}

      <div style={{ marginBottom: '32px' }}>
        <h2>Ingredients</h2>
        <ul style={{ lineHeight: '1.8' }}>
          {recipe.ingredients?.map((ing) => (
            <li key={ing.id}>
              {ing.quantity} {ing.unit} {ing.ingredient_name}
            </li>
          ))}
        </ul>
      </div>

      <div style={{ marginBottom: '32px' }}>
        <h2>Instructions</h2>
        <ol style={{ lineHeight: '1.8' }}>
          {recipe.instructions?.map((inst) => (
            <li key={inst.id} style={{ marginBottom: '12px' }}>
              {inst.instruction}
            </li>
          ))}
        </ol>
      </div>

      {recipe.storage_notes && (
        <div style={{ marginBottom: '32px' }}>
          <h3>Storage Notes</h3>
          <p>{recipe.storage_notes}</p>
        </div>
      )}

      {recipe.source && (
        <div style={{ marginBottom: '32px' }}>
          <h3>Source</h3>
          <p>{recipe.source}</p>
        </div>
      )}

      {recipe.notes && (
        <div>
          <h3>Notes</h3>
          <p>{recipe.notes}</p>
        </div>
      )}
    </div>
  )
}
