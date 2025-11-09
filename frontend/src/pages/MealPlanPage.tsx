import { useState } from 'react'
import { useMealPlansByDateRange } from '../hooks/useMealPlan'
import { MealType } from '../lib/api'

const MEAL_TYPES: MealType[] = ['breakfast', 'lunch', 'snack', 'dinner', 'dessert']

function formatDate(date: Date): string {
  return date.toISOString().split('T')[0]
}

function getWeekDates(startDate: Date): Date[] {
  const dates: Date[] = []
  for (let i = 0; i < 7; i++) {
    const date = new Date(startDate)
    date.setDate(startDate.getDate() + i)
    dates.push(date)
  }
  return dates
}

export function MealPlanPage() {
  const [weekStart, setWeekStart] = useState(() => {
    const today = new Date()
    const dayOfWeek = today.getDay()
    const diff = dayOfWeek === 0 ? -6 : 1 - dayOfWeek // Adjust when day is Sunday
    const monday = new Date(today)
    monday.setDate(today.getDate() + diff)
    return monday
  })

  const weekDates = getWeekDates(weekStart)
  const startDateStr = formatDate(weekDates[0])
  const endDateStr = formatDate(weekDates[6])

  const { data, isLoading, error } = useMealPlansByDateRange(startDateStr, endDateStr)

  const goToPreviousWeek = () => {
    const newStart = new Date(weekStart)
    newStart.setDate(weekStart.getDate() - 7)
    setWeekStart(newStart)
  }

  const goToNextWeek = () => {
    const newStart = new Date(weekStart)
    newStart.setDate(weekStart.getDate() + 7)
    setWeekStart(newStart)
  }

  const goToToday = () => {
    const today = new Date()
    const dayOfWeek = today.getDay()
    const diff = dayOfWeek === 0 ? -6 : 1 - dayOfWeek
    const monday = new Date(today)
    monday.setDate(today.getDate() + diff)
    setWeekStart(monday)
  }

  if (isLoading) {
    return <div style={{ padding: '20px' }}>Loading meal plans...</div>
  }

  if (error) {
    return (
      <div style={{ padding: '20px', color: 'red' }}>
        Error loading meal plans: {error.message}
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
        <h1 style={{ margin: 0 }}>Meal Plan</h1>
        <div style={{ display: 'flex', gap: '8px' }}>
          <button
            onClick={goToPreviousWeek}
            style={{
              padding: '8px 16px',
              backgroundColor: '#f0f0f0',
              border: 'none',
              borderRadius: '4px',
              cursor: 'pointer',
            }}
          >
            Previous Week
          </button>
          <button
            onClick={goToToday}
            style={{
              padding: '8px 16px',
              backgroundColor: '#3A7D44',
              color: 'white',
              border: 'none',
              borderRadius: '4px',
              cursor: 'pointer',
            }}
          >
            Today
          </button>
          <button
            onClick={goToNextWeek}
            style={{
              padding: '8px 16px',
              backgroundColor: '#f0f0f0',
              border: 'none',
              borderRadius: '4px',
              cursor: 'pointer',
            }}
          >
            Next Week
          </button>
        </div>
      </div>

      <div style={{ overflowX: 'auto' }}>
        <table style={{ width: '100%', borderCollapse: 'collapse', minWidth: '800px' }}>
          <thead>
            <tr>
              <th style={{ padding: '12px', textAlign: 'left', borderBottom: '2px solid #ddd' }}>
                Meal
              </th>
              {weekDates.map((date) => (
                <th
                  key={date.toISOString()}
                  style={{
                    padding: '12px',
                    textAlign: 'center',
                    borderBottom: '2px solid #ddd',
                    minWidth: '120px',
                  }}
                >
                  <div>{date.toLocaleDateString('en-US', { weekday: 'short' })}</div>
                  <div style={{ fontSize: '0.875rem', color: '#666' }}>
                    {date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })}
                  </div>
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {MEAL_TYPES.map((mealType) => (
              <tr key={mealType}>
                <td
                  style={{
                    padding: '12px',
                    borderBottom: '1px solid #eee',
                    fontWeight: 'bold',
                    textTransform: 'capitalize',
                  }}
                >
                  {mealType}
                </td>
                {weekDates.map((date) => {
                  const dateStr = formatDate(date)
                  const mealPlan = data?.data.find((mp) => mp.plan_date.startsWith(dateStr))
                  const meals = mealPlan?.meals.filter((m) => m.meal_type === mealType) || []

                  return (
                    <td
                      key={date.toISOString()}
                      style={{
                        padding: '8px',
                        borderBottom: '1px solid #eee',
                        verticalAlign: 'top',
                      }}
                    >
                      {meals.length === 0 ? (
                        <div
                          style={{
                            padding: '8px',
                            textAlign: 'center',
                            color: '#999',
                            fontSize: '0.875rem',
                            cursor: 'pointer',
                            border: '1px dashed #ddd',
                            borderRadius: '4px',
                          }}
                        >
                          +
                        </div>
                      ) : (
                        <div style={{ display: 'flex', flexDirection: 'column', gap: '4px' }}>
                          {meals.map((meal) => (
                            <div
                              key={meal.id}
                              style={{
                                padding: '8px',
                                backgroundColor: meal.out_of_kitchen ? '#fff4e6' : '#f0f9f4',
                                borderRadius: '4px',
                                fontSize: '0.875rem',
                                cursor: 'pointer',
                              }}
                            >
                              {meal.out_of_kitchen ? (
                                <div style={{ fontStyle: 'italic', color: '#666' }}>
                                  Out of Kitchen
                                </div>
                              ) : (
                                <>
                                  <div style={{ fontWeight: 'bold' }}>
                                    {meal.recipe?.title || 'Recipe'}
                                  </div>
                                  {meal.servings && (
                                    <div style={{ color: '#666', fontSize: '0.75rem' }}>
                                      {meal.servings} serving{meal.servings !== 1 ? 's' : ''}
                                    </div>
                                  )}
                                </>
                              )}
                            </div>
                          ))}
                        </div>
                      )}
                    </td>
                  )
                })}
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <div
        style={{
          marginTop: '24px',
          padding: '16px',
          backgroundColor: '#f8f9fa',
          borderRadius: '8px',
        }}
      >
        <h3 style={{ margin: '0 0 12px 0', fontSize: '1rem' }}>Legend</h3>
        <div style={{ display: 'flex', gap: '24px', fontSize: '0.875rem' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
            <div
              style={{
                width: '20px',
                height: '20px',
                backgroundColor: '#f0f9f4',
                borderRadius: '4px',
              }}
            />
            <span>Planned Meal</span>
          </div>
          <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
            <div
              style={{
                width: '20px',
                height: '20px',
                backgroundColor: '#fff4e6',
                borderRadius: '4px',
              }}
            />
            <span>Out of Kitchen</span>
          </div>
        </div>
        <p style={{ margin: '12px 0 0 0', color: '#666', fontSize: '0.875rem' }}>
          Click on a meal slot to edit or add meals. (Editing functionality coming soon)
        </p>
      </div>
    </div>
  )
}
