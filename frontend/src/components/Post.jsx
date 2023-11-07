import { useNavigate } from "react-router-dom"
import styles from "../pages/posts/Posts.module.scss"
import like_icon from "../icons/like.svg"
import dislike_icon from "../icons/dislike.svg"
import { myAxios } from "../api/axios"
import Image from "./Image"
export const Post = ({
  id,
  title,
  content,
  author,
  likes,
  dislikes,
  hit,
  image,
  actionURL,
}) => {
  const navigate = useNavigate()
  const handlePostClick = () => {
    navigate(actionURL)
  }

  //action urls: /post/id/lie | /post/id/dislike | /group/id/post/id
  const handleLikeClick = async () => {
    await myAxios
      .post(
        actionURL + "/like",
        JSON.stringify({
          postlikedislike: "like",
        })
      )
      .then(() => hit(true))
  }
  const handleDislikeClick = async () => {
    await myAxios
      .post(
        actionURL + "/dislike",
        JSON.stringify({
          postlikedislike: "dislike",
        })
      )
      .then(() => hit(true))
      .catch((err) => console.log(err))
  }
  const handleAuthorClick = (id) => {
    navigate(`/user/${id}`)
  }

  return (
    <div className={styles.post_card}>
      <div onClick={() => handlePostClick()}>
        {image ? <Image image={image} className={"post-image"} /> : null}
      </div>

      <h1 className={styles.post_title} onClick={() => handlePostClick(id)}>
        {title}
      </h1>
      <div className={styles.post_content} onClick={() => handlePostClick(id)}>
        {content}
      </div>
      <div
        className={styles.post_author}
        onClick={() => handleAuthorClick(author?.id)}
      >
        By:{" "}
        {author?.username
          ? author?.username
          : author?.firstname + " " + author?.lastname}
      </div>
      <div className={styles.post_likes_dislikes}>
        <div className={styles.like_dislike}>
          <span onClick={handleLikeClick}>
            <img src={like_icon} alt="" />
          </span>
          {likes?.length > 0 ? likes?.length : 0}
        </div>
        <div className={styles.like_dislike}>
          <span onClick={handleDislikeClick}>
            <img src={dislike_icon} alt="" />
          </span>
          {dislikes?.length ? dislikes?.length : 0}
        </div>
      </div>
    </div>
  )
}
