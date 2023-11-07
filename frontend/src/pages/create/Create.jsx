import React, { useEffect, useRef, useState } from 'react';
import styles from './Create.module.scss';
import { useNavigate } from 'react-router-dom';
import { useContext } from 'react';
import { UserContext } from '../../components/utils/UserContext';
import { toast } from 'react-toastify';
import { options } from '../../components/utils/toast/options';
import { uploadAxios } from '../../api/axios';

import { checkFileType } from '../../components/utils/checkFileType';
import { useGet } from '../../hooks/useGet';
import { ImageUploadButton } from '../../components/ImageUploadButton';

const CreatePostPage = () => {
  return (
    <div>
      <CreateForm />
    </div>
  );
};

const CreateForm = () => {
  const navigate = useNavigate();
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [postCategories, setPostCategories] = useState([]);
  const me = useContext(UserContext);
  const [category, setCategory] = useState('Cars');
  const [privacy, setPrivacy] = useState('public');
  const [allowedFollowers, setAllowedFollowers] = useState([]);
  const form = useRef(null);
  const [selectedImage, setSelectedImage] = useState(null);
  const [file, setFile] = useState(null);

  useEffect(() => {
    if (me?.user?.privacy === 'private') {
      setPrivacy('private');
    }
  }, []);
  let fileType = file?.name.split('.')[1];
  const handleMultiPart = async (e, url) => {
    e.preventDefault(e);
    const data = new FormData(form.current);
    if (postCategories.length === 0) {
      postCategories.push('Misc');
    }
    if (file !== null) {
      if (!checkFileType(file)) return;
    }
    data.append('title', title);
    data.append('content', content);
    data.append('privacy', privacy);
    data.append('allowed', allowedFollowers);
    data.append('categories', postCategories);
    data.append('image_type', 'post-image');
    data.append('file_type', fileType);
    await uploadAxios
      .post(url, data)
      .then(() => navigate('/posts'))
      .catch((error) =>
        toast.error(error.response['data']['error_description'], options)
      );
  };

  return (
    <div className="content-container">
      <div className={styles.create_form_container}>
        <form
          encType="multipart/form-data"
          onSubmit={(e) => handleMultiPart(e, '/post/create')}
          ref={form}
        >
          <TitleInput setTitle={setTitle} />
          <CategoriesAndImageSection
            postCategories={postCategories}
            setPostCategories={setPostCategories}
            category={category}
            setCategory={setCategory}
            setFile={setFile}
            setSelectedImage={setSelectedImage}
          />
          <CategoriesList
            postCategories={postCategories}
            setPostCategories={setPostCategories}
          />
          <PostContentInput setContent={setContent} />
          <ShowChosenImage
            selectedImage={selectedImage}
            setSelectedImage={setSelectedImage}
          />
          <BottomSection
            privacy={privacy}
            setPrivacy={setPrivacy}
            me={me}
            allowedFollowers={allowedFollowers}
            setAllowedFollowers={setAllowedFollowers}
          />
        </form>
      </div>
    </div>
  );
};

const CategoriesAndImageSection = ({
  postCategories,
  setPostCategories,
  category,
  setCategory,
  setFile,
  setSelectedImage,
}) => {
  return (
    <div className={styles.categories_image_form}>
      <ChooseCategoriesList
        postCategories={postCategories}
        setPostCategories={setPostCategories}
        category={category}
        setCategory={setCategory}
      />
      <ImageUploadButton
        setFile={setFile}
        setSelectedImage={setSelectedImage}
        formName={'post-image'}
      />
    </div>
  );
};

const ChooseCategoriesList = ({
  postCategories,
  setPostCategories,
  category,
  setCategory,
}) => {
  const { data: categories } = useGet('/categories');
  const handleCategory = (e, new_category) => {
    e.preventDefault();
    let check = postCategories.find((cat) => cat === new_category);

    if (check === undefined) {
      setPostCategories([...postCategories, new_category]);
    }
  };

  return (
    <div>
      <select onChange={(e) => setCategory(e.target.value)}>
        {categories?.map((category) => (
          <option key={category.id} value={category.name}>
            {category.name}
          </option>
        ))}
      </select>
      <button onClick={(e) => handleCategory(e, category)}>+</button>
    </div>
  );
};

