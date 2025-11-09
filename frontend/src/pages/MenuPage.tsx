import { useState } from 'react'
import { useNavigate } from '@tanstack/react-router'
import { useRecommendRecipes } from '../hooks/useMenu'
import { RecipeRecommendation } from '../lib/api'

export function MenuPage() {
  const [ingredients, setIngredients] = useState<string[]>([])
  const [currentInput, setCurrentInput] = useState('')
  const navigate = useNavigate()

  const recommendRecipes = useRecommendRecipes()

  const handleAddIngredient = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter' && currentInput.trim()) {
      e.preventDefault()
      setIngredients([...ingredients, currentInput.trim()])
      setCurrentInput('')
    }
  }

  const handleRemoveIngredient = (index: number) => {
    setIngredients(ingredients.filter((_, i) => i !== index))
  }

  const handleRecommend = () => {
    if (ingredients.length > 0) {
      recommendRecipes.mutate({ ingredients })
    }
  }

  const handleRecipeClick = (recipeId: string) => {
    navigate({ to: `/recipes/${recipeId}` })
  }

  return (
    <div style={{ minHeight: '100vh', backgroundColor: '#faf8f5' }}>
      {/* Header */}
      <div
        style={{
          backgroundColor: '#3A7D44',
          color: 'white',
          padding: '32px 20px',
          textAlign: 'center',
          borderBottom: '3px solid #2a5d34',
        }}
      >
        <h1
          style={{
            fontFamily: '"Playfair Display", serif',
            fontSize: '2.5rem',
            margin: '0 0 8px 0',
            fontWeight: 700,
          }}
        >
          Menu Recommendation
        </h1>
        <p
          style={{
            fontFamily: 'Inter, sans-serif',
            fontSize: '1.125rem',
            margin: 0,
            opacity: 0.95,
          }}
        >
          Tell us what you have, we'll suggest what to cook
        </p>
      </div>

      <div style={{ maxWidth: '900px', margin: '0 auto', padding: '40px 20px' }}>
        {/* Input Section */}
        <div
          style={{
            backgroundColor: 'white',
            borderRadius: '12px',
            padding: '32px',
            marginBottom: '40px',
            boxShadow: '0 2px 8px rgba(0, 0, 0, 0.08)',
          }}
        >
          <h2
            style={{
              fontFamily: '"Playfair Display", serif',
              fontSize: '1.5rem',
              margin: '0 0 20px 0',
              color: '#2a5d34',
            }}
          >
            Your Ingredients
          </h2>

          <div style={{ marginBottom: '16px' }}>
            <input
              type="text"
              value={currentInput}
              onChange={(e) => setCurrentInput(e.target.value)}
              onKeyDown={handleAddIngredient}
              placeholder="Type an ingredient and press Enter..."
              style={{
                width: '100%',
                padding: '12px 16px',
                fontSize: '1rem',
                border: '2px solid #ddd',
                borderRadius: '8px',
                outline: 'none',
                fontFamily: 'Inter, sans-serif',
              }}
              onFocus={(e) => (e.target.style.borderColor = '#3A7D44')}
              onBlur={(e) => (e.target.style.borderColor = '#ddd')}
            />
          </div>

          {/* Ingredient Chips */}
          {ingredients.length > 0 && (
            <div
              style={{
                display: 'flex',
                flexWrap: 'wrap',
                gap: '8px',
                marginBottom: '24px',
              }}
            >
              {ingredients.map((ing, index) => (
                <div
                  key={index}
                  style={{
                    backgroundColor: '#e8f5e9',
                    color: '#2a5d34',
                    padding: '8px 12px',
                    borderRadius: '20px',
                    display: 'flex',
                    alignItems: 'center',
                    gap: '8px',
                    fontSize: '0.9375rem',
                    fontFamily: 'Inter, sans-serif',
                  }}
                >
                  <span>{ing}</span>
                  <button
                    onClick={() => handleRemoveIngredient(index)}
                    style={{
                      background: 'none',
                      border: 'none',
                      color: '#2a5d34',
                      cursor: 'pointer',
                      padding: '0',
                      fontSize: '1.25rem',
                      lineHeight: 1,
                    }}
                  >
                    ×
                  </button>
                </div>
              ))}
            </div>
          )}

          <button
            onClick={handleRecommend}
            disabled={ingredients.length === 0 || recommendRecipes.isPending}
            style={{
              width: '100%',
              padding: '14px 24px',
              backgroundColor: ingredients.length > 0 ? '#3A7D44' : '#ccc',
              color: 'white',
              border: 'none',
              borderRadius: '8px',
              fontSize: '1.125rem',
              fontWeight: 600,
              cursor: ingredients.length > 0 ? 'pointer' : 'not-allowed',
              fontFamily: 'Inter, sans-serif',
              transition: 'background-color 0.2s',
            }}
          >
            {recommendRecipes.isPending ? 'Finding Recipes...' : 'Find Recipes'}
          </button>
        </div>

        {/* Results Section */}
        {recommendRecipes.data && (
          <div>
            <h2
              style={{
                fontFamily: '"Playfair Display", serif',
                fontSize: '2rem',
                margin: '0 0 24px 0',
                color: '#2a5d34',
                textAlign: 'center',
              }}
            >
              Recommended Menu
            </h2>

            {recommendRecipes.data.data.length === 0 ? (
              <div
                style={{
                  backgroundColor: 'white',
                  borderRadius: '12px',
                  padding: '40px',
                  textAlign: 'center',
                  boxShadow: '0 2px 8px rgba(0, 0, 0, 0.08)',
                }}
              >
                <p style={{ fontSize: '1.125rem', color: '#666', margin: 0 }}>
                  No recipes found with these ingredients. Try adding more ingredients!
                </p>
              </div>
            ) : (
              <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
                {recommendRecipes.data.data.map((recommendation: RecipeRecommendation) => (
                  <MenuCard
                    key={recommendation.recipe.id}
                    recommendation={recommendation}
                    onClick={() => handleRecipeClick(recommendation.recipe.id)}
                  />
                ))}
              </div>
            )}
          </div>
        )}

        {recommendRecipes.isError && (
          <div
            style={{
              backgroundColor: '#fee',
              color: '#c33',
              padding: '16px',
              borderRadius: '8px',
              textAlign: 'center',
            }}
          >
            Failed to get recommendations. Please try again.
          </div>
        )}
      </div>
    </div>
  )
}

