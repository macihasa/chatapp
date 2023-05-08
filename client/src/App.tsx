import './App.css';
import Header from './Components/Header/Header';
import Chat from './Components/Chat/Chat';
import { useAuth0 } from '@auth0/auth0-react';
import { useEffect, useState } from 'react';

function App() {
  const {
    // Auth state:
    isAuthenticated,
    user,
    // Auth methods:
    getAccessTokenSilently,
    loginWithRedirect,
    loginWithPopup,
    logout,
  } = useAuth0();

  const [localUser, setLocalUser] = useState<UserModel | null>(null);

  interface UserModel {
    _id: string;
    auth0id: string;
    username: string;
    email: string;
  }

  useEffect(() => {
    async function getUserInformation() {
      try {
        const token = await getAccessTokenSilently();
        let response = await fetch('http://localhost:5000/users/login', {
          headers: {
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/json',
          },
          method: 'POST',
          body: JSON.stringify(user),
        });
        const newobject = (await response.json()) as UserModel;
        setLocalUser(newobject);
      } catch (err) {
        console.error(err);
      }
    }
    if (isAuthenticated) {
      getUserInformation();
    }
    return () => {};
  }, [getAccessTokenSilently, isAuthenticated]);

  return (
    <div className="App">
      <Header />
      <button onClick={() => loginWithRedirect()}>Login</button>
      <button onClick={() => logout()}>logout</button>
      <div className="AppContainer">
        {isAuthenticated ? (
          <div>
            <Chat />
            {localUser ? <h2>{localUser?.auth0id}</h2> : ''}
          </div>
        ) : null}
      </div>
    </div>
  );
}

export default App;
