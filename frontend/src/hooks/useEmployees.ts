import { useEffect, useMemo, useReducer, useState } from "react";
import toast from "react-hot-toast";
import {
  getEmployees,
  createEmployee,
  updateEmployee,
  deleteEmployee as deleteEmployeeApi,
} from "../api/employees";
import type { Employee, EmployeeFormData } from "../types/employee";

type FetchState = {
  employees: Employee[];
  loading: boolean;
  error: string;
};

type FetchAction =
  | { type: "FETCH_START" }
  | { type: "FETCH_SUCCESS"; data: Employee[] }
  | { type: "FETCH_ERROR"; error: string };

function fetchReducer(state: FetchState, action: FetchAction): FetchState {
  switch (action.type) {
    case "FETCH_START":
      return { ...state, loading: true, error: "" };
    case "FETCH_SUCCESS":
      return { employees: action.data, loading: false, error: "" };
    case "FETCH_ERROR":
      return { ...state, loading: false, error: action.error };
  }
}

// Custom hook to manage employee state, fetching, searching, and CRUD operations
export function useEmployees(isAuthenticated: boolean) {
  const [state, dispatch] = useReducer(fetchReducer, {
    employees: [],
    loading: true,
    error: "",
  });
  const [search, setSearch] = useState("");
  const [saving, setSaving] = useState(false);
  const [deleting, setDeleting] = useState(false);
  const [fetchKey, refetch] = useReducer((c: number) => c + 1, 0);

  useEffect(() => {
    let ignore = false;
    if (!isAuthenticated) return;

    dispatch({ type: "FETCH_START" });
    getEmployees()
      .then((data) => {
        if (!ignore) dispatch({ type: "FETCH_SUCCESS", data: data ?? [] });
      })
      .catch((err) => {
        if (!ignore) {
          const errMsg = err?.response?.data?.message || err.message || "Failed to load employees";
          dispatch({ type: "FETCH_ERROR", error: errMsg });
          toast.error(errMsg);
        }
      });

    return () => {
      ignore = true;
    };
  }, [fetchKey, isAuthenticated]);

  const filtered = useMemo(() => {
    if (!search.trim()) return state.employees;
    const q = search.toLowerCase();
    return state.employees.filter(
      (e) =>
        e.name.toLowerCase().includes(q) ||
        e.email.toLowerCase().includes(q) ||
        e.department.toLowerCase().includes(q)
    );
  }, [state.employees, search]);

  const handleCreate = async (data: EmployeeFormData) => {
    setSaving(true);
    try {
      await createEmployee(data);
      refetch();
      return true;
    } catch {
      return false;
    } finally {
      setSaving(false);
    }
  };

  const handleUpdate = async (id: string, data: EmployeeFormData) => {
    setSaving(true);
    try {
      await updateEmployee(id, data);
      refetch();
      return true;
    } catch {
      return false;
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async (id: string) => {
    setDeleting(true);
    try {
      await deleteEmployeeApi(id);
      refetch();
      return true;
    } catch {
      return false;
    } finally {
      setDeleting(false);
    }
  };

  return {
    employees: state.employees,
    filtered,
    loading: state.loading,
    error: state.error,
    search,
    setSearch,
    saving,
    deleting,
    refetch,
    handleCreate,
    handleUpdate,
    handleDelete,
  };
}

