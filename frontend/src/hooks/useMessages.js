import axios from 'axios';
import { useState } from 'react';
import { useEffect } from 'react';
import { myAxios } from '../api/axios';

export const useMessages = (id, index = 10, ws, chatBox) => {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);
  const [messages, setMessages] = useState([]);
  const [hasMore, setHasMore] = useState(false);
  const [user, setUser] = useState({});

  if (messages.length <= 10) {
    chatBox?.scrollTo(0, chatBox?.scrollHeight);
  }

  useEffect(() => {
    setMessages([]);
  }, [id]);

  useEffect(() => {
    myAxios.get(`/user/${id}`).then((res) => setUser(res.data));
  }, [id]);

  useEffect(() => {
    setLoading(true);
    let cancel;
    myAxios
      .get(`/chat/${id}`, {
        cancelToken: new axios.CancelToken((c) => (cancel = c)),
      })

      .then((res) => {
        setMessages([...messages, ...res.data.slice(index, index + 10)]);
        setHasMore(messages.length !== res.data.length);
        setLoading(false);
      })
      .catch((e) => {
        if (axios.isCancel(e)) return;
        setError(true);
      });

    return () => cancel();
  }, [id, index]);

  useEffect(() => {
    const messageHandler = async () => {
      await myAxios
        .get(`/chat/${id}`)

        .then((res) => {
          let msgs = res.data;
          setMessages([...msgs]);
        });
    };

    ws?.addEventListener('message', messageHandler);
    return () => {
      ws?.removeEventListener('message', messageHandler);
    };
  }, [ws, id]);

  return { loading, error, messages, hasMore, user };
};
