'use client'

import { useSession, signIn, signOut } from 'next-auth/react'
import { useState } from 'react'

export default function Home() {
  const { data: session, status } = useSession()
  const [calendarData, setCalendarData] = useState(null)



  const fetchCalendarData = async () => {
    if (!session?.accessToken) return

    try {
      const response = await fetch('http://localhost:8080/api/calendar', {
        headers: {
          'Authorization': `Bearer ${session.accessToken}`
        }
      })
      const data = await response.json()
      setCalendarData(data)
    } catch (error) {
      console.error('Error fetching calendar data:', error)
    }
  }

  if (status === 'loading') return <div>Loading...</div>

  
  // Slack OAuth認証URLを生成
  // stateで状態を渡せる
  // Incoming Webhooksをonにする
  // Manage DistributionをPublic Distributionにする
  // Redirect URLsも設定とOAuth Tokensをonにする
  const slackAuthUrl = `https://slack.com/oauth/v2/authorize?client_id=2921580744176.9115187898726&scope=incoming-webhook&state=test`;


  return (
    <div className="min-h-screen p-8">
      <div className="max-w-4xl mx-auto">
        <h1 className="text-3xl font-bold mb-8">Calendar & Slack Integration</h1>
        
        {!session ? (
          <div className="text-center space-y-6">
            <div className="space-y-4">
              <button
                onClick={() => {
                  const myToken = "my-custom-token";
                  signIn("google", {
                    callbackUrl: `/?myToken=${myToken}`,
                  });
                }}
                className="bg-blue-500 text-white px-6 py-3 rounded-lg hover:bg-blue-600 mr-4"
              >
                Sign in with Google
              </button>
              <a
                href={slackAuthUrl}
                target="_blank"
                className="bg-green-500 text-white px-6 py-3 rounded-lg hover:bg-green-600"
              >
                Sign in with Slack
              </a>
            </div>
            <div className="mt-8">
              <a
                href="/meetings"
                className="bg-purple-500 text-white px-6 py-3 rounded-lg hover:bg-purple-600 inline-block"
              >
                Create Meetings (Zoom available without sign-in)
              </a>
            </div>
          </div>
        ) : (
          <div>
            <div className="mb-6">
              <p className="text-lg">Welcome, {session.user?.name}!</p>
              <button
                onClick={() => signOut()}
                className="bg-red-500 text-white px-4 py-2 rounded mt-2 hover:bg-red-600"
              >
                Sign out
              </button>
            </div>
            
            <div className="mb-6 space-x-4">
              <button
                onClick={fetchCalendarData}
                className="bg-blue-500 text-white px-6 py-3 rounded-lg hover:bg-blue-600"
              >
                Fetch Calendar Data
              </button>
              <a
                href="/meetings"
                className="bg-purple-500 text-white px-6 py-3 rounded-lg hover:bg-purple-600 inline-block"
              >
                Create Meetings
              </a>
            </div>

            {calendarData && (
              <div className="bg-gray-100 p-6 rounded-lg">
                <h2 className="text-xl font-semibold mb-4">Calendar Events</h2>
                <pre className="text-sm overflow-auto">
                  {JSON.stringify(calendarData, null, 2)}
                </pre>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  )
}
