import {
  createUserWithEmailAndPassword,
  getAuth,
  GithubAuthProvider,
  GoogleAuthProvider,
  onAuthStateChanged,
  signInWithEmailAndPassword,
  signInWithPopup,
  signOut,
  updateProfile,
  UserCredential
} from 'firebase/auth';
import React, { createContext, useContext, useEffect, useState } from 'react';
import '../firebase';
import { useRouter } from 'next/navigation';

type AuthContextType = {
  currentUser: any;
  signup: (email: string, password: string, username: string) => Promise<void>;
  login: (email: string, password: string) => Promise<UserCredential>;
  logout: () => Promise<void>;
  loginWithGithub: (provider: any) => Promise<UserCredential>;
  loginWithGoogle: () => Promise<UserCredential>;
};

const AuthContext = createContext({} as AuthContextType);
export const useAuthContext = () => {
  return useContext(AuthContext);
};

export const AuthProvider = ({ children }) => {
  const [loading, setLoading] = useState(true);
  const [currentUser, setCurrentUser] = useState();
  const router = useRouter();

  useEffect(() => {
    const auth = getAuth();
    const unsubscribe = onAuthStateChanged(auth, (user) => {
      setCurrentUser(user as any);
      setLoading(false);
    });

    return unsubscribe;
  }, []);

  // signup function
  const signup = async (email: string, password: string, username: string) => {
    const auth = getAuth();
    await createUserWithEmailAndPassword(auth, email, password);

    // update profile
    await updateProfile(auth.currentUser, {
      displayName: username
    });

    const user = auth.currentUser;
    setCurrentUser({
      ...user
    } as any);
  };

  // login function
  const login = async (email: string, password: string) => {
    const auth = getAuth();
    return signInWithEmailAndPassword(auth, email, password);
  };

  // logout function
  const logout = async () => {
    const auth = getAuth();
    await signOut(auth);
    router.push('/login');
  };
  const loginWithGithub = async (provider: any) => {
    const auth = getAuth();
    return signInWithPopup(auth, new GithubAuthProvider());
  };

  const loginWithGoogle = async () => {
    const auth = getAuth();
    return signInWithPopup(auth, new GoogleAuthProvider());
  };

  const value: AuthContextType = {
    currentUser,
    signup,
    login,
    logout,
    loginWithGithub,
    loginWithGoogle
  };

  return <AuthContext.Provider value={value}>{!loading && children}</AuthContext.Provider>;
};
