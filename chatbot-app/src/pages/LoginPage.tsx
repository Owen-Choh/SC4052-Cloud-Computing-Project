import useAuth from "../auth/useAuth";
import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import Tab from "../components/ui/Tab";
import TabPanel from "../components/ui/TabPanel";
import axios from "axios";
import { registerApi } from "../api/apiConfig";

/**
 * LoginPage Component
 *
 * This component provides a user interface for both logging in and registering.
 * It utilizes tabs to switch between the login and register forms.
 */
const LoginPage: React.FC = () => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [password2, setPassword2] = useState("");
  const [error, setError] = useState("");
  const { currentUser, login } = useAuth();
  const [activeTab, setActiveTab] = useState("login");

  const navigate = useNavigate();

  /**
   * validateUsername
   * Validates the username to ensure it is alphanumeric and at least 3 characters long.
   * @param username The username to validate.
   * @returns True if the username is valid, false otherwise.
   */
  const validateUsername = (username: string) => {
    const regex = /^[a-zA-Z0-9]{3,}$/; // Alphanumeric, at least 3 characters
    return regex.test(username);
  };

  /**
   * validatePassword
   * Validates the password to ensure it is at least 8 characters long.
   * @param password The password to validate.
   * @returns True if the password is valid, false otherwise. Also sets an error message if invalid.
   */
  const validatePassword = (password: string) => {
    if (password.length < 8) {
      setError("Password must be at least 8 characters long");
      return false;
    }
    return true;
  };

  /**
   * submitLoginForm
   * Handles the submission of the login form.  It prevents the default form submission,
   * validates the username and password, and then calls the login function from the useAuth hook.
   * @param event The form event.
   */
  const submitLoginForm = async (event: React.FormEvent) => {
    event.preventDefault();
    setError("");

    if (!validateUsername(username)) {
      setError(
        "Invalid username. Only alphanumeric characters and underscores are allowed, and it must be at least 3 characters long."
      );
      return;
    }
    if (!validatePassword(password)) {
      setError(
        "Invalid password. Password must be at least 8 characters long."
      );
      return;
    }

    const formData = new FormData(event.target as HTMLFormElement);
    try {
      await login(formData);
    } catch (error) {
      if (axios.isAxiosError(error)) {
        let errormsg = error.message;
        if (error.response?.data.error) {
          errormsg = errormsg + " " + error.response?.data.error;
        }
        setError(errormsg);
      } else {
        setError("Unknown error occurred");
      }
    }
  };

  /**
   * submitRegisterForm
   * Handles the submission of the registration form. It prevents the default form submission,
   * validates the username, password, and password confirmation, and then calls the register API.
   * @param event The form event.
   */
  const submitRegisterForm = async (event: React.FormEvent) => {
    event.preventDefault();
    setError("");
    if (password !== password2) {
      setError("Passwords do not match");
      return;
    }

    if (!validateUsername(username)) {
      setError(
        "Invalid username. Only alphanumeric characters are allowed, and it must be at least 3 characters long."
      );
      return;
    }
    if (!validatePassword(password)) {
      setError(
        "Invalid password. Password must be at least 8 characters long."
      );
      return;
    }

    const formData = new FormData(event.target as HTMLFormElement);
    try {
      const response = await registerApi.post("", formData);
      if (response.status === 201) {
        await login(formData);
      }
    } catch (error) {
      if (axios.isAxiosError(error)) {
        // console.log("error", error);
        var errormsg = error.message;
        if (errormsg === "Network Error") {
          errormsg = "Network Error: Please check your connection";
        } else if (error.response?.data.error.includes("invalid payload ")) {
          errormsg =
            errormsg +
            " " +
            "There may be invalid inputs, please check your username or password";
        } else if (error.response?.data.error) {
          errormsg = errormsg + " " + error.response?.data.error;
        }
        setError(errormsg);
      } else {
        setError("Unknown error occurred");
      }
    }
  };

  /**
   * useEffect hook to handle navigation upon successful login.
   * It checks if the currentUser state is populated (i.e., the user is logged in),
   * and if so, navigates to the '/Dashboard' route.
   */
  useEffect(() => {
    // console.log("login check user", currentUser);
    if (currentUser) {
      // console.log("currentUser updated:", currentUser);
      navigate("/Dashboard");
    }
  }, [currentUser, navigate]); // Runs when currentUser changes

  return (
    <div className="flex flex-col w-full h-full p-4 bg-gray-800 m-auto items-center">
      <div className="w-1/2 p-4 bg-blue-900 text-center m-auto rounded-lg">
        <h1 className="text-5xl">Welcome to SimpleChat</h1>
        <div className="border-b-2 py-2 border-gray-700"></div>

        <p className="text-lg pb-2">Login or Register below. </p>
        <p>
          <a href="https://ec2-54-179-162-106.ap-southeast-1.compute.amazonaws.com/chat/SimpleChat/FAQ" target="_blank" rel="noopener noreferrer" className="bg-gray-200 rounded-lg p-1 font-bold underline hover:bg-gray-300">
            Click here
          </a>{" "}
          to understand how to use our app
        </p>
      </div>
      <div className="flex flex-col w-1/2 h-full p-4 bg-gray-900 m-auto rounded-lg">
        <div className="flex gap-1 text-xl font-bold p-4">
          <Tab
            label="Login"
            isActive={activeTab === "login"}
            onClick={() => {
              setActiveTab("login");
              setError("");
            }}
          />
          <Tab
            label="Register"
            isActive={activeTab === "register"}
            onClick={() => {
              setActiveTab("register");
              setError("");
            }}
          />
        </div>
        <div className="border-b-2 border-gray-700"></div>
        <div>
          <TabPanel activeTab={activeTab} tabKey="login">
            <form
              onSubmit={submitLoginForm}
              className="flex flex-col space-y-4"
            >
              <input
                type="text"
                name="username"
                placeholder="Username"
                className="border rounded p-2"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                required
              />
              <input
                type="password"
                name="password"
                placeholder="Password"
                className="border rounded p-2"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
              />
              <button
                type="submit"
                className="bg-blue-800 text-white p-2 rounded mt-2"
              >
                Login
              </button>
            </form>
          </TabPanel>
          <TabPanel activeTab={activeTab} tabKey="register">
            <form
              onSubmit={submitRegisterForm}
              className="flex flex-col space-y-4"
            >
              <input
                type="text"
                name="username"
                placeholder="Username"
                className="border rounded p-2"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                required
              />
              <input
                type="password"
                name="password"
                placeholder="Password"
                className="border rounded p-2"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
              />
              <input
                type="password"
                name="password-confirm"
                placeholder="Confirm Password"
                className="border rounded p-2"
                value={password2}
                onChange={(e) => setPassword2(e.target.value)}
                required
              />
              <button
                type="submit"
                className="bg-blue-800 text-white p-2 rounded mt-2"
              >
                Register
              </button>
            </form>
          </TabPanel>
        </div>

        {error && <p className="text-red-500 m-auto p-4">{error}</p>}
      </div>
    </div>
  );
};

export default LoginPage;