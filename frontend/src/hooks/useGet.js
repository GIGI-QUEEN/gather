//custom hook to send GET requests
import { useEffect, useState } from 'react';
import { myAxios } from '../api/axios';

export const useGet = (url) => {
  const [data, setData] = useState(null);
  const [error, setError] = useState(null);
  const [hit, setHit] = useState(false);
  //follow state to handle follow button
  const [follow, setFollow] = useState(false);
  useEffect(() => {
    setHit(false);
    setFollow(false);
    myAxios
      .get(url)

      .then((res) => {
        setData(res.data);
      })
      .catch((error) => {
        setError(error);
      });
  }, [url, hit, follow]);

  return { data, error, setHit, setFollow };
};
