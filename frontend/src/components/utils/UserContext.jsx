import { createContext, useEffect, useState } from 'react';
import { myAxios } from '../../api/axios';
export const UserContext = createContext();

export const UserContextProvider = (props) => {
  const [user, setUser] = useState({});
  const [error, setError] = useState(null);
  const [logged, setLogged] = useState(false);
  const [isSignUp, setIsSignUp] = useState(false);

  useEffect(() => {
    myAxios
      .get('/me')
      .then((res) => {
        setLogged(true);
        setUser(res.data);
      })
      .catch((err) => {
        setLogged(false);
        setError(err);
      });
  }, [logged]);

  const value = { user, error, logged, setLogged, isSignUp, setIsSignUp };

  return (
    <UserContext.Provider value={value}>{props.children}</UserContext.Provider>
  );
};
