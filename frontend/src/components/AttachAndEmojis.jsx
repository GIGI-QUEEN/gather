import React, { useState } from 'react';
import attachIcon from '../icons/attach-icon.svg';
import emojiIcon from '../icons/emoji-icon.svg';
import EmojiPicker from 'emoji-picker-react';
import { onImageChange } from './utils/onImageChange';

export const AttachAndEmojis = React.forwardRef(
  ({ formName, setFile, setMsg, setSelectedImage }) => {
    return (
      <div className="attach-and-emojis">
        <Emoji setMsg={setMsg} />
        <Attach
          formName={formName}
          setFile={setFile}
          setSelectedImage={setSelectedImage}
        />
      </div>
    );
  }
);

export const Emoji = ({ setMsg }) => {
  const [isPickerVisible, setPickerVisible] = useState(false);

  const emojiClick = (e) => {
    setMsg((prev) => prev + e.emoji);
  };

  return (
    <label htmlFor="" onClick={() => setPickerVisible(!isPickerVisible)}>
      <img src={emojiIcon} alt="" />
      <div className="emoji-picker">
        {isPickerVisible ? (
          <EmojiPicker
            lazyLoadEmojis={true}
            width={350}
            onEmojiClick={emojiClick}
          />
        ) : null}
      </div>
    </label>
  );
};

const Attach = React.forwardRef(({ formName, setFile, setSelectedImage }) => {
  return (
    <label htmlFor="img-attach">
      <img src={attachIcon} alt="" />

      <input
        type="file"
        name={formName}
        id="img-attach"
        onChange={(e) => {
          setFile(e.target.files[0]);
          onImageChange(e, setSelectedImage);
        }}
      />
    </label>
  );
});
