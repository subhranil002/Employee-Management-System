import { useState } from "react";
import type { EmployeeFormData } from "../types/employee";

type Props = {
  initialData?: EmployeeFormData;
  onSubmit: (data: EmployeeFormData) => void;
  onCancel: () => void;
  saving: boolean;
  title: string;
};

const emptyForm: EmployeeFormData = {
  empID: "",
  name: "",
  email: "",
  department: "",
  salary: 0,
};

export default function EmployeeForm({
  initialData,
  onSubmit,
  onCancel,
  saving,
  title,
}: Props) {
  const [form, setForm] = useState<EmployeeFormData>(initialData ?? emptyForm);
  const [errors, setErrors] = useState<Record<string, string>>({});

  function validate(): boolean {
    const e: Record<string, string> = {};
    if (!form.empID.trim()) {
      e.empID = "Employee ID is required";
    } else if (!/^EMP-\d{3}$/.test(form.empID.trim())) {
      e.empID = "Employee ID must be in format EMP-123";
    }
    if (!form.name.trim()) e.name = "Name is required";
    if (!form.email.trim()) e.email = "Email is required";
    if (form.salary < 0) e.salary = "Salary must be >= 0";
    setErrors(e);
    return Object.keys(e).length === 0;
  }

  function handleSubmit(ev: React.FormEvent) {
    ev.preventDefault();
    if (validate()) onSubmit(form);
  }

  function update(field: keyof EmployeeFormData, value: string | number) {
    setForm((f) => ({ ...f, [field]: value }));
  }

  const inputClass = (field?: string) =>
    `w-full border rounded-lg px-3 py-2.5 text-sm bg-white transition-shadow focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 ${
      field && errors[field]
        ? "border-red-400 focus:ring-red-500 focus:border-red-500"
        : "border-gray-300"
    }`;

  return (
    <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-xl p-6 w-full max-w-md shadow-xl">
        <div className="flex items-center justify-between mb-5">
          <h2 className="text-lg font-semibold text-gray-900">{title}</h2>
          <button
            type="button"
            onClick={onCancel}
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
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1.5">
              Employee ID <span className="text-red-500">*</span>
            </label>
            <input
              className={`${inputClass("empID")} ${initialData ? "bg-gray-50 text-gray-500 cursor-not-allowed" : ""}`}
              placeholder="EMP-123"
              value={form.empID}
              onChange={(e) => update("empID", e.target.value)}
              disabled={!!initialData}
            />
            {errors.empID && (
              <p className="text-red-500 text-xs mt-1.5">{errors.empID}</p>
            )}
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1.5">
              Name <span className="text-red-500">*</span>
            </label>
            <input
              className={inputClass("name")}
              placeholder="John Doe"
              value={form.name}
              onChange={(e) => update("name", e.target.value)}
            />
            {errors.name && (
              <p className="text-red-500 text-xs mt-1.5">{errors.name}</p>
            )}
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1.5">
              Email <span className="text-red-500">*</span>
            </label>
            <input
              type="email"
              className={inputClass("email")}
              placeholder="john@example.com"
              value={form.email}
              onChange={(e) => update("email", e.target.value)}
            />
            {errors.email && (
              <p className="text-red-500 text-xs mt-1.5">{errors.email}</p>
            )}
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1.5">
              Department
            </label>
            <input
              className={inputClass()}
              placeholder="Engineering"
              value={form.department}
              onChange={(e) => update("department", e.target.value)}
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1.5">
              Salary
            </label>
            <div className="relative">
              <span className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 text-sm">
                $
              </span>
              <input
                type="number"
                className={`${inputClass("salary")} pl-7`}
                placeholder="0"
                value={form.salary}
                onChange={(e) => update("salary", Number(e.target.value))}
              />
            </div>
            {errors.salary && (
              <p className="text-red-500 text-xs mt-1.5">{errors.salary}</p>
            )}
          </div>
          <div className="flex justify-end gap-3 pt-3 border-t border-gray-100">
            <button
              type="button"
              onClick={onCancel}
              className="px-4 py-2.5 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
              disabled={saving}
            >
              Cancel
            </button>
            <button
              type="submit"
              className="px-5 py-2.5 text-sm font-medium bg-blue-600 text-white rounded-lg hover:bg-blue-700 active:bg-blue-800 disabled:opacity-50 transition-colors shadow-sm"
              disabled={saving}
            >
              {saving ? (
                <span className="inline-flex items-center gap-2">
                  <span className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                  Saving...
                </span>
              ) : (
                "Save"
              )}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
