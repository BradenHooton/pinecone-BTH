import { Recipe } from '../../lib/api'

interface RecipeCardProps {
  recipe: Recipe
  onClick?: () => void
}

export function RecipeCard({ recipe, onClick }: RecipeCardProps) {
  return (
    <div
      className="recipe-card"
      onClick={onClick}
      style={{
        cursor: onClick ? 'pointer' : 'default',
        border: '1px solid #e0e0e0',
        borderRadius: '8px',
        padding: '16px',
        backgroundColor: 'white',
        transition: 'box-shadow 0.2s',
      }}
      onMouseEnter={(e) => {
        if (onClick) {
          e.currentTarget.style.boxShadow = '0 4px 12px rgba(0,0,0,0.1)'
        }
      }}
      onMouseLeave={(e) => {
        e.currentTarget.style.boxShadow = 'none'
      }}
    >
      {recipe.image_url && (
        <img
          src={recipe.image_url}
          alt={recipe.title}
          style={{
            width: '100%',
            height: '200px',
            objectFit: 'cover',
            borderRadius: '4px',
            marginBottom: '12px',
          }}
        />
      )}
      <h3 style={{ margin: '0 0 8px 0', fontSize: '1.25rem' }}>
        {recipe.title}
      </h3>
      <div
        style={{
          display: 'flex',
          gap: '16px',
          fontSize: '0.875rem',
          color: '#666',
        }}
      >
        <span>Servings: {recipe.servings}</span>
        {recipe.total_time_minutes > 0 && (
          <span>Time: {recipe.total_time_minutes} min</span>
        )}
      </div>
      {recipe.tags && recipe.tags.length > 0 && (
        <div style={{ marginTop: '8px', display: 'flex', gap: '4px', flexWrap: 'wrap' }}>
          {recipe.tags.map((tag) => (
            <span
              key={tag.id}
              style={{
                padding: '2px 8px',
                backgroundColor: '#3A7D44',
                color: 'white',
                borderRadius: '4px',
                fontSize: '0.75rem',
              }}
            >
              {tag.tag_name}
            </span>
          ))}
        </div>
      )}
    </div>
  )
}
