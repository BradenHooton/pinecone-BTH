import { useState, useEffect } from 'react'
import { useNutritionSearch } from '../../hooks/useNutrition'
import { NutritionSearchResult } from '../../lib/api'

interface NutritionSearchModalProps {
  isOpen: boolean
  onClose: () => void
  onSelect: (result: NutritionSearchResult) => void
}

export function NutritionSearchModal({ isOpen, onClose, onSelect }: NutritionSearchModalProps) {
  const [query, setQuery] = useState('')
  const [debouncedQuery, setDebouncedQuery] = useState('')

  // Debounce search query (300ms)
  useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedQuery(query)
    }, 300)

    return () => clearTimeout(timer)
  }, [query])

  const { data, isLoading, error } = useNutritionSearch(debouncedQuery, isOpen)

  if (!isOpen) return null

  const handleSelect = (result: NutritionSearchResult) => {
    onSelect(result)
    setQuery('')
    onClose()
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
          maxWidth: '600px',
          width: '90%',
          maxHeight: '80vh',
          overflow: 'auto',
        }}
        onClick={(e) => e.stopPropagation()}
      >
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px' }}>
          <h2 style={{ margin: 0 }}>Search Nutrition Data</h2>
          <button
            onClick={onClose}
            style={{
              background: 'none',
              border: 'none',
              fontSize: '1.5rem',
              cursor: 'pointer',
              padding: '0',
              width: '32px',
              height: '32px',
            }}
          >
            ×
          </button>
        </div>

        <input
          type="text"
          placeholder="Search for foods..."
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          autoFocus
          style={{
            width: '100%',
            padding: '10px',
            fontSize: '1rem',
            border: '1px solid #ddd',
            borderRadius: '4px',
            marginBottom: '16px',
          }}
        />

        {isLoading && (
          <div style={{ textAlign: 'center', padding: '20px', color: '#666' }}>
            Searching...
          </div>
        )}

        {error && (
          <div style={{ padding: '20px', color: '#dc3545', textAlign: 'center' }}>
            Error: {error.message}
          </div>
        )}

        {data && data.data.length === 0 && query.length >= 2 && !isLoading && (
          <div style={{ textAlign: 'center', padding: '20px', color: '#666' }}>
            No results found. Try a different search term.
          </div>
        )}

        {data && data.data.length > 0 && (
          <div>
            <div style={{ marginBottom: '8px', fontWeight: 'bold', color: '#666' }}>
              {data.meta.total} results
            </div>
            <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
              {data.data.map((result) => (
                <div
                  key={result.fdc_id}
                  onClick={() => handleSelect(result)}
                  style={{
                    padding: '12px',
                    border: '1px solid #e0e0e0',
                    borderRadius: '4px',
                    cursor: 'pointer',
                    transition: 'background-color 0.2s',
                  }}
                  onMouseEnter={(e) => {
                    e.currentTarget.style.backgroundColor = '#f5f5f5'
                  }}
                  onMouseLeave={(e) => {
                    e.currentTarget.style.backgroundColor = 'white'
                  }}
                >
                  <div style={{ fontWeight: 'bold', marginBottom: '4px' }}>
                    {result.description}
                  </div>
                  <div style={{ fontSize: '0.875rem', color: '#666' }}>
                    {result.data_type}
                  </div>
                  {(result.calories || result.protein_g || result.carbs_g || result.fat_g) && (
                    <div style={{ fontSize: '0.875rem', marginTop: '8px', display: 'flex', gap: '16px' }}>
                      {result.calories && <span>Cal: {result.calories.toFixed(0)}</span>}
                      {result.protein_g && <span>Protein: {result.protein_g.toFixed(1)}g</span>}
                      {result.carbs_g && <span>Carbs: {result.carbs_g.toFixed(1)}g</span>}
                      {result.fat_g && <span>Fat: {result.fat_g.toFixed(1)}g</span>}
                    </div>
                  )}
                </div>
              ))}
            </div>
          </div>
        )}

        {query.length < 2 && (
          <div style={{ textAlign: 'center', padding: '20px', color: '#666' }}>
            Enter at least 2 characters to search
          </div>
        )}
      </div>
    </div>
  )
}
