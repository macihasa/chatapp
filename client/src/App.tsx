import './App.css';
import Header from './Components/Header/Header';
import { useState } from 'react';

let socket = new WebSocket('ws://localhost:5000/ws');

function App() {
  const [messages, Setmessages] = useState(Array<String>);
  const [messageToSend, SetmessageToSend] = useState(String);

  // Socket states
  socket.onopen = (event) => {
    socket.send('Connection opened.. ' + Date.now().toString());
  };

  socket.onmessage = (event: MessageEvent<any>) => {
    Setmessages((prev) => [...prev, event.data]);
  };

  // Handlers
  function handleKeyPress(event: React.KeyboardEvent<HTMLInputElement>) {
    if (event.key === 'Enter') {
      socket.send(messageToSend);
      SetmessageToSend('');
    }
  }
  function handleChange(e: React.FormEvent<HTMLInputElement>) {
    SetmessageToSend(e.currentTarget.value);
  }

  // Main function/JSX
  return (
    <div className="App">
      <Header />
      <div className="AppContainer">
        <h2>Work in progress...</h2>
        <input
          type="text"
          value={messageToSend}
          onChange={handleChange}
          onKeyDown={handleKeyPress}
        />
        {messages ? (
          <div>
            {messages.map((current: String, index: number) => (
              <div key={index}>{current}</div>
            ))}
          </div>
        ) : null}
      </div>
    </div>
  );
}

export default App;
