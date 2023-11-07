import styles from './Posts.module.scss';
import { Post } from '../../components/Post';
import { useGet } from '../../hooks/useGet';

const Feed = () => {
  const { data: posts, error, setHit } = useGet('/post/all');
  if (error) return <div>Error!</div>;
  if (posts?.length === 0) return <div>No posts</div>;
  return (
    <div className={styles.feedContainer}>
      {posts?.map((post) => (
        <Post
          key={post.id}
          id={post.id}
          title={post.title}
          content={post.content}
          author={post.user}
          likes={post.likes}
          dislikes={post.dislikes}
          hit={setHit}
          image={post.image}
          actionURL={`/post/${post.id}`}
        />
      ))}
    </div>
  );
};

export { Feed };
