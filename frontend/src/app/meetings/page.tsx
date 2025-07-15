"use client";

import { useSession } from "next-auth/react";
import { useState } from "react";
import GoogleMeetCreator from "../components/GoogleMeetCreator";
import ZoomCreator from "../components/ZoomCreator";

export default function MeetingsPage() {
  const { data: session, status } = useSession();
  const [activeTab, setActiveTab] = useState<"google" | "zoom">("google");

  if (status === "loading") {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-lg">Loading...</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-4xl mx-auto px-4">
        <div className="text-center mb-8">
          <h1 className="text-4xl font-bold text-gray-900 mb-4">
            Meeting Creator
          </h1>
          <p className="text-lg text-gray-600">
            Create Google Meet or Zoom meetings easily
          </p>
        </div>

        {!session ? (
          <div className="bg-white rounded-lg shadow-md p-8 text-center">
            <h2 className="text-2xl font-semibold mb-4">
              Authentication Required
            </h2>
            <p className="text-gray-600 mb-6">
              Please sign in with Google to create Google Meet meetings. Zoom
              meetings can be created without authentication.
            </p>
            <a
              href="/"
              className="bg-blue-600 text-white px-6 py-3 rounded-lg hover:bg-blue-700 inline-block"
            >
              Go to Sign In
            </a>
          </div>
        ) : (
          <div className="mb-6 text-center">
            <p className="text-lg text-gray-700">
              Welcome,{" "}
              <span className="font-semibold">{session.user?.name}</span>!
            </p>
          </div>
        )}

        {/* Tab Navigation */}
        <div className="bg-white rounded-lg shadow-md mb-6">
          <div className="border-b border-gray-200">
            <nav className="flex">
              <button
                onClick={() => setActiveTab("google")}
                className={`px-6 py-4 text-sm font-medium border-b-2 ${
                  activeTab === "google"
                    ? "border-blue-500 text-blue-600 bg-blue-50"
                    : "border-transparent text-gray-500 hover:text-gray-700 hover:bg-gray-50"
                }`}
              >
                <div className="flex items-center space-x-2">
                  <span>🎥</span>
                  <span>Google Meet</span>
                </div>
              </button>
              <button
                onClick={() => setActiveTab("zoom")}
                className={`px-6 py-4 text-sm font-medium border-b-2 ${
                  activeTab === "zoom"
                    ? "border-blue-500 text-blue-600 bg-blue-50"
                    : "border-transparent text-gray-500 hover:text-gray-700 hover:bg-gray-50"
                }`}
              >
                <div className="flex items-center space-x-2">
                  <span>📹</span>
                  <span>Zoom</span>
                </div>
              </button>
            </nav>
          </div>

          {/* Tab Content */}
          <div className="p-6">
            {activeTab === "google" && (
              <div>
                {session ? (
                  <GoogleMeetCreator />
                ) : (
                  <div className="text-center py-8">
                    <p className="text-gray-600 mb-4">
                      Google Meet creation requires Google authentication.
                    </p>
                    <a
                      href="/"
                      className="bg-blue-600 text-white px-6 py-3 rounded-lg hover:bg-blue-700 inline-block"
                    >
                      Sign in with Google
                    </a>
                  </div>
                )}
              </div>
            )}

            {activeTab === "zoom" && (
              <div>
                <ZoomCreator />
              </div>
            )}
          </div>
        </div>

        {/* Information Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="bg-white p-6 rounded-lg shadow-md">
            <h3 className="text-xl font-semibold mb-3 text-gray-800">
              🎥 Google Meet
            </h3>
            <ul className="text-sm text-gray-600 space-y-2">
              <li>• Integrates with Google Calendar</li>
              <li>• Automatically creates calendar events</li>
              <li>• Supports attendee invitations</li>
              <li>• Requires Google authentication</li>
            </ul>
          </div>

          <div className="bg-white p-6 rounded-lg shadow-md">
            <h3 className="text-xl font-semibold mb-3 text-gray-800">
              📹 Zoom
            </h3>
            <ul className="text-sm text-gray-600 space-y-2">
              <li>• Creates instant Zoom meetings</li>
              <li>• Provides join and host URLs</li>
              <li>• Configurable meeting settings</li>
              <li>• No user authentication required</li>
            </ul>
          </div>
        </div>

        {/* Navigation */}
        <div className="mt-8 text-center">
          <a href="/" className="text-blue-600 hover:text-blue-800 underline">
            ← Back to Home
          </a>
        </div>
      </div>
    </div>
  );
}
