import type { ReactNode } from "react";
import { useAuth } from "react-oidc-context";
import { Toaster } from "react-hot-toast";

interface AuthGuardProps {
  children: ReactNode;
}

// Ensure user is authenticated before rendering children
export default function AuthGuard({ children }: AuthGuardProps) {
  const auth = useAuth();

  // Show loading indicator while authentication state initializes
  if (auth.isLoading) {
    return (
      <div className="min-h-screen flex flex-col items-center justify-center bg-gray-50">
        <div className="w-8 h-8 border-2 border-gray-200 border-t-blue-600 rounded-full animate-spin mb-4" />
        <p className="text-gray-500 text-sm">Loading authentication...</p>
      </div>
    );
  }

  // Display sign-in prompt when unauthenticated
  if (!auth.isAuthenticated) {
    return (
      <div className="min-h-screen bg-gray-50 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
        <div className="sm:mx-auto sm:w-full sm:max-w-md">
          <div className="bg-white py-12 px-4 shadow sm:rounded-lg sm:px-10 border border-gray-200 text-center">
            <h2 className="text-2xl font-bold text-gray-900 mb-2">Welcome</h2>
            <p className="text-gray-500 text-sm mb-8">Please sign in to access the employee directory.</p>
            <button
              onClick={() => auth.signinRedirect()}
              className="w-full flex justify-center py-2.5 px-4 border border-transparent rounded-lg shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-colors"
            >
              Sign in with Cognito
            </button>
          </div>
        </div>
        <Toaster position="bottom-right" />
      </div>
    );
  }

  return <>{children}</>;
}
