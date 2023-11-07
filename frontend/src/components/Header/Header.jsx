import React from 'react';
import styles from './Header.module.scss';
import './Header.module.scss';
import { useNavigate } from 'react-router-dom';
import { useState, useEffect, useContext } from 'react';
import { UserContext } from '../utils/UserContext';
import { myAxios } from '../../api/axios';
import createIcon from '../../icons/create-icon.svg';
import messagesIcon from '../../icons/messages-icon.svg';
import feedIcon from '../../icons/feed-icon.svg';
import profileIcon from '../../icons/profile-icon.svg';
import groupsIcon from '../../icons/groups-icon.svg';
import notificationsIcon from '../../icons/notifications-icon.svg';
export const Header = () => {
  const [setOpenChat] = useState(undefined);
  const { user } = useContext(UserContext);
  const nav = useNavigate();
  const toggleClass = (elementId) => {
    let l = document.getElementById(elementId);
    const elements = document.querySelectorAll('li');
    elements.forEach((element) => {
      element.classList.remove(styles.current_page);
    });
    l?.classList.add(styles.current_page);
  };

  useEffect(() => {
    myAxios
      .get('/chat/')

      .then((res) => {
        setOpenChat(res.data);
      });
  }, [user, setOpenChat]);
  return (
    <div className={styles.headerContainer}>
      <div className={styles.title}>Gather</div>
      <nav className={styles.nav}>
        <ul>
          <li
            id="li-1"
            onClick={() => {
              nav('/create');
              toggleClass('li-1');
            }}
          >
            <img src={createIcon} alt="create link" /> create
          </li>
          <li
            id="li-2"
            onClick={() => {
              toggleClass('li-2');
              nav('/chat');
            }}
          >
            <img src={messagesIcon} alt="chat link" /> messages
          </li>{' '}
          <li
            onClick={() => {
              nav('/posts');
              toggleClass('li-3');
            }}
            id="li-3"
          >
            <img src={feedIcon} alt="feed link" /> feed
          </li>{' '}
          <li
            onClick={() => {
              nav('/me');
              toggleClass('li-4');
            }}
            id="li-4"
          >
            <img src={profileIcon} alt="profile link" />
            my profile
          </li>{' '}
          <li
            onClick={() => {
              nav('/groups');
              toggleClass('li-5');
            }}
            id="li-5"
          >
            <img src={groupsIcon} alt="groups link" /> groups
          </li>{' '}
          <li
            id="li-6"
            onClick={() => {
              nav('/notifications');
              toggleClass('li-6');
            }}
          >
            <img src={notificationsIcon} alt="notifications link" /> center
          </li>
        </ul>
      </nav>
      <LogoutButton />
    </div>
  );
};

const LogoutButton = () => {
  const { setLogged } = useContext(UserContext);

  const handleClick = async () => {
    await myAxios.post('/signout', null).then(() => setLogged(false));
  };
  return (
    <button onClick={handleClick} className={styles.logout_button}>
      Logout
    </button>
  );
};
