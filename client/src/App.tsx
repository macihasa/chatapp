import './App.css';
import Header from './Components/Header/Header';
import Chat from './Components/Chat/Chat';
import { useAuth0 } from '@auth0/auth0-react';

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

  async function sendRequest() {
    let token;
    try {
      token = await getAccessTokenSilently();
    } catch (error) {
      console.log('Failed to retrieve Access Token:');
      console.error(error);
    }
    try {
      const response = await fetch('http://localhost:5000/', {
        method: 'GET',
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      console.log(await response.json());
    } catch (error) {
      console.log('Failed authentication request to server');
      console.error(error);
    }
  }

  return (
    <div className="App">
      <Header />
      <button onClick={() => loginWithRedirect()}>Login</button>
      <button onClick={() => logout()}>logout</button>
      <button onClick={sendRequest}>sendrequest</button>
      <div className="AppContainer">
        {isAuthenticated ? (
          <div>
            <Chat />
            {user ? <h2>{user.email}</h2> : ''}
          </div>
        ) : null}
      </div>
    </div>
  );
}

export default App;
