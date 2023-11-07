import React, { Fragment } from 'react';
import { useGet } from '../../hooks/useGet';
import './notifications.scss';
import { Button } from '../groups/Groups';
import { longDate } from '../../components/utils/dateFormatter';
import { useNavigate } from 'react-router-dom';
import { myAxios } from '../../api/axios';
import { usernameOrName } from '../../components/utils/usernameOrName';
import closeIcon from '../../icons/close-icon.svg';
import acceptIcon from '../../icons/accept-icon.svg';
import rejectIcon from '../../icons/reject-icon.svg';
import { errorToast } from '../../components/utils/toast/errorToast';
const NotificationsCenter = () => {
  const { data, setHit } = useGet('/notifications');
  return (
    <div className="notifications-container">
      <ClearButtons />
      <GroupPostCommentNotifications
        notifications={data?.group_post_comment_notification}
        setHit={setHit}
      />
      <GroupJoinInvites
        notifications={data?.group_join_invites}
        setHit={setHit}
      />
      <GroupJoinRequests
        notifications={data?.group_join_requests}
        setHit={setHit}
      />
      <NewGroupEventsNotifications
        notifications={data?.group_events_created}
        setHit={setHit}
      />
      <FollowNotifications
        notifications={data?.users_want_to_follow}
        setHit={setHit}
      />
    </div>
  );
};
//clear notifications buttons || styles are taken from groups.scss
const ClearButtons = () => {
  return (
    <div className="filter-buttons clear">
      <Button text={'notifications'} />
      {/*  <hr />
      <Button text={"select"} /> */}
    </div>
  );
};

//to see user group invites
const GroupJoinInvites = ({ notifications }) => {
  const nav = useNavigate();
  const handleClick = async (groupId) => {
    nav(`/group/${groupId}`);
  };
  return (
    <>
      {notifications?.map((notification, index) => (
        <div className="notification-container" key={index}>
          <p>
            <span className="ntf-type">Invitation: </span>
            <span className="user">
              {usernameOrName(notification.user_invited) + ' '}
            </span>
            invites you to join the{' '}
            <span
              className="group"
              onClick={() => handleClick(notification.group.id)}
            >
              {notification.group.title}
            </span>
          </p>

          <span className="date">{longDate(notification?.created_date)}</span>
        </div>
      ))}
    </>
  );
};

//to see group join requests for admin | ONLY GROUP ADMIN WILL SEE IT
const GroupJoinRequests = ({ notifications, setHit }) => {
  const nav = useNavigate();
  const handleReject = async (groupId, user_id) => {
    await myAxios
      .post(
        `group/${groupId}/reject`,
        JSON.stringify({
          id: Number(user_id),
        })
      )
      .then(() => {
        setHit(true);
      })
      .catch((error) => {
        return;
      });
  };
  const handleAccept = async (groupId, user_id) => {
    await myAxios
      .post(
        `group/${groupId}/approve`,
        JSON.stringify({
          id: Number(user_id),
        })
      )
      .then(() => {
        setHit(true);
      })
      .catch((error) => {
        return;
      });
  };
  return (
    <>
      {notifications?.map((notification, index) => (
        <Fragment key={index}>
          {notification?.users_requested_join.map((user, idx) => (
            <div className="notification-box" key={idx}>
              <div className="notification-container comment">
                <p>
                  <span className="ntf-type">Join request: </span>
                  <span
                    className="user"
                    onClick={() => nav(`/user/${user.id}`)}
                  >
                    {usernameOrName(user) + ' '}
                  </span>
                  wants to join{' '}
                  <span
                    className="group"
                    onClick={() => nav(`/group/${notification?.group.id}`)}
                  >
                    {notification?.group.title}
                  </span>
                </p>
                <div className="accept-reject">
                  <img
                    src={rejectIcon}
                    onClick={() =>
                      handleReject(notification?.group.id, user.id)
                    }
                    alt=""
                  />
                  <img
                    src={acceptIcon}
                    onClick={() =>
                      handleAccept(notification?.group.id, user.id)
                    }
                    alt=""
                  />
                </div>

                <span className="date">{longDate(user?.created_date)}</span>
              </div>
            </div>
          ))}
        </Fragment>
      ))}
    </>
  );
};

