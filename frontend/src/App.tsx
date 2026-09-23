import { useState } from "react";
import toast, { Toaster } from "react-hot-toast";
import { useAuth } from "react-oidc-context";
import type { Employee, EmployeeFormData } from "./types/employee";
import EmployeeTable from "./components/EmployeeTable";
import EmployeeForm from "./components/EmployeeForm";
import EmployeeView from "./components/EmployeeView";
import DeleteConfirm from "./components/DeleteConfirm";
import AuthGuard from "./components/AuthGuard";
import Navbar from "./components/Navbar";
import { useAuthInterceptor } from "./hooks/useAuthInterceptor";
import { useEmployees } from "./hooks/useEmployees";
import { useProfile } from "./hooks/useProfile";

type Modal =
  | { type: "add" }
  | { type: "edit"; employee: Employee }
  | { type: "view"; employee: Employee }
  | { type: "delete"; employee: Employee }
  | null;

function AppContent() {
  const auth = useAuth();
  useAuthInterceptor(auth);

  const { profile, loading: profileLoading } = useProfile(auth.isAuthenticated);

  const {
    employees,
    filtered,
    loading,
    error,
    search,
    setSearch,
    saving,
    deleting,
    refetch,
    handleCreate,
    handleUpdate,
    handleDelete,
  } = useEmployees(auth.isAuthenticated);

  const [modal, setModal] = useState<Modal>(null);

  const onCreateSubmit = async (data: EmployeeFormData) => {
    const ok = await handleCreate(data);
    if (ok) setModal(null);
  };

  const onUpdateSubmit = async (id: string, data: EmployeeFormData) => {
    const ok = await handleUpdate(id, data);
    if (ok) setModal(null);
  };

  const onDeleteConfirm = async (id: string) => {
    const ok = await handleDelete(id);
    if (ok) setModal(null);
  };

  // Sign out user and redirect through Cognito logout endpoint
  const handleLogout = async () => {
    try {
      await auth.removeUser();
      const clientId = import.meta.env.VITE_COGNITO_CLIENT_ID;
      const logoutUri = import.meta.env.VITE_COGNITO_REDIRECT_URI;
      const cognitoDomain = import.meta.env.VITE_COGNITO_DOMAIN;
      window.location.href = `${cognitoDomain}/logout?client_id=${clientId}&logout_uri=${encodeURIComponent(logoutUri)}`;
    } catch {
      toast.error("Failed to clear session");
    }
  };

  return (
    <div className="min-h-screen">
      <Navbar
        totalEmployees={employees.length}
        userName={profile?.name}
        userEmail={profile?.email}
        profileLoading={profileLoading}
        showCount={!loading && !error}
        onAddClick={() => setModal({ type: "add" })}
        onLogoutClick={handleLogout}
      />

      <main className="max-w-5xl mx-auto px-4 sm:px-6 py-6">
        {/* Search filter input */}
        <div className="relative mb-5">
          <svg
            className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400 pointer-events-none"
            fill="none"
            stroke="currentColor"
            strokeWidth={2}
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="m21 21-5.197-5.197m0 0A7.5 7.5 0 1 0 5.196 5.196a7.5 7.5 0 0 0 10.607 10.607Z"
            />
          </svg>
          <input
            type="text"
            placeholder="Search by name, email, or department..."
            className="w-full border border-gray-300 rounded-lg pl-10 pr-4 py-2.5 text-sm bg-white placeholder:text-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-shadow"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
          {search && (
            <button
              onClick={() => setSearch("")}
              className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
              aria-label="Clear search"
            >
              <svg
                className="w-4 h-4"
                fill="none"
                stroke="currentColor"
                strokeWidth={2}
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M6 18 18 6M6 6l12 12"
                />
              </svg>
            </button>
          )}
        </div>

        {/* Employee directory view */}
        {loading ? (
          <div className="flex flex-col items-center justify-center py-16">
            <div className="w-8 h-8 border-2 border-gray-200 border-t-blue-600 rounded-full animate-spin mb-3" />
            <p className="text-gray-500 text-sm">Loading employees...</p>
          </div>
        ) : error ? (
          <div className="flex flex-col items-center justify-center py-16 bg-white rounded-lg border border-gray-200">
            <svg
              className="w-10 h-10 text-red-400 mb-3"
              fill="none"
              stroke="currentColor"
              strokeWidth={1.5}
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126ZM12 15.75h.007v.008H12v-.008Z"
              />
            </svg>
            <p className="text-gray-700 font-medium mb-1">Something went wrong</p>
            <p className="text-red-500 text-sm mb-4">{error}</p>
            <button
              onClick={refetch}
              className="inline-flex items-center gap-1.5 px-4 py-2 text-sm font-medium text-blue-600 border border-blue-300 rounded-lg hover:bg-blue-50 transition-colors"
            >
              <svg
                className="w-4 h-4"
                fill="none"
                stroke="currentColor"
                strokeWidth={2}
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0 3.181 3.183a8.25 8.25 0 0 0 13.803-3.7M4.031 9.865a8.25 8.25 0 0 1 13.803-3.7l3.181 3.182"
                />
              </svg>
              Retry
            </button>
          </div>
        ) : employees.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-16 bg-white rounded-lg border border-gray-200">
            <svg
              className="w-12 h-12 text-gray-300 mb-3"
              fill="none"
              stroke="currentColor"
              strokeWidth={1.5}
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M15 19.128a9.38 9.38 0 0 0 2.625.372 9.337 9.337 0 0 0 4.121-.952 4.125 4.125 0 0 0-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.106A12.318 12.318 0 0 1 8.624 21c-2.331 0-4.512-.645-6.374-1.766l-.001-.109a6.375 6.375 0 0 1 11.964-3.07M12 6.375a3.375 3.375 0 1 1-6.75 0 3.375 3.375 0 0 1 6.75 0Zm8.25 2.25a2.625 2.625 0 1 1-5.25 0 2.625 2.625 0 0 1 5.25 0Z"
              />
            </svg>
            <p className="text-gray-700 font-medium mb-1">No employees yet</p>
            <p className="text-gray-400 text-sm mb-4">Add your first employee to get started.</p>
            <button
              onClick={() => setModal({ type: "add" })}
              className="inline-flex items-center gap-1.5 px-4 py-2 text-sm font-medium text-blue-600 border border-blue-300 rounded-lg hover:bg-blue-50 transition-colors"
            >
              <svg
                className="w-4 h-4"
                fill="none"
                stroke="currentColor"
                strokeWidth={2}
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M12 4.5v15m7.5-7.5h-15"
                />
              </svg>
              Add Employee
            </button>
          </div>
        ) : filtered.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-16 bg-white rounded-lg border border-gray-200">
            <svg
              className="w-10 h-10 text-gray-300 mb-3"
              fill="none"
              stroke="currentColor"
              strokeWidth={1.5}
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="m21 21-5.197-5.197m0 0A7.5 7.5 0 1 0 5.196 5.196a7.5 7.5 0 0 0 10.607 10.607Z"
              />
            </svg>
            <p className="text-gray-700 font-medium mb-1">No results found</p>
            <p className="text-gray-400 text-sm">
              No employees match &ldquo;{search}&rdquo;. Try a different search term.
            </p>
          </div>
        ) : (
          <div className="bg-white rounded-lg border border-gray-200 overflow-hidden shadow-sm">
            <EmployeeTable
              employees={filtered}
              onView={(emp) => setModal({ type: "view", employee: emp })}
              onEdit={(emp) => setModal({ type: "edit", employee: emp })}
              onDelete={(emp) => setModal({ type: "delete", employee: emp })}
            />
          </div>
        )}
      </main>

      {/* Modal dialogs */}
      {modal?.type === "add" && (
        <EmployeeForm
          title="Add Employee"
          onSubmit={onCreateSubmit}
          onCancel={() => setModal(null)}
          saving={saving}
        />
      )}

      {modal?.type === "edit" && (
        <EmployeeForm
          title="Edit Employee"
          initialData={{
            empID: modal.employee.empID,
            name: modal.employee.name,
            email: modal.employee.email,
            department: modal.employee.department,
            salary: modal.employee.salary,
          }}
          onSubmit={(data) => onUpdateSubmit(modal.employee._id, data)}
          onCancel={() => setModal(null)}
          saving={saving}
        />
      )}

      {modal?.type === "view" && (
        <EmployeeView employee={modal.employee} onClose={() => setModal(null)} />
      )}

      {modal?.type === "delete" && (
        <DeleteConfirm
          name={modal.employee.name}
          onConfirm={() => onDeleteConfirm(modal.employee._id)}
          onCancel={() => setModal(null)}
          deleting={deleting}
        />
      )}

      <Toaster position="bottom-right" />
    </div>
  );
}

export default function App() {
  return (
    <AuthGuard>
      <AppContent />
    </AuthGuard>
  );
}
