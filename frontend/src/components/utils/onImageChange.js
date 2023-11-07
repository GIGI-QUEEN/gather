export const onImageChange = (e, setSelectedImage) => {
  if (e.target.files && e.target.files[0]) {
    setSelectedImage(URL.createObjectURL(e.target.files[0]))
  }
}
