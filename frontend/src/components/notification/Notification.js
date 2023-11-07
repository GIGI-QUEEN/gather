import { useEffect, useContext } from 'react';
import axios from 'axios';
import { UserContext } from '../utils/UserContext';
import { toast } from 'react-toastify';
import { options } from '../utils/toast/options';

const Notification = () => {
  const { ws } = useContext(UserContext);
  const me = useContext(UserContext);
  useEffect(() => {
    const messageHandler = (msg) => {
      const wsMessage = JSON.parse(msg.data);

      axios
        .get('http://localhost:8080/chat/', { withCredentials: true })
        .then((res) => {
          if (
            res.data[0]?.receiverid === me.user.id &&
            res.data[0].content === wsMessage.body
          ) {
            toast(`${wsMessage.SenderUserName}: ${wsMessage.body}`, options);
          }
        });
    };
    ws?.addEventListener('message', messageHandler);
    return () => {
      ws?.removeEventListener('message', messageHandler);
    };
  }, [ws, me.user.id]);
};

export default Notification;
