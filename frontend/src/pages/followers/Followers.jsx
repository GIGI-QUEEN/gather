import { useParams } from 'react-router-dom';
import UsersList from '../../components/UsersList';
import { useGet } from '../../hooks/useGet';

const Followers = () => {
  const { id } = useParams();

  const { data: user } = useGet(`/user/${id}`);

  return (
    <div>
      <UsersList users={user?.followers} title={'Followers'} pageOwner={user} />
    </div>
  );
};

export default Followers;
