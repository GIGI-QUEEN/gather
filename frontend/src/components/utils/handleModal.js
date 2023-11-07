export const handleModal = (element) => {
  let modal = document.getElementById("modal")
  let modal_form = document.getElementById(element)
  modal_form?.classList.toggle("create-visible")
  modal?.classList.toggle("visible")
}
