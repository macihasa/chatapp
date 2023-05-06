import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';
import './index.css';
import { Auth0Provider } from '@auth0/auth0-react';

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <Auth0Provider
      domain="dev-klg37k4khu3qm746.us.auth0.com"
      clientId="NQaqIyHFLh9wgrYVIxl3Kr3NLuLeRN1C"
      authorizationParams={{
        redirect_uri: window.location.origin,
        audience: 'http://localhost:5000/',
        scope: 'openid profile email',
      }}
    >
      <App />
    </Auth0Provider>
  </React.StrictMode>
);
