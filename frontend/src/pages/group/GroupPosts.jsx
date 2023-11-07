import { useState, useRef } from 'react';
import { Post } from '../../components/Post';
import { uploadAxios } from '../../api/axios';
import { clearInput } from '../../components/utils/clearInput';
import { toast } from 'react-toastify';
import { options } from '../../components/utils/toast/options';
import { handleModal } from '../../components/utils/handleModal';
import { ImageUploadButton } from '../../components/ImageUploadButton';
import closeIcon from '../../icons/close-icon.svg';

export const GroupPosts = ({ group, setHit }) => {
  if (!group?.group_posts) return null;

  return (
    <div className="group-posts">
      {group?.group_posts?.map((post) => (
        <Post
          key={post.post_id}
          id={post.post_id}
          title={post.title}
          content={post.content}
          author={post.user}
          hit={setHit}
          likes={post.likes}
          dislikes={post.dislikes}
          actionURL={`/group/${group.id}/post/${post.post_id}`}
          image={post.image}
        />
      ))}
    </div>
  );
};

export const CreateGroupPost = ({ group_id, setHit }) => {
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [file, setFile] = useState(null);
  const [selectedImage, setSelectedImage] = useState(null);
  const fileType = file?.name.split('.')[1];
  const form = useRef(null);
  const handleMultipart = async (e) => {
    e.preventDefault();
    const data = new FormData(form.current);
    data.append('title', title);
    data.append('content', content);
    data.append('image_type', 'group-post-image');
    data.append('file_type', fileType);

    await uploadAxios
      .post(`/group/${group_id}`, data)
      .then(() => {
        setHit(true);
        handleModal('create');
        clearInput('input', setTitle);
        clearInput('textarea', setContent);
        setSelectedImage(null);
        setFile(null);
      })
      .catch((error) => {
        toast.error(error.response['data']['error_description'], options);
        return;
      });
  };

  return (
    <div className="create-modal-container" id="create">
      <form encType="multipart/form-data" ref={form} onSubmit={handleMultipart}>
        <div className="title-and-close">
          <input
            type="text"
            placeholder="title"
            onChange={(e) => setTitle(e.target.value)}
          />
          <span
            onClick={() => {
              handleModal('create');
            }}
          >
            <img src={closeIcon} alt="" />
          </span>
        </div>
        <textarea
          name=""
          id=""
          placeholder="share your thoughts..."
          onChange={(e) => setContent(e.target.value)}
        ></textarea>
        <div className="upload-and-create">
          <ImageUploadButton
            setFile={setFile}
            setSelectedImage={setSelectedImage}
            formName={'group-post-image'}
          />
          {file ? <img src={selectedImage} alt="" /> : null}

          <button>create</button>
        </div>
      </form>
    </div>
  );
};
