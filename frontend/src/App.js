import 'react-toastify/dist/ReactToastify.css';
import AllRoutes from './components/Routes/Routes';
import { Header } from './components/Header/Header';
import './styles/index.scss';
import Notification from './components/notification/Notification';
import { ToastContainer } from 'react-toastify';
import { useContext } from 'react';
import { UserContext } from './components/utils/UserContext';
import { AuthForms } from './pages/login/Login';
import { WebSocketContextProvider } from './components/utils/WebSocketContext';

function App() {
  const { logged } = useContext(UserContext);
  if (!logged) return <AuthForms />;
  return (
    <div className="main-container">
      <WebSocketContextProvider>
        <Header />
        <Notification />
        <AllRoutes />
        <ToastContainer />
      </WebSocketContextProvider>
    </div>
  );
}

export default App;
