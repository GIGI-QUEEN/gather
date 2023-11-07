import { useContext, useEffect, useState } from 'react';
import {
  sendChatMessagesRequest,
  sendGroupChatMessagesRequest,
} from './backend-api-calls';
import { WebSocketContext } from '../../components/utils/WebSocketContext';
import { Input, GroupInput } from './Input';
import { ChatHeader } from './ChatHeader';
import { Messages } from './Messages';
import { useNewMessage } from '../../hooks/useNewMessage';
import './chat.scss';
import Image from '../../components/Image';
import { useGet } from '../../hooks/useGet';
import { shortDate } from '../../components/utils/dateFormatter';
import { useNavigate } from 'react-router-dom';
import groupAvatar from '../../icons/group-avatar.svg';

const ChatBlock = ({}) => {
  const [selectedUser, setSelectedUser] = useState(null);
  const [conversationsMap, setConversationsMap] = useState(null);

  const [selectedGroup, setSelectedGroup] = useState(null);
  const [groupConversationsMap, setGroupConversationsMap] = useState(null);

  const [messageToSend, setMessageToSend] = useState('');

  const [typeOfChat, setTypeOfChat] = useState('private'); //could be "private" if it's simple 1 on 1 chat or "group-chat" if it's group-chat
  const { webSocket } = useContext(WebSocketContext);
  const { data: me } = useGet('/me');

  useEffect(() => {
    if (selectedGroup) {
      setTypeOfChat('group-chat');
      setSelectedUser(null);
    }
    if (selectedUser) {
      setTypeOfChat('private');
      setSelectedGroup(null);
    }
  }, [selectedGroup, selectedUser]);

  useEffect(() => {
    const fetchData = async () => {
      if (selectedUser && !conversationsMap?.get(selectedUser.id)) {
        await sendChatMessagesRequest(selectedUser.id)
          .then((res) => {
            const messages = res.data.reverse();
            setConversationsMap((prevMap) => {
              const map = prevMap ? new Map(prevMap) : new Map();
              if (!map.has(selectedUser.id)) {
                map.set(selectedUser.id, messages);
              }
              return map;
            });
          })
          .catch((err) => console.log(err.message));
      }

      if (selectedGroup && !groupConversationsMap?.get(selectedGroup?.id)) {
        await sendGroupChatMessagesRequest(selectedGroup?.id)
          .then((res) => {
            const messages = res.data.reverse();
            setGroupConversationsMap((prevMap) => {
              const map = prevMap ? new Map(prevMap) : new Map();
              if (!map.has(selectedGroup?.id)) {
                map.set(selectedGroup?.id, messages);
              }
              return map;
            });
          })
          .catch((err) => console.log(err.message));
      }
    };

    fetchData();
  }, [selectedUser, selectedGroup]);
  //---------------------------------------------------------------------------------------------------
  //puts new messages into map
  useNewMessage(webSocket, setConversationsMap, setGroupConversationsMap);

  const conversationMessages = conversationsMap?.get(selectedUser?.id);
  const groupConversationMessages = groupConversationsMap?.get(
    selectedGroup?.id
  );

  return (
    <div className="chat-box">
      <ChatsList
        user={me}
        selectedUser={selectedUser}
        setSelectedUser={setSelectedUser}
        setSelectedGroup={setSelectedGroup}
      />
      <div className="right-section">
        <ChatHeader header={selectedUser ? selectedUser : selectedGroup} />
        {typeOfChat === 'private' ? (
          <Messages
            messages={conversationMessages}
            user={me}
            typeOfChat={typeOfChat}
          />
        ) : (
          <Messages
            messages={groupConversationMessages}
            user={me}
            typeOfChat={typeOfChat}
          />
        )}
        {selectedGroup === null && selectedUser === null ? null : (
          <>
            {typeOfChat === 'private' ? (
              <Input
                messageToSend={messageToSend}
                setMessageToSend={setMessageToSend}
                setConversationsMap={setConversationsMap}
                webSocket={webSocket}
                recipient={selectedUser}
                user={me}
              />
            ) : (
              <GroupInput
                messageToSend={messageToSend}
                setMessageToSend={setMessageToSend}
                setGroupConversationsMap={setGroupConversationsMap}
                webSocket={webSocket}
                group={selectedGroup}
                user={me}
              />
            )}
          </>
        )}
      </div>
    </div>
  );
};

