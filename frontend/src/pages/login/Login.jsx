import React, { useState, useContext, useRef } from 'react';
import { ToastContainer, toast } from 'react-toastify';
import { useNavigate } from 'react-router-dom';
import './Login.scss';
import { UserContext } from '../../components/utils/UserContext';
import { options } from '../../components/utils/toast/options';
import { myAxios, uploadAxios } from '../../api/axios';
import upload_icon_white from '../../icons/upload-icon-white.svg';
import { clearInput } from '../../components/utils/clearInput';
export const AuthForms = () => {
  return (
    <div className="main-container login">
      <div className="primal-container">
        <Form />
        <ToastContainer />
      </div>
    </div>
  );
};

const Form = React.forwardRef(({}) => {
  //states to set user credentials to login - email or username
  const [user, setUser] = useState('');
  const [password, setPassword] = useState('');

  //states to set user credentials to signup
  const [firsName, setFirsName] = useState('');
  const [lastName, setLastname] = useState('');
  const [email, setEmail] = useState('');
  const [username, setUsername] = useState('');
  const [about, setAbout] = useState('');
  const [age, setAge] = useState(0);
  const [gender, setGender] = useState('other');
  const [signUpPassword, setSignUpPassword] = useState('');
  const [confirm, setConfirm] = useState('');
  const [file, setFile] = useState(null);
  const fileType = file?.name.split('.')[1];
  //declare ref to get access to multipart data from
  const form = useRef(null);
  const navigate = useNavigate();
  const { setLogged } = useContext(UserContext);
  const { isSignUp, setIsSignUp } = useContext(UserContext);
  const handleTransition = () => {
    let f = document.getElementById('form-cont');
    f?.classList.toggle('transition');
  };
  const handleChange = () => {
    setIsSignUp(!isSignUp);
    handleTransition();
  };

  const handleLogin = async () => {
    await myAxios
      .post(
        '/signin',
        JSON.stringify({
          email: user,
          password,
        })
      )
      .then(() => {
        setLogged(true);

        navigate('/me');
      })

      .catch((error) => {
        setLogged(false);
        toast.error(error.response['data']['error_description'], options);
      });
  };

  const handleSignUp = async (e) => {
    e.preventDefault();
    if (confirm !== signUpPassword) {
      toast.error("Passwords don't match");
      return;
    }
    const data = new FormData(form.current);
    data.append('firstname', firsName);
    data.append('lastname', lastName);
    data.append('email', email);
    data.append('username', username);
    data.append('about', about);
    data.append('age', age);
    data.append('gender', gender);
    data.append('password', signUpPassword);
    data.append('image_type', 'avatar');
    data.append('file_type', fileType);

    await uploadAxios
      .post('/signup', data)
      .then(() => {
        handleChange();
        toast.success('Successfully signed up!', options);
        clearInput(
          ['input', 'textarea', 'select'],
          [
            setFirsName,
            setLastname,
            setUsername,
            setEmail,
            setAbout,
            setAge,
            setGender,
            setPassword,
            setConfirm,
          ]
        );
        setFile(null);
      })
      .catch((error) => {
        toast.error(error.response['data']['error_description'], options);
        return;
      });
  };

  return (
    <div className="forms-container" id="form-cont">
      {!isSignUp ? (
        <LoginForm
          handleChange={handleChange}
          setLogged={setLogged}
          setUsername={setUser}
          setPassword={setPassword}
          handleLogin={handleLogin}
        />
      ) : (
        <SignUpForm
          handleChange={handleChange}
          ref={form}
          setFirsName={setFirsName}
          setLastname={setLastname}
          setEmail={setEmail}
          setUsername={setUsername}
          setAbout={setAbout}
          setAge={setAge}
          setGender={setGender}
          setPassword={setSignUpPassword}
          setConfirm={setConfirm}
          handleSignUp={handleSignUp}
          file={file}
          setFile={setFile}
        />
      )}
    </div>
  );
});

