import { myAxios } from '../../api/axios';
import { options } from '../../components/utils/toast/options';
import { toast } from 'react-toastify';
import { handleModal } from '../../components/utils/handleModal';
import closeIcon from '../../icons/close-icon.svg';
import { useContext } from 'react';
import { WebSocketContext } from '../../components/utils/WebSocketContext';
import { UserContext } from '../../components/utils/UserContext';
import { usernameOrName } from '../../components/utils/usernameOrName';

export const FollowersToInviteList = ({ followers, group_id, group_title }) => {
  const { webSocket } = useContext(WebSocketContext);
  const { user } = useContext(UserContext);
  const handleInvite = async (user_id) => {
    await myAxios
      .post(
        `/group/${group_id}/invite`,
        JSON.stringify({
          id: user_id,
        })
      )
      .then(() => {
        const groupJoinRequestEvent = {
          event_type: 'ws_group_join_invite_event',
          group_title: group_title,
          user_inviting: user.username ? user.username : user.firstname,
          user_to_invite: user_id,
        };
        webSocket?.send(JSON.stringify(groupJoinRequestEvent));
        toast.success('Successfully invited');
        handleModal('followers-list');
      })
      .catch((error) => {
        toast.error(error.response['data']['error_description'], options);
        return;
      });
  };

  return (
    <div className="followers-list-modal-container" id="followers-list">
      <div className="list-container">
        <div className="title-and-close">
          <h3>followers</h3>
          <button
            className="close-btn"
            onClick={() => handleModal('followers-list')}
          >
            <img src={closeIcon} alt="" />
          </button>
        </div>

        {followers?.map((follower) => (
          <div className="follower-card" key={follower.id}>
            <img src={`http://localhost:8080/${follower?.avatar}`} alt="" />
            {usernameOrName(follower)}
            <button onClick={() => handleInvite(follower.id)}>invite</button>
          </div>
        ))}
      </div>
    </div>
  );
};
