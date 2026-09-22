import type { Employee } from "../types/employee";

type Props = {
  employee: Employee;
  onClose: () => void;
};

export default function EmployeeView({ employee, onClose }: Props) {
  const fields: [string, string][] = [
    ["Employee ID", employee.empID],
    ["Name", employee.name],
    ["Email", employee.email],
    ["Department", employee.department || "—"],
    ["Salary", `$${employee.salary.toLocaleString()}`],
    ["DB ID", employee._id],
  ];

  return (
    <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-xl w-full max-w-md shadow-xl">
        <div className="flex items-center justify-between px-6 py-4 border-b border-gray-100">
          <h2 className="text-lg font-semibold text-gray-900">
            Employee Details
          </h2>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600 transition-colors"
            aria-label="Close"
          >
            <svg
              className="w-5 h-5"
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
        </div>

        <dl className="px-6 py-4 divide-y divide-gray-100">
          {fields.map(([label, value]) => (
            <div
              key={label}
              className="flex items-center justify-between py-3 first:pt-0 last:pb-0"
            >
              <dt className="text-sm font-medium text-gray-500">{label}</dt>
              <dd
                className={`text-sm text-right ${
                  label === "DB ID" || label === "Employee ID"
                    ? "font-mono text-xs text-gray-400"
                    : "text-gray-900"
                }`}
              >
                {value}
              </dd>
            </div>
          ))}
        </dl>

        <div className="flex justify-end px-6 py-4 border-t border-gray-100">
          <button
            onClick={onClose}
            className="px-4 py-2.5 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
          >
            Close
          </button>
        </div>
      </div>
    </div>
  );
}