const LoginForm = ({
  handleChange,
  setUsername,
  setPassword,
  setLogged,
  handleLogin,
}) => {
  return (
    <div className="form-container login-form">
      <FormTitle title={'Login'} />
      <FormInput
        type={'text'}
        placeholder={'email/username'}
        setter={setUsername}
      />
      <FormInput
        type={'password'}
        placeholder={'password'}
        setter={setPassword}
      />
      <FormButton
        text={'login'}
        setter={setLogged}
        actionHandler={handleLogin}
      />
      <ChangeForm
        text={"Don't have an account?"}
        link={'Sign up'}
        handleChange={handleChange}
      />
    </div>
  );
};

const SignUpForm = React.forwardRef(
  (
    {
      handleChange,
      setFirsName,
      setLastname,
      setEmail,
      setUsername,
      setAbout,
      setAge,
      setGender,
      setPassword,
      setConfirm,
      handleSignUp,
      file,
      setFile,
    },
    ref
  ) => {
    return (
      <form action="" ref={ref} encType="multipart/form-data">
        <div className="form-container signup-form">
          <FormTitle title={'Sign up'} />
          <div className="inputs-container">
            <div className="left-inputs">
              <FormInput
                type={'text'}
                placeholder={'firstname'}
                setter={setFirsName}
              />
              <FormInput
                type={'text'}
                placeholder={'lastname'}
                setter={setLastname}
              />
              <FormInput
                type={'text'}
                placeholder={'email'}
                setter={setEmail}
              />
              <FormInput
                type={'text'}
                placeholder={'username (optional)'}
                setter={setUsername}
              />
              <FormInput
                type={'text'}
                placeholder={'about you (optional)'}
                setter={setAbout}
              />
            </div>
            <div className="right-inputs">
              <FormInput type={'number'} placeholder={'age'} setter={setAge} />
              <GenderInput setter={setGender} />
              <AvatarUpload file={file} setFile={setFile} />
              <FormInput
                type={'password'}
                placeholder={'password'}
                setter={setPassword}
              />
              <FormInput
                type={'password'}
                placeholder={'confirm password'}
                setter={setConfirm}
              />
            </div>
          </div>
          <FormButton text={'sign up'} actionHandler={handleSignUp} />
          <ChangeForm
            text={'Already have an account?'}
            link={'Sign in'}
            handleChange={handleChange}
          />
        </div>
      </form>
    );
  }
);

//Form components

const FormTitle = ({ title }) => {
  return <h2 className="form-title">{title}</h2>;
};

const FormInput = ({ type, placeholder, setter }) => {
  return (
    <input
      type={type}
      placeholder={placeholder}
      onChange={(e) => setter(e.target.value)}
    />
  );
};

const AvatarUpload = ({ file, setFile }) => {
  let shortfileName;
  if (file?.name.length > 20) {
    shortfileName = file?.name.slice(0, 20) + '...';
  } else {
    shortfileName = file?.name;
  }
  return (
    <div className="avatar-upload">
      <label htmlFor="avtr">
        {file ? shortfileName : 'avatar (optional)'}{' '}
        <img src={upload_icon_white} alt="" />
      </label>
      <input
        type="file"
        id="avtr"
        onChange={(e) => setFile(e.target.files[0])}
        name="avatar"
      />
    </div>
  );
};

const FormButton = ({ actionHandler, text }) => {
  return (
    <button
      onClick={(e) => {
        actionHandler(e);
      }}
    >
      {text}
    </button>
  );
};

const GenderInput = ({ setter }) => {
  return (
    <select onChange={(e) => setter(e.target.value)}>
      <option value="" hidden>
        gender
      </option>
      <option value="Male">male</option>
      <option value="Female">female</option>
      <option value="Other">other</option>
    </select>
  );
};

const ChangeForm = ({ text, link, handleChange }) => {
  return (
    <p>
      {text}
      <span onClick={handleChange}>{link}</span>
    </p>
  );
};
