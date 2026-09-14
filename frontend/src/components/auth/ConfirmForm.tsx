import { useState } from "react";
import toast from "react-hot-toast";
import { confirmSignup, signin } from "../../api/auth";
import { useAuth } from "../../context/AuthContext";

type Props = {
  email: string;
  password?: string;
  onGoToLogin: () => void;
};

export default function ConfirmForm({ email, password, onGoToLogin }: Props) {
  const [code, setCode] = useState("");
  const [saving, setSaving] = useState(false);
  const { checkAuth } = useAuth();

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!code) return toast.error("Please enter the verification code");

    setSaving(true);
    try {
      await confirmSignup({ email, code });
      
      if (password) {
        toast.success("Account confirmed! Signing in...");
        await signin({ email, password });
        await checkAuth(); // Triggers reload into dashboard
      } else {
        toast.success("Account confirmed! Please log in.");
        onGoToLogin();
      }
    } catch (err: any) {
      toast.error(err?.response?.data?.message || "Failed to confirm account");
    } finally {
      setSaving(false);
    }
  }

  const inputClass =
    "w-full border border-gray-300 rounded-lg px-3 py-2.5 text-sm bg-white focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 text-center tracking-widest font-mono text-lg";

  return (
    <div className="w-full max-w-sm">
      <h2 className="text-2xl font-bold text-gray-900 text-center mb-2">
        Verify Email
      </h2>
      <p className="text-sm text-gray-600 text-center mb-6">
        We sent a verification code to <br />
        <span className="font-medium text-gray-900">{email}</span>
      </p>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1.5 text-center">
            Verification Code
          </label>
          <input
            type="text"
            required
            className={inputClass}
            placeholder="123456"
            value={code}
            onChange={(e) => setCode(e.target.value)}
          />
        </div>
        <button
          type="submit"
          disabled={saving}
          className="w-full px-5 py-2.5 text-sm font-medium bg-blue-600 text-white rounded-lg hover:bg-blue-700 active:bg-blue-800 disabled:opacity-50 transition-colors shadow-sm"
        >
          {saving ? "Verifying..." : "Verify Account"}
        </button>
      </form>
      <button
        type="button"
        onClick={onGoToLogin}
        className="w-full mt-4 text-sm text-gray-500 hover:text-gray-700 transition-colors"
      >
        Back to Sign In
      </button>
    </div>
  );
}

