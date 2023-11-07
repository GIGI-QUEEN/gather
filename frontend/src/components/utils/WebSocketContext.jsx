import React, { createContext, useContext, useEffect, useState } from 'react';
import { toast } from 'react-toastify';
import { options } from './toast/options';

const WebSocketContext = createContext({
  webSocket: null,
});

const WebSocketContextProvider = ({ children }) => {
  const [webSocket, setWebSocket] = useState(null);

  useEffect(() => {
    const ws = new WebSocket('ws://localhost:8080/ws');

    ws.onopen = () => {};
    ws.onmessage = (msg) => {
      const eventData = JSON.parse(msg.data);
      const eventType = eventData.event_type;
      switch (eventType) {
        case 'ws_post_comment_event':
          toast.info(
            `${eventData.author_username} commented on your post:\n"${eventData.content}"`,
            options
          );
          break;
        case 'ws_msg_event':
          toast.info(`New message from ${eventData.sender_username}`, options);
          break;
        case 'ws_group_event_created_event':
          toast.info(
            `New event in one of your groups: "${eventData.group_event_title}"`,
            options
          );
          break;
        case 'ws_group_join_invite_event':
          toast.info(
            `${eventData.user_inviting} sent you invite to join group: ${eventData.group_title}`,
            options
          );
          break;
        case 'ws_group_join_request_event':
          toast.info(
            `${eventData.user_to_join} wants to join your group: ${eventData.group_title}`,
            options
          );
          break;
        case 'ws_follow_request_event':
          toast.info(`${eventData.follower_name} wants to follow you`, options);
          break;
        case 'ws_group_msg_event':
          toast.info(`New message in: ${eventData.group_title}`, options);
          break;
        default:
      }
    };
    ws.onerror = () => {};

    setWebSocket(ws);

    return () => {
      ws.close();
    };
  }, []);

  return (
    <WebSocketContext.Provider value={{ webSocket }}>
      {children}
    </WebSocketContext.Provider>
  );
};

const useWebSocket = () => useContext(WebSocketContext);

export { WebSocketContext, WebSocketContextProvider, useWebSocket };
