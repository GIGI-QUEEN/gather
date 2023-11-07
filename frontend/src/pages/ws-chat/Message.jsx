import { longDate } from '../../components/utils/dateFormatter';
import { usernameOrName } from '../../components/utils/usernameOrName';

export const Message = ({ message, user, typeOfChat }) => {
  const check = message.recipient !== user?.id && message.sender === user?.id;
  const messageBody = (
    <div className="message-box">
      <div className="content">
        {typeOfChat === 'group-chat' ? (
          <span className="name">
            {check
              ? usernameOrName(user)
              : message.sender_username
              ? message.sender_username
              : message.firstname}
          </span>
        ) : null}
        <p className="message-body">{message.message}</p>
      </div>
      <div className="message-date">{longDate(message.created_date)}</div>
    </div>
  );

  const content = (
    <div
      className={
        message.recipient !== user?.id && message.sender === user?.id
          ? 'message owner'
          : 'message'
      }
    >
      {messageBody}
    </div>
  );

  return content;
};
