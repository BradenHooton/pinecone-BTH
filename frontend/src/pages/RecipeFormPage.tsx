import { useState } from 'react'
import { useParams, useNavigate } from '@tanstack/react-router'
import { useRecipe, useCreateRecipe, useUpdateRecipe } from '../hooks/useRecipes'
import { CreateIngredientRequest, CreateInstructionRequest } from '../lib/api'

export function RecipeFormPage() {
  const { id } = useParams({ strict: false }) as { id?: string }
  const navigate = useNavigate()
  const isEditing = !!id

  const { data: recipeData } = useRecipe(id)
  const createMutation = useCreateRecipe()
  const updateMutation = useUpdateRecipe()

  const recipe = recipeData?.data

  const [title, setTitle] = useState(recipe?.title || '')
  const [servings, setServings] = useState(recipe?.servings || 1)
  const [servingSize, setServingSize] = useState(recipe?.serving_size || '')
  const [prepTime, setPrepTime] = useState(recipe?.prep_time_minutes || 0)
  const [cookTime, setCookTime] = useState(recipe?.cook_time_minutes || 0)
  const [ingredients, setIngredients] = useState<CreateIngredientRequest[]>(
    recipe?.ingredients?.map(ing => ({
      ingredient_name: ing.ingredient_name,
      quantity: ing.quantity,
      unit: ing.unit,
      department: ing.department,
    })) || []
  )
  const [instructions, setInstructions] = useState<CreateInstructionRequest[]>(
    recipe?.instructions?.map(inst => ({
      step_number: inst.step_number,
      instruction: inst.instruction,
    })) || []
  )

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    const data = {
      title,
      servings,
      serving_size: servingSize,
      prep_time_minutes: prepTime || undefined,
      cook_time_minutes: cookTime || undefined,
      ingredients,
      instructions,
    }

    try {
      if (isEditing && id) {
        await updateMutation.mutateAsync({ id, data })
      } else {
        await createMutation.mutateAsync(data)
      }
      navigate({ to: '/recipes' })
    } catch (err) {
      alert('Failed to save recipe')
    }
  }

  const addIngredient = () => {
    setIngredients([
      ...ingredients,
      { ingredient_name: '', quantity: 0, unit: '', department: 'other' },
    ])
  }

  const removeIngredient = (index: number) => {
    setIngredients(ingredients.filter((_, i) => i !== index))
  }

  const addInstruction = () => {
    setInstructions([
      ...instructions,
      { step_number: instructions.length + 1, instruction: '' },
    ])
  }

  const removeInstruction = (index: number) => {
    setInstructions(instructions.filter((_, i) => i !== index))
  }

  return (
    <div style={{ padding: '20px', maxWidth: '800px', margin: '0 auto' }}>
      <h1>{isEditing ? 'Edit Recipe' : 'Create Recipe'}</h1>

      <form onSubmit={handleSubmit}>
        <div style={{ marginBottom: '16px' }}>
          <label style={{ display: 'block', marginBottom: '4px', fontWeight: 'bold' }}>
            Title *
          </label>
          <input
            type="text"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            required
            style={{ width: '100%', padding: '8px', fontSize: '1rem' }}
          />
        </div>

        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '16px', marginBottom: '16px' }}>
          <div>
            <label style={{ display: 'block', marginBottom: '4px', fontWeight: 'bold' }}>
              Servings *
            </label>
            <input
              type="number"
              value={servings}
              onChange={(e) => setServings(parseInt(e.target.value))}
              required
              min="1"
              style={{ width: '100%', padding: '8px', fontSize: '1rem' }}
            />
          </div>
          <div>
            <label style={{ display: 'block', marginBottom: '4px', fontWeight: 'bold' }}>
              Serving Size *
            </label>
            <input
              type="text"
              value={servingSize}
              onChange={(e) => setServingSize(e.target.value)}
              required
              placeholder="e.g., 1 cup"
              style={{ width: '100%', padding: '8px', fontSize: '1rem' }}
            />
          </div>
        </div>

        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '16px', marginBottom: '24px' }}>
          <div>
            <label style={{ display: 'block', marginBottom: '4px', fontWeight: 'bold' }}>
              Prep Time (minutes)
            </label>
            <input
              type="number"
              value={prepTime}
              onChange={(e) => setPrepTime(parseInt(e.target.value))}
              min="0"
              style={{ width: '100%', padding: '8px', fontSize: '1rem' }}
            />
          </div>
          <div>
            <label style={{ display: 'block', marginBottom: '4px', fontWeight: 'bold' }}>
              Cook Time (minutes)
            </label>
            <input
              type="number"
              value={cookTime}
              onChange={(e) => setCookTime(parseInt(e.target.value))}
              min="0"
              style={{ width: '100%', padding: '8px', fontSize: '1rem' }}
            />
          </div>
        </div>

        <div style={{ marginBottom: '24px' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '8px' }}>
            <h2 style={{ margin: 0 }}>Ingredients</h2>
            <button
              type="button"
              onClick={addIngredient}
              style={{ padding: '6px 12px', backgroundColor: '#3A7D44', color: 'white', border: 'none', borderRadius: '4px', cursor: 'pointer' }}
            >
              Add Ingredient
            </button>
          </div>
          {ingredients.map((ing, index) => (
            <div key={index} style={{ display: 'grid', gridTemplateColumns: '1fr 1fr 1fr auto', gap: '8px', marginBottom: '8px' }}>
              <input
                type="text"
                placeholder="Name"
                value={ing.ingredient_name}
                onChange={(e) => {
                  const newIng = [...ingredients]
                  newIng[index].ingredient_name = e.target.value
                  setIngredients(newIng)
                }}
                required
                style={{ padding: '8px' }}
              />
              <input
                type="number"
                placeholder="Quantity"
                value={ing.quantity}
                onChange={(e) => {
                  const newIng = [...ingredients]
                  newIng[index].quantity = parseFloat(e.target.value)
                  setIngredients(newIng)
                }}
                required
                min="0"
                step="0.01"
                style={{ padding: '8px' }}
              />
              <input
                type="text"
                placeholder="Unit"
                value={ing.unit}
                onChange={(e) => {
                  const newIng = [...ingredients]
                  newIng[index].unit = e.target.value
                  setIngredients(newIng)
                }}
                required
                style={{ padding: '8px' }}
              />
              <button
                type="button"
                onClick={() => removeIngredient(index)}
                style={{ padding: '8px', backgroundColor: '#dc3545', color: 'white', border: 'none', borderRadius: '4px', cursor: 'pointer' }}
              >
                Remove
              </button>
            </div>
          ))}
        </div>

        <div style={{ marginBottom: '24px' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '8px' }}>
            <h2 style={{ margin: 0 }}>Instructions</h2>
            <button
              type="button"
              onClick={addInstruction}
              style={{ padding: '6px 12px', backgroundColor: '#3A7D44', color: 'white', border: 'none', borderRadius: '4px', cursor: 'pointer' }}
            >
              Add Instruction
            </button>
          </div>
          {instructions.map((inst, index) => (
            <div key={index} style={{ display: 'flex', gap: '8px', marginBottom: '8px', alignItems: 'flex-start' }}>
              <span style={{ fontWeight: 'bold', padding: '8px' }}>{index + 1}.</span>
              <textarea
                placeholder="Instruction"
                value={inst.instruction}
                onChange={(e) => {
                  const newInst = [...instructions]
                  newInst[index].instruction = e.target.value
                  newInst[index].step_number = index + 1
                  setInstructions(newInst)
                }}
                required
                rows={2}
                style={{ flex: 1, padding: '8px', fontSize: '1rem', resize: 'vertical' }}
              />
              <button
                type="button"
                onClick={() => removeInstruction(index)}
                style={{ padding: '8px', backgroundColor: '#dc3545', color: 'white', border: 'none', borderRadius: '4px', cursor: 'pointer' }}
              >
                Remove
              </button>
            </div>
          ))}
        </div>

        <div style={{ display: 'flex', gap: '12px' }}>
          <button
            type="submit"
            disabled={createMutation.isPending || updateMutation.isPending}
            style={{
              padding: '10px 24px',
              backgroundColor: '#3A7D44',
              color: 'white',
              border: 'none',
              borderRadius: '4px',
              cursor: 'pointer',
              fontSize: '1rem',
            }}
          >
            {createMutation.isPending || updateMutation.isPending ? 'Saving...' : 'Save Recipe'}
          </button>
          <button
            type="button"
            onClick={() => navigate({ to: '/recipes' })}
            style={{
              padding: '10px 24px',
              backgroundColor: '#f0f0f0',
              border: 'none',
              borderRadius: '4px',
              cursor: 'pointer',
              fontSize: '1rem',
            }}
          >
            Cancel
          </button>
        </div>
      </form>
    </div>
  )
}
