import type { Employee } from "../types/employee";

type Props = {
  employees: Employee[];
  onView: (emp: Employee) => void;
  onEdit: (emp: Employee) => void;
  onDelete: (emp: Employee) => void;
};

export default function EmployeeTable({
  employees,
  onView,
  onEdit,
  onDelete,
}: Props) {
  return (
    <div className="overflow-x-auto">
      <table className="w-full text-sm text-left">
        <thead>
          <tr className="border-b border-gray-200 bg-gray-50 text-gray-500 uppercase text-xs tracking-wide">
            <th className="px-5 py-3 font-semibold">Emp ID</th>
            <th className="px-5 py-3 font-semibold">Name</th>
            <th className="px-5 py-3 font-semibold">Email</th>
            <th className="px-5 py-3 font-semibold">Department</th>
            <th className="px-5 py-3 font-semibold text-right">Salary</th>
            <th className="px-5 py-3 font-semibold text-right">Actions</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-gray-100">
          {employees.map((emp) => (
            <tr
              key={emp._id}
              className="hover:bg-gray-50/80 transition-colors"
            >
              <td className="px-5 py-3.5 font-mono text-xs text-gray-500">
                {emp.empID}
              </td>
              <td className="px-5 py-3.5 font-medium text-gray-900">
                {emp.name}
              </td>
              <td className="px-5 py-3.5 text-gray-600">{emp.email}</td>
              <td className="px-5 py-3.5">
                <span className="inline-block bg-gray-100 text-gray-700 text-xs font-medium px-2.5 py-1 rounded-full">
                  {emp.department}
                </span>
              </td>
              <td className="px-5 py-3.5 text-right text-gray-900 tabular-nums">
                ${emp.salary.toLocaleString()}
              </td>
              <td className="px-5 py-3.5 text-right">
                <div className="inline-flex items-center gap-1">
                  <button
                    onClick={() => onView(emp)}
                    className="px-2.5 py-1.5 text-xs font-medium text-gray-600 rounded-md hover:bg-gray-100 transition-colors"
                  >
                    View
                  </button>
                  <button
                    onClick={() => onEdit(emp)}
                    className="px-2.5 py-1.5 text-xs font-medium text-blue-600 rounded-md hover:bg-blue-50 transition-colors"
                  >
                    Edit
                  </button>
                  <button
                    onClick={() => onDelete(emp)}
                    className="px-2.5 py-1.5 text-xs font-medium text-red-600 rounded-md hover:bg-red-50 transition-colors"
                  >
                    Delete
                  </button>
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
