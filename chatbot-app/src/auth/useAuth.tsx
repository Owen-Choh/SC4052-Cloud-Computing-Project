import { useState, useContext, createContext } from "react";
import { LoginResponse, User } from "./auth";
import { loginApi } from "../api/apiConfig";

export interface AuthContextType {
  currentUser: User | null;
  login: (formData: FormData) => Promise<void>;
  isAuthenticated: boolean;
  doLogout: () => void;
  token: string;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ 
  children 
}) => {
  const [currentUser, setCurrentUser] = useState<User | null>(null);
  const [token, setToken] = useState<string>("");
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);

  const login = async (formData: FormData) => {
      console.log("sending formData to loginApi:", loginApi.getUri());
      const loginResponse = await loginApi.post("", formData, {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      });
      if (loginResponse.status === 200) {
        const logindata: LoginResponse = loginResponse.data;
        
        setToken(logindata.token);
        setCurrentUser(logindata.user);
        setIsAuthenticated(true);
      } 
    
  }

  const doLogout = () => {
    setCurrentUser(null);
    setToken("");
    setIsAuthenticated(false);
  }

  return <AuthContext.Provider
    value={{ currentUser, login, isAuthenticated, doLogout, token }}
  >
    {children}
  </AuthContext.Provider>
}

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};

export default useAuth;
