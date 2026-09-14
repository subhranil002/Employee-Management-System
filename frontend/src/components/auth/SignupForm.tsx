import { useState } from "react";
import toast from "react-hot-toast";
import { signup, checkStatus } from "../../api/auth";

type Props = {
  onGoToLogin: (email?: string) => void;
  onSignupSuccess: (email: string, password: string) => void;
};

export default function SignupForm({ onGoToLogin, onSignupSuccess }: Props) {
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [saving, setSaving] = useState(false);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!name || !email || !password)
      return toast.error("Please fill all fields");

    setSaving(true);
    try {
      await signup({ name, email, password });
      toast.success("Registration successful! Check your email for OTP.");
      onSignupSuccess(email, password);
    } catch (err: any) {
      const msg = err?.response?.data?.message || "";
      if (msg.toLowerCase().includes("exist")) {
        try {
          const status = await checkStatus(email);
          if (status === "CONFIRMED") {
            toast.error("Account already exists. Please sign in.");
            onGoToLogin(email);
          } else if (status === "UNCONFIRMED") {
            toast.error("Account exists but is unconfirmed. Please verify.");
            onSignupSuccess(email, password);
          } else {
            toast.error(msg);
          }
        } catch (statusErr) {
          toast.error(msg);
        }
      } else {
        toast.error(msg || "Failed to sign up");
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
        Sign Up
      </h2>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1.5">
            Full Name
          </label>
          <input
            type="text"
            required
            className={inputClass}
            placeholder="John Doe"
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
        </div>
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
          <p className="text-xs text-gray-500 mt-1">
            At least 8 chars, 1 uppercase, 1 lowercase, 1 number, 1 special
            character.
          </p>
        </div>
        <button
          type="submit"
          disabled={saving}
          className="w-full px-5 py-2.5 text-sm font-medium bg-blue-600 text-white rounded-lg hover:bg-blue-700 active:bg-blue-800 disabled:opacity-50 transition-colors shadow-sm mt-2"
        >
          {saving ? "Signing up..." : "Sign Up"}
        </button>
      </form>
      <p className="text-sm text-center text-gray-600 mt-6">
        Already have an account?{" "}
        <button
          type="button"
          onClick={() => onGoToLogin(email)}
          className="text-blue-600 font-medium hover:underline"
        >
          Sign In
        </button>
      </p>
    </div>
  );
}

