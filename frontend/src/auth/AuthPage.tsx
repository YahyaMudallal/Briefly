import React, { useState } from "react";
import styles from "./AuthPage.module.css";
import logoDarkHalf from "../assets/logoDarkHalf.png";
import logoBrightFull from "../assets/logoBrightFull.png";


type Mode = "login" | "signup";

export default function AuthPage() {
  const [mode, setMode] = useState<Mode>("login");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [repeatPassword, setReapeatPassword] = useState("");
  const [username, setUsername] = useState("");
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  // Optional: You can keep the window resize listener if you have specific JS logic, 
  // but the CSS handles the responsive layout automatically now.

  const toggle = () => {
    setMode(m => (m === "login" ? "signup" : "login"));
    setError("");
  };

  const validate = (): string => {
    if (mode === "signup" && firstName == "" )
      return "First name cannot be empty."
    if (mode === "signup" && lastName == "" )
      return "Last name cannot be empty."
    if (username.trim().length < 3)
      return "Username must be at least 3 characters.";
    if (mode ==="signup" && (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)))
      return "Please enter a valid email address.";
    if (mode=="signup" && (password.length < 4))
      return "Password must be at least 4 characters.";
    if (mode ==="signup" && password != repeatPassword)
      return "Passwords do not match."
    
    return "";
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const err = validate();
    if (err) { setError(err); return; }
    
    setLoading(true);
    setError("");
    
    // Replace with your real API call
    await new Promise(r => setTimeout(r, 900));
    setLoading(false);
  };

  // Modern geometric logo SVG
  const LogoMark = () => (
    <svg width="32" height="32" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg" className={styles.logoIcon}>
      <rect x="4" y="4" width="10" height="10" rx="3" fill="currentColor" opacity="0.8"/>
      <rect x="18" y="4" width="10" height="24" rx="3" fill="currentColor"/>
      <rect x="4" y="18" width="10" height="10" rx="3" fill="currentColor" opacity="0.4"/>
    </svg>
  );

  return (
    <div className={styles.page}>
      {/* Left panel — hidden on mobile via CSS */}
      <div className={styles.left}>
        <div className={styles.leftContent}>
          <div className={styles.logoRow}>
            {/*<LogoMark />*/}
            <img className={styles.logoIcon} src={logoDarkHalf} />
            <div className={styles.logoText}>
              <span className={styles.logoTitle}>Briefly</span>
              <span className={styles.logoSubTitle}>Smart News Community</span>
            </div>
          </div>
          <div className={styles.tagline}>
            <h1 className={styles.taglineMain}>Stay informed.<br />Stay ahead.</h1>
            <p className={styles.taglineSub}>
              Top headlines, AI-powered summaries, and community credibility scores — all in one feed.
            </p>
          </div>
        </div>
      </div>

      {/* Right panel — form */}
      <div className={styles.right}>
        <div className={styles.formWrapper}>
          
          {/* Mobile Header */}
          <div className={styles.mobileHeader}>
            <img className={styles.logoIcon} src={logoBrightFull} />
          </div>

          <div className={styles.formHeader}>
            <h2 className={styles.formTitle}>
              {mode === "login" ? "Welcome back" : "Create an account"}
            </h2>
            <p className={styles.formSubTitle}>
              {mode === "login" ? "Don't have an account? " : "Already have an account? "}
              <button type="button" className={styles.formSubLink} onClick={toggle}>
                {mode === "login" ? "Sign up for free" : "Log in"}
              </button>
            </p>
          </div>

          <form className={styles.form} onSubmit={handleSubmit} noValidate>
            
            {mode === "signup" && (
              <div className={styles.nameRow}>
                <div className={styles.halfField}>
                  <label className={styles.label}>First name</label>
                  <input
                    className={styles.input}
                    type="text"
                    value={firstName}
                    onChange={e => setFirstName(e.target.value)}
                    placeholder="Jean"
                  />
                </div>
                <div className={styles.halfField}>
                  <label className={styles.label}>Last name</label>
                  <input
                    className={styles.input}
                    type="text"
                    value={lastName}
                    onChange={e => setLastName(e.target.value)}
                    placeholder="Paul"
                  />
                </div>
              </div>
            )}

            {mode === "signup" && (
              <Field
                label="Email"
                type="email"
                value={email}
                onChange={setEmail}
                placeholder="name@example.com"
              />
            )}

            <Field
              label="Username"
              type="text"
              value={username}
              onChange={setUsername}
              placeholder="jean_paul"
            />

            <Field
              label="Password"
              type="password"
              value={password}
              onChange={setPassword}
              placeholder="••••••••"
            />


            {mode == "signup" &&(
              <Field
                label="Repeat password"
                type="password"
                value={repeatPassword}
                onChange={setReapeatPassword}
                placeholder="••••••••"
              />
            )}

            {error && (
              <div className={styles.errorBox}>
                <span className={styles.errorIcon}>!</span>
                <p className={styles.errorText}>{error}</p>
              </div>
            )}

            <button
              type="submit"
              className={styles.submitBtn}
              disabled={loading}
            >
              {loading
                ? "Please wait…"
                : mode === "login" ? "Sign in" : "Create account"}
            </button>
          </form>
        </div>
      </div>
    </div>
  );
}

/* ── Reusable field ────────────────────────────────────────── */
interface FieldProps {
  label: string;
  type: string;
  value: string;
  onChange: (v: string) => void;
  placeholder: string;
}

function Field({ label, type, value, onChange, placeholder }: FieldProps) {
  return (
    <div className={styles.field}>
      <label className={styles.label}>{label}</label>
      <input
        className={styles.input}
        type={type}
        value={value}
        onChange={e => onChange(e.target.value)}
        placeholder={placeholder}
      />
    </div>
  );
}