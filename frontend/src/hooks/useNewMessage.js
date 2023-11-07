import { useEffect } from 'react';
import { sendGroupChatMessagesRequest } from '../pages/ws-chat/backend-api-calls';

//puts new messages into map
export const useNewMessage = (
  webSocket,
  setConversationsMap,
  setGroupConversationMap
) => {
  useEffect(() => {
    const messageHandler = async (event) => {
      const handleNewMessageReceived = (message) => {
        setConversationsMap((prevMap) => {
          const map = prevMap ? new Map(prevMap) : new Map();
          if (map.has(message.sender)) {
            const existingMessages = map.get(message.sender);
            map.set(message.sender, [...existingMessages, message]);
          } else {
            map.set(message.sender, [message]);
          }
          return map;
        });
      };

      const handleNewGroupMessageReceived = async (message) => {
        setGroupConversationMap((prevMap) => {
          const map = prevMap ? new Map(prevMap) : new Map();
          if (map.has(message.group_id)) {
            const existingMessages = map.get(message.group_id);
            map.set(message.group_id, [...existingMessages, message]);
          } else {
            sendGroupChatMessagesRequest(message?.group_id)
              .then((res) => {
                const messages = res.data.reverse();
                map.set(message?.group_id, messages);
              })
              .catch((err) => console.log(err.message));
            map.set(message.group_id, [message]);
          }
          return map;
        });
      };

      const message = JSON.parse(event.data);
      if (message.event_type === 'ws_msg_event') {
        handleNewMessageReceived(message);
      }
      if (message.event_type === 'ws_group_msg_event') {
        handleNewGroupMessageReceived(message);
      }
    };
    webSocket?.addEventListener('message', messageHandler);
    return () => {
      webSocket?.removeEventListener('message', messageHandler);
    };
  }, [webSocket, setGroupConversationMap, setConversationsMap]);
};
