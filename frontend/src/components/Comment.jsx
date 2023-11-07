import React from "react"
import "../styles/comment.scss"
import likeIcon from "../icons/like.svg"
import dislikeIcon from "../icons/dislike.svg"
import { myAxios } from "../api/axios"
import Image from "./Image"
import { useNavigate } from "react-router-dom"
import { usernameOrName } from "./utils/usernameOrName"
const Comment = ({
  id,
  content,
  user,
  likeCount,
  dislikeCount,
  post_id,
  hit,
  image,
  groupComment,
  groupid,
}) => {
  const nav = useNavigate()

  let url
  if (groupComment) {
    url = `/group/${groupid}/post/${post_id}/comment/${id}`
  } else {
    url = `/post/${id}/comment/${id}`
  }
  const handleLikeClick = async () => {
    await myAxios
      .post(
        `${url}/like`,
        JSON.stringify({
          commentlikedislike: "like",
        })
      )
      .then(() => hit(true))
  }
  const handleDislikeClick = async () => {
    await myAxios
      .post(
        `${url}/dislike`,
        JSON.stringify({
          commentlikedislike: "dislike",
        })
      )
      .then(() => hit(true))
  }

  return (
    <div className="comment-card">
      {image ? <Image image={image} className={"comment-image"} /> : null}

      <div className="comment-author">
        <h1 onClick={() => nav(`/user/${user?.id}`)}>{usernameOrName(user)}</h1>
      </div>
      <div className="comment-content">
        <p>{content}</p>
      </div>
      <div className="comment-likes_dislikes">
        <div className="like-dislike">
          <span onClick={handleLikeClick}>
            <img src={likeIcon} alt="like button" />
          </span>
          {likeCount}
        </div>
        <div className="like-dislike">
          <span onClick={handleDislikeClick}>
            <img src={dislikeIcon} alt="dislike button" />
          </span>
          {dislikeCount}
        </div>
      </div>
    </div>
  )
}

export default Comment