interface MenuCardProps {
  recommendation: RecipeRecommendation
  onClick: () => void
}

function MenuCard({ recommendation, onClick }: MenuCardProps) {
  const { recipe, match_score, matched_ingredients, missing_ingredients } = recommendation

  // Determine match quality
  const matchQuality =
    match_score >= 80
      ? { label: 'Excellent Match', color: '#2a5d34' }
      : match_score >= 60
      ? { label: 'Good Match', color: '#3A7D44' }
      : match_score >= 40
      ? { label: 'Fair Match', color: '#7a9d54' }
      : { label: 'Partial Match', color: '#a0a0a0' }

  return (
    <div
      onClick={onClick}
      style={{
        backgroundColor: 'white',
        borderRadius: '12px',
        padding: '28px',
        boxShadow: '0 2px 8px rgba(0, 0, 0, 0.08)',
        cursor: 'pointer',
        transition: 'transform 0.2s, box-shadow 0.2s',
        border: '1px solid #f0f0f0',
      }}
      onMouseEnter={(e) => {
        e.currentTarget.style.transform = 'translateY(-2px)'
        e.currentTarget.style.boxShadow = '0 4px 12px rgba(0, 0, 0, 0.12)'
      }}
      onMouseLeave={(e) => {
        e.currentTarget.style.transform = 'translateY(0)'
        e.currentTarget.style.boxShadow = '0 2px 8px rgba(0, 0, 0, 0.08)'
      }}
    >
      <div
        style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'flex-start',
          marginBottom: '16px',
        }}
      >
        <h3
          style={{
            fontFamily: '"Playfair Display", serif',
            fontSize: '1.75rem',
            margin: 0,
            color: '#2a5d34',
            flex: 1,
          }}
        >
          {recipe.title}
        </h3>
        <div
          style={{
            backgroundColor: matchQuality.color,
            color: 'white',
            padding: '6px 14px',
            borderRadius: '20px',
            fontSize: '0.875rem',
            fontWeight: 600,
            whiteSpace: 'nowrap',
            marginLeft: '16px',
          }}
        >
          {Math.round(match_score)}%
        </div>
      </div>

      <div
        style={{
          display: 'inline-block',
          padding: '4px 12px',
          backgroundColor: '#f5f5f5',
          borderRadius: '4px',
          fontSize: '0.875rem',
          color: '#666',
          marginBottom: '16px',
        }}
      >
        {matchQuality.label}
      </div>

      {recipe.prep_time_minutes && recipe.cook_time_minutes && (
        <div
          style={{
            fontSize: '0.875rem',
            color: '#666',
            marginBottom: '16px',
            fontFamily: 'Inter, sans-serif',
          }}
        >
          ⏱️ Prep: {recipe.prep_time_minutes}min · Cook: {recipe.cook_time_minutes}min · Serves:{' '}
          {recipe.servings}
        </div>
      )}

      {/* Matched Ingredients */}
      {matched_ingredients.length > 0 && (
        <div style={{ marginBottom: '12px' }}>
          <div
            style={{
              fontSize: '0.875rem',
              fontWeight: 600,
              color: '#2a5d34',
              marginBottom: '6px',
              fontFamily: 'Inter, sans-serif',
            }}
          >
            ✓ You have ({matched_ingredients.length}):
          </div>
          <div
            style={{
              fontSize: '0.875rem',
              color: '#555',
              lineHeight: '1.6',
              fontFamily: 'Inter, sans-serif',
            }}
          >
            {matched_ingredients.join(', ')}
          </div>
        </div>
      )}

      {/* Missing Ingredients */}
      {missing_ingredients.length > 0 && (
        <div>
          <div
            style={{
              fontSize: '0.875rem',
              fontWeight: 600,
              color: '#c33',
              marginBottom: '6px',
              fontFamily: 'Inter, sans-serif',
            }}
          >
            You'll need ({missing_ingredients.length}):
          </div>
          <div
            style={{
              fontSize: '0.875rem',
              color: '#777',
              lineHeight: '1.6',
              fontFamily: 'Inter, sans-serif',
            }}
          >
            {missing_ingredients.join(', ')}
          </div>
        </div>
      )}
    </div>
  )
}
