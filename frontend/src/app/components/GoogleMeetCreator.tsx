'use client'

import { useState } from 'react'
import { useSession } from 'next-auth/react'

interface GoogleMeetResponse {
  eventId: string
  meetURL: string
  title: string
  description: string
  startTime: string
  endTime: string
  attendees: string[]
  createdAt: string
}

export default function GoogleMeetCreator() {
  const { data: session } = useSession()
  const [loading, setLoading] = useState(false)
  const [meetData, setMeetData] = useState<GoogleMeetResponse | null>(null)
  const [error, setError] = useState('')

  const [formData, setFormData] = useState({
    title: '',
    description: '',
    startTime: '',
    endTime: '',
    attendees: ''
  })

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!session?.accessToken) {
      setError('Please sign in first')
      return
    }

    setLoading(true)
    setError('')
    setMeetData(null)

    try {
      const attendeesArray = formData.attendees
        .split(',')
        .map(email => email.trim())
        .filter(email => email.length > 0)

      fetch('http://localhost:8080/health/google', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${session.accessToken}`
        },
      })

      const response = await fetch('http://localhost:8080/api/google-meet/create', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${session.accessToken}`
        },
        body: JSON.stringify({
          title: formData.title,
          description: formData.description,
          startTime: new Date(formData.startTime).toISOString(),
          endTime: new Date(formData.endTime).toISOString(),
          attendees: attendeesArray
        })
      })

      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.error || 'Failed to create Google Meet')
      }

      const data = await response.json()
      setMeetData(data)
      
      // Reset form
      setFormData({
        title: '',
        description: '',
        startTime: '',
        endTime: '',
        attendees: ''
      })
    } catch (error) {
      console.error('Error creating Google Meet:', error)
      setError(error instanceof Error ? error.message : 'Failed to create Google Meet')
    } finally {
      setLoading(false)
    }
  }

  const handleInputChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => {
    const { name, value } = e.target
    setFormData(prev => ({ ...prev, [name]: value }))
  }

  return (
    <div className="bg-white p-6 rounded-lg shadow-md">
      <h2 className="text-2xl font-bold mb-6 text-gray-800">Create Google Meet</h2>
      
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label htmlFor="title" className="block text-sm font-medium text-gray-700 mb-2">
            Meeting Title *
          </label>
          <input
            type="text"
            id="title"
            name="title"
            value={formData.title}
            onChange={handleInputChange}
            required
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="Enter meeting title"
          />
        </div>

        <div>
          <label htmlFor="description" className="block text-sm font-medium text-gray-700 mb-2">
            Description
          </label>
          <textarea
            id="description"
            name="description"
            value={formData.description}
            onChange={handleInputChange}
            rows={3}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="Enter meeting description (optional)"
          />
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label htmlFor="startTime" className="block text-sm font-medium text-gray-700 mb-2">
              Start Time *
            </label>
            <input
              type="datetime-local"
              id="startTime"
              name="startTime"
              value={formData.startTime}
              onChange={handleInputChange}
              required
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <div>
            <label htmlFor="endTime" className="block text-sm font-medium text-gray-700 mb-2">
              End Time *
            </label>
            <input
              type="datetime-local"
              id="endTime"
              name="endTime"
              value={formData.endTime}
              onChange={handleInputChange}
              required
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
        </div>

        <div>
          <label htmlFor="attendees" className="block text-sm font-medium text-gray-700 mb-2">
            Attendees (Email addresses)
          </label>
          <input
            type="text"
            id="attendees"
            name="attendees"
            value={formData.attendees}
            onChange={handleInputChange}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="email1@example.com, email2@example.com"
          />
          <p className="text-sm text-gray-500 mt-1">
            Separate multiple emails with commas
          </p>
        </div>

        <button
          type="submit"
          disabled={loading || !session}
          className="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed"
        >
          {loading ? 'Creating Meeting...' : 'Create Google Meet'}
        </button>
      </form>

      {error && (
        <div className="mt-4 p-4 bg-red-100 border border-red-400 text-red-700 rounded-md">
          {error}
        </div>
      )}

      {meetData && (
        <div className="mt-6 p-4 bg-green-100 border border-green-400 rounded-md">
          <h3 className="text-lg font-semibold text-green-800 mb-3">
            Google Meet Created Successfully!
          </h3>
          <div className="space-y-2 text-sm">
            <p><strong>Title:</strong> {meetData.title}</p>
            <p><strong>Event ID:</strong> {meetData.eventId}</p>
            <p><strong>Start:</strong> {new Date(meetData.startTime).toLocaleString()}</p>
            <p><strong>End:</strong> {new Date(meetData.endTime).toLocaleString()}</p>
            {meetData.meetURL && (
              <p>
                <strong>Meet URL:</strong>{' '}
                <a 
                  href={meetData.meetURL} 
                  target="_blank" 
                  rel="noopener noreferrer"
                  className="text-blue-600 hover:underline"
                >
                  {meetData.meetURL}
                </a>
              </p>
            )}
          </div>
        </div>
      )}
    </div>
  )
}