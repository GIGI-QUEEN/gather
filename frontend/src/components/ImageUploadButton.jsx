import { onImageChange } from "./utils/onImageChange"
import { uploadIcon } from "../icons/upload-icon"

export const ImageUploadButton = ({ setFile, setSelectedImage, formName }) => {
  return (
    <div>
      <label htmlFor="image-inp">
        Image/GIF
        {uploadIcon}
      </label>
      <input
        type="file"
        name={formName}
        id="image-inp"
        onChange={(e) => {
          onImageChange(e, setSelectedImage)
          setFile(e.target.files[0])
        }}
      />
    </div>
  )
}
