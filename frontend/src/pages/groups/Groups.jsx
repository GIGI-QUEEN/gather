import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Modal from '../../components/Modal';
import { handleModal } from '../../components/utils/handleModal';
import { useGet } from '../../hooks/useGet';
import plusIcon from '../../icons/plus-icon-white.svg';
import closeIcon from '../../icons/close-icon.svg';
import './groups.scss';
import { myAxios } from '../../api/axios';
import { toast, ToastContainer } from 'react-toastify';
import { options } from '../../components/utils/toast/options';
import { clearInput } from '../../components/utils/clearInput';
const GroupsPage = () => {
  const { data: groups, setHit } = useGet('/groups');
  return (
    <div className="groups-container">
      <FilterButtons />
      <GroupsList groups={groups} />
      <CreateNewGroupModal setHit={setHit} />
      <Modal />
      <ToastContainer />
    </div>
  );
};

const GroupsList = ({ groups }) => {
  return (
    <div className="groups-list-container">
      {groups?.map((group) => (
        <Group key={group?.id} group={group} />
      ))}
    </div>
  );
};

const Group = ({ group }) => {
  const navigate = useNavigate();
  return (
    <div
      className="banner list-banner"
      onClick={() => navigate(`/group/${group?.id}`)}
    >
      <h2>{group?.title}</h2>
      <p>{group?.description}</p>
    </div>
  );
};

const FilterButtons = () => {
  return (
    <div className="filter-buttons">
      <Button text={'all groups'} />
      <hr />
      <CreateGroupButton />
    </div>
  );
};

export const Button = ({ text }) => {
  return <button className="button">{text}</button>;
};

const CreateNewGroupModal = ({ setHit }) => {
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const handleClick = async (e) => {
    e.preventDefault();
    await myAxios
      .post(
        '/group/create',
        JSON.stringify({
          title,
          description,
        })
      )
      .then(() => {
        setHit(true);
        handleModal('create-group');
        clearInput('input', setTitle);
        clearInput('textarea', setDescription);
      })
      .catch((error) => {
        toast.error(error.response['data']['error_description'], options);
        return;
      });
  };
  return (
    <div className="create-modal-container" id="create-group">
      <div className="title-and-close">
        <input
          type="text"
          placeholder="title"
          onChange={(e) => setTitle(e.target.value)}
        />
        <span
          onClick={() => {
            handleModal('create-group');
          }}
        >
          <img src={closeIcon} alt="" />
        </span>
      </div>
      <textarea
        name=""
        id=""
        placeholder="describe your group..."
        onChange={(e) => setDescription(e.target.value)}
      ></textarea>
      <div className="upload-and-create group-create">
        <button onClick={handleClick}>create</button>
      </div>
    </div>
  );
};

const CreateGroupButton = () => {
  return (
    <button
      className="button plus"
      onClick={() => {
        handleModal('create-group');
      }}
    >
      new <img src={plusIcon} alt="" />
    </button>
  );
};

export default GroupsPage;
