import React from "react"

const Image = ({ image, className }) => {
  return (
    <div className={`image` + ` ${className}`}>
      {image ? (
        <div>
          <img src={`http://localhost:8080/${image}`} alt="" />
        </div>
      ) : null}
    </div>
  )
}

export default Image
