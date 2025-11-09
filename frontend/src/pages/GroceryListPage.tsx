import { useState } from 'react'
import {
  useGroceryLists,
  useCreateGroceryList,
  useUpdateItemStatus,
  useAddManualItem,
  useDeleteGroceryList,
} from '../hooks/useGroceryList'
import {
  GroceryDepartment,
  GroceryItemStatus,
  GroceryListItem,
  CreateManualItemRequest,
} from '../lib/api'

const DEPARTMENT_ORDER: GroceryDepartment[] = [
  'produce',
  'meat',
  'seafood',
  'dairy',
  'bakery',
  'frozen',
  'pantry',
  'spices',
  'beverages',
  'other',
]

const DEPARTMENT_NAMES: Record<GroceryDepartment, string> = {
  produce: 'Produce',
  meat: 'Meat & Poultry',
  seafood: 'Seafood',
  dairy: 'Dairy & Eggs',
  bakery: 'Bakery',
  frozen: 'Frozen Foods',
  pantry: 'Pantry Staples',
  spices: 'Spices & Seasonings',
  beverages: 'Beverages',
  other: 'Other',
}

function formatDate(date: Date): string {
  return date.toISOString().split('T')[0]
}

function getWeekRange(startDate: Date): { start: string; end: string } {
  const end = new Date(startDate)
  end.setDate(startDate.getDate() + 6)
  return {
    start: formatDate(startDate),
    end: formatDate(end),
  }
}

function groupByDepartment(items: GroceryListItem[]): Map<GroceryDepartment, GroceryListItem[]> {
  const grouped = new Map<GroceryDepartment, GroceryListItem[]>()

  // Initialize all departments
  DEPARTMENT_ORDER.forEach(dept => {
    grouped.set(dept, [])
  })

  // Group items
  items.forEach(item => {
    const dept = item.department || 'other'
    const existing = grouped.get(dept) || []
    grouped.set(dept, [...existing, item])
  })

  return grouped
}

export function GroceryListPage() {
  const [weekStart, setWeekStart] = useState(() => {
    const today = new Date()
    const dayOfWeek = today.getDay()
    const diff = dayOfWeek === 0 ? -6 : 1 - dayOfWeek
    const monday = new Date(today)
    monday.setDate(today.getDate() + diff)
    return monday
  })

  const [showAddItemModal, setShowAddItemModal] = useState(false)
  const [selectedListId, setSelectedListId] = useState<string | null>(null)

  const { data, isLoading, error } = useGroceryLists()
  const createGroceryList = useCreateGroceryList()
  const updateItemStatus = useUpdateItemStatus()
  const addManualItem = useAddManualItem()
  const deleteGroceryList = useDeleteGroceryList()

  const handleCreateList = () => {
    const range = getWeekRange(weekStart)
    createGroceryList.mutate({
      start_date: range.start,
      end_date: range.end,
    })
  }

  const handleStatusChange = (itemId: string, status: GroceryItemStatus) => {
    updateItemStatus.mutate({ itemId, data: { status } })
  }

  const handleAddManualItem = (data: CreateManualItemRequest) => {
    if (!selectedListId) return

    addManualItem.mutate(
      { listId: selectedListId, data },
      {
        onSuccess: () => {
          setShowAddItemModal(false)
        },
      }
    )
  }

  const handleDeleteList = (id: string) => {
    if (confirm('Are you sure you want to delete this grocery list?')) {
      deleteGroceryList.mutate(id)
    }
  }

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
    return <div style={{ padding: '20px' }}>Loading grocery lists...</div>
  }

  if (error) {
    return (
      <div style={{ padding: '20px', color: 'red' }}>
        Error loading grocery lists: {error.message}
      </div>
    )
  }

  const currentWeekRange = getWeekRange(weekStart)

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
        <h1 style={{ margin: 0 }}>Grocery Lists</h1>
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
            This Week
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
          <button
            onClick={handleCreateList}
            disabled={createGroceryList.isPending}
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
            {createGroceryList.isPending ? 'Creating...' : 'Generate List'}
          </button>
        </div>
      </div>

      <div style={{ marginBottom: '16px', color: '#666' }}>
        Week of {currentWeekRange.start} to {currentWeekRange.end}
      </div>

      {data?.data && data.data.length === 0 && (
        <div
          style={{
            padding: '40px',
            textAlign: 'center',
            backgroundColor: '#f8f9fa',
            borderRadius: '8px',
          }}
        >
          <p>No grocery lists yet. Generate one to get started!</p>
        </div>
      )}

      {data?.data && data.data.map((list) => {
        const groupedItems = groupByDepartment(list.items || [])

        return (
          <div
            key={list.id}
            style={{
              marginBottom: '32px',
              border: '1px solid #ddd',
              borderRadius: '8px',
              overflow: 'hidden',
            }}
          >
            <div
              style={{
                padding: '16px',
                backgroundColor: '#f8f9fa',
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
              }}
            >
              <div>
                <h2 style={{ margin: 0, fontSize: '1.25rem' }}>
                  Grocery List: {list.start_date} to {list.end_date}
                </h2>
                <p style={{ margin: '4px 0 0 0', color: '#666', fontSize: '0.875rem' }}>
                  {list.items?.length || 0} items
                </p>
              </div>
              <div style={{ display: 'flex', gap: '8px' }}>
                <button
                  onClick={() => {
                    setSelectedListId(list.id)
                    setShowAddItemModal(true)
                  }}
                  style={{
                    padding: '8px 16px',
                    backgroundColor: '#3A7D44',
                    color: 'white',
                    border: 'none',
                    borderRadius: '4px',
                    cursor: 'pointer',
                  }}
                >
                  Add Item
                </button>
                <button
                  onClick={() => handleDeleteList(list.id)}
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

            <div style={{ padding: '16px' }}>
              {DEPARTMENT_ORDER.map((dept) => {
                const items = groupedItems.get(dept) || []
                if (items.length === 0) return null

                return (
                  <div key={dept} style={{ marginBottom: '24px' }}>
                    <h3
                      style={{
                        margin: '0 0 12px 0',
                        fontSize: '1.125rem',
                        color: '#3A7D44',
                        borderBottom: '2px solid #3A7D44',
                        paddingBottom: '4px',
                      }}
                    >
                      {DEPARTMENT_NAMES[dept]}
                    </h3>
                    <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
                      {items.map((item) => (
                        <label
                          key={item.id}
                          style={{
                            display: 'flex',
                            alignItems: 'center',
                            gap: '12px',
                            padding: '8px',
                            backgroundColor:
                              item.status === 'bought'
                                ? '#d4edda'
                                : item.status === 'have_on_hand'
                                ? '#fff3cd'
                                : '#fff',
                            borderRadius: '4px',
                            cursor: 'pointer',
                            transition: 'background-color 0.2s',
                          }}
                        >
                          <input
                            type="checkbox"
                            checked={item.status === 'bought'}
                            onChange={(e) => {
                              const newStatus: GroceryItemStatus = e.target.checked
                                ? 'bought'
                                : 'pending'
                              handleStatusChange(item.id, newStatus)
                            }}
                            style={{ cursor: 'pointer', width: '18px', height: '18px' }}
                          />
                          <div style={{ flex: 1 }}>
                            <div
                              style={{
                                textDecoration: item.status === 'bought' ? 'line-through' : 'none',
                                fontWeight: 500,
                              }}
                            >
                              {item.item_name}
                              {item.quantity && item.unit && (
                                <span style={{ color: '#666', marginLeft: '8px' }}>
                                  ({item.quantity.toFixed(2)} {item.unit})
                                </span>
                              )}
                            </div>
                            {item.is_manual && (
                              <div style={{ fontSize: '0.75rem', color: '#666' }}>
                                Manual item
                              </div>
                            )}
                          </div>
                          <select
                            value={item.status}
                            onChange={(e) =>
                              handleStatusChange(item.id, e.target.value as GroceryItemStatus)
                            }
                            onClick={(e) => e.stopPropagation()}
                            style={{
                              padding: '4px 8px',
                              borderRadius: '4px',
                              border: '1px solid #ddd',
                              fontSize: '0.875rem',
                            }}
                          >
                            <option value="pending">Pending</option>
                            <option value="bought">Bought</option>
                            <option value="have_on_hand">Have on Hand</option>
                          </select>
                        </label>
                      ))}
                    </div>
                  </div>
                )
              })}
            </div>
          </div>
        )
      })}

      {showAddItemModal && (
        <AddManualItemModal
          onClose={() => setShowAddItemModal(false)}
          onSubmit={handleAddManualItem}
          isLoading={addManualItem.isPending}
        />
      )}
    </div>
  )
}

