interface NavbarProps {
  totalEmployees: number;
  userName?: string;
  userEmail?: string;
  profileLoading: boolean;
  showCount: boolean;
  onAddClick: () => void;
  onLogoutClick: () => void;
}

// Navigation header component
export default function Navbar({
  totalEmployees,
  userName,
  userEmail,
  profileLoading,
  showCount,
  onAddClick,
  onLogoutClick,
}: NavbarProps) {
  return (
    <header className="bg-white border-b border-gray-200">
      <div className="max-w-5xl mx-auto px-4 sm:px-6 py-5 flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Employees</h1>
          {showCount && (
            <p className="text-sm text-gray-500 mt-0.5">
              {totalEmployees} {totalEmployees === 1 ? "employee" : "employees"} total
            </p>
          )}
        </div>
        <div className="flex items-center gap-4">
          <div className="text-right hidden sm:block min-w-[100px]">
            {profileLoading ? (
              <div className="flex flex-col items-end gap-1">
                <div className="h-3.5 w-24 bg-gray-200 rounded animate-pulse" />
                <div className="h-3 w-32 bg-gray-100 rounded animate-pulse" />
              </div>
            ) : (
              <>
                <p className="text-sm font-medium text-gray-900">{userName || "—"}</p>
                <p className="text-xs text-gray-500">{userEmail || ""}</p>
              </>
            )}
          </div>
          <button
            onClick={onAddClick}
            className="inline-flex items-center gap-1.5 px-4 py-2.5 text-sm font-medium bg-blue-600 text-white rounded-lg hover:bg-blue-700 active:bg-blue-800 transition-colors shadow-sm"
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
          <button
            onClick={onLogoutClick}
            className="inline-flex items-center gap-1.5 px-4 py-2.5 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
          >
            Logout
          </button>
        </div>
      </div>
    </header>
  );
}