//to see group post notifications
const GroupPostCommentNotifications = ({ notifications, setHit }) => {
  const nav = useNavigate();
  const handleClick = async (groupId, postId, commentId, shouldNavigate) => {
    await myAxios
      .put(`/notifications/clear`, {
        comment_id: [commentId],
      })
      .catch((error) => console.log(error))
      .finally(() => {
        if (shouldNavigate) {
          nav(`/group/${groupId}/post/${postId}`);
          setHit(true);
        } else {
          setHit(true);
        }
      });
  };
  return (
    <>
      {notifications?.map((notification) => (
        <div className="notification-box" key={notification?.comment_id}>
          <div
            className="notification-container comment"
            onClick={() =>
              handleClick(
                notification.group_id,
                notification.post_id,
                notification.comment_id,
                true
              )
            }
          >
            <p>
              <span className="ntf-type">Comment:</span>{' '}
              {notification?.comment_content}
            </p>

            <span className="date">{longDate(notification?.created_date)}</span>
          </div>

          <img
            className="clear"
            src={closeIcon}
            alt=""
            onClick={() =>
              handleClick(
                notification.group_id,
                notification.post_id,
                notification.comment_id,
                false
              )
            }
          />
        </div>
      ))}
    </>
  );
};

const NewGroupEventsNotifications = ({ notifications, setHit }) => {
  const nav = useNavigate();
  const handleClear = async (groupId, eventId, shouldNavigate) => {
    await myAxios.put(`/group/${groupId}/event/${eventId}`).then(() => {
      if (shouldNavigate) {
        nav(`/group/${groupId}`);
        setHit(true);
      } else {
        setHit(true);
      }
    });
  };
  return (
    <>
      {notifications?.map((notification, index) => (
        <div className="notification-box" key={index}>
          <div
            className="notification-container"
            onClick={() =>
              handleClear(
                notification.group.id,
                notification.group_event.event_id,
                true
              )
            }
          >
            <p>
              <span className="ntf-type">Event:</span>
              new event -{' '}
              <span className="group">{notification.group_event.title}</span> -
              in group{' '}
              <span className="group">{notification?.group.title}</span>
            </p>
            <span className="date">
              {longDate(notification?.group_event?.created_date)}
            </span>
          </div>
          <img
            className="clear"
            src={closeIcon}
            alt=""
            onClick={() =>
              handleClear(
                notification.group.id,
                notification.group_event.event_id,
                false
              )
            }
          />
        </div>
      ))}
    </>
  );
};

const FollowNotifications = ({ notifications, setHit }) => {
  const nav = useNavigate();
  const handleReject = async (followerId) => {
    await myAxios
      .post(
        '/me/settings/reject-follow',
        JSON.stringify({
          follower: { id: followerId },
        })
      )
      .then(() => {
        setHit(true);
      })
      .catch((error) => {
        return errorToast(error);
      });
  };

  const handleAccept = async (followerId) => {
    await myAxios
      .post(
        '/me/settings/accept-follow',
        JSON.stringify({
          follower: { id: followerId },
        })
      )
      .then(() => {
        setHit(true);
      })
      .catch((error) => {
        return errorToast(error);
      });
  };

  return (
    <>
      {notifications?.map((notification, index) => (
        <div key={index} className="notification-container">
          <p>
            <span className="ntf-type">Follow request: </span>
            <span
              className="user"
              onClick={() => nav(`/user/${notification.id}`)}
            >
              {notification?.username
                ? notification?.username + ' '
                : notification?.firstname + ' '}
            </span>
            wants to follow you
          </p>
          <div className="accept-reject">
            <img
              src={rejectIcon}
              onClick={() => handleReject(notification.id)}
              alt=""
            />
            <img
              src={acceptIcon}
              onClick={() => handleAccept(notification.id)}
              alt=""
            />
          </div>
        </div>
      ))}
    </>
  );
};

export default NotificationsCenter;
