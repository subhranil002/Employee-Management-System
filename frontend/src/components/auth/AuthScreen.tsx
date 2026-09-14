import { useState } from "react";
import LoginForm from "./LoginForm";
import SignupForm from "./SignupForm";
import ConfirmForm from "./ConfirmForm";

type ViewState = "login" | "signup" | "confirm";

export default function AuthScreen() {
  const [view, setView] = useState<ViewState>("login");
  const [registeredEmail, setRegisteredEmail] = useState("");
  const [registeredPassword, setRegisteredPassword] = useState("");

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
      <div className="sm:mx-auto sm:w-full sm:max-w-md">
        <div className="bg-white py-8 px-4 shadow sm:rounded-lg sm:px-10 border border-gray-200 flex justify-center">
          {view === "login" && (
            <LoginForm 
              onGoToSignup={() => setView("signup")} 
              initialEmail={registeredEmail}
              onGoToConfirm={(email, password) => {
                setRegisteredEmail(email);
                setRegisteredPassword(password);
                setView("confirm");
              }}
            />
          )}
          {view === "signup" && (
            <SignupForm
              onGoToLogin={(email) => {
                if (email) setRegisteredEmail(email);
                setView("login");
              }}
              onSignupSuccess={(email, password) => {
                setRegisteredEmail(email);
                setRegisteredPassword(password);
                setView("confirm");
              }}
            />
          )}
          {view === "confirm" && (
            <ConfirmForm
              email={registeredEmail}
              password={registeredPassword}
              onGoToLogin={() => setView("login")}
            />
          )}
        </div>
      </div>
    </div>
  );
}

