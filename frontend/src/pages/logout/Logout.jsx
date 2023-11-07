import { api } from '../../api/axios';

const LOGOUT_URL = '/signout';

const Logout = async () => {
  try {
    await api.post(LOGOUT_URL, null, {
      headers: { 'Content-Type': 'application/json' },
      withCredentials: true,
    });
  } catch (error) {}
};

export default Logout;