const TitleInput = ({ setTitle }) => {
  return (
    <div className={styles.title_input}>
      <input
        type="text"
        placeholder="Title"
        onChange={(e) => setTitle(e.target.value)}
      />
    </div>
  );
};

const CategoriesList = ({ postCategories, setPostCategories }) => {
  const handleRemove = (category_to_remove) => {
    let updated_categories = postCategories.filter(
      (category) => category !== category_to_remove
    );

    setPostCategories([...updated_categories]);
  };
  return (
    <div className={styles.categories}>
      {postCategories.map((category, index) => (
        <div
          className={styles.category}
          key={index}
          onClick={() => handleRemove(category)}
        >
          {category}
        </div>
      ))}
    </div>
  );
};

const PostContentInput = ({ setContent }) => {
  return (
    <div className={styles.post_content}>
      <textarea onChange={(e) => setContent(e.target.value)} />
    </div>
  );
};

const ShowChosenImage = ({ selectedImage, setSelectedImage }) => {
  return (
    <div className={styles.image_section}>
      <PostImage
        selectedImage={selectedImage}
        setSelectedImage={setSelectedImage}
      />
    </div>
  );
};

const BottomSection = ({
  privacy,
  setPrivacy,
  me,
  allowedFollowers,
  setAllowedFollowers,
}) => {
  return (
    <div className={styles.bottom_section}>
      <PostPrivacy
        privacy={privacy}
        setPrivacy={setPrivacy}
        me={me}
        allowedFollowers={allowedFollowers}
        setAllowedFollowers={setAllowedFollowers}
      />
      <CreateButton />
    </div>
  );
};

const PostPrivacy = ({
  privacy,
  setPrivacy,
  me,
  allowedFollowers,
  setAllowedFollowers,
}) => {
  const { user } = useContext(UserContext);
  return (
    <div className={styles.post_privacy}>
      <select onChange={(e) => setPrivacy(e.target.value)}>
        {user.privacy === 'private' ? null : (
          <option value="public">Public</option>
        )}{' '}
        {/* prevent private user from creating public posts */}
        <option value="private">Private</option>
        <option value="almost private">Almost private</option>
      </select>
      {privacy === 'almost private' ? (
        <ChooseFollowers
          followers={me?.user.followers}
          setAllowedFollowers={setAllowedFollowers}
          allowedFollowers={allowedFollowers}
        />
      ) : null}
    </div>
  );
};

const CreateButton = () => {
  return (
    <div className={styles.create_button}>
      <button type="submit">Create</button>
    </div>
  );
};

const PostImage = React.forwardRef(({ selectedImage }) => {
  if (selectedImage === null) return null;
  return (
    <div className={styles.post_image}>
      <img src={selectedImage} alt="" />
    </div>
  );
});

const ChooseFollowers = ({
  followers,
  allowedFollowers,
  setAllowedFollowers,
}) => {
  const handleFollowers = (e, f) => {
    e.preventDefault();
    let check = allowedFollowers.find((id) => id === f.id);
    if (check === undefined) {
      setAllowedFollowers([...allowedFollowers, f.id]);
    } else {
      let updated_followers = allowedFollowers.filter((id) => id !== f.id);
      setAllowedFollowers([...updated_followers]);
    }
  };

  return (
    <ul>
      {followers?.map((follower) => (
        <li key={follower.id}>
          <span className={styles.username}>{follower.username}</span>{' '}
          <button onClick={(e) => handleFollowers(e, follower)}>
            {allowedFollowers.find((id) => {
              return id === follower.id && allowedFollowers.length !== 0;
            })
              ? '-'
              : '+'}
          </button>{' '}
        </li>
      ))}
    </ul>
  );
};

export default CreatePostPage;
