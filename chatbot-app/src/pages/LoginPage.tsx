import useAuth from "../auth/useAuth";
import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import Tab from "../components/ui/Tab";
import TabPanel from "../components/ui/TabPanel";

const LoginPage: React.FC = () => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const { currentUser, login } = useAuth();
  const [activeTab, setActiveTab] = useState("login");

  const navigate = useNavigate();

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
    setError("");

    const formData = new FormData(event.target as HTMLFormElement);
    try {
      await login(formData);
    } catch (error) {
      if (error instanceof Error) {
        setError(error.message);
      } else {
        setError("Unknown error occurred");
      }
    }
  };

  useEffect(() => {
    console.log("login check user", currentUser);
    if (currentUser) {
      console.log("currentUser updated:", currentUser);
      navigate("/Dashboard");
    }
  }, [currentUser, navigate]); // Runs when currentUser changes

  return (
    <div className="flex flex-col w-fit h-full p-4 bg-gray-900 m-auto rounded-lg">
      <div className="flex gap-1 text-xl font-bold mb-4">
        <Tab
          label="Login"
          isActive={activeTab === "login"}
          onClick={() => setActiveTab("login")}
        />
        <Tab
          label="Register"
          isActive={activeTab === "register"}
          onClick={() => setActiveTab("register")}
        />
      </div>
      <div className="border-b-2 border-gray-700"></div>
      <div>
        <TabPanel activeTab={activeTab} tabKey="login">
          <form onSubmit={handleSubmit} className="flex flex-col space-y-4">
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
              className="bg-blue-800 text-white p-2 rounded"
            >
              Login
            </button>
          </form>
        </TabPanel>
        <TabPanel activeTab={activeTab} tabKey="register">
          <form onSubmit={handleSubmit} className="flex flex-col space-y-4">
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
              className="bg-blue-800 text-white p-2 rounded"
            >
              Register
            </button>
          </form>
        </TabPanel>
      </div>

      {error && <p className="text-red-500 mt-2">{error}</p>}
    </div>
  );
};

export default LoginPage;
