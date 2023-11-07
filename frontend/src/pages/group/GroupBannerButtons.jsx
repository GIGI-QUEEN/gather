import { useContext } from 'react';
import cancelIcon from '../../icons/cancel-icon.svg';
import wrenchIcon from '../../icons/wrench-icon.svg';
import joinIncon from '../../icons/plus-icon.svg';
import leaveIcon from '../../icons/leave-icon.svg';
import { useGet } from '../../hooks/useGet';
import { useNavigate } from 'react-router-dom';
import { myAxios } from '../../api/axios';
import { errorToast } from '../../components/utils/toast/errorToast';
import inviteIcon from '../../icons/invite-icon.svg';
import { handleModal } from '../../components/utils/handleModal';
import { WebSocketContext } from '../../components/utils/WebSocketContext';
import { UserContext } from '../../components/utils/UserContext';

export const JoinLeaveCancelButtons = ({ group, setHit }) => {
  const { data: me } = useGet('/me');
  if (group?.admin.id === me?.id) {
    return <SettingsButton group_id={group?.id} />;
  }

  if (group?.user_invited === true) {
    return <AcceptRejectInviteButtons group_id={group?.id} setHit={setHit} />;
  }

  if (group?.join_requested === false && group?.join_approved === false) {
    return (
      <JoinButton
        group_id={group?.id}
        setHit={setHit}
        group_title={group?.title}
      />
    );
  }
  if (group?.join_requested === true && group?.join_approved === false) {
    return <CancelJoinButton group_id={group?.id} setHit={setHit} />;
  }
  if (group?.join_requested === true && group?.join_approved === true) {
    return <LeaveButton group_id={group?.id} setHit={setHit} />;
  }
};

const JoinButton = ({ group_id, setHit, group_title }) => {
  const { webSocket } = useContext(WebSocketContext);
  const { user } = useContext(UserContext);
  const handleJoin = (group_id) => {
    myAxios
      .post(`/group/${group_id}/join`)
      .then(() => {
        const groupJoinRequestEvent = {
          event_type: 'ws_group_join_request_event',
          group_id: parseInt(group_id),
          group_title: group_title,
          user_to_join: user.username ? user.username : user.firstname,
        };
        webSocket?.send(JSON.stringify(groupJoinRequestEvent));
        setHit(true);
      })
      .catch((error) => {
        errorToast(error);
        return;
      });
  };
  return (
    <div className="container_1 action-button">
      <button onClick={() => handleJoin(group_id)}>
        join <img src={joinIncon} alt="" />
      </button>
    </div>
  );
};

const LeaveButton = ({ group_id, setHit }) => {
  const handleLeave = async (group_id) => {
    await myAxios
      .post(`/group/${group_id}/leave`)
      .then(() => {
        setHit(true);
      })
      .catch((error) => {
        errorToast(error);

        return;
      });
  };
  return (
    <div className="container_1 action-button">
      <button onClick={() => handleLeave(group_id)}>
        leave <img src={leaveIcon} alt="leave icon" />
      </button>
    </div>
  );
};

const CancelJoinButton = ({ group_id, setHit }) => {
  const handleLeave = async (group_id) => {
    await myAxios
      .post(`/group/${group_id}/leave`)
      .then(() => {
        setHit(true);
      })
      .catch((error) => {
        errorToast(error);

        return;
      });
  };
  return (
    <div className="container_1 action-button">
      <button onClick={() => handleLeave(group_id)}>
        cancel <img src={cancelIcon} alt="" />
      </button>
    </div>
  );
};

const SettingsButton = ({ group_id }) => {
  const nav = useNavigate();
  return (
    <div className="container_1 action-button">
      <button onClick={() => nav(`/group/${group_id}/settings`)}>
        settings <img src={wrenchIcon} alt="settings icon" />
      </button>
    </div>
  );
};

const AcceptRejectInviteButtons = ({ group_id, setHit }) => {
  const handleInviteReject = async () => {
    await myAxios
      .post(`group/${group_id}/reject-invite`)
      .then(() => {
        setHit(true);
      })
      .catch((error) => {
        errorToast(error);

        return;
      });
  };
  const handleInviteAccept = async () => {
    await myAxios
      .post(`group/${group_id}/accept-invite`)
      .then(() => {
        setHit(true);
      })
      .catch((error) => {
        errorToast(error);

        return;
      });
  };
  return (
    <div>
      <div className="container_1 action-button reject-accept-invite">
        <button onClick={() => handleInviteReject()}>reject invite</button>
      </div>
      <div className="container_1 action-button reject-accept-invite">
        <button onClick={() => handleInviteAccept()}>accept invite</button>
      </div>
    </div>
  );
};

export const InviteButton = ({ group }) => {
  return (
    <div
      className={
        group?.join_approved === true && group?.join_requested === true
          ? 'container_1 action-button invite-button'
          : 'container_1 action-button invite-button disabled'
      }
    >
      <button onClick={() => handleModal('followers-list')}>
        invite <img src={inviteIcon} alt="" />
      </button>
    </div>
  );
};
