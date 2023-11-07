//this func can either accept single input and setValue as arguments and array of inputs and setValues
//e.g. clearInput("textarea", setDescription) || clearInput(["input", "textarea"], [setTtile, setDescription, setEventDate])

export const clearInput = (tagName, valueSetter) => {
  if (Array.isArray(tagName)) {
    tagName.map((tag) => {
      let inp = document.getElementsByTagName(tag)
      for (let i = 0; i < inp.length; i++) {
        inp[i].value = ""
      }
    })
  } else {
    let inp = document.getElementsByTagName(tagName)
    for (let i = 0; i < inp.length; i++) {
      inp[i].value = ""
    }
  }

  if (Array.isArray(valueSetter)) {
    valueSetter.map((setValue) => setValue(""))
  } else {
    valueSetter("")
  }
}
