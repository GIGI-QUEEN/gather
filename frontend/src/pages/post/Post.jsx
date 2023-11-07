import { useParams } from 'react-router-dom';
import { Post } from '../../components/Post';
import Comment from '../../components/Comment';
import CommentForm from '../../components/CommentForm';
import { useGet } from '../../hooks/useGet';
import { PostCategories } from '../../components/PostCategories';

const PostPage = ({ fetch_url }) => {
  const { groupid, id } = useParams();
  const { data: post, setHit } = useGet(`${fetch_url}/${id}`);
  let actionURL;
  let imageType;
  let isGroupAction = false;
  if (groupid) {
    actionURL = `/group/${groupid}/post/${post?.post_id}`;
    imageType = 'group-post-comment-image';
    isGroupAction = true;
  } else {
    actionURL = `post/${post?.id}`;
    imageType = 'comment-image';
  }
  return (
    <div className="content-container">
      <Post
        id={post?.id}
        title={post?.title}
        content={post?.content}
        author={post?.user}
        likes={post?.likes}
        dislikes={post?.dislikes}
        hit={setHit}
        image={post?.image}
        actionURL={actionURL}
      />
      {post?.categories ? (
        <PostCategories categories={post?.categories} />
      ) : null}
      <CommentForm
        groupId={groupid}
        postId={post?.id}
        postAuthor={post?.user.id}
        hit={setHit}
        url={actionURL + '/comment'}
        image_type={imageType}
      />
      {post?.comments ? (
        <div className="comments-container">
          {post?.comments?.map((comment) => (
            <Comment
              key={comment.id}
              id={comment.id}
              content={comment.content}
              user={comment.user}
              likeCount={comment.likeCount}
              dislikeCount={comment.dislikeCount}
              post_id={comment.postId}
              hit={setHit}
              image={comment.image}
              groupComment={isGroupAction}
            />
          ))}
        </div>
      ) : null}
    </div>
  );
};

export { PostPage };
