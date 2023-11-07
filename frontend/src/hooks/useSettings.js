import { useEffect, useState } from 'react';
import { myAxios } from '../api/axios';

export const useSettings = () => {
  const [me, setMe] = useState(null);
  const [username, setUsername] = useState('');
  const [about, setAbout] = useState('');
  const [privacy, setPrivacy] = useState('public');
  const [followRequests, setFollowRequests] = useState([]);
  const url = '/me';
  useEffect(() => {
    const fetchMe = async () => {
      await myAxios
        .get(url)
        .then((res) => {
          setMe(res.data);
          setUsername(res.data.username);
          setAbout(res.data.about);
          setPrivacy(res.data.privacy);
          setFollowRequests(res.data.follow_requests);
        })
        .catch((err) => console.log('err in useSettings fetch', err));
    };
    fetchMe();
  }, []);
  return {
    me,
    privacy,
    about,
    username,
    followRequests,
    setPrivacy,
    setAbout,
    setUsername,
    setFollowRequests,
  };
};
