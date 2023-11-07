import React, { useEffect } from 'react';
import { Emoji } from '../../components/AttachAndEmojis';
export const Input = ({
  messageToSend,
  setMessageToSend,
  setConversationsMap,
  webSocket,
  recipient,
  user,
}) => {
  let inp = document.getElementById('input');
  const handleNewMessageSent = (message) => {
    setConversationsMap((prevMap) => {
      const map = prevMap ? new Map(prevMap) : new Map();
      if (map.has(recipient.id)) {
        const existingMessages = map.get(recipient.id);
        map.set(recipient.id, [...existingMessages, message]);
      } else {
        map.set(recipient.id, [message]);
      }
      return map;
    });
  };
  const handleKeyUp = (e) => {
    if (e.key === 'Enter') {
      handleMessageSubmit(e);
    }
  };
  useEffect(() => {
    inp?.addEventListener('keyup', handleKeyUp);
    return () => inp?.removeEventListener('keyup', handleKeyUp);
  }, [inp?.value]);

  const handleMessageSubmit = (event) => {
    event.preventDefault();
    const messageEvent = {
      event_type: 'ws_msg_event',
      sender: user.id,
      sender_username: user.username,
      recipient: recipient.id,
      message: messageToSend,
      created_date: Math.floor(Date.now() / 1000),
    };
    if (messageEvent.message !== '') {
      webSocket?.send(JSON.stringify(messageEvent));
      handleNewMessageSent(messageEvent);
      setMessageToSend('');
    }
  };

  return (
    <div className="input-box-container">
      <input
        type="text"
        onChange={(e) => setMessageToSend(e.target.value)}
        placeholder="type your message..."
        value={messageToSend}
        id="input"
      />
      <div className="attach-and-emojis">
        <Emoji setMsg={setMessageToSend} />
      </div>
      <button onClick={(e) => handleMessageSubmit(e)}>send</button>
    </div>
  );
};

// FOR GROUP
export const GroupInput = ({
  messageToSend,
  setMessageToSend,
  setGroupConversationsMap,
  webSocket,
  group,
  user,
}) => {
  let inp = document.getElementById('input');
  const handleNewGroupMessageSent = (message) => {
    setGroupConversationsMap((prevMap) => {
      const map = prevMap ? new Map(prevMap) : new Map();
      if (map.has(group?.id)) {
        const existingMessages = map.get(group?.id);
        map.set(group?.id, [...existingMessages, message]);
      } else {
        map.set(group?.id, [message]);
      }
      return map;
    });
  };
  const handleKeyUp = (e) => {
    if (e.key === 'Enter') {
      handleMessageSubmit(e);
    }
  };
  useEffect(() => {
    inp?.addEventListener('keyup', handleKeyUp);
    return () => inp?.removeEventListener('keyup', handleKeyUp);
  }, [inp?.value]);

  const handleMessageSubmit = (event) => {
    event.preventDefault();
    const messageEvent = {
      event_type: 'ws_group_msg_event',
      sender: user.id,
      sender_username: user.username,
      group_id: group?.id,
      group_title: group?.title,
      message: messageToSend,
      created_date: Math.floor(Date.now() / 1000),
    };
    if (messageEvent.message !== '') {
      webSocket?.send(JSON.stringify(messageEvent));
      handleNewGroupMessageSent(messageEvent);
      setMessageToSend('');
    }
  };

  return (
    <div className="input-box-container">
      <input
        type="text"
        onChange={(e) => setMessageToSend(e.target.value)}
        placeholder="type your message..."
        value={messageToSend}
        id="input"
      />
      <div className="attach-and-emojis">
        <Emoji setMsg={setMessageToSend} />
      </div>
      <button onClick={(e) => handleMessageSubmit(e)}>send</button>
    </div>
  );
};
