import { useState } from 'react'
import { useNavigate } from '@tanstack/react-router'
import {
  useCookbooks,
  useCreateCookbook,
  useDeleteCookbook,
} from '../hooks/useCookbooks'
import { Cookbook, CreateCookbookRequest } from '../lib/api'

export function CookbooksPage() {
  const [showCreateModal, setShowCreateModal] = useState(false)
  const navigate = useNavigate()

  const { data, isLoading, error } = useCookbooks()
  const createCookbook = useCreateCookbook()
  const deleteCookbook = useDeleteCookbook()

  const handleCreateCookbook = (data: CreateCookbookRequest) => {
    createCookbook.mutate(data, {
      onSuccess: () => {
        setShowCreateModal(false)
      },
    })
  }

  const handleDeleteCookbook = (id: string) => {
    if (confirm('Are you sure you want to delete this cookbook?')) {
      deleteCookbook.mutate(id)
    }
  }

  const handleCookbookClick = (id: string) => {
    navigate({ to: `/cookbooks/${id}` })
  }

  if (isLoading) {
    return <div style={{ padding: '20px' }}>Loading cookbooks...</div>
  }

  if (error) {
    return (
      <div style={{ padding: '20px', color: 'red' }}>
        Error loading cookbooks: {error.message}
      </div>
    )
  }

  return (
    <div style={{ padding: '20px', maxWidth: '1200px', margin: '0 auto' }}>
      <div
        style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          marginBottom: '24px',
        }}
      >
        <h1 style={{ margin: 0 }}>My Cookbooks</h1>
        <button
          onClick={() => setShowCreateModal(true)}
          style={{
            padding: '10px 20px',
            backgroundColor: '#3A7D44',
            color: 'white',
            border: 'none',
            borderRadius: '6px',
            cursor: 'pointer',
            fontSize: '1rem',
            fontWeight: 600,
          }}
        >
          + New Cookbook
        </button>
      </div>

      {data?.data && data.data.length === 0 && (
        <div
          style={{
            padding: '60px 20px',
            textAlign: 'center',
            backgroundColor: '#f8f9fa',
            borderRadius: '8px',
          }}
        >
          <p style={{ fontSize: '1.125rem', color: '#666', margin: 0 }}>
            No cookbooks yet. Create your first cookbook to get started!
          </p>
        </div>
      )}

      <div
        style={{
          display: 'grid',
          gridTemplateColumns: 'repeat(auto-fill, minmax(300px, 1fr))',
          gap: '20px',
        }}
      >
        {data?.data && data.data.map((cookbook: Cookbook) => (
          <CookbookCard
            key={cookbook.id}
            cookbook={cookbook}
            onClick={() => handleCookbookClick(cookbook.id)}
            onDelete={() => handleDeleteCookbook(cookbook.id)}
          />
        ))}
      </div>

      {showCreateModal && (
        <CreateCookbookModal
          onClose={() => setShowCreateModal(false)}
          onSubmit={handleCreateCookbook}
          isLoading={createCookbook.isPending}
        />
      )}
    </div>
  )
}

interface CookbookCardProps {
  cookbook: Cookbook
  onClick: () => void
  onDelete: () => void
}

function CookbookCard({ cookbook, onClick, onDelete }: CookbookCardProps) {
  return (
    <div
      style={{
        border: '1px solid #ddd',
        borderRadius: '8px',
        padding: '20px',
        backgroundColor: 'white',
        cursor: 'pointer',
        transition: 'transform 0.2s, box-shadow 0.2s',
      }}
      onMouseEnter={(e) => {
        e.currentTarget.style.transform = 'translateY(-2px)'
        e.currentTarget.style.boxShadow = '0 4px 8px rgba(0, 0, 0, 0.1)'
      }}
      onMouseLeave={(e) => {
        e.currentTarget.style.transform = 'translateY(0)'
        e.currentTarget.style.boxShadow = 'none'
      }}
    >
      <div onClick={onClick}>
        <h3
          style={{
            margin: '0 0 8px 0',
            fontSize: '1.25rem',
            color: '#2a5d34',
            fontFamily: '"Playfair Display", serif',
          }}
        >
          {cookbook.name}
        </h3>
        {cookbook.description && (
          <p
            style={{
              margin: '0 0 16px 0',
              color: '#666',
              fontSize: '0.9375rem',
              lineHeight: '1.5',
            }}
          >
            {cookbook.description}
          </p>
        )}
        <div
          style={{
            display: 'flex',
            alignItems: 'center',
            gap: '6px',
            color: '#3A7D44',
            fontSize: '0.875rem',
            fontWeight: 500,
          }}
        >
          <span>📚</span>
          <span>{cookbook.recipe_count} recipe{cookbook.recipe_count !== 1 ? 's' : ''}</span>
        </div>
      </div>
      <div
        style={{
          marginTop: '16px',
          paddingTop: '16px',
          borderTop: '1px solid #eee',
          display: 'flex',
          gap: '8px',
        }}
      >
        <button
          onClick={(e) => {
            e.stopPropagation()
            onClick()
          }}
          style={{
            flex: 1,
            padding: '8px 16px',
            backgroundColor: '#3A7D44',
            color: 'white',
            border: 'none',
            borderRadius: '4px',
            cursor: 'pointer',
            fontSize: '0.875rem',
          }}
        >
          View
        </button>
        <button
          onClick={(e) => {
            e.stopPropagation()
            onDelete()
          }}
          style={{
            padding: '8px 16px',
            backgroundColor: '#dc3545',
            color: 'white',
            border: 'none',
            borderRadius: '4px',
            cursor: 'pointer',
            fontSize: '0.875rem',
          }}
        >
          Delete
        </button>
      </div>
    </div>
  )
}

interface CreateCookbookModalProps {
  onClose: () => void
  onSubmit: (data: CreateCookbookRequest) => void
  isLoading: boolean
}

function CreateCookbookModal({ onClose, onSubmit, isLoading }: CreateCookbookModalProps) {
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')

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
        <h2 style={{ margin: '0 0 20px 0' }}>Create New Cookbook</h2>

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
              {isLoading ? 'Creating...' : 'Create Cookbook'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
