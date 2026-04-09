import React, { useState } from "react";
import styles from "./AuthPage.module.css";
import logoDarkHalf from "../assets/logoDarkHalf.png";
import logoBrightFull from "../assets/logoBrightFull.png";


type Mode = "login" | "signup";
type PageProps = {
  onNavigate : (v: "authPage" | "homePage") => void;
}

export default function AuthPage({onNavigate}:PageProps) {
  const [mode, setMode] = useState<Mode>("login");
  const [email, setEmail] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [repeatPassword, setReapeatPassword] = useState<string>("");
  const [username, setUsername] = useState<string>("");
  const [firstName, setFirstName] = useState<string>("");
  const [lastName, setLastName] = useState<string>("");
  const [error, setError] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(false);

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
    const url = "http://localhost:8080";
    const endpoint = (mode === "login") ? "/api/users/login" : "/api/users/signup";
    const body = (mode === "login") ? 
                {username, password} : 
                {firstName, lastName, username, email, password};
    try {
      const res = await fetch(url+endpoint, {
        method : "POST",
        headers : { "Content-Type": "application/json" },
        body: JSON.stringify(body),
      });

      if(!res.ok){
        const msg = await res.text();
        setError(msg);
        return;
      }

      //WE CAN CREATE A COOKIE SESSION

      onNavigate("homePage");
      console.log("Success");
    }catch {
      setError("Network error, please try again!");
    }finally {
      setLoading(false);
    }
  };


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