import { useState } from "react";
import toast from "react-hot-toast";
import { signin } from "../../api/auth";
import { useAuth } from "../../context/AuthContext";

type Props = {
  onGoToSignup: () => void;
  initialEmail?: string;
  onGoToConfirm: (email: string, password: string) => void;
};

export default function LoginForm({ onGoToSignup, initialEmail, onGoToConfirm }: Props) {
  const [email, setEmail] = useState(initialEmail || "");
  const [password, setPassword] = useState("");
  const [saving, setSaving] = useState(false);
  const { checkAuth } = useAuth();

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!email || !password) return toast.error("Please fill all fields");

    setSaving(true);
    try {
      await signin({ email, password });
      toast.success("Successfully logged in!");
      // Re-fetch profile to update context state
      await checkAuth();
    } catch (err: any) {
      const msg = err?.response?.data?.message || "";
      if (msg.toLowerCase().includes("not confirmed")) {
        toast.error("Please confirm your account first.");
        onGoToConfirm(email, password);
      } else {
        toast.error(msg || "Failed to log in");
      }
    } finally {
      setSaving(false);
    }
  }

  const inputClass =
    "w-full border border-gray-300 rounded-lg px-3 py-2.5 text-sm bg-white focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500";

  return (
    <div className="w-full max-w-sm">
      <h2 className="text-2xl font-bold text-gray-900 text-center mb-6">
        Sign In
      </h2>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1.5">
            Email
          </label>
          <input
            type="email"
            required
            className={inputClass}
            placeholder="john@example.com"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1.5">
            Password
          </label>
          <input
            type="password"
            required
            className={inputClass}
            placeholder="••••••••"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
        </div>
        <button
          type="submit"
          disabled={saving}
          className="w-full px-5 py-2.5 text-sm font-medium bg-blue-600 text-white rounded-lg hover:bg-blue-700 active:bg-blue-800 disabled:opacity-50 transition-colors shadow-sm"
        >
          {saving ? "Signing in..." : "Sign In"}
        </button>
      </form>
      <p className="text-sm text-center text-gray-600 mt-6">
        Don't have an account?{" "}
        <button
          type="button"
          onClick={onGoToSignup}
          className="text-blue-600 font-medium hover:underline"
        >
          Sign Up
        </button>
      </p>
    </div>
  );
}

