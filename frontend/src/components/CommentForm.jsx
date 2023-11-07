import React, { useRef, useContext, useState } from 'react';
import '../styles/comment-form.scss';
import { uploadAxios } from '../api/axios';
import { AttachAndEmojis } from './AttachAndEmojis';
import { checkFileType } from './utils/checkFileType';
import { WebSocketContext } from './utils/WebSocketContext';
import { UserContext } from './utils/UserContext';

const CommentForm = ({ groupId, postId, postAuthor, hit, url, image_type }) => {
  const { user } = useContext(UserContext);
  const { webSocket } = useContext(WebSocketContext);
  const [content, setContent] = useState('');
  const [file, setFile] = useState(null);
  const form = useRef(null);
  const fileType = file?.name?.split('.')[1];
  const handleSubmit = async (e) => {
    e.preventDefault();
    const data = new FormData(form.current);
    if (file || content) {
      if (file != null) {
        if (!checkFileType(file)) return;
      }
      data.append('content', content);
      data.append('file_type', fileType);
      data.append('image_type', image_type);
      await uploadAxios.post(url, data).then(() => {
        const commentOnGroupPostEvent = {
          event_type: 'ws_post_comment_event',
          group_id: groupId,
          post_id: postId,
          post_author: postAuthor,
          author_username: user.username ? user.username : user.firstname,
          content: content,
          created_date: Math.floor(Date.now() / 1000),
        };
        webSocket?.send(JSON.stringify(commentOnGroupPostEvent));
        hit(true);
        setContent('');
        setFile(null);
        form.current.reset(); // Reset the form
      });
    }
  };
  return (
    <form ref={form} encType="multipart/form-data" onSubmit={handleSubmit}>
      <div className="comment-form">
        <input
          value={content}
          placeholder="Comment..."
          onChange={(e) => setContent(e.target.value)}
        ></input>
        <AttachAndEmojis
          formName={image_type}
          setFile={setFile}
          setMsg={setContent}
        />
        <button type="submit">comment</button>
      </div>
    </form>
  );
};
export default CommentForm;
