import React from 'react';
import { useParams } from 'react-router-dom';
import { useGet } from '../../hooks/useGet';
import '../settings/settings.scss';
import { myAxios } from '../../api/axios';
import { BackButton } from '../../components/BackButton';
const GroupSettingsPage = () => {
  const { groupid: id } = useParams();
  const { data: group, setHit } = useGet(`group/${id}`);
  return (
    <div className="settings-wrapper">
      <div className="sections">
        {group?.join_requests?.length > 0 ? (
          <JoinRequests
            requests={group?.join_requests}
            groupId={group?.id}
            setHit={setHit}
          />
        ) : (
          <div>No join requests</div>
        )}
        <BackButton url={`/group/${id}`} />
      </div>
    </div>
  );
};

const JoinRequests = ({ requests, groupId, setHit }) => {
  const handleReject = async (user_id) => {
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
      .catch(() => {
        return;
      });
  };
  const handleApprove = async (user_id) => {
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
      .catch(() => {
        return;
      });
  };
  return (
    <div className="container users-to-approve">
      {requests?.map((request) => (
        <div key={request?.id} className="user-to-approve">
          <p>
            {request?.username
              ? request?.username
              : `${request?.firstName} ${request?.lastname}`}
          </p>
          <div className="approve-reject-buttons">
            <button onClick={() => handleReject(request?.id)}>reject</button>
            <button onClick={() => handleApprove(request?.id)}>approve</button>
          </div>
        </div>
      ))}
    </div>
  );
};

export default GroupSettingsPage;