interface AddManualItemModalProps {
  onClose: () => void
  onSubmit: (data: CreateManualItemRequest) => void
  isLoading: boolean
}

function AddManualItemModal({ onClose, onSubmit, isLoading }: AddManualItemModalProps) {
  const [itemName, setItemName] = useState('')
  const [quantity, setQuantity] = useState('')
  const [unit, setUnit] = useState('')
  const [department, setDepartment] = useState<GroceryDepartment>('other')

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()

    onSubmit({
      item_name: itemName,
      quantity: quantity ? parseFloat(quantity) : undefined,
      unit: unit || undefined,
      department,
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
        <h2 style={{ margin: '0 0 20px 0' }}>Add Manual Item</h2>

        <form onSubmit={handleSubmit}>
          <div style={{ marginBottom: '16px' }}>
            <label style={{ display: 'block', marginBottom: '4px', fontWeight: 500 }}>
              Item Name *
            </label>
            <input
              type="text"
              value={itemName}
              onChange={(e) => setItemName(e.target.value)}
              required
              style={{
                width: '100%',
                padding: '8px',
                borderRadius: '4px',
                border: '1px solid #ddd',
              }}
            />
          </div>

          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '16px', marginBottom: '16px' }}>
            <div>
              <label style={{ display: 'block', marginBottom: '4px', fontWeight: 500 }}>
                Quantity
              </label>
              <input
                type="number"
                step="any"
                value={quantity}
                onChange={(e) => setQuantity(e.target.value)}
                style={{
                  width: '100%',
                  padding: '8px',
                  borderRadius: '4px',
                  border: '1px solid #ddd',
                }}
              />
            </div>

            <div>
              <label style={{ display: 'block', marginBottom: '4px', fontWeight: 500 }}>
                Unit
              </label>
              <input
                type="text"
                value={unit}
                onChange={(e) => setUnit(e.target.value)}
                placeholder="lbs, oz, cups, etc."
                style={{
                  width: '100%',
                  padding: '8px',
                  borderRadius: '4px',
                  border: '1px solid #ddd',
                }}
              />
            </div>
          </div>

          <div style={{ marginBottom: '24px' }}>
            <label style={{ display: 'block', marginBottom: '4px', fontWeight: 500 }}>
              Department
            </label>
            <select
              value={department}
              onChange={(e) => setDepartment(e.target.value as GroceryDepartment)}
              style={{
                width: '100%',
                padding: '8px',
                borderRadius: '4px',
                border: '1px solid #ddd',
              }}
            >
              {DEPARTMENT_ORDER.map((dept) => (
                <option key={dept} value={dept}>
                  {DEPARTMENT_NAMES[dept]}
                </option>
              ))}
            </select>
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
              {isLoading ? 'Adding...' : 'Add Item'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
