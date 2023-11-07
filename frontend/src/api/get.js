import useSWR from 'swr';
import { getURL } from './axios';

const fetcher = (urls) => {
  return Promise.all(urls.map((url) => getURL(url)));
};
export function useUsers(ids) {
  const { data, error } = useSWR(
    ids?.map((id) => `/user/${id}`),
    ids ? fetcher : null
  );
  return {
    users: data?.map((data) => data?.data),
    error: error,
  };
}
export function useChats() {
  const { data, error } = useSWR('/chat/', getURL);
  let chats = data?.data;
  const { users } = useUsers(chats?.map((chat) => Math.abs(chat.to_from_user)));
  chats = chats?.map((chat, index) => {
    chat.user = users ? users[index] : undefined;
    chat.id = chat.user?.id;
    return chat;
  });
  return {
    chats: chats,
    error: error,
  };
}
