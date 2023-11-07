import React, { useRef, useState } from 'react';
import { myAxios, uploadAxios } from '../../api/axios';
import { options } from '../../components/utils/toast/options';
import { toast } from 'react-toastify';
import './settings.scss';
import { useSettings } from '../../hooks/useSettings';
import { checkFileType } from '../../components/utils/checkFileType';
import { useGet } from '../../hooks/useGet';
import { errorToast } from '../../components/utils/toast/errorToast';

const SettingsPage = () => {
  const [save, setSave] = useState('save changes');
  const [disabled, setDisabled] = useState(true);
  const [file, setFile] = useState(null);
  const form = useRef(null);

  const {
    me,
    privacy,
    about,
    username,
    followRequests,
    setPrivacy,
    setAbout,
    setUsername,
    setFollowRequests,
  } = useSettings();

  let fileType = file?.name.split('.')[1];

  const handleClick = async (e) => {
    e.preventDefault();
    await myAxios.post('/me/settings/privacy', JSON.stringify({ privacy }));
    if (about?.length > 100) {
      toast.error('About is too long', options);
      return;
    } else {
      await myAxios.post('/me/settings/about', JSON.stringify({ about }));
    }
    if (username?.length > 16) {
      toast.error('Username is too long', options);
      return;
    } else {
      await myAxios.post('/me/settings/username', JSON.stringify({ username }));
    }
    if (file !== null) {
      if (checkFileType(file)) {
        const data = new FormData(form.current);
        data.append('image_type', 'avatar');
        data.append('file_type', fileType);
        await uploadAxios.post('/changeavatar', data).catch((error) => {
          toast.error(error.response['data']['error_description'], options);

          return;
        });
      } else {
        return;
      }
    }

    setSave('saved');
    setDisabled(true);
    toast.success('Changes succesfully made!');
  };

  return (
    <div className="settings-wrapper">
      <div className="save-button">
        <button
          onClick={handleClick}
          className="save-changes"
          disabled={disabled}
        >
          {save}
        </button>
      </div>
      <div className="sections">
        <FollowRequests
          followRequests={followRequests}
          setFollowRequests={setFollowRequests}
        />
        <About
          about={about}
          changeAbout={setAbout}
          setSave={setSave}
          setDisabled={setDisabled}
        />
        <Username
          username={username}
          setUsername={setUsername}
          setSave={setSave}
          setDisabled={setDisabled}
        />
        <Privacy
          privacy={privacy}
          changePrivacy={setPrivacy}
          setSave={setSave}
          setDisabled={setDisabled}
        />
        <Avatar
          me={me}
          ref={form}
          setSave={setSave}
          setDisabled={setDisabled}
          setFile={setFile}
        />
      </div>
    </div>
  );
};

const About = ({ about, changeAbout, setSave, setDisabled }) => {
  let limit = document.getElementById('limit');

  if (about?.length > 100) {
    limit.classList.add('limit');
  } else if (about?.length <= 100) {
    limit?.classList.remove('limit');
  }
  return (
    <div className="container about">
      <h2>About me</h2>
      <textarea
        value={about}
        onChange={(e) => {
          changeAbout(e.target.value);
          setSave('save changes');
          setDisabled(false);
        }}
      ></textarea>
      <label id="limit">{about?.length ? about?.length : '0'}/100</label>
    </div>
  );
};

const Privacy = ({ privacy, changePrivacy, setSave, setDisabled }) => {
  let publicProfile = document.getElementById('public');
  let privateProfile = document.getElementById('private');
  if (privacy === 'public') {
    privateProfile?.classList.remove('toggeled');

    publicProfile?.classList.add('toggeled');
  }
  if (privacy === 'private') {
    publicProfile?.classList.remove('toggeled');

    privateProfile?.classList.add('toggeled');
  }

  const handlePublicToggle = () => {
    privateProfile.classList.remove('toggeled');
    publicProfile.classList.add('toggeled');
    changePrivacy('public');
    setSave('save changes');
    setDisabled(false);
  };
  const handlePrivateToggle = () => {
    publicProfile.classList.remove('toggeled');
    privateProfile.classList.add('toggeled');
    changePrivacy('private');
    setSave('save changes');
    setDisabled(false);
  };
  return (
    <div className="container privacy">
      <h2>Account privacy</h2>
      <div className="toggle-button">
        <div
          id="public"
          className="toggler"
          onClick={() => handlePublicToggle()}
        >
          public
        </div>
        <div
          id="private"
          className="toggler"
          onClick={() => handlePrivateToggle()}
        >
          private
        </div>
      </div>
    </div>
  );
};

const Username = ({ username, setUsername, setSave, setDisabled }) => {
  const limit = document.getElementById('username-limit');
  if (username?.length > 16) {
    limit?.classList.add('limit');
  } else {
    limit?.classList.remove('limit');
  }

  return (
    <div className="container username">
      <h2>Username</h2>
      <div className="username_input">
        <label htmlFor="">@</label>
        <input
          type="text"
          value={username ? username : ''}
          placeholder="enter your username"
          onChange={(e) => {
            setUsername(e.target.value);
            setSave('save changes');
            setDisabled(false);
          }}
        />
      </div>
      <span htmlFor="" id="username-limit">
        {username?.length ? username?.length : 0}/16
      </span>
    </div>
  );
};

const Avatar = React.forwardRef(
  ({ me, setFile, setSave, setDisabled }, ref) => {
    const [image, setImage] = useState(null);

    const onImageChange = (e) => {
      if (e.target.files && e.target.files[0]) {
        setImage(URL.createObjectURL(e.target.files[0]));
      }
      setSave('save changes');
      setDisabled(false);
    };

    const src = `http://localhost:8080/${me?.avatar}`;
    return (
      <div className="container avatar">
        <div className="user-avatar">
          <img src={image !== null ? image : src} alt="" />
        </div>
        <form action="" encType="multipart/form-data" ref={ref}>
          <input
            type="file"
            name="avatar"
            id="avatar-input"
            onChange={(e) => {
              setFile(e.target.files[0]);
              onImageChange(e);
            }}
          />
        </form>
      </div>
    );
  }
);

const FollowRequests = () => {
  const { data, setHit } = useGet('/me');

  const handleReject = async (followerId) => {
    await myAxios
      .post(
        '/me/settings/reject-follow',
        JSON.stringify({
          follower: { id: followerId },
        })
      )
      .then(() => {
        setHit(true);
      })
      .catch((error) => {
        return errorToast(error);
      });
  };

  const handleAccept = async (followerId) => {
    await myAxios
      .post(
        '/me/settings/accept-follow',
        JSON.stringify({
          follower: { id: followerId },
        })
      )
      .then(() => {
        setHit(true);
      })
      .catch((error) => {
        return errorToast(error);
      });
  };

  return (
    <div className="container follow-requests">
      <h2>Follow requests</h2>
      {!data?.follow_requests ? <p>No follow requests</p> : null}
      {data?.follow_requests?.map((request) => (
        <div key={request.id} className="request-container">
          <span className="name">
            {request.username
              ? request.username
              : `${request.firstname} ${request.lastname}`}
          </span>
          <div className="buttons">
            <button onClick={() => handleReject(request.id)}>reject</button>
            <button onClick={() => handleAccept(request.id)}>accept</button>
          </div>
        </div>
      ))}
    </div>
  );
};

export default SettingsPage;