const ChatsList = ({
  user,
  selectedUser,
  setSelectedUser,
  setSelectedGroup,
}) => {
  const { data } = useGet('/chat/');

  const allIds = data
    ?.map((obj) => obj.recipient.id)
    .concat(data.map((obj) => obj.sender.id));
  const uniqueIds = [...new Set(allIds)].filter((userId) => userId != user?.id);

  return (
    <div className="chats-list">
      <div className="first-block">
        <div className="chats-list-header"></div>
        <div className="divider">Private chats</div>
        <ListOfChats
          chats={data}
          selectedUser={selectedUser}
          setSelectedUser={setSelectedUser}
          setSelectedGroup={setSelectedGroup}
          user={user}
          me={user}
        />
      </div>
      <div className="second-block">
        <ListOfFollowingsToChat
          user={user}
          setSelectedUser={setSelectedUser}
          uniqueIds={uniqueIds}
        />
        <div className="divider">Group chats</div>
        <ListOfGroupChats
          user={user}
          setSelectedGroup={setSelectedGroup}
          setSelectedUser={setSelectedUser}
        />
      </div>
    </div>
  );
};

const ListOfChats = ({ user, chats, setSelectedUser, setSelectedGroup }) => {
  if (!chats) return null;
  const handleChatWith = (chat) => {
    let chatWith;
    if (chat?.recipient?.id !== user?.id) {
      chatWith = chat.recipient;
    } else if (chat.sender.id !== user?.id) {
      chatWith = chat.sender;
    }
    return chatWith;
  };
  return (
    <div>
      {chats?.map((chat, index) => (
        <ChatWithUser
          key={index}
          user={handleChatWith(chat)}
          content={chat.message}
          date={chat.created_date}
          setSelectedUser={setSelectedUser}
          setSelectedGroup={setSelectedGroup}
          sender={chat.sender}
          recipient={chat.recipient}
        />
      ))}
    </div>
  );
};

const ChatWithUser = ({
  user,
  content,
  date,
  setSelectedUser,
  setSelectedGroup,
}) => {
  const handleClick = () => {
    setSelectedUser(user);
    setSelectedGroup(null);
  };
  if (content?.length >= 8) {
    content = content.slice(0, 7) + '...';
  }
  return (
    <div className="chat-with-user_container" onClick={() => handleClick()}>
      <div className="avatar">
        <Image image={user?.avatar} className={'chat-list-avatar'} />
        <div
          className={user?.status === 1 ? 'status online' : 'status offline'}
        ></div>
      </div>
      <div className="chat-info">
        <h2 className="chat-info__user">{user?.username}</h2>
        <div className="chat-info__content">
          <p>{content}</p>
          <p>{shortDate(date)}</p>
        </div>
      </div>
    </div>
  );
};

const FollowSomeoneToChat = () => {
  const navigate = useNavigate();
  return (
    <div
      className="chat-with-user_container follow-someone"
      onClick={() => {
        navigate('/user/all');
      }}
    >
      <h2>Follow someone</h2>
      <p>to start chat</p>
    </div>
  );
};

const AvailableContact = ({ user, setSelectedUser }) => {
  return (
    <div
      className="available-contact-container"
      onClick={() => setSelectedUser(user)}
    >
      <Image image={user.avatar} className={'available-contact-image'} />
      <p>{user.username}</p>
    </div>
  );
};

const ListOfFollowingsToChat = ({ user, setSelectedUser, uniqueIds }) => {
  let copiedFollowings = user?.followings?.filter(
    (obj) => !uniqueIds.includes(obj.id)
  );
  if (user?.followings) {
    return (
      <>
        <div>
          {copiedFollowings.map((following) => (
            <AvailableContact
              user={following}
              key={user.id}
              setSelectedUser={setSelectedUser}
            />
          ))}
        </div>
      </>
    );
  } else {
    return <FollowSomeoneToChat />;
  }
};

const ListOfGroupChats = ({ user, setSelectedGroup, setSelectedUser }) => {
  if (user?.member_in_groups) {
    return (
      <>
        {user?.member_in_groups.map((group) => (
          <AvailableGroupChat
            group={group}
            key={group.id}
            setSelectedGroup={setSelectedGroup}
            setSelectedUser={setSelectedUser}
          />
        ))}
      </>
    );
  } else {
    return <JoinGroupToChat />;
  }
};

const AvailableGroupChat = ({ group, setSelectedGroup, setSelectedUser }) => {
  return (
    <div
      className="available-contact-container"
      onClick={() => {
        setSelectedGroup(group);
        setSelectedUser(null);
      }}
    >
      <img className="group-avatar" src={groupAvatar} alt="" />
      <p>{group.title}</p>
    </div>
  );
};

const JoinGroupToChat = () => {
  const navigate = useNavigate();
  return (
    <div
      className="chat-with-user_container follow-someone"
      onClick={() => {
        navigate('/groups');
      }}
    >
      <h2>Join the group</h2>
      <p>and chat there</p>
    </div>
  );
};

export default ChatBlock;
