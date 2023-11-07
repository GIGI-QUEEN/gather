import { useContext } from 'react';
import { useNavigate } from 'react-router-dom';
import { UserContext } from './utils/UserContext';
import { Post } from './Post';
import { myAxios } from '../api/axios';
import { WebSocketContext } from './utils/WebSocketContext';
import { errorToast } from './utils/toast/errorToast';
export const UserProfile = ({ user, id, follow, hit }) => {
  const me = useContext(UserContext);
  const nav = useNavigate();
  const findMe = user?.followers?.find(
    (follower) => follower?.id === me.user.id
  );
  return (
    <div className="content-container">
      <div className="profile-info-container">
        <div className="left-section">
          <div className="user-avatar">
            {user?.avatar ? (
              <img src={`http://localhost:8080/${user?.avatar}`} alt="Avatar" />
            ) : (
              <div className="empty-avatar user-profile-avatar"></div>
            )}
          </div>
        </div>
        <div className="right-section">
          <div className="stats">
            <ul>
              <li>
                Posts: {user?.posts?.length > 0 ? user?.posts?.length : 0}
              </li>
              <li
                onClick={() => nav(`/user/${user?.id}/followers`)}
                className="stats-link"
              >
                Followers:{' '}
                {user?.followers?.length > 0 ? user?.followers?.length : 0}
              </li>
              <li
                onClick={() => nav(`/user/${user?.id}/followings`)}
                className="stats-link"
              >
                Following:{' '}
                {user?.followings?.length > 0 ? user?.followings?.length : 0}
              </li>
            </ul>
          </div>
          <div className="user-info">
            <h1>
              {user?.firstname} {user?.lastname}
            </h1>
            {user?.username ? <p>@{user?.username}</p> : null}
            {!findMe &&
            user?.id !== me?.user?.id &&
            user?.privacy === 'private' ? null : (
              <>
                <p>Email: {user?.email}</p>
                <p>Age: {user?.age}</p>
                <p className="about">{user?.about}</p>
              </>
            )}
          </div>
        </div>
      </div>
      <ActionButtons user={user} me={me} follow={follow} />
      <div className="users-posts">
        {user?.posts?.map((post) => (
          <Post
            key={post.id}
            id={post.id}
            title={post.title}
            content={post.content}
            author={post.user}
            likes={post.likes}
            dislikes={post.dislikes}
            hit={hit}
            image={post.image}
            actionURL={`/post/${post?.id}`}
          />
        ))}
      </div>
    </div>
  );
};

const ActionButtons = ({ user, me, follow }) => {
  const navigate = useNavigate();
  const handleClick = () => {
    navigate(`/chat/`);
  };
  const { webSocket } = useContext(WebSocketContext);

  const find = user?.followers?.find(
    (user) => user.username === me.user.username
  );
  const handleFollow = async (e) => {
    await myAxios
      .post(
        `/user/${user.id}/follow`,
        JSON.stringify({
          followed_id: user.id,
        })
      )
      .then(() => {
        follow(true);
        if (user?.privacy === 'private') {
          const followRequestEvent = {
            event_type: 'ws_follow_request_event',
            user_to_follow: Number(user.id),
            follower_id: Number(me.user.id),
            follower_name: me.user.username,
          };
          webSocket?.send(JSON.stringify(followRequestEvent));
        }
      })
      .catch((err) => {
        errorToast(err);
        return;
      });
  };

  const handleUnfollow = async () => {
    await myAxios
      .post(`/user/${user.id}/unfollow`)
      .then(() => {
        follow(true);
      })
      .catch((err) => {
        errorToast(err);
        return;
      });
  };

  return (
    <div>
      {user?.id !== me?.user.id ? (
        <div className="action-buttons">
          {find !== undefined || user?.follow_requested === true ? (
            <button onClick={handleUnfollow}>Unfollow</button>
          ) : (
            <button onClick={handleFollow}>Follow</button>
          )}

          <button onClick={handleClick}>Write message</button>
        </div>
      ) : (
        <div className="action-buttons">
          <button onClick={() => navigate('/me/settings')}>Settings</button>
        </div>
      )}
    </div>
  );
};
