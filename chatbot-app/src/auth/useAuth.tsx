import { useState, useEffect, useContext, createContext } from "react";
import { User } from "./User";
import { loginApi } from "../api/apiConfig";
import { useNavigate } from "react-router-dom";

export interface AuthContextType {
  currentUser: User | null;
  login: (formData: FormData) => Promise<void>;
  isAuthenticated: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ 
  children 
}) => {
  const [currentUser, setCurrentUser] = useState<User | null>(null);
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);

  const navigate = useNavigate();
  const login = async (formData: FormData) => {
      console.log("login called with formData", formData);
      const oAuth2FormData = new FormData();
      for (const [key, value] of formData.entries()) {
        // console.log("key", key, "value", value);
        oAuth2FormData.append(key, value);
      }

      console.log("sending oAuth2FormData", oAuth2FormData, "to loginApi", loginApi.getUri());
      const loginResponse = await loginApi.post("/", oAuth2FormData, {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      });
      if (loginResponse.status === 200) {
        const user: User = loginResponse.data;
        setCurrentUser(user);
        setIsAuthenticated(true);
      } 
    
  }

  return <AuthContext.Provider
    value={{ currentUser, login, isAuthenticated }}
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
