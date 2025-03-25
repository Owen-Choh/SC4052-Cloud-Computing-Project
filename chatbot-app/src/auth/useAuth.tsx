import { useState, useContext, createContext, useEffect } from "react";
import { LoginResponse, User } from "./auth";
import { checkAuthApi, loginApi, logoutApi } from "../api/apiConfig";

export interface AuthContextType {
  currentUser: User | null;
  login: (formData: FormData) => Promise<void>;
  isAuthenticated: boolean;
  doLogout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [currentUser, setCurrentUser] = useState<User | null>(null);
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
  // const [token, setToken] = useState<string>("");

  const login = async (formData: FormData) => {
    console.log("sending formData to loginApi:", loginApi.getUri());
    const loginResponse = await loginApi.post("", formData, {
      headers: {
        "Content-Type": "multipart/form-data",
      },
      withCredentials: true,
    });
    if (loginResponse.status === 200) {
      const logindata: LoginResponse = loginResponse.data;

      // setToken(logindata.token);
      setCurrentUser(logindata.user);
      setIsAuthenticated(true);
    }
  };

  const doLogout = () => {
    setCurrentUser(null);
    setIsAuthenticated(false);
    // setToken("");
    logoutApi.get("", { withCredentials: true });
  };

  const checkAuth = async () => {
    try {
      const response = await checkAuthApi.get("", {
        withCredentials: true,
      });
      if (response.status === 200) {
        const logindata: LoginResponse = response.data;
        setCurrentUser(logindata.user);
        setIsAuthenticated(true);
      } else {
        setIsAuthenticated(false);
        setCurrentUser(null);
        doLogout();
      }
    } catch (error) {
      setIsAuthenticated(false);
      setCurrentUser(null);
      doLogout();
    }
  };

  useEffect(() => {
    checkAuth();
  }, []);

  return (
    <AuthContext.Provider
      value={{ currentUser, login, isAuthenticated, doLogout }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};

export default useAuth;
