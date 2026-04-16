import React, { useState } from "react";
import { useUser } from "../context/UserContext";
import { Link, useNavigate } from 'react-router-dom';
import type { User } from '../types/types'; 

import styles from "../css/AuthPage.module.css";
import logoDarkHalf from "../assets/logoDarkHalf.png";
import logoBrightFull from "../assets/logoBrightFull.png";


type Mode = "login" | "signup";

export default function AuthPage() {
  const [mode, setMode] = useState<Mode>("login");
  const [email, setEmail] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [repeatPassword, setReapeatPassword] = useState<string>("");
  const [firstName, setFirstName] = useState<string>("");
  const [lastName, setLastName] = useState<string>("");
  const [error, setError] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(false);

  const { login } = useUser();
  const navigate = useNavigate();

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
    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email))
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
    const endpoint = (mode === "login") ? "/v1/users/login" : "/v1/users";
    const body = (mode === "login") ? 
                {email, password} : 
                {firstName, lastName, email, password};
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

      const data = await res.json(); // <-- this is the actual response body

      //construct the User object based on the expected response structure
      //this user is shared across the app via context
      const user: User = {
        id: data.user.id,
        email: data.user.email,
        firstName: data.user.firstName,
        lastName: data.user.lastName,
        initials: data.user.firstName.charAt(0).toUpperCase() + data.user.lastName.charAt(0).toUpperCase(),
        isAdmin: data.user.isAdmin
      };
      
      // Save user and token in context and localStorage
      login(user, data.token);

      // Redirect to home page after successful login/signup
      navigate("/");

      console.log("login success");
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
            {/*<Logo />*/}
            <Link to="/" className={styles.logoRow}>
              <img className={styles.logoIcon} src={logoDarkHalf} />
              <div className={styles.logoText}>
                <span className={styles.logoTitle}>Briefly</span>
                <span className={styles.logoSubTitle}>Smart News Community</span>
              </div>
            </Link>
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

            <Field
              label="Email"
              type="email"
              value={email}
              onChange={setEmail}
              placeholder="name@example.com"
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